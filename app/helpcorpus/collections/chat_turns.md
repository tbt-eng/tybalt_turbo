<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: chat_turns

`chat_turns` is a collection of records. Turbo's screens reach it through PocketBase's API (`/api/collections/chat_turns/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): chat_turns

```
@request.auth.id != '' && uid = @request.auth.id
```

## Who can view one record (view rule): chat_turns

```
@request.auth.id != '' && uid = @request.auth.id
```

## Who can create a record (create rule): chat_turns

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can update a record (update rule): chat_turns

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can delete a record (delete rule): chat_turns

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Fields: chat_turns

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; min 15; max 15; pattern "^[a-z0-9]+$" |
| `uid` | relation | required; maxSelect 1; to `users`; cascadeDelete |
| `chat` | text | required; max 15 |
| `turn` | number | onlyInt |
| `question` | text | required; max 4000 |
| `answer` | text | max 60000 |
| `sources` | json |  |
| `provider` | text | max 50 |
| `model` | text | max 100 |
| `effort` | text | max 20 |
| `build` | text | max 100 |
| `corpus` | text | max 64 |
| `claims` | json |  |
| `tool_calls` | json |  |
| `route` | text | max 200 |
| `record` | text | max 100 |
| `surface` | select | values ["help","popover"]; maxSelect 1 |
| `status` | select | required; values ["running","answered","gap","refused","cancelled","error"]; maxSelect 1 |
| `gap` | json |  |
| `stop_reason` | text | max 50 |
| `error` | text | max 2000 |
| `request_id` | text | max 100 |
| `tokens_in` | number | onlyInt |
| `tokens_cache_read` | number | onlyInt |
| `tokens_cache_write` | number | onlyInt |
| `tokens_out` | number | onlyInt |
| `ms_first_token` | number | onlyInt |
| `ms_total` | number | onlyInt |
| `rating` | select | values ["up","down"]; maxSelect 1 |
| `rating_note` | text | max 1000 |
| `opened_sources` | json |  |
| `copied` | bool |  |
| `created` | autodate | onCreate |
| `updated` | autodate | onCreate; onUpdate |

## Indexes: chat_turns

A unique index means that two records can't have the same values in its fields.

- `CREATE UNIQUE INDEX `idx_chat_turns_chat_turn` ON `chat_turns` (`chat`, `turn`)`
- `CREATE INDEX `idx_chat_turns_uid_created` ON `chat_turns` (`uid`, `created`)`
