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
		new_last_active := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "erbuqwm0",
			"name": "last_active",
			"type": "date",
			"required": false,
			"unique": false,
			"options": {
				"min": "",
				"max": ""
			}
		}`), new_last_active)
		collection.Schema.AddField(new_last_active)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("px00yjig95x0mcw")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("erbuqwm0")

		return dao.SaveCollection(collection)
	})
}
