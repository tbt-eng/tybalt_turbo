<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/categories.md","sha256":"4e4e45c2e95f611ec9e47e71c0c7ca90dabb711e5908da594349de1191c549ba","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/categories.md

## Who can list and view categories

Lines 7-17

Asked as: see categories without logging in; view category list; who can read categories; categories visible to staff

Any signed-in person can list categories and view a single category. Someone who is not signed in cannot.

## Who can create a category

Lines 19-24

Asked as: add a new category; make category; category not allowed to create; procurement create category

Only a signed-in person who holds the job claim (procurement) can create a category.

## Nobody can update a category through the API

Lines 26-28

Asked as: edit category; change category name; can't update category; rename category API

The categories collection has no update rule, so nobody can update a category through the API. Only a PocketBase superuser or Turbo's own server code can change one.

## Who can delete a category, and when deletion is blocked

Lines 30-44

Asked as: remove category; delete category fails; can't delete category; category in use error

Only a signed-in person with the job claim (procurement) can delete a category. Deletion is blocked while a time entry, a purchase order or an expense still references the category.

## Category name and job requirements

Lines 48-52

Asked as: category name too short; category without job; minimum name length category; category needs job

A category needs a name of at least 3 characters and exactly one job (job).

## Deleting a job deletes its categories

Lines 52-52

Asked as: job deleted categories gone; removing job removes categories; categories vanish with job; cascade delete job

The job field is set to cascade delete, so deleting a job also deletes its categories.

## Category names are unique within a job

Lines 59-61

Asked as: duplicate category name; category name already taken; same category name twice; name exists in job

Two categories of the same job cannot have the same name, because a unique index covers job and name. Different jobs can reuse a name.
