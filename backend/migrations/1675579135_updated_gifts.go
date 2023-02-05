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

		collection, err := dao.FindCollectionByNameOrId("z3ucff2sh922dhz")
		if err != nil {
			return err
		}

		// add
		new_is_remittable := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "n9lck8lq",
			"name": "is_remittable",
			"type": "bool",
			"required": false,
			"unique": false,
			"options": {}
		}`), new_is_remittable)
		collection.Schema.AddField(new_is_remittable)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("z3ucff2sh922dhz")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("n9lck8lq")

		return dao.SaveCollection(collection)
	})
}
