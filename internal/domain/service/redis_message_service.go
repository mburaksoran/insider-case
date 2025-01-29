package service

import (
	"context"
	"github.com/mburaksoran/insider-case/internal/app/service"
	"github.com/mburaksoran/insider-case/internal/domain/models"
	"github.com/mburaksoran/insider-case/internal/domain/repository"
	"go.uber.org/zap"
)

type redisMessageService struct {
	redisMessageRepository repository.RedisMessageRepositoryInterface
	logger                 *zap.SugaredLogger
}

func NewRedisMessageService(redisMessageRepository repository.RedisMessageRepositoryInterface, lgr *zap.SugaredLogger) service.RedisMessageServiceInterface {
	return &redisMessageService{
		redisMessageRepository: redisMessageRepository,
		logger:                 lgr,
	}
}

func (s *redisMessageService) SetMessageReceiveHistory(ctx context.Context, msg models.MessageReceiveHistory) error {
	err := s.redisMessageRepository.Set(ctx, msg.ID.String(), msg)
	if err != nil {
		s.logger.Errorf("[SetMessageReceiveHistory] Error while setting message receive history: %v", err)
		return err
	}

	return nil
}
