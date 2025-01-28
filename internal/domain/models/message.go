package models

import (
	"github.com/google/uuid"
)

type Message struct {
	ID                   uuid.UUID `json:"id"`
	Content              string    `json:"content"`
	RecipientPhoneNumber string    `json:"recipient_phone_number"`
	Status               string    `json:"status"`
	MessageReceivedId    uuid.UUID `json:"message_received_id"`
}
