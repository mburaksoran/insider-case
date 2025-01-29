package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/mburaksoran/insider-case/internal/domain/models"
)

type MessageServiceInterface interface {
	CreateMessage(ctx context.Context, job models.Message) error
	GetMessageToSend(ctx context.Context) ([]*models.Message, error)
	UpdateMessageStatus(ctx context.Context, uuid uuid.UUID, status string) error
	GetSendMessage(ctx context.Context) ([]*models.Message, error)
}
