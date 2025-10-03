package app

import (
	"chat-grpc-go/pb/proto/chat"
	"fmt"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

// NewChat returns a new server
func NewChat(tr opentracing.Tracer) *Chat {
	return &Chat{
		tracer:   tr,
		chatRoom: make(map[string]map[string]chan *chat.Chat),
	}
}

// Server implements the geo service
type Chat struct {
	chat.UnimplementedChatServiceServer
	tracer   opentracing.Tracer
	chatRoom map[string]map[string]chan *chat.Chat
}

func (s *Chat) Run(port int) error {
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(
			otgrpc.OpenTracingServerInterceptor(s.tracer),
		),
		grpc.StreamInterceptor(
			otgrpc.OpenTracingStreamServerInterceptor(s.tracer),
		),
	)

	chat.RegisterChatServiceServer(srv, s)

	reflection.Register(srv)

	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", viper.GetString("app.host"), port))
	if err != nil {
		return fmt.Errorf("failed to listen : %v", err)
	}
	log.Println("gRPC started!")
	log.Printf("server is running on %v:%v \n", viper.GetString("app.host"), viper.GetString("app.port"))
	if err := srv.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve rpc : %v", err)
	}

	return nil
}
