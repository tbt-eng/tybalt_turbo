# Help chat

Help chat answers questions about Turbo. A signed-in user asks a question in Help, and a language model answers from the user docs, the design notes and the code. The model reads only. It has no tool that writes.

This file says how the server side works. The steps to build it, with diagrams, are in the design doc that this work follows.

## Status

| Part | State |
| --- | --- |
| Corpus generator (`app/helpcorpus`, `app/cmd/helpcorpus`) | Built |
| Search index with the claims filter (`app/helpcorpus/index.go`) | Built |
| Chat endpoint and log (`app/helpchat`, `chat_turns`) | Built, tested with a stub model |
| Ask in the UI | Not built |
| Records in answers, and drafts | Later phases |

The chat is off until a superuser adds the `help_chat` setting. See "Set up the model" below.

## Set up the model

The chat needs two things: a key on the server, and a setting in the database. The key is a secret. The setting is not.

### 1. Choose a provider

| Provider | `provider` value | Secret for the key | Models |
| --- | --- | --- | --- |
| Claude, from Anthropic | `anthropic` | `ANTHROPIC_API_KEY` | `claude-haiku-5-5` (the default) or `claude-sonnet-5-5` |
| OpenAI | `openai` | `OPENAI_API_KEY` | Any model id from OpenAI's model list |
| OpenRouter | `openrouter` | `OPENROUTER_API_KEY` | Any model id from OpenRouter's model list |

Use `anthropic` for a Claude model, and not `openrouter`. Only the `anthropic` client marks the shared prompt for Claude's cache. Without the mark, each question pays the full price for the handbook.

### 2. Put the key on the server

In production, the key is a Fly secret:

```bash
fly secrets set ANTHROPIC_API_KEY=your-key
```

Fly restarts the app when a secret changes. You can set the keys of more than one provider, and then change between them with the setting alone.

In development, set the key in the terminal that runs the server:

```bash
export ANTHROPIC_API_KEY=your-key
```

```bash
go run main.go serve
```

Do not put a key in the repo, in `ui/.env`, or in `app_config`. Every signed-in user can read `app_config`.

### 3. Add the setting

1. Open the PocketBase admin at `/_/` and sign in as a superuser.
2. Open the `app_config` collection.
3. Add a record. Set `key` to `help_chat`, and set `value` to one of these.

```json
{ "enabled": true }
```

```json
{ "enabled": true, "model": "claude-sonnet-5-5" }
```

```json
{ "enabled": true, "provider": "openai", "model": "gpt-5-mini" }
```

```json
{ "enabled": true, "provider": "openrouter", "model": "openai/gpt-5-mini" }
```

The first one uses Claude Haiku, which is the default. The model ids for OpenAI and OpenRouter are examples: copy the id from the provider's own list. [app_config.md](app_config.md) lists every property.

### 4. Change the model, or turn the chat off

Edit the `value` of the `help_chat` record. The next question uses the new values. The app does not restart.

- To change the model, change `model`.
- To change the provider, change `provider` and `model` together. The server must have that provider's key already (step 2).
- To turn the chat off, set `enabled` to `false`.
- `effort` is `"low"` unless you set it. For an OpenAI model that has no reasoning setting, set `"effort": ""`.

Each row of `chat_turns` records the provider and the model that answered, so you can compare two models on real questions.

### 5. Check that it works

Until Help has the Ask button, send a question from a terminal. Copy your token from the browser: it is the `token` in the `pocketbase_auth` item of the site's local storage.

```bash
curl -N -X POST http://localhost:8090/api/help/ask -H "Authorization: YOUR_TOKEN" -H "Content-Type: application/json" -d '{"question":"How do I submit my week?"}'
```

The answer arrives as `text` events, and the last event is `done`.

| What you get | What it means |
| --- | --- |
| HTTP 401 | You are not signed in. |
| HTTP 403 | The chat is off: there is no `help_chat` record, or `enabled` is false. |
| HTTP 503 | The setting has a bad value, or the server has no key for the provider. The server log has a line that starts "help chat" and gives the reason. |
| `done` with status `error` | The provider refused the request or failed. The `error` field of the turn's row in `chat_turns` holds the provider's message. |

