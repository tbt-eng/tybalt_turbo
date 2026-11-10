<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/purchase_orders_legacy.go","sha256":"66f93440b0582761e9cab5ba162fbc4dfd7dda9aeb2aa5ab530b212722e08fc0","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/purchase_orders_legacy.go

## Legacy PO number format

Lines 18-18

A legacy purchase order's number (po_number) must match YYMM-NNNN. YY is 24, 25 or 26, MM is 01 to 12, and NNNN is 5 followed by three digits. Otherwise Turbo refuses with "Legacy PO number must match YYMM-NNNN with YY 24/25/26 and NNNN in the 5XXX range." Turbo trims spaces from the number before it checks.

## Fields editable on a legacy PO

Lines 20-34

Only these fields can be sent to the legacy purchase order endpoint: uid, approver, po_number, date, division, description, payment_type, total, vendor, type, kind, branch and job. A request with any other field is refused with "this field cannot be edited through the legacy purchase order endpoint".

## Required fields on a legacy PO

Lines 36-51

A legacy purchase order needs these fields, each non-blank after trimming. Turbo names the first missing one: staff member ("Staff member is required."), approver, legacy PO number, date, division, description, payment type, vendor, type, kind and branch. Each message has the form "<Field> is required." The job is not required.

## Legacy PO access requires sign-in

Lines 53-65

Every legacy purchase order action needs a signed-in user. Without one, Turbo refuses with "legacy purchase order access requires authentication".

## Legacy PO blocked when expense editing is off

Lines 67-78

Turbo refuses legacy purchase order actions when the check requireExpensesEditing for purchase_orders fails. The refusal carries the message of that check. Another file defines when that check fails.

## Legacy PO create/update switch

Lines 80-104

Turbo refuses legacy purchase order actions with "legacy purchase order create/update is disabled" while the setting behind IsLegacyPOCreateUpdateEnabled is off. If Turbo cannot read the setting, it refuses with "unable to validate legacy purchase order configuration".

## Legacy PO needs a dedicated claim

Lines 106-130

Only staff with the legacy PO create/update claim (LEGACY_PO_CREATE_UPDATE_CLAIM_NAME) can view for edit, create or update legacy purchase orders. Others see "you do not have permission to create or edit legacy purchase orders". If the claim check fails, Turbo says "unable to validate legacy purchase order permissions".

## Legacy PO payload rejects fields that are not editable

Lines 135-153

When a create or update request for a legacy purchase order contains a field outside the editable list, Turbo refuses it with 'legacy purchase order payload contains a field that is not editable'. The field's own error reads 'this field cannot be edited through the legacy purchase order endpoint' (code field_not_allowed).

## Only editable fields from the request are applied to a legacy PO

Lines 155-161

Turbo copies only the editable fields from the request onto the legacy purchase order record. It ignores any other value.

## Legacy PO missing a required field is refused

Lines 163-177

Turbo checks the required fields in order and refuses the first one that is empty or only spaces. The error is 'legacy purchase order is missing a required field' plus that field's message, such as 'Approver is required.' or 'Legacy PO number is required.' A field made only of spaces counts as empty.

## Legacy PO total must be positive

Lines 179-190

A legacy purchase order is refused with "Total must be greater than 0." when its total is 0 or less. Turbo checks the other required fields first.

## Legacy PO values Turbo sets itself

Lines 196-217

On every create or update, Turbo sets the legacy purchase order to status Active and flags it as a manual legacy entry (legacy_manual_entry true, _imported false). It clears the second approver and its priority, the second approval, and the rejection fields. It clears the cancel and close fields and sets closed_by_system false. It clears end date, frequency, parent PO, category and attachment with its hash, and sets covered_within_project_budget false. Whatever the user sends for these is overwritten.

## Legacy PO number is trimmed before the format check

Lines 196-196

Turbo trims spaces from the legacy PO number (po_number) before it checks the format. If the format fails, Turbo refuses with 'legacy purchase order number is invalid'. The field message is 'Legacy PO number must match YYMM-NNNN with YY 24/25/26 and NNNN in the 5XXX range.'

## Legacy PO number format check runs after Turbo resets fields

Lines 217-230

Turbo resets the system-managed fields first, then checks the PO number against the pattern. A bad number is refused with code invalid_legacy_po_number.

## Legacy PO type restricted

Lines 232-245

A legacy purchase order's type must be One-Time or Cumulative. Any other type is refused with "Type must be One-Time or Cumulative for legacy purchase orders." Recurring is not allowed.

## Legacy PO approved date set on create

Lines 247-249

When a legacy purchase order is created, Turbo sets its approved date to the current time. An update does not change the approved date.

## Legacy PO validation errors

Lines 254-281

Field validation failures on a legacy purchase order come back as "validation failed", with a message and code for each field. The code is validation_error when the check gives none.

## Only legacy POs open in legacy editor

Lines 283-309

The legacy endpoint refuses a purchase order that is not a legacy manual entry. It says "legacy create/update can only be used with legacy manual purchase orders". If the purchase order does not exist, Turbo says "purchase order not found".

## Closed or cancelled legacy POs cannot be edited

Lines 311-340

A legacy purchase order with status Closed cannot be updated: "closed legacy purchase orders cannot be edited". One with status Cancelled cannot be updated either: "cancelled legacy purchase orders cannot be edited". The other statuses can be updated.

## Legacy PO errors without a specific type use a generic code

Lines 342-357

When a legacy purchase order action fails, Turbo shows the error's own message and code. An error that is neither a hook error nor a code error gets the code "legacy_purchase_order_error" and its message text. A code error keeps its own code and message and takes the status the handler set for that step.

## Load legacy PO for editing

Lines 359-387

Opening a legacy purchase order for edit needs the route access checks. The ID is required ("purchase order id is required"). An unknown ID gives "purchase order not found". A purchase order that is not legacy is refused. The view does not refuse Closed or Cancelled ones.

## Legacy PO create checks

Lines 389-460

When Turbo creates a legacy purchase order, it checks in this order: route access, allowed fields, required fields, then number and type. After that it runs the standard PO cleaning (CleanPurchaseOrder). Then it checks that the user can use the record's branch (EnsureUserCanUseRecordBranch) and that the division is active (EnsureActiveDivision). Finally it runs PO validation (ValidatePurchaseOrder). The save happens in one transaction, so any failure saves nothing.

## Legacy PO update checks

Lines 462-546

When Turbo updates a legacy purchase order, it checks route access and allowed fields first. A missing purchase order gives "purchase order not found". Then it refuses a non-legacy, Closed or Cancelled one. It applies the payload to the saved record, then runs the same required-field, normalization, cleaning, branch, active-division and validation checks as a create. The save is one transaction.

## Legacy PO update needs an ID

Lines 465-471

An update request without a purchase order ID is refused with "purchase order id is required".
