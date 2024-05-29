package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:artists,alias:p"`

	Id        string    `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	Name      string    `bun:"name,notnull" json:"name"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
}
