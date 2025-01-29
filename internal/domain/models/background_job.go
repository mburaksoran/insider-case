package models

import (
	"github.com/google/uuid"
	"time"
)

type BackgroundJob struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Handler       string    `json:"handler"`
	Interval      int32     `json:"interval"`
	Status        string    `json:"status"`
	LastTriggered time.Time `json:"last_triggered"`
}
