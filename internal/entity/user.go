package entity

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	UUID      uuid.UUID `db:"uuid"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}
