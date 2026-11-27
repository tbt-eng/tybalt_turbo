import test from "node:test";
import assert from "node:assert/strict";
import {
  bootDriveLabel,
  canViewComputers,
  computerAssignment,
  computerListSearch,
  formatBytes,
  lastUserLabel,
  localDay,
  localTime,
  loginReportedName,
  loginSearchText,
  loginTurboUserName,
  modelLabel,
  networkAdapters,
  rawLoginReportText,
  rawLoginSearchText,
  systemTypeLabel,
} from "../src/lib/computers.ts";

test("the computers claim and admin see computers; others don't", () => {
  assert.equal(canViewComputers(["computers"]), true);
  assert.equal(canViewComputers(["admin"]), true);
  assert.equal(canViewComputers(["it", "time"]), false);
  assert.equal(canViewComputers([]), false);
});

test("a computer assigned to nobody can be assigned to its last user", () => {
  assert.deepEqual(computerAssignment({ assigned_to: "", last_user: "u1", retired: "" }), {
    state: "unassigned",
    canAssign: true,
  });
});

test("a computer assigned to someone else can be assigned to its last user", () => {
  assert.deepEqual(computerAssignment({ assigned_to: "u2", last_user: "u1", retired: "" }), {
    state: "other",
    canAssign: true,
  });
});

test("a computer assigned to its last user needs nothing", () => {
  assert.deepEqual(computerAssignment({ assigned_to: "u1", last_user: "u1", retired: "" }), {
    state: "last_user",
    canAssign: false,
  });
});

test("assign is not offered without a matched last user", () => {
  assert.deepEqual(computerAssignment({ assigned_to: "", last_user: "", retired: "" }), {
    state: "unassigned",
    canAssign: false,
  });
  assert.deepEqual(computerAssignment({ assigned_to: "u2", last_user: "", retired: "" }), {
    state: "other",
    canAssign: false,
  });
});

test("assign is not offered on a retired computer", () => {
  const retired = "2026-06-01 00:00:00.000Z";
  assert.equal(computerAssignment({ assigned_to: "", last_user: "u1", retired }).canAssign, false);
});

test("the last user is the matched user, else the reported name", () => {
  const reported = { user_given_name: "Fake", user_surname: "Person" };
  assert.equal(lastUserLabel({ last_user_name: "Tester Time", ...reported }), "Tester Time");
  assert.equal(lastUserLabel({ last_user_name: "", ...reported }), "Fake Person");
  assert.equal(
    lastUserLabel({ last_user_name: "", user_given_name: "", user_surname: "" }),
    "Unknown",
  );
});

test("the maker is left out when the model already names it", () => {
  assert.equal(modelLabel({ mfg: "HP", model: "HP EliteBook 840" }), "HP EliteBook 840");
  assert.equal(modelLabel({ mfg: "Dell Inc.", model: "Latitude 5440" }), "Dell Inc. Latitude 5440");
  assert.equal(modelLabel({ mfg: "LENOVO", model: "" }), "LENOVO");
});

test("sizes are in powers of 1024 with Windows units", () => {
  assert.equal(formatBytes(17179869184), "16 GB");
  assert.equal(formatBytes(510770802688), "475.7 GB");
  assert.equal(formatBytes(512), "512 B");
  assert.equal(formatBytes(1536), "1.5 KB");
  assert.equal(formatBytes(0), "");
});

test("the boot drive shows its size, free space and share free", () => {
  assert.equal(
    bootDriveLabel({
      boot_drive: "C:",
      boot_drive_fs: "NTFS",
      boot_drive_cap: 510770802688,
      boot_drive_free: 212345678848,
    }),
    "C: (NTFS) 475.7 GB, 197.8 GB free (41%)",
  );
  assert.equal(
    bootDriveLabel({ boot_drive: "C:", boot_drive_fs: "", boot_drive_cap: 0, boot_drive_free: 0 }),
    "C:",
  );
});

test("system types have their Windows names", () => {
  assert.equal(systemTypeLabel(1), "Desktop");
  assert.equal(systemTypeLabel(2), "Mobile");
  assert.equal(systemTypeLabel(0), "");
});

