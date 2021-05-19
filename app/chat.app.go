package app

import (
	"fmt"
	"github.com/rifqiakrm/chat-grpc-go/pb/chat"
	"io"
)

//Server streaming
func (s *Chat) JoinRoom(req *chat.Room, stream chat.ChatService_JoinRoomServer) error {
	msgChannel := make(chan *chat.Chat)
	if s.chatRoom[req.GetRoomId()] == nil {
		s.chatRoom[req.GetRoomId()] = make(map[string]chan *chat.Chat)
	}
	s.chatRoom[req.GetRoomId()][req.GetUserId()] = msgChannel

	for {
		select {
		case <-stream.Context().Done():
			close(s.chatRoom[req.GetRoomId()][req.GetUserId()])
			delete(s.chatRoom[req.GetRoomId()], req.GetUserId())
			return nil
		case msg := <-msgChannel:
			fmt.Println("Channels : ", s.chatRoom[req.GetRoomId()])
			fmt.Printf("GO ROUTINE (got message): %v \n", msg)
			if err := stream.Send(msg); err != nil {
				fmt.Printf("failed to send message): %v \n", err)
				return err
			}
		}
	}
}

//Client Streaming
func (s *Chat) SendMessage(stream chat.ChatService_SendMessageServer) error {
	msg, err := stream.Recv()

	if err == io.EOF {
		return nil
	}

	if err != nil {
		return err
	}

	ack := chat.ChatAck{Status: "SENT"}
	if err := stream.SendAndClose(&ack); err != nil {
		return err
	}

	go func() {
		streams := s.chatRoom[msg.GetRoomId()]
		fmt.Println("client :", s.chatRoom[msg.GetRoomId()])
		for _, ch := range streams {
			ch <- msg
		}
	}()

	return nil
}
