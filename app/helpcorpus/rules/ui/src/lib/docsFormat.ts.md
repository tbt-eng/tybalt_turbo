<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/docsFormat.ts","sha256":"19aca3c2d47a540eda5cbacc393eeb613295a8578889e8362d5f06bc46c5863f","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/docsFormat.ts

## Links in help render as styled links

Lines 33-43

Markdown links in help render as <a> tags with the "link" class. Detail text in flow steps and note text also support inline Markdown.

## Commit links point to the GitHub repository

Lines 45-49

A commit link goes to https://github.com/tbt-eng/tybalt_turbo/commit/ followed by the commit hash (commitUrl).

## Section links use the slug of the title

Lines 51-57

Turbo makes a section id (id) from its title in lowercase. It removes apostrophes (' and ’) and turns each run of other characters into one hyphen. It then trims a hyphen at the start or end.

## Word counts use words that start with a letter

Lines 59-61

Turbo counts a word as a run that starts with a letter (A-Z or a-z) and continues with letters, digits, underscores, apostrophes or hyphens. Numbers alone are not counted. A section's count uses the raw text of its tokens; a page's count uses the whole source, including headings.

## Each flow line needs Actor | Step | Detail

Lines 64-78

Each line in a flow block must have exactly three non-empty parts split by "|": actor, step and detail. Otherwise Turbo refuses the page with: A flow line needs "Actor | Step | Detail": <line>. Blank lines are skipped.

## Flow lines starting with Note: become notes

Lines 69-72

A flow line that starts with "Note:" is not a step. Turbo adds the text after it to the flow's notes instead, and does not check it for three parts.

## Flow code blocks become step lists

Lines 82-100

A code block marked "flow" in a help page becomes a flow of steps and notes. Other Markdown around it is rendered as normal HTML blocks. Blank space tokens are dropped.

## Help docs need a level-1 title heading

Lines 106-121

Turbo refuses to load a help page whose Markdown has no level-1 heading ("# Title"). The error says: Doc <slug> needs a "# Title" line. Only the first level-1 heading becomes the page title.

## Level-2 headings split a help page into sections

Lines 106-125

Each level-2 heading ("##") starts a new section on a help page. Content before the first level-2 heading becomes the page intro. Later level-1 headings are not titles and go into the current section or the intro.

## What's new heading ends with its last commit

Lines 109-116

A section heading that ends with "(to <hex>)" gives the section a commit (commit). The hex value must have 8 to 40 characters from 0-9 and a-f. Turbo removes this suffix from the section title. A heading that does not match keeps its full text and has no commit.

## Revisions page is listed under Purchase orders

Lines 137-160

The page "revise-a-purchase-order" is a child of "purchase-orders". Help shows it under that parent with the label "Revisions" and an indent in the side navigation. Other pages use their own title with no indent.

## Help page order in the list

Lines 142-154

Help lists Start here (slug start-here) first and What's new (whats-new) last. All other pages sit between them, sorted by title. A child page is not listed on its own but follows its parent.

## Help popover links open in a new tab

Lines 162-175

In a help popover over a form, links open in a new tab (target="_blank", rel="noopener"), so staff do not lose unsaved work. This applies to HTML blocks and to flow step details and notes.

## Help references need a page and a section

Lines 177-186

A help reference such as "project-requests#approve-a-request" finds a page by slug and a section by id. If the page or the section does not exist, Turbo finds nothing (null). A reference with no "#" part finds no section.
