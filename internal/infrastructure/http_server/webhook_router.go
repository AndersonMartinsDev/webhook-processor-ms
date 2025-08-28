package http_server

import "webhook-processor-ms/internal/application/services"

// UserHandler é o handler HTTP para as requisições de usuário.
type WebhookHandler struct {
	webhookService *services.WebhookService
}

// NewUserHandler cria uma nova instância de UserHandler.
func NewWebhookHandler(s *services.WebhookService) *WebhookHandler {
	return &WebhookHandler{
		webhookService: s,
	}
}
