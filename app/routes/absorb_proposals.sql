SELECT
  p.id,
  p.collection_name,
  p.target_id,
  p.target_name,
  p.client_id,
  p.client_name,
  p.requester,
  p.reason,
  p.status,
  p.decided_at,
  p.decision_comment,
  p.created,
  COALESCE(TRIM(requester_profile.given_name || ' ' || requester_profile.surname), '') AS requester_name,
  COALESCE(TRIM(decider_profile.given_name || ' ' || decider_profile.surname), '') AS decided_by_name,
  EXISTS (SELECT 1 FROM absorb_actions a WHERE a.proposal = p.id) AS in_progress
FROM absorb_proposals p
LEFT JOIN profiles requester_profile ON requester_profile.uid = p.requester
LEFT JOIN profiles decider_profile ON decider_profile.uid = p.decided_by
