SELECT
  r.id,
  r.expense,
  r.status,
  r.outcome,
  r.note,
  r.created,
  r.requester,
  r.decided_at,
  r.decision_comment,
  r.employee_name,
  r.po_number,
  r.vendor_name,
  r.description,
  r.currency_code,
  r.payment_type,
  CAST(r.total AS REAL) AS total,
  r.committed,
  r.pay_period_ending,
  COALESCE(reason.name, '') AS reason_name,
  COALESCE(TRIM(requester_profile.given_name || ' ' || requester_profile.surname), '') AS requester_name,
  COALESCE(TRIM(decider_profile.given_name || ' ' || decider_profile.surname), '') AS decided_by_name
FROM expense_uncommit_requests r
LEFT JOIN expense_uncommit_reasons reason ON reason.id = r.reason
LEFT JOIN profiles requester_profile ON requester_profile.uid = r.requester
LEFT JOIN profiles decider_profile ON decider_profile.uid = r.decided_by
