package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// A revised purchase order returns to Unapproved but keeps its po_number.
// The creator may cancel it once it is rejected, but must not delete it, so
// that a number once issued never disappears. Only unnumbered purchase orders
// may be deleted.
const (
	purchaseOrdersOldDeleteRule = "@request.auth.id = uid && status = 'Unapproved'"
	purchaseOrdersNewDeleteRule = "@request.auth.id = uid && status = 'Unapproved' && po_number = ''"
)

func init() {
	m.Register(func(app core.App) error {
		return updateRule(app, "purchase_orders", "deleteRule", purchaseOrdersNewDeleteRule)
	}, func(app core.App) error {
		return updateRule(app, "purchase_orders", "deleteRule", purchaseOrdersOldDeleteRule)
	})
}
