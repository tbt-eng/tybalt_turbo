<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/absorb/suggest/dismissals.go","sha256":"812243168a4940bd0222929276bf6f270972c490d0c884429e35c50c89472bba","model":"claude-sonnet-5-5"} -->

# Rules in app/absorb/suggest/dismissals.go

## Bucket hidden when a 'Not duplicates' set holds all its records

Lines 26-46

Asked as: suggestion hidden after not duplicates; dismissed duplicate suggestion gone; hide duplicate suggestion permanently; not duplicates keeps disappearing

Turbo hides an absorb suggestion bucket when one dismissed set (a 'Not duplicates' set) holds every record in the bucket. The dismissed set may also hold other records. The bucket is hidden if the bucket's records are a subset of that one set (Covers).

## Bucket shows again when it has a new record

Lines 26-46

Asked as: duplicate suggestion came back; new record reappears as duplicate; dismissed suggestion showing again; not duplicates reappeared

A new record in a bucket is in no dismissed set with the other records. Turbo then shows the bucket again. Turbo does not combine several dismissed sets to cover a bucket: one set must hold all its records.

## Empty bucket is never hidden

Lines 29-32

Asked as: empty duplicate group still showing; blank absorb suggestion; no records suggestion stays; empty suggestion not hidden

A bucket with no records is never covered by a dismissed set, so Turbo does not hide it (Covers returns false).
