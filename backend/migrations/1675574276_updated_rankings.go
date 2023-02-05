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

		collection, err := dao.FindCollectionByNameOrId("ocpdx07v34h97tx")
		if err != nil {
			return err
		}

		// update
		edit_recipient := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "xovugjt0",
			"name": "recipient",
			"type": "relation",
			"required": true,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"collectionId": "px00yjig95x0mcw",
				"cascadeDelete": true
			}
		}`), edit_recipient)
		collection.Schema.AddField(edit_recipient)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("ocpdx07v34h97tx")
		if err != nil {
			return err
		}

		// update
		edit_recipient := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "xovugjt0",
			"name": "user",
			"type": "relation",
			"required": true,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"collectionId": "px00yjig95x0mcw",
				"cascadeDelete": true
			}
		}`), edit_recipient)
		collection.Schema.AddField(edit_recipient)

		return dao.SaveCollection(collection)
	})
}
