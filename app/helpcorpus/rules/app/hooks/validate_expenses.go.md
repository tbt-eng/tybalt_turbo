<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/validate_expenses.go","sha256":"3b949dce54ca30270d01ea681a12fa2b14ba676a21fb3807d85f75f71f03c183","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/validate_expenses.go

## Unparseable PO dates stop expense validation

Lines 37-64

Asked as: purchase order date unreadable; error parsing PO date; recurring PO end date broken; expense won't validate bad PO date

Turbo fails the expense with 'error parsing purchase order date' if the PO date cannot be read. For a Recurring PO it fails with 'error parsing purchase order end date' if the end date cannot be read.

## Expense currency must match its purchase order currency

Lines 66-84

Asked as: expense currency differs from PO; wrong currency on purchase order; PO currency mismatch; currency doesn't match the PO

When an expense is linked to a purchase order, Turbo compares the effective currency of the expense with that of the PO. If they differ, it refuses the expense with 'currency must match purchase order currency'. This check runs before the PO total checks.

## Cumulative PO overflow refuses the expense

Lines 92-129

Asked as: PO is full; cumulative PO over limit; expenses exceed PO total; PO already used up

For a Cumulative purchase order, Turbo adds the existing expenses total and this expense's total, compared in whole cents. If the sum is above the PO total, it refuses the expense with 'cumulative expenses exceed purchase order total'. The error carries the PO, its number, its total and the overflow amount. An expense that exactly fills the PO is accepted.

## PO required at or above the no-PO limit

Lines 141-166

Asked as: need a PO for this amount; no PO over limit; expense too big without purchase order; PO required above threshold

If LIMIT_NON_PO_AMOUNTS is on, an expense with no PO is refused when its comparable total is at or above the no-PO limit. The message is 'a purchase order is required for expenses of $X or more'. The equal value is refused. The comparable total is the settled total if above 0, else the total, or an indicative home-currency amount for a foreign currency. Mileage, FuelCard, PersonalReimbursement and Allowance are exempt. OnAccount is exempt when the total limit is bypassed.

## Unknown currency is refused

Lines 142-154

Asked as: currency not found; unknown currency on expense; can't find expense currency; invalid currency code

If the expense currency cannot be found, Turbo refuses with 'referenced currency not found'.

## Attachment required for most payment types

Lines 168-190

Asked as: receipt missing; no attachment on expense; upload receipt required; can't submit without receipt

Every payment type except Allowance, Mileage and PersonalReimbursement needs an attachment, or Turbo refuses with 'attachment is required'. Any one of these satisfies it: an attachment document, a source expense, a newly uploaded file, or a missing-attachment reason. The missing-attachment reason counts only on an existing expense, not a new one.

## Expense with a PO needs the PO record

Lines 192-204

Asked as: PO record missing; expense linked to PO but no PO; purchase order not loaded; can't validate against PO

If an expense has a purchase order but the PO record is not supplied, Turbo cannot validate it. It refuses with 'an expense against a purchase_orders record cannot be validated without a corresponding purchase order record'.

## Expense job must match its PO job

Lines 206-222

Asked as: job differs from PO job; expense job wrong for PO; PO job mismatch; different job than purchase order

An expense linked to a purchase order must have the same job as the PO. Both may be empty. Otherwise Turbo refuses with 'job must match purchase order job'.

## Expense kind is required and must exist

Lines 224-260

Asked as: no expense kind; missing expenditure type; invalid kind on expense; expense type not recognised

A new expense must have a kind (expenditure kind), or Turbo refuses with 'kind is required'. An existing legacy expense with no kind gets a default kind based on whether it has a job. A kind that does not exist is refused with 'invalid expenditure kind'.

## Expense date rules

Lines 263-275

Asked as: expense date before PO date; date outside PO window; wrong expense date; date after PO end

The date is required and must be a valid date (YYYY-MM-DD). With a PO, it must be on or after the PO date. With a Recurring PO, it must also be on or before the PO end date.

## Description rules

