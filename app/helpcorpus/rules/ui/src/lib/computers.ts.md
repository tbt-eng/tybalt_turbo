<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/computers.ts","sha256":"6dd8c55f45bc2771710644a1ea72f5e599ca78527e18d5ebaede89581ef87307","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/computers.ts

## Server checks each computer action again

Lines 1-3

Asked as: button shows but action denied; hidden button still works; server rechecks computer action; button visible no permission

The Computers pages decide in the browser which actions to offer. The server checks each action again (app/routes/computers.go), so a hidden or shown button does not grant or deny access by itself.

## Last user is empty when no Turbo user matches

Lines 27-29

Asked as: last user blank; no user for sign-in; computer last user empty; sign in matched nobody

A computer's last user (last_user) is the Turbo user that the last sign-in matched. It is an empty value when the sign-in matched no Turbo user.

## Who can see the Computers pages and use their actions

Lines 39-42

Asked as: can't see computers page; no access to computers; computers menu missing; who may use computers

A person with the computers claim or the admin claim can see the Computers pages, the list routes and the actions (canViewComputers). Other people cannot. The server checks each action again.

## Assignment state of a computer

Lines 44-63

Asked as: computer assignment status; unassigned or assigned meaning; computer assigned to whom; assignment label computer

A computer with no assignee is 'unassigned'. A computer assigned to its last user is 'last_user'. A computer assigned to anyone else is 'other' (computerAssignment).

## When Assign to the last user is available

Lines 48-67

Asked as: assign to last user missing; can't assign computer; assign button greyed out; give computer to last user

Assign sets the assignee to the last user. It is available only when the computer is not already assigned to the last user, a sign-in matched a Turbo user as the last user, and the computer is not retired. An unassigned computer and a computer assigned to someone else can both be assigned.

## How the last user is named

Lines 70-80

Asked as: last user shows unknown; last user name wrong; who is last user name; reported name for computer

Turbo shows the name of the matched Turbo user. If no user matched, it shows the given name and surname that the computer reported. If both are empty, it shows 'Unknown' (lastUserLabel).

## How maker and model are shown together

Lines 82-90

Asked as: maker and model repeated; HP model name shown twice; model text duplicates maker; computer model display

Turbo shows the maker and model together. When the model starts with the first word of the maker, case ignored, it shows the model alone, such as 'HP EliteBook' for maker 'HP Inc.' (modelLabel).

## How sizes are formatted

Lines 92-105

Asked as: disk size in GB wrong; bytes shown as gigabytes; memory size display; why size shows blank

Turbo formats sizes in powers of 1024 with units B, KB, MB, GB, TB and PB, as Windows does. 17179869184 bytes shows as '16 GB'. Turbo rounds values above bytes to one decimal. A size that is zero, negative or not a number shows as empty (formatBytes).

## Boot drive label with free space

Lines 107-125

Asked as: C drive free space; boot disk shows no size; drive label missing free space; system drive percent free

The boot drive shows its letter, its file system in parentheses, its size, and its free space with the percent free rounded down. If the size is empty or zero, Turbo shows only the letter and file system (bootDriveLabel).

## System type names

Lines 127-141

Asked as: system type blank; desktop or laptop type number; PC type shows nothing; what does system type mean

Turbo names the system type from the Win32 PCSystemType number: 1 Desktop, 2 Mobile, 3 Workstation, 4 Enterprise Server, 5 SOHO Server, 6 Appliance PC, 7 Performance Server, 8 Maximum. Any other number shows as empty.

## Network config display by adapter

Lines 153-194

Asked as: network adapters list order; MAC address adapters shown; IP config not showing; network settings missing fields

Turbo lists the network config with one adapter for each MAC address, sorted by MAC. Known properties come first with readable labels: DHCP, DHCP server, IP addresses, Subnets, Gateways, DNS host name, DNS servers. Other properties follow in name order and keep their own names. Fields with empty values are left out. A value that is not an object shows as a field named 'value'. A config that is not an object shows no adapters.

## Network values are shown as plain text

Lines 164-170

Asked as: network true false display; yes no network values; list values comma separated; network field shows json

Turbo shows true as 'Yes' and false as 'No'. It joins list values with commas and drops empty items. It shows null as empty and nested objects as JSON text (networkValue).

## Dates and times use the user's device

Lines 196-211

Asked as: wrong time shown; timezone of computer time; time differs from my clock; last seen time wrong

Turbo shows the day and time of a stored datetime in the time zone of the user's device, such as '8:00 AM'. An empty or invalid datetime shows as empty (localDay, localTime).

## What the computer list search matches

Lines 213-228

Asked as: computer search not finding; search by serial or email; which fields search covers; search computers by owner

The computer list search matches the computer name, serial, maker, model, reported given name and surname, email, UPN, last user name, assignee name and OS version.

## Rejected report is shown as text

Lines 246-301

Asked as: rejected sign-in report text; raw login json view; why report looks like code; failed login details display

Turbo shows the report of a rejected sign-in as indented JSON text, or as is if it is a string. It treats the report as untrusted. Only the single-report request fills it in, not the list (rawLoginReportText).

## Sign-in and rejected report list searches

Lines 257-279

Asked as: search rejected logins; sign-in search finds nothing; blank search returns all; search login results

The sign-in list search matches the computer name, serial, reported name and matched Turbo user name. The rejected report (raw login) search matches the computer name and the error. A blank search after trimming returns everything. Otherwise the server searches with the trimmed text (computerListSearch).

## Names shown for a sign-in

Lines 281-294

Asked as: sign-in name wrong; why name says unknown; matched user name differs; name shown at sign-in

A sign-in shows the name it reported, or 'Unknown' if it reported none. Turbo shows the matched Turbo user's name only when a user matched and the name differs from the reported name, case ignored (loginReportedName, loginTurboUserName).
