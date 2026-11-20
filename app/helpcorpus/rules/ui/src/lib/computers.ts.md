<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"ui/src/lib/computers.ts","sha256":"212ed4055130e13396f86e69abbe265d8ee7feefa561a49d4a8ebf57ec75b686","model":"claude-sonnet-5-5"} -->

# Rules in ui/src/lib/computers.ts

## The server checks each computer action again

Lines 1-3

Asked as: button shows but save fails; action allowed on screen but errors; server refuses computer change; page lets me but it fails

The pages show the computer rules, but the server checks each action again (app/routes/computers.go). The pages alone do not grant an action.

## Last user is the Turbo user matched at sign-in

Lines 27-29

Asked as: who last logged in; computer last signed in by; last user is blank; wrong person listed as user

A computer's last user (last_user) is the Turbo user that the last sign-in matched. It is empty when no Turbo user matched. In that case the list shows the name the computer reported (user_given_name, user_surname).

## Assignment fields on a computer

Lines 31-33

Asked as: who has this laptop; who is the computer given to; assigned time missing; computer owner shown

A computer record holds the person it is assigned to (assigned_to, assigned_to_name) and the time of assignment (assigned_time). An empty assigned_to means nobody has the computer.

## Retired computers carry a retired date

Lines 34-34

Asked as: no assign button for old computer; decommissioned laptop cannot assign; computer retired so no assign; retired machine options missing

Each computer record has a retired value (retired). An empty value means the computer is in use, and a value means it is retired. Turbo then does not offer Assign to last user for it.

## Who can see the computers pages and actions

Lines 39-42

Asked as: can't see computers menu; no access to computers; missing computers permission; need admin rights for computers

Staff with the computers claim or the admin claim can see the computers pages, the list routes and the actions. Staff without either claim cannot. The server checks each action again.

## Computer assignment states

Lines 44-63

Asked as: what unassigned status means; computer status with no owner; meaning of last-user state; assigned to someone else status

A computer is unassigned when nobody is assigned (assigned_to is empty). It is in the last-user state when the assignee is the last user. Otherwise it is assigned to someone other than the last user.

## When Assign to last user is available

Lines 48-67

Asked as: assign button missing; can't assign to last user; assign option greyed out; why no assign option

Assign sets the assignee to the last user. It is available only when the assignee is not already the last user, the last sign-in matched a Turbo user (last_user is not empty), and the computer is not retired.

## How the last user is named

Lines 70-80

Asked as: last user shows Unknown; wrong name for last user; name shown for user is wrong; no name shown for user

Turbo shows the name of the matched Turbo user (last_user_name). If there is none, it shows the given name and surname that the computer reported. If both are empty, it shows "Unknown".

## Maker and model display

Lines 82-90

Asked as: model name repeats maker; manufacturer shown twice; device name duplicated brand; HP appears twice in name

Turbo shows the maker and the model together. When the model starts with the first word of the maker, compared without regard to case (such as "HP EliteBook" with maker "HP Inc."), it shows the model alone. Empty parts are left out.

## Size units and rounding

Lines 92-105

Asked as: disk size shows wrong number; GB display looks wrong; memory size shown incorrectly; storage shows zero size

Turbo shows sizes in powers of 1024 with the units B, KB, MB, GB, TB and PB, so 17179869184 shows as "16 GB". Values above the byte level round to one decimal place. A size of zero, a negative size or an invalid number shows as empty.

## Boot drive display

Lines 107-125

Asked as: C drive free space; boot drive shows wrong info; system drive percent free; drive size missing on computer

Turbo shows the boot drive letter, its file system in parentheses, its size, and its free space with the percent free rounded down. When the size is missing, it shows only the letter and file system. The free part is left out when the free space is missing or zero.

## System type names

Lines 127-141

Asked as: system type number meaning; desktop or laptop label; what type of machine number; server type shows blank

Turbo names the system type from the number: 1 Desktop, 2 Mobile, 3 Workstation, 4 Enterprise Server, 5 SOHO Server, 6 Appliance PC, 7 Performance Server, 8 Maximum. Any other number shows as empty.

## Network configuration display

Lines 153-194

Asked as: network adapters list; DHCP label missing; IP addresses hidden; adapters in wrong order

Turbo shows the network config with one adapter for each MAC address, sorted by MAC address. Known properties come first, with the labels DHCP, DHCP server, IP addresses, Subnets, Gateways, DNS host name and DNS servers. Other properties follow in name order with their own names. Properties with an empty value are hidden.

## Network value formatting

Lines 164-175

Asked as: true shows as yes; list of IPs with commas; network values show blank; adapters not showing up

Turbo shows true as "Yes" and false as "No". It shows lists as values separated by commas, without empty items. It shows empty values as blank. A config that is not an object shows no adapters.

## Dates and times shown in the device's time zone

Lines 196-211

Asked as: wrong time on computer; last seen time is off; date in my time zone; time shows a different hour

Turbo shows the day and the time of day of a computer date in the time zone of the user's device. The time looks like "8:00 AM". An empty or invalid date shows as empty.

## What the computers list search matches

Lines 213-228

Asked as: search by serial number; find computer by name; search box not finding computer; search by user email

The computers list search matches the computer name, serial, maker, model, reported given name and surname, email, UPN, last user name, assignee name and OS version.
