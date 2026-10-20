package migrations

import (
	"database/sql"
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

const (
	invoicingClaimName = "invoicing"
	invoicingClaimDesc = "Can record client invoices and link billed-back expenses and purchase orders to them"

	invoicingRule       = `@request.auth.user_claims_via_uid.cid.name ?= "invoicing"`
	invoicingReportRule = `@request.auth.user_claims_via_uid.cid.name ?= "invoicing" || @request.auth.user_claims_via_uid.cid.name ?= "report"`
)

// invoices records the invoices we send clients, and invoice_lines the costs
// billed back on them. Issuing goes through POST /api/invoices/{id}/issue; the
// hooks keep status and the server-set fields out of the collection API. See
// descriptions/invoices.md.
func init() {
	m.Register(func(app core.App) error {
		if err := ensureInvoicingClaim(app); err != nil {
			return err
		}
		clients, err := app.FindCollectionByNameOrId("clients")
		if err != nil {
			return err
		}
		users, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}
		expenses, err := app.FindCollectionByNameOrId("expenses")
		if err != nil {
			return err
		}
		purchaseOrders, err := app.FindCollectionByNameOrId("purchase_orders")
		if err != nil {
			return err
		}

		invoices := core.NewBaseCollection("invoices")
		invoices.ListRule = types.Pointer(invoicingReportRule)
		invoices.ViewRule = types.Pointer(invoicingReportRule)
		invoices.CreateRule = types.Pointer(invoicingRule)
		invoices.UpdateRule = types.Pointer(invoicingRule + ` && status = "Draft"`)
		invoices.DeleteRule = types.Pointer(invoicingRule + ` && status = "Draft"`)
		invoices.Fields.Add(
			&core.RelationField{Name: "client", CollectionId: clients.Id, Required: true, MaxSelect: 1},
			// YYMMNNN with a real month. Empty only on a draft.
			&core.TextField{Name: "number", Pattern: `^(\d{2}(0[1-9]|1[0-2])\d{3})?$`},
			// YYYY-MM-DD. Empty only on a draft.
			&core.TextField{Name: "date", Pattern: `^(\d{4}-\d{2}-\d{2})?$`},
			&core.NumberField{Name: "total", Min: types.Pointer(0.0)},
			&core.SelectField{Name: "status", Required: true, MaxSelect: 1, Values: []string{"Draft", "Issued"}},
			&core.RelationField{Name: "creator", CollectionId: users.Id, MaxSelect: 1},
			&core.RelationField{Name: "issued_by", CollectionId: users.Id, MaxSelect: 1},
			&core.TextField{Name: "issued_at"},
			&core.AutodateField{Name: "created", OnCreate: true},
			&core.AutodateField{Name: "updated", OnCreate: true, OnUpdate: true},
		)
		invoices.AddIndex("idx_invoices_number", true, "`number`", "`number` != ''")
		invoices.AddIndex("idx_invoices_client", false, "`client`", "")
		if err := app.Save(invoices); err != nil {
			return err
		}

		lines := core.NewBaseCollection("invoice_lines")
		lines.ListRule = types.Pointer(invoicingReportRule)
		lines.ViewRule = types.Pointer(invoicingReportRule)
		lines.CreateRule = types.Pointer(invoicingRule)
		lines.DeleteRule = types.Pointer(invoicingRule)
		lines.Fields.Add(
			&core.RelationField{Name: "invoice", CollectionId: invoices.Id, Required: true, MaxSelect: 1},
			&core.SelectField{Name: "kind", Required: true, MaxSelect: 1, Values: []string{"expense", "po"}},
			&core.RelationField{Name: "expense", CollectionId: expenses.Id, MaxSelect: 1},
			&core.RelationField{Name: "purchase_order", CollectionId: purchaseOrders.Id, MaxSelect: 1},
			&core.NumberField{Name: "amount", Required: true, Min: types.Pointer(0.01)},
			&core.RelationField{Name: "creator", CollectionId: users.Id, MaxSelect: 1},
			&core.AutodateField{Name: "created", OnCreate: true},
			&core.AutodateField{Name: "updated", OnCreate: true, OnUpdate: true},
		)
		lines.AddIndex("idx_invoice_lines_invoice", false, "`invoice`", "")
		// An expense or PO can be billed on only one line.
		lines.AddIndex("idx_invoice_lines_expense", true, "`expense`", "`expense` != ''")
		lines.AddIndex("idx_invoice_lines_purchase_order", true, "`purchase_order`", "`purchase_order` != ''")
		return app.Save(lines)
	}, func(app core.App) error {
		for _, name := range []string{"invoice_lines", "invoices"} {
			c, err := app.FindCollectionByNameOrId(name)
			if err != nil {
				return err
			}
			if err := app.Delete(c); err != nil {
				return err
			}
		}
		claim, err := app.FindFirstRecordByFilter("claims", "name={:name}", dbx.Params{"name": invoicingClaimName})
		if err != nil {
			return nil
		}
		return app.Delete(claim)
	})
}

func ensureInvoicingClaim(app core.App) error {
	existing, err := app.FindFirstRecordByFilter("claims", "name={:name}", dbx.Params{"name": invoicingClaimName})
	if err == nil && existing != nil {
		return nil
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	collection, err := app.FindCollectionByNameOrId("claims")
	if err != nil {
		return err
	}
	record := core.NewRecord(collection)
	record.Set("id", "invoicingclm001")
	record.Set("name", invoicingClaimName)
	record.Set("description", invoicingClaimDesc)
	return app.Save(record)
}
