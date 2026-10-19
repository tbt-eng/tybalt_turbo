package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// Vendor invoice details live on the expense document that substantiates an
// expense. See the vendor invoice section of descriptions/expense_documents.md.
func init() {
	m.Register(func(app core.App) error {
		documents, err := app.FindCollectionByNameOrId("expense_documents")
		if err != nil {
			return err
		}
		vendors, err := app.FindCollectionByNameOrId("vendors")
		if err != nil {
			return err
		}
		currencies, err := app.FindCollectionByNameOrId("currencies")
		if err != nil {
			return err
		}

		documents.Fields.Add(
			&core.RelationField{Name: "vendor", CollectionId: vendors.Id, MaxSelect: 1},
			&core.TextField{Name: "vendor_invoice_number", Max: 100},
			// Lowercase letters and digits only, set by the server for the
			// duplicate check.
			&core.TextField{Name: "vendor_invoice_number_key", Max: 100},
			&core.TextField{Name: "vendor_invoice_date", Pattern: `^(\d{4}-\d{2}-\d{2})?$`},
			&core.NumberField{Name: "vendor_invoice_total"},
			&core.RelationField{Name: "vendor_invoice_currency", CollectionId: currencies.Id, MaxSelect: 1},
		)
		documents.AddIndex("idx_expense_documents_vendor_invoice", true, "`vendor`, `vendor_invoice_number_key`", "`vendor_invoice_number_key` != ''")
		return app.Save(documents)
	}, func(app core.App) error {
		documents, err := app.FindCollectionByNameOrId("expense_documents")
		if err != nil {
			return err
		}
		documents.RemoveIndex("idx_expense_documents_vendor_invoice")
		for _, name := range []string{"vendor", "vendor_invoice_number", "vendor_invoice_number_key", "vendor_invoice_date", "vendor_invoice_total", "vendor_invoice_currency"} {
			documents.Fields.RemoveByName(name)
		}
		return app.Save(documents)
	})
}
