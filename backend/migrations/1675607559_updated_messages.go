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

		// update
		edit_gifts := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "uuyjuypn",
			"name": "gifts",
			"type": "relation",
			"required": false,
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
	})
}
