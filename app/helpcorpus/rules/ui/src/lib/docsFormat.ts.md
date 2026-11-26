<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/docsFormat.ts","sha256":"8c140f3d3456253a177ca7628d42774a7aab93299831d7bf0805da039f589f7a","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/docsFormat.ts

## Links in Help docs get the link style

Lines 33-39

Asked as: help doc link styling; links look plain in help; link class on help links; help link anchor style

When Turbo shows a Help doc, it renders each Markdown link as an anchor with the class "link" (the custom link renderer). Turbo adds no target or rel attribute to the link at this step.

## Commit links point to the GitHub repository

Lines 45-49

Asked as: what's new commit link; commit link goes to wrong place; github commit url; link to commit in github

A What's new commit link goes to https://github.com/tbt-eng/tybalt_turbo/commit followed by a slash and the commit id (commitUrl).

## Section ids come from the title

Lines 51-57

Asked as: section anchor id wrong; help section id from title; section link id missing; apostrophe in section id

The id of a section is its title in lowercase, with apostrophes removed and each run of other characters replaced by one hyphen, and leading and trailing hyphens cut.

## Word count of a Help section counts only English-letter words

Lines 59-61

Asked as: word count wrong in help; help section word count; numbers not counted as words; how words are counted

Turbo counts a word as a run that starts with a letter A-Z or a-z and continues with letters, digits, underscores, apostrophes or hyphens (countWords). Numbers on their own do not count. A text with no such word counts 0.

## Flow lines need exactly three parts

Lines 63-80

Asked as: flow line error; flow block syntax error; flow line has too many parts; actor step detail format

Each line in a flow block is 'Actor | Step | Detail' or 'Note: text'. A line with other than three parts, or with an empty part, raises the error: A flow line needs "Actor | Step | Detail". Blank lines are skipped.

## Flow note lines are kept apart from steps

Lines 67-72

Asked as: note line in flow; note shows as step; note: text in diagram; flow note not showing

In a flow block, a line that starts with "Note:" becomes a note, not a step. Turbo trims the text after "Note:" and renders it as inline Markdown. Blank lines are skipped.

## Flow step fields are rendered differently

Lines 73-78

Asked as: markdown not rendering in flow; flow step text formatting; bold in flow actor; detail markdown in flow

In a flow step, Turbo keeps the actor and the step as plain text. It renders only the detail as inline Markdown.

## Only code blocks marked flow become flow diagrams

Lines 82-100

Asked as: code block not a diagram; flow code fence not working; help code block shows as html; flow diagram not drawn

A code block whose language is "flow" becomes a flow block with steps and notes. Turbo renders all other content as HTML, and it ignores blank-space tokens. A flow block splits the HTML around it into separate blocks.

## A doc page needs a level-1 title

Lines 102-121

Asked as: doc page won't load; help page missing title error; doc needs title line; page fails no heading

Turbo refuses to load a doc page that has no '# Title' line. The error reads: Doc <slug> needs a "# Title" line. Only the first level-1 heading is the title.

## Only the first level-1 heading is the page title

Lines 102-120

Asked as: wrong title on help page; second h1 heading; which heading is title; intro text before sections

Turbo takes the first level-1 heading as the doc title. A later level-1 heading is not a title and stays in the content. Content before the first level-2 heading is the page intro. Each level-2 heading starts a new section.

## What's new headings carry a commit id

Lines 109-116

Asked as: what's new commit missing; commit id in heading; last commit not shown; heading with commit hash

A level-2 heading that ends with '(to <8 to 40 hex characters>)' sets the section's last commit and removes that text from the title. Other level-2 headings have no commit.

## Section and page word counts are calculated differently

Lines 126-134

Asked as: page word count differs; section word count off; headings counted in words; word count mismatch section page

A section's word count comes from the raw text of its own content, without its heading. The page's word count comes from the whole source, headings included.

## Help page order in the side navigation

Lines 137-198

Asked as: help page order wrong; side navigation order; pages out of order in help; where does my help page appear

Help lists top-level pages in the order of docPageOrder: Start here, Time, Purchase Orders, Expenses, Proposals, Project Requests, Rate Sheets, Clients, Absorb Proposals, Computers, Profile, Glossary, What's new. A page that is not in the list goes just before What's new, sorted by title.

## Child pages appear under their parent with a short label

Lines 161-204

Asked as: child page not nested; sub page label in help; help child page indent; child page under parent

Help lists a child page right after its parent, indented, with a shorter label (for example 'Close or cancel' for close-a-purchase-order). Children of one parent keep the order of the childDocPages list.

## Help popover links open in a new tab

Lines 206-219

Asked as: help popover link opens; popover links new tab; link in help popover behaviour; help popup links new window

In a help popover over a form, every link opens in a new tab (target=_blank, rel=noopener), including links in flow steps and notes. This keeps unsaved work.

## Help references find a page section by slug and id

Lines 221-230

Asked as: help link to section broken; help reference shows nothing; hash link to help section; help anchor not working

A help reference such as 'project-requests#approve-a-request' returns the page and section only when both exist. If either is missing, it returns nothing.
