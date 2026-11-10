<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/helpcorpus/collections/chat_turns.md","sha256":"56da4c7a616bfa528ed7a952e5d8636c9d797fc2085b6b5961af573f4ee5df15","model":"claude-sonnet-5-5"} -->

# Rules in app/helpcorpus/collections/chat_turns.md

## Who can list chat turns

Lines 7-11

Asked as: list my chats; see my chat history; other people's chat turns; who can list chats

A signed-in person can list only the chat turns (chat_turns records) that belong to them, where the record's user (uid) is that person. Nobody can list another person's chat turns through the API.

## Who can view one chat turn

Lines 13-17

Asked as: view someone else's chat; open another chat turn; can't see chat record; read a colleague's chat

A signed-in person can view a chat turn (chat_turns record) only when the record's user (uid) is that person. Other people's chat turns are not visible through the API.

## Nobody can create chat turns through the API

Lines 19-21

Asked as: create chat turn via api; add a chat record; insert chat message; api can't post chats

Nobody can create a chat turn (chat_turns record) through the API, except a PocketBase superuser. Turbo's own server code can still create them.

## Nobody can update chat turns through the API

Lines 23-25

Asked as: change my chat rating; edit rating note; update chat answer; rating won't save

Nobody can update a chat turn (chat_turns record) through the API, except a PocketBase superuser. This includes the rating and the rating note. Turbo's own server code can still update them.

## Nobody can delete chat turns through the API

Lines 27-29

Asked as: delete my chat history; remove chat turns; user deleted chats gone; erase conversation records

Nobody can delete a chat turn (chat_turns record) through the API, except a PocketBase superuser. Turbo's own code can still delete one. Deleting the user (users record) also deletes that user's chat turns (cascadeDelete on uid).

## Each chat has one record per turn number

Lines 37-74

Asked as: duplicate chat turn number; same turn twice error; chat id too long; turn number must be whole

A unique index on chat and turn (idx_chat_turns_chat_turn) means two chat turns cannot have the same chat and the same turn number. The chat (chat) is required text of at most 15 characters, and the turn number (turn) must be a whole number.

## Chat turn text length limits

Lines 39-64

Asked as: question too long; answer text limit; rating note too long; chat message character limit

A chat turn's question (question) is required and has at most 4000 characters. The answer (answer) has at most 60000 characters. The rating note (rating_note) has at most 1000 characters, and the error (error) at most 2000.

## Chat turn status and surface allowed values

Lines 51-63

Asked as: invalid chat status; surface value rejected; rating only up or down; chat status options

A chat turn must have a status (status) of running, answered, gap, refused, cancelled or error. The surface (surface) is optional and can only be help or popover. The rating (rating) is optional and can only be up or down.
