-- The computers for the Computers pages, with the names of the last user and
-- of the person a computer is assigned to. A holder of the computers claim
-- can read computers but not the profiles of other people, so the names come
-- from here and not from an expand.
-- A blank id selects every computer. retired picks the list: 1 for the
-- retired computers, 0 for the others. It is ignored when an id is given.
SELECT
  c.id,
  COALESCE(c.computer_name, '') AS computer_name,
  COALESCE(c.serial, '') AS serial,
  COALESCE(c.mfg, '') AS mfg,
  COALESCE(c.model, '') AS model,
  COALESCE(c.system_type, 0) AS system_type,
  COALESCE(c.ram, 0) AS ram,
  COALESCE(c.os_version, '') AS os_version,
  COALESCE(c.os_sku, 0) AS os_sku,
  COALESCE(c.os_arch, '') AS os_arch,
  COALESCE(c.boot_drive, '') AS boot_drive,
  COALESCE(c.boot_drive_fs, '') AS boot_drive_fs,
  COALESCE(c.boot_drive_cap, 0) AS boot_drive_cap,
  COALESCE(c.boot_drive_free, 0) AS boot_drive_free,
  COALESCE(c.radiator_version, 0) AS radiator_version,
  COALESCE(NULLIF(c.network_config, ''), 'null') AS network_config,
  COALESCE(c.user_given_name, '') AS user_given_name,
  COALESCE(c.user_surname, '') AS user_surname,
  COALESCE(c.email, '') AS email,
  COALESCE(c.upn, '') AS upn,
  COALESCE(c.last_user, '') AS last_user,
  COALESCE(
    NULLIF(TRIM(COALESCE(lp.given_name, '') || ' ' || COALESCE(lp.surname, '')), ''),
    NULLIF(lu.name, ''),
    lu.email,
    ''
  ) AS last_user_name,
  COALESCE(c.last_login, '') AS last_login,
  COALESCE(c.assigned_to, '') AS assigned_to,
  COALESCE(
    NULLIF(TRIM(COALESCE(ap.given_name, '') || ' ' || COALESCE(ap.surname, '')), ''),
    NULLIF(au.name, ''),
    au.email,
    ''
  ) AS assigned_to_name,
  COALESCE(c.assigned_time, '') AS assigned_time,
  COALESCE(c.retired, '') AS retired,
  c.created,
  c.updated
FROM computers c
LEFT JOIN users lu ON lu.id = c.last_user
LEFT JOIN profiles lp ON lp.uid = c.last_user
LEFT JOIN users au ON au.id = c.assigned_to
LEFT JOIN profiles ap ON ap.uid = c.assigned_to
WHERE
  CASE
    WHEN {:id} != '' THEN c.id = {:id}
    WHEN {:retired} = 1 THEN COALESCE(c.retired, '') != ''
    ELSE COALESCE(c.retired, '') = ''
  END
-- The retired list shows the newest retired first.
ORDER BY
  CASE WHEN {:retired} = 1 THEN c.retired END DESC,
  c.computer_name,
  c.id;
