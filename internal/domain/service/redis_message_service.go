package service

import (
	"context"
	"github.com/mburaksoran/insider-case/internal/app/service"
	"github.com/mburaksoran/insider-case/internal/domain/models"
	"github.com/mburaksoran/insider-case/internal/domain/repository"
	"log"
)

type redisMessageService struct {
	redisMessageRepository repository.RedisMessageRepositoryInterface
}

func NewRedisMessageService(redisMessageRepository repository.RedisMessageRepositoryInterface) service.RedisMessageServiceInterface {
	return &redisMessageService{
		redisMessageRepository: redisMessageRepository,
	}
}

func (s *redisMessageService) SetMessageReceiveHistory(ctx context.Context, msg models.MessageReceiveHistory) error {
	err := s.redisMessageRepository.Set(ctx, msg.ID.String(), msg)
	if err != nil {
		log.Printf("Error while setting message receive history: %v", err)
		return err
	}

	return nil
}
