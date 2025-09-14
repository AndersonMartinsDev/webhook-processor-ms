package grpc_client

import (
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GrcpConnection(service_url string) *grpc.ClientConn {
	conn, err := grpc.NewClient(service_url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error(fmt.Sprintf("Não foi possível conectar ao servidor Webhook Processor: %v", err))
		panic("Não foi possível conectar ao servidor Webhook Processor")
	}
	slog.Info(fmt.Sprintf("Connection Sucessfull with MS %s", service_url))
	return conn
}
