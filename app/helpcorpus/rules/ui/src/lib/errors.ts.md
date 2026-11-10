<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/errors.ts","sha256":"6dba6669c75bbc8f51a56d2eaf89231cf67accd63da5f1b55e14a03174848913","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/errors.ts

## Which error message Turbo shows when an action fails

Lines 1-3

When the server returns an error, Turbo shows the server's error text (response.error) first. If there is none, it shows the server's message text (response.message). If the response has neither, Turbo shows the fallback message that the screen supplies.
