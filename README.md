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
go get chat-grpc-go
```

## Generate Protocol Buffer

Install prerequisites

Ubuntu/Debian:
```
sudo apt update
sudo apt install -y protobuf-compiler
protoc --version

```

macOS (Homebrew):
```
brew install protobuf
protoc --version
```

Windows

Download protoc-<version>-win64.zip from protobuf releases\

Extract the zip (e.g., to C:\protoc).

Add C:\protoc\bin to your PATH in Environment Variables.

Open a new Command Prompt or PowerShell and check:

```
protoc --version
```

Then, install

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

And add Go bin to PATH:
```
export PATH="$PATH:$(go env GOPATH)/bin"
```

Use this command to generate protocol buffer manually :

```
protoc \
  --go_out=pb --go_opt=paths=source_relative \
  --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
  proto/chat/chat.proto
```
or you can do a simple command like `./generate.sh`, but make sure that `generate.sh` is an executable file. If you faced an error while executing the file try to run `chmod +x generate.sh` then run `./generate.sh` again.

## Run as the server

Simply run `go run main.go` and the server will be up in no time!

## Run as a client

To run as a client you need to run
```
go run client/example_client.go -sender rifqiakrm -room default
```

## Thanks to

This repository is inspired by [Dipesh Dulal's](https://github.com/dipeshdulal/grpc-samples) github repo. Please check him out!
