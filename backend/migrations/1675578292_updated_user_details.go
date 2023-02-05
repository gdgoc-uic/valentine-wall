package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("px00yjig95x0mcw")
		if err != nil {
			return err
		}

		// add
		new_user := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "ukhwq1aw",
			"name": "user",
			"type": "relation",
			"required": true,
			"unique": true,
			"options": {
				"maxSelect": 1,
				"collectionId": "_pb_users_auth_",
				"cascadeDelete": true
			}
		}`), new_user)
		collection.Schema.AddField(new_user)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("px00yjig95x0mcw")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("ukhwq1aw")

		return dao.SaveCollection(collection)
	})
}
