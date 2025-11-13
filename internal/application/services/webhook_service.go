package services

import (
	"context"
	"log/slog"
	"sync"
	"time" // Importação necessária
	"webhook-processor-ms/internal/domain/message"
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
func (s *WebhookService) ReadMessages(ctx context.Context) {
	var wg sync.WaitGroup
	// --- 1. Iniciar Processamento PF (WWEBJS) ---
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.processMessageFeature(ctx, NewWWEBJSProcessMessage(s.ConsumerPF, *s.SessionService, s.AgentModelService, s.Publisher))
	}()
	// --- 2. Iniciar Processamento PJ (META) ---
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.processMessageFeature(ctx, NewMetaProcessMessage(s.ConsumerPJ, *s.SessionService, s.AgentModelService, s.Publisher))
	}()
	slog.Info("Todos os consumidores iniciados. Esperando por cancelamento do contexto...")
	wg.Wait()
	slog.Info("Todos os processos de consumo encerrados.")
}

func (s *WebhookService) processMessageFeature(ctx context.Context, processor ProcessMessageInterface) {
	for {
		select {
		case <-ctx.Done():
			slog.Info("Consumidor: Contexto cancelado, encerrando.")
			return
		default:
			err := processor.ProcessMessages(ctx)
			if err != nil {
				slog.Error("Consumidor: Falha na conexão ou consumo, tentando reconectar em 5 segundos...", "error", err)
				time.Sleep(5 * time.Second)
			} else {
				// Evita loops rápidos se a função retornar inesperadamente sem erro
				time.Sleep(1 * time.Second)
			}
		}
	}

}
