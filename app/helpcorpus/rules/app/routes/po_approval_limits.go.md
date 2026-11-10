<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/po_approval_limits.go","sha256":"8b83fa0defa1d07906fe85c1249e2da42f00db392ecfeb38b017342542929b7d","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/po_approval_limits.go

## Who can view PO approval limits

Lines 50-59

Asked as: view PO approval limits; who can see approval limits; no permission approval limits report; PO approver claim access

Turbo lets a person view the PO approval limits report only if they hold the report claim (report) or the PO approver claim (po_approver). Anyone else gets the message "you do not have permission to view PO approval limits". The report claim is checked first, and the PO approver claim is checked only if the first check does not allow access.

## Access check failure message

Lines 54-56

Asked as: failed to check approval-limit access; cannot check claims; access check error approval limits; approval limits access error

If Turbo cannot check the person's claims, it refuses the request with "failed to check approval-limit access".

## Only active admin-profile PO approvers are listed

Lines 66-83

Asked as: approver missing from limits report; inactive admin approver not listed; why is approver not showing; approver list sorting order

The report lists one row for each holder of the PO approver claim who has a profile and an active admin profile (admin_profiles.active = 1). A PO approver whose admin profile is not active, or who has none, does not appear. Rows are sorted by surname, then given name, ignoring case.

## Unconfigured approvers show zero limits

Lines 68-75

Asked as: limits show zero; approver limits all 0; limits not configured; no divisions shown for approver

For each approver, the report shows whether limits are set up (configured). It is true only when the approver has a PO approver properties record (po_approver_props). With no record, every limit shows 0: max_amount, project_max, sponsorship_max, staff_and_social_max, media_and_event_max and computer_max. Turbo shows no divisions in that case.

## Approver limits and divisions shown per row

Lines 69-96

Asked as: approver limit amounts per row; division list for approver; empty divisions in report; approver maximum amount columns

Each approver row shows the overall maximum amount, the project, sponsorship, staff and social, media and event, and computer limits, and the list of divisions. If the divisions value is empty, Turbo shows an empty list.

## Approval limits load failure

Lines 84-86

Asked as: failed to load PO approval limits; approval limits report won't load; approver rows error; limits report error

If Turbo cannot load the approver rows, the report fails with "failed to load PO approval limits".

## Division list read failure

Lines 90-92

Asked as: failed to read PO approval divisions; approver divisions read error; cannot read approver divisions; divisions unreadable approver

If an approver's stored divisions cannot be read, the report fails with "failed to read PO approval divisions".

## Categories with second approval thresholds

Lines 98-111

Asked as: second approval threshold by category; expenditure kind thresholds; PO category limit field; failed to resolve PO approval category

The report lists every expenditure kind as a category, with its English label and its second approval threshold (second_approval_threshold). Categories are ordered by ui_order, then id. Turbo also gives each category the name of the approver limit field that applies to it, found by ResolvePOApproverLimitColumn. If that lookup fails, the report fails with "failed to resolve PO approval category"; if loading fails, with "failed to load PO approval categories".

## Divisions listed with active flag

Lines 112-114

Asked as: list of divisions inactive; division codes and names report; inactive divisions shown; division active flag

The report lists all divisions, active or not, with their code, name and active flag. They are sorted by code, then name. If loading fails, the report fails with "failed to load PO approval divisions".
