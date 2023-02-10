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

		collection, err := dao.FindCollectionByNameOrId("caqiysan7yf0wve")
		if err != nil {
			return err
		}

		// add
		new_replies_count := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "xjc2zwan",
			"name": "replies_count",
			"type": "number",
			"required": false,
			"unique": false,
			"options": {
				"min": 0,
				"max": null
			}
		}`), new_replies_count)
		collection.Schema.AddField(new_replies_count)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("caqiysan7yf0wve")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("xjc2zwan")

		return dao.SaveCollection(collection)
	})
}
