<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/utilities/branch_claims.go","sha256":"0210faf54b50efff3a2a01f85ac54daefe47b706950087032672db0a58658b59","model":"claude-sonnet-5-5"} -->

# Rules in app/utilities/branch_claims.go

## Record branch check is skipped when the branch is inherited

Lines 30-36

Asked as: branch inherited from job; why no branch check; skip branch permission; branch taken from parent

When a record's branch is not empty and is inherited, Turbo does not check the user's branch claims for it. Otherwise Turbo applies the branch claim check to the record's branch (field "branch"). A record with no branch goes to the claim check, which passes an empty branch.

## A job with the same branch makes the branch inherited

Lines 39-44

Asked as: charging time to job; job on restricted branch; can't log time to job; job branch mismatch

A record that has a job counts as inheriting its branch when the job's branch equals the record's branch. Charging time or costs to a job on a restricted branch then needs no branch claim. If the record's branch differs from the job's branch, the claim check still applies. If the job cannot be found, the branch is not inherited from it.

## Purchase order in a group keeps its original branch

Lines 45-50

Asked as: change branch on purchase order; PO branch locked in group; can't edit PO branch; grouped PO branch change

An existing purchase order (not a new one) in a PO group counts as inheriting its branch when the branch is unchanged from the saved record (Original). Changing the branch of such a PO triggers the claim check. A new PO in a group does not get this exemption.

## Expense against a group purchase order inherits the PO branch

Lines 51-58

Asked as: expense on group PO; expense linked to purchase order; PO group expense branch; expense branch doesn't match PO

An expense linked to a purchase order counts as inheriting its branch when that PO belongs to a PO group and the PO's branch equals the expense's branch. An expense whose branch differs from the PO's branch, or whose PO has no group, is not exempt this way.

## Empty branch passes the branch claim check

Lines 64-68

Asked as: no branch selected; leave branch blank; branch field empty; save without branch

When the branch is empty or only spaces, Turbo does not check branch claims for the user and allows the action.

## Unknown branch is refused

Lines 70-82

Asked as: branch not found; invalid branch error; branch doesn't exist; can't find my branch

When the branch cannot be found, Turbo refuses with "specified branch could not be found" (code invalid_branch) on the branch field.

## Misconfigured branch claim list is refused

Lines 84-96

Asked as: branch claims error; branch misconfigured; allowed claims broken; branch claim list invalid

When a branch's allowed claims (allowed_claims) cannot be read as a list, Turbo refuses with "branch claim requirements are misconfigured" (code branch_claim_config_invalid). A blank value or null counts as no claims. Blank entries in the list are ignored.

## Branch with no allowed claims is open to all

Lines 97-99

Asked as: branch open to everyone; anyone can use branch; branch has no restrictions; branch without permissions

A branch with no allowed claims (allowed_claims) is unrestricted, so any user can use it.

## Failure to load user claims is refused

Lines 101-120

Asked as: can't check my permissions; claims won't load; permission check failed; validate branch claims error

When Turbo cannot load the user's claims, it refuses with "unable to validate branch claim requirements" (code branch_claim_check_failed).

## One matching claim is enough for a restricted branch

Lines 122-126

Asked as: need only one claim; any one permission enough; partial claims branch access; one permission unlocks branch

A user may use a restricted branch when they hold at least one of the branch's allowed claims. They do not need all of them.

## Missing branch claim message lists required claims

Lines 128-152

Asked as: requires at least one of these claims; which claim do I need; branch access denied; why can't I use branch

When the user holds none of the allowed claims, Turbo refuses with "<branch name> requires at least one of these claims: <claims>" (code branch_claim_required). The claims are listed by name in alphabetical order, separated by commas. A claim with no name or that cannot be found is listed by its ID.

## Blank entries in a branch's allowed claims are ignored

Lines 155-177

Asked as: blank claim in branch; empty entry in claims list; whitespace claim ignored; branch claims with blanks

When Turbo reads a branch's allowed claims (allowed_claims), it drops blank or whitespace-only entries. A branch whose list holds only blank entries counts as having no allowed claims, so every user can use it. A list value that is not valid, such as text that is not a JSON list, makes Turbo refuse with the misconfigured branch claim error.
