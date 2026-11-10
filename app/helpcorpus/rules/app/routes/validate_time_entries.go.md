<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/validate_time_entries.go","sha256":"ec63c221f73e601e2c96893e14eb208d2023f719ea12cfbfa8f6767541860eba","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/validate_time_entries.go

## Time entry validation fails if time types cannot be loaded

Lines 13-21

Asked as: time type not loading; time sheet error_expanding_time_type_relations; time type missing on entries; cannot load time type

Turbo checks a time sheet's entries as a group. If it cannot load the time type of the entries (time_type), it refuses the time sheet with the error code error_expanding_time_type_relations.

## Time sheet week comes from the first entry

Lines 23-24

Asked as: wrong week ending date; week ending taken from first entry; which week date is used; time sheet week date wrong

Turbo reads the week ending date (week_ending) from the first entry only. It uses that date for the opening date and balance checks on the whole time sheet.

## Admin profile settings drive time sheet checks

Lines 26-33

Asked as: admin profile settings; check against my profile; profile balances used in checks; staff profile limits time sheet

Turbo checks a time sheet against the staff member's admin profile: salary, opening date, opening OP and OV balances (opening_op, opening_ov), off-rotation permission, untracked time off, skip minimum time check (skip_min_time_check) and work week hours (work_week_hours).

## A work record can appear on only one entry

Lines 47-58

Asked as: duplicate work record; same work record twice; work record in multiple entries; work record repeated on sheet

Turbo refuses a time sheet when the same work record (work_record) appears on more than one entry. The message is "work record <number> appears in multiple entries". Entries with an empty work record are not checked.

## Only one off-rotation (OR) entry per date

Lines 66-76

Asked as: two OR entries same day; double off rotation same date; more than one OR on a date; two off rotation days same day

Turbo refuses a time sheet that has two or more off-rotation (OR) entries on the same date. The message is "more than one OR entry exists for the date: <date>".

## Salaried staff cannot claim a full week off (OW)

Lines 77-84

Asked as: salaried cannot take full week off; OW not allowed for salary; full week off rejected; salaried week off error

Turbo refuses a time sheet with a full-week-off (OW) entry when the staff member is salaried. The message is "salaried staff cannot claim full week off, use OP or OV".

## An OW entry must be the only entry on the time sheet

Lines 85-102

Asked as: OW with other entries; full week off alongside other entries; OW must be alone; two week off entries

Turbo refuses a time sheet when a full-week-off (OW) entry sits next to any other entry. The message is "if present, an OW entry must be the only entry on a time sheet". Only one OW entry is allowed on a time sheet, with the message "only one off-rotation week entry can exist on a time sheet".

## Salaried staff cannot request overtime payout (OTO)

Lines 103-111

Asked as: salaried overtime payout; salaried request overtime pay; OTO refused salaried; overtime payout not allowed

Turbo refuses an overtime payout request (OTO) entry from salaried staff. The message is "salaried staff cannot request overtime payouts".

## Only one overtime payout request per time sheet

Lines 112-120

Asked as: two overtime payout requests; multiple OTO entries; more than one payout request; second overtime payout on sheet

Turbo refuses a time sheet with more than one payout request (OTO) entry. The message is "only one payout request entry can exist on a time sheet".

## Salaried staff cannot bank overtime (RB)

Lines 121-129

Asked as: salaried bank overtime; salaried overtime to bank refused; RB not allowed salary; cannot add overtime to bank

Turbo refuses an add-overtime-to-bank (RB) entry from salaried staff. The message is "salaried staff cannot bank overtime".

## Only one overtime banking entry per time sheet

Lines 130-139

Asked as: two overtime banking entries; multiple RB entries; more than one bank entry; second overtime bank on sheet

Turbo refuses a time sheet with more than one overtime banking (RB) entry. The message is "only one overtime banking entry can exist on a time sheet". The hours of the RB entry count as banked hours.

## Regular entries count as job or non-job hours

Lines 140-147

Asked as: regular hours counted as job; where do regular hours go; non-job hours from R entry; RT hours totals

For regular (R) and regular time (RT) entries, Turbo adds the hours to job hours when the entry has a job. Without a job, it adds them to non-job hours.

## Other time types need non-zero hours

Lines 148-161

Asked as: zero hours error; time entry missing hours; entry with no hours refused; hours required for time type

For time types other than OR, OW, OTO, RB, R and RT, Turbo refuses an entry with zero hours. The message is "a time entry is missing hours". Turbo totals the hours of these entries per time type code.

## Banked hours cannot drop worked hours below 44

Lines 166-173

Asked as: banked hours too many; worked hours below 44; bank overtime drops total; 44 hours minimum banking

When an entry banks hours, Turbo refuses the time sheet if job hours plus non-job hours minus banked hours is under 44. Exactly 44 passes. The message is "banked hours cannot bring your total worked hours below 44 hours on a time sheet".

