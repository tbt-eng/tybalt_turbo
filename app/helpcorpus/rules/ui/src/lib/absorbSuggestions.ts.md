<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/absorbSuggestions.ts","sha256":"2dfc71e719c5eda11bd0f9d7359821fc0707bdabf5d0c46ff79276ee845574bd","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/absorbSuggestions.ts

## Absorb suggestions cover only clients and vendors

Lines 1-11

Asked as: duplicate clients and vendors; merge suggestions for vendors; which clients are duplicates; similar records to combine

Absorb suggestions exist for two collections: clients and vendors. Each suggestion is a bucket of probable duplicates, listed best first, for the users who absorb them.

## Three matching levels for absorb suggestions

Lines 9-17

Asked as: strict or loose matching; match level for duplicates; how fuzzy the duplicate search is; change duplicate match sensitivity

A user picks one of three levels for matching duplicates: Strict, Normal or Loose (strict, normal, loose). The server finds the suggestions again on each request, using the chosen level.

## What an absorb suggestion member shows

Lines 19-29

Asked as: match score on duplicate record; what shows for each suggested record; alias of possible duplicate; linked jobs and POs for duplicate

Each record in a bucket shows its name, alias and a match score from 0 to 1 (the target scores 1). It also shows its references: for a client, the jobs where it is the client or the owner; for a vendor, its POs and expenses. A member also shows any open absorb proposal that names it, or an empty value if there is none.

## Bucket confidence, moved references and rank

Lines 31-38

Asked as: confidence of duplicate group; how many references would move; order of duplicate groups; which record is the suggested keeper

Each bucket has a suggested target record, a confidence value, a count of references that would move (moved_references) and a rank. The rank orders the buckets, best first.

## Pending absorb action waits for commit or undo

Lines 40-57

Asked as: absorb action not finished; undo my merge; who absorbed these records; commit pending merge

The suggestions response can show one pending absorb action for the collection. It waits for a commit or an undo. It lists the target and the absorbed records. Turbo does not record who made the action.

## Dismissed buckets are hidden and counted

Lines 54-55

Asked as: hidden duplicate groups; where did my dismissed duplicates go; count of hidden matches; groups marked not duplicates missing

Buckets that someone marked "Not duplicates" do not appear in the list. The response gives the number of hidden buckets (dismissed).

## Suggestions page keeps the tab in the URL

Lines 63-67

Asked as: page resets to clients tab; keep vendors tab after refresh; suggestions tab lost on reload; tab stays after refresh

The absorb suggestions page (/absorb/suggestions) puts the selected tab, clients or vendors, in the URL (tab=). A reload keeps the tab.

## Not duplicates mark applies to everyone

Lines 78-88

Asked as: hid group for everyone; my not duplicates mark affects others; why did my colleague's duplicate group vanish; group stays hidden for all users

Marking records "Not duplicates" hides their bucket for all users, not only the person who marked it. The bucket stays hidden until a new record joins it.

## Absorbing from suggestions uses the normal absorb route

Lines 90-100

Asked as: merge duplicates from suggestions; which route absorbs suggested records; absorb button separate process; send selected records to merge

Absorbing from a suggestion sends the selected records (ids_to_absorb) to the normal absorb route of the target client or vendor. It has no separate absorb process.

## Check a suggested record in its details page

Lines 102-108

Asked as: open suggested record to check; view duplicate details before merging; new tab link for duplicate; inspect candidate record first

A suggested record links to its details page, so a user can check it in a new tab before absorbing.
