<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/docsFormat.ts","sha256":"e7cd139f9371a872e1acecdcbc1876fe63cef52ce235096c3686cbbc0a22060f","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/docsFormat.ts

## Doc links render with the link style

Lines 33-39

Asked as: link looks like plain text; docs link style; link colour in help; markdown link not styled

Turbo renders every Markdown link in a docs page as an anchor with the class "link". The link text is parsed as inline Markdown.

## Commit links point to the GitHub commit page

Lines 45-49

Asked as: commit link goes nowhere; link to github commit; open commit on github; commit hash link

Turbo builds a commit link from the base address https://github.com/tbt-eng/tybalt_turbo/commit and the commit hash.

## Section ids come from titles

Lines 51-57

Asked as: section anchor id; link to a heading; help link to section; section link broken

Each docs section gets an id (slugify) from its title: lowercase, apostrophes removed, other runs of non-alphanumeric characters turned into one hyphen, and leading or trailing hyphens trimmed. Help references like 'page#section-id' use this id.

## Word count rule for docs

Lines 59-61

Asked as: how words are counted; word count wrong; docs word count; numbers not counted as words

Turbo counts a word as a run that starts with a letter and may contain letters, digits, underscores, apostrophes and hyphens. Numbers alone do not count. A section's count uses its raw text, and a page's count uses the whole source.

## Flow line format in docs

Lines 63-80

Asked as: flow line error; flow step missing field; flow line format; step needs three parts

Each line of a docs flow block must have three parts, 'Actor | Step | Detail', with none of them empty. Otherwise the parser throws 'A flow line needs "Actor | Step | Detail": <line>'. A line that starts with 'Note:' becomes a note instead. Blank lines are skipped.

## Flow code blocks become step blocks

Lines 82-100

Asked as: flow code block not showing; flow block shows as code; step by step block; code block turned into steps

A code block marked 'flow' in a docs page becomes a step-by-step flow block. All other content becomes HTML. Blank space tokens are dropped.

## Docs page needs a level-1 title

Lines 102-121

Asked as: docs page will not load; missing title error; page has no heading; doc needs title line

Turbo's docs parser (parseDoc) refuses a docs page that has no level-1 heading. It throws the error 'Doc <slug> needs a "# Title" line'. Only the first level-1 heading becomes the page title.

## Text before the first level-2 heading is the page intro

Lines 104-120

Asked as: text above first section; page introduction text; content before first heading; intro paragraph missing

Turbo puts content that comes before the first level-2 heading into the page intro. Content after a level-2 heading goes into that section.

## Only the first level-1 heading becomes the page title

Lines 106-120

Asked as: wrong page title; second h1 heading; title taken from heading; extra level one heading

Turbo takes the first level-1 heading ("# Title") as the page title. A later level-1 heading is not a title, and it goes into the intro or the current section.

## What's new heading carries last commit

Lines 109-116

Asked as: what's new commit missing; commit suffix in heading; release heading commit hash; section shows commit id

A level-2 heading in a docs page that ends with '(to <hex>)' sets the section's last commit (commit). The hex must have 8 to 40 characters from 0-9 and a-f. The heading title drops that suffix. A heading without it has no commit.

## A page's word count covers its whole source

Lines 126-133

Asked as: page word count total; whole page word count; count of page words; word count for whole doc

Turbo counts the words of a page from its whole source text, including the title, intro and all sections. Each section counts the words of its own content only.

## Help page order

Lines 137-183

Asked as: help page order; pages in wrong order; help menu sequence; where does page appear

Help lists the pages in this order: start-here, time, purchase-orders, expenses, proposals, project-requests, rate-sheets, profile, glossary, whats-new (docPageOrder). A page missing from the list goes just before What's new, sorted by title.

## Child pages follow their parent in Help

Lines 158-196

Asked as: sub page not indented; child page under parent; help nested page; page label shorter in help

Some pages appear under a parent with a shorter label and an indent: review-time (Review and track), revise-a-purchase-order (Revisions), find-a-purchase-order (Find and print), child-purchase-orders (Child POs), commit-expenses (Review and commit), proposed-rates (Proposed rates), admin-profiles (Admin profiles), rate-sheet-periods (Periods), rate-sheet-absorb (Duplicates). Children follow the order of that list under their parent. Other pages use their title as the label.

## Help popover links open in a new tab

Lines 198-211

Asked as: help link opens same tab; popover link new window; link closes my form; help link lose unsaved work

In a help popover, every docs link opens in a new tab (target="_blank", rel="noopener"), so a person does not lose unsaved work in the form. This applies to text blocks, flow step details and notes.

## Help reference lookup

Lines 213-222

Asked as: help reference not found; section link shows nothing; help link to section empty; broken help reference

A help reference has the form 'page-slug#section-id'. Turbo shows the section only when both the page and the section exist. Otherwise it finds nothing.
