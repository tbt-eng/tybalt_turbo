<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/docsFormat.ts","sha256":"9e61fcb2f8e132822edd4cdd9cb4046ad79c35375d14b397d60926b9de4a5218","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/docsFormat.ts

## Help links open the page with the same link style

Lines 33-39

Asked as: help link style; link looks wrong in help; markdown link in help; help link formatting

In Help docs, Markdown links render as an anchor with the class "link" and the given address. The link text keeps its inline formatting.

## Commit links point to the tbt-eng tybalt_turbo repository

Lines 45-49

Asked as: commit link goes nowhere; link to commit in help; github commit link; wrong repository for commit

A commit link in Help goes to https://github.com/tbt-eng/tybalt_turbo/commit followed by the commit hash.

## Section ids and word counts

Lines 51-61

Asked as: section anchor id; heading slug wrong; section id from title; word count wrong in section

A section id is its title in lowercase. Apostrophes are removed and every other run of characters outside a-z and 0-9 becomes one hyphen. Hyphens at the start or end are trimmed. A section word count counts words that start with a letter. A page word count covers the whole source.

## Help docs flow lines need three parts

Lines 63-80

Asked as: flow line error; flow step format error; help doc parse error; missing part in flow line

A flow block line in a help doc must have the form "Actor | Step | Detail", with three non-empty parts. Otherwise the docs parser throws the error: A flow line needs "Actor | Step | Detail": followed by the line. A line that starts with "Note:" becomes a note and is not a step. Blank lines are skipped.

## Flow blocks become steps and notes

Lines 82-100

Asked as: flow code block not rendering; steps and notes in help; help block shows as plain text; flow block shows as code

In a Help doc, a code block marked "flow" becomes a flow block of steps and notes, not ordinary HTML. The text around it is split into separate HTML blocks before and after it. Blank space tokens are dropped.

## Help doc needs a level-1 title

Lines 102-121

Asked as: help doc has no title; help page missing heading; doc needs title line; help page won't load title

A help doc must have a "# Title" heading. The first level-1 heading becomes the page title, and a later level-1 heading is not used as a title. Without one, the parser throws: Doc <slug> needs a "# Title" line.

## Only the first level-1 heading is the page title

Lines 102-120

Asked as: second title in help doc; two h1 headings help; which heading is the page title; extra big heading in help

The first level-1 heading is the page title. Level-2 headings start new sections. Content before the first level-2 heading goes into the page intro. A second level-1 heading is not a title and goes into the current section or intro.

## What's new headings carry the last commit

Lines 106-120

Asked as: what's new commit missing; week commit at heading; commit hash in whats new title; last commit for what's new week

A level-2 heading starts a new section of a help page. If the heading ends with "(to " plus a commit hash of 8 to 40 lowercase hex characters and ")", the hash is stored as the section's commit (the last commit a What's new week covers) and is removed from the title. Content before the first level-2 heading is the page intro.

## Help page word count covers the whole source

Lines 122-135

Asked as: help page word count; word count off in help; how words are counted in docs; section word total wrong

A Help page's word count counts the whole source text. Each section counts only its own content, without its heading.

## Help page order and child pages

Lines 145-192

Asked as: help menu page order; help page in wrong place; child page not under parent; where help page appears in list

Help lists pages in the order of docPageOrder: start-here, time, purchase-orders, expenses, proposals, project-requests, rate-sheets, profile, glossary, whats-new. A page not in that list goes just before What's new, sorted by title. Child pages (review-time, revise-a-purchase-order, find-a-purchase-order, commit-expenses, admin-profiles) appear right after their parent, in the order of the child list, with a shorter label and an indent.

## Help popover links open in a new tab

Lines 194-207

Asked as: help link opens same tab; help popover link new window; link in popup closes form; unsaved work lost clicking help link

In a help popover over a form, links open in a new tab (target="_blank", rel="noopener") so unsaved work is not lost. This applies to the text blocks, the flow step details and the flow notes.

## Help references need a page and a section

Lines 209-218

Asked as: help reference shows nothing; link to help section empty; help section not found; help anchor blank

A help reference has the form "page-slug#section-id", such as project-requests#approve-a-request. Turbo returns nothing if the page or the section is not found, or if the reference has no section id.
