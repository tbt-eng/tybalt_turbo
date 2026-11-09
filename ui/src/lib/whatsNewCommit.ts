// What's new lists the newest week first, and each heading ends with the last
// commit of that week: "## November 1 to 7, 2026 (to 47c7396e)".
export function latestWhatsNewCommit(source: string): string | null {
  return source.match(/^## .* \(to ([0-9a-f]{8,40})\)$/m)?.[1] ?? null;
}

// The weeks to mark as new, given the commits on the page, newest first, and the
// commit this browser saw last. A first visit marks only the newest week. A commit
// that has left the page means that every week is new.
export function newWhatsNewWeeks(commits: string[], seen: string | null): Set<string> {
  if (seen === null) return new Set(commits.slice(0, 1));
  const index = commits.indexOf(seen);
  return new Set(index === -1 ? commits : commits.slice(0, index));
}