test("the network config is one adapter for each MAC address", () => {
  const adapters = networkAdapters({
    "00:11:22:33:44:67": { dhcpEnabled: false, ips: ["192.168.1.30"] },
    "00:11:22:33:44:66": {
      zeta: "z",
      dnsOrder: ["10.0.0.2"],
      ips: ["10.0.0.30", "fe80::1"],
      dhcpEnabled: true,
      empty: [],
    },
  });
  assert.deepEqual(adapters, [
    {
      mac: "00:11:22:33:44:66",
      fields: [
        { label: "DHCP", value: "Yes" },
        { label: "IP addresses", value: "10.0.0.30, fe80::1" },
        { label: "DNS servers", value: "10.0.0.2" },
        { label: "zeta", value: "z" },
      ],
    },
    {
      mac: "00:11:22:33:44:67",
      fields: [
        { label: "DHCP", value: "No" },
        { label: "IP addresses", value: "192.168.1.30" },
      ],
    },
  ]);
  assert.deepEqual(networkAdapters(null), []);
  assert.deepEqual(networkAdapters([]), []);
});

test("a PocketBase datetime becomes the day on the user's device", () => {
  const previous = process.env.TZ;
  process.env.TZ = "America/Toronto";
  try {
    assert.equal(localDay("2026-11-01 03:00:00.000Z"), "2026-10-31");
    assert.equal(localDay(""), "");
    assert.equal(localTime("2026-11-01 13:00:00.000Z"), "8:00 AM");
    assert.equal(localTime(""), "");
  } finally {
    process.env.TZ = previous;
  }
});

test("a search becomes a URL search, and no search none", () => {
  assert.equal(computerListSearch(""), "");
  assert.equal(computerListSearch("  "), "");
  assert.equal(computerListSearch(" pc "), "?q=pc");
  assert.equal(computerListSearch("a&b"), "?q=a%26b");
  assert.equal(new URLSearchParams(computerListSearch("tester time")).get("q"), "tester time");
});

test("the list search matches the fields that the server searches", () => {
  const login = {
    id: "l1",
    created: "",
    computer: "c1",
    computer_name: "FAKE-PC-0003",
    serial: "FAKESERIAL0003",
    given_name: "Guest",
    surname: "Account",
    user: "",
    user_name: "",
  };
  const text = loginSearchText(login).toLowerCase();
  for (const q of ["fake-pc-0003", "fakeserial0003", "guest account"]) {
    assert.ok(text.includes(q), q);
  }
  assert.ok(loginSearchText({ ...login, user_name: "Tester Time" }).includes("Tester Time"));
  const raw = { id: "r1", created: "", computer_name: "FAKE-PC-0098", error: "upn: is required" };
  assert.ok(rawLoginSearchText(raw).includes("FAKE-PC-0098"));
  assert.ok(rawLoginSearchText(raw).includes("upn: is required"));
});

test("a sign-in shows its reported name and a Turbo user only when the names differ", () => {
  const login = { given_name: "Tester", surname: "Time", user: "u1", user_name: "Tester Time" };
  assert.equal(loginReportedName(login), "Tester Time");
  assert.equal(loginTurboUserName(login), "");
  assert.equal(loginTurboUserName({ ...login, user_name: "tester time " }), "");
  assert.equal(loginTurboUserName({ ...login, user_name: "Fatty Maclean" }), "Fatty Maclean");
  assert.equal(loginTurboUserName({ ...login, user: "" }), "");
  assert.equal(loginReportedName({ given_name: "", surname: "" }), "Unknown");
});

test("a rejected report shows as indented JSON text", () => {
  assert.equal(
    rawLoginReportText({ computerName: "<b>PC</b>", ram: 8 }),
    '{\n  "computerName": "<b>PC</b>",\n  "ram": 8\n}',
  );
  assert.equal(rawLoginReportText(null), "null");
  assert.equal(rawLoginReportText(undefined), "");
  assert.equal(rawLoginReportText("not json"), "not json");
});
