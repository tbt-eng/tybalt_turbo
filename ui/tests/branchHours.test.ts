import assert from "node:assert/strict";
import test from "node:test";
import {
  branchHoursChoices,
  branchHoursProjects,
  branchHoursTotals,
  canViewBranchHours,
  yearToDate,
  type BranchHoursRow,
} from "../src/lib/reports/branchHours.ts";

function row(overrides: Partial<BranchHoursRow>): BranchHoursRow {
  return {
    direction: "received",
    jobId: "job1",
    jobNumber: "26-001",
    jobDescription: "Culvert",
    client: "City",
    workBranch: "Kitchener-Waterloo",
    uid: "u1",
    surname: "Doe",
    givenName: "Jane",
    staffBranch: "Thunder Bay",
    hours: 1,
    ...overrides,
  };
}

const rows: BranchHoursRow[] = [
  row({ uid: "u1", hours: 2 }),
  row({ uid: "u2", surname: "Roe", staffBranch: "Ottawa", hours: 5 }),
  row({ jobId: "job2", jobNumber: "26-002", jobDescription: "Bridge", hours: 4 }),
  row({ direction: "supplied", workBranch: "Toronto", hours: 3 }),
  row({
    direction: "internal",
    jobId: "",
    jobNumber: "",
    jobDescription: "",
    client: "",
    hours: 1.5,
  }),
];

test("totals hours by direction", () => {
  assert.deepEqual(branchHoursTotals(rows), { received: 11, supplied: 3, internal: 1.5 });
});

test("groups one direction into projects, largest first, with people by hours", () => {
  const projects = branchHoursProjects(rows, "received");
  assert.deepEqual(
    projects.map((project) => [project.jobNumber, project.hours, project.people.length]),
    [
      ["26-001", 7, 2],
      ["26-002", 4, 1],
    ],
  );
  assert.deepEqual(projects[0].staffBranches, ["Ottawa", "Thunder Bay"]);
  assert.deepEqual(
    projects[0].people.map((person) => person.uid),
    ["u2", "u1"],
  );
});

test("keeps the same job separate for each work branch", () => {
  const projects = branchHoursProjects(
    [row({ direction: "supplied" }), row({ direction: "supplied", workBranch: "Toronto" })],
    "supplied",
  );
  assert.deepEqual(
    projects.map((project) => project.workBranch),
    ["Kitchener-Waterloo", "Toronto"],
  );
});

test("search matches project details and people, and finds no-job time", () => {
  assert.deepEqual(
    branchHoursProjects(rows, "received", "roe ottawa").map((project) => project.jobNumber),
    ["26-001"],
  );
  assert.deepEqual(
    branchHoursProjects(rows, "received", "bridge").map((project) => project.jobNumber),
    ["26-002"],
  );
  assert.equal(branchHoursProjects(rows, "internal", "no project").length, 1);
});

test("year to date runs from January 1 to the reference day", () => {
  assert.deepEqual(yearToDate(new Date(2026, 9, 22)), { start: "2026-01-01", end: "2026-10-22" });
});

test("a branch manager can select only the branches they manage", () => {
  const branches = [
    { id: "kw", manager: "u1" },
    { id: "tb", manager: "u2" },
    { id: "to", manager: "u1" },
    { id: "ot", manager: "" },
  ];
  assert.deepEqual(
    branchHoursChoices([], branches, "u1").map((branch) => branch.id),
    ["kw", "to"],
  );
  assert.deepEqual(branchHoursChoices(["report"], branches, "u3"), []);
  assert.deepEqual(branchHoursChoices([], branches, ""), []);
});

test("kpi and admin can select any branch, even when they also manage one", () => {
  const branches = [
    { id: "kw", manager: "u1" },
    { id: "tb", manager: "u2" },
  ];
  for (const claim of ["kpi", "admin"]) {
    assert.equal(branchHoursChoices([claim], branches, "u1").length, 2);
    assert.equal(canViewBranchHours([claim]), true);
  }
  assert.equal(canViewBranchHours([], true), true);
  assert.equal(canViewBranchHours(["report"]), false);
});
