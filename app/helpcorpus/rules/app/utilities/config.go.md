<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/utilities/config.go","sha256":"819a472d8760a6e4d9ac577c77853b18c220cefe1e1a3175c3e29865e8aac25b","model":"claude-sonnet-5-5"} -->

# Rules in app/utilities/config.go

## Missing or unreadable config domain returns nothing

Lines 17-31

Asked as: config value not showing; setting shows blank; default value used instead; config lookup returns nothing; app config empty

Turbo's config lookup (GetConfigValue) returns no value when the app_config row for a domain key does not exist. It also returns no value when the row cannot be found for any reason, and each setting then falls back to its own default. A stored value that is not valid JSON gives an error.

## Boolean config setting falls back to its default

Lines 33-53

Asked as: true setting ignored; toggle not applying; yes flag not working; boolean setting stuck; setting reverts to default

A boolean setting (GetConfigBool) uses its default when the config row is missing, when the property is absent, or when the property is not true or false. A lookup error returns the default along with the error. Only a real boolean overrides the default.

## Deploy group limits a feature to its members

Lines 55-87

Asked as: feature not showing for me; only some users see feature; not in rollout group; feature hidden on my account; beta feature missing for me

A feature whose settings name a deploy group (deploy_group) is on for a user only if the feature's own switch is on and the user is a member of that group (deploy_group_memberships). If no group is named, the switch alone decides. A user with no ID is never in a named group.

## Unknown deploy group gives the feature to nobody

Lines 66-87

Asked as: feature off for everyone; wrong group id disables feature; rollout group typo; feature vanished after group change; group id mistake turned off

A deploy group ID that matches no group has no members, so the feature is off for every user. A mistake in the ID can turn a feature off, but it cannot turn it on for everyone.

## Missing deploy group setting lets everyone in

Lines 99-122

Asked as: feature open to everyone; rollout not restricted; everyone can see feature; no group set access; why does everyone get it

If a feature has no settings row, or its settings have no deploy_group (or it is null), the deploy group check lets every user in. A deploy_group that is not text refuses with the error "<key>.deploy_group must be the id of a deploy group". A database read error also fails the check instead of letting users in.

## Job creating, editing and absorbing default to on

Lines 124-129

Asked as: cannot create job; job absorb blocked; job edits greyed out; cannot absorb job; job changes not allowed

Job creation, editing and absorbing are allowed unless the jobs setting create_edit_absorb is set to false. If the setting is missing or is not true/false, they stay allowed.

## Project request approval is off by default

Lines 131-153

Asked as: job request needs approval; project request approval missing; approval step skipped for projects; no approval on project requests; why requests not approved

Turbo requires project request approval only when the jobs setting require_job_request_approval is true. If the jobs row or the setting is missing, approval is not required. A value that is not true/false gives the error "jobs.require_job_request_approval must be Boolean".

## Message when job editing is disabled

Lines 155-162

Asked as: job editing disabled message; cannot edit jobs now; jobs locked during transition; job edit blocked error; job changes refused

When job editing is switched off, Turbo refuses with "job editing is currently disabled". The detail says "job creation and editing is disabled during transition".

## Expense, PO and vendor editing default to on

Lines 164-169

Asked as: cannot edit expense; vendor edit blocked; purchase order editing locked; can't absorb expense; expense changes not saved

Creating, editing and absorbing expenses, vendors and purchase orders is allowed unless the expenses setting create_edit_absorb is false. If the setting is missing, editing stays allowed.

## Time editing defaults to on

Lines 171-177

Asked as: cannot edit timesheet entry; time amendment blocked; cannot add time entry; timesheet entry locked; time changes refused

Time entry creation and editing, and some time sheet actions, are allowed unless the time setting create_edit is false. If the setting is missing, they stay allowed.

## Notifications default to off

Lines 179-184

Asked as: no notification emails; notification not sending; template notifications turned off; alerts not arriving; notifications silent

A notification template is on only if the notifications setting has true under that template's code. If the setting is missing or is not true/false, the notification is off.

## PO second approver timeout defaults to 24 hours

Lines 186-215

Asked as: second approver timeout hours; PO approval expires too soon; why PO approval timing out; purchase order reminder time; second stage deadline

The second-stage timeout for a purchase order is the purchase_orders setting second_stage_timeout_hours. It must be a number above 0. If it is missing, zero, negative or not a number, Turbo uses 24 hours.

## Legacy PO create/update flow is off by default

Lines 217-222

Asked as: old PO form missing; legacy purchase order screen; cannot create old style PO; hidden PO editor; old PO flow gone

The hidden legacy purchase order create/update flow is on only if the purchase_orders setting enable_legacy_po_create_update is true. If the setting is missing, it is off.

## PO expense excess settings and defaults

Lines 248-292

Asked as: allowed excess over PO; expense over PO total; PO overspend tolerance; excess percent setting; PO excess settings

