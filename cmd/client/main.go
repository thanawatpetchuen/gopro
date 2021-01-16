package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	helloPb "github.com/thanawatpetchuen/gopro/generated/proto/hello"
	pingpongPb "github.com/thanawatpetchuen/gopro/generated/proto/pingpong"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func SendPing(client pingpongPb.PingPongClient, helloClient helloPb.HelloClient) (*pingpongPb.Pong, error) {
	// Timeout 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ping := pingpongPb.Ping{
		Id:      1,
		Message: "Ping",
	}
	pong, err := client.StartPing(ctx, &ping)
	statusCode := status.Code(err)
	if statusCode != codes.OK {
		return nil, err
	}

	// Hello Service
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	helloMassage, err := helloClient.Helloing(ctx, &helloPb.HelloParams{})
	statusCode = status.Code(err)
	if statusCode != codes.OK {
		return nil, err
	}

	// return
	pong.Message = helloMassage.GetMessage()
	return pong, err
}

func NewPingPongClient() pingpongPb.PingPongClient {
	opts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial("127.0.0.1:9001", opts...)
	if err != nil {
		panic(err)
	}

	client := pingpongPb.NewPingPongClient(conn)
	return client
}

func NewHelloClient() helloPb.HelloClient {
	opts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial("127.0.0.1:9000", opts...)
	if err != nil {
		panic(err)
	}

	client := helloPb.NewHelloClient(conn)
	return client
}

func main() {
	client := NewPingPongClient()
	helloClient := NewHelloClient()

	router := mux.NewRouter()
	router.Methods(http.MethodGet).Path("/ping").HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		msg, err := SendPing(client, helloClient)
		if err != nil {
			panic(err)
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(rw).Encode(msg); err != nil {
			log.Fatal(err)
		}
	})

	if err := http.ListenAndServe(":3000", router); err != nil {
		log.Fatal(err)
	}
}
