<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_pos_api.go","sha256":"18f77c571d823ca7bc6cb7141dd187c75640de86f00c3df230e20afe1af3e5ae","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_pos_api.go

## Job PO list covers active and under-revision purchase orders

Lines 42-47

The job purchase order list is documented as returning purchase orders that are active or under revision (the 'job_pos.sql' query applies the status filter). The list can be filtered by division, type and user.

## Job purchase order list requires a job id

Lines 49-52

The list of purchase orders (POs) for a job refuses a request that has no job id. Turbo shows the message "id is required".

## Job PO list filters

Lines 54-58

The job PO list can be filtered by division, branch, PO type and user (uid). The filter values go to the job_pos query. The query file decides how an empty filter is treated, and this file does not show it. The code comment says the list holds active and under-revision POs, but the query file sets that.

## Job PO list default page and page size

Lines 60-75

The job PO list returns page 1 with 50 POs per page unless the request asks for something else. A page or limit that is not a whole number above 0 is ignored, and Turbo uses the default.

## Job PO list page size capped at 200

Lines 66-74

The job PO list returns at most 200 POs per page. A requested limit above 200 becomes 200.

## Job PO list reports total and page count

Lines 87-105

The job PO list response gives the page, the limit, the total count of matching POs and the number of pages. The number of pages is the total divided by the limit, rounded up. Each row shows the PO number, status, date, total, type, branch, division and the surname and given name of the person.
