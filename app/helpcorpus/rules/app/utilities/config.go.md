<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/utilities/config.go","sha256":"a37b7ec7f56f8b645565cbfa621f3662163d40e855b07959ef87c142864fd66d","model":"claude-sonnet-5-5"} -->

# Rules in app/utilities/config.go

## Unreadable config JSON gives an error with the default

Lines 16-52

Asked as: config value broken; settings won't parse; bad config json error; config shows error default

When the stored value of a config row (app_config) is not valid JSON, reading a Boolean switch returns the default value along with an error. A missing config row is not an error: it simply gives the default.

## Non-Boolean config switch value falls back to its default

Lines 32-52

Asked as: setting holds a string; config switch wrong value; switch ignored default used; config value not true false

When a config switch (a property in app_config) holds something other than true or false, such as a string, Turbo ignores it and uses the switch's default. The same happens when the property is absent or the config row for the domain is missing.

## Job editing switch defaults to on when config is missing

Lines 54-59

Asked as: can still edit jobs; jobs editable without setting; job editing allowed by default; job creation open by default

Turbo allows job creation, editing and absorb unless the jobs setting create_edit_absorb is false. If the jobs config row or that property is missing, or is not true/false, Turbo allows it (fail-open).

## Project request approval requirement defaults to off

Lines 61-83

Asked as: job request approval not required; no approval for project requests; jobs must be boolean error; project request approval default

Turbo requires project request approval only when the jobs setting require_job_request_approval is true. If the jobs config row or the property is missing, approval is not required. A value that is not true/false produces the error 'jobs.require_job_request_approval must be Boolean'.

## Job editing refused while disabled

Lines 85-92

Asked as: cannot create job; job editing turned off; jobs locked during transition; job editing is disabled

When job editing is switched off, Turbo refuses job creation and editing with 'job editing is currently disabled'. The detail reads 'job creation and editing is disabled during transition' (jobs_editing_disabled).

## Expense, PO and vendor editing switch defaults to on

Lines 94-99

Asked as: expense edits allowed by default; vendor creation always allowed; PO absorb allowed by default; expense setting missing allowed

Turbo allows expense, vendor and purchase order creation, editing and absorb unless the expenses setting create_edit_absorb is false. A missing config or property means allowed (fail-open).

## Time editing switch defaults to on

Lines 101-107

Asked as: time entries editable by default; timesheet changes allowed default; time editing open by default; time setting missing allowed

Turbo allows time entry creation and editing, and selected time sheet changes, unless the time setting create_edit is false. A missing config or property means allowed (fail-open).

## Notification features default to off

Lines 109-114

Asked as: notification not sending; email template disabled; notifications off by default; notification code missing

A notification template is enabled only when the notifications setting has its template code set to true. If the config or the code is missing or not true/false, the notification is off (fail-closed).

## PO second approver timeout defaults to 24 hours

Lines 116-145

Asked as: second approver deadline; PO approval timeout hours; how long PO second approval waits; purchase order timeout default

The purchase order second-stage timeout is the purchase_orders setting second_stage_timeout_hours. Turbo uses 24 hours when the setting is missing, zero, negative or not a number.

## Legacy PO create/update flow is off by default

Lines 147-152

Asked as: old PO screen missing; legacy purchase order form; hidden PO create page; old PO update flow

The hidden legacy purchase order create/update flow is enabled only when purchase_orders.enable_legacy_po_create_update is true. If the setting is missing, it is off.

## Allowed expense excess over PO: defaults and validation

Lines 178-222

Asked as: PO excess percent setting; allowed overspend over PO; invalid excess percentage config; PO overage default values

The allowed excess of expenses over a PO total comes from expenses.po_expense_allowed_excess. Turbo falls back to the built-in defaults (constants MAX_PURCHASE_ORDER_EXCESS_PERCENT and MAX_PURCHASE_ORDER_EXCESS_VALUE, mode lesser_of) when it is missing or not an object. Each field is checked separately and an invalid one keeps its default. Percent must be 0 to 100 (entered as 5 for 5%), value must be 0 or more, and mode must be lesser_of or greater_of.

## PO currency and rate for the excess limit

Lines 224-250

Asked as: foreign currency PO rate; exchange rate for PO limit; PO in USD conversion; currency rate missing on PO

A PO in the home currency is never converted, and a currency that can't be found is treated as the home currency. For a foreign PO the rate is approval_total_home divided by approval_total. If either total is missing or zero, Turbo uses the currency's current rate, or 1 if it has none or it is not positive.

## How the maximum expense total for a PO is calculated

Lines 252-308

Asked as: max expense on PO; how expense limit is computed; PO overage limit calculation; maximum expenses against purchase order

The limit is the PO total plus the percent excess or the fixed amount. The fixed amount is in home currency (CAD) and is divided by the PO's home rate on a foreign PO. In lesser_of mode (default) the fixed amount is used only when it is strictly less than the percent excess, so an equal value uses the percent. In greater_of mode the fixed amount is used when it is greater than or equal to the percent excess. The limit is floored to the cent.

## Zero-total PO reports a percentage excess

Lines 261-264

Asked as: zero total PO shows percent; $0 PO percentage message; PO with no total excess

When the PO total is 0, the percent excess is 0 and the percent branch is chosen in lesser_of mode, so the text shows a percentage even though the allowed excess is $0.

## Excess text in the over-limit message

Lines 283-301

Asked as: over limit message format; how excess is displayed; PO over limit wording; excess amount shown wrong

The excess is shown as a percentage like '5.00%' when the percent applies, or as '$100.00' when the fixed amount applies on a home-currency PO. On a foreign PO it is shown as '100.00 CAD (74.07 USD)'. The PO total is shown as '$1500.00', or in its own currency on a foreign PO.

## Text in a numeric config setting is rejected

Lines 310-326

Asked as: number saved as text; quoted number config ignored; config number in quotes; numeric setting rejected

Numeric config settings, such as the PO excess percent and value and the no-PO expense limit, accept only JSON numbers. A value written as text, such as "5", is rejected and Turbo uses the default for that setting.

## No-PO expense limit default and zero meaning

Lines 328-341

Asked as: no PO expense threshold; expenses without PO limit; zero no PO limit meaning; expense amount needs PO

The no-PO expense limit comes from expenses.no_po_expense_limit. Turbo uses the built-in default (NO_PO_EXPENSE_LIMIT) when the config or key is missing, the value is not a number, or it is negative. A value of 0 means all non-exempt expenses require a PO.

## Expense editing refused while disabled

Lines 343-350

Asked as: cannot edit expenses; expense editing locked; expenses frozen during transition; expense changes blocked

When expense editing is switched off, Turbo refuses with 'expense editing is currently disabled'. The detail reads 'expense, purchase order, and vendor creation and editing is disabled during transition' (expenses_editing_disabled).

## Time editing refused while disabled

Lines 352-359

Asked as: cannot edit time; time entries locked; timesheet submission disabled; time editing blocked

When time editing is switched off, Turbo refuses with 'time editing is currently disabled'. The detail says time entry and time amendment creation/editing, plus time sheet submission/approval, are disabled during transition (time_editing_disabled).
