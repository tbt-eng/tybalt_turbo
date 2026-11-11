<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/validate_expenses.go","sha256":"274d50c86a59c72035a591a63a7a0c51ff889df7ec8299fedb4ae124e2fa5dee","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/validate_expenses.go

## Purchase order date cannot be read

Lines 36-63

Asked as: cannot read PO date; error parsing purchase order date; PO end date invalid; purchase order date broken; recurring PO end date error

If the date of the linked purchase order cannot be parsed, Turbo refuses with "error parsing purchase order date". For a Recurring PO whose end date cannot be parsed, it refuses with "error parsing purchase order end date".

## Excess limit for One-Time and Recurring POs is computed from PO total and currency

Lines 65-69

Asked as: how is PO limit calculated; expense allowed amount on PO; PO currency limit; excess allowance on purchase order; what can I claim against PO

For a purchase order, Turbo works out the allowed expense total from the PO total, the PO currency and the excess configuration (GetPOExpenseExcessConfig, CalculatePOExpenseTotalLimit). That limit is used in the One-Time and Recurring total check.

## Cumulative PO refuses expense that overflows the PO total

Lines 71-107

Asked as: cumulative PO over total; expenses exceed PO total; PO already used up; too much spent on cumulative PO; overflow on purchase order

When an expense is against a Cumulative purchase order, Turbo adds the expense total to the existing expenses on that PO, compared in whole cents. If the sum is more than the PO total, Turbo refuses with "cumulative expenses exceed purchase order total" (cumulative_po_overflow). The error gives the PO, its number, its total and the overflow amount. An expense that exactly fills the PO is accepted.

## Cumulative PO overflow compares whole cents

Lines 78-84

Asked as: expense exactly fills PO; rounding error PO total; cents rounding purchase order; PO total exact match refused; float rounding expense

Turbo checks the Cumulative PO overflow in whole cents, so an expense that exactly fills the PO total is accepted. A float rounding error does not cause a refusal.

## Currency must exist

Lines 121-133

Asked as: currency not found; unknown currency on expense; wrong currency code; currency missing from list; invalid currency expense

If the expense currency cannot be found, Turbo refuses with "referenced currency not found" (not_found on currency).

## PO required for expenses at or above the no-PO limit

Lines 135-145

Asked as: need PO for big expense; no PO expense too large; purchase order required amount; expense over no-PO limit; why must I have a PO

When LIMIT_NON_PO_AMOUNTS is on, an expense with no PO is refused if its comparable total is at or above the no-PO limit. The refusal reads "a purchase order is required for expenses of $X or more". Mileage, FuelCard, PersonalReimbursement and Allowance are exempt. An OnAccount expense is also exempt when byPassTotalLimit is set. The comparable total is the settled total if it is above 0, else the total for home currency, else an indicative home-currency amount. This message replaces the other total errors.

## Attachment required except for some payment types

Lines 147-169

Asked as: receipt missing from expense; upload receipt required; no attachment on expense; cannot submit without receipt; missing receipt reason

Every expense needs an attachment, except Allowance, Mileage and PersonalReimbursement. Turbo accepts an attachment document, a newly uploaded file, a source expense (source_expense), or a missing-attachment reason (attachment_missing_reason). A missing-attachment reason counts only on an existing expense, not a new one. Otherwise Turbo refuses with "attachment is required".

## Expense with a PO must have its PO record

Lines 171-183

Asked as: PO record missing; expense linked to deleted PO; cannot validate PO expense; purchase order not loaded; PO link broken

An expense that names a purchase order cannot be validated without the PO record. Turbo refuses with "an expense against a purchase_orders record cannot be validated without a corresponding purchase order record".

## Expense job and currency must match its PO

Lines 185-216

Asked as: expense job differs from PO; wrong job on expense; expense currency different from PO; PO job mismatch; expense does not match PO

An expense linked to a purchase order must have the same job as the PO; both may be empty. Otherwise Turbo refuses with "job must match purchase order job". The expense currency must also equal the PO currency, otherwise Turbo refuses with "currency must match purchase order currency".

## Expenditure kind required, with legacy fallback

Lines 218-254

Asked as: no expense type chosen; missing expenditure type; which kind of expense; kind field empty; unknown expense type

An expense needs an expenditure kind (kind). An existing expense whose original record also has no kind gets a default kind, chosen by whether it has a job. A new expense with no kind is refused with "kind is required". An unknown kind is refused with "invalid expenditure kind".

## Expense date rules

Lines 257-269

Asked as: expense date before PO; wrong date on expense; date after PO ends; bad expense date format; expense date outside PO dates

The date is required ("date is required") and must be a valid YYYY-MM-DD date ("must be a valid date"). With a PO, the date must not be before the PO date. For a Recurring PO, the date must also not be after the PO end date. Exact error text comes from DateStringLimit, which is in another file.

