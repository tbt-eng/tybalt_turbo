-- The jobs that use one contact or one invoicing profile of a client. The
-- caller binds the contact or the profile and leaves the other empty. A job
-- uses a contact as its project contact, through an invoicing profile that
-- sends invoices to the contact, or both; it is listed once with both roles.
-- These are the uses client_details.sql counts.
WITH profiles AS (
  SELECT id FROM client_invoicing_information
  WHERE client = {:client}
    AND (id = {:profile} OR ({:contact} != '' AND contact = {:contact}))
)
SELECT
  j.id,
  j.number,
  j.description,
  j.status,
  ({:contact} != '' AND j.contact = {:contact}) AS project_contact,
  (j.invoicing_information IN (SELECT id FROM profiles)) AS invoice_recipient
FROM jobs j
WHERE ({:contact} != '' AND j.contact = {:contact})
  OR j.invoicing_information IN (SELECT id FROM profiles)
ORDER BY j.number DESC, j.id
