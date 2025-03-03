package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {

	config := &kafka.ConfigMap{
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
	c, err := kafka.NewConsumer(config)

	if err != nil {
		log.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}

	topic := "zookie-outbox"
	err = c.SubscribeTopics([]string{topic}, nil)
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
			ev, err := c.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// Errors are informational and automatically handled by the consumer

				continue
			}
			_, err = c.Commit()
			if err != nil {
				log.Printf("Error on commit: %v", err)
			}
			log.Printf("Consumed event from topic %s, partition %d at offset %s: key = %-10s value = %s\n",
				*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset, string(ev.Key), string(ev.Value))
		}
	}
	c.Close()
}
