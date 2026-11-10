<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/helpChatFormat.ts","sha256":"dfaaf3780bceba694cc4137bc2144fb83a5193874f67f25c325b0d0baac70d74","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/helpChatFormat.ts

## Help answer links work only for user docs pages

Lines 38-62

A Markdown link in a Help answer becomes a link only when its address matches /docs/<slug> or /docs/<slug>#<anchor>. The slug and anchor use only lowercase letters, digits and hyphens. Any other link shows as its text with no link, including external links.

## Raw HTML in a Help answer shows as plain text

Lines 44-58

Turbo escapes raw HTML in a Help chat answer (the html renderer), so it shows as text and not as markup. It escapes the characters & < > " and '.

## Images in a Help answer show as alt text

Lines 63-65

Turbo does not show an image in a Help chat answer. It shows the image's alt text as escaped plain text.

## Source marks become numbers only after sources arrive

Lines 69-86

While an answer streams, the sources are null, so Turbo leaves out the chunk marks such as [S3]. After the server sends the sources, each mark that matches a chunk source's ref (ref) becomes a superscript with that source's position in the list, counted from 1. A mark such as [S3, S4] becomes one number for each ref that matches.

## Marks with no matching source are dropped

Lines 73-85

A chunk mark in an answer whose ref is not among the chunk sources shows nothing. Doc sources get no number in the answer text, because only chunk sources are mapped, but they still count in the list position.

## Copy gives answer text without source marks

Lines 88-94

The Copy text of an answer (answerText) removes the chunk marks such as [S3]. It also removes spaces or tabs before . , ; or :, and trims the ends.

## Invalid sources are left out of the source list

Lines 102-137

Turbo leaves out a source from the list under the answer when it is not valid. This covers a doc without a url or with a url outside the docs pattern, and a chunk without a ref or path. A chunk whose path has characters outside letters, digits and _ . / ( ) [ ] + @ - is also left out. Other sources keep their original numbers.

## Doc source link goes to the docs page

Lines 104-115

A doc source with a url that matches the docs link pattern becomes an internal link that opens in the same tab. Its label is the docs page title for the slug, or the slug when no title is known. Feedback sends the url as the key.

## Code or notes source links to the GitHub repository

Lines 116-134

A chunk source with a ref and a safe path becomes an external link (opens in a new tab) to the tbt-eng/tybalt_turbo repository on GitHub (main branch). Its label starts with "Design note:" when the type is notes, and with "Code:" otherwise. Feedback sends the ref as the key.

## Chunk source line range needs both start and end

Lines 122-122

A chunk source link gets a line anchor (#Lstart-Lend) only when both start and end are set and non-zero. Otherwise it links to the whole file.

## Chunk source label uses heading only for notes

Lines 125-130

For a notes chunk with a heading, the label shows the heading. For code, or for notes without a heading, it shows the file path.
