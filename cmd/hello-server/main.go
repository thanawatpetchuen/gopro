package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	pb "github.com/thanawatpetchuen/gopro/generated/hello/proto"
)

// Implement pb.HelloServer
type HelloServerImpl struct {
}

func (s *HelloServerImpl) Helloing(ctx context.Context, params *pb.HelloParams) (*pb.HelloMessage, error) {
	fmt.Println("Hello Received")

	resp := pb.HelloMessage{
		Id:      1,
		Message: "HELLO WORLD",
	}

	return &resp, nil
}

func StartHelloServer() {
	server := HelloServerImpl{}

	lis, err := net.Listen("tcp", "localhost:9000")

	grpcServer := grpc.NewServer()
	pb.RegisterHelloServer(grpcServer, &server)

	// Start grpcServer
	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	exit := make(chan os.Signal)
	signal.Notify(exit, os.Interrupt)

	fmt.Println("Hello started")

	<-exit
	fmt.Println("Stop Hello server")
}

func main() {
	StartHelloServer()
}