The allowed excess of expenses over a PO total uses the expenses setting po_expense_allowed_excess. Each part is checked alone: percent must be 0 to 100 (entered as 5 for 5%), value must be 0 or more, and mode must be "lesser_of" or "greater_of". A part that is missing or invalid falls back to its default. The defaults are the constants MAX_PURCHASE_ORDER_EXCESS_PERCENT and MAX_PURCHASE_ORDER_EXCESS_VALUE, with mode lesser_of.

## Percent excess config must be between 0 and 100

Lines 273-279

Asked as: excess percent invalid; percent over 100 ignored; enter 5 for percent; percentage not accepted; excess percent rejected

In the expenses config, the percent of allowed PO excess is entered as a whole percentage, such as 5 for 5%. Turbo accepts it only from 0 to 100 inclusive and divides it by 100. Any other or non-numeric value is ignored and the default percent stays.

## Fixed excess value must not be negative

Lines 281-283

Asked as: negative excess amount; fixed PO excess ignored; excess value rejected; excess dollar amount wrong; excess fixed amount refused

The fixed PO excess amount (value) in the expenses config is accepted when it is a number of 0 or more. A negative or non-numeric value is ignored and the default stays. The amount is in the home currency (CAD).

## Excess mode accepts only lesser_of or greater_of

Lines 285-289

Asked as: lesser of greater of; excess mode ignored; PO excess mode wrong; malformed excess setting; excess mode reset

The PO excess mode is used only when it is the text lesser_of or greater_of. Any other value is ignored and the mode stays lesser_of. A malformed excess setting as a whole gives all the defaults.

## Foreign PO currency rate for the excess limit

Lines 294-320

Asked as: PO exchange rate wrong; foreign currency PO limit; currency conversion PO total; rate used for USD PO; PO rate for limit

For a PO in a foreign currency, Turbo takes the rate from the PO's stored totals: approval_total_home divided by approval_total. If either total is not above 0, it uses the currency's current rate. If that rate is missing or 0 or less, it uses 1. A PO in the home currency, or one whose currency cannot be found, is treated as home currency and not converted.

## Excess text in the limit message

Lines 331-371

Asked as: limit message shows percent; excess text wrong; zero total PO message; limit message currency amount; excess message confusing

The excess text reports the percent (such as "5.00%") when the percent limit applies. It reports the fixed amount (such as "$100.00") when the fixed amount applies. On a foreign PO it shows the CAD amount with the PO-currency amount in parentheses. When the PO total is 0, the percent branch is always chosen in lesser_of mode, even though the allowed excess is $0.

## How the PO expense total limit is calculated

Lines 335-378

Asked as: maximum expense on PO; PO limit calculation; how expense limit computed; why limit rounded down; PO expense cap amount

The largest allowed expense total is the PO total plus the excess. The percent limit is the PO total times (1 + percent). The fixed limit is the PO total plus the fixed amount. In lesser_of mode, Turbo uses the fixed amount only if it is strictly less than the percent excess. In greater_of mode, it uses the fixed amount if it is greater than or equal to the percent excess. Otherwise it uses the percent. The limit is floored to the cent.

## Fixed excess is converted on a foreign PO

Lines 335-342

Asked as: fixed excess foreign PO; CAD excess in USD PO; excess conversion wrong; home currency excess limit; excess not converted

The fixed excess amount is in the home currency (CAD). On a foreign PO, Turbo divides it by the PO's home rate before comparing it with the percent and before adding it, so the limit is in the PO's currency. A rate of 0 or less leaves the amount unconverted.

## Number settings accept only JSON numbers

Lines 380-396

Asked as: number setting text ignored; setting stored as string; numeric value not taken; quoted number config; number config rejected

The number coercion (CoerceFloat64) accepts only float64, int and int64 values. A missing value or any other type, such as text, counts as invalid. Settings that use it then fall back to their default.

## No-PO expense limit default and zero value

Lines 398-411

Asked as: expense needs PO over limit; no PO limit amount; zero limit every expense PO; expense without PO threshold; no PO expense cutoff

The no-PO expense limit comes from the expenses setting no_po_expense_limit. If it is missing, not a number or negative, Turbo uses the constant NO_PO_EXPENSE_LIMIT. A value of 0 means every non-exempt expense needs a PO.

## Message when expense editing is disabled

Lines 413-420

Asked as: expense editing disabled; cannot edit expenses now; expense locked message; PO vendor creation blocked; expense changes refused

When expense editing is switched off, Turbo refuses with "expense editing is currently disabled". The detail says "expense, purchase order, and vendor creation and editing is disabled during transition".

## Message when time editing is disabled

Lines 422-429

Asked as: time editing disabled; timesheet submission blocked; time amendment locked; time approval disabled transition; cannot submit timesheet now

When time editing is switched off, Turbo refuses with "time editing is currently disabled". The detail says time entry and time amendment creation/editing, plus time sheet submission/approval, are disabled during transition.