## OV and OP cannot raise hours above the work week

Lines 181-200

Asked as: vacation exceeds work week; OV over work week hours; PPTO too many hours in week; hours above weekly limit

When a time sheet claims vacation (OV) or PPTO (OP) hours, Turbo refuses it if total hours (job, non-job and all other time-off) exceed the staff member's work week hours (work_week_hours). Equal passes. The message is "you cannot claim OV or OP entries that increase hours beyond <hours>".

## Salaried staff need permission for OR entries

Lines 202-209

Asked as: salaried off rotation not allowed; OR entry salaried permission; no permission for off rotation; salaried OR refused

Turbo refuses a time sheet with off-rotation (OR) entries from salaried staff whose profile does not permit off rotation (off_rotation_permitted). The message is "salaried staff need permission to claim OR entries".

## Salaried staff must meet minimum hours

Lines 211-222

Asked as: not enough hours salaried; minimum hours on sheet; salaried under work week hours; minimum time check failed

Turbo refuses a salaried staff member's time sheet when hours total under the work week hours (work_week_hours). Each OR date counts as 8 hours. The check applies only when skip_min_time_check is "no" and untracked time off is off. The message is "you must have a minimum of <hours> hours on your time sheet".

## Salaried staff cannot claim sick time (OS)

Lines 224-231

Asked as: salaried sick time; OS refused salaried; sick leave salaried not allowed; claim sick day salary

Turbo refuses a time sheet with a sick time (OS) entry from salaried staff. The message is "salaried staff cannot claim OS. Please use OP or OV instead".

## Untracked time off limits salaried staff to R and RT

Lines 233-244

Asked as: untracked time off error; untracked staff OB OH; only R and RT allowed; untracked staff vacation refused

Turbo refuses OB, OH, OP or OV entries from salaried staff with untracked time off (untracked_time_off). The message is "staff with untracked time off are only permitted to create R or RT entries".

## Invalid opening date on profile

Lines 246-270

Asked as: bad opening date; opening date invalid profile; admin profile date wrong; opening date format error

Turbo refuses the time sheet when the opening_date on the admin profile is not a valid date (YYYY-MM-DD). The message is "your admin_profile has an invalid opening_date, contact support". The opening date does not have to be a Sunday, because that check is disabled.

## Invalid week ending date

Lines 272-280

Asked as: invalid week ending; week ending format wrong; entry date not valid; week ending error contact support

Turbo refuses the time sheet when an entry's week_ending is not a valid date (YYYY-MM-DD). The message is "an entry has an invalid week_ending date, contact support".

## Time off on a sheet before the opening date

Lines 282-292

Asked as: claim before opening balances; time sheet before opening date; prior period vacation claim; balances effective later date

When a time sheet claims OV or OP and the opening date is after the week ending, Turbo refuses it. A week ending equal to the opening date passes. The message is "your opening balances were set effective <date> but you are submitting a time sheet for a prior period, contact support".

## Opening balances out of date for a new period

Lines 294-309

Asked as: new year opening balances; contact accounting opening balances; payroll year end balance update; balances out of date

When a time sheet claims OV or OP and the payroll year end date is after the profile's opening date, Turbo refuses it until accounting updates the opening balances. The message says to contact accounting to update opening balances for the new period. Time sheets without OV or OP skip this check.

## Used OV and OP hours are summed after the opening date

Lines 311-367

Asked as: used vacation after opening; sum of used OV and OP; saved vacation hours total; previously claimed vacation hours

Turbo sums a staff member's saved OV and OP hours with week ending after the opening date (not equal to it) and up to and including this week ending. Entries in the opening date's own week are already in the opening balance. The sum is 0 when there are none.

## Vacation claim cannot exceed the balance

Lines 369-378

Asked as: vacation claim too big; not enough vacation balance; OV exceeds balance; vacation balance insufficient

Turbo refuses a time sheet with an OV entry when used vacation hours exceed the opening vacation balance (opening_ov). A time sheet with no OV entry skips the check, even if the balance is negative. The message is "your vacation claim exceeds your available vacation balance".

## PPTO claim cannot exceed the balance

Lines 380-387

Asked as: PPTO too big; not enough PPTO balance; OP exceeds balance; PPTO balance insufficient

Turbo refuses a time sheet with an OP entry when used PPTO hours exceed the opening PPTO balance (opening_op). A time sheet with no OP entry skips the check. The message is "your PPTO claim exceeds your available PPTO balance".

## Vacation must be used up before PPTO

Lines 389-396

Asked as: use vacation before PPTO; vacation still left PPTO refused; PPTO before vacation used up; exhaust vacation first

Turbo refuses a time sheet with an OP entry while opening vacation minus used vacation is above 0. The message is "exhaust your vacation balance (<hours> hours) prior to claiming PPTO".
