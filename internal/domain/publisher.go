package message

import "context"

// MessagePublisher define a interface para publicar mensagens.
type MessagePublisher interface {
	Publish(ctx context.Context, payload string) error
	Close() error
}
