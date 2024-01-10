# kv-raft
A distributed key-value storage powered by Raft protocol developed in golang

## How to use

1. Run Node 1: Assume node 1 with NodeID=node1 and run it as single Raft node
```
go run main.go --node-id node1 --raft-port 2222 --http-port 8222
```

2. Run Node 1: Assume node 1 with NodeID=node1 and run it as single Raft node
```
go run main.go --node-id node1 --raft-port 2223 --http-port 8223
```

3. Join Node 2 as follower node to Node 1 as leader
```
curl 'localhost:8222/join?followerAddr=localhost:2223&followerId=node2'
```

4. Store a simple Key/Value by Node1
```
curl -X POST 'localhost:8222/store' -d '{"key": "x", "value": "323"}' -H 'content-type: application/json'
```

5. Restore the Value of the key from Node1 and Node2
```
curl 'localhost:8222/store?key=x'
```