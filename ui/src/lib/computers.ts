// The computers that report at sign-in, as GET /api/computers and
// GET /api/computers/{id} return them, and the rules the pages show. The
// server checks each action again (app/routes/computers.go).
import { localIsoDate } from "./localDate.ts";

export interface ComputerRow {
  id: string;
  computer_name: string;
  serial: string;
  mfg: string;
  model: string;
  system_type: number;
  ram: number;
  os_version: string;
  os_sku: number;
  os_arch: string;
  boot_drive: string;
  boot_drive_fs: string;
  boot_drive_cap: number;
  boot_drive_free: number;
  radiator_version: number;
  network_config: Record<string, unknown> | null;
  user_given_name: string;
  user_surname: string;
  email: string;
  upn: string;
  // The Turbo user that the last sign-in matched, or "".
  last_user: string;
  last_user_name: string;
  last_login: string;
  assigned_to: string;
  assigned_to_name: string;
  assigned_time: string;
  retired: string;
  created: string;
  updated: string;
}

// The pages, the list routes and the actions (the computers claim or admin).
export function canViewComputers(claims: string[]): boolean {
  return claims.includes("computers") || claims.includes("admin");
}

// unassigned: nobody. other: someone who is not the last user.
// last_user: the last user, the state that needs nothing.
export type AssignmentState = "unassigned" | "other" | "last_user";

export interface ComputerAssignment {
  state: AssignmentState;
  // Assign sets the assignee to the last user. It needs a matched last user
  // and a computer that is not retired.
  canAssign: boolean;
}

export function computerAssignment(
  computer: Pick<ComputerRow, "assigned_to" | "last_user" | "retired">,
): ComputerAssignment {
  const state: AssignmentState =
    computer.assigned_to === ""
      ? "unassigned"
      : computer.assigned_to === computer.last_user
        ? "last_user"
        : "other";
  return {
    state,
    canAssign: state !== "last_user" && computer.last_user !== "" && computer.retired === "",
  };
}

// The person of the last sign-in: the Turbo user it matched, or the name it
// reported.
export function lastUserLabel(
  computer: Pick<ComputerRow, "last_user_name" | "user_given_name" | "user_surname">,
): string {
  return (
    computer.last_user_name ||
    `${computer.user_given_name} ${computer.user_surname}`.trim() ||
    "Unknown"
  );
}

// The maker and model. A model that starts with the maker's first word, such
// as "HP EliteBook" from "HP Inc.", is shown alone.
export function modelLabel(computer: Pick<ComputerRow, "mfg" | "model">): string {
  const mfg = computer.mfg.trim();
  const model = computer.model.trim();
  const maker = mfg.split(/\s+/)[0]?.toLowerCase() ?? "";
  if (model && maker && model.toLowerCase().startsWith(maker)) return model;
  return [mfg, model].filter(Boolean).join(" ");
}

// A size in bytes, in powers of 1024 with the units that Windows shows:
// 17179869184 is "16 GB".
export function formatBytes(bytes: number): string {
  if (!Number.isFinite(bytes) || bytes <= 0) return "";
  const units = ["B", "KB", "MB", "GB", "TB", "PB"];
  let value = bytes;
  let unit = 0;
  while (value >= 1024 && unit < units.length - 1) {
    value /= 1024;
    unit += 1;
  }
  const rounded = unit === 0 ? value : Math.round(value * 10) / 10;
  return `${rounded} ${units[unit]}`;
}

// The boot drive: its letter, file system, size and free space.
export function bootDriveLabel(
  computer: Pick<
    ComputerRow,
    "boot_drive" | "boot_drive_fs" | "boot_drive_cap" | "boot_drive_free"
  >,
): string {
  const name = [computer.boot_drive, computer.boot_drive_fs && `(${computer.boot_drive_fs})`]
    .filter(Boolean)
    .join(" ");
  const size = formatBytes(computer.boot_drive_cap);
  if (!size) return name;
  const free = formatBytes(computer.boot_drive_free);
  const share =
    computer.boot_drive_cap > 0
      ? ` (${Math.floor((100 * computer.boot_drive_free) / computer.boot_drive_cap)}%)`
      : "";
  return `${name} ${size}${free ? `, ${free} free${share}` : ""}`.trim();
}

// Win32_ComputerSystem.PCSystemType.
const systemTypes: Record<number, string> = {
  1: "Desktop",
  2: "Mobile",
  3: "Workstation",
  4: "Enterprise Server",
  5: "SOHO Server",
  6: "Appliance PC",
  7: "Performance Server",
  8: "Maximum",
};

export function systemTypeLabel(systemType: number): string {
  return systemTypes[systemType] ?? "";
}

export interface NetworkField {
  label: string;
  value: string;
}

export interface NetworkAdapter {
  mac: string;
  fields: NetworkField[];
}

// The properties radiator sends for each adapter (radiator.ps1).
const networkLabels: Record<string, string> = {
  dhcpEnabled: "DHCP",
  dhcpServer: "DHCP server",
  ips: "IP addresses",
  subnets: "Subnets",
  gateways: "Gateways",
  dnsHostname: "DNS host name",
  dnsOrder: "DNS servers",
};

