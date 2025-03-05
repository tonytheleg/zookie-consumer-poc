package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	ResourceOne   = "my_resource_one"
	ResourceTwo   = "another_resource_here"
	ResourceThree = "third_resource_this_is"
)

var config = kafka.ConfigMap{
	// Settings are based on running Kafka using the provided compose files
	// See https://github.com/confluentinc/librdkafka/blob/master/CONFIGURATION.md for details
	"bootstrap.servers":     "kafka:9093",
	"group.id":              "zookie-consumer", // defines the consumer group name
	"session.timeout.ms":    "45000",           // ensures broker reblances consumer group
	"heartbeat.interval.ms": "3000",            // when consumer heartbeats are missing
	"max.poll.interval.ms":  "300000",          // Max time between calls for messages before considered dead
	"auto.offset.reset":     "earliest",        // Action to take when there is no initial offset in offset store or the desired offset is out of range
	"enable.auto.commit":    "false",           // disables auto commit of messages
}

type App struct {
	db       *gorm.DB
	consumer *kafka.Consumer
	topic    string
}

type Resource struct {
	gorm.Model
	ResourceID       string `gorm:"primaryKey" json:"resource_id"`
	ConsistencyToken string `json:"consistency_token"`
}

type ResourceHistory struct {
	gorm.Model
	ResourceID    string `json:"resource_id"`
	CurrentToken  string `json:"current_token"`
	PreviousToken string `json:"previous_token"`
}

func (a *App) NewApp() {
	var err error

	dsn := "host=invdatabase user=postgres password=tonyisawesome dbname=inventory-db port=5432 sslmode=disable TimeZone=America/New_York"
	a.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	err = a.db.AutoMigrate(&Resource{}, &ResourceHistory{})
	if err != nil {
		log.Fatal("failed to migrate database", err)
	}

	// Create initial values to test updates and ordering
	for _, resource := range []string{ResourceOne, ResourceTwo, ResourceThree} {
		a.db.Create(&Resource{ResourceID: resource, ConsistencyToken: "myrandomconsistencytoken"})
		a.db.Create(&ResourceHistory{ResourceID: resource, CurrentToken: "myrandomconsistencytoken"})
	}

	// Setup consumer and topic
	a.consumer, err = kafka.NewConsumer(&config)
	if err != nil {
		log.Fatal("Failed to create consumer: ", err)
	}
	a.topic = "zookie-outbox"
}

func (a *App) WriteDB(msg *kafka.Message) error {
	var resource Resource
	var resourceHistory ResourceHistory

	// a place to put values from message and not muck up our object from DB
	var msgInput struct {
		ResourceID       string `json:"resource_id"`
		ConsistencyToken string `json:"consistency_token"`
	}

	// Key and Message are separate JSON values -- must unmarshal each
	err := json.Unmarshal(msg.Key, &msgInput)
	if err != nil {
		return err
	}
	err = json.Unmarshal(msg.Value, &msgInput)
	if err != nil {
		return err
	}

	// Get resource from DB
	a.db.First(&resource, "resource_id = ?", msgInput.ResourceID)

	// update history for token changes -- this should confirm ordering
	resourceHistory.ResourceID = resource.ResourceID
	resourceHistory.PreviousToken = resource.ConsistencyToken
	resourceHistory.CurrentToken = msgInput.ConsistencyToken

	// update resource to new token
	resource.ConsistencyToken = msgInput.ConsistencyToken

	// write to db
	a.db.Save(&resource)
	a.db.Save(&resourceHistory)
	return nil
}

func main() {
	var app App
	app.NewApp()

	err := app.consumer.SubscribeTopics([]string{app.topic}, nil)
	if err != nil {
		log.Fatal("Failed to subscribe to topic: ", err)
	}
	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Process messages
	run := true
	log.Println("Waiting for messages")
	log.Println("Configuration settings: ", config)
	for run {
		select {
		case sig := <-sigchan:
			log.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			// Most examples seem to use Poll() and batch commit after a specific
			// number of messages. This example does so after every message
			ev, err := app.consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// Errors are informational and automatically handled by the consumer
				continue
			}
			err = app.WriteDB(ev)
			if err != nil {
				log.Printf("Failed to write to db: %v", err)
			}
			_, err = app.consumer.Commit()
			if err != nil {
				log.Printf("Error on commit: %v", err)
			}
			log.Printf("Consumed event from topic %s, partition %d at offset %s: key = %-10s value = %s\n",
				*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset, string(ev.Key), string(ev.Value))
		}
	}
	app.consumer.Close()
}
