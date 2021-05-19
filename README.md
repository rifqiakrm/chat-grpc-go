# Chat GRPC Service

Chat GRPC service with client and server side streaming. This code is actually being used in real life project. Tested OK with Android and iOS app.

It comes pre-configured with :

1. GRPC (google.golang.org/grpc)
2. GRPC OpenTracing (github.com/grpc-ecosystem/grpc-opentracing)
3. OpenTracing for Go (github.com/opentracing/opentracing-go)
4. Viper (github.com/spf13/viper)
5. Jaeger (github.com/uber/jaeger-client-go)


## Setup

Use this command to install the blueprint

```bash
go get github.com/rifqiakrm/chat-grpc-go
```

## Generate Protocol Buffer

Use this command to generate protocol buffer manually :

```
protoc -I $GOPATH/src --go_out=$GOPATH/src $GOPATH/src/github.com/rifqiakrm/{project_name}/proto/{proto_dir}/{your_proto}.proto
protoc -I $GOPATH/src --go-grpc_out=$GOPATH/src $GOPATH/src/github.com/rifqiakrm/{project_name}/proto/{proto_dir}/{your_proto}.proto
```
or you can do a simple command like `./generate.sh`, but make sure that `generate.sh` is an executable file. If you faced an error while executing the file try to run `chmod +x generate.sh` then run `./generate.sh` again.

## Run as the server

Simply run `go run main.go` and the server will be up in no time!

## Run as a client

To run as a client you need to run
```
go run client/example.go -name rifqiakrm -channel default
```

In the server, 
- the request is received `chan` is constructed and put into the map for later use when multiple users are connected in a channel.

- On another user connecting to the channel, the list of channel is appended with another `chan` variable.

- And then the function `JoinChannel` is waiting for message to be present in `chan`

    ```go
    case msg := <-msgChannel:
        fmt.Printf("GO ROUTINE (got message): %v \n", msg)
        msgStream.Send(msg)
    ```

- When one client sends the message, `SendMessage` function receives the message and puts the message in `chan` object present in the map looping over.

    ```go
    streams := s.channel[msg.Channel.Name]
    for _, msgChan := range streams {
        msgChan <- msg
    }
    ```

In this way, server receives and routes the message using `grpc`.

## Thanks to

This repository is inspired by [Dipesh Dulal's](https://github.com/dipeshdulal/grpc-samples) github repo. Please check him out!
