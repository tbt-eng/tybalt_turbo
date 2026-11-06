# User docs: how to write them

This folder holds all user docs. Each file is one topic, and its name sets the URL: `project-requests.md` opens at `/docs/project-requests`. Turbo doesn't show this README.

Update these docs in the same commit as any change to what a user sees or does. If a user-facing commit needs no docs change, add a `Docs: none` line to its message. `scripts/check-docs-commit.mjs` warns about commits that have neither.

## Limits

`tests/docs.test.ts` fails if a page breaks these limits. `npm run check:ui` runs it.

- A section (`##`) has at most 60 words, because it must fit in a help popover.
- A page has at most 600 words. Aim for far fewer.
- A flow has at most 5 steps.

If a page doesn't fit, split it into two topics. Don't raise the limits.

## What to write

Write only what a person does and what happens next. Add a rule only if staff need it to finish the task. Don't write background, history or how the code works.

## Style

We use the writing rules of Simplified Technical English (ASD-STE100), plus our own word list in `glossary.md`.

- Use active voice.
- Start each step with a verb: "Select **Approve**", not "The request can be approved".
- Put one instruction in each sentence.
- Use at most 20 words in a step and 25 in a description.
- Use one word for one meaning. Use the glossary word, and add new words there.
- Use the exact labels from the screen, in bold.
- Describe people by what they do ("Procurement"), not by claim names.
- Don't state a configurable value as fixed. Link to the page that shows it.

`npm run check:ui` runs `scripts/check-docs-style.mjs`. It checks sentence length, steps that start with a verb, one instruction in a step, passive voice and the glossary's replaced words.

## What's new

Add a section for each complete week that has new features, with the newest week first. Weeks run from Sunday to Saturday, by the date a commit landed on `main`. Don't add a week until its Saturday has passed.

- Start each section with the date range and the last commit that it covers, such as `## November 1 to 7, 2026 (to af436ec9)`. Use the first 8 characters of the commit. Turbo shows it as a link to GitHub, and the next entry starts after it.
- Write one point for each feature, from the user's point of view. End it with its version tags.
- Leave out internal changes and anything behind a feature flag that is off.
- Remove weeks older than the last 12 complete weeks.

The word limits don't apply to this page, because its 12 weeks limit it.

## Flows

A `flow` block draws cards that show who acts and what happens next. Each line has an actor, a step and one line of detail. A `Note:` line covers an exception.

````markdown
```flow
You | Create the request | Enter the details and division hours.
Approver | Approve or reject | Check the request.
Note: Rejected? Edit the request and submit it again.
```
````

## Help in the app

`DSHelp` shows one section in a popover and links to the full page. Keep the text here, not in the component.

```svelte
<DSHelp doc="project-requests#approve-a-request" label="Help" />
```

The reference is the file name, then the section heading in lowercase with hyphens. Write it as plain text, not `doc={...}`, so `tests/docs.test.ts` can check it. The test fails if a reference or a link between docs doesn't resolve.
