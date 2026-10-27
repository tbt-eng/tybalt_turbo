import assert from "node:assert/strict";
import test from "node:test";
import { localIsoDate } from "../src/lib/localDate.ts";

// 03:30 UTC on Oct 27 is still the evening of Oct 26 in the Americas and
// already the afternoon of Oct 27 in Tokyo.
const instant = new Date(Date.UTC(2026, 9, 27, 3, 30));

function inZone(zone: string, fn: () => void) {
  const previous = process.env.TZ;
  process.env.TZ = zone;
  try {
    fn();
  } finally {
    if (previous === undefined) delete process.env.TZ;
    else process.env.TZ = previous;
  }
}

test("uses the device's date, not the UTC date", () => {
  inZone("America/Toronto", () => assert.equal(localIsoDate(instant), "2026-10-26"));
  inZone("America/Vancouver", () => assert.equal(localIsoDate(instant), "2026-10-26"));
  inZone("Asia/Tokyo", () => assert.equal(localIsoDate(instant), "2026-10-27"));
  inZone("UTC", () => assert.equal(localIsoDate(instant), "2026-10-27"));
});

test("local midnight stays on its own date east of UTC", () => {
  inZone("Europe/Berlin", () => assert.equal(localIsoDate(new Date(2026, 9, 27)), "2026-10-27"));
});
