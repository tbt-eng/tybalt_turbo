<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/client_notes_api.go","sha256":"2cc6ba65b8577f9cd05915992b831937381a2861c20d06b0fff6be762e2dd286","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/client_notes_api.go

## Client notes require a client id

Lines 34-37

When a person asks for a client's notes without a client id, Turbo refuses the request. The message is "client id is required".

## Client notes list newest first, with the linked job

Lines 39-55

Turbo lists all notes of a client, newest first (by created date). Each note shows its author and the job it is linked to, if there is one.

## Note author name falls back to empty text

Lines 46-53

Each note shows its author's id and email. If the author has no profile name, the given name and surname show as empty text and the note is still listed.

## Client notes load failure message

Lines 70-72

When Turbo cannot load a client's notes, it shows the message "failed to load client notes".

## Client note without a job shows no job

Lines 76-83

A client note that is not linked to a job shows an empty job (null). Turbo gives the job's id, number and description only when the note has a job.

## Each client note shows its text, time and author

Lines 84-96

Each note in a client's note list carries its id, its created time, the note text, its linked job (if any) and its author. The author entry holds the author's id, email, given name and surname.

## Job notes require a job id

Lines 104-107

When a person asks for a job's notes without a job id, Turbo refuses the request. The message is "job id is required".

## Job notes list newest first

Lines 109-151

Turbo lists the notes linked to a job, newest first (by created date). Each note shows its author. The job field of each note is always empty (null).

## Job notes author name falls back to empty text

Lines 113-116

Each job note shows its author's id and email. If the author has no profile name, the given name and surname show as empty text.

## Job notes load failure message

Lines 133-135

When Turbo cannot load a job's notes, it shows the message "failed to load job notes".
