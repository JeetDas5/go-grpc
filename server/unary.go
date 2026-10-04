package main

import (
	"context"

	pb "github.com/jeetdas5/go-grpc/proto"
)

func (s *helloServer) SayHello(ctx context.Context, req *pb.NoParam) (*pb.HelloResponse, error) {
	res := &pb.HelloResponse{
		Message: "Hello World!",
	}
	return res, nil
}
