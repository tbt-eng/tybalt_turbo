<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/docsFormat.ts","sha256":"ce100fb19a5098946288e59e19938196624db0bca333127cb03c5e3703fb378a","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/docsFormat.ts

## Docs links get the link style

Lines 33-39

Asked as: link style in docs; docs links look plain; anchor class link; markdown link styling; links not styled in help

Turbo renders every Markdown link in the docs as an anchor with the class "link". The link text keeps its inline formatting.

## Commit links point to the GitHub repository

Lines 45-49

Asked as: whats new commit link; link to github commit; commit hash url; github repository link; commit link wrong repo

The link for a What's new commit is https://github.com/tbt-eng/tybalt_turbo/commit/ followed by the commit hash.

## Second-level headings split a doc into sections

Lines 51-57

Asked as: section split in doc; second level heading sections; doc page sections; section id from heading; heading creates new section

Each second-level heading starts a new section. Content before the first one is the page intro. A heading that ends with "(to <8 to 40 hex characters>)" sets the section's last commit (commit), and Turbo removes that suffix from the title. A section id is its title in lowercase, with apostrophes removed and other runs of characters that are not letters or digits turned into hyphens.

## Docs word counts count only words that start with a letter

Lines 59-61

Asked as: word count docs; numbers not counted words; section word count; page word total; how words are counted

Turbo counts a word as a run that starts with a letter and may contain letters, digits, underscores, apostrophes and hyphens. Numbers alone do not count. A section's count uses its own content, and the page count uses the whole source.

## Flow blocks need three parts per line

Lines 63-80

Asked as: flow line error; flow step needs three parts; actor step detail format; flow line missing part; note line in flow block

A code block marked flow lists steps as "Actor | Step | Detail". A line with other than three parts, or with an empty part, throws the error: A flow line needs "Actor | Step | Detail": <line>. A line that starts with "Note:" is a note, not a step. Blank lines are skipped.

## A flow code block becomes a flow block

Lines 82-100

Asked as: flow code block; flow block not rendering; code block language flow; steps shown as code; flow steps display

A code block in a doc page whose language is "flow" becomes a flow block with steps and notes. Turbo renders all other content as HTML. Blank space between blocks is dropped.

## A doc page without a title fails to load

Lines 102-121

Asked as: doc page fails to load; doc missing title error; page will not load no title; missing title line; doc slug error

A doc page needs a first-level heading ("# Title"). Without one, Turbo throws the error: Doc <slug> needs a "# Title" line. Only the first first-level heading is the title.

## Only the first level-1 heading is the page title

Lines 106-108

Asked as: two h1 headings; second title in doc; multiple top level headings; which heading is page title; extra level one heading

The first level-1 heading becomes the page title. Any later level-1 heading is treated as ordinary content.

## A What's new heading ending in (to commit) sets the section's commit

Lines 109-116

Asked as: commit suffix on heading; heading shows to commit; whats new heading hash; commit removed from title; whats new commit missing

A second-level heading that ends with "(to " plus 8 to 40 lowercase hex characters and ")" sets the section's last commit (commit). Turbo removes that suffix from the section title. A heading without it has no commit.

## Content before the first section is the page intro

Lines 117-119

Asked as: text before first heading; intro above sections; page intro text; content above first heading; doc introduction paragraph

Content that comes before the first second-level heading becomes the page intro. Content after a second-level heading belongs to that section.

## Section ids come from the section title

Lines 126-133

Asked as: section anchor id; slug from section title; section link id wrong; section id changes; commit text in section id

Turbo sets each section's id (id) by slugifying its title, with commit text removed. Each section's word count covers only that section's own text. The page word count covers the whole source file.

## Help page order is fixed, with unlisted pages before What's new

Lines 137-156

Asked as: help menu page order; help pages sorted wrong; page order in help; unlisted page position; pages before whats new

Help lists top-level pages in the order of docPageOrder: start-here, time, purchase-orders, expenses, proposals, project-requests, rate-sheets, profile, glossary, whats-new. A page that is not in the list goes just before What's new, sorted by title.

## Child pages show under their parent page in Help

Lines 158-194

Asked as: subpage under help page; child page in side menu; nested help pages; indented help navigation; review and track placement

Seven pages show under a parent with a shorter label: Review and track (time), Revisions, Find and print and Child POs (purchase orders), Review and commit (expenses), Admin profiles (profile), and Duplicates (rate sheets). A child page follows its parent, in the order of the child list, and is indented in the side navigation.

## Help popover links open in a new tab

Lines 196-209

Asked as: help popover link opens; link in help tooltip new tab; popover link losing form; help link same window; open help link safely

Links inside a help popover open in a new tab (target=_blank, rel=noopener). A person does not lose unsaved work on the form under the popover. This applies to text blocks, flow step details and flow notes.

## Help reference needs a page and a section

Lines 211-220

Asked as: help reference not found; help link to section fails; page#section link broken; section link does nothing; help anchor not working

A help reference such as "project-requests#approve-a-request" finds a section only when both the page slug and the section id exist. Otherwise Turbo finds nothing, and a reference with no section id also finds nothing.
