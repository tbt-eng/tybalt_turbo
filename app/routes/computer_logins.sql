-- Every login for the Logins page and the computer details page,
-- newest first, with the computer name and the name of the matched user. A
-- holder of the computers claim can read logins but not the profiles of other
-- people, so the names come from here and not from an expand.
-- computer limits the rows to one computer when it is not blank. search is a
-- LIKE pattern, or blank for every row; it matches the computer name, the
-- serial, the reported name and the name of the matched user as the row shows it
-- (the profile name, else the user name, else the email).
SELECT
  l.id,
  l.created,
  l.computer,
  COALESCE(c.computer_name, '') AS computer_name,
  COALESCE(c.serial, '') AS serial,
  COALESCE(l.given_name, '') AS given_name,
  COALESCE(l.surname, '') AS surname,
  COALESCE(l.user, '') AS user,
  COALESCE(
    NULLIF(TRIM(COALESCE(p.given_name, '') || ' ' || COALESCE(p.surname, '')), ''),
    NULLIF(u.name, ''),
    u.email,
    ''
  ) AS user_name
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
  )
ORDER BY l.created DESC, l.id DESC;
