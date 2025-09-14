package repository

import "webhook-processor-ms/internal/domain/model"

type RedisRepository interface {
	Get(keyId string) (*model.SessionModel, error)
	Save(session *model.SessionModel) (string, error)
	Delete(keyId string) error
}
