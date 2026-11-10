<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_po_summary.sql","sha256":"db804418406d9049e9b21f718b44cfa06f0f35b2c7902885929ebcd720f71f77","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_po_summary.sql

## Total amount adds up active purchase orders only

Lines 5-5

The summary total (total_amount) adds the total of each PO with status 'Active'. A PO under revision adds 0 to the total, because it has no approved amount yet.

## Revision count counts POs that are not active

Lines 6-6

The summary revision_count counts each PO in the summary whose status is not 'Active'. These are the POs under revision, and the total leaves them out.

## Summary shows earliest and latest PO dates

Lines 7-8

The summary gives the earliest date (earliest_po) and the latest date (latest_po) of the POs it includes. Both active and under-revision POs count here.

## Summary lists the distinct branches, divisions, types and people of the POs

Lines 9-24

The summary lists each distinct branch (id and code), division (id and code), PO type (name) and PO owner (uid and given name plus surname) among the included POs. A PO with no branch or division still counts, because those lookups are left joins.

## Job PO summary covers only active and under-revision purchase orders

Lines 25-26

The job purchase order summary counts only purchase orders (POs) of the job whose status is 'Active' or that are under revision (the __PO_IS_REVISION__ condition, defined outside this file). POs in any other status are left out of every figure.

## Optional filters narrow the job PO summary

Lines 27-30

The summary can be filtered by branch, division, type and owner (uid). A filter that is empty or not given is ignored. A filter that is given keeps only POs that match it exactly. Filters combine, so a PO must pass all of them.
