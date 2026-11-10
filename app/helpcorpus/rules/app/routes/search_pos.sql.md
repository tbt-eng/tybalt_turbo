<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/search_pos.sql","sha256":"8337cf49055c96ade8195b6a5184f38d35637223480e49e3ba527b8ad93f248f","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/search_pos.sql

## PO search shows only POs that pass a visibility flag

Lines 77-84

Asked as: PO missing from search; can't find purchase order; PO not showing in search; hidden purchase order

The PO search returns a purchase order only when at least one visibility flag is 1. The flags are is_active_visible, is_closed_cancelled_visible, is_legacy_visible and is_revision_visible. The flags come from the visibility base query (__PO_VISIBILITY_BASE__), which is not in this file, so this file does not show who sees which PO.

## PO search lists only Active, Closed and Cancelled POs, plus revisions

Lines 85-85

Asked as: PO status not showing; draft PO missing; search skips some POs; pending PO not in results

The PO search returns a purchase order only if its status is 'Active', 'Closed' or 'Cancelled', or if is_revision_visible is 1. A PO with any other status does not appear unless it is visible as a revision. This status check applies in addition to the visibility flags.

## PO search sorts newest PO date first

Lines 86-88

Asked as: PO search order; latest PO first; sort purchase orders by date; newest POs at top

The PO search sorts results by PO date, newest first. POs with the same date are sorted by last update (updated), most recent first.
