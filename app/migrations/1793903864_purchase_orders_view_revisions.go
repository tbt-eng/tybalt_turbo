package migrations

import (
	"strings"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// A revised purchase order returns to Unapproved but keeps its po_number. It
// was visible to everyone while Active, so hiding it during the revision makes
// it look deleted. Any signed-in user can now see it, until it is rejected. A
// rejected revision stays limited to the people involved, because the collection
// API can't hide its rejection_reason from filters. The /visible routes still
// show it, with the reason blanked.
const purchaseOrdersRevisionVisibleClause = `(status = "Unapproved" && po_number != "" && legacy_manual_entry = false && rejected = "" && @request.auth.id != "") || `

// addPurchaseOrdersRevisionClause puts the clause inside the opening parenthesis
// of a rule when it has one, because reverting the purchase order groups
// migration expects that parenthesis and its own clause at the ends of the rule.
func addPurchaseOrdersRevisionClause(rule string) string {
	if rest, ok := strings.CutPrefix(rule, "(\n"); ok {
		return "(\n" + purchaseOrdersRevisionVisibleClause + rest
	}
	return purchaseOrdersRevisionVisibleClause + rule
}

func removePurchaseOrdersRevisionClause(rule string) string {
	return strings.Replace(rule, purchaseOrdersRevisionVisibleClause, "", 1)
}

func changePurchaseOrdersReadRules(app core.App, change func(string) string) error {
	collection, err := app.FindCollectionByNameOrId("purchase_orders")
	if err != nil {
		return err
	}
	for _, rule := range []**string{&collection.ListRule, &collection.ViewRule} {
		if *rule != nil {
			changed := change(**rule)
			*rule = &changed
		}
	}
	return app.Save(collection)
}

func init() {
	m.Register(func(app core.App) error {
		return changePurchaseOrdersReadRules(app, addPurchaseOrdersRevisionClause)
	}, func(app core.App) error {
		return changePurchaseOrdersReadRules(app, removePurchaseOrdersRevisionClause)
	})
}
