// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: queries.sql

package sqlc_db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createJob = `-- name: CreateJob :one
INSERT INTO jobs (name, handler, interval, status, last_triggered)
VALUES ($1, $2, $3, $4, $5)
    RETURNING id, name, handler, interval, status, last_triggered
`

type CreateJobParams struct {
	Name          string
	Handler       string
	Interval      int32
	Status        string
	LastTriggered sql.NullTime
}

func (q *Queries) CreateJob(ctx context.Context, arg CreateJobParams) (Job, error) {
	row := q.db.QueryRowContext(ctx, createJob,
		arg.Name,
		arg.Handler,
		arg.Interval,
		arg.Status,
		arg.LastTriggered,
	)
	var i Job
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Handler,
		&i.Interval,
		&i.Status,
		&i.LastTriggered,
	)
	return i, err
}

const createMessage = `-- name: CreateMessage :one
INSERT INTO messages (id, content, recipient_phone_number, status, message_received_id)
VALUES ($1, $2, $3, $4, $5)
    RETURNING id, content, recipient_phone_number, status, message_received_id
`

type CreateMessageParams struct {
	ID                   uuid.UUID
	Content              string
	RecipientPhoneNumber string
	Status               string
	MessageReceivedID    uuid.NullUUID
}

type CreateMessageRow struct {
	ID                   uuid.UUID
	Content              string
	RecipientPhoneNumber string
	Status               string
	MessageReceivedID    uuid.NullUUID
}

func (q *Queries) CreateMessage(ctx context.Context, arg CreateMessageParams) (CreateMessageRow, error) {
	row := q.db.QueryRowContext(ctx, createMessage,
		arg.ID,
		arg.Content,
		arg.RecipientPhoneNumber,
		arg.Status,
		arg.MessageReceivedID,
	)
	var i CreateMessageRow
	err := row.Scan(
		&i.ID,
		&i.Content,
		&i.RecipientPhoneNumber,
		&i.Status,
		&i.MessageReceivedID,
	)
	return i, err
}

const getDueJobs = `-- name: GetDueJobs :many
SELECT id, name, handler, interval, status, last_triggered FROM jobs
WHERE status = 'active' AND last_triggered + interval * interval '1 second' < NOW()
    FOR UPDATE SKIP LOCKED
`

func (q *Queries) GetDueJobs(ctx context.Context) ([]Job, error) {
	rows, err := q.db.QueryContext(ctx, getDueJobs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Job
	for rows.Next() {
		var i Job
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Handler,
			&i.Interval,
			&i.Status,
			&i.LastTriggered,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getJobs = `-- name: GetJobs :many
SELECT id, name, handler, interval, status, last_triggered FROM jobs
`

func (q *Queries) GetJobs(ctx context.Context) ([]Job, error) {
	rows, err := q.db.QueryContext(ctx, getJobs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Job
	for rows.Next() {
		var i Job
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Handler,
			&i.Interval,
			&i.Status,
			&i.LastTriggered,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNotSendedMessages = `-- name: GetNotSendedMessages :many
SELECT id, content, recipient_phone_number, status, message_received_id
FROM messages
WHERE status = 'pending'
ORDER BY created_at ASC
    LIMIT 2
FOR UPDATE SKIP LOCKED
`

type GetNotSendedMessagesRow struct {
	ID                   uuid.UUID
	Content              string
	RecipientPhoneNumber string
	Status               string
	MessageReceivedID    uuid.NullUUID
}

func (q *Queries) GetNotSendedMessages(ctx context.Context) ([]GetNotSendedMessagesRow, error) {
	rows, err := q.db.QueryContext(ctx, getNotSendedMessages)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetNotSendedMessagesRow
	for rows.Next() {
		var i GetNotSendedMessagesRow
		if err := rows.Scan(
			&i.ID,
			&i.Content,
			&i.RecipientPhoneNumber,
			&i.Status,
			&i.MessageReceivedID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSendedMessages = `-- name: GetSendedMessages :many
SELECT id, content, recipient_phone_number, status, message_received_id
FROM messages
WHERE status = 'sent'
`

type GetSendedMessagesRow struct {
	ID                   uuid.UUID
	Content              string
	RecipientPhoneNumber string
	Status               string
	MessageReceivedID    uuid.NullUUID
}

func (q *Queries) GetSendedMessages(ctx context.Context) ([]GetSendedMessagesRow, error) {
	rows, err := q.db.QueryContext(ctx, getSendedMessages)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetSendedMessagesRow
	for rows.Next() {
		var i GetSendedMessagesRow
		if err := rows.Scan(
			&i.ID,
			&i.Content,
			&i.RecipientPhoneNumber,
			&i.Status,
			&i.MessageReceivedID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAllJobsStatus = `-- name: UpdateAllJobsStatus :exec
UPDATE jobs
SET status = $1
`

func (q *Queries) UpdateAllJobsStatus(ctx context.Context, status string) error {
	_, err := q.db.ExecContext(ctx, updateAllJobsStatus, status)
	return err
}

const updateJobLastTriggered = `-- name: UpdateJobLastTriggered :exec
UPDATE jobs
SET last_triggered = $1
WHERE id = $2
`

type UpdateJobLastTriggeredParams struct {
	LastTriggered sql.NullTime
	ID            uuid.UUID
}

func (q *Queries) UpdateJobLastTriggered(ctx context.Context, arg UpdateJobLastTriggeredParams) error {
	_, err := q.db.ExecContext(ctx, updateJobLastTriggered, arg.LastTriggered, arg.ID)
	return err
}

const updateJobStatus = `-- name: UpdateJobStatus :exec
UPDATE jobs
SET status = $1
WHERE id = $2
`

type UpdateJobStatusParams struct {
	Status string
	ID     uuid.UUID
}

func (q *Queries) UpdateJobStatus(ctx context.Context, arg UpdateJobStatusParams) error {
	_, err := q.db.ExecContext(ctx, updateJobStatus, arg.Status, arg.ID)
	return err
}

const updateMessageStatus = `-- name: UpdateMessageStatus :exec
UPDATE messages
SET status = $1
WHERE id = $2
`

type UpdateMessageStatusParams struct {
	Status string
	ID     uuid.UUID
}

func (q *Queries) UpdateMessageStatus(ctx context.Context, arg UpdateMessageStatusParams) error {
	_, err := q.db.ExecContext(ctx, updateMessageStatus, arg.Status, arg.ID)
	return err
}
