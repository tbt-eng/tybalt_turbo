<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/absorbSuggestions.ts","sha256":"ae5ae7d762176bc2e550fd756a8038a6b9bfa7acfa218ef8e98bdb1e6d17e2ae","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/absorbSuggestions.ts

## Absorb suggestions cover only clients and vendors

Lines 8-69

Asked as: duplicate clients tab; vendors tab missing; suggestions page tabs; reload loses tab

Absorb suggestions exist for two collections: clients and vendors. Each has its own tab on the suggestions page. The tab is kept in the URL, so a reload keeps it.

## Absorb suggestions have three matching levels

Lines 9-78

Asked as: match strictness setting; how loose matching; change match level; fuzzy match threshold

A person can pick a matching level for suggestions: Strict, Normal or Loose (strict, normal, loose). Turbo's server finds the suggestions again on each request for the chosen level.

## What a suggestion bucket and its members show

Lines 19-38

Asked as: duplicate group members; match score meaning; why records grouped; open proposal on client

A bucket (a group of probable duplicates) lists its members from best match first. Each member has a score from 0 to 1 against the suggested target, and the target scores 1. Each member also shows an open absorb proposal that names it, if there is one. For a client, references count its jobs as client or owner. For a vendor, they count its POs and expenses.

## Pending absorb action shown with its absorber

Lines 40-59

Asked as: who merged records; pending merge blank user; undo absorb who did it; waiting merge action

The suggestions response includes the collection's absorb action that waits for a commit or an undo. It gives the target, the absorbed records and who absorbed them. The absorber is blank for an action made before Turbo recorded it.

## Marking Not duplicates hides a bucket for everyone

Lines 56-90

Asked as: hide duplicate suggestion; dismissed duplicates count; stop suggesting these records; mark not same

Marking records "Not duplicates" applies to everyone, not only to the person who marked them. Their bucket stays hidden until a new record joins it. The response counts the hidden buckets (dismissed).

## Absorbing from a suggestion uses the normal absorb route

Lines 92-102

Asked as: merge from suggestions; combine duplicates button; absorb suggested group; merge duplicate vendors

Absorbing a bucket from the suggestions page sends the records to the target through the normal absorb route of that collection. It does not use a separate suggestion-only action.
