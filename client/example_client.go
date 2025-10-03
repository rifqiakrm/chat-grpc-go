package main

import (
	"bufio"
	"chat-grpc-go/pb/proto/chat"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
)

/*
Run as user 1 :
go run example_client.go -sender rifqiakrm -room room-1

Run as user 2 :
go run example_client.go -sender another_user -room room-1
*/

var roomId = flag.String("room", "room-1", "Room ID for chatting")
var senderId = flag.String("sender", "rifqiakrm", "Senders name")
var tcpServer = flag.String("server", "localhost:8080", "Tcp server")

func joinRoom(ctx context.Context, client chat.ChatServiceClient) {
	channel := chat.Room{
		RoomId: *roomId,
		UserId: *senderId,
	}
	stream, err := client.JoinRoom(ctx, &channel)
	if err != nil {
		log.Fatalf("JoinRoom throws: %v", err)
	}

	fmt.Printf("Joined room: %v \n", *roomId)

	ch := make(chan struct{})

	defer stream.CloseSend()

	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(ch)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive message from channel joining. \nErr: %v", err)
			}

			if *senderId != in.UserId {
				fmt.Printf("MESSAGE: (%v) -> %v \n", in.UserId, in.Message)
			}
		}
	}()

	<-ch
}

func sendMessage(ctx context.Context, client chat.ChatServiceClient, message string) {
	stream, err := client.SendMessage(ctx)
	if err != nil {
		log.Printf("Error while sending message: %v", err)
	}
	msg := chat.Chat{
		RoomId:     *roomId,
		UserId:     *senderId,
		Message:    message,
		Additional: nil,
		CreatedAt:  time.Now().String(),
		Type:       chat.Type_TEXT,
	}
	if err := stream.Send(&msg); err != nil {
		fmt.Println("err :", err)
	}

	ack, err := stream.CloseAndRecv()

	if err != nil {
		fmt.Println("err :", err)
	}

	fmt.Printf("Message sent: %v \n", ack)
}

func main() {
	flag.Parse()

	fmt.Println("--- CLIENT ---")
	fmt.Println("Connected to", *tcpServer)
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock(), grpc.WithInsecure())

	conn, err := grpc.Dial(*tcpServer, opts...)
	if err != nil {
		log.Fatalf("Fail to dial rpc: %v", err)
	}

	defer conn.Close()

	ctx := context.Background()
	client := chat.NewChatServiceClient(conn)

	go joinRoom(ctx, client)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		go sendMessage(ctx, client, scanner.Text())
	}

}
