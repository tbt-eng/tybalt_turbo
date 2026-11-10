<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/user_claims_summary.md","sha256":"e49baae81b5568cc556eb731d98793d4247896f5131e9e17ba51f12c84b809f9","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/user_claims_summary.md

## Nobody can create, update or delete claims summary records through the API

Lines 5-5

Asked as: cannot edit claims summary; update claims summary via api; delete claims summary record; read only claims view; create claims summary fails

The claims summary (user_claims_summary) is a view that Turbo only reads. The file sets no create, update or delete rule, so nobody can do these actions through the API. Turbo's own server code is not limited by these rules.

## Who can list the claims summary records

Lines 7-11

Asked as: see other people's claims; list all claims summaries; only my own claims; why only one record; claims list shows just me

A signed-in person can list only their own record in the claims summary view (user_claims_summary). The list rule requires the signed-in person's id to equal the record's id (`@request.auth.id = id`). Nobody sees another person's claims through this list.

## Who can view one claims summary record

Lines 13-17

Asked as: open someone else's claims; view colleague claims record; access another user's claims; claims record forbidden; view claims of other person

A signed-in person can view only their own claims summary record (user_claims_summary). The view rule requires the signed-in person's id to equal the record's id. A request for another person's record does not pass.

## Claims summary lists the claim names of each user

Lines 21-34

Asked as: which claims do I have; list of my claim names; empty claims array; claims names per user; user's claim list

Each user has one record in the view, with the user's id and a JSON array of the names of the claims (claims.name) that the user holds in user_claims. A user with no claims gets an empty array '[]'.
