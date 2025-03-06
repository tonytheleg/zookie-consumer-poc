# Zookie Consumer PoC

Super simple PoC of leveraging a Kafka Consumer for consuming messages from a topic containing Zookies to update in 
a database. The goal is to show a basic setup of how we can monitor a topic, consume messages of a specific 
prescribed format, and update the resources in a database AND ensure updates are ordered for consistency. 

### PoC Run Through

#### 1. Spin up the Kafka cluster and Consumers: `make poc-up`

> note, you'll likely see errors and things fail as there are dependencies on services. 
> This is because it takes a minute for zookeeper and postgres to setup and I didn't implement health checks to make this pretty.

#### 2. Clean logs for better viewing later (requires sudo permissions): `make clean-logs`

#### 3. Access the Kafka Node to initialize the test

```shell
# Access the container
docker exec -it zookie-consumer-poc-kafka-1 bash

# produce multiple messages to the topic in a loop
# This example uses resource_id as a message key, and the token as the message value. 
# Key and token are separated by Kafka using the defined separator '|' since it 
# should never appear in the payload (not standard for JSON). 
#
# look at this beaut.....
for i in `seq 1 25`; do for RESOURCE in my_resource_one another_resource_here third_resource_this_is; do \
  TOKEN=$(tr -dc A-Za-z0-9 </dev/urandom | head -c 13); \
  echo '{"resource_id":'"\"${RESOURCE}"'"}|{"consistency_token":'"\"${TOKEN}"'"}' | /opt/kafka/bin/kafka-console-producer.sh --bootstrap-server kafka:9093 --topic zookie-outbox --property parse.key=true --property key.separator='|'; \
  done; done
```

#### 4. Review consumer pod logs: `make consumer-logs`

A review of the logs will show a few things:

```shell
# example logs
===zookie-consumer-poc-zookie-consumer-1===
2025/03/06 18:27:21 Consumed event from topic zookie-outbox, partition 2 at offset 10: key = {"resource_id":"another_resource_here"} value = {"consistency_token":"liWW2ZQm0tllD"}
2025/03/06 18:27:23 Consumed event from topic zookie-outbox, partition 2 at offset 11: key = {"resource_id":"third_resource_this_is"} value = {"consistency_token":"iV1cLFiZ8WVdr"}
2025/03/06 18:27:26 Consumed event from topic zookie-outbox, partition 2 at offset 12: key = {"resource_id":"another_resource_here"} value = {"consistency_token":"BpokhbTwq1FWZ"}

===zookie-consumer-poc-zookie-consumer-2===

===zookie-consumer-poc-zookie-consumer-3===
2025/03/06 18:27:20 Consumed event from topic zookie-outbox, partition 0 at offset 5: key = {"resource_id":"my_resource_one"} value = {"consistency_token":"uzqZzOaqSYlAp"}
2025/03/06 18:27:24 Consumed event from topic zookie-outbox, partition 0 at offset 6: key = {"resource_id":"my_resource_one"} value = {"consistency_token":"pl6KMCdx6WAON"}
 ```

* When multiple partitions are configured for a topic (this PoC uses 3) and there are multiple Consumers in a Consumer group (3 again), the Kafka Broker will divide up the partitions amongst the Consumer group
* Each Consumer is currently assigned a partition in this PoC automatically, but what messages go to what partition are based on hashing of the message so it's not always a perfect split across each partition
  * The above shows that the hashing of the resource ID's caused Kafka to leverage partition 2 for both resources `another_resource_here` and `third_resource_this_is` and partition 0 for `my_resource_one`
* Since the app code is commiting messages on each message received the offset is always updated immediately which ensures if a consumer crashed it could pick back up again (seen in offset output)

