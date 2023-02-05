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

		// add
		new_college_department := &schema.SchemaField{}
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
		}`), new_college_department)
		collection.Schema.AddField(new_college_department)

		// add
		new_sex := &schema.SchemaField{}
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
		}`), new_sex)
		collection.Schema.AddField(new_sex)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("_pb_users_auth_")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("8yz5dhcb")

		// remove
		collection.Schema.RemoveField("6im4sml3")

		return dao.SaveCollection(collection)
	})
}
