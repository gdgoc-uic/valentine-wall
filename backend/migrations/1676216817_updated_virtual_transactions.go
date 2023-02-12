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

		collection, err := dao.FindCollectionByNameOrId("rocfs9e910vqq5g")
		if err != nil {
			return err
		}

		// update
		edit_wallet := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "gmnhl7zd",
			"name": "wallet",
			"type": "relation",
			"required": false,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"collectionId": "rs07r3dxff0hxbz",
				"cascadeDelete": false
			}
		}`), edit_wallet)
		collection.Schema.AddField(edit_wallet)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("rocfs9e910vqq5g")
		if err != nil {
			return err
		}

		// update
		edit_wallet := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "gmnhl7zd",
			"name": "wallet",
			"type": "relation",
			"required": true,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"collectionId": "rs07r3dxff0hxbz",
				"cascadeDelete": false
			}
		}`), edit_wallet)
		collection.Schema.AddField(edit_wallet)

		return dao.SaveCollection(collection)
	})
}
