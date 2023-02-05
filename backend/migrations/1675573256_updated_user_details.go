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
		new_sex := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "2gizlajc",
			"name": "sex",
			"type": "select",
			"required": true,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"values": [
					"unknown",
					"male",
					"female"
				]
			}
		}`), new_sex)
		collection.Schema.AddField(new_sex)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("px00yjig95x0mcw")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("2gizlajc")

		return dao.SaveCollection(collection)
	})
}
