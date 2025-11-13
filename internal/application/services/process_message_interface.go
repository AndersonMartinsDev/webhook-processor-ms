package services

import "context"

type ProcessMessageInterface interface {
	ProcessMessages(ctx context.Context) error
}