function networkValue(value: unknown): string {
  if (value === null || value === undefined) return "";
  if (typeof value === "boolean") return value ? "Yes" : "No";
  if (Array.isArray(value)) return value.map(networkValue).filter(Boolean).join(", ");
  if (typeof value === "object") return JSON.stringify(value);
  return String(value);
}

// The network config, one adapter for each MAC address, in MAC order. Known
// properties come first with readable labels; others keep their names.
export function networkAdapters(config: unknown): NetworkAdapter[] {
  if (config === null || typeof config !== "object" || Array.isArray(config)) return [];
  const known = Object.keys(networkLabels);
  return Object.entries(config as Record<string, unknown>)
    .sort(([a], [b]) => a.localeCompare(b))
    .map(([mac, properties]) => {
      const entries =
        properties !== null && typeof properties === "object" && !Array.isArray(properties)
          ? Object.entries(properties as Record<string, unknown>)
          : [["value", properties] as [string, unknown]];
      const rank = (key: string) => {
        const at = known.indexOf(key);
        return at < 0 ? known.length : at;
      };
      const fields = entries
        .sort(([a], [b]) => rank(a) - rank(b) || a.localeCompare(b))
        .map(([key, value]) => ({ label: networkLabels[key] ?? key, value: networkValue(value) }))
        .filter((field) => field.value !== "");
      return { mac, fields };
    });
}

// The day of a PocketBase datetime on the user's device, as YYYY-MM-DD for
// shortDate().
export function localDay(datetime: string): string {
  if (!datetime) return "";
  const date = new Date(datetime.replace(" ", "T"));
  return Number.isNaN(date.getTime()) ? "" : localIsoDate(date);
}

// The time of day of a PocketBase datetime on the user's device, such as
// "8:00 AM".
export function localTime(datetime: string): string {
  if (!datetime) return "";
  const date = new Date(datetime.replace(" ", "T"));
  if (Number.isNaN(date.getTime())) return "";
  return date.toLocaleTimeString("en-US", { hour: "numeric", minute: "2-digit" });
}

// What the list search matches.
export function computerSearchText(computer: ComputerRow): string {
  return [
    computer.computer_name,
    computer.serial,
    computer.mfg,
    computer.model,
    computer.user_given_name,
    computer.user_surname,
    computer.email,
    computer.upn,
    computer.last_user_name,
    computer.assigned_to_name,
    computer.os_version,
  ].join(" ");
}

// One sign-in, as GET /api/computers/logins returns it. user_name is the Turbo
// user that the sign-in matched, which the server joins because a holder of
// the computers claim can't read other people's profiles.
export interface ComputerLoginRow {
  id: string;
  created: string;
  computer: string;
  computer_name: string;
  serial: string;
  given_name: string;
  surname: string;
  // The Turbo user that the sign-in matched, or "".
  user: string;
  user_name: string;
}

// One rejected report, as GET /api/computers/raw_logins returns it. Only
// GET /api/computers/raw_logins/{id} fills report: the report as radiator
// sent it. It is untrusted, so the page shows it as text.
export interface ComputerRawLoginRow {
  id: string;
  created: string;
  computer_name: string;
  error: string;
  report?: unknown;
}

// One page of a server-paged list, in the shape of a PocketBase list.
export interface ComputerPage<T> {
  items: T[];
  page: number;
  perPage: number;
  totalItems: number;
  totalPages: number;
}

// The search and the page of a server-paged list, as the URL holds them.
export interface ComputerListState {
  q: string;
  page: number;
  perPage: number;
}

// The page sizes that DSList offers. Another size in the URL gets the default.
const listPageSizes = [10, 20, 50];
const defaultPerPage = 20;

export function computerListState(params: URLSearchParams): ComputerListState {
  const page = Number(params.get("page"));
  const perPage = Number(params.get("perPage"));
  return {
    q: params.get("q")?.trim() ?? "",
    page: Number.isSafeInteger(page) && page > 0 ? page : 1,
    perPage: listPageSizes.includes(perPage) ? perPage : defaultPerPage,
  };
}

// The URL search for a list state, without the defaults: "" for the first
// page of everything, "?q=pc&page=2" for the second page of a search.
export function computerListSearch(state: ComputerListState): string {
  const params = new URLSearchParams();
  const q = state.q.trim();
  if (q) params.set("q", q);
  if (state.page > 1) params.set("page", String(state.page));
  if (state.perPage !== defaultPerPage) params.set("perPage", String(state.perPage));
  const search = params.toString();
  return search ? `?${search}` : "";
}

// The name that a sign-in reported.
export function loginReportedName(login: Pick<ComputerLoginRow, "given_name" | "surname">): string {
  return `${login.given_name} ${login.surname}`.trim() || "Unknown";
}

// The Turbo user that a sign-in matched, when its name differs from the name
// the sign-in reported. "" when it matched no user or the names are the same.
export function loginTurboUserName(
  login: Pick<ComputerLoginRow, "given_name" | "surname" | "user" | "user_name">,
): string {
  if (!login.user || !login.user_name) return "";
  const same = login.user_name.trim().toLowerCase() === loginReportedName(login).toLowerCase();
  return same ? "" : login.user_name;
}

// A rejected report as indented JSON text.
export function rawLoginReportText(report: unknown): string {
  if (report === undefined) return "";
  if (typeof report === "string") return report;
  return JSON.stringify(report, null, 2) ?? "";
}
