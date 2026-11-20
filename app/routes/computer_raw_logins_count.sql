-- The number of rows computer_raw_logins.sql pages through, with the same
-- FROM and WHERE.
SELECT COUNT(*)
FROM computer_raw_logins r
WHERE
  {:search} = ''
  OR COALESCE(r.computer_name, '') LIKE {:search} ESCAPE '\'
  OR COALESCE(r.error, '') LIKE {:search} ESCAPE '\';
