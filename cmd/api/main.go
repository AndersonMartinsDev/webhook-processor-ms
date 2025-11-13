package main

import (
	"context"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"webhook-processor-ms/cmd"
	"webhook-processor-ms/internal/infrastructure/configuration"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// configuration.LoadEnv()
	configuration.LoadLogger()
	configuration.LoadRedis()

	// 1. Cria um WaitGroup para gerenciar o encerramento das goroutines
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ms_compose := cmd.NewMSCompose()

	// Conexão com o RabbitMQ
	conn, err := amqp.Dial(ms_compose.RabbitMQURL)
	if err != nil {
		log.Fatalf("Falha ao conectar no RabbitMQ: %v", err)
	}
	// Fecha a conexão somente depois que todas as goroutines estiverem encerradas
	defer conn.Close()

	// Injeção de dependências e inicialização dos serviços
	webhook_service := ms_compose.MessageProcessorConfiguration(conn)

	// 2. Inicia o consumidor do RabbitMQ em uma goroutine e o adiciona ao WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		webhook_service.ReadMessages(ctx)
	}()

	// 3. Prepara e inicia o servidor gRPC
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Falha ao iniciar o servidor gRPC: %v", err)
	}
	s := grpc.NewServer()
	reflection.Register(s)

	wg.Add(1)
	go func() {
		defer wg.Done()
		slog.Info("Servidor gRPC do Webhook Processor iniciado na porta 50051...")
		if err := s.Serve(lis); err != nil {
			slog.Error("Falha ao servir", "error", err)
		}
	}()

	// 4. Bloqueia a execução da main até que um sinal de interrupção seja recebido
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// 5. Lógica de encerramento
	slog.Info("Sinal de interrupção recebido, iniciando o encerramento...")

	s.GracefulStop() // Encerra o servidor gRPC
	cancel()         // Cancela o contexto, sinalizando para as goroutines de background pararem

	// Espera as goroutines finalizarem
	wg.Wait()

	slog.Info("Todos os serviços foram encerrados com sucesso.")
}
