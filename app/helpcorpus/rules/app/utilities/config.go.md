<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/utilities/config.go","sha256":"156d9030feb420047916524039d556305019d6cbdbed3b81ec6ae788c5853adc","model":"claude-sonnet-5-5"} -->

# Rules in app/utilities/config.go

## Config lookup returns nothing when a setting group is missing

Lines 16-30

Asked as: setting not found; no config row; config group missing; app config empty

When Turbo finds no app_config row for a setting group (domain key), it treats the group as absent and not as an error. Callers then use their own defaults. A stored value that is not valid JSON is an error (GetConfigValue).

## Yes/no setting falls back to its default

Lines 32-52

Asked as: yes no setting wrong; true false setting ignored; boolean setting default used; toggle shows wrong value

A yes/no setting (GetConfigBool) returns the caller's default when the group is missing, the property is missing, or the value is not true or false (for example a string). It returns the stored value only when it is a real boolean.

## Yes/no setting errors on unreadable config

Lines 34-38

Asked as: config json broken; bad json in settings; settings error on toggle; unreadable setting error

When the setting group exists but its stored value cannot be parsed as JSON, GetConfigBool returns the default together with an error. It does not return the default silently.

## Job editing switch defaults to on when config missing

Lines 54-59

Asked as: cannot edit jobs; jobs editing on by default; job creation allowed by default; absorb job switch

Turbo allows job creation, editing and absorb unless the jobs setting create_edit_absorb is false. If the jobs config, the property or a Boolean value is missing, Turbo treats editing as enabled. A config that cannot be read as JSON returns an error and not the default.

## Project request approval defaults to not required

Lines 61-83

Asked as: project request needs approval; job request approval not required; approval for project requests; who approves project requests

Turbo requires approval of project requests only when the jobs setting require_job_request_approval is true. It treats a missing jobs config or a missing property as not required. A value that is not a Boolean is an error: "jobs.require_job_request_approval must be Boolean".

## Message when job editing is disabled

Lines 85-92

Asked as: job editing disabled message; cannot edit jobs during transition; why can't I edit job; job editing blocked

When job editing is off, Turbo refuses with "job editing is currently disabled". The detail says "job creation and editing is disabled during transition" (code jobs_editing_disabled).

## Expense, PO and vendor editing switch defaults to on

Lines 94-99

Asked as: cannot edit expenses; vendor editing blocked; PO editing allowed by default; expense absorb switch

Turbo allows creating, editing and absorbing expenses, vendors and purchase orders unless the expenses setting create_edit_absorb is false. A missing config or property means enabled.

## Time editing switch defaults to on

Lines 101-107

Asked as: cannot edit time entries; time amendment allowed; timesheet change blocked; time entry switch default

Turbo allows time entry creation and editing, and selected time sheet changes, unless the time setting create_edit is false. A missing config or property means enabled.

## Notification features default to off

Lines 109-114

Asked as: notification not sending; template notifications disabled; emails not sent; notifications turned off

A notification template is enabled only when the notifications config has its template code set to true. If the config or the key is missing, Turbo treats the notification as disabled.

## PO second approver timeout defaults to 24 hours

Lines 116-145

Asked as: second approver timeout; PO approval deadline 24 hours; how long purchase order waits; second stage timeout hours

The second-stage timeout for purchase orders is the purchase_orders setting second_stage_timeout_hours. Turbo uses 24 hours when the config or key is missing, the value is zero or negative, or it is not a number.

## Legacy PO create/update flow is off by default

Lines 147-152

Asked as: old PO form missing; legacy purchase order screen; hidden PO create page; legacy PO flow off

The hidden legacy purchase order create/update flow is enabled only when purchase_orders.enable_legacy_po_create_update is true. If it is missing, Turbo treats it as disabled.

## PO expense excess settings and defaults

Lines 169-194

Asked as: expense over PO total; PO excess allowance setting; expense exceeds purchase order; how much over PO allowed

