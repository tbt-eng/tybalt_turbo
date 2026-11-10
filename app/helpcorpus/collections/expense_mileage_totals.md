<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: expense_mileage_totals

`expense_mileage_totals` is a view: a saved query that Turbo reads and never writes. Turbo's screens reach it through PocketBase's API (`/api/collections/expense_mileage_totals/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): expense_mileage_totals

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can view one record (view rule): expense_mileage_totals

No rule. Nobody can do this through the API, except a PocketBase superuser.

## The query of the view: expense_mileage_totals

```sql
-- 1) Explode + window up front
WITH rates_expanded AS (
  SELECT
    m.effective_date,
    CAST(t.key   AS INTEGER) AS tier_lower,
    LEAD(CAST(t.key AS INTEGER)) OVER (
      PARTITION BY m.effective_date
      ORDER BY CAST(t.key AS INTEGER)
    ) AS tier_upper,
    CAST(t.value AS REAL)    AS tier_rate

  FROM expense_rates m
  CROSS JOIN json_each(m.mileage) AS t
),

-- 2) Base: compute cumulative & interval boundaries
base AS (
  SELECT
    e.id,
    e.uid,
    e.date,
    -- This ends up being faster than joining to mileage_reset_dates
    (
      SELECT MAX(r.date)
      FROM mileage_reset_dates r
      WHERE r.date <= e.date
    ) AS reset_mileage_date,
    e.distance,
    -- interval end = cumulative distance
    SUM(e.distance) OVER (
      PARTITION BY e.uid, (
        SELECT MAX(r.date)
        FROM mileage_reset_dates r
        WHERE r.date <= e.date
      )
      ORDER BY e.date
    ) AS end_distance,
    (
      SELECT MAX(m.effective_date)
      FROM expense_rates m
      WHERE m.effective_date <= e.date
    ) AS effective_date
  FROM expenses e
  WHERE e.payment_type = 'Mileage'
),

/* 
3) Join each expense to its tiers, filtering to only those that overlap We are
pairing each expense’s [start_distance, end_distance) interval with only those
tier intervals [tier_lower, tier_upper) that intersect it. By exploding each
tier into its own row, an expense whose kilometres cross a tier boundary will
join to two (or more, in theory) tier‐rows. Thus we can expect that the number
of rows in overlaps can be more than the number of rows in expenses.
*/
overlaps AS (
  SELECT
    b.id,
    b.end_distance - b.distance AS start_distance,
    b.end_distance,
    r.tier_lower,
    COALESCE(r.tier_upper, 1e9) AS tier_upper,
    r.tier_rate
  FROM base b
  JOIN rates_expanded r
    ON r.effective_date = b.effective_date
  WHERE b.end_distance > r.tier_lower
    AND (r.tier_upper IS NULL OR (b.end_distance - b.distance) < r.tier_upper)
),

/*
4) Compute overlap length per tier
for each expense × tier row, compute how many kilometres of the expense fall into that tier by:
  1. Determining the overlap interval between the expense’s [start_distance, end_distance)
     and the tier’s [tier_lower, tier_upper) interval:
       • overlap_start = max(start_distance, tier_lower)
       • overlap_end   = tier_upper IS NULL
                           ? end_distance
                           : min(end_distance, tier_upper)
  2. Calculating overlap_km = max(0, overlap_end − overlap_start)
       • Ensures negative values (no overlap) are clipped to zero
  3. Carrying along tier_rate so we can later multiply overlap_km × tier_rate
*/
tier_calcs AS (
  SELECT
    id,
    max(0,
      min(end_distance, tier_upper)
      - max(start_distance, tier_lower)
    ) AS overlap_km,
    tier_rate
  FROM overlaps
)

-- 5) Sum up reimbursement per expense
SELECT
  b.id,
  b.uid,
  b.date,
  b.reset_mileage_date,
  b.distance,
  b.end_distance AS cumulative,
  b.effective_date,
  ROUND(COALESCE(
    -- sum up this expense’s (overlap × rate) directly
    (SELECT SUM(overlap_km * tier_rate)
     FROM tier_calcs tc
     WHERE tc.id = b.id),
    0
  ), 2) AS mileage_total
FROM base b
ORDER BY b.date;
```

## Fields: expense_mileage_totals

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; pattern "^[a-z0-9]+$" |
| `uid` | json | maxSize 1 |
| `date` | json | maxSize 1 |
| `reset_mileage_date` | json | maxSize 1 |
| `distance` | json | maxSize 1 |
| `cumulative` | json | maxSize 1 |
| `effective_date` | json | maxSize 1 |
| `mileage_total` | json | maxSize 1 |
