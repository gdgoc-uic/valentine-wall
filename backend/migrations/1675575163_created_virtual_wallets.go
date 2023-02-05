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
			"id": "rs07r3dxff0hxbz",
			"created": "2023-02-05 05:32:43.476Z",
			"updated": "2023-02-05 05:32:43.476Z",
			"name": "virtual_wallets",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "jx6fdcnf",
					"name": "user",
					"type": "relation",
					"required": true,
					"unique": false,
					"options": {
						"maxSelect": 1,
						"collectionId": "_pb_users_auth_",
						"cascadeDelete": true
					}
				},
				{
					"system": false,
					"id": "0abw1isr",
					"name": "balance",
					"type": "number",
					"required": true,
					"unique": false,
					"options": {
						"min": null,
						"max": null
					}
				}
			],
			"listRule": null,
			"viewRule": null,
			"createRule": null,
			"updateRule": null,
			"deleteRule": null,
			"options": {}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("rs07r3dxff0hxbz")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
