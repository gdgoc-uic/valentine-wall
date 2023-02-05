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

		collection, err := dao.FindCollectionByNameOrId("_pb_users_auth_")
		if err != nil {
			return err
		}

		// update
		edit_details := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "lz0uypgf",
			"name": "details",
			"type": "relation",
			"required": false,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"collectionId": "px00yjig95x0mcw",
				"cascadeDelete": true
			}
		}`), edit_details)
		collection.Schema.AddField(edit_details)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("_pb_users_auth_")
		if err != nil {
			return err
		}

		// update
		edit_details := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "lz0uypgf",
			"name": "details",
			"type": "relation",
			"required": true,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"collectionId": "px00yjig95x0mcw",
				"cascadeDelete": true
			}
		}`), edit_details)
		collection.Schema.AddField(edit_details)

		return dao.SaveCollection(collection)
	})
}
