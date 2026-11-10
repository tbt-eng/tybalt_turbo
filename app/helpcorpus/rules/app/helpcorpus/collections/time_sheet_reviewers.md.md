<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/time_sheet_reviewers.md","sha256":"e7d32947fc3e3987301fbbc654d4619aa0a3a5c11e88e488506da162da4403cf","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/time_sheet_reviewers.md

## Nobody can list time sheet reviewer records through the API

Lines 7-9

Asked as: cannot list reviewers; list time sheet reviewers API; reviewer records not showing; fetch all reviewers

The time_sheet_reviewers collection has no list rule. Nobody can list its records through the API. Turbo's own server code can still read them.

## Nobody can view one time sheet reviewer record through the API

Lines 11-13

Asked as: cannot view reviewer record; view single reviewer API; open reviewer details; reviewer record forbidden

The time_sheet_reviewers collection has no view rule. Nobody can view a single record through the API. Turbo's own server code can still read it.

## Only the time sheet's approver can add a reviewer

Lines 15-22

Asked as: add reviewer to time sheet; cannot add reviewer; assign reviewer to timesheet; not allowed to add reviewer

A signed-in person can create a time_sheet_reviewers record only if they are the approver (the approver field) of the time sheet named in the record's time_sheet field. Everyone else is refused.

## Nobody can update a time sheet reviewer record through the API

Lines 24-26

Asked as: change reviewed date; edit reviewer record; update reviewer API; set reviewed date

The time_sheet_reviewers collection has no update rule. Nobody can change a record through the API, for example to set the reviewed date. Turbo's own server code can still update it.

## Only the time sheet's approver can remove a reviewer

Lines 28-35

Asked as: remove reviewer from timesheet; delete reviewer cannot; take off reviewer; cannot delete reviewer

A signed-in person can delete a time_sheet_reviewers record only if they are the approver (the approver field) of the time sheet that the record points to. Everyone else is refused.

## A time sheet can list the same reviewer only once

Lines 39-52

Asked as: duplicate reviewer on timesheet; same reviewer twice; added reviewer already there; reviewer added twice

A unique index on time_sheet and reviewer means the same person cannot be added twice as reviewer of the same time sheet. Each record needs a time sheet and a reviewer. The reviewed date is optional.

## Deleting a time sheet deletes its reviewer records

Lines 42-42

Asked as: deleted timesheet reviewers gone; timesheet delete removes reviewers; reviewers vanished after delete; delete timesheet what happens reviewers

The time_sheet relation uses cascadeDelete. When a time sheet is deleted, its reviewer records are deleted too.
