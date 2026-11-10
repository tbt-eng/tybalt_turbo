<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/currencies.md","sha256":"3a01e90237b79f9e399c50b665cf6222fa451ba3ff3c775d8b9ab92032c41f21","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/currencies.md

## Signed-in staff can list and view currencies

Lines 7-17

Asked as: view currency list; see currencies as staff; currency api read access; list currencies not signed in

Any signed-in person can list the currencies collection and view one currency record through the API. A request with no signed-in person fails the filter (`@request.auth.id != ""`).

## Only admins can create, update or delete currencies

Lines 19-35

Asked as: who can add currency; cannot edit currency; admin only currency changes; delete currency permission denied

Only a person with the admin claim can create, update or delete a currency record through the API. Other signed-in people can only read currencies. Turbo's own server code is not limited by these rules.

## Currency code must be three capital letters and unique

Lines 42-55

Asked as: currency code rejected; duplicate currency code; three letter code error; currency code already exists

A currency's code (`code`) is required and must be exactly three capital letters A-Z, such as USD. Two currencies cannot have the same code, because of a unique index on `code`.

## Currency symbol format

Lines 43-43

Asked as: currency symbol invalid; which symbol for currency; symbol field error; currency sign format

A currency's symbol (`symbol`) is required. It must be one currency symbol character (Unicode class Sc, such as $) or three capital letters.

## Currency icon must be one SVG file

Lines 44-44

Asked as: upload currency flag image; currency icon not accepted; svg upload for currency; currency logo file type

A currency's icon (`icon`) is required. It accepts one file only, and the file must be an SVG image (image/svg+xml).

## Currency rate is required; rate date and sort order are optional

Lines 45-47

Asked as: currency exchange rate missing; rate field required; currency rate date blank; sort currencies order

A currency's rate (`rate`) is required. The rate date (`rate_date`) and the sort order (`ui_sort`) are optional.
