package rabbitmq

import (
	"context"
	"fmt"

	"webhook-processor-ms/internal/domain/message"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Consumer implementa a interface domain/message.MessageConsumer.
type Consumer struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue amqp.Queue
}

// NewConsumer cria uma nova instância do consumidor RabbitMQ.
func NewConsumer(conn *amqp.Connection, queueName string) (message.MessageConsumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir um canal: %w", err)
	}

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("falha ao declarar a fila: %w", err)
	}

	return &Consumer{
		conn:  conn,
		ch:    ch,
		queue: q,
	}, nil
}

// Consume retorna um canal de entrega para o consumo de mensagens.
func (c *Consumer) Consume(ctx context.Context) (<-chan []byte, error) {
	msgs, err := c.ch.Consume(
		c.queue.Name,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, fmt.Errorf("falha ao consumir mensagens: %w", err)
	}

	// Canal para enviar as mensagens para o serviço
	messageChan := make(chan []byte)

	go func() {
		defer close(messageChan)
		for d := range msgs {
			select {
			case <-ctx.Done():
				return
			case messageChan <- d.Body:
			}
		}
	}()

	return messageChan, nil
}

// Close fecha o canal e a conexão.
func (c *Consumer) Close() error {
	if c.ch != nil {
		c.ch.Close()
	}
	return nil
}
