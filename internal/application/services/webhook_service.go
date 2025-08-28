package services

import (
	"context"
	message "webhook-processor-ms/internal/domain"
)

// WebhookService orquestra o processamento do webhook.
type WebhookService struct {
	publisher message.MessagePublisher
}

// NewWebhookService cria uma nova instância de WebhookService.
func NewWebhookService(publisher message.MessagePublisher) *WebhookService {
	return &WebhookService{
		publisher: publisher,
	}
}

// ProcessWebhook recebe o request gRPC e publica na fila.
func (s *WebhookService) ProcessWebhook(ctx context.Context, payload string) error {
	return s.publisher.Publish(ctx, payload)
}
