package message

import "context"

// MessageConsumer define a interface para consumir mensagens de uma fila.
type MessageConsumer interface {
	Consume(ctx context.Context) (<-chan []byte, error)
	Close() error
}
