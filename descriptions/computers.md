# Computers

Each Windows computer runs a PowerShell script called radiator at sign-in (`stamler/psFUN`, `radiator/radiator.ps1`). It POSTs one JSON report about the computer and the user. Turbo stores the reports in three collections. This replaces the `rawLogins` function of the legacy system, which is the reference for the behaviour. See issue 211.

The endpoint is `app/routes/radiator.go`. The collections are added by `app/migrations/1795200837_computers.go`.

---

## Endpoint: `POST /api/radiator`

No user auth. The client sends a machine secret with the role `radiator`:

```
Authorization: Bearer <secret>
Content-Type: application/json
```

An admin makes the secret with `POST /api/machine_secrets/create` (`{"days": 365, "role": "radiator"}`). A `legacy_writeback` secret is refused. The legacy header `Authorization: TYBALT <secret>` is refused too: radiator 10 sends `Bearer`.

The checks run in this order. Each one stops the request.

| Check | Reply |
| --- | --- |
| The method is not POST | `405`, with `Allow: POST` |
| No `Bearer` token, or no unexpired `radiator` secret matches it | `401` |
| `Content-Type` is not `application/json` (a charset is fine) | `415` |
| The body is over 64 KB | `413` |
| The body is empty, `{}`, `null` or `[]` | `400`, nothing stored |
| The body is not JSON, or not a JSON object | `400`, nothing stored |
| The report is invalid | `202`, stored as a raw login |
| The report is valid | `202`, stored as a computer and a login |
| A database error | `500`, with no details |

The method, the token and the Content-Type are checked before any of the body is read. The body limit is `apis.BodyLimit`: it refuses a larger `Content-Length` at once and stops reading a body without one just past the limit. It replaces PocketBase's default limit of 32 MB for this route.

### Validation

The rules are those of the legacy `RawLogins.schema.json`, as ajv 6 applied them:

- Required: `serial`, `upn`, `networkConfig`, `userSourceAnchor`, `radiatorVersion`, `mfg`, `userGivenName`, `userSurname`.
- Lengths and patterns: `serial` 4 to 48 characters, `mfg` 2 to 32, `userGivenName` 1 to 24, `userSurname` 2 to 35, `osArch` 2 to 16, `bootDriveFS` 2 to 32, `model` at most 48. `upn` must hold an `@`. `userSourceAnchor` is 32 hex characters. `bootDrive` is a letter and a colon. `osVersion` is dotted numbers.
- `networkConfig` is an object with at least one key, and each key is a MAC address (`AA:BB:CC:DD:EE:FF`).
- `systemType`, `osSku`, `ram`, `bootDriveCap`, `bootDriveFree` and `radiatorVersion` are integers.
- Unknown properties are dropped.
- Types are coerced, because old PowerShell sends some numbers as strings: `" 5."` becomes `5`, a number sent for a string becomes its text, `null` becomes `""` or `0`, a boolean becomes `1` or `0`.
- An invalid `email` or `computerName` is dropped, not a reject. An empty email is the common case.

A valid report can still be a reject when its serial or manufacturer is too short once cleaned for the key (see below), or when a value is too long for its field.

### A valid report

One transaction:

1. Upsert the computer by its slug. The computer fields of the report replace the stored ones; a field the report leaves out keeps its stored value. The user fields are always replaced, because they describe the last report. `last_login` is now. `assigned_to`, `assigned_time` and `retired` are never written.
2. Match the user (see below) and set `last_user`.
3. Add one row to `computer_logins`.
4. Delete the raw logins with the same `computer_name`. Legacy left this as a TODO.

The anchor, the email and the UPN are stored in lowercase.

### An invalid report

It is stored in `computer_raw_logins` as it was received, with the reasons in `error`. Nothing else changes. Nothing reprocesses raw logins.

---

## The computer key

`slug` is a port of `makeSlug` in the legacy `utilities.ts`:

- The serial loses white space, `/` and `,`.
- The manufacturer is lowercased, loses `/`, `.`, `,`, `inc` and `ltd` (anywhere), is trimmed and has spaces replaced by `_`.
- The key is `serial,manufacturer`. The serial must keep at least 4 characters and the manufacturer 2.

So `"kQDS 9/452,At"` from `"Dell, Inc. Ltd"` is `kQDS9452At,dell`.

---

## Matching the user

1. The admin profile whose `user_source_anchor` equals the report's anchor. The anchor is the AD `mS-DS-ConsistencyGuid`, as 32 lowercase hex characters.
2. Otherwise, a user whose email equals the report's email, then its UPN, ignoring case.
3. Otherwise no user. `last_user` and the login's `user` stay empty.

The endpoint never creates a user. `admin_profiles.user_source_anchor` has a unique index that skips empty values, so one anchor can't match two people. Legacy had a reject for that case; Turbo doesn't need it.

Nothing sets the anchor yet. It is to be filled one time from the legacy Profiles; see issue 211.

---

## Collections

Only server code writes them: none has a create, update or delete rule. A holder of the `computers` claim or of `admin` can list and view all three.

### `computers`

One row per computer, updated in place.

| Field | Notes |
| --- | --- |
| `slug` | Unique key, see above |
| `serial`, `mfg`, `model`, `computer_name`, `system_type`, `ram`, `os_version`, `os_sku`, `os_arch`, `boot_drive`, `boot_drive_fs`, `boot_drive_cap`, `boot_drive_free`, `radiator_version` | From the report |
| `network_config` | JSON, one entry per MAC address |
| `user_given_name`, `user_surname`, `email`, `upn`, `user_source_anchor` | The user of the last report |
| `last_user` | The matched user, or empty |
| `last_login` | Time of the last valid report |
| `assigned_to`, `assigned_time` | Set by people, not reports |
| `retired` | Set by people, not reports |

### `computer_logins`

One row per valid report: `computer` (deleted with the computer), `given_name`, `surname`, `user_source_anchor`, `user` (the matched user, or empty) and `created`.

### `computer_raw_logins`

One row per rejected report: `report` (the body as received), `computer_name` (from the report, when valid), `error` (the reasons, separated by `; `) and `created`.

---

## Not done yet

The pages, assign and retire, the 30 day retention job, and the user docs are later steps of issue 211.
