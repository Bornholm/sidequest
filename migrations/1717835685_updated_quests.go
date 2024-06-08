package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("502ykv83hmmw4xu")
		if err != nil {
			return err
		}

		// add
		new_solution := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "792hlsfg",
			"name": "solution",
			"type": "text",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_solution); err != nil {
			return err
		}
		collection.Schema.AddField(new_solution)

		// add
		new_characters := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "ncbdv2kx",
			"name": "characters",
			"type": "text",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_characters); err != nil {
			return err
		}
		collection.Schema.AddField(new_characters)

		// add
		new_clues := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "t7mjwmgo",
			"name": "clues",
			"type": "text",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_clues); err != nil {
			return err
		}
		collection.Schema.AddField(new_clues)

		// add
		new_description := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "3doqpyn3",
			"name": "description",
			"type": "text",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_description); err != nil {
			return err
		}
		collection.Schema.AddField(new_description)

		// add
		new_author := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "1nnhamnm",
			"name": "author",
			"type": "relation",
			"required": true,
			"presentable": false,
			"unique": false,
			"options": {
				"collectionId": "_pb_users_auth_",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": 1,
				"displayFields": null
			}
		}`), new_author); err != nil {
			return err
		}
		collection.Schema.AddField(new_author)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("502ykv83hmmw4xu")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("792hlsfg")

		// remove
		collection.Schema.RemoveField("ncbdv2kx")

		// remove
		collection.Schema.RemoveField("t7mjwmgo")

		// remove
		collection.Schema.RemoveField("3doqpyn3")

		// remove
		collection.Schema.RemoveField("1nnhamnm")

		return dao.SaveCollection(collection)
	})
}
