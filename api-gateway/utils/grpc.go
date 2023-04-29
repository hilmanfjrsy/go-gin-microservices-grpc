package utils

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GRPCClient(target string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Println("Can't connect to gRPC client", err)
		return nil, err
	}

	return conn, nil
}
