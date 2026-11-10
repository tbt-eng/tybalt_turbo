<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/utilities/payroll_placeholders.go","sha256":"46b13308dbaa6d68db72538688555bda688fefc9f335be053e1e0bc5bd378071","model":"claude-sonnet-5-5"} -->

# Rules in app/utilities/payroll_placeholders.go

## Generated placeholder payroll ID format

Lines 5-16

Asked as: payroll ID looks fake; 9 digit payroll number; placeholder payroll number; why is my payroll ID wrong

Turbo treats a payroll ID as a generated placeholder only when it is exactly 9 characters long. The first character must be 9 and all 9 characters must be digits (generatedPayrollPlaceholderGlob). A payroll ID of any other length or pattern is not a placeholder.

## Empty payroll ID is not a generated placeholder

Lines 10-16

Asked as: blank payroll ID; payroll ID is empty; no payroll number showing; missing payroll ID

When the payroll ID field is empty or null, Turbo reads it as an empty string. An empty string is not 9 characters long, so the placeholder check does not match it (GeneratedPayrollPlaceholderSQLCondition).
