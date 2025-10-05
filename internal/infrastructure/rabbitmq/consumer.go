package rabbitmq

import (
	"context"
	"fmt"
	"log/slog"

	"webhook-processor-ms/internal/domain/message"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Consumer implementa a interface domain/message.MessageConsumer.
type Consumer struct {
	conn      *amqp.Connection
	ch        *amqp.Channel
	queue     amqp.Queue
	queueName string
}

// NewConsumer cria uma nova instância do consumidor RabbitMQ e tenta reconectar em caso de falha.
func NewConsumer(conn *amqp.Connection, queueName string) (message.MessageConsumer, error) {
	consumer := &Consumer{
		conn:      conn,
		queueName: queueName,
	}

	// Tenta se conectar e abrir um canal.
	// Se falhar, retorna um erro para que a lógica de reconexão externa lide com isso.
	if err := consumer.connect(); err != nil {
		return nil, err
	}
	return consumer, nil
}

// connect tenta abrir um canal e declarar a fila.
func (c *Consumer) connect() error {
	ch, err := c.conn.Channel()
	if err != nil {
		return fmt.Errorf("falha ao abrir um canal: %w", err)
	}

	q, err := ch.QueueDeclare(
		c.queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		return fmt.Errorf("falha ao declarar a fila: %w", err)
	}

	c.ch = ch
	c.queue = q
	return nil
}

// Consume retorna um canal de entrega para o consumo de mensagens.
func (c *Consumer) Consume(ctx context.Context) (<-chan []byte, error) {
	// A reconexão agora acontece aqui, garantindo que o canal esteja sempre vivo.
	msgs, err := c.ch.Consume(
		c.queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
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
				slog.Info("Sinal de encerramento recebido, parando o consumo.")
				return
			case messageChan <- d.Body:
			}
		}
	}()

	return messageChan, nil
}

// Close fecha o canal e a conexão.
func (c *Consumer) Close() error {
	// A reconexão da conexão principal deve ser feita no main.
	// Aqui, apenas fechamos o canal.
	if c.ch != nil {
		return c.ch.Close()
	}
	return nil
}
