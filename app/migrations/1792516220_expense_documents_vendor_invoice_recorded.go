package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// Records who last set or changed an expense document's vendor invoice details,
// and when, so the people working with a shared invoice can see who set it.
// This is the latest change only, not a history.
func init() {
	m.Register(func(app core.App) error {
		documents, err := app.FindCollectionByNameOrId("expense_documents")
		if err != nil {
			return err
		}
		users, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}
		documents.Fields.Add(
			&core.RelationField{Name: "vendor_invoice_recorded_by", CollectionId: users.Id, MaxSelect: 1},
			&core.DateField{Name: "vendor_invoice_recorded"},
		)
		return app.Save(documents)
	}, func(app core.App) error {
		documents, err := app.FindCollectionByNameOrId("expense_documents")
		if err != nil {
			return err
		}
		documents.Fields.RemoveByName("vendor_invoice_recorded_by")
		documents.Fields.RemoveByName("vendor_invoice_recorded")
		return app.Save(documents)
	})
}