## Description required except for allowances

Lines 270-276

Asked as: description too short; empty description on expense; what to write as description; description missing error; expense description required

For every payment type except Allowance, the description is required ("required for non-allowance expenses") and must be at least 4 characters ("must be at least 4 characters").

## Vendor required except allowance, personal reimbursement, mileage

Lines 277-283

Asked as: no supplier on expense; merchant field empty; vendor blank error; who was paid missing; expense needs a supplier

The vendor is required ("required for this expense type") unless the payment type is Allowance, PersonalReimbursement or Mileage.

## Credit card last 4 digits

Lines 284-292

Asked as: card last digits; last four card numbers; credit card number on expense; card digits wrong; cc digits not allowed

A CorporateCreditCard expense requires cc_last_4_digits ("required for corporate credit card expenses") of exactly 4 characters ("must be 4 digits"). For any other payment type the field must be empty, otherwise Turbo refuses with "cc_last_4_digits is not applicable for non-corporate credit card expenses".

## Expense total must be positive

Lines 293-296

Asked as: zero amount expense; negative expense amount; amount too small; expense cannot be free; total must be above zero

The expense total must be at least 0.01. Otherwise Turbo refuses with "must be greater than 0".

## One-Time and Recurring PO excess limit on expense total

Lines 297-304

Asked as: expense over PO amount; too big for purchase order; PO excess limit exceeded; expense more than PO allows; over PO by too much

For an expense on a One-Time or Recurring PO, Turbo compares the total in cents with a limit from the PO total, the PO currency and the excess configuration (GetPOExpenseExcessConfig). An expense above that limit is refused with "expense exceeds purchase order total of ... by more than ...". An amount equal to the limit passes. Cumulative POs use the overflow check instead.

## Mileage requires distance

Lines 311-316

Asked as: mileage miles missing; no distance entered; how many miles; distance needed for claim; mileage field empty

A Mileage expense requires a distance ("required for mileage expenses").

## PO required for expenses with a job

Lines 317-322

Asked as: job expense without PO; need PO for job; job expense no order; PO missing on job expense; why does job need PO

An expense with a job must have a purchase order ("required for all expenses with a job"), except for Mileage, FuelCard, PersonalReimbursement and Allowance.

## Allowance requires allowance types

Lines 323-328

Asked as: allowance type missing; which allowance to pick; allowance kind blank; no allowance selected; pick allowance category

An Allowance expense requires allowance_types ("required for allowance expenses").

## Expense without a PO is refused at or above the no-PO limit

Lines 331-343

Asked as: expense with no PO too high; no order expense over limit; cannot claim without order; non-PO spending limit; large expense without purchase order

Turbo refuses an expense with no purchase order when its comparable total is at or above the no-PO limit (GetNoPOExpenseLimit). The message is "a purchase order is required for expenses of $X or more". The comparison includes the equal value. The check runs only when LIMIT_NON_PO_AMOUNTS is on. Mileage, fuel card, personal reimbursement and allowance expenses are exempt.

## On Account expenses can bypass the no-PO limit

Lines 331-338

Asked as: on account skip limit; bypass no-PO limit; on account expense too large; skip purchase order requirement; OnAccount limit override

The no-PO limit check is skipped when the caller sets the bypass (byPassTotalLimit) and the payment type is OnAccount. Other payment types still face the limit.

## Foreign-currency out-of-pocket expense needs a settled total

Lines 345-359

Asked as: foreign currency settled amount; settled total zero; card charged in other currency; exchange amount doesn't match; out of pocket foreign expense

For an Expense payment type in a non-home currency, the settled total must be greater than 0 ("settled total must be greater than 0"). It must also be within a tolerance of the total, otherwise Turbo refuses with a tolerance message.

## Vendor must exist and be Active

Lines 361-369

Asked as: supplier inactive; vendor not active; unknown supplier id; vendor disabled; pick different vendor

A vendor on an expense must exist ("invalid vendor reference") and have status Active ("vendor must be active").

## Category needs a job and must belong to it

Lines 371-383

Asked as: category not allowed on job; expense category needs job; wrong category for job; category blank without job; category not in job

A category on an expense requires a job ("category requires a job"). The category must exist ("invalid category reference") and belong to the expense's job ("category must belong to the selected job").

## Job must exist and be Active; division must be allocated

Lines 385-405

Asked as: job is inactive; division not on job; job not found on expense; division not allocated; cannot pick this job

The job on an expense must exist ("invalid job reference") and have status Active ("Job status must be Active"). Turbo also checks that the division is allocated to the job, when the expense is new or its job or division changes (shouldValidateJobDivisionAllocationOnRecord). Other updates to legacy expenses skip this check.
