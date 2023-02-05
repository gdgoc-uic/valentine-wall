package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `{
			"id": "35mnuyxwxc8xvs6",
			"created": "2023-01-21 02:52:37.027Z",
			"updated": "2023-01-21 02:52:37.027Z",
			"name": "message_replies",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "f7vpglss",
					"name": "message",
					"type": "relation",
					"required": true,
					"unique": false,
					"options": {
						"maxSelect": 1,
						"collectionId": "caqiysan7yf0wve",
						"cascadeDelete": true
					}
				},
				{
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
				}
			],
			"listRule": "@request.auth.user.id = message.user.id || @request.auth.user.id = user.id",
			"viewRule": "@request.auth.user.id = message.user.id || @request.auth.user.id = user.id",
			"createRule": "@request.auth.user.id = user.id",
			"updateRule": "@request.auth.user.id = user.id",
			"deleteRule": "@request.auth.user.id = user.id",
			"options": {}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("35mnuyxwxc8xvs6")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
