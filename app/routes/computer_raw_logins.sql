-- One page of raw logins for the Raw Logins page, newest first.
-- The report itself is left out: it can be 64 KB, and the page loads it for
-- one row at a time. search is a LIKE pattern, or blank for every row; it
-- matches the computer name and the error.
-- computer_raw_logins_count.sql must keep the same FROM and WHERE.
SELECT
  r.id,
  r.created,
  COALESCE(r.computer_name, '') AS computer_name,
  COALESCE(r.error, '') AS error
FROM computer_raw_logins r
WHERE
  {:search} = ''
  OR COALESCE(r.computer_name, '') LIKE {:search} ESCAPE '\'
  OR COALESCE(r.error, '') LIKE {:search} ESCAPE '\'
ORDER BY r.created DESC, r.id DESC
LIMIT {:limit} OFFSET {:offset};
