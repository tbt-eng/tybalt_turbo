<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/docsFormat.ts","sha256":"ee54482a851fcb4d60c5759ad241f691af75a1cfb8cb1294829ea8fd9954a424","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/docsFormat.ts

## Links in help pages get the link style

Lines 33-39

Asked as: link style in help; help links look plain; anchor class link; link color in docs

The Markdown renderer turns every link in a doc into an anchor with the class "link". The link keeps its address (href) and its text.

## Commit links go to GitHub

Lines 45-49

Asked as: commit link goes where; link to commit on github; open commit in github; github commit url

A commit link opens https://github.com/tbt-eng/tybalt_turbo/commit/ followed by the commit hash.

## Section ids come from titles

Lines 51-57

Asked as: section anchor id; heading id from title; why section id changed; apostrophe in section link

The id of a section is its title in lower case. Apostrophes are removed, other runs of characters that are not a-z or 0-9 become one hyphen, and a hyphen at the start or end is dropped. For example, "What's new" becomes "whats-new".

## Word counts for pages and sections

Lines 59-61

Asked as: word count wrong; how words are counted; numbers not counted as words; page word count

Turbo counts a word as a run that starts with a letter and continues with letters, digits, underscores, apostrophes or hyphens. Numbers alone do not count. A section counts the words of its own text. A page counts the words of its whole source.

## Flow line format

Lines 63-80

Asked as: flow step format; flow line looks wrong; how to write flow steps; flow note line

A flow block has one step on each line, written "Actor | Step | Detail". A line that starts with "Note:" is a note. A line that does not have exactly three non-empty parts fails with the error: A flow line needs "Actor | Step | Detail". Empty lines are ignored.

## Flow lines need three non-empty parts

Lines 64-80

Asked as: flow line error; a flow line needs; empty part in flow; too many bars in flow

A line in a flow block must have exactly three parts split by "|": Actor, Step and Detail. A line with another number of parts, or with an empty part, causes the error 'A flow line needs "Actor | Step | Detail": <line>'. Blank lines are skipped.

## Flow blocks come from code blocks marked flow

Lines 82-100

Asked as: flow code block not showing; code block shows as code; flow block rendering; text before flow block

A code block whose language is "flow" becomes a flow block of steps and notes, not a code listing. The text before it becomes a separate HTML block. Blank lines (space tokens) are dropped.

## A doc needs a title; text before the first section is intro

Lines 102-135

Asked as: doc needs title; missing page title error; no title line in doc; text before first section

The first level 1 heading is the title of the page. A page without one fails with the error "Doc <slug> needs a "# Title" line". Content before the first level 2 heading is intro. Every level 2 heading starts a new section.

## Only the first level-1 heading is the page title

Lines 102-121

Asked as: second h1 heading; extra level one heading; multiple page titles; later title heading ignored

The first heading of depth 1 is the page title. Later level-1 headings are not titles and go to the intro or the current section as ordinary content.

## Level-2 headings start sections

Lines 106-120

Asked as: new section heading; h2 starts section; where section ends; section content boundaries

Each level-2 heading starts a new section. Everything after it, until the next level-2 heading, belongs to that section.

## What's new heading ends with its last commit

Lines 109-116

Asked as: commit hash in heading; what's new commit marker; last commit on whats new; heading shows commit id

A level 2 heading that ends with "(to " followed by 8 to 40 hex characters and ")" sets the last commit of that section. The title does not include the commit part. A heading without it has no commit.

## Help page order and unlisted pages

Lines 137-194

Asked as: help page order; page missing from help menu; where pages appear in help; help menu sorting

Help lists pages in the order of the list docPageOrder: start-here, time, purchase-orders, expenses, proposals, project-requests, rate-sheets, computers, profile, glossary, whats-new. A page that is not in the list and not a child page goes just before What's new, sorted by title. Pages in the same place sort by title.

## Child pages sit under their parent page

Lines 159-200

Asked as: subpage under parent; indented help page; child page label; review and track menu

Help shows some pages under a parent page with a shorter label and an indent, for example review-time under time as "Review and track". Child pages follow their parent in the order of the childDocPages list. A page that is not a child keeps its own title as its label and has no indent.

## Help popover links open in a new tab

Lines 202-215

Asked as: popover link opens new tab; help link in form opens; links in help popup; link leaves form unsaved

In a help popover over a form, every link opens in a new tab (target="_blank" rel="noopener"). This keeps the person from losing unsaved work. It applies to text blocks and to the details and notes of flow blocks.

## Help reference needs a page and a section

Lines 217-226

Asked as: help reference not found; broken help link; section link finds nothing; page#section link not working

A help reference such as "project-requests#approve-a-request" finds a section only when both the page slug and the section id exist. If either is missing, or the reference has no section id, Turbo finds nothing.
