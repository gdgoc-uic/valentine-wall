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

		collection, err := dao.FindCollectionByNameOrId("rs07r3dxff0hxbz")
		if err != nil {
			return err
		}

		// update
		edit_balance := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "0abw1isr",
			"name": "balance",
			"type": "number",
			"required": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null
			}
		}`), edit_balance)
		collection.Schema.AddField(edit_balance)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("rs07r3dxff0hxbz")
		if err != nil {
			return err
		}

		// update
		edit_balance := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "0abw1isr",
			"name": "balance",
			"type": "number",
			"required": true,
			"unique": false,
			"options": {
				"min": null,
				"max": null
			}
		}`), edit_balance)
		collection.Schema.AddField(edit_balance)

		return dao.SaveCollection(collection)
	})
}
