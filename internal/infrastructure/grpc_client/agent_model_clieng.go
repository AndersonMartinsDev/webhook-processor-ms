package grpc_client

import (
	"context"
	"webhook-processor-ms/proto"

	"google.golang.org/grpc"
)

type AgentModelClient struct {
	client proto.AIAgentServiceClient
}

// NewWebhookClient cria uma nova instância do cliente gRPC.
func NewAgentModelClient(conn *grpc.ClientConn) *AgentModelClient {
	return &AgentModelClient{
		client: proto.NewAIAgentServiceClient(conn),
	}
}

func (agent *AgentModelClient) GetAgentModelByPhone(ctx context.Context, in *proto.AgentRequest, opts ...grpc.CallOption) (*proto.AgentResponse, error) {
	return agent.client.GetAgentModelByPhone(ctx, in, opts...)
}
