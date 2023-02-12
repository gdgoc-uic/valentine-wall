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
		edit_user := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "nfu7q7wd",
			"name": "user",
			"type": "relation",
			"required": false,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"collectionId": "px00yjig95x0mcw",
				"cascadeDelete": false
			}
		}`), edit_user)
		collection.Schema.AddField(edit_user)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("caqiysan7yf0wve")
		if err != nil {
			return err
		}

		// update
		edit_user := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "nfu7q7wd",
			"name": "user",
			"type": "relation",
			"required": true,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"collectionId": "px00yjig95x0mcw",
				"cascadeDelete": false
			}
		}`), edit_user)
		collection.Schema.AddField(edit_user)

		return dao.SaveCollection(collection)
	})
}
