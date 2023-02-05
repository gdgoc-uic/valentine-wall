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
			"id": "ocpdx07v34h97tx",
			"created": "2023-01-29 11:33:55.096Z",
			"updated": "2023-01-29 11:33:55.096Z",
			"name": "rankings",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "dvavc1rl",
					"name": "recipient_id",
					"type": "text",
					"required": true,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				},
				{
					"system": false,
					"id": "ylooys31",
					"name": "department",
					"type": "text",
					"required": true,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				},
				{
					"system": false,
					"id": "lgeckfbq",
					"name": "sex",
					"type": "text",
					"required": true,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				},
				{
					"system": false,
					"id": "c5zgfdxg",
					"name": "total_coins",
					"type": "text",
					"required": true,
					"unique": false,
					"options": {
						"min": 0,
						"max": null,
						"pattern": ""
					}
				}
			],
			"listRule": "",
			"viewRule": "",
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

		collection, err := dao.FindCollectionByNameOrId("ocpdx07v34h97tx")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
