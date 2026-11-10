<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/config.ts","sha256":"beba1b36219f33e455d095c7fd38114db26fb6b9707f6a2712ff8a630b2a9b81","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/config.ts

## Session token refreshes 5 minutes before it expires

Lines 3-4

Asked as: session expiring while I work; token expires logged out; lost work expired session; refresh sign in before timeout

Turbo refreshes the sign-in token 5 minutes before it expires (TOKEN_REFRESH_BUFFER_MINUTES). This prevents staff from saving work with an expired token. The code also holds this buffer as 300,000 milliseconds (TOKEN_REFRESH_BUFFER_MS).

## Mouse, keyboard and touch activity can restart or refresh a session

Lines 6-8

Asked as: logged out after clicking; session restarts when I click; random sign in prompt; app wakes up session

Turbo treats a mousedown, keypress, touchstart or click as a sign that the app is awake. On these events the auth store can restart or refresh the session without a sign-in prompt (SESSION_CHECK_EVENTS).

## Token refresh buffer is converted from minutes to milliseconds

Lines 11-13

Asked as: token refresh buffer milliseconds; convert minutes to ms; refresh buffer 300000; timeout value conversion

Turbo computes the token refresh buffer in milliseconds (TOKEN_REFRESH_BUFFER_MS) as the 5-minute setting times 60 times 1000, which gives 300,000 ms.
