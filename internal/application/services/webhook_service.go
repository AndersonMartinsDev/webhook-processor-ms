package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time" // Importação necessária
	"webhook-processor-ms/internal/domain/message"
	"webhook-processor-ms/internal/domain/model"
)

// WebhookService orquestra o processamento do webhook.
type WebhookService struct {
	Publisher         message.MessagePublisher
	ConsumerPF        message.MessageConsumer
	SessionService    *SessionService
	AgentModelService *AgentModelService
}

// NewWebhookService cria uma nova instância de WebhookService.
func NewWebhookService(publisher message.MessagePublisher, consumer message.MessageConsumer, sessionService *SessionService,
	agentModelService *AgentModelService) *WebhookService {
	return &WebhookService{
		Publisher:         publisher,
		ConsumerPF:        consumer,
		SessionService:    sessionService,
		AgentModelService: agentModelService,
	}
}

// ProcessMessages inicia o loop de processamento com reconexão.
func (s *WebhookService) ProcessMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			slog.Info("Contexto cancelado, encerrando o loop de consumo.")
			return
		default:
			err := s.processMessagesPF(ctx)
			if err != nil {
				slog.Error("Falha na conexão ou consumo, tentando reconectar em 5 segundos...", "error", err)
				time.Sleep(5 * time.Second)
			} else {
				// Se a função retornar sem erro (o que é improvável em um loop de consumo),
				// a gente dá um pequeno tempo para evitar um loop muito rápido.
				time.Sleep(1 * time.Second)
			}
		}
	}
}

// processMessagesPF contém a lógica de consumo de um ciclo.
func (s *WebhookService) processMessagesPF(ctx context.Context) error {
	// Consome da fila de webhooks brutos (o nome da fila deve ser o mesmo usado no gateway-ms)
	msgs, err := s.ConsumerPF.Consume(ctx)
	if err != nil {
		slog.Error("Falha ao iniciar o consumo de webhooks", "error", err)
		return err // Retorna o erro para o loop principal.
	}
	slog.Info("Consumidor de webhook iniciado, esperando por mensagens...")

	// Loop para processar as mensagens recebidas
	for msgPayload := range msgs {
		var webhookEvent model.WWEBJSPayload
		err := json.Unmarshal(msgPayload, &webhookEvent)
		if err != nil {
			slog.Error("Falha ao desserializar o payload do webhook", "error", err)
			continue
		}

		if webhookEvent.Body != "" {
			From := strings.Replace(webhookEvent.From, "@c.us", "", 1)
			To := strings.Replace(webhookEvent.To, "@c.us", "", 1)

			var wg sync.WaitGroup
			wg.Add(2)

			var agentId uint64
			var session model.SessionModel

			go func() {
				defer wg.Done()
				agentId = s.AgentModelService.GetAgentId(To)
			}()

			go func() {
				defer wg.Done()
				session, _ = s.SessionService.GetHistory(From)
			}()

			wg.Wait() // Adicionado: espera ambas as goroutines terminarem

			if session.IDSession == "" {
				s.SessionService.SaveHistory(model.SessionModel{
					IDSession:  From,
					FontNumber: strings.Replace(webhookEvent.To, "@c.us", "", 1),
					AgentId:    agentId, // Agora 'agentId' terá o valor correto
				})
			}
			session.AgentId = agentId

			slog.Info("Webhook processado, criando payload simplificado...", "user_id", From)

			// Crie o payload simplificado
			simplifiedPayload := model.SimplifiedMessagePayload{
				SessionKey:  From,
				AgentID:     session.AgentId,
				FontNumber:  strings.Replace(webhookEvent.To, "@c.us", "", 1),
				MessageType: webhookEvent.Type,
				Message:     webhookEvent.Body,
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
	// Retorna um erro para o loop principal, indicando que o consumo terminou.
	// Isso sinaliza que o canal foi fechado e que uma reconexão é necessária.
	return fmt.Errorf("canal de consumo fechado")
}
