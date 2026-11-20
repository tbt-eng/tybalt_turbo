<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/computers.ts","sha256":"656c0bb5f9558c26d35a88f328cec63e1a09ed61a32944731862b8d85b585676","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/computers.ts

## The server checks each computer action again

Lines 1-3

Asked as: page allows but server rejects; action fails after clicking; server ignores button; browser check vs server check

The computers pages apply the action rules in the browser. The server checks each action again (app/routes/computers.go), so a page that allows an action does not guarantee that the server accepts it.

## Last user is the Turbo user the last sign-in matched

Lines 27-28

Asked as: who last signed in; last user empty; sign-in matched no user; computer last signed in by

A computer's last user (last_user) is the Turbo user that the most recent sign-in matched. It is empty when the sign-in matched no Turbo user.

## Who can see the computers pages and actions

Lines 39-42

Asked as: cannot see computers page; no access to computers; missing computers claim; who can assign computers

Staff with the computers claim or the admin claim can open the computers pages, the list routes and the actions. Other staff cannot (canViewComputers). The server checks each action again.

## A computer's assignment state

Lines 44-63

Asked as: computer assigned to someone else; unassigned computer status; assignment state meaning; assigned but not last user

A computer is unassigned when nobody is assigned to it. It is in the last-user state when it is assigned to its last user. It is in the other state when it is assigned to someone who is not the last user (computerAssignment).

## When Assign to the last user is allowed

Lines 48-67

Asked as: assign button greyed out; cannot assign to last user; assign to last user disabled; retired computer assign

Assign sets the assignee to the last user. It is allowed only if the computer is not already assigned to the last user, a sign-in matched a last user, and the computer is not retired (canAssign). A computer with no assignee, or an assignee who is another person, can be assigned.

## Last user name fallback

Lines 70-80

Asked as: shows Unknown as last user; last user name missing; last user shows wrong name; name not showing for last user

Turbo shows the last user as the matched Turbo user's name. If no user matched, it shows the given name and surname that the sign-in reported. If both are empty, it shows "Unknown".

## Maker and model display

Lines 82-90

Asked as: maker repeated in model; HP HP EliteBook; duplicate brand in model name; model name shows twice

Turbo shows the maker and the model together. If the model starts with the maker's first word, compared without case (such as "HP EliteBook" from "HP Inc."), it shows the model alone.

## Size formatting

Lines 92-105

Asked as: memory shows wrong size; disk size in GB wrong; bytes shown as number; size displays blank

Turbo shows sizes in powers of 1024 with the units B, KB, MB, GB, TB and PB, rounded to one decimal place (17179869184 is "16 GB"). A size that is zero, negative or not a finite number shows as empty.

## Boot drive display

Lines 107-125

Asked as: C drive shows no free space; boot disk display missing; free percent wrong; drive letter and size missing

Turbo shows the boot drive as its letter, its file system in parentheses, its size, and its free space with the percent free rounded down. With no capacity, only the letter and file system show. The free text is left out when free space is empty.

## System type names

Lines 127-141

Asked as: system type shows blank; computer type number; what does type 3 mean; desktop or laptop code

Turbo names the system type from the codes 1 to 8: Desktop, Mobile, Workstation, Enterprise Server, SOHO Server, Appliance PC, Performance Server and Maximum. Any other code shows as empty.

## Network configuration display

Lines 153-194

Asked as: network adapters missing; MAC address order; ip address list wrong; adapter settings not showing

Turbo lists one adapter for each MAC address, sorted by MAC. Known properties come first with readable labels (DHCP, DHCP server, IP addresses, Subnets, Gateways, DNS host name, DNS servers). Other properties keep their names and follow in name order. Booleans show as Yes or No, lists are joined by commas, and empty values are left out. A config that is not an object gives no adapters.

## Dates and times use the viewer's device time

Lines 196-211

Asked as: wrong time on computer; timezone off by hours; last seen time wrong; date shows wrong day

Turbo shows the day and the time of day of a computer date in the time zone of the viewer's device. An empty or invalid date shows as empty.

## What the computers list search matches

Lines 213-228

Asked as: search by email not working; cannot find computer by serial; search misses computers; search by user name

The list search matches the computer name, serial, maker, model, the reported given name, surname, email and UPN, the last user's name, the assignee's name and the OS version.

## Raw login report display

Lines 246-319

Asked as: login report shows odd text; raw login json; login details display; sign-in report blank

The raw login report shows as text. A string shows as it is, and other values show as indented JSON. A missing report shows as empty.

## List page sizes and defaults

Lines 266-297

Asked as: page size options; list shows wrong count per page; page number invalid in URL; search text spaces trimmed

The computers list offers page sizes 10, 20 and 50. Any other perPage in the URL gets 20. A page that is not a positive whole number gets page 1. Turbo trims the search text. The URL leaves out default values.

## Sign-in reported name and matched user name

Lines 299-312

Asked as: sign-in name differs from user; Unknown name on sign-in; reported name not matching; which name is shown for sign-in

A sign-in shows the reported given name and surname, or "Unknown" if both are empty. Turbo shows the matched Turbo user's name only when a user matched and the name differs from the reported name, ignoring case and spaces at the ends.
