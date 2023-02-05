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
			"id": "55bmnoqmzjhnng0",
			"created": "2023-01-21 02:47:07.014Z",
			"updated": "2023-01-21 02:47:07.014Z",
			"name": "user_connections",
			"type": "auth",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "gzjhpahn",
					"name": "user",
					"type": "relation",
					"required": true,
					"unique": false,
					"options": {
						"maxSelect": 1,
						"collectionId": "_pb_users_auth_",
						"cascadeDelete": true
					}
				}
			],
			"listRule": null,
			"viewRule": null,
			"createRule": null,
			"updateRule": null,
			"deleteRule": null,
			"options": {
				"allowEmailAuth": false,
				"allowOAuth2Auth": true,
				"allowUsernameAuth": false,
				"exceptEmailDomains": [],
				"manageRule": null,
				"minPasswordLength": 8,
				"onlyEmailDomains": [],
				"requireEmail": false
			}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("55bmnoqmzjhnng0")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
