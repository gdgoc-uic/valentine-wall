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
			"id": "z3ucff2sh922dhz",
			"created": "2023-01-21 02:37:51.004Z",
			"updated": "2023-01-21 02:37:51.004Z",
			"name": "gifts",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "m3d9nf2a",
					"name": "uid",
					"type": "text",
					"required": false,
					"unique": true,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				},
				{
					"system": false,
					"id": "dbx2dmds",
					"name": "label",
					"type": "text",
					"required": false,
					"unique": true,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				},
				{
					"system": false,
					"id": "2or8do7n",
					"name": "price",
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

		collection, err := dao.FindCollectionByNameOrId("z3ucff2sh922dhz")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
