# Help chat

Help chat answers questions about Turbo. A signed-in user asks a question in Help, and a language model answers from the user docs, the design notes and the code. The model reads only. It has no tool that writes.

This file is the technical record: how to set the chat up, how each part works, and how we measure it. The design doc "Turbo Help Chat" has the diagrams, the safety rules, the threats and the plans for later phases.

## Status

| Part | State |
| --- | --- |
| Search of the user docs in Help (`ui/src/lib/docsSearch.ts`) | Live for every user |
| Corpus generator, rule files and collection files (`app/helpcorpus`, `app/cmd/helpcorpus`, `app/cmd/helprules`, `app/cmd/helpcollections`) | Built. The Docker build runs them. |
| Search index with the claims filter (`app/helpcorpus/index.go`) | Built |
| Chat endpoint and log (`app/helpchat`, `chat_turns`) | Built. Used with Claude Haiku on a copy of the production data. |
| Ask in the UI (`ui/src/routes/docs/+layout.svelte`, `DSChat`, `ui/src/lib/stores/helpChat.ts`) | Built |
| Benchmark (`app/cmd/helpeval`) | Built. Not run since the prompt changes of 2026-11-10. |
| Records in answers, and drafts | Later phases |

The code has been in production since 2026-11-10 (version 0.3.176). The chat is off until a superuser adds the `help_chat` setting. See "Set up the model" below.

## Set up the model

The chat needs two things: a key on the server, and a setting in the database. The key is a secret. The setting is not.

### 1. Choose a provider

| Provider | `provider` value | Secret for the key | Models |
| --- | --- | --- | --- |
| Claude, from Anthropic | `anthropic` | `ANTHROPIC_API_KEY` | `claude-haiku-5-5` (the default) or `claude-sonnet-5-5` |
| OpenAI | `openai` | `OPENAI_API_KEY` | Any model id from OpenAI's model list |
| OpenRouter | `openrouter` | `OPENROUTER_API_KEY` | Any model id from OpenRouter's model list |

Use `anthropic` for a Claude model, and not `openrouter`. Only the `anthropic` client marks the shared prompt for Claude's cache. Without the mark, each question pays the full price for the handbook.

### 2. Put the key on the server

In production, the key is a Fly secret:

```bash
fly secrets set ANTHROPIC_API_KEY=your-key
```

Fly restarts the app when a secret changes. You can set the keys of more than one provider, and then change between them with the setting alone.

In development, set the key in the terminal that runs the server:

```bash
export ANTHROPIC_API_KEY=your-key
```

```bash
go run main.go serve
```

Do not put a key in the repo, in `ui/.env`, or in `app_config`. Every signed-in user can read `app_config`.

### 3. Add the setting

1. Open the PocketBase admin at `/_/` and sign in as a superuser.
2. Open the `app_config` collection.
3. Add a record. Set `key` to `help_chat`, and set `value` to one of these.

```json
{ "enabled": true }
```

```json
{ "enabled": true, "model": "claude-sonnet-5-5" }
```

```json
{ "enabled": true, "provider": "openai", "model": "gpt-5-mini" }
```

```json
{ "enabled": true, "provider": "openrouter", "model": "openai/gpt-5-mini" }
```

The first one uses Claude Haiku, which is the default. The model ids for OpenAI and OpenRouter are examples: copy the id from the provider's own list. [app_config.md](app_config.md) lists every property.

### 4. Change the model, or turn the chat off

Edit the `value` of the `help_chat` record. The next question uses the new values. The app does not restart.

- To change the model, change `model`.
- To change the provider, change `provider` and `model` together. The server must have that provider's key already (step 2).
- To turn the chat off, set `enabled` to `false`.
- To give the chat to a test group first, see the next step.
- `effort` is `"low"` unless you set it. For an OpenAI model that has no reasoning setting, set `"effort": ""`.

Each row of `chat_turns` records the provider and the model that answered, so you can compare two models on real questions.

A property that the chat does not know is a bad setting, and the chat is then off for everyone (HTTP 503). So a mistake in the name of a property can't remove a limit.

### 4a. Give the chat to a test group first

