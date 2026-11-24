<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/docsFormat.ts","sha256":"a31111c6541a847a96e60673f07b2325234d49bc37a70a7504821c215ec7f5f3","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/docsFormat.ts

## Links in help pages get the link style

Lines 33-39

Asked as: help page link style; link looks plain; link color in help; help link open same tab

Every Markdown link in a help page renders as an anchor with the CSS class "link" and the given address. Turbo does not add a new-tab setting here.

## What's new commit links point to GitHub

Lines 45-49

Asked as: commit link goes where; whats new commit github; link to commit hash; github commit url

A commit link for a What's new section goes to https://github.com/tbt-eng/tybalt_turbo/commit/ followed by the commit hash (commitUrl).

## Section ids and word counts

Lines 51-61

Asked as: section id from heading; anchor name for section; word count wrong; how words are counted

A section id is its title in lowercase, with apostrophes removed and each run of other characters replaced by one hyphen, and no hyphen at the start or end. Word counts count only words that start with a letter.

## Flow block lines need three parts

Lines 63-99

Asked as: flow block error; flow line needs three parts; actor step detail format; flow step missing field

In a doc, a code block marked "flow" lists steps as "Actor | Step | Detail". A line that does not have exactly three non-empty parts fails with the error "A flow line needs "Actor | Step | Detail": <line>". Blank lines are skipped. A line that starts with "Note:" becomes a note.

## Doc page needs a level-1 title

Lines 102-121

Asked as: doc missing title; doc parser error no title; page needs heading one; doc fails to load

Turbo's doc parser uses the first level-1 heading as the page title. A doc without one fails with the error "Doc <slug> needs a "# Title" line".

## What's new heading carries its last commit

Lines 106-120

Asked as: whats new commit suffix; heading shows commit hash; to commit in heading; section last commit

A level-2 heading that ends with "(to <8 to 40 hex characters>)" gives the section its last commit (commit), and the title drops that suffix. A level-2 heading without that ending has no commit. Text before the first level-2 heading goes to the page intro.

## Text before the first section heading is the page intro

Lines 122-126

Asked as: page intro text; text before first heading; help page introduction; content above first section

A help page's intro (intro) holds the content between its level-1 title and its first level-2 heading. Each level-2 heading starts a new section. The page keeps its slug and title.

## Section id comes from the section title

Lines 126-130

Asked as: section anchor id; heading to section id; section link not working; trailing to commit heading

Each section's id is its title converted by slugify. The title is the heading text without any trailing "(to commit)" part. A reference like page#section-id matches this id.

## Section and page word counts cover different text

Lines 131-133

Asked as: page word count total; section word count only; word count includes title; count words per section

A section's word count (words) counts only the text of that section. A page's word count counts its whole source, including the title, the intro and all sections.

## Help page order and unlisted pages

Lines 145-196

Asked as: help page order; pages listed out of order; where page appears in help; unlisted help page position

Help lists the pages in the order of docPageOrder: start-here, time, purchase-orders, expenses, proposals, project-requests, rate-sheets, clients, absorb-proposals, computers, profile, glossary, whats-new. A page that is not in the list goes just before What's new, sorted by title.

## Child pages sit under their parent in Help

Lines 163-202

Asked as: child page under parent; indented help page; sub page label in help; review and track page

Some pages show under a parent page with a shorter label, for example review-time under time as "Review and track" and rate-sheet-delete under rate-sheets as "Delete". Help shows a child right after its parent, indented, in the order of the childDocPages list. A page with no entry in that list is not indented and uses its own title.

## Help popover links open in a new tab

Lines 204-217

Asked as: help popover link new tab; link in popover opens tab; lose form work following link; popover link keeps form

In a help popover, Turbo adds target="_blank" and rel="noopener" to links in text, flow step details and notes. Staff keep their unsaved form work when they follow a link.

## Help reference lookup by page and section

Lines 219-228

Asked as: help reference not found; link to page section broken; project-requests anchor missing; help lookup finds nothing

A help reference such as "project-requests#approve-a-request" finds a page by its slug and a section by its id (the slugified heading). If the page or the section does not exist, Turbo finds nothing.
