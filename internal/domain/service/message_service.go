package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"

	"github.com/mburaksoran/insider-case/internal/app/service"
	"github.com/mburaksoran/insider-case/internal/domain/models"
	"github.com/mburaksoran/insider-case/internal/domain/repository"
	"github.com/mburaksoran/insider-case/internal/shared/sqlc_db"
)

type messageService struct {
	messageRepository repository.MessageRepositoryInterface
}

func NewMessageService(messageRepo repository.MessageRepositoryInterface) service.MessageServiceInterface {
	return &messageService{
		messageRepository: messageRepo,
	}
}

func (s *messageService) CreateMessage(ctx context.Context, msg models.Message) error {
	_, err := s.messageRepository.WithoutTransaction(ctx, func(queries *sqlc_db.Queries) (interface{}, error) {
		result, err := s.messageRepository.CreateMessage(ctx, queries, msg)
		if err != nil {
			log.Printf("Error while creating Job: %v", err)
			return nil, err
		}
		return result, nil
	})

	return err
}

func (s *messageService) GetMessageToSend(ctx context.Context) ([]*models.Message, error) {
	result, err := s.messageRepository.WithTransaction(ctx, func(q *sqlc_db.Queries) (interface{}, error) {
		jobs, err := s.messageRepository.GetMessageNotSent(ctx, q)
		if err != nil {
			log.Printf("Error while getting messages: %v", err)
			return nil, err
		}
		if len(jobs) < 1 {
			return nil, nil
		}
		for _, job := range jobs {
			updateErr := s.messageRepository.UpdateMessageStatus(ctx, q, job.ID, "processing") //TODO create new const for status
			if updateErr != nil {
				return nil, errors.New(fmt.Sprintf("errors while updating message status JobId: %s error: %s", job.ID, err.Error()))
			}
		}
		return jobs, nil
	})
	if err != nil {
		return nil, err
	}
	if result != nil {
		msgList := result.([]*models.Message)
		return msgList, nil
	}
	return nil, nil
}

func (s *messageService) UpdateMessageStatus(ctx context.Context, uuid uuid.UUID, status string) error {
	_, err := s.messageRepository.WithoutTransaction(ctx, func(queries *sqlc_db.Queries) (interface{}, error) {
		err := s.messageRepository.UpdateMessageStatus(ctx, queries, uuid, status)
		if err != nil {
			log.Printf("Error while creating Job: %v", err)
			return nil, err
		}
		return nil, nil
	})

	return err
}
