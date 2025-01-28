package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/mburaksoran/insider-case/internal/domain/models"
	"github.com/mburaksoran/insider-case/internal/shared/sqlc_db"
)

type MessageRepositoryInterface interface {
	CreateMessage(ctx context.Context, queries *sqlc_db.Queries, msg models.Message) (bool, error)
	GetMessageNotSent(ctx context.Context, queries *sqlc_db.Queries) ([]*models.Message, error)
	UpdateMessageStatus(ctx context.Context, queries *sqlc_db.Queries, id uuid.UUID, status string) error
	WithTransaction(ctx context.Context, fn func(*sqlc_db.Queries) (interface{}, error)) (interface{}, error)
	WithoutTransaction(ctx context.Context, fn func(*sqlc_db.Queries) (interface{}, error)) (interface{}, error)
}
