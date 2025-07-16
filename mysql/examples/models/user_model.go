package models

import (
	"database/sql"
	"github.com/gflydev/db/v2"
)

// User model
type User struct {
	// Table meta data
	MetaData db.MetaData `db:"-" model:"table:users"`

	// Table fields
	Id   int            `db:"id" model:"type:serial,primary"`
	Name sql.NullString `db:"name" model:"type:varchar(255)"`
	Age  uint8          `db:"age" model:"type:numeric"`
}
