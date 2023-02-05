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
			"id": "rocfs9e910vqq5g",
			"created": "2023-02-05 05:35:28.633Z",
			"updated": "2023-02-05 05:35:28.633Z",
			"name": "virtual_transactions",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "gmnhl7zd",
					"name": "wallet",
					"type": "relation",
					"required": true,
					"unique": false,
					"options": {
						"maxSelect": 1,
						"collectionId": "rs07r3dxff0hxbz",
						"cascadeDelete": false
					}
				},
				{
					"system": false,
					"id": "nlkxzwwf",
					"name": "amount",
					"type": "number",
					"required": true,
					"unique": false,
					"options": {
						"min": null,
						"max": null
					}
				},
				{
					"system": false,
					"id": "c64fp4e6",
					"name": "description",
					"type": "text",
					"required": true,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
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

		collection, err := dao.FindCollectionByNameOrId("rocfs9e910vqq5g")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
