# Zookie Consumer POC

Super simple PoC of leveraging a Kafka Consumer for consuming messages from a topic containing Zookies to update in 
a database. The goal is to show a basic setup of how we can monitor a topic, consume messages of a specific 
prescribed format, and update the resources in a database AND ensure updates are ordered for consistency. 

### Build

`make build`

### Run Locally (Compose is better but..)

`make run`

### Spin up Kafka cluster with Consumers

`make poc-up`

> note, you'll likely see errors and things fail as there are dependencies on services. This is because it takes a minute for zookeeper and postgres to setup and I didnt implement healthchecks to make this pretty.

### Print all consumer logs (for looking at message consumption output)
make consumer-logs

### If you want clean logs to start with...
make clean-logs # will prompt for your password for sudo

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
# This example uses resource_id as a key, and the token as the value. Key and token are separated by Kafka using the defined
# separator '|' since it should never appear in the payload. Consumer logs should show that any messages with the same key will
# always go to the same partition and therefore handled by same consumer (ordered)
# to test multiple messages, you can change the resource_id or the token to show updates

# look at this beaut.....
for i in `seq 1 5`; do for RESOURCE in my_resource_one, another_resource_here, third_resource_this_is; do TOKEN=$(tr -dc A-Za-z0-9 </dev/urandom | head -c 13); echo '{"resource_id":'"\"${RESOURCE}"'"}|{"consistency_token":'"\"${TOKEN}"'"}' | /opt/kafka/bin/kafka-console-producer.sh --bootstrap-server kafka:9093 --topic zookie-outbox --property parse.key=true --property key.separator='|'; done; done
```
