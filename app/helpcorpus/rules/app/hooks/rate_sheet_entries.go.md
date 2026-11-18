<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/rate_sheet_entries.go","sha256":"3ff6faa43fb347aec36f6d5f0c12940a1ca06419fb723ae995990626eb6e6ec9","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/rate_sheet_entries.go

## Overtime rate check tolerates float rounding error

Lines 19-25

Asked as: overtime rate rounding error; float rounding overtime rate; 1.255 overtime rate refused; overtime rate cents check

Turbo counts an overtime rate as a whole number of cents when it is within a millionth of a cent of a whole cent. So 1.15 passes despite float error, and 1.255 is refused. The check multiplies the value by 100 and compares it with the rounded result.

## Overtime rate limited to two decimal places on rate sheet entry create

Lines 27-43

Asked as: overtime rate too many decimals; overtime rate decimal places error; validation_too_many_decimals; overtime_rate two decimal places

Turbo refuses a new rate sheet entry created through the collection API when its overtime rate (overtime_rate) has more than two decimal places. The message is "overtime_rate must have at most two decimal places" (code validation_too_many_decimals). The error is attached to the overtime_rate field.

## Other rate sheet entry paths check overtime rate themselves

Lines 27-31

Asked as: rate sheet entry correction decimals; proposed rate overtime decimals; overtime rate update not checked; edit rate sheet entry overtime

The create hook runs only for creates through the collection API. The create route, the entry correction route and proposed rates each check the overtime rate decimals themselves. The collection has no update rule, so this hook does not cover updates.
