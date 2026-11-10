<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/helpChatStream.ts","sha256":"203037d522fa32707447856d6a9069984b625db34a07df82b6ffaa6edb5cc033","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/helpChatStream.ts

## Help chat answer events can arrive split across pieces

Lines 11-25

Asked as: help chat answer cut off; streamed reply arrives in pieces; help reply shows partly; chat answer missing text stream

Turbo collects the answer text and handles an event only when it is complete, which is after a blank line. An event that arrives in several pieces is handled once. Turbo treats Windows line endings (\r\n) like \n.

## Help chat event name defaults to "message"

Lines 20-28

Asked as: help chat event without name; stream message default event; help reply lines joined wrong; chat stream data lines broken

When an event block has no event line, Turbo names the event "message". Several data lines in one block are joined with a newline before Turbo reads them as JSON.

## Help chat ignores events without data or with invalid JSON

Lines 26-31

Asked as: help chat blank reply no error; help stream empty event skipped; chat answer invalid json; help reply disappears silently

Turbo skips an event block that has no data lines. It also skips a block whose data is not valid JSON. The skipped block gives no error and no event.

## Help chat refusal shows the server's message or a default

Lines 46-73

Asked as: help chat failed error message; help can't answer now; help unavailable message shown; chat request no response text

When the Help chat request fails or has no response body, Turbo shows the server's message if it is a non-empty text. If the server gives no usable message, Turbo shows "Help can't answer now. Try again in a minute." (helpUnavailable).

## Help chat handles a last event with no closing blank line

Lines 82-82

Asked as: last help reply missing; help chat final message not shown; stream ended without blank line; help answer stuck at end

When the stream ends, Turbo adds a blank line to the leftover text. A final event that the server did not close with a blank line is still handled.
