package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	pb "github.com/thanawatpetchuen/gopro/generated/pingpong/proto"
)

// Implement pb.PingPongServer
type PingPongServerImpl struct {
}

func (s *PingPongServerImpl) StartPing(ctx context.Context, ping *pb.Ping) (*pb.Pong, error) {
	fmt.Println("Ping Received")

	resp := pb.Pong{
		Id:      ping.Id,
		Message: ping.Message,
		User: &pb.User{
			Id:   2,
			Name: "Thanawat",
		},
	}

	return &resp, nil
}

func StartPingPongServer() {
	server := PingPongServerImpl{}

	lis, err := net.Listen("tcp", "localhost:9001")

	grpcServer := grpc.NewServer()
	pb.RegisterPingPongServer(grpcServer, &server)

	// Start grpcServer
	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	exit := make(chan os.Signal)
	signal.Notify(exit, os.Interrupt)

	fmt.Println("pingpong started")

	<-exit
	fmt.Println("Stop pingpong server")
}

func main() {
	StartPingPongServer()
}
