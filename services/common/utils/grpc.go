package utils

import (
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGrpcClient(addr string) *grpc.ClientConn {

	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%s", addr), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Fatalln("Failed to connect: %v", err)
	}

	return conn
}
