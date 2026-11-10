<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/validate_expenses.go","sha256":"607fd2b3a1fe0166021a1ec122a6f3b07cd9e3a26cf85e44a8e30d3cdfb4881a","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/validate_expenses.go

## Unreadable PO dates stop expense validation

Lines 35-62

If the purchase order's date cannot be parsed, Turbo refuses with "error parsing purchase order date". For a Recurring PO, an unreadable end date gives "error parsing purchase order end date".

## Cumulative PO refuses expenses that overflow its total

Lines 69-105

For a Cumulative purchase order, Turbo adds the existing expenses total to the new expense total, compared in whole cents. If the sum is above the PO total, it refuses with "cumulative expenses exceed purchase order total" (code cumulative_po_overflow). The error gives the overflow amount, the PO id, the PO number and the PO total. An expense that exactly fills the PO is accepted.

## Expense currency must exist

Lines 119-131

If the expense's currency cannot be resolved, Turbo refuses with "referenced currency not found" (field currency, code not_found).

## Settled total fallback for limit comparisons

Lines 133-143

When the settled total (settled_total) is 0 or less, Turbo compares the expense total instead. In the home currency it uses the total as entered. In another currency it uses the indicative home-currency amount (IndicativeHomeAmount).

## Attachment required except for allowance, mileage, personal reimbursement

Lines 145-167

Every payment type except Allowance, Mileage and PersonalReimbursement needs an attachment, or Turbo refuses with "attachment is required". Any one of these satisfies it: a document attachment (attachment_document), a source expense (source_expense), a newly uploaded file, or a missing-attachment reason (attachment_missing_reason) on an existing expense. The missing reason does not count on a new expense.

## Expense with a PO needs the PO record

Lines 169-181

If an expense names a purchase order but the PO record is not supplied, Turbo cannot validate it and refuses with "an expense against a purchase_orders record cannot be validated without a corresponding purchase order record".

## Expense job and currency must match its PO

Lines 183-214

An expense linked to a PO must have the same job as the PO. Two empty jobs also match. Otherwise Turbo says "job must match purchase order job". The expense currency must also equal the PO currency, or Turbo says "currency must match purchase order currency". Both comparisons use the effective currency code.

## Expenditure kind required, with legacy fallback on update

Lines 216-252

An expense needs a kind, or Turbo says "kind is required". When an existing expense whose original record also has no kind is updated, Turbo sets a default kind (DefaultExpenditureKindIDForJob) from whether it has a job. A new expense never gets this fallback. A kind that does not exist gives "invalid expenditure kind".

## Expense date rules

Lines 255-267

The date is required and must be a valid date (YYYY-MM-DD). With a PO, it must be on or after the PO date. With a Recurring PO, it must also be on or before the PO end date.

## Description required and at least 4 characters except allowances

Lines 268-274

For every payment type except Allowance, the description is required and needs at least 4 characters. Otherwise Turbo says "required for non-allowance expenses" or "must be at least 4 characters".

## Vendor required except allowance, personal reimbursement, mileage

Lines 275-281

The vendor is required, with the message "required for this expense type", unless the payment type is Allowance, PersonalReimbursement or Mileage.

## Credit card last 4 digits rules

Lines 282-290

A CorporateCreditCard expense needs cc_last_4_digits of exactly 4 characters, or Turbo says "required for corporate credit card expenses" or "must be 4 digits". For other payment types the field must be empty, or Turbo says it is not applicable.

## Expense total must be above zero

Lines 291-294

The expense total must be at least 0.01, or Turbo says "must be greater than 0".

## One-Time and Recurring PO expense total limit

Lines 295-308

For an expense on a One-Time or Recurring PO, the total in cents must not exceed the PO total plus the allowed excess. The excess is calculated by CalculatePOExpenseTotalLimit from the app's configuration. Otherwise Turbo says "expense exceeds purchase order total of $X by more than ...". Turbo does not stop a second expense on a One-Time PO.

## Mileage requires distance

Lines 309-314

A Mileage expense needs a distance, or Turbo says "required for mileage expenses".

## Purchase order required for job expenses

Lines 315-320

An expense with a job needs a purchase order, or Turbo says "required for all expenses with a job". Mileage, FuelCard, PersonalReimbursement and Allowance expenses are exempt.

## Allowance types required for allowances

Lines 321-326

An Allowance expense needs allowance types (allowance_types), or Turbo says "required for allowance expenses".

## Purchase order required above the no-PO limit

Lines 329-341

When the LIMIT_NON_PO_AMOUNTS constant is on, an expense with no PO whose comparable total is at or above the no-PO limit is refused. Turbo says "a purchase order is required for expenses of $X or more". The limit includes the equal value. Mileage, FuelCard, PersonalReimbursement and Allowance are exempt. OnAccount is exempt only when the total limit is bypassed (byPassTotalLimit). This message replaces other total errors.

## Foreign-currency out-of-pocket expense needs a settled total

Lines 343-357

For an out-of-pocket (Expense) payment in a non-home currency, the settled total must be above 0, or Turbo says "settled total must be greater than 0". It must also be within the tolerance of the total (IsSettledTotalWithinTolerance), or Turbo gives a tolerance message.

## Vendor must exist and be Active

Lines 359-367

If a vendor is set, it must exist, or Turbo says "invalid vendor reference". It must have status Active, or Turbo says "vendor must be active".

## Category must belong to the expense's job

Lines 369-381

A category needs a job, or Turbo says "category requires a job". The category must exist ("invalid category reference") and belong to the expense's job ("category must belong to the selected job").

## Job must exist and be Active; division allocation check is gated

Lines 383-403

The job must exist ("invalid job reference") and have status Active ("Job status must be Active"). Turbo then checks that the division is allocated to the job, but only if shouldValidateJobDivisionAllocationOnRecord says so. This is for new expenses, or when the job or division changes. Legacy expenses are otherwise left alone.
