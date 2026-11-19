<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/rate_sheet_duplicates.go","sha256":"8e21cb8e04e5a7a7631835d72491634e7badb737e4829d1b88706bf9c3664287","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/rate_sheet_duplicates.go

## Near-match limits: at most 5 roles, 3 sheets shown

Lines 16-21

Asked as: how many roles can differ; similar rate sheets shown; close rate sheet limit; how many similar sheets

When Turbo shows rate sheets that are close to a set of rates, it counts a sheet as close only if it differs in at most 5 roles (RateSheetNearMatchMaxRoles). It shows at most 3 such sheets (RateSheetNearMatchLimit).

## Closest near matches are trimmed by role count and limit

Lines 54-68

Asked as: closest rate sheets list; stops after few matches; sheets with too many differences; near match cutoff

Closest returns near-match rate sheets in order, closest first. It stops at the given limit. It also stops at the first sheet that differs in more roles than the maximum. A sheet that differs in exactly the maximum number of roles is kept.

## Roles missing from the entries can be ignored in near matches

Lines 70-84

Asked as: missing role ignored; role with no rate yet; proposed rates missing role; ignore empty roles in match

IgnoringRolesLeftOut removes from each near sheet's difference count the roles that the sheet has and the entries lack, such as a role with no rate yet in proposed rates. Turbo then sorts the near matches again. Exact matches stay unchanged.

## Rates are compared in whole cents

Lines 99-111

Asked as: rate rounding cents; same rate looks different; float rounding rates; decimal rate mismatch

Turbo rounds each rate and overtime rate to whole cents (rate times 100, rounded) before it compares them. Float rounding on input does not make two sheets differ.

## Rates key identifies identical rates on sheets of any state

Lines 122-139

Asked as: same rates key; rate sheets same key; compare inactive sheet rates; empty rates no key

RateSheetRatesKey gives two sets of entries the same key only when they have the same roles and the same rate and overtime rate in cents for each role. It works on sheets of any state, active or inactive. It gives an empty key for no entries.

## Only active sheets not owned by a request are compared

Lines 147-197

Asked as: which sheets are compared; draft sheet excluded from check; inactive sheets not checked; request draft ignored in matching

Comparison loads only active rate sheets (active = true) and leaves out any sheet that a request owns, such as a request's draft. It also leaves out every revision of the excluded code. A sheet with no entries is still loaded, with no roles.

## Identical means same roles and same rates; differences are counted

Lines 199-240

Asked as: exact match same rates; what counts as different; sheet has extra role; differing roles counted

CompareWithActiveRateSheets compares entries with each active sheet. A role counts as differing if the sheet lacks it or has a different rate or overtime rate. A role that only the sheet has also counts as differing. A sheet with 0 differing roles is an exact match. Every other sheet is a near match. Inactive sheets are never compared. Pass an empty code for a sheet with no code yet. Otherwise the revisions of that code are skipped.

## Lookup of active sheets by known rates

Lines 242-302

Asked as: find sheet by my rates; look up rate sheet; sheets with these rates; overtime rate any match

FindActiveRateSheetsByRates returns the active sheets that have every wanted role at the same rate in cents. If the person gave an overtime rate for a role, the sheet must match it too. If they gave none, any overtime rate matches. Roles not listed are ignored. With no wanted rate, Turbo returns no match. The result is newest code first, then highest revision first.

## Match lists are sorted by difference, then code, then revision

Lines 304-315

Asked as: order of matching sheets; sort rate sheet results; which match comes first; newest revision first

Exact and near matches are sorted by fewest differing roles first, then by code ascending, then by revision ascending.

## Duplicate rate sheet error message

Lines 331-351

Asked as: duplicate rate sheet refused; rates already active error; identical rates error message; sheet already exists same rates

When a rate sheet's rates are identical to an active sheet, Turbo refuses it with "a rate sheet with identical rates is already active". The field error reads "the rates are identical to active rate sheet <label>; use that sheet instead". The label comes from the first match, and the error lists every match (code duplicate_rate_sheet).
