package repository

import (
	"context"
	"github.com/mburaksoran/insider-case/internal/domain/models"
)

type RedisMessageRepositoryInterface interface {
	Set(ctx context.Context, key string, value models.MessageReceiveHistory) error
}