### What happens if a Consumer crashes? 
* Since Key's are explicitly being set in the produce calls (`'{"resource_id":'"\"${RESOURCE}"'"}`), Kafka will ensure all messages for a single key are always delivered to the same Consumer
  * This ensures two consumers don't make changes to the same resource in parallel and ensures ordering of messages
  * You can see from the logs, all entries for a resource will always be consumed by the same consumer
* Killing and Restarting Consumers would show partitions get rebalanced and the partition would be resumed from the last offset accurately

### Testing Consumer Rebalancing

After going through Steps 1 through 3 above:

> Note: if you want to start fresh, `make poc-down` will tear down all the containers and you can start back at Step 1 above

#### 1. Print logs and take note of what `resource_ids` are assigned to what consumer and the offsets

```shell
$ make consumer-logs 
for i in zookie-consumer-poc-zookie-consumer-1 zookie-consumer-poc-zookie-consumer-2 zookie-consumer-poc-zookie-consumer-3; do echo "===${i}===" && docker logs ${i} && echo ""; done
===zookie-consumer-poc-zookie-consumer-1===
2025/03/06 19:06:04 Consumed event from topic zookie-outbox, partition 2 at offset 0: key = {"resource_id":"another_resource_here"} value = {"consistency_token":"jyAgIHSUpCtn4"}
2025/03/06 19:06:06 Consumed event from topic zookie-outbox, partition 2 at offset 1: key = {"resource_id":"third_resource_this_is"} value = {"consistency_token":"V5r5FTIV3ZitV"}
2025/03/06 19:06:08 Consumed event from topic zookie-outbox, partition 2 at offset 2: key = {"resource_id":"another_resource_here"} value = {"consistency_token":"26pjslRt473HM"}

===zookie-consumer-poc-zookie-consumer-2===
2025/03/06 19:06:03 Consumed event from topic zookie-outbox, partition 0 at offset 0: key = {"resource_id":"my_resource_one"} value = {"consistency_token":"ik8putTvInHwA"}
2025/03/06 19:06:07 Consumed event from topic zookie-outbox, partition 0 at offset 1: key = {"resource_id":"my_resource_one"} value = {"consistency_token":"YYhXvtYm8ygVu"}

===zookie-consumer-poc-zookie-consumer-3===

```
#### 2. Whack a Consumer: `make kill-a-consumer`

When it completes, all three pods will be back up and running, but you'll notice messages that were originally going to Consumer 1 are now going to another Consumer

```shell
$ make consumer-logs 
for i in zookie-consumer-poc-zookie-consumer-1 zookie-consumer-poc-zookie-consumer-2 zookie-consumer-poc-zookie-consumer-3; do echo "===${i}===" && docker logs ${i} && echo ""; done
===zookie-consumer-poc-zookie-consumer-1===
2025/03/06 19:06:04 Consumed event from topic zookie-outbox, partition 2 at offset 0: key = {"resource_id":"another_resource_here"} value = {"consistency_token":"jyAgIHSUpCtn4"}
2025/03/06 19:06:06 Consumed event from topic zookie-outbox, partition 2 at offset 1: key = {"resource_id":"third_resource_this_is"} value = {"consistency_token":"V5r5FTIV3ZitV"}
2025/03/06 19:06:08 Consumed event from topic zookie-outbox, partition 2 at offset 2: key = {"resource_id":"another_resource_here"} value = {"consistency_token":"26pjslRt473HM"}
2025/03/06 19:06:10 Consumed event from topic zookie-outbox, partition 2 at offset 3: key = {"resource_id":"third_resource_this_is"} value = {"consistency_token":"KO8tPq2qwKDLt"}
2025/03/06 19:06:12 Consumed event from topic zookie-outbox, partition 2 at offset 4: key = {"resource_id":"another_resource_here"} value = {"consistency_token":"SMcLeO5Zi5ges"}
########## CONSUMER 1 DIES AFTER OFFSET 4 #####################
2025/03/06 19:06:12 Caught signal terminated: terminating
2025/03/06 19:06:19 Waiting for messages
2025/03/06 19:06:19 Configuration settings:  map[auto.offset.reset:earliest bootstrap.servers:kafka:9093 enable.auto.commit:false group.id:zookie-consumer heartbeat.interval.ms:3000 max.poll.interval.ms:300000 session.timeout.ms:45000]
2025/03/06 19:06:22 Consumed event from topic zookie-outbox, partition 2 at offset 9: key = {"resource_id":"third_resource_this_is"} value = {"consistency_token":"BienXb6v3cw8p"}

===zookie-consumer-poc-zookie-consumer-2===
2025/03/06 19:06:03 Consumed event from topic zookie-outbox, partition 0 at offset 0: key = {"resource_id":"my_resource_one"} value = {"consistency_token":"ik8putTvInHwA"}
2025/03/06 19:06:07 Consumed event from topic zookie-outbox, partition 0 at offset 1: key = {"resource_id":"my_resource_one"} value = {"consistency_token":"YYhXvtYm8ygVu"}
2025/03/06 19:06:11 Consumed event from topic zookie-outbox, partition 0 at offset 2: key = {"resource_id":"my_resource_one"} value = {"consistency_token":"RIIvMg9JGaQN4"}
2025/03/06 19:06:15 Consumed event from topic zookie-outbox, partition 0 at offset 3: key = {"resource_id":"my_resource_one"} value = {"consistency_token":"j7Xlu7UqzFYYw"}
2025/03/06 19:06:19 Consumed event from topic zookie-outbox, partition 0 at offset 4: key = {"resource_id":"my_resource_one"} value = {"consistency_token":"yx6QhJZws3X0A"}

===zookie-consumer-poc-zookie-consumer-3===
######### CONSUMER 3, WHO PREVIOUSLY WAS NOT RECEVING MESSAGES PICKED UP WHERE CONSUMER 1 LEFTOFF ###################
2025/03/06 19:06:15 Consumed event from topic zookie-outbox, partition 2 at offset 5: key = {"resource_id":"third_resource_this_is"} value = {"consistency_token":"CLH2lvqJKnubP"}
2025/03/06 19:06:17 Consumed event from topic zookie-outbox, partition 2 at offset 6: key = {"resource_id":"another_resource_here"} value = {"consistency_token":"Lh21bdUFdslKt"}
2025/03/06 19:06:18 Consumed event from topic zookie-outbox, partition 2 at offset 7: key = {"resource_id":"third_resource_this_is"} value = {"consistency_token":"X9bJdSIduX7YP"}
2025/03/06 19:06:21 Consumed event from topic zookie-outbox, partition 2 at offset 8: key = {"resource_id":"another_resource_here"} value = {"consistency_token":"K32tJ3Uvi9MoG"}
```

#### Useful Commands

While exec'd into the Kafka container, here are some useful discovery commands

```shell
# list topics
/opt/kafka/bin/kafka-topics.sh --list --zookeeper zookeeper:2181

# describe topics
/opt/kafka/bin/kafka-topics.sh --describe --topic zookie-outbox --zookeeper zookeeper:2181


# describe a consumer group and see partition assignments
/opt/kafka/bin/kafka-consumer-groups.sh --bootstrap-server kafka:9093 --describe --group zookie-consumer

# describe a consumer group and see its members
/opt/kafka/bin/kafka-consumer-groups.sh --bootstrap-server kafka:9093 --describe --group zookie-consumer --members
```

### Confirmed Ordering

Also part of this PoC, I created a History table to show the updates of Zookies over time for each resource. After each message is consumed, this database is updated to:
* Set the token that was intialy set on the resource as `previous_token`
* Update the `current_token` with the token received in the Kafka message

As each message is consumed, this table shows a nice history of each resource and its zookie changes over time. Each entry for a specific resource you can trace the `previous_token` back to the last entry where it was set as the `current_token`.


1. Access the database: `psql postgres://postgres:tonyisawesome@localhost:5432/inventory-db`
2. Review the `resource_histories` table
```shell
inventory-db=# select * from resource_histories where resource_id='my_resource_one';
 id |          created_at           |          updated_at           | deleted_at |   resource_id   |      current_token       |      previous_token      
----+-------------------------------+-------------------------------+------------+-----------------+--------------------------+--------------------------
  1 | 2025-03-06 17:51:27.455649+00 | 2025-03-06 17:51:27.455649+00 |            | my_resource_one | myrandomconsistencytoken | 
 10 | 2025-03-06 19:06:03.386184+00 | 2025-03-06 19:06:03.386184+00 |            | my_resource_one | ik8putTvInHwA            | myrandomconsistencytoken
 13 | 2025-03-06 19:06:07.505977+00 | 2025-03-06 19:06:07.505977+00 |            | my_resource_one | YYhXvtYm8ygVu            | ik8putTvInHwA
 16 | 2025-03-06 19:06:11.610468+00 | 2025-03-06 19:06:11.610468+00 |            | my_resource_one | RIIvMg9JGaQN4            | YYhXvtYm8ygVu
 18 | 2025-03-06 19:06:15.703661+00 | 2025-03-06 19:06:15.703661+00 |            | my_resource_one | j7Xlu7UqzFYYw            | RIIvMg9JGaQN4
 22 | 2025-03-06 19:06:19.882244+00 | 2025-03-06 19:06:19.882244+00 |            | my_resource_one | yx6QhJZws3X0A            | j7Xlu7UqzFYYw
 25 | 2025-03-06 19:06:23.956993+00 | 2025-03-06 19:06:23.956993+00 |            | my_resource_one | LmlbXhZDjmJzZ            | yx6QhJZws3X0A
 28 | 2025-03-06 19:06:28.007393+00 | 2025-03-06 19:06:28.007393+00 |            | my_resource_one | GDZEl1MsghU2E            | LmlbXhZDjmJzZ
 31 | 2025-03-06 19:06:32.148935+00 | 2025-03-06 19:06:32.148935+00 |            | my_resource_one | 26RZtCTdE6Vp0            | GDZEl1MsghU2E
```

### Tear Down the PoC

To tear it all down and start fresh: `make poc-down`
