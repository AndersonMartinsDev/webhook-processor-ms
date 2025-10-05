package cmd

import (
	"log"
	"os"

	"webhook-processor-ms/internal/application/services"
	"webhook-processor-ms/internal/infrastructure/grpc_client"
	"webhook-processor-ms/internal/infrastructure/rabbitmq"
	"webhook-processor-ms/internal/infrastructure/repository"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MSCompose struct {
	AgentModelMSUrl       string
	WebHookProcessorMSURL string
	RabbitMQURL           string
}

func NewMSCompose() *MSCompose {
	return &MSCompose{
		AgentModelMSUrl: os.Getenv("AGENT_MODEL_URL"),
		RabbitMQURL:     os.Getenv("RABBITMQ_URL"),
	}
}

func (manager MSCompose) MessageProcessorConfiguration(conn *amqp.Connection) *services.WebhookService {

	agent_model_client := grpc_client.NewAgentModelClient(grpc_client.GrcpConnection(manager.AgentModelMSUrl))
	publisher, err := rabbitmq.NewPublisher(conn, "webhook-whatsapp-messages")
	if err != nil {
		log.Fatalf("Falha ao criar publicador RabbitMQ: %v", err)
	}
	consumerPf, err := rabbitmq.NewConsumer(conn, "whatsapp-webhooks-pf-raw")
	if err != nil {
		log.Fatalf("Falha ao criar consumidor RabbitMQ: %v", err)
	}
	redisRepository := repository.NewRedisRepository()
	sessionService := services.NewSessionService(redisRepository)
	agentModelService := services.NewAgentModelService(agent_model_client)

	return services.NewWebhookService(publisher, consumerPf, sessionService, agentModelService)
}
