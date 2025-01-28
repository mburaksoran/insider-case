// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package sqlc_db

import (
	"database/sql"

	"github.com/google/uuid"
)

type Job struct {
	ID            uuid.UUID
	Name          string
	Handler       string
	Interval      int32
	Status        string
	LastTriggered sql.NullTime
}

type Message struct {
	ID                   uuid.UUID
	Content              string
	RecipientPhoneNumber string
	Status               string
	MessageReceivedID    uuid.NullUUID
}
