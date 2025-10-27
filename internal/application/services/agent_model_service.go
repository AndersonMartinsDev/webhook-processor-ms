package services

import (
	"context"
	"time"
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

func (s *AgentModelService) GetAgentId(phoneNumber string) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	req := proto.AgentRequest{
		PhoneNumber: phoneNumber,
	}
	res, err := s.client.GetAgentModelByPhone(ctx, &req)
	return res.AgentId, err
}
