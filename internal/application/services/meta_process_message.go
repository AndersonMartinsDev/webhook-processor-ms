package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"webhook-processor-ms/internal/domain/message"
	"webhook-processor-ms/internal/domain/model"
)

type MetaProcessMessage struct {
	Consumer          message.MessageConsumer
	AgentModelService *AgentModelService
	SessionService    *SessionService
	Publisher         message.MessagePublisher
}

func NewMetaProcessMessage(consumer message.MessageConsumer, sessionService SessionService, agentModelService *AgentModelService, publisher message.MessagePublisher) *MetaProcessMessage {
	return &MetaProcessMessage{Consumer: consumer, SessionService: &sessionService, AgentModelService: agentModelService, Publisher: publisher}
}

// processMessagesPF contém a lógica de consumo de um ciclo.
func (s *MetaProcessMessage) ProcessMessages(ctx context.Context) error {
	// Consome da fila de webhooks brutos (o nome da fila deve ser o mesmo usado no gateway-ms)
	msgs, err := s.Consumer.Consume(ctx)
	if err != nil {
		slog.Error("Falha ao iniciar o consumo de webhooks", "error", err)
		return err // Retorna o erro para o loop principal.
	}
	slog.Info("Consumidor de webhook PJ iniciado, esperando por mensagens...")

	// Loop para processar as mensagens recebidas
	for msgPayload := range msgs {
		var webhookEvent model.MetaWebhookPayload
		if err := json.Unmarshal(msgPayload, &webhookEvent); err != nil {
			slog.Error("Falha ao desserializar o payload do webhook", "error", err)
			continue
		}
		slog.Info("Lendo Mensagem origem do META")

		// Garante que a estrutura do webhook é válida antes de tentar acessar os índices
		if len(webhookEvent.Entry) == 0 || len(webhookEvent.Entry[0].Changes) == 0 || len(webhookEvent.Entry[0].Changes[0].Value.Messages) == 0 {
			slog.Warn("Payload de webhook recebido, mas sem mensagens válidas para processar.")
			continue
		}

		message := webhookEvent.Entry[0].Changes[0].Value.Messages[0]
		from := message.From
		to := webhookEvent.Entry[0].Changes[0].Value.Metadata.DisplayPhoneNumber

		var messageContent string
		var processMessage bool = true

		// Usando switch para tratar diferentes tipos de mensagem
		switch message.Type {
		case "text":
			messageContent = message.Text.Body
		case "image":
			messageContent = message.Image.Caption
			if messageContent == "" {
				messageContent = message.Image.ID
			}
		case "audio":
			messageContent = message.Audio.ID
		case "interactive":
			if message.Interactive.Type == "button_reply" {
				messageContent = message.Interactive.ButtonReply.Title
			} else if message.Interactive.Type == "list_reply" {
				messageContent = message.Interactive.ListReply.Title
			}
		default:
			slog.Info("Tipo de mensagem não suportado recebido", "type", message.Type, "from", from)
			processMessage = false // Não processa tipos não suportados
		}

		// Se a mensagem for de um tipo que queremos processar
		if processMessage && messageContent != "" {
			// A lógica de goroutines para buscar agentId e session continua a mesma
			var wg sync.WaitGroup
			wg.Add(2)

			var agentId uint64
			var session model.SessionModel

			go func() {
				defer wg.Done()
				agentId, err = s.AgentModelService.GetAgentId(to)
				if err != nil {
					slog.Error("Error ao recuperar agent correto! %v", err)
				}
			}()

			go func() {
				defer wg.Done()
				session, err = s.SessionService.GetHistory(from)
				if err != nil {
					slog.Error("Error ao recuperar o historico! %v", err)
				}
			}()

			wg.Wait()

			if session.IDSession == "" {
				s.SessionService.SaveHistory(model.SessionModel{
					IDSession:  from,
					FontNumber: to,
					AgentId:    agentId,
				})
			}
			session.AgentId = agentId
			slog.Info("Webhook processado, criando payload simplificado...", "user_id", from, "type", message.Type)

			// Crie o payload simplificado
			simplifiedPayload := model.SimplifiedMessagePayload{
				OriginMessage: "META",
				SessionKey:    from,
				AgentID:       session.AgentId,
				FontNumber:    to,
				MessageType:   message.Type,
				Message:       messageContent,
				MediaURL:      messageContent,
			}

			payloadBytes, err := json.Marshal(simplifiedPayload)
			if err != nil {
				slog.Error("Falha ao serializar payload simplificado", "error", err)
				continue
			}

			if err := s.Publisher.Publish(ctx, payloadBytes); err != nil {
				slog.Error("Falha ao publicar mensagem para ai-integration-ms", "error", err)
			}

		}
	}
	return fmt.Errorf("canal de consumo fechado")
}
