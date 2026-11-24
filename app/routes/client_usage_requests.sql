-- The job requests that use one contact or one invoicing profile of a client,
-- with the same roles as client_usage_jobs.sql. The creator and approver let
-- the handler leave out requests the caller may not read.
WITH profiles AS (
  SELECT id FROM client_invoicing_information
  WHERE client = {:client}
    AND (id = {:profile} OR ({:contact} != '' AND contact = {:contact}))
)
SELECT
  r.id,
  r.description,
  r.status,
  r.creator,
  COALESCE(r.approver, '') AS approver,
  ({:contact} != '' AND r.contact = {:contact}) AS project_contact,
  (r.invoicing_information IN (SELECT id FROM profiles)) AS invoice_recipient
FROM job_requests r
WHERE ({:contact} != '' AND r.contact = {:contact})
  OR r.invoicing_information IN (SELECT id FROM profiles)
ORDER BY r.created DESC, r.id
