package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("35mnuyxwxc8xvs6")
		if err != nil {
			return err
		}

		collection.ListRule = types.Pointer("@request.auth.user.id = message.user.id || @request.auth.user.id = sender.id")

		collection.ViewRule = types.Pointer("@request.auth.user.id = message.user.id || @request.auth.user.id = sender.id")

		collection.CreateRule = types.Pointer("@request.auth.user.id = sender.id")

		collection.UpdateRule = types.Pointer("@request.auth.user.id = sender.id")

		collection.DeleteRule = types.Pointer("@request.auth.user.id = sender.id")

		// add
		new_content := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "zphlknfx",
			"name": "content",
			"type": "text",
			"required": true,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_content)
		collection.Schema.AddField(new_content)

		// update
		edit_sender := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "yzdg2alv",
			"name": "sender",
			"type": "relation",
			"required": true,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"collectionId": "_pb_users_auth_",
				"cascadeDelete": true
			}
		}`), edit_sender)
		collection.Schema.AddField(edit_sender)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("35mnuyxwxc8xvs6")
		if err != nil {
			return err
		}

		collection.ListRule = nil

		collection.ViewRule = nil

		collection.CreateRule = nil

		collection.UpdateRule = nil

		collection.DeleteRule = nil

		// remove
		collection.Schema.RemoveField("zphlknfx")

		// update
		edit_sender := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "yzdg2alv",
			"name": "user",
			"type": "relation",
			"required": false,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"collectionId": "_pb_users_auth_",
				"cascadeDelete": true
			}
		}`), edit_sender)
		collection.Schema.AddField(edit_sender)

		return dao.SaveCollection(collection)
	})
}
