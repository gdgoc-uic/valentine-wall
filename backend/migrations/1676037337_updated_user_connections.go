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

		collection, err := dao.FindCollectionByNameOrId("55bmnoqmzjhnng0")
		if err != nil {
			return err
		}

		// add
		new_provider := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "v9kwnaqi",
			"name": "provider",
			"type": "text",
			"required": true,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_provider)
		collection.Schema.AddField(new_provider)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("55bmnoqmzjhnng0")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("v9kwnaqi")

		return dao.SaveCollection(collection)
	})
}
