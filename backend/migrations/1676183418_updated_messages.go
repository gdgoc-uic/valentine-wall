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
		edit_recipient := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "idbaffkv",
			"name": "recipient",
			"type": "text",
			"required": true,
			"unique": false,
			"options": {
				"min": 8,
				"max": 12,
				"pattern": "[0-9]{11,12}|everyone"
			}
		}`), edit_recipient)
		collection.Schema.AddField(edit_recipient)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("caqiysan7yf0wve")
		if err != nil {
			return err
		}

		// update
		edit_recipient := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "idbaffkv",
			"name": "recipient",
			"type": "text",
			"required": true,
			"unique": false,
			"options": {
				"min": 11,
				"max": 12,
				"pattern": ""
			}
		}`), edit_recipient)
		collection.Schema.AddField(edit_recipient)

		return dao.SaveCollection(collection)
	})
}
