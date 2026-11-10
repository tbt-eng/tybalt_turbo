<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/claims_list.sql","sha256":"9d0f3b6e3ff87b80587c5641f3f9487edf7ec8af5f2756add63d3723e85d6695","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/claims_list.sql

## Claims list shows how many people hold each claim

Lines 1-8

Asked as: how many users have a claim; claim holder count; number of people with permission; claims list zero holders

The claims list shows each claim's id, name and description, with a count of its holders (holder_count). The count comes from the user_claims rows that match the claim. A claim that nobody holds still appears, with a count of 0.

## Claims list is sorted by claim name

Lines 9-9

Asked as: claims order alphabetical; sort claims by name; claims list sorting; why are claims in this order

The claims list is ordered by claim name (c.name). The order does not depend on the holder count.
