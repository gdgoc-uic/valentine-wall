package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("35mnuyxwxc8xvs6")
		if err != nil {
			return err
		}

		collection.ListRule = types.Pointer("message.recipient = \"everyone\" || @request.auth.details.id = message.user.id || @request.auth.details.id = sender.id")

		collection.ViewRule = types.Pointer("message.recipient = \"everyone\" || @request.auth.details.id = message.user.id || @request.auth.details.id = sender.id")

		collection.CreateRule = types.Pointer("@request.auth.details.id = @request.data.sender.id || (message.recipient = \"everyone\" && @request.auth.id != \"\")")

		collection.UpdateRule = types.Pointer("@request.auth.details.id = sender.id")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("35mnuyxwxc8xvs6")
		if err != nil {
			return err
		}

		collection.ListRule = nil

		collection.ViewRule = nil

		collection.CreateRule = nil

		collection.UpdateRule = nil

		return dao.SaveCollection(collection)
	})
}
