package service

import (
	"context"
	"fmt"
	"github.com/mburaksoran/insider-case/internal/domain/types"

	"github.com/google/uuid"
	"github.com/pkg/errors"

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
			return nil, errors.Wrapf(err, "[CreateMessage] - CreateMessage Error")
		}
		return result, nil
	})

	return err
}

func (s *messageService) GetMessageToSend(ctx context.Context) ([]*models.Message, error) {
	result, err := s.messageRepository.WithTransaction(ctx, func(q *sqlc_db.Queries) (interface{}, error) {
		jobs, err := s.messageRepository.GetMessageNotSent(ctx, q)
		if err != nil {
			return nil, errors.Wrapf(err, "[GetMessageToSend] - GetMessageNotSent Error")
		}
		if len(jobs) < 1 {
			return nil, nil
		}
		for _, job := range jobs {
			updateErr := s.messageRepository.UpdateMessageStatus(ctx, q, job.ID, types.MessageStatusInProgress) //TODO create new const for status
			if updateErr != nil {
				return nil, errors.Wrapf(err, fmt.Sprintf("[UpdateMessageStatus] - UpdateMessageStatus error status JobId: %s", job.ID))
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
			return nil, errors.Wrapf(err, "[UpdateMessageStatus] - UpdateMessageStatus Error")
		}
		return nil, nil
	})

	return err
}

func (s *messageService) GetSendMessage(ctx context.Context) ([]*models.Message, error) {
	result, err := s.messageRepository.WithoutTransaction(ctx, func(queries *sqlc_db.Queries) (interface{}, error) {
		res, err := s.messageRepository.GetSendMessages(ctx, queries)
		if err != nil {
			return nil, errors.Wrapf(err, "[GetSendMessage] - GetSendMessages Error")
		}
		return res, nil
	})
	var messageList []*models.Message

	if result != nil {
		messageList = result.([]*models.Message)
		return messageList, err
	}
	return nil, err
}
