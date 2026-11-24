import { downloadActiveJobsCsv, downloadTimeEntryBranchMismatchesCsv } from "$lib/navActions";
import { canViewWIPReports } from "$lib/reports/wipReports";
import { canViewBranchHours } from "$lib/reports/branchHours";
import { canViewPOApprovers } from "$lib/poApprovalLimits";
import { canViewPOGroups } from "$lib/poGroups";
import { canViewJobCreateQueue } from "$lib/jobRequests";
import { canViewComputers } from "$lib/computers";

// What a canView check may look at.
export interface NavContext {
  claims: string[];
  isBranchManager: boolean;
  jobRequestsEnabled: boolean;
  attentionCounts: Record<string, number>;
}

// Decides whether the signed-in user sees a link or button. An entry with no
// canView is shown to everyone. Show All UI skips every check.
export type NavCheck = (ctx: NavContext) => boolean;

// Shown to holders of any of the named claims.
export const anyClaim =
  (...names: string[]): NavCheck =>
  ({ claims }) =>
    names.some((name) => claims.includes(name));

// Shown only while the project requests feature is on.
const jobRequests =
  (check: NavCheck): NavCheck =>
  (ctx) =>
    ctx.jobRequestsEnabled && check(ctx);

const wipReports: NavCheck = ({ claims, isBranchManager }) =>
  canViewWIPReports(claims, isBranchManager);

export interface NavButton {
  action: string | (() => void | Promise<void>);
  icon: string;
  title: string;
  color?: string;
  canView?: NavCheck;
}

export interface NavLink {
  label: string;
  href: string;
  canView?: NavCheck;
}

export interface NavItem extends NavLink {
  buttons: NavButton[];
  children?: NavLink[];
}

export interface NavSection {
  title: string;
  items: NavItem[];
}

// The sections, links and buttons the user sees. With Show All UI on, that is
// all of them, whatever the user's claims; the server still checks access.
// Without it, each entry's canView decides, and an empty section is hidden.
export function visibleNavSections(
  sections: NavSection[],
  ctx: NavContext,
  showAllUi: boolean,
): NavSection[] {
  if (showAllUi) return sections;
  const allowed = (entry: { canView?: NavCheck }) => !entry.canView || entry.canView(ctx);
  return sections
    .map((section) => ({
      ...section,
      items: section.items
        .filter(allowed)
        .map((item) => ({
          ...item,
          buttons: item.buttons.filter(allowed),
          ...(item.children && { children: item.children.filter(allowed) }),
        }))
        .filter((item) => !item.children || item.children.length > 0),
    }))
    .filter((section) => section.items.length > 0);
}

