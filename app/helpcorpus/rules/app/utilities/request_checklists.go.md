<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/utilities/request_checklists.go","sha256":"95d6145087d21df4eaa44f4ee82604604db137677428e04112694f4bd73724ae","model":"claude-sonnet-5-5"} -->

# Rules in app/utilities/request_checklists.go

## Project request checklists are separate for approval and creation

Lines 12-17

A project request checklist belongs to the stage that answers it: approval or creation (RequestChecklistScopeApproval, RequestChecklistScopeCreation). Each scope numbers its own versions.

## Approver or job holder must answer every checklist question true

Lines 19-24

Each project request checklist question has a key and a text. An approver or job holder must answer it true to continue.

## A checklist needs at least one question

Lines 27-34

Turbo rejects a request checklist question set that is not a list. It also rejects an empty set with the message "a request checklist needs at least one question".

## Each checklist question needs a key and text

Lines 36-39

Turbo rejects a checklist question whose key or text is empty or only spaces. The message is "each request checklist question needs a key and text".

## Checklist question keys must be unique

Lines 40-43

Turbo rejects a question set in which one key appears twice. The message is "request checklist question <key> appears twice".

## Latest checklist version is the highest saved, or zero

Lines 48-58

The latest version for a scope is the highest version number saved for it. When no version is saved for the scope, the latest version is zero.

## Current checklist uses the highest version's questions

Lines 61-77

The questions a scope answers now are those of its highest saved version. Turbo checks the shape of that version's questions when it loads them.

## No saved checklist gives version zero and no questions

Lines 61-71

When a scope has no saved checklist, Turbo returns version zero and no questions. A request can still be viewed, but its action cannot proceed.

## A checklist version exists only if it has saved questions

Lines 79-87

Turbo treats a version of a scope as existing when at least one checklist record is saved for that scope and version.
