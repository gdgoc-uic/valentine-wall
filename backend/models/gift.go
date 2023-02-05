package models

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

var _ models.Model = (*Gift)(nil)

type Gifts []Gift

// NOTE: if you choose money, the money will be added to virtual coins

type Gift struct {
	models.BaseModel

	UID   string  `db:"uid" json:"uid"`
	Label string  `db:"label" json:"label"`
	Price float32 `db:"price" json:"price"`
}

func (gift *Gift) TableName() string {
	return "gifts"
}

func GiftQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&Gift{})
}
