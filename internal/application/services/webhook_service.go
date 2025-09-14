package services

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"webhook-processor-ms/internal/domain/message"
	"webhook-processor-ms/internal/domain/model"
)

// WebhookService orquestra o processamento do webhook.
type WebhookService struct {
	Publisher         message.MessagePublisher
	Consumer          message.MessageConsumer
	SessionService    *SessionService
	AgentModelService *AgentModelService
}

// NewWebhookService cria uma nova instância de WebhookService.
func NewWebhookService(publisher message.MessagePublisher, consumer message.MessageConsumer, sessionService *SessionService,
	agentModelService *AgentModelService) *WebhookService {
	return &WebhookService{
		Publisher:         publisher,
		Consumer:          consumer,
		SessionService:    sessionService,
		AgentModelService: agentModelService,
	}
}

// ProcessMessages recebe o request gRPC e publica na fila.
func (s *WebhookService) ProcessMessages(ctx context.Context) {
	// Consome da fila de webhooks brutos (o nome da fila deve ser o mesmo usado no gateway-ms)
	msgs, err := s.Consumer.Consume(ctx)
	if err != nil {
		slog.Error("Falha ao iniciar o consumo de webhooks", "error", err)
		return
	}
	slog.Info("Consumidor de webhook iniciado, esperando por mensagens...")

	// Loop para processar as mensagens recebidas
	for msgPayload := range msgs {
		var webhookEvent model.MetaWebhookPayload
		err := json.Unmarshal(msgPayload, &webhookEvent)
		if err != nil {
			slog.Error("Falha ao desserializar o payload do webhook", "error", err)
			continue
		}

		if len(webhookEvent.Entry) > 0 && len(webhookEvent.Entry[0].Changes) > 0 && len(webhookEvent.Entry[0].Changes[0].Value.Messages) > 0 {
			messageData := webhookEvent.Entry[0].Changes[0].Value.Messages[0]
			metadata := webhookEvent.Entry[0].Changes[0].Value.Metadata
			userID := messageData.From

			var wg sync.WaitGroup
			wg.Add(2)

			var agentId uint64
			var session model.SessionModel

			go func() {
				defer wg.Done()
				agentId = s.AgentModelService.GetAgentId(userID)
			}()

			go func() {
				defer wg.Done()
				session, _ = s.SessionService.GetHistory(userID)
			}()

			wg.Wait() // Adicionado: espera ambas as goroutines terminarem

			if session.IDSession == "" {
				s.SessionService.SaveHistory(model.SessionModel{
					IDSession:  userID,
					FontNumber: metadata.DisplayPhoneNumber,
					AgentId:    agentId, // Agora 'agentId' terá o valor correto
				})
			}
			session.AgentId = agentId

			slog.Info("Webhook processado, criando payload simplificado...", "user_id", userID)

			// Crie o payload simplificado
			simplifiedPayload := model.SimplifiedMessagePayload{
				SessionKey:  userID,
				AgentID:     session.AgentId,
				FontNumber:  metadata.DisplayPhoneNumber,
				MessageType: messageData.Type,
				Message:     messageData.Text.Body,
			}

			// Serializar e publicar na próxima fila
			payloadBytes, err := json.Marshal(simplifiedPayload)
			if err != nil {
				slog.Error("Falha ao serializar payload simplificado", "error", err)
				continue
			}

			// Publicar na fila de entrada do ai-integration-ms
			// O nome da fila deve ser o que o ai-integration-ms está consumindo
			if err := s.Publisher.Publish(ctx, payloadBytes); err != nil {
				slog.Error("Falha ao publicar mensagem para ai-integration-ms", "error", err)
			}
		}
	}
	slog.Info("Consumidor de webhook encerrado.")
}
