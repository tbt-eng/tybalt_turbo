<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/utilities/config.go","sha256":"156d9030feb420047916524039d556305019d6cbdbed3b81ec6ae788c5853adc","model":"claude-sonnet-5-5"} -->

# Rules in app/utilities/config.go

## Config lookup returns nothing when a setting group is missing

Lines 16-30

When Turbo finds no app_config row for a setting group (domain key), it treats the group as absent and not as an error. Callers then use their own defaults. A stored value that is not valid JSON is an error (GetConfigValue).

## Yes/no setting falls back to its default

Lines 32-52

A yes/no setting (GetConfigBool) returns the caller's default when the group is missing, the property is missing, or the value is not true or false (for example a string). It returns the stored value only when it is a real boolean.

## Yes/no setting errors on unreadable config

Lines 34-38

When the setting group exists but its stored value cannot be parsed as JSON, GetConfigBool returns the default together with an error. It does not return the default silently.

## Job editing switch defaults to on when config missing

Lines 54-59

Turbo allows job creation, editing and absorb unless the jobs setting create_edit_absorb is false. If the jobs config, the property or a Boolean value is missing, Turbo treats editing as enabled. A config that cannot be read as JSON returns an error and not the default.

## Project request approval defaults to not required

Lines 61-83

Turbo requires approval of project requests only when the jobs setting require_job_request_approval is true. It treats a missing jobs config or a missing property as not required. A value that is not a Boolean is an error: "jobs.require_job_request_approval must be Boolean".

## Message when job editing is disabled

Lines 85-92

When job editing is off, Turbo refuses with "job editing is currently disabled". The detail says "job creation and editing is disabled during transition" (code jobs_editing_disabled).

## Expense, PO and vendor editing switch defaults to on

Lines 94-99

Turbo allows creating, editing and absorbing expenses, vendors and purchase orders unless the expenses setting create_edit_absorb is false. A missing config or property means enabled.

## Time editing switch defaults to on

Lines 101-107

Turbo allows time entry creation and editing, and selected time sheet changes, unless the time setting create_edit is false. A missing config or property means enabled.

## Notification features default to off

Lines 109-114

A notification template is enabled only when the notifications config has its template code set to true. If the config or the key is missing, Turbo treats the notification as disabled.

## PO second approver timeout defaults to 24 hours

Lines 116-145

The second-stage timeout for purchase orders is the purchase_orders setting second_stage_timeout_hours. Turbo uses 24 hours when the config or key is missing, the value is zero or negative, or it is not a number.

## Legacy PO create/update flow is off by default

Lines 147-152

The hidden legacy purchase order create/update flow is enabled only when purchase_orders.enable_legacy_po_create_update is true. If it is missing, Turbo treats it as disabled.

## PO expense excess settings and defaults

Lines 169-194

The expenses setting po_expense_allowed_excess controls how far expenses may exceed a PO total, using percent, value and mode. Turbo starts from the default constants (MAX_PURCHASE_ORDER_EXCESS_PERCENT and MAX_PURCHASE_ORDER_EXCESS_VALUE) with mode lesser_of. Each field is overridden only if valid.

## Validation of excess percent, value and mode

Lines 196-210

The excess percent is accepted only from 0 to 100 inclusive, and is entered as a percentage (5 means 5%). The excess value is accepted only if it is 0 or more. The mode is accepted only as "lesser_of" or "greater_of". Invalid fields keep their defaults.

## Excess percent is stored as a percentage and converted

Lines 196-200

A percent setting is entered as a percentage, such as 5 for 5%, and Turbo converts it to a fraction (0.05). Turbo accepts values from 0 to 100, both included. A setting outside this range keeps the default percent.

## Excess value must be zero or more

Lines 202-204

Turbo accepts a dollar excess value of 0 or more. A negative or non-numeric value keeps the default value.

## Excess text shown in limit messages

Lines 215-253

The excess text names the rule applied, either a dollar amount such as "$100.00" or a percentage such as "5.00%". The limit is converted to whole cents by CurrencyLimitToCents. With a PO total of 0, the percent branch is chosen under lesser_of, so the text shows a percentage although the allowed excess is $0.

## Expense limit on a PO: lesser_of or greater_of

Lines 223-247

Turbo's maximum expense total for a PO is the PO total plus either the percent excess or the dollar excess. In lesser_of mode (the default), it uses the dollar excess only when it is strictly less than the percent excess; on a tie it uses the percent. In greater_of mode, it uses the dollar excess when it is equal to or greater than the percent excess.

## Numeric settings accept only numbers

Lines 255-271

When Turbo reads a numeric setting (CoerceFloat64), it accepts numbers only. A missing, null or non-number value (such as text) is rejected, and the caller then uses its default.

## No-PO expense limit default and zero meaning

Lines 273-286

The no-PO expense limit comes from the expenses setting no_po_expense_limit. Turbo falls back to the constant NO_PO_EXPENSE_LIMIT when the config or key is missing, the value is not a number, or it is negative. A value of 0 means every non-exempt expense needs a PO.

## No-PO expense limit ignores a negative or non-numeric value

Lines 277-286

The no-PO expense limit comes from the expenses setting no_po_expense_limit. Turbo uses it when it is a number of 0 or more. A negative, non-numeric or missing value, or a missing expenses setting, gives the default NO_PO_EXPENSE_LIMIT.

## Message when expense editing is disabled

Lines 288-295

When expense editing is off, Turbo refuses with "expense editing is currently disabled". The detail says expense, purchase order and vendor creation and editing is disabled during transition (code expenses_editing_disabled).

## Message when time editing is disabled

Lines 297-304

When time editing is off, Turbo refuses with "time editing is currently disabled". The detail says time entry and time amendment creation/editing, plus time sheet submission/approval, are disabled during transition (code time_editing_disabled).
