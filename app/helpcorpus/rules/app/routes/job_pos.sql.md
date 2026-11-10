<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_pos.sql","sha256":"dc04254c6b5335fc4c734a49fe2bc3f735d1e3e9bb7e92c0bcc6f4adea300a56","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_pos.sql

## POs with no branch, division or profile still appear with blank names

Lines 2-17

A PO without a branch, division or profile still appears in the job's list. Its branch code, division code, surname and given name show as blank text instead of empty values.

## Job purchase order list shows only active and under-revision POs

Lines 18-19

The job's purchase order list includes only POs with status Active (status 'Active') or POs that count as under revision (the __PO_IS_REVISION__ condition). POs in other statuses do not appear. The list covers one job only (po.job).

## Job PO list filters by branch, division, type and person

Lines 20-23

The list can be narrowed by branch, division, PO type and the person who created the PO (uid). A filter that is missing or blank is ignored. Each filter that is set must match exactly, and all set filters apply together.

## Job PO list shows newest POs first and is paged

Lines 24-25

The list is sorted by PO date, newest first (po.date DESC). It returns one page at a time, using the limit and offset that the request gives.
