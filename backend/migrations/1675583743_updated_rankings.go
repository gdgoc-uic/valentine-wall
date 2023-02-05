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

		// remove
		collection.Schema.RemoveField("xovugjt0")

		// add
		new_recipient := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "optuzp59",
			"name": "recipient",
			"type": "text",
			"required": true,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_recipient)
		collection.Schema.AddField(new_recipient)

		// add
		new_college_department := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "j0a0c0dg",
			"name": "college_department",
			"type": "relation",
			"required": true,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"collectionId": "yhv9suo8ru0esf0",
				"cascadeDelete": false
			}
		}`), new_college_department)
		collection.Schema.AddField(new_college_department)

		// add
		new_sex := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "tnebrg5b",
			"name": "sex",
			"type": "text",
			"required": true,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_sex)
		collection.Schema.AddField(new_sex)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("ocpdx07v34h97tx")
		if err != nil {
			return err
		}

		// add
		del_recipient := &schema.SchemaField{}
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
		}`), del_recipient)
		collection.Schema.AddField(del_recipient)

		// remove
		collection.Schema.RemoveField("optuzp59")

		// remove
		collection.Schema.RemoveField("j0a0c0dg")

		// remove
		collection.Schema.RemoveField("tnebrg5b")

		return dao.SaveCollection(collection)
	})
}
