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
		collection.Schema.RemoveField("dvavc1rl")

		// remove
		collection.Schema.RemoveField("ylooys31")

		// remove
		collection.Schema.RemoveField("lgeckfbq")

		// add
		new_user := &schema.SchemaField{}
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
		}`), new_user)
		collection.Schema.AddField(new_user)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("ocpdx07v34h97tx")
		if err != nil {
			return err
		}

		// add
		del_recipient_id := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "dvavc1rl",
			"name": "recipient_id",
			"type": "text",
			"required": true,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), del_recipient_id)
		collection.Schema.AddField(del_recipient_id)

		// add
		del_department := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "ylooys31",
			"name": "department",
			"type": "text",
			"required": true,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), del_department)
		collection.Schema.AddField(del_department)

		// add
		del_sex := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "lgeckfbq",
			"name": "sex",
			"type": "text",
			"required": true,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), del_sex)
		collection.Schema.AddField(del_sex)

		// remove
		collection.Schema.RemoveField("xovugjt0")

		return dao.SaveCollection(collection)
	})
}