export const navSections: NavSection[] = [
  {
    title: "Time Management",
    items: [
      {
        label: "Entries",
        href: "/time/entries/list",
        buttons: [
          {
            action: "/time/entries/add",
            icon: "feather:plus-circle",
            title: "New Entry",
            color: "green",
          },
        ],
      },
      { label: "Sheets", href: "/time/sheets/list", buttons: [] },
      {
        label: "Pending My Approval",
        href: "/time/sheets/pending",
        buttons: [],
        canView: anyClaim("tapr"),
      },
      {
        label: "Approved By Me",
        href: "/time/sheets/approved",
        buttons: [],
        canView: anyClaim("tapr"),
      },
      { label: "Shared with Me", href: "/time/sheets/shared", buttons: [] },
      {
        label: "Tracking",
        href: "/time/tracking",
        buttons: [],
        canView: anyClaim("report", "admin"),
      },
      {
        label: "Work Records",
        href: "/time/work-records",
        buttons: [],
        canView: anyClaim("report", "work_record"),
      },
      {
        label: "Amendments",
        href: "/time/amendments/pending",
        canView: anyClaim("tame", "report"),
        buttons: [
          {
            action: "/time/amendments/list",
            icon: "feather:list",
            title: "List All Amendments",
            color: "purple",
          },
          {
            action: "/time/amendments/add",
            icon: "feather:plus-circle",
            title: "New Amendment",
            color: "green",
          },
        ],
      },
      { label: "Time Off", href: "/time/off", buttons: [] },
    ],
  },
  {
    title: "Purchase Orders",
    items: [
      {
        label: "My Purchase Orders",
        href: "/pos/list",
        buttons: [
          {
            action: "/pos/add",
            icon: "feather:plus-circle",
            title: "New PO",
            color: "green",
          },
        ],
      },
      {
        label: "PO Groups",
        href: "/pos/groups",
        canView: ({ claims }) => canViewPOGroups(claims),
        buttons: [
          {
            action: "/pos/groups/add",
            canView: anyClaim("super_po"),
            icon: "feather:plus-circle",
            title: "New PO Group",
            color: "green",
          },
        ],
      },
      { label: "Rejected (Mine)", href: "/pos/rejected", buttons: [] },
      {
        label: "Pending My Approval",
        href: "/pos/pending",
        buttons: [],
        canView: anyClaim("tapr"),
      },
      // remove all active (test replacing with search)
      //{ label: "All Active", href: "/pos/active", buttons: [] },
      { label: "Search", href: "/pos/search", buttons: [] },
      // Company-wide lists for payables staff. Everyone else sees a label on
      // the purchase order in their own list.
      { label: "Stale", href: "/pos/stale", buttons: [], canView: anyClaim("payables_admin") },
      {
        label: "Expiring (Recurring)",
        href: "/pos/expiring",
        buttons: [],
        canView: anyClaim("payables_admin"),
      },
    ],
  },
  {
    title: "Expenses",
    items: [
      {
        label: "My Expenses",
        href: "/expenses/list",
        buttons: [
          {
            action: "/expenses/add",
            icon: "feather:plus-circle",
            title: "New Expense",
            color: "green",
          },
        ],
      },
      { label: "My Reimbursements", href: "/expenses/reimbursements", buttons: [] },
      {
        label: "Pending My Approval",
        href: "/expenses/pending",
        buttons: [],
        canView: anyClaim("tapr"),
      },
      {
        label: "Approved By Me",
        href: "/expenses/approved",
        buttons: [],
        canView: anyClaim("tapr"),
      },
      {
        label: "Tracking",
        href: "/expenses/tracking",
        buttons: [],
        canView: anyClaim("report", "admin"),
      },
      {
        label: "Settlement",
        href: "/expenses/settlement",
        buttons: [],
        canView: anyClaim("payables_admin"),
      },
      {
        label: "Commit Queue",
        href: "/reports/expense/queue",
        buttons: [],
        canView: anyClaim("commit"),
      },
      {
        label: "Uncommit Requests",
        href: "/expenses/uncommit-requests",
        buttons: [],
        canView: anyClaim("commit", "admin"),
      },
    ],
  },
  {
    title: "Business",
    items: [
      {
        label: "Jobs",
        href: "/jobs/list",
        buttons: [
          {
            action: "/jobs/unused",
            canView: anyClaim("job"),
            icon: "mdi:cancel",
            title: "Unused Active Jobs",
            color: "blue",
          },
          {
            action: "/jobs/stale",
            canView: anyClaim("job"),
            icon: "mdi:timeline-clock-outline",
            title: "Stale Jobs",
            color: "blue",
          },
          {
            action: "/jobs/latest",
            canView: anyClaim("job"),
            icon: "mdi:timeline-text-outline",
            title: "Latest Jobs",
            color: "blue",
          },
          {
            action: "/jobs/map",
            icon: "feather:map",
            title: "Jobs Map",
            color: "blue",
          },
          // The one way to start a project. Procurement creates jobs and
          // proposals; while project requests are on, /jobs/add sends staff
          // with the time claim to a new request.
          {
            action: "/jobs/add",
            canView: ({ claims, jobRequestsEnabled }) =>
              claims.includes("job") || (jobRequestsEnabled && claims.includes("time")),
            icon: "feather:plus-circle",
            title: "New Project or Proposal",
            color: "green",
          },
        ],
      },
      {
        label: "WIP",
        href: "",
        buttons: [],
        children: [
          { label: "My WIP", href: "/jobs/wip/my" },
          { label: "Branch", href: "/jobs/wip/branch", canView: wipReports },
          { label: "Division", href: "/jobs/wip/division", canView: wipReports },
        ],
      },
      {
        label: "Project Requests",
        href: "/job-requests",
        canView: jobRequests(anyClaim("time", "job", "project_authorization")),
        buttons: [],
      },
      {
        label: "Request Approvals",
        href: "/job-requests?queue=approval",
        buttons: [],
        canView: jobRequests(anyClaim("project_authorization")),
      },
      {
        label: "Job Create Queue",
        href: "/job-requests?queue=creation",
        buttons: [],
        canView: (ctx) =>
          canViewJobCreateQueue(ctx.claims, ctx.jobRequestsEnabled, ctx.attentionCounts),
      },
      {
        label: "Value Approvals",
        href: "/jobs/value-approvals",
        buttons: [],
        canView: anyClaim("project_authorization"),
      },
      {
        label: "Abandoned Requests",
        href: "/job-requests?queue=abandoned",
        buttons: [],
        canView: anyClaim("admin"),
      },
      {
        label: "Rate Sheets",
        href: "/rate-sheets/list",
        // Procurement, admins and revisers work on sheets. Requesters and
        // approvers of project requests read them, to find the sheet that
        // has a client's rates.
        canView: (ctx) =>
          anyClaim("job", "admin", "rate_sheet_revise")(ctx) ||
          jobRequests(anyClaim("time", "project_authorization"))(ctx),
        buttons: [],
      },
      {
        label: "Clients",
        href: "/clients/list",
        buttons: [
          {
            action: "/clients/add",
            canView: anyClaim("job"),
            icon: "feather:plus-circle",
            title: "New Client",
            color: "green",
          },
        ],
      },
      {
        label: "Vendors",
        href: "/vendors/list",
        buttons: [
          {
            action: "/vendors/add",
            canView: anyClaim("payables_admin"),
            icon: "feather:plus-circle",
            title: "New Vendor",
            color: "green",
          },
        ],
      },
      {
        label: "Absorb Actions",
        href: "/absorb/actions",
        buttons: [],
        canView: anyClaim("absorb"),
      },
      {
        label: "Absorb Proposals",
        href: "/absorb/proposals",
        buttons: [],
      },
      {
        label: "PO Approvers",
        href: "/pos/approval-limits",
        buttons: [],
        canView: ({ claims }) => canViewPOApprovers(claims),
      },
      {
        label: "Admin Profiles",
        href: "/admin_profiles/list",
        buttons: [],
        canView: anyClaim("admin", "hr", "time_off_manager", "it"),
      },
      {
        label: "Computers",
        href: "/computers/list",
        canView: ({ claims }) => canViewComputers(claims),
        buttons: [
          {
            action: "/computers/retired",
            icon: "mdi:archive-outline",
            title: "Retired Computers",
            color: "blue",
          },
          {
            action: "/computers/logins",
            icon: "mdi:login",
            title: "Logins",
            color: "blue",
          },
          {
            action: "/computers/raw_logins",
            icon: "mdi:alert-circle-outline",
            title: "Raw Logins",
            color: "blue",
          },
        ],
      },
      { label: "Claims", href: "/claims/list", buttons: [] },
      { label: "App Config", href: "/app_config", buttons: [], canView: anyClaim("admin") },
      { label: "UI Components", href: "/dev/components", buttons: [], canView: anyClaim("admin") },
    ],
  },
  {
    title: "Reports",
    items: [
      { label: "Payroll", href: "/reports/payroll", buttons: [], canView: anyClaim("report") },
      { label: "Weekly", href: "/reports/weekly", buttons: [], canView: anyClaim("report") },
      { label: "Time Off", href: "/time/off/all", buttons: [], canView: anyClaim("report") },
      {
        label: "Payables Spreadsheet",
        href: "/reports/payables-spreadsheet",
        buttons: [],
        canView: anyClaim("report"),
      },
      {
        label: "Branch Mismatches CSV",
        href: "",
        canView: anyClaim("report"),
        buttons: [
          {
            action: downloadTimeEntryBranchMismatchesCsv,
            icon: "mdi:file-download-outline",
            title: "Download CSV",
            color: "blue",
          },
        ],
      },
      {
        label: "Active Jobs",
        href: "",
        canView: anyClaim("report"),
        buttons: [
          {
            action: downloadActiveJobsCsv,
            icon: "mdi:file-download-outline",
            title: "Download CSV",
            color: "blue",
          },
        ],
      },
      {
        label: "Branch Hours",
        href: "/reports/branch-hours",
        buttons: [],
        canView: ({ claims, isBranchManager }) => canViewBranchHours(claims, isBranchManager),
      },
      { label: "KPI", href: "/reports/kpi", buttons: [], canView: anyClaim("kpi") },
    ],
  },
  {
    title: "Settings",
    items: [
      { label: "Currencies", href: "/currencies", buttons: [], canView: anyClaim("admin") },
      { label: "Time Types", href: "/timetypes", buttons: [], canView: anyClaim("admin") },
      { label: "Divisions", href: "/divisions", buttons: [], canView: anyClaim("admin") },
      {
        label: "Machine Secrets",
        href: "/machine_secrets",
        buttons: [],
        canView: anyClaim("admin"),
      },
      {
        label: "Attachment Audit",
        href: "/attachment_audit",
        buttons: [],
        canView: anyClaim("admin"),
      },
    ],
  },
];
