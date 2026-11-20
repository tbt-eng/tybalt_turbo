-- The number of rows computer_logins.sql pages through, with the same FROM
-- and WHERE.
SELECT COUNT(*)
FROM computer_logins l
JOIN computers c ON c.id = l.computer
LEFT JOIN users u ON u.id = l.user
LEFT JOIN profiles p ON p.uid = l.user
WHERE
  ({:computer} = '' OR l.computer = {:computer})
  AND (
    {:search} = ''
    OR COALESCE(c.computer_name, '') LIKE {:search} ESCAPE '\'
    OR COALESCE(c.serial, '') LIKE {:search} ESCAPE '\'
    OR TRIM(COALESCE(l.given_name, '') || ' ' || COALESCE(l.surname, '')) LIKE {:search} ESCAPE '\'
    OR COALESCE(
      NULLIF(TRIM(COALESCE(p.given_name, '') || ' ' || COALESCE(p.surname, '')), ''),
      NULLIF(u.name, ''),
      u.email,
      ''
    ) LIKE {:search} ESCAPE '\'
  );
