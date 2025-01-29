package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type MessageReceiveHistory struct {
	ID          uuid.UUID
	SendingTime time.Time
}

func (m MessageReceiveHistory) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary metodu, []byte'ı MessageReceiveHistory türüne dönüştürür.
func (m *MessageReceiveHistory) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}
