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
		jsonData := `[
			{
				"id": "z3ucff2sh922dhz",
				"created": "2023-01-21 02:37:51.004Z",
				"updated": "2023-02-06 09:35:28.839Z",
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
					},
					{
						"system": false,
						"id": "n9lck8lq",
						"name": "is_remittable",
						"type": "bool",
						"required": false,
						"unique": false,
						"options": {}
					}
				],
				"listRule": "",
				"viewRule": "",
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "yhv9suo8ru0esf0",
				"created": "2023-01-21 02:38:51.020Z",
				"updated": "2023-02-06 09:35:28.730Z",
				"name": "college_departments",
				"type": "base",
				"system": false,
				"schema": [
					{
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
					},
					{
						"system": false,
						"id": "yae8npoj",
						"name": "label",
						"type": "text",
						"required": true,
						"unique": true,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					}
				],
				"listRule": "",
				"viewRule": "",
				"createRule": "",
				"updateRule": "",
				"deleteRule": "",
				"options": {}
			},
			{
				"id": "caqiysan7yf0wve",
				"created": "2023-01-21 02:43:38.501Z",
				"updated": "2023-02-10 13:46:20.991Z",
				"name": "messages",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "orteybh2",
						"name": "content",
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
						"id": "nfu7q7wd",
						"name": "user",
						"type": "relation",
						"required": true,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"collectionId": "px00yjig95x0mcw",
							"cascadeDelete": true
						}
					},
					{
						"system": false,
						"id": "uuyjuypn",
						"name": "gifts",
						"type": "relation",
						"required": false,
						"unique": false,
						"options": {
							"maxSelect": 3,
							"collectionId": "z3ucff2sh922dhz",
							"cascadeDelete": false
						}
					},
					{
						"system": false,
						"id": "8v4wh5mo",
						"name": "deleted",
						"type": "date",
						"required": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					},
					{
						"system": false,
						"id": "idbaffkv",
						"name": "recipient",
						"type": "text",
						"required": true,
						"unique": false,
						"options": {
							"min": 11,
							"max": 12,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "xjc2zwan",
						"name": "replies_count",
						"type": "number",
						"required": false,
						"unique": false,
						"options": {
							"min": 0,
							"max": null
						}
					}
				],
				"listRule": "",
				"viewRule": "",
				"createRule": "@request.auth.id != \"\"",
				"updateRule": "@request.auth.id = user.id",
				"deleteRule": "@request.auth.id = user.id",
				"options": {}
			},
			{
				"id": "55bmnoqmzjhnng0",
				"created": "2023-01-21 02:47:07.014Z",
				"updated": "2023-02-10 13:56:17.657Z",
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
					},
					{
						"system": false,
						"id": "v9kwnaqi",
						"name": "provider",
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
				"listRule": "@request.auth.user.id = user.id",
				"viewRule": "@request.auth.user.id = user.id",
				"createRule": null,
				"updateRule": null,
				"deleteRule": "@request.auth.user.id = user.id",
				"options": {
					"allowEmailAuth": false,
					"allowOAuth2Auth": true,
					"allowUsernameAuth": false,
					"exceptEmailDomains": [],
					"manageRule": "@request.auth.user.id = user.id",
					"minPasswordLength": 8,
					"onlyEmailDomains": [],
					"requireEmail": false
				}
			},
			{
				"id": "35mnuyxwxc8xvs6",
				"created": "2023-01-21 02:52:37.027Z",
				"updated": "2023-02-11 04:52:52.643Z",
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
						"name": "sender",
						"type": "relation",
						"required": true,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"collectionId": "px00yjig95x0mcw",
							"cascadeDelete": true
						}
					},
					{
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
					},
					{
						"system": false,
						"id": "bgot1d9q",
						"name": "liked",
						"type": "bool",
						"required": false,
						"unique": false,
						"options": {}
					}
				],
				"listRule": "@request.auth.details = message.user.id || @request.auth.details = sender.id",
				"viewRule": "@request.auth.id = message.user.id || @request.auth.id = sender.id",
				"createRule": "@request.auth.details = @request.data.sender",
				"updateRule": "@request.auth.details = sender.id",
				"deleteRule": "@request.auth.details = sender.id",
				"options": {}
			},
			{
				"id": "ocpdx07v34h97tx",
				"created": "2023-01-29 11:33:55.096Z",
				"updated": "2023-02-06 09:35:28.868Z",
				"name": "rankings",
				"type": "base",
				"system": false,
				"schema": [
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
					},
					{
						"system": false,
						"id": "optuzp59",
						"name": "recipient",
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
						"id": "j0a0c0dg",
						"name": "college_department",
						"type": "relation",
						"required": true,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"collectionId": "yhv9suo8ru0esf0",
							"cascadeDelete": false
						}
					},
					{
						"system": false,
						"id": "tnebrg5b",
						"name": "sex",
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
				"listRule": "",
				"viewRule": "",
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "px00yjig95x0mcw",
				"created": "2023-02-05 05:00:28.767Z",
				"updated": "2023-02-09 12:55:43.638Z",
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
					},
					{
						"system": false,
						"id": "2gizlajc",
						"name": "sex",
						"type": "select",
						"required": true,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"values": [
								"unknown",
								"male",
								"female"
							]
						}
					},
					{
						"system": false,
						"id": "oaqps5fh",
						"name": "student_id",
						"type": "text",
						"required": true,
						"unique": false,
						"options": {
							"min": 11,
							"max": 12,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "erbuqwm0",
						"name": "last_active",
						"type": "date",
						"required": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					},
					{
						"system": false,
						"id": "ukhwq1aw",
						"name": "user",
						"type": "relation",
						"required": true,
						"unique": true,
						"options": {
							"maxSelect": 1,
							"collectionId": "_pb_users_auth_",
							"cascadeDelete": true
						}
					}
				],
				"listRule": "",
				"viewRule": "",
				"createRule": "@request.auth.id = user.id",
				"updateRule": "@request.auth.id = user.id",
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "rs07r3dxff0hxbz",
				"created": "2023-02-05 05:32:43.476Z",
				"updated": "2023-02-10 12:17:48.206Z",
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
						"unique": true,
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
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null
						}
					}
				],
				"listRule": "@request.auth.id != \"\"",
				"viewRule": "@request.auth.id = user.id",
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "rocfs9e910vqq5g",
				"created": "2023-02-05 05:35:28.633Z",
				"updated": "2023-02-10 02:46:12.051Z",
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
				"listRule": "@request.auth.id = wallet.user.id",
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "_pb_users_auth_",
				"created": "2023-02-10 11:44:13.368Z",
				"updated": "2023-02-10 11:44:13.372Z",
				"name": "users",
				"type": "auth",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "users_avatar",
						"name": "avatar",
						"type": "file",
						"required": false,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"maxSize": 5242880,
							"mimeTypes": [
								"image/jpg",
								"image/jpeg",
								"image/png",
								"image/svg+xml",
								"image/gif",
								"image/webp"
							],
							"thumbs": null
						}
					},
					{
						"system": false,
						"id": "lz0uypgf",
						"name": "details",
						"type": "relation",
						"required": false,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"collectionId": "px00yjig95x0mcw",
							"cascadeDelete": true
						}
					}
				],
				"listRule": "id = @request.auth.id",
				"viewRule": "id = @request.auth.id",
				"createRule": "",
				"updateRule": "id = @request.auth.id",
				"deleteRule": "id = @request.auth.id",
				"options": {
					"allowEmailAuth": false,
					"allowOAuth2Auth": true,
					"allowUsernameAuth": false,
					"exceptEmailDomains": null,
					"manageRule": null,
					"minPasswordLength": 8,
					"onlyEmailDomains": null,
					"requireEmail": false
				}
			}
		]`

		collections := []*models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collections); err != nil {
			return err
		}

		return daos.New(db).ImportCollections(collections, true, nil)
	}, func(db dbx.Builder) error {
		return nil
	})
}
