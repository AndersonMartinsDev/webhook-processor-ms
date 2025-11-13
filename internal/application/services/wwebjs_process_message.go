package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"webhook-processor-ms/internal/domain/message"
	"webhook-processor-ms/internal/domain/model"
)

type WWEBJSProcessMessage struct {
	Consumer          message.MessageConsumer
	AgentModelService *AgentModelService
	SessionService    *SessionService
	Publisher         message.MessagePublisher
}

func NewWWEBJSProcessMessage(consumer message.MessageConsumer, sessionService SessionService, agentModelService *AgentModelService, publisher message.MessagePublisher) *WWEBJSProcessMessage {
	return &WWEBJSProcessMessage{Consumer: consumer, SessionService: &sessionService, AgentModelService: agentModelService, Publisher: publisher}
}

// processMessagesPF contém a lógica de consumo de um ciclo.
func (s *WWEBJSProcessMessage) ProcessMessages(ctx context.Context) error {
	msgs, err := s.Consumer.Consume(ctx)
	if err != nil {
		slog.Error("Falha ao iniciar o consumo de webhooks", "error", err)
		return err
	}
	slog.Info("Consumidor de webhook WEBJS iniciado, esperando por mensagens...")

	// Loop para processar as mensagens recebidas
	for msgPayload := range msgs {
		var webhookEvent model.WWEBJSPayload
		err := json.Unmarshal(msgPayload, &webhookEvent)
		if err != nil {
			slog.Error("Falha ao desserializar o payload do webhook", "error", err)
			continue
		}
		slog.Info("Lendo Mensagem origem WEBJS")

		// --- Extração de Números ---
		From := strings.Replace(webhookEvent.From, "@c.us", "", 1)
		To := strings.Replace(webhookEvent.To, "@c.us", "", 1)

		// --- Lógica de Extração de Conteúdo Refatorada ---
		var messageContent string
		var processMessage bool = true

		// O whatsapp-web.js usa 'type' para classificar a mensagem.
		switch webhookEvent.Type {
		case "chat": // Mensagens de texto simples
			messageContent = webhookEvent.Body

		case "ptt", "audio":
			webhookEvent.Type = "audio"
			messageContent = webhookEvent.Data.DeprecatedMms3Url

		case "image", "sticker":
			if webhookEvent.Body != "" {
				messageContent = webhookEvent.Body // Captura legenda (caption)
			} else {
				messageContent = webhookEvent.Data.DeprecatedMms3Url
			}

		case "list_response":
		case "button_response":
		case "notification_template":
		default:
			// Mensagens que não são de interesse ou não são suportadas
			slog.Info("Tipo de mensagem não suportado recebido", "type", webhookEvent.Type, "from", From)
			processMessage = false
		}

		if processMessage && messageContent != "" {
			var wg sync.WaitGroup
			wg.Add(2)

			var agentId uint64
			var session model.SessionModel

			go func() {
				defer wg.Done()
				agentId, err = s.AgentModelService.GetAgentId(To)
				if err != nil {
					slog.Error("Error ao recuperar agent correto! ", err)
				}
			}()

			go func() {
				defer wg.Done()
				session, err = s.SessionService.GetHistory(From)
				if err != nil {
					slog.Error("Error ao recuperar o historico! ", err)
				}
			}()

			wg.Wait()

			if session.IDSession == "" {
				s.SessionService.SaveHistory(model.SessionModel{
					IDSession:  From,
					FontNumber: To,
					AgentId:    agentId,
				})
			}
			session.AgentId = agentId

			slog.Info("Webhook processado, criando payload simplificado...", "user_id", From, "type", webhookEvent.Type)

			// --- Criação do Payload Simplificado Refatorada ---
			simplifiedPayload := model.SimplifiedMessagePayload{
				OriginMessage: "WWEBJS",
				SessionKey:    From,
				AgentID:       session.AgentId,
				FontNumber:    To,
				MessageType:   webhookEvent.Type,
				Message:       messageContent,
				MessageID:     webhookEvent.Data.ID.Serialized,
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
	slog.Info("Consumidor de webhook encerrado.")
	return fmt.Errorf("canal de consumo fechado")
}
