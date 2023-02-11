package models

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

var _ models.Model = (*CollegeDepartment)(nil)

type CollegeDepartment struct {
	models.BaseModel

	UID   string `db:"uid" json:"uid"`
	Label string `db:"label" json:"label"`
}

func (dept *CollegeDepartment) TableName() string {
	return "college_departments"
}

func DepartmentQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&CollegeDepartment{})
}
