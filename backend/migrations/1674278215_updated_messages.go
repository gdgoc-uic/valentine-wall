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
		new_field1 := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "fzt2anh4",
			"name": "field1",
			"type": "text",
			"required": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_field1)
		collection.Schema.AddField(new_field1)

		// update
		edit_gifts := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "uuyjuypn",
			"name": "gifts",
			"type": "relation",
			"required": true,
			"unique": false,
			"options": {
				"maxSelect": null,
				"collectionId": "z3ucff2sh922dhz",
				"cascadeDelete": false
			}
		}`), edit_gifts)
		collection.Schema.AddField(edit_gifts)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("caqiysan7yf0wve")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("fzt2anh4")

		// update
		edit_gifts := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "uuyjuypn",
			"name": "field",
			"type": "relation",
			"required": true,
			"unique": false,
			"options": {
				"maxSelect": null,
				"collectionId": "z3ucff2sh922dhz",
				"cascadeDelete": false
			}
		}`), edit_gifts)
		collection.Schema.AddField(edit_gifts)

		return dao.SaveCollection(collection)
	})
}
