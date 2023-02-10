package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("55bmnoqmzjhnng0")
		if err != nil {
			return err
		}

		collection.ListRule = types.Pointer("@request.auth.user.id = user.id")

		collection.ViewRule = types.Pointer("@request.auth.user.id = user.id")

		collection.DeleteRule = types.Pointer("@request.auth.user.id = user.id")

		options := map[string]any{}
		json.Unmarshal([]byte(`{
			"allowEmailAuth": false,
			"allowOAuth2Auth": true,
			"allowUsernameAuth": false,
			"exceptEmailDomains": [],
			"manageRule": "@request.auth.user.id = user.id",
			"minPasswordLength": 8,
			"onlyEmailDomains": [],
			"requireEmail": false
		}`), &options)
		collection.SetOptions(options)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("55bmnoqmzjhnng0")
		if err != nil {
			return err
		}

		collection.ListRule = nil

		collection.ViewRule = nil

		collection.DeleteRule = nil

		options := map[string]any{}
		json.Unmarshal([]byte(`{
			"allowEmailAuth": false,
			"allowOAuth2Auth": true,
			"allowUsernameAuth": false,
			"exceptEmailDomains": [],
			"manageRule": null,
			"minPasswordLength": 8,
			"onlyEmailDomains": [],
			"requireEmail": false
		}`), &options)
		collection.SetOptions(options)

		return dao.SaveCollection(collection)
	})
}
