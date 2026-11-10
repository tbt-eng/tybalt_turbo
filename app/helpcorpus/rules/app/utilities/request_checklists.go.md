<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/utilities/request_checklists.go","sha256":"95d6145087d21df4eaa44f4ee82604604db137677428e04112694f4bd73724ae","model":"claude-sonnet-5-5"} -->

# Rules in app/utilities/request_checklists.go

## Project request checklists are separate for approval and creation

Lines 12-17

Asked as: approval and creation checklist mixed; same checklist for both stages; checklist showing in wrong stage; separate checklist per stage

A project request checklist belongs to the stage that answers it: approval or creation (RequestChecklistScopeApproval, RequestChecklistScopeCreation). Each scope numbers its own versions.

## Approver or job holder must answer every checklist question true

Lines 19-24

Asked as: can't continue until all boxes ticked; checkbox blocks my approval; must tick every checklist item; answer yes to every question

Each project request checklist question has a key and a text. An approver or job holder must answer it true to continue.

## A checklist needs at least one question

Lines 27-34

Asked as: empty checklist error; checklist saved with no questions; no questions in checklist; checklist has nothing in it

Turbo rejects a request checklist question set that is not a list. It also rejects an empty set with the message "a request checklist needs at least one question".

## Each checklist question needs a key and text

Lines 36-39

Asked as: blank checklist question; question text left empty; checklist item has no label; question missing its name

Turbo rejects a checklist question whose key or text is empty or only spaces. The message is "each request checklist question needs a key and text".

## Checklist question keys must be unique

Lines 40-43

Asked as: duplicate checklist question; same question listed twice; question appears twice error; repeated checklist item key

Turbo rejects a question set in which one key appears twice. The message is "request checklist question <key> appears twice".

## Latest checklist version is the highest saved, or zero

Lines 48-58

Asked as: which checklist version is latest; checklist version number shows zero; how version numbers are picked; newest checklist version number

The latest version for a scope is the highest version number saved for it. When no version is saved for the scope, the latest version is zero.

## Current checklist uses the highest version's questions

Lines 61-77

Asked as: old checklist questions still showing; which questions apply now; newest checklist questions used; checklist questions out of date

The questions a scope answers now are those of its highest saved version. Turbo checks the shape of that version's questions when it loads them.

## No saved checklist gives version zero and no questions

Lines 61-71

Asked as: request has no checklist; can't approve without checklist; no checklist so button disabled; request stuck without questions

When a scope has no saved checklist, Turbo returns version zero and no questions. A request can still be viewed, but its action cannot proceed.

## A checklist version exists only if it has saved questions

Lines 79-87

Asked as: checklist version not appearing; version missing after saving checklist; unsaved checklist version; checklist version disappeared

Turbo treats a version of a scope as existing when at least one checklist record is saved for that scope and version.
