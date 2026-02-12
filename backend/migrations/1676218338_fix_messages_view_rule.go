package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("caqiysan7yf0wve")
		if err != nil {
			return err
		}

		collection.ViewRule = types.Pointer("gifts:length = 0 || recipient = \"everyone\" || (@request.auth.details.id = user.id || @request.auth.details.student_id = recipient)")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("caqiysan7yf0wve")
		if err != nil {
			return err
		}

		collection.ViewRule = types.Pointer("gifts:length = 0 || recipient = \"everyone\" || (@request.auth.details.id = user.user.id || @request.auth.details.student_id = recipient)")

		return dao.SaveCollection(collection)
	})
}
