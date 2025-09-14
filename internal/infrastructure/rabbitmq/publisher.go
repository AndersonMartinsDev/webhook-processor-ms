package rabbitmq

import (
	"context"
	"fmt"
	"webhook-processor-ms/internal/domain/message"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publisher implementa a interface domain/message.MessagePublisher.
type Publisher struct {
	conn      *amqp.Connection
	ch        *amqp.Channel
	queueName string
}

// NewPublisher cria uma nova instância do publicador RabbitMQ.
func NewPublisher(conn *amqp.Connection, queueName string) (message.MessagePublisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir um canal: %w", err)
	}
	// Declara a fila (cria se não existir)
	_, err = ch.QueueDeclare(
		queueName, // nome da fila
		true,      // durável
		false,     // auto-delete
		false,     // exclusiva
		false,     // noWait
		nil,       // argumentos
	)
	if err != nil {
		return nil, fmt.Errorf("falha ao declarar a fila: %w", err)
	}

	return &Publisher{
		conn:      conn,
		ch:        ch,
		queueName: queueName,
	}, nil
}

// Publish envia uma mensagem para a fila do RabbitMQ.
func (p *Publisher) Publish(ctx context.Context, payload []byte) error {
	return p.ch.PublishWithContext(
		ctx,
		"",          // exchange
		p.queueName, // routing key (nome da fila)
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
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
