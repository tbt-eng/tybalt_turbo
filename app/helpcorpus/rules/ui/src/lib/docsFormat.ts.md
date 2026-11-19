<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/docsFormat.ts","sha256":"bdc0aeda580385a2377567c3474b23444e976f66797ea7b990bd788bc394e658","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/docsFormat.ts

## Help links open in the same tab with a styled link

Lines 33-39

Asked as: help link opens same tab; link looks like plain text; clicking link in help; help link not styled

Turbo renders every Markdown link in a help page as an anchor with the class "link". The anchor keeps the link's href and the link's text.

## Commit links point to the tbt-eng tybalt_turbo repository

Lines 45-49

Asked as: commit link wrong repo; link to commit on github; whats new commit link; commit link goes nowhere

A What's new commit links to https://github.com/tbt-eng/tybalt_turbo/commit followed by the commit hash (commitUrl).

## Section ids come from titles

Lines 51-57

Asked as: section anchor link broken; section link id wrong; heading id from title; section id format

Each section's id is its title in lowercase. Apostrophes are removed, and each run of other characters becomes one hyphen. Hyphens at the start and end are cut off.

## Word count of docs counts only words that start with a letter

Lines 59-61

Asked as: word count ignores numbers; doc word count wrong; word count skips numbers; why word count low

A word counts if it starts with a letter. Numbers alone do not count. Each section counts the words in its own text, and a page counts the words in its whole source.

## Flow lines in the docs need three parts

Lines 63-80

Asked as: flow line error; help page fails to load; flow step missing part; pipe separated line error

A flow block line in a help page must have the form "Actor | Step | Detail", with all three parts filled in. A line with another number of parts, or with an empty part, makes the page fail to load with the error 'A flow line needs "Actor | Step | Detail"'. A line that starts with "Note:" is a note, not a step.

## Flow notes start with Note: and skip blank lines

Lines 64-80

Asked as: note in flow not showing; blank lines in flow; how to add note; note prefix flow

In a flow block, Turbo ignores blank lines. A line that starts with "Note:" becomes a note, and Turbo trims the text after the prefix. The detail of a step and the text of a note accept inline Markdown, such as links.

## Flow blocks come from code blocks marked flow

Lines 82-100

Asked as: flow block not rendering; flow code block; steps not showing as flow; flow language code fence

A fenced code block with the language "flow" becomes a step-by-step flow block. Turbo renders the other content as HTML, and it groups consecutive tokens into one HTML block. Blank spacing tokens are dropped.

## Only the first level-1 heading is the page title

Lines 102-120

Asked as: two titles help page; wrong page title; heading becomes section; which heading is title

The first level-1 heading becomes the page title. A level-2 heading starts a new section. Content before the first level-2 heading goes in the page intro, and the rest goes to the section above it.

## A help page needs a level-1 title

Lines 106-121

Asked as: help page won't load; needs title line error; missing title help page; doc fails no heading

A help page must have a "# Title" heading. Only the first level-1 heading is the title. Without one, the page fails to load with the error 'Doc <slug> needs a "# Title" line'.

## What's new heading ends with its last commit

Lines 109-116

Asked as: hash shown in heading; whats new version hash; commit suffix in heading; heading with commit id

A level-2 heading that ends with "(to <hex commit>)" gives the section a commit. The commit must be 8 to 40 lowercase hex characters. The section title drops that suffix. Any other level-2 heading keeps its full text and has no commit.

## Page and section word counts

Lines 126-134

Asked as: section word count; word count includes heading; words counted per section; page word total

Each page counts the words of its whole source, including the headings. Each section counts the words of its own content (the raw text), without its heading.

## Order of pages in Help

Lines 137-192

Asked as: help page order; menu order wrong; pages sorted wrong; help list sequence

Help lists pages in this order: start-here, time, purchase-orders, expenses, proposals, project-requests, rate-sheets, profile, glossary, whats-new. A child page follows its parent, in the order of the child list. A page in neither list goes just before What's new, sorted by title.

## Child pages have short labels and indent in navigation

Lines 158-198

Asked as: shorter label in menu; review and track label; sub page indent; child page short name

Some pages show under a parent with a shorter label, for example review-time as "Review and track" under time. Other pages show with their own title and no indent.

## Pages not in the order list go before What's new

Lines 176-192

Asked as: page missing from menu; new page position help; page above whats new; unlisted page sorted

A top-level page that docPageOrder does not list goes just before What's new. Turbo sorts such pages by title. Listed pages keep the order of the list, and ties also sort by title.

## Child pages sit under their parent in listed order

Lines 182-191

Asked as: child page not under parent; subpage position in menu; children order wrong; sub page listed after parent

Turbo lists a child page right after its parent page, not among the top-level pages. Children of one parent follow the order of the childDocPages list, not alphabetical order.

## Links in help popovers open in a new tab

Lines 200-213

Asked as: popover link new tab; help popup link leaves form; link in popup opens; unsaved work lost link

In a help popover, every link opens in a new tab, so that people do not lose unsaved work on the form. This covers links in text, in flow step details and in notes.

## Help references need a page and a section

Lines 215-224

Asked as: help link finds nothing; broken help reference; section reference not found; page anchor link empty

A help reference such as "project-requests#approve-a-request" finds a section only when both the page slug and the section id exist. If either is missing, Turbo finds nothing.
