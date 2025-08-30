# Rediska

Rediska is a Redis implementation built from scratch in Go. It supports core Redis commands and handles connections over TCP using the RESP.
This project is both an educational journey and a complicated challenge.

## Key Features

- RESP support
- RDB persistence
- Replication
- Multi type storage
- String data storage
- List data storage
- Stream data storage
- Sorted set data storage
- Pub/Sub messaging
- Transactions
- Geospatial commands
- Other commands

## How to run

1. Ensure you have `go (1.24)` installed locally
2. Run `./your_program.sh` to run your Redis server, you can pass such args like:

- `--host`
- `--port`
- `--dir`
- `--dbfilename`
- `--replicaof`

### To run master server:

```
./your_program.sh --port=6380
```

### To run replica server:

```
./your_program.sh --port=6381 --replicaof="127.0.0.1 6380"
```

### To run server with rdb file:

```
./your_program.sh --dir "." --dbfilename "dump.rdb" --port=6380
```

### To run client

I recommend you to install redis-cli and run your client like that:

```
redis-cli -h 127.0.0.1 -p 6380
```

## Features overview:

### RESP

RESP (redis serialization protocol) is based on 2 version of original RESP and consists of next types:

- Array (recursive type, can be nil)
- BulkString (can be nil)
- Integer
- SimpleString
- SimpleError

Some commands in redis can return multiple responses without wrapping it in one array (e.g. `SUBSCRIBE chan1 chan2`). In my case, i wrap it up in one final RESP Array.

### RDB persistence

Persistence ensures data is not lost.

In original Redis there are 2 main persistence strategies: `RDB (redis database)` file and `AOF (append only file)`. You can combine them both, use only one or not use at all. RDB file is a small binary file that encodes the whole redis storage. AOF strategy logs every write operation received by the server.

This project has only RDB file persistence. Once the server starts and if arguments are provided, the RDB file seeds the initial server storage and you can see decoded RDB file in server logs.

Limitations:

- no interval writing and saving the new state of server storage into RDB file
- sending RDB file allowed from master to replica when the handshake between them is in process
- only String storage type can be decoded and seeded into server storage

### Replication

Replication is a concept where you have some data and you want to clone it into another place. It increases up database durability. There are 2 main roles: `Master` and `Replica`. Master is only the one, who `propagates` (repeats) special commands to Replicas and they silently execute them on their side.

In original Redis, the logic is complicated. So, some features were omitted. First, Replica can't be a Master to other Replicas and in same time a Replica to another master. Second, no promote system (when servers shut down and we need to change out master). Third, to really be in sync, both Master and Server must track special data, called `replication offset`, there is only correct replication offset tracking for Replica.

To connect Replica to Master you need to pass valid `--replicaof` argument, that contains host and port of running Master server. Then, there will be a handshake with RDB file transfer from Master to Replica

List of commands, related to this extension:

- WAIT
- INFO
- REPLCONF
- PSYNC

### Multi type storage

The project storage is `divided into small storages`. Each small storage stores its `own data type` and related to it `methods`. In original Redis, there is only one main storage (map) that just contains different storage types. My variant leads to a small overhead for commands that need to scan the whole storage. But i didn't decide to do it like in original Redis, because it would cause a lot of edge cases checks, type checking and type assertion overhead. And also i don't block the whole storage, i block one type of storage at a time, so there can be parallel XADD and INCR commands for instance.

List of commands, related to this extension:

- KEYS
- DEL
- TYPE

### String type storage

Represents usual key-value map. Value type is always string, so you can set some numbers or words.

List of commands and features, related to this extension:

- GET
- SET (with or without expiration, both in MS and S)
- INCR
- expired keys regular cleanup (1/hour)

### List data storage

Represents key-value map where key is a string and value is a doubly linked list.

List of commands, related to this extension:

- RPUSH / LPUSH (with many values)
- RPOP / LPOP
- LLEN
- LRANGE
- BRPOP / BLPOP

### Stream data storage

Streams are needed to make Redis work like a message broker. In original Redis, Stream is radix trie. But it is a complex data structure and it is not needed considering the commands being developed. So, in my case, it's a map where `key` is a `stream name` and value has a `stream` type.
Stream represents a map where `key` is `stream id` (for each entry it is different, even in one stream) and value has `entry` type. Entry is a key-value map where key and value are both strings (here you store your payload). You can associate this structure with Kafka. For example, Kafka's topic == stream, Kafka's partition == entry.

List of commands, related to this extension:

- XADD
- XRANGE (with '+', '-' support)
- XREAD (with '$' support and blocking mode)

### Sorted set data storage

Sorted set is data structure that stores unique values (`members`) and also these values are sorted in ascending order by their `score`.
In original Redis, sorted set consists of hashmap and `skip list`. I also decided to create the same data structures as a challenge. The main reason to use skip list is for `fast insert, delete and search` operations.

Skip list consists of multiple linked lists, that are separated by levels. Each level presents original members of our sorted set with their scores. The smallest level consists of all elements in the sorted set. But as we know, to traverse a linked list we need O(N) time. The feauture of higher levels is that the higher level we choose, the less elements it has. The gap between elements increases with level and we can traverse it faster. Sometimes, we can miss element with that approach but we can go down and find it on a lower level. Because of ascending order, we always go right and down, and such manipulations generally lead us to O(logN) time complexity.

List of commands, related to this extension:

- ZADD (with many members)
- ZREM (with many members)
- ZRANK
- ZRANGE (with optional WITHSCORES argument)
- ZCARD
- ZSCORE

### Pub/Sub messaging

Pub/Sub is a pattern where there is one `source` that `sends messages` to all `subscribers` (something like Observer). As far as i know, in original Redis commands related to this extension can be propagated, but i didn't make it.

List of commands, related to this extension:

- PING (another behavior in subscribed mode)
- PUBLISH
- SUBSCRIBE (single/multiple channels)
- UNSUBCRIBE (single/multiple channels)

### Transactions

Transactions look like `deferred sequence (queue) of commands`, you write a lot of commands and than you can decide: execute them or discard. If you chose execution, they will be executed sequentially.

In original Redis, in transaction mode command is validated and only then is pushed to the queue. In mine, commands are pushed in queue, even if there is some invalid command, then it will just return corresponding RESP SimpleError. It is because my commands don't separate validation and execution and also errors can occur in execution, thus there need to be a lot of refactoring.

List of commands, related to this extension:

- MULTI
- EXEC
- DISCARD

### Geospatial commands

Geospatial commands help to store coordinates and to find a distance between 2 points on Earth. They are based on Sorted set data storage.

The most interesting, but the complicated parts to implement are encoding and decoding coordinate longitude and latitude to one value (it will be the score of member, referencing to the sorted set). But, the idea is simple, we take our longitude and latitude in binary formats and `interleave them with each other`. No less diffucult is the distance between 2 locations (points). It is a math thing, so it it cool, but i don't understand it quite well :) (but i've tried my best to make code readable)

List of commands, related to this extension:

- GEOADD (with many points in time)
- GEOPOS
- GEODIST
- GEOSEARCH (by radius, with 'km', 'ft', 'm', 'mi' support)

### Other commands

List of general commands:

- CONFIG GET
- ECHO
- PING

## Afterword

This project is currently one of the biggest projects i've ever made and the biggest one written in Go. I learned and practiced a lot while coding, both practicing new patterns, coding, designing project architecture overall.

A Redis server was built by completing all stages of the [Codecrafters Redis course](https://app.codecrafters.io/courses/redis/overview). Thank you so much guys for such an awesome challenge)
