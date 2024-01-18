package entity

import (
	"github.com/google/uuid"
	"time"
)

type Transaction struct {
	Sender    uuid.UUID `json:"from,required" db:"sender"`
	Receiver  uuid.UUID `json:"to,required" db:"receiver"`
	CreatedAt time.Time `json:"time" db:"created_at"`
	Amount    float64   `json:"amount,required" db:"amount"`
}