Lines 276-282

Asked as: missing description; description too short; no description entered; expense note blank

Except for Allowance expenses, the description is required and must be at least 4 characters. The messages are 'required for non-allowance expenses' and 'must be at least 4 characters'.

## Vendor required except some payment types

Lines 283-289

Asked as: no vendor on expense; supplier missing; vendor blank; who was paid field empty

The vendor is required ('required for this expense type') unless the payment type is Allowance, PersonalReimbursement or Mileage.

## Credit card last 4 digits rules

Lines 290-298

Asked as: card digits wrong; last four card numbers; credit card 4 digits; card number field error

For CorporateCreditCard expenses, cc_last_4_digits is required and must be exactly 4 digits ('must be 4 digits'). For all other payment types it must be empty ('cc_last_4_digits is not applicable for non-corporate credit card expenses').

## Total must be positive and within PO limit

Lines 299-316

Asked as: amount over PO limit; expense more than PO; zero amount expense; total too high for PO

The expense total must be at least 0.01 ('must be greater than 0'). For a One-Time or Recurring PO, a total above the PO's limit is refused with 'expense exceeds purchase order total of X by more than Y'. The limit comes from the PO total plus the configured excess allowance. Amounts are compared in cents.

## Mileage requires distance

Lines 317-322

Asked as: mileage distance missing; no miles entered; mileage claim needs miles; distance blank on mileage

A Mileage expense must have a distance ('required for mileage expenses').

## PO required for expenses with a job

Lines 323-328

Asked as: job needs PO; expense with job no purchase order; job expense missing PO; PO needed for job

An expense with a job must have a purchase order ('required for all expenses with a job'). Mileage, FuelCard, PersonalReimbursement and Allowance are exempt.

## Allowance expenses need allowance types

Lines 329-334

Asked as: allowance types missing; no allowance selected; allowance needs type; pick allowance category

An Allowance expense must list its allowance types ('required for allowance expenses').

## Expense without PO is refused at or above the no-PO limit

Lines 337-349

Asked as: no PO over limit; expense over no-PO threshold; purchase order needed above limit; large expense without PO

An expense with no purchase order is refused when its comparable total is equal to or greater than the no-PO expense limit. The message is "a purchase order is required for expenses of $X or more", where X is the limit. Mileage, fuel card, personal reimbursement and allowance expenses are exempt. An on-account expense is exempt when the total-limit bypass is on (byPassTotalLimit). The check runs only when LIMIT_NON_PO_AMOUNTS is on.

## No-PO limit message replaces other total errors

Lines 345-348

Asked as: total error hidden; only PO message shown; wrong total error shown; amount message replaced

When the no-PO limit check refuses an expense, its message replaces any other error on the total. The person sees only the "a purchase order is required" message for the total, even if the total also failed another rule.

## Foreign-currency out-of-pocket expenses need settled total

Lines 351-365

Asked as: settled total missing; foreign currency expense settled amount; out of pocket settled total; settled amount tolerance

For a payment type of Expense in a non-home currency, the settled total must be greater than 0 ('settled total must be greater than 0'). It must also be within tolerance of the total, or Turbo refuses with a tolerance message.

## Vendor must exist and be active

Lines 367-375

Asked as: vendor not found; supplier inactive; vendor not active; bad vendor selected

A vendor that does not exist is refused with 'invalid vendor reference'. A vendor whose status is not Active is refused with 'vendor must be active'.

## Category must belong to the expense's job

Lines 377-389

Asked as: category wrong for job; category not on job; no category for job; category doesn't match job

A category needs a job ('category requires a job'). It must exist ('invalid category reference') and belong to the selected job ('category must belong to the selected job').

## Job must exist, be Active and allocate the division

Lines 391-411

Asked as: job inactive; job not found; division not allocated to job; job closed on expense

The job must exist ('invalid job reference') and have status Active ('Job status must be Active'). On create, or when the job or division changes, the division must also be allocated to the job. Other updates to older expenses skip the allocation check.
