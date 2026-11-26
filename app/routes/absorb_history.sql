SELECT
  h.id,
  h.collection_name,
  h.target_id,
  h.target_name,
  h.absorbed_records,
  h.updated_references,
  h.proposal,
  h.absorbed_by,
  h.absorbed_at,
  h.committed_by,
  h.created AS committed_at,
  COALESCE(TRIM(absorber.given_name || ' ' || absorber.surname), '') AS absorbed_by_name,
  COALESCE(TRIM(committer.given_name || ' ' || committer.surname), '') AS committed_by_name
FROM absorb_history h
LEFT JOIN profiles absorber ON absorber.uid = h.absorbed_by
LEFT JOIN profiles committer ON committer.uid = h.committed_by
