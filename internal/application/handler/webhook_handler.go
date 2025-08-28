package handler

import (
	"context"
	"fmt"
	"log/slog"
	"webhook-processor-ms/internal/application/services"
	pb "webhook-processor-ms/proto"
)

// WebhookHandler lida com as requisições gRPC.
type WebhookHandler struct {
	pb.UnimplementedWebhookProcessorServiceServer
	webhookService *services.WebhookService // Depende da camada de serviço
}

// NewWebhookHandler cria um novo handler.
func NewWebhookHandler(s *services.WebhookService) *WebhookHandler {
	return &WebhookHandler{
		webhookService: s,
	}
}

// ProcessWebhook implementa a lógica do serviço gRPC.
func (h *WebhookHandler) ProcessWebhook(ctx context.Context, req *pb.ProcessWebhookRequest) (*pb.ProcessWebhookResponse, error) {
	slog.Info("Webhook gRPC recebido", "id", req.GetId())

	// Chama o serviço de aplicação para processar e publicar na fila
	err := h.webhookService.ProcessWebhook(ctx, req.GetPayload())
	if err != nil {
		slog.Error("Falha ao processar e publicar webhook", "error", err)
		return nil, fmt.Errorf("falha interna ao processar o webhook")
	}

	return &pb.ProcessWebhookResponse{
		Status:  "OK",
		Message: "Webhook processado e enviado para a fila.",
	}, nil
}
