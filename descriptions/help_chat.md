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

The chat is off until an admin adds the `help_chat` setting. See [app_config.md](app_config.md).

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
5. Run the tools that the model calls: `search_turbo`, `read_source` and `report_gap`.
6. Stream the answer, then the sources that it cites.
7. Write the result to the row.

Events: `meta` (chat id and turn id), `status`, `text`, `sources`, `done`.

The answer cites a docs page as a link, such as `[Time](/docs/time)`, and a chunk from a tool as `[S3]`. The server keeps a source only if it knows it.

Other routes: `POST /api/help/turns/{id}/feedback` saves a rating, a note, a source that the user opened, or a copy. `POST /api/help/warm` asks the model to cache the shared prompt.

## The log

`chat_turns` holds one row for each turn: the question and answer, the sources, the model, effort, build and corpus hash, the user's claims, the tool calls with the ids of what they returned, the route and record, the status, token counts, times, and the user's feedback.

The log holds chunk ids and not chunk text. The corpus hash says which text an id means, and the generator gives the same corpus for the same commit.

A user can list and view their own rows through the collection API. Nobody can create, change or delete a row there. Feedback goes through its own route, which changes the feedback fields only.

## The model

`app/helpchat/anthropic.go` is the only file that uses the Claude SDK. The rest of the package talks to the `Model` interface, and the tests use a stub. `anthropic_test.go` checks the request body against a local server that speaks the API's streaming format.

The setting picks `claude-haiku-5-5` or `claude-sonnet-5-5`. The request sets no sampling values, no forced tool choice and no thinking mode, because these models refuse them.
