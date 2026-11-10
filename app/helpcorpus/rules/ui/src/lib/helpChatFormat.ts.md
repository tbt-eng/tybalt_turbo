<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/helpChatFormat.ts","sha256":"dfaaf3780bceba694cc4137bc2144fb83a5193874f67f25c325b0d0baac70d74","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/helpChatFormat.ts

## Help answer links work only for user docs pages

Lines 38-62

Asked as: link not clickable in answer; help link does not work; external link shows as text; link appears as plain text

A Markdown link in a Help answer becomes a link only when its address matches /docs/<slug> or /docs/<slug>#<anchor>. The slug and anchor use only lowercase letters, digits and hyphens. Any other link shows as its text with no link, including external links.

## Raw HTML in a Help answer shows as plain text

Lines 44-58

Asked as: html tags showing in answer; angle brackets visible in help; code tags shown as text; markup appears instead of formatting

Turbo escapes raw HTML in a Help chat answer (the html renderer), so it shows as text and not as markup. It escapes the characters & < > " and '.

## Images in a Help answer show as alt text

Lines 63-65

Asked as: picture missing from help answer; image not showing in chat; alt text instead of image; screenshot does not appear

Turbo does not show an image in a Help chat answer. It shows the image's alt text as escaped plain text.

## Source marks become numbers only after sources arrive

Lines 69-86

Asked as: [S3] showing in answer; citation numbers not appearing; little numbers after answer missing; source numbers appear late

While an answer streams, the sources are null, so Turbo leaves out the chunk marks such as [S3]. After the server sends the sources, each mark that matches a chunk source's ref (ref) becomes a superscript with that source's position in the list, counted from 1. A mark such as [S3, S4] becomes one number for each ref that matches.

## Marks with no matching source are dropped

Lines 73-85

Asked as: source tag disappeared from answer; bracket tag with no number; missing citation in answer; doc source has no number

A chunk mark in an answer whose ref is not among the chunk sources shows nothing. Doc sources get no number in the answer text, because only chunk sources are mapped, but they still count in the list position.

## Copy gives answer text without source marks

Lines 88-94

Asked as: copied answer has brackets; copy button leaves source tags; clipboard text spaces before punctuation; copied text has citations

The Copy text of an answer (answerText) removes the chunk marks such as [S3]. It also removes spaces or tabs before . , ; or :, and trims the ends.

## Invalid sources are left out of the source list

Lines 102-137

Asked as: source missing from list under answer; reference gone from sources; sources list skips an item; why is one source hidden

Turbo leaves out a source from the list under the answer when it is not valid. This covers a doc without a url or with a url outside the docs pattern, and a chunk without a ref or path. A chunk whose path has characters outside letters, digits and _ . / ( ) [ ] + @ - is also left out. Other sources keep their original numbers.

## Doc source link goes to the docs page

Lines 104-115

Asked as: source link opens docs page; doc reference opens in same tab; documentation link label wrong; clicking source goes to docs

A doc source with a url that matches the docs link pattern becomes an internal link that opens in the same tab. Its label is the docs page title for the slug, or the slug when no title is known. Feedback sends the url as the key.

## Code or notes source links to the GitHub repository

Lines 116-134

Asked as: code link opens github; design note link goes to github; source opens repository in new tab; where does code reference link

A chunk source with a ref and a safe path becomes an external link (opens in a new tab) to the tbt-eng/tybalt_turbo repository on GitHub (main branch). Its label starts with "Design note:" when the type is notes, and with "Code:" otherwise. Feedback sends the ref as the key.

## Chunk source line range needs both start and end

Lines 122-122

Asked as: link opens whole file not lines; github link no line numbers; line highlight missing in file link; jump to lines not working

A chunk source link gets a line anchor (#Lstart-Lend) only when both start and end are set and non-zero. Otherwise it links to the whole file.

## Chunk source label uses heading only for notes

Lines 125-130

Asked as: source label shows file path; design note title instead of path; heading not shown on source; label is just file name

For a notes chunk with a heading, the label shows the heading. For code, or for notes without a heading, it shows the file path.
