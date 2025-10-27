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
	ConsumerPJ        message.MessageConsumer
	SessionService    *SessionService
	AgentModelService *AgentModelService
}

// NewWebhookService cria uma nova instância de WebhookService.
func NewWebhookService(publisher message.MessagePublisher, consumer message.MessageConsumer, consumerPj message.MessageConsumer, sessionService *SessionService,
	agentModelService *AgentModelService) *WebhookService {
	return &WebhookService{
		Publisher:         publisher,
		ConsumerPF:        consumer,
		ConsumerPJ:        consumerPj,
		SessionService:    sessionService,
		AgentModelService: agentModelService,
	}
}

// ProcessMessages inicia o loop de processamento com reconexão.
func (s *WebhookService) ProcessMessages(ctx context.Context) {
	// Usamos um sync.WaitGroup para esperar que todas as goroutines de processamento terminem,
	// embora neste caso, elas só terminem se o contexto for cancelado.
	var wg sync.WaitGroup

	// --- 1. Iniciar Processamento PF (WWEBJS) ---
	wg.Add(1)
	go func() {
		defer wg.Done()
		// O loop 'for' de reconexão deve estar DENTRO desta goroutine
		for {
			select {
			case <-ctx.Done():
				slog.Info("Consumidor PF: Contexto cancelado, encerrando.")
				return
			default:
				err := s.processMessagesPF(ctx)
				if err != nil {
					slog.Error("Consumidor PF: Falha na conexão ou consumo, tentando reconectar em 5 segundos...", "error", err)
					time.Sleep(5 * time.Second)
				} else {
					// Evita loops rápidos se a função retornar inesperadamente sem erro
					time.Sleep(1 * time.Second)
				}
			}
		}
	}()

	// --- 2. Iniciar Processamento PJ (META) ---
	wg.Add(1)
	go func() {
		defer wg.Done()
		// O loop 'for' de reconexão deve estar DENTRO desta goroutine
		for {
			select {
			case <-ctx.Done():
				slog.Info("Consumidor PJ: Contexto cancelado, encerrando.")
				return
			default:
				err := s.processMessagesPJ(ctx)
				if err != nil {
					slog.Error("Consumidor PJ: Falha na conexão ou consumo, tentando reconectar em 5 segundos...", "error", err)
					time.Sleep(5 * time.Second)
				} else {
					// Evita loops rápidos se a função retornar inesperadamente sem erro
					time.Sleep(1 * time.Second)
				}
			}
		}
	}()

	slog.Info("Todos os consumidores iniciados. Esperando por cancelamento do contexto...")
	// Bloqueia e espera que o Contexto seja cancelado e as goroutines terminem
	wg.Wait()
	slog.Info("Todos os processos de consumo encerrados.")
}

// processMessagesPF contém a lógica de consumo de um ciclo.
func (s *WebhookService) processMessagesPF(ctx context.Context) error {
	// Consome da fila de webhooks brutos (o nome da fila deve ser o mesmo usado no gateway-ms)
	msgs, err := s.ConsumerPF.Consume(ctx)
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
		slog.Info("Lendo mensagem!!")

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
					slog.Error("Error ao recuperar agent correto! %v", err)
				}
			}()

			go func() {
				defer wg.Done()
				session, err = s.SessionService.GetHistory(From)
				if err != nil {
					slog.Error("Error ao recuperar o historico! %v", err)
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
				MediaURL:      webhookEvent.Data.DeprecatedMms3Url,
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

func (s *WebhookService) processMessagesPJ(ctx context.Context) error {
	// Consome da fila de webhooks brutos (o nome da fila deve ser o mesmo usado no gateway-ms)
	msgs, err := s.ConsumerPJ.Consume(ctx)
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
