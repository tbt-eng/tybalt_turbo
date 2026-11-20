<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/computers.sql","sha256":"021dd862f8cf6e4e30135bab5f21933c7fdf6a14fb2f58a0db0c4b13975dbf8c","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/computers.sql

## Computer names come from a query, not from profile access

Lines 1-4

Asked as: names not showing on computers; computer page missing user name; last user blank; can't see who has the computer

The Computers pages get the last user's name and the assigned person's name from this query. A holder of the computers claim can read computers but not other people's profiles, so Turbo does not use an expand for these names.

## Which computers the query returns

Lines 5-57

Asked as: retired computers list; show retired machines; computer lookup by id; which computers are returned

When an id is given, Turbo returns only that computer, whether or not it is retired, and ignores the retired setting. With a blank id and retired = 1, it returns the computers that have a retired value. With a blank id and any other retired value, it returns only computers with no retired value.

## Missing computer details appear as blank or zero

Lines 9-44

Asked as: computer field shows 0; blank computer details; network config says null; empty value on computer

A computer with no value in a text field shows a blank. A computer with no value in a number field (system type, RAM, OS SKU, boot drive capacity, boot drive free space, radiator version) shows 0. Retired and assigned time also show blank when empty. A blank or missing network config shows as the text 'null'.

## Last user name fallback order

Lines 29-34

Asked as: last user shows email; wrong last user name; who last used the computer; last user name missing

Turbo shows the last user's name (last_user_name) from the first of these that is not blank. First comes the given name and surname on that person's profile. Then comes the name on the user record, then the user's email. If none exists, the name is blank.

## Assigned person name fallback order

Lines 37-42

Asked as: assigned person shows email; assignee name blank; computer owner name wrong; who the computer is assigned to

Turbo shows the name of the person a computer is assigned to (assigned_to_name) from the first of these that is not blank. First comes the given name and surname on that person's profile. Then comes the name on the user record, then the user's email. If none exists, the name is blank.

## Computer list sort order

Lines 58-62

Asked as: computers out of order; retired list order; sort computers by name; most recently retired first

The retired list shows the most recently retired computers first. Computers with the same retired value, and all other lists, sort by computer name and then by id.
