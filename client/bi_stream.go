package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/jeetdas5/go-grpc/proto"
)

func callSayHelloBiDirectionalStreaming(client pb.GreetServiceClient, names *pb.NamesList) {
	log.Printf("SayHelloBiDirectionalStreaming start")

	stream, err := client.SayHelloBiDirectionalStreaming(context.Background())

	if err != nil {
		log.Fatalf("Failed to call SayHelloBiDirectionalStreaming: %v", err)
	}

	waitc := make(chan struct{})

	go func() {
		for {
			message, err := stream.Recv()

			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Failed to receive response: %v", err)
			}
			log.Printf("Received response: %v", message)
		}
		close(waitc)
	}()

	for _, name := range names.Names {
		req := &pb.HelloRequest{
			Name: name,
		}

		if err := stream.Send(req); err != nil {
			log.Fatalf("Failed to send request: %v", err)
		}

		time.Sleep(2 * time.Second)
	}
	stream.CloseSend()
	<-waitc
	log.Println("SayHelloBiDirectionalStreaming done")
}
