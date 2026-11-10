<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/app_config.md","sha256":"7409644e77aa45edde9fe1b5f57582d95b7af7bcca048dfd975726cf38c51b80","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/app_config.md

## Signed-in staff can list and view app configuration records

Lines 7-17

Asked as: list app settings; view config records; who can read settings; see app configuration

Any signed-in person can list the app_config records and view one record through the API. A request with no signed-in person is refused.

## Nobody can create app configuration records through the API

Lines 19-21

Asked as: add a new setting; create config record; cannot add app config; insert setting via api

Nobody can create an app_config record through the API. Only a PocketBase superuser can, and Turbo's own server code can still do it.

## Nobody can update app configuration records through the API

Lines 23-25

Asked as: change a setting value; edit app config; update config via api; modify saved settings

Nobody can update an app_config record through the API. Only a PocketBase superuser can, and Turbo's own server code can still do it.

## Nobody can delete app configuration records through the API

Lines 27-29

Asked as: remove a setting; delete config record; cannot delete setting; get rid of config entry

Nobody can delete an app_config record through the API. Only a PocketBase superuser can, and Turbo's own server code can still do it.

## Each app configuration key is unique

Lines 33-46

Asked as: duplicate setting key; key already exists; same config key twice; required fields for setting

Two app_config records can't have the same key (unique index on key). The key and the value (a JSON value) are required on every record. The description is optional.
