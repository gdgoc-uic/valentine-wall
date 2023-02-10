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
			"id": "0z6clug00f33eli",
			"created": "2023-02-10 13:38:36.519Z",
			"updated": "2023-02-10 13:38:36.519Z",
			"name": "message_replies_count",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "bx4yf3hr",
					"name": "message",
					"type": "relation",
					"required": true,
					"unique": true,
					"options": {
						"maxSelect": 1,
						"collectionId": "caqiysan7yf0wve",
						"cascadeDelete": true
					}
				},
				{
					"system": false,
					"id": "bsubls6d",
					"name": "count",
					"type": "number",
					"required": false,
					"unique": false,
					"options": {
						"min": 0,
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

		collection, err := dao.FindCollectionByNameOrId("0z6clug00f33eli")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
