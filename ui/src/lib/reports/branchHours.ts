export type BranchHoursDirection = "received" | "supplied" | "internal";

export const branchHoursDirections: { key: BranchHoursDirection; label: string }[] = [
  { key: "received", label: "Received" },
  { key: "supplied", label: "Supplied" },
  { key: "internal", label: "Internal" },
];

export type BranchHoursRow = {
  direction: BranchHoursDirection;
  jobId: string;
  jobNumber: string;
  jobDescription: string;
  client: string;
  workBranch: string;
  uid: string;
  payrollId: string;
  surname: string;
  givenName: string;
  staffBranch: string;
  hours: number;
};

export type BranchHoursReport = {
  branch: string;
  rows: BranchHoursRow[];
};

// One table row: a job (or no-job time) for one work branch, with the people
// who recorded its hours.
export type BranchHoursProject = {
  key: string;
  jobId: string;
  jobNumber: string;
  jobDescription: string;
  client: string;
  workBranch: string;
  staffBranches: string[];
  hours: number;
  people: BranchHoursRow[];
};

export function canViewBranchHours(claims: string[], isBranchManager = false): boolean {
  return isBranchManager || seesAllBranchHours(claims);
}

export function seesAllBranchHours(claims: string[]): boolean {
  return claims.includes("kpi") || claims.includes("admin");
}

// The kpi and admin claims may select any branch. Without them, a branch
// manager may select only the branches they manage, as the API enforces.
export function branchHoursChoices<T extends { id: string; manager?: string }>(
  claims: string[],
  branches: T[],
  userId: string,
): T[] {
  if (seesAllBranchHours(claims)) return branches;
  return branches.filter((branch) => !!userId && branch.manager === userId);
}

export function branchHoursTotals(rows: BranchHoursRow[]): Record<BranchHoursDirection, number> {
  const totals = { received: 0, supplied: 0, internal: 0 };
  for (const row of rows) totals[row.direction] += row.hours;
  return totals;
}

// Group one direction's rows into projects, largest first. Every search word
// must appear in the project's details or one of its people.
export function branchHoursProjects(
  rows: BranchHoursRow[],
  direction: BranchHoursDirection,
  search = "",
): BranchHoursProject[] {
  const projects = new Map<string, BranchHoursProject>();
  for (const row of rows) {
    if (row.direction !== direction) continue;
    const key = `${row.jobId}|${row.workBranch}`;
    let project = projects.get(key);
    if (!project) {
      project = {
        key,
        jobId: row.jobId,
        jobNumber: row.jobNumber,
        jobDescription: row.jobDescription,
        client: row.client,
        workBranch: row.workBranch,
        staffBranches: [],
        hours: 0,
        people: [],
      };
      projects.set(key, project);
    }
    project.hours += row.hours;
    project.people.push(row);
    if (row.staffBranch && !project.staffBranches.includes(row.staffBranch)) {
      project.staffBranches.push(row.staffBranch);
    }
  }

  const words = search.trim().toLocaleLowerCase().split(/\s+/).filter(Boolean);
  return [...projects.values()]
    .filter((project) => {
      const text = [
        project.jobNumber || "No project",
        project.jobDescription,
        project.client,
        project.workBranch,
        ...project.people.flatMap((person) => [
          person.givenName,
          person.surname,
          person.staffBranch,
        ]),
      ]
        .join(" ")
        .toLocaleLowerCase();
      return words.every((word) => text.includes(word));
    })
    .map((project) => ({
      ...project,
      staffBranches: project.staffBranches.sort(),
      people: project.people.sort(
        (a, b) =>
          b.hours - a.hours ||
          a.surname.localeCompare(b.surname) ||
          a.givenName.localeCompare(b.givenName),
      ),
    }))
    .sort(
      (a, b) =>
        b.hours - a.hours ||
        Number(!a.jobId) - Number(!b.jobId) ||
        a.jobNumber.localeCompare(b.jobNumber) ||
        a.workBranch.localeCompare(b.workBranch),
    );
}

export function formatHours(hours: number): string {
  return hours.toLocaleString("en-CA", { minimumFractionDigits: 0, maximumFractionDigits: 2 });
}

export function yearToDate(reference = new Date()): { start: string; end: string } {
  const pad = (value: number) => String(value).padStart(2, "0");
  const year = reference.getFullYear();
  return {
    start: `${year}-01-01`,
    end: `${year}-${pad(reference.getMonth() + 1)}-${pad(reference.getDate())}`,
  };
}
