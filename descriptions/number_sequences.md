# Number Sequences

This note defines where Turbo should use the `number_sequences` table and where
it should not.

## Purpose

`number_sequences` exists to serialize number assignment for "hot" numbering
streams where concurrent creates are plausible and duplicate-number races are a
real operational risk.

It is not intended to mirror every parent/child numbering relationship in the
system.

## Rule of Thumb

Use `number_sequences` only for shared top-level or parent-style sequences with
a bounded number of scopes.

Do not use `number_sequences` for child numbering where the scope count grows
with individual parent records unless production evidence shows that collisions
are a real problem.

Internal scope keys should follow:

- `<collection>:<scope_kind>:<display_prefix>`

Current examples:

- `work_records:monthly:W2604-`
- `purchase_orders:monthly_parent:2604-`

## Why

### Good fit

- monthly work record numbers such as `WYYMM-NNNN`
- monthly parent purchase order numbers such as `YYMM-NNNN`
- yearly top-level job numbers such as `YY-NNNN`
- yearly top-level proposal numbers such as `PYY-NNNN`

These are good fits because:

- concurrent creation is plausible
- the number of sequence scopes stays small
- the sequence table acts as a true coordination primitive rather than as a
  shadow copy of business data

### Bad fit

- child purchase order suffixes such as `YYMM-NNNN-XX`
- child job suffixes such as `YY-NNNN-XX`

These are poor fits because:

- contention is expected to be rare
- each parent would need its own sequence row
- the number of scopes would grow with ordinary data volume
- the table would start to duplicate state that can already be derived from the
  child records themselves

## Seeding and Manual Inserts

Each scope seeds itself once from existing records the first time that scope is
used, and after that the counter row becomes the source of truth for the next
auto-generated value.

That means post-seed manual inserts inside the auto-generated range are not
noticed by the sequence row. If manual numbers are later inserted into that
same range, the next auto-generated value may eventually collide with them.

Operationally, that means:

- manual/imported numbers should stay outside the auto-generated range whenever
  possible
- parent purchase orders already do this by reserving `5000+` for
  manual/imported values
- work records currently assume there will not be post-seed manual inserts into
  the normal monthly auto-generated range

Future consideration:

- if work records or another sequence family later need in-range manual inserts,
  add an explicit reseed/bump tool or a reserved manual range before relying on
  that workflow

## Current Direction

- `work_records` should use `number_sequences`
- parent `purchase_orders` should use `number_sequences`
- child `purchase_orders` should continue using the existing scan-based approach
- top-level `jobs` may adopt `number_sequences` later if we choose to harden job
  numbering against concurrency
- child `jobs` should continue using the existing derived-from-children approach
  unless operational evidence justifies something stronger

## Revisit Criteria

We should only extend `number_sequences` to child numbering if we see one of the
following in production:

- repeated unique-index collisions on child numbers
- meaningful user-facing failures caused by concurrent child creation
- a workflow change that makes concurrent child creation common rather than rare

Until then, keeping child numbering derived from the existing child rows is the
preferred tradeoff.
