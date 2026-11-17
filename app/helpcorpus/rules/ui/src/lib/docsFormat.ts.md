<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/docsFormat.ts","sha256":"ff88968b9b1fdab3e0b9206960fb39303ce66f9c65cb15f569f6e61b8496150d","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/docsFormat.ts

## Help links render as styled links

Lines 33-39

Asked as: help link styling; link opens same tab; docs link looks plain; anchor class link

Help renders each Markdown link in the docs as an anchor with the class "link". The link keeps its address (href) and its text. The anchor has no target attribute here, so it opens in the same tab. Popovers add the new-tab behavior separately.

## Commit link address

Lines 45-49

Asked as: commit link url; link to github commit; commit hash link; where does commit link go

A commit links to https://github.com/tbt-eng/tybalt_turbo/commit/ plus the commit hash.

## Section id comes from its title

Lines 51-57

Asked as: section anchor id; heading id from title; section id wrong; url fragment for section

A section id is its title in lowercase. Apostrophes are removed, other runs of characters that are not a-z or 0-9 become one hyphen, and leading or trailing hyphens are trimmed.

## Word counts

Lines 59-61

Asked as: word count wrong; how words are counted; page word count; section word count

Turbo counts a word as a run that starts with a letter and may contain letters, digits, underscores, apostrophes and hyphens. A section counts only its own content. A page counts its whole source.

## Flow block line format

Lines 63-80

Asked as: flow block error; flow line format; flow line needs three parts; flow step parse error

Each line of a 'flow' code block must be 'Actor | Step | Detail' with exactly three non-empty parts. Otherwise parsing fails with 'A flow line needs "Actor | Step | Detail": <line>'. Blank lines are skipped.

## Flow lines need three non-empty parts

Lines 67-78

Asked as: flow line error; flow step missing part; pipe separated line error; empty part in flow

A flow line that is not a note must have exactly three parts split by "|" (Actor | Step | Detail). A line with fewer or more parts, or with an empty part, makes Turbo throw the error: A flow line needs "Actor | Step | Detail" followed by the line. Blank lines are skipped.

## Flow notes

Lines 69-72

Asked as: flow note; note line in flow; note not a step; flow notes shown

A flow line that starts with 'Note:' is a note, not a step. Turbo shows it with the notes of that flow.

## Flow code blocks become flow steps

Lines 82-100

Asked as: flow code block; flow diagram not rendering; code block shows as html; fenced flow block

In a Help page, a fenced code block with the language "flow" becomes a flow block of steps and notes instead of HTML. The Markdown before it is closed into its own HTML block first. Empty space tokens are dropped, and the remaining content is grouped into HTML blocks.

## Only the first level-1 heading is the page title

Lines 106-111

Asked as: second title heading; page title wrong; level one heading title; which heading is title

The first level-1 heading ("# ") becomes the page title. A later level-1 heading does not become the title and is kept as normal content. A level-2 heading ("## ") starts a new section.

## A Help page needs a title

Lines 107-121

Asked as: help page no title; doc missing title error; page fails to parse; needs title line

A page must have a level-1 heading ('# Title'). Without one, parsing fails with 'Doc <slug> needs a "# Title" line'. Only the first level-1 heading is the title.

## What's new heading carries its last commit

Lines 109-116

Asked as: whats new commit in heading; commit shown in heading; heading commit hash; release commit in title

A level-2 heading that ends with '(to ' plus 8 to 40 lowercase hex characters and ')' sets the commit of that section. The commit is removed from the displayed title. Other headings have no commit.

## Content before the first section is the intro

Lines 117-119

Asked as: text before first section; page intro text; content above sections; intro paragraph missing

Content before the first level-2 heading is the page intro. Content after a level-2 heading belongs to that section.

## Parsed page keeps slug, intro and sections

Lines 122-134

Asked as: parsed help page structure; page sections list; slug and sections; help page data shape

A parsed Help page keeps the slug it was given. It holds the intro blocks and one section per level-2 heading, in file order. Each section has an id, a title, an optional commit and its content blocks.

## Section word count excludes its heading

Lines 126-133

Asked as: section words excluding heading; heading words counted; word count includes title; section count off

Each section's word count (words) counts only the raw text of the section's content tokens, not its heading. The page's word count (words) counts the whole source file, including headings.

## Help page order and child pages

Lines 137-191

Asked as: help page order; pages listed wrong; child page position; navigation order of help

Help lists the pages in the order of docPageOrder: start-here, time, purchase-orders, expenses, proposals, project-requests, rate-sheets, profile, glossary, whats-new. A child page (childDocPages) shows right after its parent, in the order of the childDocPages list.

## Child pages use a shorter label and indent in navigation

Lines 158-197

Asked as: sub page indent; child page label in nav; short name in sidebar; side menu indentation

In the side navigation, a child page shows its short label (for example 'Revisions' or 'Child POs') and is indented. Other pages show their own title with no indent.

## Unlisted Help pages go before What's new

Lines 173-191

Asked as: new help page position; page missing from menu; unlisted page placement; page sorted before whats new

A page that is in neither docPageOrder nor childDocPages goes just before What's new. Such pages are sorted by title among themselves.

## Help popover links open in a new tab

Lines 199-212

Asked as: popover link new tab; help link closes form; link in popup opens tab; unsaved form lost on link

In a help popover, every link opens in a new tab (target _blank, rel noopener), in text, flow steps and notes, so that unsaved form work is kept.

## Help reference lookup returns nothing if page or section is missing

Lines 214-223

Asked as: help link not found; broken help reference; section link does nothing; missing help anchor

A help reference such as 'project-requests#approve-a-request' finds a page by slug and a section by id. If either is missing, or the reference has no '#id', Turbo finds nothing.
