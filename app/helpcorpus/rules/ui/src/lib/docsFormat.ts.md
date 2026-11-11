<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/docsFormat.ts","sha256":"b2c98b9457c8aa9d120fa60ee8d58bc8b66485011df60ae6a9a6402fe1b3331d","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/docsFormat.ts

## Help links use the link style

Lines 33-39

Asked as: link looks wrong in help; help link style; anchor styling in docs; markdown link appearance

Markdown links in Help docs render as anchors with the class "link" and the given address.

## Commit links point to the tbt-eng tybalt_turbo repository

Lines 45-49

Asked as: commit link in help; link to commit on github; help commit url; github commit address

A commit link in Help is the address https://github.com/tbt-eng/tybalt_turbo/commit followed by a slash and the commit hash (commitUrl).

## Section ids come from the title text

Lines 51-57

Asked as: section anchor id; heading id in help; why section link broken; help section id format

A section id is its title lowercased, with apostrophes removed, other runs of non-alphanumeric characters turned into a single hyphen, and leading or trailing hyphens trimmed.

## Help word counts

Lines 59-61

Asked as: help word count; how words are counted; numbers not counted as words; word count rule

Turbo counts a word as a run that starts with a letter and continues with letters, digits, underscores, apostrophes or hyphens. Numbers alone are not counted.

## Flow lines need exactly three parts or a Note

Lines 63-80

Asked as: flow line error; flow parsing failed; actor step detail format; flow line needs three parts

In a flow block, each non-empty line must be a "Note:" line or have exactly three non-empty parts: "Actor | Step | Detail". Otherwise parsing fails with: A flow line needs "Actor | Step | Detail": <line>. Note lines become notes. Blank lines are skipped.

## Flow code blocks become step lists in Help pages

Lines 82-100

Asked as: flow code block not showing; flow block shows as code; empty space in help page; flow steps not rendering

A code block marked "flow" in a doc becomes a flow block of steps and notes (parseFlow). Other content between flow blocks is rendered as Markdown HTML. Blank space tokens are dropped.

## Only the first level-1 heading is the page title

Lines 102-120

Asked as: two titles on help page; second h1 heading; which heading is page title; level one heading in body

The first level-1 heading in a doc is the page title. A later level-1 heading goes into the page body (the intro or the current section). Each level-2 heading starts a new section, and content before the first level-2 heading is the intro.

## Help page needs a level-1 title

Lines 107-121

Asked as: help page missing title; doc parse error no title; help page won't load; missing h1 in doc

A help page must have a level-1 heading ("# Title"). Without one, parsing fails with the error: Doc <slug> needs a "# Title" line. Only the first level-1 heading is the title.

## Help section headings can carry a commit reference

Lines 109-116

Asked as: commit hash in heading; what's new commit reference; heading with commit hash; section linked to commit

A heading ending with "(to " plus 8 to 40 lowercase hex characters and ")" becomes a section whose title leaves out that suffix. The hex value is stored as the last commit that the What's new week covers. A heading without it keeps its full text and has no commit.

## Section word count excludes the heading

Lines 131-133

Asked as: section word count wrong; word count excludes heading; page word count includes title; why word count differs

A section's word count (words) counts the words in its body tokens only. The page word count counts the words of the whole source, headings included.

## Help pages follow a fixed order

Lines 137-175

Asked as: help page order; sidebar page order; where new help page appears; help pages sorted how

Help lists pages in this order: start-here, time, purchase-orders, expenses, proposals, project-requests, profile, glossary, whats-new. A page that is not in the list goes just before What's new. Unlisted pages are sorted by title.

## Child pages show under their parent page

Lines 157-188

Asked as: revisions page under purchase orders; child page in side navigation; admin profiles under profile; sub page indent sidebar

The page revise-a-purchase-order shows under Purchase Orders with the label "Revisions". The page admin-profiles shows under Profile with the label "Admin profiles". A child page follows its parent and is indented in the side navigation. Children appear in the order of the child list.

## Help popover links open in a new tab

Lines 190-203

Asked as: help link opens new tab; popover link closes form; help link losing form data; links in help popover

In a help popover, every link opens in a new tab (target="_blank" rel="noopener") so unsaved form work is not lost. This applies to text blocks, flow step details and flow notes.

## Help references find a section by page and id

Lines 205-214

Asked as: help reference not found; link to help section broken; hash anchor in help doc; page section reference missing

A help reference such as "project-requests#approve-a-request" finds the page by slug and the section by id. If the page or the section is missing, nothing is found. A reference with no "#" part also finds nothing.
