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

		// remove
		collection.Schema.RemoveField("qwdmzp1u")

		// add
		new_recipient := &schema.SchemaField{}
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
		}`), new_recipient)
		collection.Schema.AddField(new_recipient)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("caqiysan7yf0wve")
		if err != nil {
			return err
		}

		// add
		del_recipient := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "qwdmzp1u",
			"name": "recipient",
			"type": "relation",
			"required": true,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"collectionId": "px00yjig95x0mcw",
				"cascadeDelete": false
			}
		}`), del_recipient)
		collection.Schema.AddField(del_recipient)

		// remove
		collection.Schema.RemoveField("idbaffkv")

		return dao.SaveCollection(collection)
	})
}
