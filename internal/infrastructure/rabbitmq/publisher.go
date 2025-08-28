package rabbitmq

import (
	"context"
	"fmt"
	message "webhook-processor-ms/internal/domain"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publisher implementa a interface domain/message.MessagePublisher.
type Publisher struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

// NewPublisher cria uma nova instância do publicador RabbitMQ.
func NewPublisher(conn *amqp.Connection) (message.MessagePublisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir um canal: %w", err)
	}

	// Declara a fila (cria se não existir)
	_, err = ch.QueueDeclare(
		"webhook_queue", // nome da fila
		true,            // durável
		false,           // auto-delete
		false,           // exclusiva
		false,           // noWait
		nil,             // argumentos
	)
	if err != nil {
		return nil, fmt.Errorf("falha ao declarar a fila: %w", err)
	}

	return &Publisher{
		conn: conn,
		ch:   ch,
	}, nil
}

// Publish envia uma mensagem para a fila do RabbitMQ.
func (p *Publisher) Publish(ctx context.Context, payload string) error {
	return p.ch.PublishWithContext(
		ctx,
		"",              // exchange
		"webhook_queue", // routing key (nome da fila)
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(payload),
		},
	)
}

// Close fecha o canal e a conexão.
func (p *Publisher) Close() error {
	if p.ch != nil {
		p.ch.Close()
	}
	return nil
}
