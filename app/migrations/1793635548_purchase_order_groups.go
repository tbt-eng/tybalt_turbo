package migrations

import (
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

const (
	purchaseOrderGroupsCollectionID = "pbc_po_groups"
	purchaseOrderGroupIndexName     = "idx_purchase_orders_po_group"
	superPOClaimID                  = "superpo00000001"

	// super_po and admin holders see every purchase order in a group, as the
	// purchase order visibility routes already allow. The rule also sends
	// them realtime events for the group page.
	groupPurchaseOrdersRule = "po_group != '' &&\n(\n  @request.auth.user_claims_via_uid.cid.name ?= 'super_po' ||\n  @request.auth.user_claims_via_uid.cid.name ?= 'admin'\n)"
)

// purchaseOrderGroupEmails tell each approver once about a whole group instead
// of once per purchase order, which for a group of every staff member would be
// hundreds of emails.
var purchaseOrderGroupEmails = []projectWorkflowEmail{
	{"pogroupapprove0", "po_group_approval_required", "Asks the approver of a purchase order group to review its purchase orders",
		"Approval needed: {{.Count}} purchase orders",
		`Hello {{.RecipientName}},

{{.UserName}} has created {{.Count}} purchase orders, one for each of several staff members, and chose you as the approver of all of them.

Description: {{.Description}}
Vendor: {{.VendorName}}
Total of all {{.Count}}: {{.Total}}

Each purchase order is approved on its own, which gives it a PO number. They are in your Pending My Approval list:

{{.ActionURL}}`},
	{"pogroup2ndappr0", "po_group_second_approval_required", "Asks the second approver of a purchase order group to give second approval",
		"Second approval needed: purchase orders for {{.Count}} staff",
		`Hello {{.RecipientName}},

{{.CreatorName}} created {{.Count}} purchase orders, one for each of several staff members, and chose you as their second approver. {{.UserName}} has started approving them, and those over the second approval threshold need your approval too.

Description: {{.Description}}
Vendor: {{.VendorName}}

This is the only email you get for the group. Each purchase order appears in your Pending My Approval list once {{.UserName}} approves it:

{{.ActionURL}}`},
}

// widenRule returns rule with clause added as another way to pass it. The
// existing rule is wrapped whole rather than edited, so this works however
// the rule is written, including one that ends in a // comment. narrowRule
// undoes it.
func widenRule(rule *string, clause string) *string {
	if rule == nil {
		return pointerString(clause)
	}
	if strings.TrimSpace(*rule) == "" {
		return rule
	}
	return pointerString("(\n" + *rule + "\n) || (\n" + clause + "\n)")
}

// narrowRule removes a clause that widenRule added. A rule that has been
// rewritten since is left as it is.
func narrowRule(rule *string, clause string) *string {
	if rule == nil || *rule == clause {
		return nil
	}
	suffix := "\n) || (\n" + clause + "\n)"
	if strings.HasPrefix(*rule, "(\n") && strings.HasSuffix(*rule, suffix) {
		return pointerString(strings.TrimSuffix(strings.TrimPrefix(*rule, "(\n"), suffix))
	}
	return rule
}

// Add purchase order groups and the super_po claim (issue #168).
//
// A super_po holder fills out one purchase order as a template and uploads a
// CSV of staff. The route creates one ordinary Unapproved purchase order per
// row, owned by that staff member, and links them all to one
// purchase_order_groups record. The group records who created the purchase
// orders, since each one's uid is the staff member rather than the creator.
//
// Groups are created and read only through the purchase order group routes,
// so the collection has no API rules. purchase_orders.po_group is set only by
// the same route; ProcessPurchaseOrder rejects it on collection API creates
// and updates. super_po and admin holders can read every purchase order in a
// group. The approver gets one po_group_approval_required
// email per group rather than a po_approval_required email per purchase order,
// and the second approver one po_group_second_approval_required email.
func init() {
	m.Register(func(app core.App) error {
		users, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}

		groups := core.NewBaseCollection("purchase_order_groups")
		groups.Id = purchaseOrderGroupsCollectionID
		groups.Fields.Add(
			&core.RelationField{
				Name:          "creator",
				CollectionId:  users.Id,
				MaxSelect:     1,
				Required:      true,
				CascadeDelete: false,
			},
			// The template's description, blank when each purchase order's
			// description came from the upload.
			&core.TextField{Name: "description"},
			&core.AutodateField{Name: "created", OnCreate: true},
			&core.AutodateField{Name: "updated", OnCreate: true, OnUpdate: true},
		)
		if err := app.Save(groups); err != nil {
			return err
		}

		purchaseOrders, err := app.FindCollectionByNameOrId("purchase_orders")
		if err != nil {
			return err
		}
		purchaseOrders.Fields.Add(&core.RelationField{
			Name:          "po_group",
			CollectionId:  groups.Id,
			MaxSelect:     1,
			CascadeDelete: false,
		})
		purchaseOrders.AddIndex(purchaseOrderGroupIndexName, false, "`po_group`", "`po_group` != ''")
		purchaseOrders.ListRule = widenRule(purchaseOrders.ListRule, groupPurchaseOrdersRule)
		purchaseOrders.ViewRule = widenRule(purchaseOrders.ViewRule, groupPurchaseOrdersRule)
		if err := app.Save(purchaseOrders); err != nil {
			return err
		}

		claims, err := app.FindCollectionByNameOrId("claims")
		if err != nil {
			return err
		}
		claim := core.NewRecord(claims)
		claim.Id = superPOClaimID
		claim.Set("name", "super_po")
		claim.Set("description", "Create purchase orders for other staff from a template and a CSV upload")
		if err := app.Save(claim); err != nil {
			return err
		}

		return replaceProjectWorkflowEmails(app, nil, purchaseOrderGroupEmails)
	}, func(app core.App) error {
		if err := replaceProjectWorkflowEmails(app, purchaseOrderGroupEmails, nil); err != nil {
			return err
		}

		// A user_claims record requires its claim, so the claim's holders
		// lose it first.
		holders, err := app.FindAllRecords("user_claims", dbx.HashExp{"cid": superPOClaimID})
		if err != nil {
			return err
		}
		for _, holder := range holders {
			if err := app.Delete(holder); err != nil {
				return err
			}
		}
		if claim, err := app.FindRecordById("claims", superPOClaimID); err == nil {
			if err := app.Delete(claim); err != nil {
				return err
			}
		}

		purchaseOrders, err := app.FindCollectionByNameOrId("purchase_orders")
		if err != nil {
			return err
		}
		purchaseOrders.ListRule = narrowRule(purchaseOrders.ListRule, groupPurchaseOrdersRule)
		purchaseOrders.ViewRule = narrowRule(purchaseOrders.ViewRule, groupPurchaseOrdersRule)
		purchaseOrders.RemoveIndex(purchaseOrderGroupIndexName)
		purchaseOrders.Fields.RemoveByName("po_group")
		if err := app.Save(purchaseOrders); err != nil {
			return err
		}

		groups, err := app.FindCollectionByNameOrId(purchaseOrderGroupsCollectionID)
		if err != nil {
			return nil
		}
		return app.Delete(groups)
	})
}
