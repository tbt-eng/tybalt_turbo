<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/notifications.md","sha256":"5e3397e01f0480d25981ac8491ca395df0d598dc1a2288b1031bad3b568fd9b3","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/notifications.md

## Nobody can list notifications through the API

Lines 7-9

Asked as: can't see notifications list; list notifications api; fetch all notifications; why can't I read notifications

The notifications collection has no list rule. Nobody can list notification records through the API, except a PocketBase superuser. Turbo's own server code can still read them.

## Nobody can view a notification through the API

Lines 11-13

Asked as: view single notification; open one notification record; get notification by id; notification access denied

The notifications collection has no view rule. Nobody can view a single notification record through the API, except a PocketBase superuser. Turbo's own server code can still read them.

## Nobody can create a notification through the API

Lines 15-17

Asked as: create notification via api; post a notification; add notification from client; insert notification manually

The notifications collection has no create rule. Nobody can create a notification through the API, except a PocketBase superuser. Turbo's own server code can still create them.

## Nobody can update a notification through the API

Lines 19-21

Asked as: edit a notification; change notification status by api; update notification record; patch notification api

The notifications collection has no update rule. Nobody can update a notification through the API, except a PocketBase superuser. Turbo's own server code can still update them.

## Nobody can delete a notification through the API

Lines 23-25

Asked as: remove a notification; delete notification api; clear notifications from api; can't delete notification

The notifications collection has no delete rule. Nobody can delete a notification through the API, except a PocketBase superuser. Turbo's own server code can still delete them.

## Notification status allows only four values

Lines 34-35

Asked as: invalid notification status; notification status value; notification missing recipient; notification status_updated required

A notification's status (status) must be one of pending, inflight, sent or error. A notification must also have a recipient, a template and a status_updated date.
