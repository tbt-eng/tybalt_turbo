package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// Make clients.business_development_lead optional in the schema.
//
// The business does not want to record a lead for now, so the UI no longer
// shows the field. A schema `required` flag is checked on every save, so it
// would reject any save of a client that has no lead, which the UI can no
// longer set. Existing values stay as they are.
//
// An optional relation lets PocketBase clear the reference when the referenced
// user is deleted, where a required one refused the delete.
//
// The down migration restores the flag. PocketBase doesn't recheck existing
// rows when a collection's schema is saved, so it succeeds even with blank
// values present, but those clients can no longer be saved afterwards.
func init() {
	m.Register(func(app core.App) error {
		return setClientsBusdevLeadRequired(app, false)
	}, func(app core.App) error {
		return setClientsBusdevLeadRequired(app, true)
	})
}

func setClientsBusdevLeadRequired(app core.App, required bool) error {
	clients, err := app.FindCollectionByNameOrId("clients")
	if err != nil {
		return err
	}
	clients.Fields.GetByName("business_development_lead").(*core.RelationField).Required = required
	return app.Save(clients)
}
