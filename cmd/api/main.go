package main

import (
	"log"
	"log/slog"
	"net"
	"os"

	"webhook-processor-ms/internal/application/handler"
	"webhook-processor-ms/internal/application/services"
	"webhook-processor-ms/internal/infrastructure/configuration"
	"webhook-processor-ms/internal/infrastructure/rabbitmq"
	pb "webhook-processor-ms/proto"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	configuration.LoadEnv()
	configuration.LoadLogger()
	// 1. Conecta ao RabbitMQ
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://rabbitmq:root@lead-docker.duckdns.org:5672/"
	}
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Falha ao conectar no RabbitMQ: %v", err)
	}
	defer conn.Close()

	// 2. Cria o publicador de mensagens (infraestrutura)
	publisher, err := rabbitmq.NewPublisher(conn)
	if err != nil {
		log.Fatalf("Falha ao criar publicador RabbitMQ: %v", err)
	}
	defer publisher.Close()

	// 3. Cria o serviço de aplicação injetando o publicador (Aplicação)
	webhookService := services.NewWebhookService(publisher)

	// 4. Cria o handler gRPC injetando o serviço (Aplicação)
	webhookHandler := handler.NewWebhookHandler(webhookService)

	// 5. Inicia o servidor gRPC
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Falha ao iniciar o servidor gRPC: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterWebhookProcessorServiceServer(s, webhookHandler)
	reflection.Register(s)

	slog.Info("Servidor gRPC do Webhook Processor iniciado na porta 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Falha ao servir: %v", err)
	}
}
