-- Summary of active and under-revision purchase orders for a job with optional
-- filters. A revision has no approved amount yet, so the total leaves it out and
-- revision_count says how many it left out.
SELECT
  SUM(CASE WHEN po.status = 'Active' THEN po.total ELSE 0 END) total_amount,
  SUM(CASE WHEN po.status = 'Active' THEN 0 ELSE 1 END) revision_count,
  MIN(po.date)                            earliest_po,
  MAX(po.date)                            latest_po,
  json_group_array(
    DISTINCT json_object('id', b.id, 'code', b.code)
  )                                        branches,
  json_group_array(
    DISTINCT json_object('id', d.id, 'code', d.code)
  )                                        divisions,
  json_group_array(
    DISTINCT json_object('name', po.type)
  )                                        types,
  json_group_array(
    DISTINCT json_object('id', p.uid, 'name', p.given_name || ' ' || p.surname)
  )                                        names
FROM   purchase_orders po
LEFT   JOIN branches  b ON po.branch  = b.id
LEFT   JOIN divisions d ON po.division = d.id
LEFT   JOIN profiles  p ON po.uid      = p.uid
WHERE  (po.status = 'Active' OR __PO_IS_REVISION__)
  AND  po.job = {:id}
  AND  ({:branch}   IS NULL OR {:branch}   = '' OR po.branch   = {:branch})
  AND  ({:division} IS NULL OR {:division} = '' OR po.division = {:division})
  AND  ({:type} IS NULL OR {:type} = '' OR po.type = {:type})
  AND  ({:uid} IS NULL OR {:uid} = '' OR po.uid = {:uid}); 