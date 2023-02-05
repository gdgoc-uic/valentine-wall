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

		// remove
		collection.Schema.RemoveField("users_name")

		// remove
		collection.Schema.RemoveField("h677c9aw")

		// remove
		collection.Schema.RemoveField("8yz5dhcb")

		// remove
		collection.Schema.RemoveField("6im4sml3")

		// remove
		collection.Schema.RemoveField("i6iyb54r")

		// add
		new_details := &schema.SchemaField{}
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
		}`), new_details)
		collection.Schema.AddField(new_details)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("_pb_users_auth_")
		if err != nil {
			return err
		}

		// add
		del_name := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "users_name",
			"name": "name",
			"type": "text",
			"required": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), del_name)
		collection.Schema.AddField(del_name)

		// add
		del_last_active := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "h677c9aw",
			"name": "last_active",
			"type": "date",
			"required": false,
			"unique": false,
			"options": {
				"min": "",
				"max": ""
			}
		}`), del_last_active)
		collection.Schema.AddField(del_last_active)

		// add
		del_college_department := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "8yz5dhcb",
			"name": "college_department",
			"type": "relation",
			"required": false,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"collectionId": "yhv9suo8ru0esf0",
				"cascadeDelete": false
			}
		}`), del_college_department)
		collection.Schema.AddField(del_college_department)

		// add
		del_sex := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "6im4sml3",
			"name": "sex",
			"type": "select",
			"required": false,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"values": [
					"unknown",
					"male",
					"female"
				]
			}
		}`), del_sex)
		collection.Schema.AddField(del_sex)

		// add
		del_student_id := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "i6iyb54r",
			"name": "student_id",
			"type": "text",
			"required": false,
			"unique": false,
			"options": {
				"min": 11,
				"max": 12,
				"pattern": ""
			}
		}`), del_student_id)
		collection.Schema.AddField(del_student_id)

		// remove
		collection.Schema.RemoveField("lz0uypgf")

		return dao.SaveCollection(collection)
	})
}
