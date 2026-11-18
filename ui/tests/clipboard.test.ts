import test from "node:test";
import assert from "node:assert/strict";
import { copyText } from "../src/lib/clipboard.ts";

test("a copy writes the text and reports success", async () => {
  const written: string[] = [];
  const clipboard = { writeText: async (text: string) => void written.push(text) };
  assert.equal(await copyText("2026-02-00001", clipboard), true);
  assert.deepEqual(written, ["2026-02-00001"]);
});

test("a browser with no clipboard reports a failure and doesn't throw", async () => {
  assert.equal(await copyText("2026-02-00001", undefined), false);
  assert.equal(await copyText("2026-02-00001", {} as Pick<Clipboard, "writeText">), false);
});

test("a refused copy reports a failure and doesn't throw", async () => {
  const clipboard = {
    writeText: async () => {
      throw new DOMException("Write permission denied.", "NotAllowedError");
    },
  };
  assert.equal(await copyText("2026-02-00001", clipboard), false);
});
