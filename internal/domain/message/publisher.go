package message

import "context"

// MessagePublisher define a interface para publicar mensagens.
type MessagePublisher interface {
	Publish(context.Context, []byte) error
	Close() error
}
