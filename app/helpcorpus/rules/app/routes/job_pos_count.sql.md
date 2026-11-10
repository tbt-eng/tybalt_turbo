<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_pos_count.sql","sha256":"0be0229070ce643fadf5c89d22e69dcd24e77276cbe900dbea3a6fe9e51a7ca8","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_pos_count.sql

## Job purchase order count includes active and under-revision POs

Lines 2-4

Asked as: PO count wrong for job; purchase orders not counted; draft POs missing from count; job PO count status

The count of a job's purchase orders (POs) includes only POs whose status is 'Active' or that are under revision (the __PO_IS_REVISION__ condition). POs in any other status are not counted. This file does not show how the under-revision condition is defined.

## Purchase order count is limited to one job

Lines 5-5

Asked as: POs from other jobs counted; count shows POs of another job; job PO total wrong; purchase orders wrong job

Turbo counts only the purchase orders (POs) that belong to the job with the given id (po.job = id).

## Optional branch filter on job PO count

Lines 6-6

Asked as: count only one branch; PO count by branch; branch filter purchase orders; POs for my branch

When a branch is given, Turbo counts only POs in that branch. A null or empty branch applies no filter.

## Optional division filter on job PO count

Lines 7-7

Asked as: count per division; division PO total; filter POs by division; division purchase orders

When a division is given, Turbo counts only POs in that division. A null or empty division applies no filter.

## Optional type filter on job PO count

Lines 8-8

Asked as: count by PO type; filter purchase orders type; only one type of PO; PO type count

When a type is given, Turbo counts only POs of that type. A null or empty type applies no filter.

## Optional user filter on job PO count

Lines 9-9

Asked as: POs created by a user; count of my purchase orders; user PO count; purchase orders by person

When a uid is given, Turbo counts only POs with that uid (po.uid). A null or empty uid applies no filter.
