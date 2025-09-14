package services

import (
	"fmt"
	"log/slog"
	"webhook-processor-ms/internal/domain/model"
	"webhook-processor-ms/internal/domain/repository"
)

type SessionService struct {
	Repository repository.RedisRepository
}

func NewSessionService(repository repository.RedisRepository) *SessionService {
	return &SessionService{
		Repository: repository,
	}
}

func (cache *SessionService) SaveHistory(session model.SessionModel) error {

	if _, err := cache.Repository.Save(&session); err != nil {
		slog.Error("erro ao salvar sessão e histórico:")
		return fmt.Errorf("erro ao salvar sessão e histórico: %w", err)
	}
	return nil
}

func (cache *SessionService) GetHistory(keyId string) (model.SessionModel, error) {

	register, _ := cache.Repository.Get(keyId)
	if register == nil {
		return model.SessionModel{}, nil
	}

	return model.SessionModel{
		IDSession: keyId,
		AgentId:   register.AgentId,
	}, nil
}
