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

		collection, err := dao.FindCollectionByNameOrId("35mnuyxwxc8xvs6")
		if err != nil {
			return err
		}

		// add
		new_liked := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "bgot1d9q",
			"name": "liked",
			"type": "bool",
			"required": false,
			"unique": false,
			"options": {}
		}`), new_liked)
		collection.Schema.AddField(new_liked)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("35mnuyxwxc8xvs6")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("bgot1d9q")

		return dao.SaveCollection(collection)
	})
}