The expenses setting po_expense_allowed_excess controls how far expenses may exceed a PO total, using percent, value and mode. Turbo starts from the default constants (MAX_PURCHASE_ORDER_EXCESS_PERCENT and MAX_PURCHASE_ORDER_EXCESS_VALUE) with mode lesser_of. Each field is overridden only if valid.

## Validation of excess percent, value and mode

Lines 196-210

Asked as: invalid excess percent; excess percent over 100; excess mode wrong value; excess setting rejected

The excess percent is accepted only from 0 to 100 inclusive, and is entered as a percentage (5 means 5%). The excess value is accepted only if it is 0 or more. The mode is accepted only as "lesser_of" or "greater_of". Invalid fields keep their defaults.

## Excess percent is stored as a percentage and converted

Lines 196-200

Asked as: percent 5 means 5; excess percent as fraction; percentage entered wrong; percent setting converted

A percent setting is entered as a percentage, such as 5 for 5%, and Turbo converts it to a fraction (0.05). Turbo accepts values from 0 to 100, both included. A setting outside this range keeps the default percent.

## Excess value must be zero or more

Lines 202-204

Asked as: negative excess value; dollar excess below zero; excess amount invalid; excess value rejected

Turbo accepts a dollar excess value of 0 or more. A negative or non-numeric value keeps the default value.

## Excess text shown in limit messages

Lines 215-253

Asked as: limit message shows percent; excess text zero PO; why percent shown with $0; limit message wording wrong

The excess text names the rule applied, either a dollar amount such as "$100.00" or a percentage such as "5.00%". The limit is converted to whole cents by CurrencyLimitToCents. With a PO total of 0, the percent branch is chosen under lesser_of, so the text shows a percentage although the allowed excess is $0.

## Expense limit on a PO: lesser_of or greater_of

Lines 223-247

Asked as: lesser of greater of; which excess applies; expense limit on PO calculation; tie excess percent or dollar

Turbo's maximum expense total for a PO is the PO total plus either the percent excess or the dollar excess. In lesser_of mode (the default), it uses the dollar excess only when it is strictly less than the percent excess; on a tie it uses the percent. In greater_of mode, it uses the dollar excess when it is equal to or greater than the percent excess.

## Numeric settings accept only numbers

Lines 255-271

Asked as: setting text not number; numeric config rejected; string value for number setting; number setting ignored

When Turbo reads a numeric setting (CoerceFloat64), it accepts numbers only. A missing, null or non-number value (such as text) is rejected, and the caller then uses its default.

## No-PO expense limit default and zero meaning

Lines 273-286

Asked as: no PO expense limit; zero no PO limit meaning; expense needs PO threshold; default PO limit amount

The no-PO expense limit comes from the expenses setting no_po_expense_limit. Turbo falls back to the constant NO_PO_EXPENSE_LIMIT when the config or key is missing, the value is not a number, or it is negative. A value of 0 means every non-exempt expense needs a PO.

## No-PO expense limit ignores a negative or non-numeric value

Lines 277-286

Asked as: negative PO limit; no PO limit missing default; expense limit setting invalid; non numeric PO limit

The no-PO expense limit comes from the expenses setting no_po_expense_limit. Turbo uses it when it is a number of 0 or more. A negative, non-numeric or missing value, or a missing expenses setting, gives the default NO_PO_EXPENSE_LIMIT.

## Message when expense editing is disabled

Lines 288-295

Asked as: expense editing disabled; cannot edit expenses transition; expense creation blocked; expense editing message

When expense editing is off, Turbo refuses with "expense editing is currently disabled". The detail says expense, purchase order and vendor creation and editing is disabled during transition (code expenses_editing_disabled).

## Message when time editing is disabled

Lines 297-304

Asked as: time editing disabled; cannot submit timesheet; time approval blocked message; time entry locked message

When time editing is off, Turbo refuses with "time editing is currently disabled". The detail says time entry and time amendment creation/editing, plus time sheet submission/approval, are disabled during transition (code time_editing_disabled).
