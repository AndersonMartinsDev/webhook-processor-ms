package ms

import (
	"context"
	"webhook-processor-ms/proto"

	"google.golang.org/grpc"
)

type AgentModelClient interface {
	GetAgentModelByPhone(ctx context.Context, in *proto.AgentRequest, opts ...grpc.CallOption) (*proto.AgentResponse, error)
}
