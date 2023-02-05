package models

import (
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

var _ models.Model = (*User)(nil)

// TODO:

type User struct {
	models.BaseModel

	ID        string         `db:"id" json:"id"`
	Recipient string         `db:"recipient" json:"recipient" validate:"required,min=6,max=12,numeric"`
	User      string         `db:"user" json:"user"`
	Content   string         `db:"content" json:"content" validate:"required,max=240"`
	Gifts     []string       `db:"gifts" json:"gifts"`
	Deleted   types.DateTime `db:"deleted" json:"deleted"`
}

func (msg *User) TableName() string {
	return "users"
}
