#!/usr/bin/env node
// Warns when a commit changes what a user sees or does but not the user docs.
// A commit that needs no docs change can say so with a "Docs: none" line. This
// only warns and always exits 0, so a false alarm can't delay a fix.
//
//   node scripts/check-docs-commit.mjs            check the latest commit
//   node scripts/check-docs-commit.mjs A..B       check a range of commits
import { execFileSync } from "node:child_process";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const repoRoot = join(dirname(fileURLToPath(import.meta.url)), "..", "..");

const userFacing = [
  /^ui\/src\/routes\/(?!dev\/)/,
  /^ui\/src\/lib\/components\//,
  /^app\/hooks\/(?!.*_test\.go$)/,
  /^app\/routes\/(?!.*_test\.go$)/,
];

/**
 * @param {string[]} files
 * @param {string} message
 */
export function needsDocsWarning(files, message) {
  if (/^Docs: none\s*$/m.test(message)) return false;
  if (files.some((file) => file.startsWith("ui/src/docs/"))) return false;
  return files.some((file) => userFacing.some((pattern) => pattern.test(file)));
}

/** @param {string[]} args */
function git(...args) {
  return execFileSync("git", ["-C", repoRoot, ...args], { encoding: "utf8" });
}

if (process.argv[1] === fileURLToPath(import.meta.url)) {
  try {
    const range = process.argv[2];
    const commits = range
      ? git("rev-list", "--no-merges", range).split("\n").filter(Boolean)
      : [git("rev-parse", "HEAD").trim()];
    for (const commit of commits) {
      const message = git("show", "-s", "--format=%B", commit);
      const files = git("show", "--name-only", "--format=", commit).split("\n").filter(Boolean);
      if (needsDocsWarning(files, message)) {
        const subject = message.split("\n")[0];
        const text = `${commit.slice(0, 8)} "${subject}" changes user-facing code but not ui/src/docs. Update the docs, or add a "Docs: none" line if it needs none.`;
        console.log(process.env.GITHUB_ACTIONS ? `::warning::${text}` : `warning: ${text}`);
      }
    }
  } catch (error) {
    console.log(`warning: check-docs-commit could not run: ${String(error).split("\n")[0]}`);
  }
}
