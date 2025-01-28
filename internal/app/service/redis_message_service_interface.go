package service

import (
	"context"

	"github.com/mburaksoran/insider-case/internal/domain/models"
)

type RedisMessageServiceInterface interface {
	SetMessageReceiveHistory(ctx context.Context, msg models.MessageReceiveHistory) error
}
