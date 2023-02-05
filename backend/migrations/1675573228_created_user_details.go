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
			"id": "px00yjig95x0mcw",
			"created": "2023-02-05 05:00:28.767Z",
			"updated": "2023-02-05 05:00:28.767Z",
			"name": "user_details",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "s2f534mf",
					"name": "college_department",
					"type": "relation",
					"required": true,
					"unique": false,
					"options": {
						"maxSelect": 1,
						"collectionId": "yhv9suo8ru0esf0",
						"cascadeDelete": true
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

		collection, err := dao.FindCollectionByNameOrId("px00yjig95x0mcw")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
