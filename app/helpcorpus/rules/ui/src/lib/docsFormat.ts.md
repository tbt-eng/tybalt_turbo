<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/docsFormat.ts","sha256":"223644d75d381d6bdd8e62024238a66583a22c2a2e93bcbfc1907fdcab6e3ab7","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/docsFormat.ts

## Links in help pages use the link style

Lines 33-39

Asked as: link style in help; help page link looks wrong; markdown link anchor class

Every Markdown link in a help page becomes an anchor with the CSS class "link" and the same address (href) that the doc source gives.

## Help links to commits use the GitHub commit address

Lines 45-49

Asked as: link to a commit; commit link in help; github commit url

A link to a commit in Help points to https://github.com/tbt-eng/tybalt_turbo/commit followed by a slash and the commit hash (commitUrl).

## Section ids come from the title

Lines 51-57

Asked as: section anchor id; section link id from title; heading id in help

A section's link id is its title in lower case, with apostrophes removed and every run of characters other than a-z or 0-9 replaced by one hyphen. Leading and trailing hyphens are dropped.

## Word counts for help pages

Lines 59-61

Asked as: word count for help page; how words are counted; help page word total

Turbo counts a word as a run that starts with a letter and may contain letters, digits, underscores, apostrophes and hyphens. A section's count uses only its own text.

## Flow line format

Lines 63-80

Asked as: flow line error; flow line needs actor step detail; flow line format wrong

Each line of a flow block must be 'Actor | Step | Detail', with exactly three non-empty parts. Otherwise Turbo refuses it with the error 'A flow line needs "Actor | Step | Detail": <line>'. A line that starts with 'Note:' becomes a note instead, and blank lines are skipped.

## Flow blocks become steps in the help page

Lines 82-100

Asked as: flow block not showing; flow code block as steps; flow block shows as code

A code block marked "flow" in a doc becomes a flow block with steps and notes, not ordinary code. The rest of the Markdown around it becomes HTML blocks in order. Blank (space) tokens are dropped.

## Doc sections split at level-2 headings; level-1 sets the title

Lines 102-120

Asked as: doc sections split; level 2 heading new section; help page title from heading

The first level-1 heading is the page title. Each level-2 heading starts a new section. Content before the first level-2 heading is the page intro. A later level-1 heading goes into the current section's content, not the title.

## Doc page needs a title

Lines 106-121

Asked as: doc page won't load; needs a title line; help page missing title

A doc page needs a level-1 heading. The first one becomes the page title. Without one, Help refuses to load the page with the error 'Doc <slug> needs a "# Title" line'.

## What's new heading carries its last commit

Lines 109-116

Asked as: whats new commit on heading; to commit in section heading; release heading commit hash

A level-2 heading that ends with '(to ' plus a hex commit of 8 to 40 characters sets that section's commit. Help removes the '(to ...)' part from the section title. A heading without it has no commit.

## Section word count excludes the heading, page count includes all

Lines 126-134

Asked as: section word count; page word count includes heading; word count excludes heading

A section's word count covers its content tokens (raw text) and not its heading. The page's word count covers the whole source.

## Help page order and unlisted pages

Lines 137-189

Asked as: help menu page order; where pages appear in help; unlisted help page position

Help lists top-level pages in the order of docPageOrder: start-here, time, purchase-orders, expenses, proposals, project-requests, rate-sheets, profile, glossary, whats-new. A page not in that list goes just before What's new, sorted by title.

## Child pages follow their parent in Help

Lines 158-195

Asked as: child page order in help; sub page under parent; indented help sub page label

Help lists a child page right after its parent, in the order of the childDocPages list. Child pages are indented in the side navigation and show a shorter label, such as 'Review and track', 'Revisions', 'Find and print', 'Child POs', 'Review and commit', 'Admin profiles', 'Periods' or 'Duplicates'.

## Help popover links open in a new tab

Lines 197-210

Asked as: help link opens new tab; popover link new window; help popover links lose form

In a help popover, every link in the text, flow details and notes opens in a new tab, so a person does not lose unsaved work in the form.

## Help references need page and section

Lines 212-221

Asked as: help link to section not working; help reference finds nothing; hash link to help section broken

A help reference such as 'project-requests#approve-a-request' finds a section only when both the page slug and the section id exist. If either is missing, or the reference has no '#', Turbo finds nothing.
