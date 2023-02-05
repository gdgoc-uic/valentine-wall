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

		collection, err := dao.FindCollectionByNameOrId("yhv9suo8ru0esf0")
		if err != nil {
			return err
		}

		// update
		edit_uid := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "sr6nfuzl",
			"name": "uid",
			"type": "text",
			"required": true,
			"unique": true,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), edit_uid)
		collection.Schema.AddField(edit_uid)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("yhv9suo8ru0esf0")
		if err != nil {
			return err
		}

		// update
		edit_uid := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "sr6nfuzl",
			"name": "abbreviation",
			"type": "text",
			"required": true,
			"unique": true,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), edit_uid)
		collection.Schema.AddField(edit_uid)

		return dao.SaveCollection(collection)
	})
}