A deploy group turns the chat on for some users and leaves it off for the rest. [app_config.md](app_config.md#deploy-groups) says how deploy groups work.

1. In the PocketBase admin, add a record to `deploy_groups` with a `name`, such as `help chat testers`. Copy its `id`.
2. Add one record to `deploy_group_memberships` for each tester: the group and the user.
3. Set the group's id in the `help_chat` record:

```json
{ "enabled": true, "deploy_group": "THE_GROUP_ID" }
```

A member sees **Ask Help** and can ask. Every other user sees Help with its search and no **Ask Help**, and `/api/help/ask` returns HTTP 403 for them. To add a tester, add a membership record: it takes effect at once. To give the chat to everyone, remove `deploy_group` from the record.

Each user still has their own limit of questions in an hour, so the group's size sets the most that a test can cost. With few users, more questions find the model's cache cold and pay to write it again. So the cost of a question in a small test is higher than it will be with all users: use the benchmark to estimate the cost.

### 5. Check that it works

Open Help and type a question in the search bar. **Ask Help** is the first row, and Enter sends the question. Then look at the turn's row in `chat_turns`: after the first call of a few minutes, `tokens_cache_read` must not be zero.

To test the route with no UI, send a question from a terminal. Copy your token from the browser: it is the `token` in the `pocketbase_auth` item of the site's local storage.

```bash
curl -N -X POST http://localhost:8090/api/help/ask -H "Authorization: YOUR_TOKEN" -H "Content-Type: application/json" -d '{"question":"How do I submit my week?"}'
```

The answer arrives as `text` events, and the last event is `done`.

| What you get | What it means |
| --- | --- |
| HTTP 401 | You are not signed in. |
| HTTP 403 | The chat is off: there is no `help_chat` record, or `enabled` is false. |
| HTTP 503 | The setting has a bad value, or the server has no key for the provider. The server log has a line that starts "help chat" and gives the reason. |
| `done` with status `error` | The provider refused the request or failed. The `error` field of the turn's row in `chat_turns` holds the provider's message. |

### Add another provider

A provider that speaks the OpenAI format needs one entry in `app/helpchat/providers.go`: its name, the secret for its key, and its address. A provider with another format needs a new file that implements `Model`.

### What is tested

Tests send real requests from both clients to local servers that speak each API's format, and a browser test runs the UI against a stand-in server. The benchmark and real chats have used Anthropic with `claude-haiku-5-5` and `claude-sonnet-5-5`. Nobody has sent a question to OpenAI or OpenRouter from this code. Do step 5 with one of them before you rely on it.

## Rules: what the code does, in plain words

The code often has none of the words of a question. The line that approves an expense at submit, when its owner is also its approver, compares three ids. It has no "self", "own" or "automatic". A search with the user's words, or the model's, does not find it.

A rule file fixes this. It says in plain words what one source file does: who may do an action, when Turbo refuses it and with what message, the limits and the defaults, and the cases that give a result a person does not expect. Each rule names the lines of the code that make it true.

- The rule files are in `app/helpcorpus/rules/`, one for each source file, with the path of the source and `.md`. They are Markdown, so a person can read them.
- A model writes them from the code: `go run ./cmd/helprules` in `app/`, with the provider's key in the shell. The files are in the repo, so the app's build calls no model.
- Files in `app/hooks`, `app/routes`, `app/utilities`, `app/reports`, `app/notifications`, `app/cron`, `app/constants` and `app/absorb` get rules, and the `.ts` files directly in `ui/src/lib`. `app/helpcorpus/rules.go` has the list.
- A rule file records the hash of the source that it was written from. The corpus takes a rule file only when the source has that hash today. A source that changed has no rules until someone runs the command again. So a rule can be missing, and can't be out of date. This is the difference from a design note.
- The command writes rules only for a source that changed, so a second run costs nothing. `go run ./cmd/helprules --check` calls no model, lists the sources whose rules are missing or old, and fails if there are any.
- A rule can still have other words than a question. So each rule also has search phrases, in an `Asked as:` line: 3 to 5 short phrases that a person could type about that case, such as "my own manager" for the rule about an owner who is also the approver. The search matches them, and the model does not see them. The same command writes them in a second step, with a cheaper model (`claude-haiku-5-5`), from the rule and not from the code.
- Each rule is one chunk of the corpus, with the path and the lines of its code. A source link for a rule opens the code. A rule that names uncommit is left out, as the chat treats a commit as final.
- One search returns up to 10 rules of the code, then up to 3 rules of the collections, then the code, then the design notes. The prompt tells the model to answer from a rule as from the code, and to read the lines when it needs a detail that the rule does not give.

After a change to a source file that has rules, run the command and commit the rule file with the change. Read the new rules: a model wrote them, and a wrong rule gives wrong answers.

A rule file is tied to its source file in two ways: it has the path of the source, with `.md`, under `app/helpcorpus/rules/`, and it records the hash of the source. Nothing writes a rule file in a build or in a test, because that needs a model and its key. So three checks make sure that a person wrote them before the change reaches production. Each one fails when a source file changed and its rule file did not, when a rule file is missing, when a rule has no search phrases, or when a rule file has no source:

- A test: `go test ./helpcorpus/` (`TestTheRuleFilesMatchTheSource`).
- The pre-push hook, `.githooks/pre-push`. It refuses a push to main, and checks the files of the commit that is pushed. It checks no other branch, so work in progress can be pushed. Turn the hooks on one time in each clone: `git config core.hooksPath .githooks`.
- The Docker build: `go run ./cmd/helpcorpus --require-rules` stops the build, so a push that skipped the hook can't deploy.

The instructions and the glossary are the same in every request, so the model caches them. A first run for all 280 source files costs a few dollars with `claude-sonnet-5-5`. After that, a run costs a few cents for each file that changed.

The command refuses a rule that is a placeholder: a title or a text of fewer than three characters. A model once wrote a rule with the title "x" among 29 real ones, and the search phrases of that file then failed on every run.

### The access rules of the collections

The UI reads and writes many collections through PocketBase's own API: profiles, time entries, expenses, jobs, vendors and more. For those requests the collection's rules decide who may list, view, create, update and delete a record. The rules are in the migrations, and the corpus leaves the migrations out, because an old migration says what a rule was and not what it is. So the chat could not answer "who can change a profile?", and it answered a question about a profile from the code of `admin_profiles`, which has far more of it.

- `go run ./cmd/helpcollections` applies every migration to an empty database and writes one file for each collection to `app/helpcorpus/collections/`: its five rules as they are at the end, its fields with their limits, its indexes, and the query of a view. It calls no model and needs no key.
- The files are a source of the corpus, as code, and each one gets a rule file like any other source. The rule file says the same in plain words.
- A search returns the rules of these files as a group of their own: up to 3, after the rules of the code. The 64 files have about 480 rules, and most give the limit of a field or the query of a view. In one group with the rules of the code they took 4 of the 10 places in a search on PO approver limits, and the rule that answers that question lost its place. The text of a collection file is not a search hit; `read_source` shows it.
- The heading of such a rule starts with its collection and the part of the file that it is about, such as `profiles collection, update rule:`. The corpus adds this from the lines of the rule, so it is the same in every build. A model wrote the title of the rule in plain words, and many collections have "profiles" in their name: without the start, a search for `profiles update rule` returned rules of three other collections.
- The prompt tells the model to search for the name of the collection and the action when the user asks who can see, add, change or delete a kind of record and a search of the code did not show who. Before that, the model asked about "profile default division manager", got rules about fields, and said that it could not confirm the answer. A first form of this line said "also search" for every such question: the model then wrote such a search for 24 of the 100 older questions, in the place of a search of the code. With the line as it is now, it does so for 10.
- Run the command after a migration that changes a collection, then `go run ./cmd/helprules`, and commit both. A test (`go test ./cmd/helpcollections/`) and the Docker build fail when the files are not what the migrations give.
- The rules in the files were compared with the schema of a copy of the production database: all 275 are the same.

## Measure the answers

`app/cmd/helpeval` is the benchmark. `questions.json` holds 104 hard questions about how Turbo behaves, 10 or 11 in each of ten areas, and 4 on the access rules of the collections. The ten areas: PO approval, the PO lifecycle, expenses, time, jobs, reports, people and permissions, clients and vendors, other spending rules, and the project workflow with its emails. Their sources are 107 of the 216 source files that have rules, and those files hold two rules in three. Each one has an ideal answer, the key points that a correct answer states, the claims that make an answer wrong, and the lines of source that settle it. An engineer checked each answer against those lines. The corpus leaves out `app/cmd`, so the chat can't find the file.

Run it from `app/`, with the provider's key in the shell:

| Command | What it measures | Cost |
| --- | --- | --- |
| `go run ./cmd/helpeval --mode search` | The question itself as one search. Did the result hold the right lines? | None. No model. |
| `go run ./cmd/helpeval --mode replay --out <run>` | The searches that the model wrote in an earlier run, sent to the search as it is now. Use it to test a change to the search or the chunks before a run that costs money. | None. No model. |
| `go run ./cmd/helpeval --mode ask` | The chat as it is: the same prompt, tools and tool loop as `/api/help/ask`. | Model and judge |
| `go run ./cmd/helpeval --mode oracle` | The chat with the right source in its prompt. This is the best that a perfect search can give. | Model and judge |
| `go run ./cmd/helpeval --mode notools` | The chat with the handbook and no tools. | Model and judge |
| `go run ./cmd/helpeval --mode judge-check` | The judge, on answers whose grade is known: the reference answers, empty answers, and the answers in `judge_cases.json`. Run it after a change to the judge or to the questions. | Judge |

`--model`, `--effort` and `--provider` choose the model that answers. `--only` and `--limit` choose questions. A run writes `results.jsonl`, `summary.md` and one trace for each question to `cmd/helpeval/runs/`, which git ignores. A second run with the same settings asks only the questions that have no grade. `--mode regrade --out <folder>` grades the saved answers again.

How to read a run:

- A second model, the judge, says for each key point if the answer states it, misses it or contradicts it, and if the answer makes one of the wrong claims that the question lists. The code makes the grade: **correct** (every key point, nothing wrong), **partial** (some key points, nothing wrong), **wrong** (it contradicts a key point or makes a listed wrong claim), or **abstained** (the chat said that it does not know). A wrong answer is worse than an abstained one.
- The judge can also **flag** an answer: a statement that it takes for an error and that is not a key point or a listed claim. A flag does not change the grade. The judge reads a few lines of the source, and a flagged statement is often true and about lines that it did not have. The summary lists the flagged answers. Read each one against the code. If the statement is false, add it to the question's `wrong_if`, so that the next run grades it.
- `--mode rescore --out <folder>` makes the grades again from the saved verdicts, with no model. Use it after a change to how a verdict becomes a grade.
- A failure of the API is an **error**. It is not in the numbers.
- `ask` compared with `oracle` shows what the search costs. `oracle` compared with a stronger model shows what the model costs. `ask` compared with `notools` shows what the tools give.
- 100 questions asked one time can't show a difference of less than about 10 points. For a smaller difference, run again and compare question by question.
- The judge reads the question, the reference, the key points and the lines of the source that settle the question. With the source it can tell a true statement that the reference does not have from a false one. It still makes mistakes. `judge_cases.json` holds answers that it graded wrong in earlier runs, with the grade that a person gave after reading the code, and answers that are wrong and that it must still catch. `judge-check` says on how many it agrees. When the judge is wrong about an answer in a run, add that answer to the file. Read the answers that it grades as wrong before you act on a number.

The first run with all 100 questions, on 2026-11-10 with `claude-haiku-5-5` at low effort, the rules and their search phrases, and the judge that reads the source: 90 correct, 10 partial, none wrong, and 94% of the key points stated. The searches returned the right lines for 96 questions. The judge flagged three answers, and each flagged statement was true. Of the first 50 questions, 42 were correct. Of the 50 new ones, 48 were. The new questions are easier for this search: their writers found the cases with the help of the rule files, so the questions and the rules use the same words.

Two runs with the rules of the collections and the four questions on access, on the same day:

| Run | Access questions correct (of 4) | The 100 older questions: correct | Partial | Wrong | Searches for a collection rule in them |
| --- | ---: | ---: | ---: | ---: | ---: |
| Before the rules of the collections (the run above) | | 90 | 10 | 0 | 0 |
| With them, and a prompt line that said "also search" for a collection rule | 4 | 84 | 15 | 1 | 25 |
| With them, and the prompt line as it is now | 4 | 87 | 11 | 2 | 10 |

- Before the collection rules had a group of their own and the prompt named them, the same four questions got 1 correct answer.
- Between the last two runs, which differ by one line of the prompt, 10 of the 100 older grades changed: 6 went up and 4 went down. So a difference of 3 between two runs is chance. The first run is also not a clean base: the design notes and the Claude client changed between it and the others.
- One loss has a known cause. A rule about a field of `client_notes` said that the job is optional, which is true of the database, and a hook asks for a job or a mark that the note is not for a job. The chat said that the job is optional in both runs. A rule about the fields of a collection now ends with a sentence that says that the code can ask for more. With it, the question was asked three more times: the search returned the same rule of the field each time, and all three answers were correct.
- The cost for a question did not change: about $0.0025.

The two modes with no model give the same number on every run, so they can show a small change to the search, which a run with a model can't. For how many of the first 50 questions does the search return the lines that settle the answer?

| The search | The question itself, as one search | The searches that the model wrote |
| --- | ---: | ---: |
| One list of 6 chunks, notes and code together | 14 | 27 |
| Code first, then notes | 23 | 45 |
| With the rules: 5 rules, then code, then notes | 35 | 46 |
| With the search phrases of the rules | 38 | 45 |
| With 10 rules and not 5 | 40 | 48 |

The last column uses the searches of the run just before each change, so its rows are not the same searches.

Add a question when Help gives a wrong answer about behavior. Read the code, write the key points from it, and cite the lines. A change to the code can make a reference answer wrong: when a run marks a good answer as wrong, check the reference first.

Each question lists as key points only what an answer must state to answer the question that was asked. Other true facts are in `extra_points`, which nothing grades.

Runs on 2026-11-10, with the first 50 questions (the first five areas). Each question was asked one time.

| Run | Graded | Correct | Wrong | Key points stated | The searches returned the right lines | Cost for a question | First word (median) |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| Haiku 5.5, low effort, with one list of 6 chunks of notes and code | 50 | 26 | 7 | 69% | 27 | $0.0018 | 4.5 s |
| Haiku 5.5, low effort, with the right source in the prompt (`oracle`) | 49 | 40 | 2 | 88% | | | |
| Haiku 5.5, low effort, with code and notes as two groups and the search rules in the prompt | 50 | 34 | 1 | 83% | 45 | $0.0029 | 5.5 s |
| Haiku 5.5, medium effort, the same search and prompt | 50 | 38 | 3 | 89% | 46 | $0.0035 | 6.4 s |
| Sonnet 5.5, low effort, the same search and prompt | 49 | 38 | 2 | 89% | 45 | $0.039 | 5.7 s |
| Haiku 5.5, low effort, with the rules: a search returns rules, then code, then notes | 49 | 37 | 0 | 84% | 47 | $0.0024 | 4.4 s |
| The same, with each question asked as the person that it names | 50 | 37 | 4 | 86% | 46 | $0.0024 | 4.6 s |
| With the search phrases of the rules, and 10 rules in a search | 50 | 38 | 1 | 89% | 49 | $0.0024 | 4.8 s |

What the runs show:

- In the first run, two chunks in three that a search returned were design notes, and almost every wrong answer came from a question where the search did not return the code. Code first, with the search rules, removed most of the wrong answers.
- The last three runs do not differ by more than chance. Sonnet costs 13 times as much as Haiku for the same result, so the model is not what limits the answers now.
- Before the rules, what was left was mostly a search that can't find the code because the code does not use the user's words. One question failed in each of those three runs: the code that approves an expense at submit, when the owner is also the approver, has none of the words "self", "own" or "automatic".
- With the rules, no answer was wrong, and the model needs fewer searches (2.2 for a question and not 2.6), so an answer starts about a second sooner and costs less. The question about your own expense now gets the right answer: the third search found the rule. 46 of the 50 answers cite a rule.
- The questions that are still partial are mostly not a search problem: the search returned the right lines for 10 of the 12. The model leaves out a fact, or says that it could not confirm it.
- Those runs asked every question as a person with no claims. A question that starts with "I'm a committer" then conflicted with the prompt, and some answers said that the user can't do the task. Each question now names the claims of the person who asks, in `held`.
- The last two rows are the same system, and one has 0 wrong answers and the other has 4. One run of 50 questions can't tell them apart: three of the four wrong answers are for questions that other runs got right. Do not read a change of a few answers as a result. Of the 18 questions that name a role, 15 were correct when asked as that role and 13 before.
- In three of those four wrong answers, the search did not return the rule that settles the question, though the rule is there. The words differ: the question says "my own manager" and the rule says "the owner is also the approver".
- With the search phrases and 10 rules in a search, those four answers are right or partly right, and the searches returned the right lines for 49 of the 50 questions. The one answer that the judge graded as wrong states every key point: the judge took a true statement about the approver list for an error.

## Rules that must hold

1. A visitor who is not signed in gets nothing from the chat. Every route under `/api/help` uses `apis.RequireAuth("users")`. The corpus is embedded from `app/helpcorpus/data`, which is not in `pb_public`, and no route serves it.
2. The model gets only text that the user may read. Every search and read of the corpus takes the user's claims.
3. The server decides who the user is. The user id and the claims come from the auth token and the database. The browser sends only the question, the chat id, the route and the record on the page.
4. A user reads only their own chats. The `chat_turns` rules and every query filter on the user id.
5. The model never writes to the database.
6. The model gets no data of a user except that user's own settings that `facts.go` lists, and never a value of another user. Most are values that the user can see in Turbo. The values of the user's own admin profile are the exception, by the owner's decision. Ask the owner before you add another value that the user cannot see.

## What the model may read

The repo is public, so the user docs, the design notes and the code are open to every signed-in user. `app/helpcorpus/audience.go` holds the rules. A rule can restrict a folder or a word to some claims later.

Uncommit is not a secret, but the chat does not raise it. The system rules tell the model to treat a commit as final, so that people do not come to rely on a reversal. A holder of `commit` or `admin` gets an answer when they ask about it. This is guidance for the answer and not access control.

When a user cannot do a task, the answer says who can, in the glossary's words and not in claim names.

## A turn

`POST /api/help/ask` runs one turn and returns server-sent events.

1. Check the setting, the question and the limits. Load the user's claims.
2. Write a `chat_turns` row with status `running`, in the same transaction as the limit checks.
3. Build the prompt. The first part is the system rules and the handbook (all the user docs). It is the same bytes for every user, so the model caches it. The second part is what the database holds for every user alike: all the claims in Turbo with their descriptions, and Turbo's settings. The model caches it too, with its own cache mark. The third part names the user's claims, the page they are on and the user's own settings. See "What the model knows about the user" below.
4. Send the earlier turns of the chat and the question. The server reads the earlier turns from `chat_turns`.
5. Run the tools that the model calls: `search_turbo`, `read_source` and `report_gap`. A reply that has text for the user and only a gap report ends the turn there: a second request made the model say the same twice. One search returns up to 10 rules of the code, then up to 3 rules of the collections, then up to 6 chunks of code, then up to 2 chunks of design notes, each marked with its type. They are separate groups because the notes have headings in plain words: in one list they outscore the code, and the code is what says how Turbo behaves today. The prompt tells the model to search the code for a question about a rule or a case, to trust the code over a note, and to say when it could not confirm a rule.
6. Stream the answer, then the sources that it cites.
7. Write the result to the row.

Events: `meta` (chat id and turn id), `status`, `text`, `sources`, `done`.

`report_gap` has three kinds. `not_covered` and `off_topic` give the turn the status `gap`, and the next turn of the chat does not get it as history. `no_steps` is for an answer that gives the rule of a task when no page of the user docs has its steps: the turn stays `answered`, and the `gap` field of the row names the task. Read those rows to find the pages that the docs lack.

### What the model knows about the user

The corpus says how a limit applies. It cannot say what one user's limit is, or what an admin has set a switch to, because those values are in the database. `facts.go` reads them for each question and writes two sections.

"This user's settings" is in the part of the prompt for one user. It has:

- From the profile: the manager, the alternate manager and the default division.
- From the admin profile: the default branch, and if Personal Reimbursement expenses are allowed.
- From the admin profile, also for every user: **Salary**, **Work Week Hours**, **Time Sheet Expected**, **Off Rotation Permitted**, **Untracked Time Off**, **Skip Min Time Check**, the time off opening date and balances, and **Personal Vehicle Insurance Expiry**. Only an admin, HR and a time off manager can open admin profiles in Turbo. The owner decided that Help can tell each user these values of their own admin profile, because they explain what Turbo asks of that user. The payroll ID and the charge-out rate stay out.
- For a holder of `po_approver`: the limit for each kind of PO, and the divisions.
- For a holder of `project_authorization`: the largest project value.

"Turbo's settings today" is the same for every user, so it is in the common part, after the list of all the claims:

- An editing switch of `app_config`, only when it is off.
- If a new project needs a project request.
- The amount at which an expense needs a PO, the amount that the expenses of a PO can go over its total, and the hours of the second approval timeout.
- The amount above which each kind of PO needs a second approval.

Each line of Turbo's settings has the name of the setting in the code, so the model can search for what the setting does.

The user can already see almost every one of these values: their profile, the defaults that `/api/user_defaults` returns, their own row of `user_po_approver_profile`, and the `app_config` and `expenditure_kinds` collections. The admin profile values above are the one exception. A project approver's limit is shown to everyone who picks an approver. The names of the user's two managers go to the model provider with the question.

The common part is about 1,200 tokens. It has its own cache mark after the shared part, so a call reads it from the cache at a tenth of the price. When an admin changes a claim or a setting, the next call reads the shared part from the cache and writes the common part again. Its lines are in a fixed order, so it is the same bytes for every user.

The part for one user is about 150 tokens for a user with no limits and about 270 for an approver. It is after the cache marks, so each model call of a turn pays for it as plain input. Before this layout, that part also held the list of all the claims, about 550 tokens, so a call now sends fewer tokens at the full price than it did. The benchmark sends the claims and neither section of settings, because its questions have no user.

A part that cannot be read is left out and logged, and the chat answers without it.

The system rules tell the model to answer with the values that the question is about and no others, to check the user's claims when the handbook says who can do a task, to treat the text of a claim as a short note and not as all that the claim allows, to give steps only from the handbook or the code of a screen, and to leave an earlier answer alone unless a source in the turn shows that it was wrong. A real chat showed each of these mistakes: the model took back a right answer because the claim `admin` had the text "Can delete jobs". In the first chat with the settings, the question was about POs and the answer also gave the limit for project requests, so each limit line now says what it covers. The line of PO limits also says that a limit is for the final approval: an approver can give the first approval of a larger PO.

The answer cites a docs page as a link, such as `[Time](/docs/time)`, and a chunk from a tool as `[S3]`. The server keeps a source only if it knows it.

Other routes: `POST /api/help/turns/{id}/feedback` saves a rating, a note, a source that the user opened, or a copy. `POST /api/help/warm` asks the model to cache the shared prompt and the common part.

## The log

`chat_turns` holds one row for each turn: the question and answer, the sources, the model, effort, build and corpus hash, the user's claims, the tool calls with the ids of what they returned, the route and record, the status, token counts, times, and the user's feedback.

The log holds chunk ids and not chunk text. The corpus hash says which text an id means, and the generator gives the same corpus for the same commit.

A user can list and view their own rows through the collection API. Nobody can create, change or delete a row there. Feedback goes through its own route, which changes the feedback fields only.

## The model clients

The chat talks to a model through the `Model` interface in `model.go`. The tests use a stub.

- `anthropic.go` calls the Claude Messages API with plain HTTP. It marks the shared prompt and the common part for Claude's cache, and it sends the model's thinking blocks back unchanged between tool calls. It does not use Anthropic's Go SDK: the compiler needs about 2.4 GB of memory for that one package, which a small build machine does not have, and the chat uses one endpoint. The request and the stream are written for API version `2026-10-21`.
- `openai.go` calls any API in the OpenAI chat completions format, with plain HTTP. It puts the shared prompt first and the common part after it, because these services cache the start of a prompt by themselves.
- `providers.go` lists the providers: the secret for each key, the address, and the models and effort values that the setting may name.

Both clients send the same system rules, handbook and tools. Neither sets sampling values or a forced tool choice.
