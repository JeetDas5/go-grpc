package main

import (
	"log"

	pb "github.com/jeetdas5/go-grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	port = ":8080"
)

func main() {
	conn, err := grpc.NewClient("localhost"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewGreetServiceClient(conn)

	names := &pb.NamesList{
		Names: []string{"Jeet", "Alice", "Bob"},
	}

	// clientSayHello(client);

	// callSayHelloServerStream(client,names)
	// callSayHelloClientStream(client, names)
	callSayHelloBiDirectionalStreaming(client, names)

}
