package main

import (
	"context"
	"log"
	"time"

	pb "github.com/jeetdas5/go-grpc/proto"
)

func callSayHelloClientStream(client pb.GreetServiceClient, names *pb.NamesList) {
	log.Printf("Client Stream Started")

	stream, err := client.SayHelloClientStreaming(context.Background())
	if err != nil {
		log.Fatalf("Failed to start client streaming: %v", err)
	}

	for _, name := range names.Names {
		req := &pb.HelloRequest{
			Name: name,
		}
		if err := stream.Send(req); err != nil {
			log.Fatalf("Error sending: %v", err)
		}

		log.Printf("Sent the request with name: %s", name)

		time.Sleep(2 * time.Second)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving: %v", err)
	}
	log.Printf("Received response: %v", res)
}
