package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/mburaksoran/insider-case/internal/domain/models"
	"github.com/mburaksoran/insider-case/internal/domain/repository"
	"github.com/mburaksoran/insider-case/internal/shared/sqlc_db"
)

type messageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) repository.MessageRepositoryInterface {
	return &messageRepository{db: db}
}

func (r *messageRepository) WithoutTransaction(ctx context.Context, fn func(*sqlc_db.Queries) (interface{}, error)) (interface{}, error) {
	q := sqlc_db.New(r.db)
	return fn(q)
}

func (r *messageRepository) WithTransaction(ctx context.Context, fn func(*sqlc_db.Queries) (interface{}, error)) (interface{}, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	q := sqlc_db.New(tx)
	res, err := fn(q)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return res, tx.Commit()
}

func (r *messageRepository) CreateMessage(ctx context.Context, queries *sqlc_db.Queries, msg models.Message) (bool, error) {

	_, err := queries.CreateMessage(ctx, sqlc_db.CreateMessageParams{
		Content:              msg.Content,
		RecipientPhoneNumber: msg.RecipientPhoneNumber,
		Status:               msg.Status,
		MessageReceivedID:    uuid.NullUUID{},
	})

	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *messageRepository) GetMessageNotSent(ctx context.Context, queries *sqlc_db.Queries) ([]*models.Message, error) {

	result, err := queries.GetNotSendedMessages(ctx)
	if err != nil {
		return nil, err
	}
	var msgList = []*models.Message{}
	for _, msg := range result {
		msgList = append(msgList, &models.Message{
			ID:                   msg.ID,
			Content:              msg.Content,
			RecipientPhoneNumber: msg.RecipientPhoneNumber,
			Status:               msg.Status,
			MessageReceivedId:    uuid.UUID{},
		})
	}

	return msgList, nil
}

func (r *messageRepository) UpdateMessageStatus(ctx context.Context, queries *sqlc_db.Queries, id uuid.UUID, status string) error {

	err := queries.UpdateMessageStatus(ctx, sqlc_db.UpdateMessageStatusParams{
		Status: status,
		ID:     id,
	})
	if err != nil {
		return err
	}
	return nil
}
