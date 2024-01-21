package entity

import (
	"github.com/google/uuid"
	"time"
)

type Transaction struct {
	Sender    uuid.UUID `json:"from" db:"sender"`
	Receiver  uuid.UUID `json:"to" db:"receiver"`
	CreatedAt time.Time `json:"time" db:"created_at"`
	Amount    float64   `json:"amount" db:"amount"`
}