### Add another provider

A provider that speaks the OpenAI format needs one entry in `app/helpchat/providers.go`: its name, the secret for its key, and its address. A provider with another format needs a new file that implements `Model`.

### What is tested

Tests send real requests from both clients to local servers that speak each API's format. No test, and nobody yet, has sent a question to Anthropic, OpenAI or OpenRouter from this code. Do step 5 with each provider before you rely on it.

## Rules that must hold

1. A visitor who is not signed in gets nothing from the chat. Every route under `/api/help` uses `apis.RequireAuth("users")`. The corpus is embedded from `app/helpcorpus/data`, which is not in `pb_public`, and no route serves it.
2. The model gets only text that the user may read. Every search and read of the corpus takes the user's claims.
3. The server decides who the user is. The user id and the claims come from the auth token and the database. The browser sends only the question, the chat id, the route and the record on the page.
4. A user reads only their own chats. The `chat_turns` rules and every query filter on the user id.
5. The model never writes to the database.

## What the model may read

The repo is public, so the user docs, the design notes and the code are open to every signed-in user. `app/helpcorpus/audience.go` holds the rules. A rule can restrict a folder or a word to some claims later.

Uncommit is not a secret, but the chat does not raise it. The system rules tell the model to treat a commit as final, so that people do not come to rely on a reversal. A holder of `commit` or `admin` gets an answer when they ask about it. This is guidance for the answer and not access control.

When a user cannot do a task, the answer says who can, in the glossary's words and not in claim names.

## A turn

`POST /api/help/ask` runs one turn and returns server-sent events.

1. Check the setting, the question and the limits. Load the user's claims.
2. Write a `chat_turns` row with status `running`, in the same transaction as the limit checks.
3. Build the prompt. The first part is the system rules and the handbook (all the user docs). It is the same bytes for every user, so the model caches it. The second part names the user's claims, the page they are on and all the claims in Turbo.
4. Send the earlier turns of the chat and the question. The server reads the earlier turns from `chat_turns`.
5. Run the tools that the model calls: `search_turbo`, `read_source` and `report_gap`. One search returns up to 6 chunks of code and then up to 2 chunks of design notes, each marked with its type. They are two groups because the notes have headings in plain words: in one list they outscore the code, and the code is what says how Turbo behaves today. The prompt tells the model to search the code for a question about a rule or a case, to trust the code over a note, and to say when it could not confirm a rule.
6. Stream the answer, then the sources that it cites.
7. Write the result to the row.

Events: `meta` (chat id and turn id), `status`, `text`, `sources`, `done`.

The answer cites a docs page as a link, such as `[Time](/docs/time)`, and a chunk from a tool as `[S3]`. The server keeps a source only if it knows it.

Other routes: `POST /api/help/turns/{id}/feedback` saves a rating, a note, a source that the user opened, or a copy. `POST /api/help/warm` asks the model to cache the shared prompt.

## The log

`chat_turns` holds one row for each turn: the question and answer, the sources, the model, effort, build and corpus hash, the user's claims, the tool calls with the ids of what they returned, the route and record, the status, token counts, times, and the user's feedback.

The log holds chunk ids and not chunk text. The corpus hash says which text an id means, and the generator gives the same corpus for the same commit.

A user can list and view their own rows through the collection API. Nobody can create, change or delete a row there. Feedback goes through its own route, which changes the feedback fields only.

## The model clients

The chat talks to a model through the `Model` interface in `model.go`. The tests use a stub.

- `anthropic.go` calls the Claude API with Anthropic's SDK. It is the only file that uses the SDK. It marks the shared prompt for Claude's cache, and it sends the model's thinking blocks back unchanged between tool calls.
- `openai.go` calls any API in the OpenAI chat completions format, with plain HTTP. It puts the shared prompt first, because these services cache the start of a prompt by themselves.
- `providers.go` lists the providers: the secret for each key, the address, and the models and effort values that the setting may name.

Both clients send the same system rules, handbook and tools. Neither sets sampling values or a forced tool choice.
