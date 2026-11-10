<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/currencies.go","sha256":"fd44be4e8c5d21bd0c07882f072f8b5ebd7ad1b7db6460d5c6806a38c8e049cb","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/currencies.go

## Admin claim required for currency admin actions

Lines 42-68

Currency initialization status, rates reload, CAD backfill and currency delete require an admin (the admin claim). A request with no login is refused with "unauthorized". A logged-in user without the claim is refused with "admin claim required".

## Currency initialization status is admin-only and counts blank currencies

Lines 70-111

Only an admin (the admin claim) can view the currency initialization status. Turbo reports whether the home currency exists (home_currency_exists, home_currency_ready) and gives its ID. It also counts purchase orders and expenses whose currency is blank (blank_purchase_orders, blank_expenses).

## Only admins can reload currency rates; result shows updated and skipped counts

Lines 113-130

Only an admin (the admin claim) can reload currency rates. Turbo runs the rate sync (SyncCurrencyRatesWithResult) and returns how many rates it updated (updated) and how many it skipped because the stored rate is newer (skipped_newer). If the sync fails, Turbo shows 'failed to reload currency rates'.

## CAD backfill requires a currency ID

Lines 138-141

When an admin starts the CAD initialization backfill without a currency ID, Turbo refuses with 'currency id is required'. Turbo trims spaces from the ID first.

## CAD backfill refuses an unknown currency

Lines 143-146

When the currency ID does not match a currency, Turbo refuses the backfill with 'currency not found'.

## CAD backfill only works on the CAD currency row

Lines 147-149

The backfill runs only if the chosen currency's code equals the home currency code (CAD). The comparison ignores case and surrounding spaces. Any other currency gets 'initialization requires the CAD currency row'.

## CAD backfill is all-or-nothing

Lines 151-200

The backfill runs in one transaction. If any step fails, Turbo shows 'failed to initialize CAD backfill'. Otherwise it returns the refreshed initialization status.

## CAD backfill fills blank currencies on purchase orders and expenses

Lines 153-176

The backfill sets the chosen CAD currency on every purchase order and every expense whose currency is blank.

## CAD backfill copies amounts into home-currency totals

Lines 161-188

The backfill sets a purchase order's home approval total (approval_total_home) to its approval total (approval_total) when the home total is empty or zero, or when the order's currency is the chosen CAD currency. It sets an expense's settled total (settled_total) to its total when the settled total is empty or zero and the currency is blank or CAD. Non-CAD expenses keep their settled total.

## Currency list needs sign-in and is ordered by sort value

Lines 212-240

Any signed-in user can list currencies; without sign-in Turbo answers 'unauthorized'. Each currency shows its rate, rate date and whether purchase orders or expenses use it. The list is ordered by display order (ui_sort), then by code. A currency without a sort value goes last. A missing rate shows as 0.

## Deleting a currency needs admin and a currency ID

Lines 244-251

Only an admin (the admin claim) can delete a currency. Without a currency ID, Turbo refuses with 'currency id is required'.

## A currency in use cannot be deleted

Lines 253-272

Turbo refuses to delete a currency that any purchase order or expense uses, and shows 'cannot delete a currency that is already in use'. Turbo runs this check before it looks for the currency. An unknown ID that nothing uses gets 'currency not found'.

## Currency code is uppercased and symbol is trimmed on save requests

Lines 282-290

When Turbo reads a currency request, it trims and uppercases the code and trims the symbol.
