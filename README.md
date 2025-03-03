# Zookie Consumer POC

Super simple Kafka Consumer setup that allows for testing consuming messages from a topic with Go. Settings are easily updated in Go or compose files to nail down configurations needed to ensure the right consistency of message processing

### Build

`make build`

### Run Locally (Compose is better but..)

`make run`

### Spin up Kafka cluster with Consumers

`make poc-up`

> note, the kafka server always fails first run, it will restart. You may see errors in consumers at first. This is because it takes a minute for zookeeper to setup and I didnt implement healthchecks to make this pretty.

### Print all consumer logs (for looking at message consumption output)
make consumer-logs

### Spin down

`make poc-down`

## Useful Commands

Here are some useful Kafka commands you can run from inside the Kafka container (`docker exec -it zookie-consumer-poc-kafka-1 bash`)

```shell
# list topics
/opt/kafka/bin/kafka-topics.sh --list --zookeeper zookeeper:2181

# create topics (the below topic already exists)
/opt/kafka/bin/kafka-topics.sh --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions 3 --topic zookie-outbox

# describe topics
/opt/kafka/bin/kafka-topics.sh --describe --topic zookie-outbox --zookeeper zookeeper:2181

# add partitions to a topic (NOTE, you can only increase, not decrease)
/opt/kafka/bin/kafka-topics.sh --alter --zookeeper zookeeper:2181 --partitions 4 --topic zookie-outbox

# delete topic
/opt/kafka/bin/kafka-topics.sh --delete --zookeeper zookeeper:2181 --topic zookie-outbox

# describe a consumer group and see parition assignments
/opt/kafka/bin/kafka-consumer-groups.sh --bootstrap-server kafka:9093 --describe --group zookie-consumer

# describe a consumer group and see its members
/opt/kafka/bin/kafka-consumer-groups.sh --bootstrap-server kafka:9093 --describe --group zookie-consumer --members

# produce a message to a topic
# This example uses resource_id as a key, and the token as the value. Key and token are separeated by Kafka using the defined
# separator '|' since it should never appear in the payload. Consumer logs should show that any messages with the same key will
# always go to the same parition and therefore handled by same consumer (ordered)
# to test multiple messages, you can change the resource_id or the token to show updates
echo "'{\"resource_id\":\"my_cluster\"}'|'{\"continuationToken\":\"1a2b3c4d=\"}'" | /opt/kafka/bin/kafka-console-producer.sh --bootstrap-server kafka:9093 --topic zookie-outbox --property parse.key=true --property key.separator='|'
```
