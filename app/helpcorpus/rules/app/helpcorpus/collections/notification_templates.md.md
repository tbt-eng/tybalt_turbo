<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/notification_templates.md","sha256":"27ce5ef55ca8095ec653dc4d5e850246c6a747fa0f3d6fbca8adc92a5452b724","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/notification_templates.md

## Nobody can list notification templates through the API

Lines 7-9

Asked as: list notification templates; can't see templates in API; templates missing from API list; fetch all templates api

The notification_templates collection has no list rule. No one can list its records through the API. Turbo's own server code can still read them.

## Nobody can view a notification template through the API

Lines 11-13

Asked as: open a notification template; view single template; template api returns forbidden; read one template record

The notification_templates collection has no view rule. No one can open a single record through the API. Turbo's own server code can still read it.

## Nobody can create a notification template through the API

Lines 15-17

Asked as: create new notification template; add a template via api; can't add template; post template api denied

The notification_templates collection has no create rule. No one can create a record through the API. Turbo's own server code can still create records.

## Nobody can update a notification template through the API

Lines 19-21

Asked as: edit notification template; change template through api; update template not allowed; template changes rejected

The notification_templates collection has no update rule. No one can change a record through the API. Turbo's own server code can still change records.

## Nobody can delete a notification template through the API

Lines 23-25

Asked as: remove notification template; delete template api; can't delete template; template deletion blocked

The notification_templates collection has no delete rule. No one can delete a record through the API. Turbo's own server code can still delete records.

## Notification template code format and minimum length

Lines 32-32

Asked as: template code rejected; invalid template code; code with capitals or digits; underscore error in code

A notification template's code is required and has at least 3 characters. It must be lowercase letters only, with words joined by single underscores, for example "po_approved". Digits, capitals, and a leading, trailing or doubled underscore are refused.

## Notification template codes must be unique

Lines 44-44

Asked as: duplicate template code; template code already exists; code taken for template; same code twice

Two notification templates can't have the same code, because a unique index covers the code field.
