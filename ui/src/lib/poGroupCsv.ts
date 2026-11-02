// CSV handling for purchase order groups (issue #168). The upload has one row
// per staff member: an email column that picks the owner, plus a column for
// each template field the creator chose to take from the upload.

export type PoGroupUploadField = "total" | "description" | "division" | "branch";

export const poGroupUploadFieldOptions: {
  id: PoGroupUploadField;
  label: string;
  hint: string;
}[] = [
  { id: "total", label: "Total", hint: "an amount, such as 250 or 250.00" },
  { id: "description", label: "Description", hint: "at least 4 characters" },
  { id: "division", label: "Division", hint: "the division code, such as CI" },
  { id: "branch", label: "Branch", hint: "the branch name or code" },
];

// The columns a CSV needs for the chosen upload fields, in a fixed order.
export function poGroupColumns(uploadFields: PoGroupUploadField[]): string[] {
  return [
    "email",
    ...poGroupUploadFieldOptions.map((o) => o.id).filter((id) => uploadFields.includes(id)),
  ];
}

// Parse CSV text into rows of cells. Handles quoted cells with commas, doubled
// quotes and line breaks, CRLF line endings, and a leading byte order mark.
export function parseCsv(text: string): string[][] {
  const rows: string[][] = [];
  let row: string[] = [];
  let cell = "";
  let quoted = false;
  const input = text.charCodeAt(0) === 0xfeff ? text.slice(1) : text;

  for (let i = 0; i < input.length; i++) {
    const char = input[i];
    if (quoted) {
      if (char === '"' && input[i + 1] === '"') {
        cell += '"';
        i++;
      } else if (char === '"') {
        quoted = false;
      } else {
        cell += char;
      }
    } else if (char === '"') {
      quoted = true;
    } else if (char === ",") {
      row.push(cell);
      cell = "";
    } else if (char === "\n" || char === "\r") {
      if (char === "\r" && input[i + 1] === "\n") i++;
      row.push(cell);
      rows.push(row);
      row = [];
      cell = "";
    } else {
      cell += char;
    }
  }
  if (cell !== "" || row.length > 0) {
    row.push(cell);
    rows.push(row);
  }
  return rows;
}

export type PoGroupCsv = {
  rows: Record<string, string>[];
  // Row i of rows is on this line of the file, for pointing at errors.
  lines: number[];
  missingColumns: string[];
};

// Read an uploaded CSV into rows keyed by column name. Column names are
// matched without regard to case or surrounding spaces, blank lines are
// skipped, and columns the upload doesn't need are ignored.
export function readPoGroupCsv(text: string, uploadFields: PoGroupUploadField[]): PoGroupCsv {
  const cells = parseCsv(text);
  const isBlank = (r: string[]) => r.every((c) => c.trim() === "");
  const headerIndex = cells.findIndex((r) => !isBlank(r));
  if (headerIndex === -1) {
    return { rows: [], lines: [], missingColumns: poGroupColumns(uploadFields) };
  }

  const headers = cells[headerIndex].map((h) => h.trim().toLowerCase());
  const wanted = poGroupColumns(uploadFields);
  const missingColumns = wanted.filter((column) => !headers.includes(column));

  const rows: Record<string, string>[] = [];
  const lines: number[] = [];
  cells.slice(headerIndex + 1).forEach((cellRow, offset) => {
    if (isBlank(cellRow)) return;
    const row: Record<string, string> = {};
    for (const column of wanted) {
      const at = headers.indexOf(column);
      row[column] = at === -1 ? "" : (cellRow[at] ?? "").trim();
    }
    rows.push(row);
    lines.push(headerIndex + offset + 2);
  });
  return { rows, lines, missingColumns };
}

// A random id in the shape PocketBase uses, made before submitting so that a
// repeated request finds the group it already created.
export function newPoGroupId(): string {
  const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789";
  const bytes = crypto.getRandomValues(new Uint8Array(15));
  return Array.from(bytes, (b) => alphabet[b % alphabet.length]).join("");
}
