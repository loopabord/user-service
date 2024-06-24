package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:artists,alias:p"`

	Id        string    `bun:"id,pk" json:"id"`
	Name      string    `bun:"name,notnull" json:"name"`
	Email     string    `bun:"email,notnull" json:"email"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
}
