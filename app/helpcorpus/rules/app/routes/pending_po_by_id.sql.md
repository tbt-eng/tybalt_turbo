<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/pending_po_by_id.sql","sha256":"3318bbfa2dd184fc31de051609fe39fbb407b470f8933324b073f7a91015d23c","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/pending_po_by_id.sql

## Pending PO lookup uses the PO visibility base

Lines 1-3

Asked as: which purchase orders can I find; PO not showing in lookup; missing PO search results; PO visibility rules

The pending purchase order lookup reads from a shared visibility base (__PO_VISIBILITY_BASE__), which is filled in elsewhere. Only POs that this base returns can be found. This file does not show which POs the base includes or excludes.

## Pending PO lookup returns PO identity, status and type

Lines 4-18

Asked as: PO details by id; PO number and status; look up purchase order fields; who created this PO

The pending purchase order lookup by id returns the PO's identity and descriptive fields: id, PO number (po_number), status, creator (uid), legacy_manual_entry, type, date, end date, frequency, division, description, payment type, attachment and group (po_group). The view it reads from decides which POs are visible. These fields come back only for a PO that passes the actionable filter.

## Pending PO lookup returns amounts, currency and spending figures

Lines 16-54

Asked as: PO total and currency; how much is left on PO; PO expenses spent so far; recurring PO remaining occurrences

The pending purchase order lookup returns the PO total, the currency details (code, symbol, icon, rate, rate date), and the expense figures. Those figures are committed_expenses_count, expenses_total, remaining_amount and print_max_amount. For a recurring PO it also returns the expected and remaining occurrences, and covered_within_project_budget.

## Pending PO lookup returns approval, rejection, cancel and close details

Lines 19-26

Asked as: who approved this PO; PO rejection or cancellation reason; who closed the purchase order; second approver of PO

The pending purchase order lookup returns the PO's approval fields (approver, approved, second_approver, second_approval, priority_second_approver, approval_total, approval_total_home). It also returns the rejection, cancellation and closure fields (rejector, rejected, rejection_reason, canceller, cancelled, closer, closed, closed_by_system). Each person is also returned by name.

## Pending PO lookup only returns actionable unapproved POs

Lines 77-80

Asked as: PO not returned as pending; why no record for PO; PO already approved so missing; PO not actionable right now

The pending purchase order lookup by id returns a PO only when its id matches and the visibility base marks it as unapproved and actionable now (is_unapproved_actionable_now = 1). Otherwise it returns no record. What makes a PO actionable is defined in the visibility base (__PO_VISIBILITY_BASE__), which is not in this file.
