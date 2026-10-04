package main

import (
	"context"
	"io"
	"log"

	pb "github.com/jeetdas5/go-grpc/proto"
)

func callSayHelloServerStream(client pb.GreetServiceClient, names *pb.NamesList) {
	log.Printf("Streaming Started")

	stream, err := client.SayHelloServerStreaming(context.Background(), names)
	if err != nil {
		log.Fatalf("Failed to call SayHelloServerStreaming: %v", err)
	}

	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed to receive message: %v", err)
		}
		log.Println(message)
	}
	log.Printf("Streaming finished")
}
