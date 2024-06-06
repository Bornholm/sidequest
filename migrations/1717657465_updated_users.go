package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("_pb_users_auth_")
		if err != nil {
			return err
		}

		collection.ListRule = types.Pointer("id = @request.auth.id")

		collection.ViewRule = types.Pointer("id = @request.auth.id")

		collection.UpdateRule = types.Pointer("id = @request.auth.id")

		collection.DeleteRule = types.Pointer("id = @request.auth.id")

		// remove
		collection.Schema.RemoveField("u4dyvd1d")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("_pb_users_auth_")
		if err != nil {
			return err
		}

		collection.ListRule = types.Pointer("id = @request.auth.id && enabled = true")

		collection.ViewRule = types.Pointer("id = @request.auth.id && enabled = true")

		collection.UpdateRule = types.Pointer("id = @request.auth.id && enabled = true")

		collection.DeleteRule = types.Pointer("id = @request.auth.id && enabled = true")

		// add
		del_enabled := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "u4dyvd1d",
			"name": "enabled",
			"type": "bool",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {}
		}`), del_enabled); err != nil {
			return err
		}
		collection.Schema.AddField(del_enabled)

		return dao.SaveCollection(collection)
	})
}
