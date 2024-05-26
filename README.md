### tcp-chat

A simple tcp chat service to run on a local network.

### Getting started

First make sure you have the [go compiler](https://go.dev/dl/) installed.

To build client run:
```bash
cd /path/to/tcp-chat
go build -o client cmd/client/main.go
```

To build server run:
```bash
cd /path/to/tcp-chat
go build -o server cmd/server/main.go
```

Then run the server with a given port as the command arg, like so:
```bash
./server 42069
```

Use the IP address printed by the server to connect with a client, like so:
```bash
./client server_ip:port
```
