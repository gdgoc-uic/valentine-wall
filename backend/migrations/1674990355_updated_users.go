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
		edit_student_id := &schema.SchemaField{}
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
		}`), edit_student_id)
		collection.Schema.AddField(edit_student_id)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("_pb_users_auth_")
		if err != nil {
			return err
		}

		// update
		edit_student_id := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "i6iyb54r",
			"name": "student_id",
			"type": "text",
			"required": false,
			"unique": false,
			"options": {
				"min": 11,
				"max": 11,
				"pattern": ""
			}
		}`), edit_student_id)
		collection.Schema.AddField(edit_student_id)

		return dao.SaveCollection(collection)
	})
}
