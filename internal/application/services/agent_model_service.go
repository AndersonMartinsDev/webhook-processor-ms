package services

import (
	"context"
	"log/slog"
	"webhook-processor-ms/internal/domain/ms"
	"webhook-processor-ms/proto"
)

type AgentModelService struct {
	client ms.AgentModelClient
}

func NewAgentModelService(client ms.AgentModelClient) *AgentModelService {
	return &AgentModelService{
		client: client,
	}
}

func (s *AgentModelService) GetAgentId(phoneNumber string) uint64 {
	ctx := context.Background()
	req := proto.AgentRequest{
		PhoneNumber: phoneNumber,
	}
	res, err := s.client.GetAgentModelByPhone(ctx, &req)
	if err != nil {
		slog.Error("Error para recuperar identificador de agente")
		return uint64(0)
	}
	return res.AgentId
}
