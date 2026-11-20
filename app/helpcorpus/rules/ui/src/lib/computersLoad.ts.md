<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/computersLoad.ts","sha256":"41300878ffa7a866f6dd018dcf195864030c5de9823da14cd809646d133418de","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/computersLoad.ts

## Computers pages need the computers or admin claim

Lines 6-13

Asked as: can't open computers page; computers page access denied; need computers claim; missing admin claim; computers message says claim required

Turbo stops a person from opening a Computers page unless they have the computers claim or the admin claim. It then shows the message "The computers or admin claim is required to see computers." The check uses the result of canViewComputers (defined in another file) on the person's claims. The server also refuses the computer data to that person.
