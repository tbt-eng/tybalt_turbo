<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->

# Collection: expense_allowance_totals

`expense_allowance_totals` is a view: a saved query that Turbo reads and never writes. Turbo's screens reach it through PocketBase's API (`/api/collections/expense_allowance_totals/records`), and the rules below decide who may do each action there. A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.

## Who can list the records (list rule): expense_allowance_totals

No rule. Nobody can do this through the API, except a PocketBase superuser.

## Who can view one record (view rule): expense_allowance_totals

No rule. Nobody can do this through the API, except a PocketBase superuser.

## The query of the view: expense_allowance_totals

```sql
SELECT e.id, 
  e.date, 
  r.effective_date allowance_rates_effective_date, 
  e.payment_type,
  e.allowance_types,
  r.breakfast breakfast_rate,
  r.lunch lunch_rate,
  r.dinner dinner_rate,
  r.lodging lodging_rate,
  r.mileage,
  ((CASE WHEN e.allowance_types LIKE '%"Breakfast"%'   THEN r.breakfast ELSE 0 END)
  + (CASE WHEN e.allowance_types LIKE '%"Lunch"%'     THEN r.lunch     ELSE 0 END)
  + (CASE WHEN e.allowance_types LIKE '%"Dinner"%'    THEN r.dinner    ELSE 0 END)
  + (CASE WHEN e.allowance_types LIKE '%"Lodging"%'   THEN r.lodging   ELSE 0 END)
  )AS allowance_total,
  RTRIM(
    (CASE WHEN e.allowance_types LIKE '%"Breakfast"%' THEN 'Breakfast ' ELSE '' END) ||
    (CASE WHEN e.allowance_types LIKE '%"Lunch"%'     THEN 'Lunch '     ELSE '' END) ||
    (CASE WHEN e.allowance_types LIKE '%"Dinner"%'    THEN 'Dinner '    ELSE '' END) ||
    (CASE WHEN e.allowance_types LIKE '%"Lodging"%'   THEN 'Lodging '   ELSE '' END)
  ) AS allowance_description
FROM expenses e 
LEFT JOIN expense_rates r ON ((r.effective_date = (SELECT MAX(i.effective_date) FROM expense_rates i WHERE (i.effective_date <= e.date))))
WHERE e.payment_type IN ('Allowance','Meals');
```

## Fields: expense_allowance_totals

| Field | Type | Details |
| --- | --- | --- |
| `id` | text | required; pattern "^[a-z0-9]+$" |
| `date` | text | required; pattern "^\\d{4}-\\d{2}-\\d{2}$" |
| `allowance_rates_effective_date` | text | required; pattern "^\\d{4}-\\d{2}-\\d{2}$" |
| `payment_type` | select | required; values ["OnAccount","Expense","CorporateCreditCard","Allowance","FuelCard","Mileage","PersonalReimbursement"]; maxSelect 1 |
| `allowance_types` | select | values ["Lodging","Breakfast","Lunch","Dinner"]; maxSelect 4 |
| `breakfast_rate` | number | required |
| `lunch_rate` | number | required |
| `dinner_rate` | number | required |
| `lodging_rate` | number | required |
| `mileage` | json | required; maxSize 2000000 |
| `allowance_total` | json | maxSize 1 |
| `allowance_description` | json | maxSize 1 |
