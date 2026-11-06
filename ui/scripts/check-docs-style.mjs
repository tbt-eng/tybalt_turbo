#!/usr/bin/env node
// Checks the writing rules in src/docs/README.md. It can't judge meaning, so it
// checks what a script can check reliably: sentence length, one instruction in
// a step, passive voice, steps that start with a verb, and the glossary's
// replaced words.
import { readFileSync, readdirSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const maxStepWords = 20;
const maxDescriptionWords = 25;
const passive =
  /\b(?:is|are|was|were|be|been|being)\s+(?:not\s+)?(?:\w+ed|made|sent|paid|set|done|given|taken|shown|held|kept|found)\b/i;
const notAVerb = /^(?:the|a|an|you|your|it|this|these|that|those)\b/i;

/** @param {string} text */
function words(text) {
  return text.match(/[A-Za-z$][\w'’$-]*/g)?.length ?? 0;
}

/** @param {string} text */
function plain(text) {
  return text
    .replace(/\[([^\]]+)\]\([^)]+\)/g, "$1")
    .replace(/[*`]/g, "")
    .trim();
}

/** @param {string} text */
function sentences(text) {
  return plain(text)
    .split(/(?<=[.?!])\s+/)
    .filter(Boolean);
}

// The third column of the glossary table lists words that the docs must not use.
/** @param {string} glossary */
export function replacedWords(glossary) {
  return glossary
    .split("\n")
    .filter((line) => line.startsWith("|"))
    .map((line) => line.split("|")[3]?.trim() ?? "")
    .filter((word) => word && word !== "Don't use" && !/^-+$/.test(word));
}

/**
 * @param {string} source
 * @param {string[]} [replaced]
 */
export function checkDocStyle(source, replaced = []) {
  /** @type {{ line: number, message: string }[]} */
  const problems = [];
  let fence = "";
  source.split("\n").forEach((raw, index) => {
    /** @param {string} message */
    const report = (message) => problems.push({ line: index + 1, message });
    const line = raw.trim();
    if (line.startsWith("```")) {
      fence = fence ? "" : line.slice(3) || "code";
      return;
    }
    if (fence && fence !== "flow") return;
    if (!line || line.startsWith("#") || /^\|[\s:|-]+\|$/.test(line)) return;

    let texts = [line];
    let limit = maxDescriptionWords;
    const step = line.match(/^\d+\.\s+(.*)$/);
    if (step) {
      texts = [step[1]];
      limit = maxStepWords;
      if (notAVerb.test(plain(step[1]))) report(`start this step with a verb: "${plain(step[1])}"`);
      if (/,\s+then\s/i.test(step[1]))
        report(`use one instruction in each sentence: "${plain(step[1])}"`);
    } else if (line.startsWith("|") || fence === "flow") {
      texts = line.split("|").map((cell) => cell.trim());
    } else {
      texts = [line.replace(/^[-*]\s+/, "").replace(/^Note:\s*/, "")];
    }

    for (const text of texts) {
      for (const sentence of sentences(text)) {
        if (words(sentence) > limit) {
          report(`${words(sentence)} words; the limit is ${limit}: "${sentence}"`);
        }
        if (passive.test(sentence)) report(`use active voice: "${sentence}"`);
      }
    }
    for (const word of replaced) {
      if (new RegExp(`(?<![\\w-])${word}(?![\\w-])`, "i").test(plain(line))) {
        report(`use the glossary word, not "${word}"`);
      }
    }
  });
  return problems;
}

if (process.argv[1] === fileURLToPath(import.meta.url)) {
  const docsDir = join(dirname(fileURLToPath(import.meta.url)), "..", "src", "docs");
  const files = readdirSync(docsDir).filter((f) => f.endsWith(".md") && f !== "README.md");
  const replaced = replacedWords(readFileSync(join(docsDir, "glossary.md"), "utf8"));
  let failures = 0;
  for (const file of files) {
    // The glossary has to name the words that it replaces.
    const list = file === "glossary.md" ? [] : replaced;
    for (const { line, message } of checkDocStyle(
      readFileSync(join(docsDir, file), "utf8"),
      list,
    )) {
      console.log(`error: src/docs/${file}:${line}: ${message}`);
      failures++;
    }
  }
  if (failures) {
    console.error(`check-docs-style: ${failures} problem(s). See src/docs/README.md.`);
    process.exit(1);
  }
  console.log("check-docs-style: OK");
}
