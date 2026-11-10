<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_pos_count.sql","sha256":"0be0229070ce643fadf5c89d22e69dcd24e77276cbe900dbea3a6fe9e51a7ca8","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_pos_count.sql

## Job purchase order count includes active and under-revision POs

Lines 2-4

The count of a job's purchase orders (POs) includes only POs whose status is 'Active' or that are under revision (the __PO_IS_REVISION__ condition). POs in any other status are not counted. This file does not show how the under-revision condition is defined.

## Purchase order count is limited to one job

Lines 5-5

Turbo counts only the purchase orders (POs) that belong to the job with the given id (po.job = id).

## Optional branch filter on job PO count

Lines 6-6

When a branch is given, Turbo counts only POs in that branch. A null or empty branch applies no filter.

## Optional division filter on job PO count

Lines 7-7

When a division is given, Turbo counts only POs in that division. A null or empty division applies no filter.

## Optional type filter on job PO count

Lines 8-8

When a type is given, Turbo counts only POs of that type. A null or empty type applies no filter.

## Optional user filter on job PO count

Lines 9-9

When a uid is given, Turbo counts only POs with that uid (po.uid). A null or empty uid applies no filter.
