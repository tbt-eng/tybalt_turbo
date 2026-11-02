import assert from "node:assert/strict";
import test from "node:test";
import { newPoGroupId, parseCsv, poGroupColumns, readPoGroupCsv } from "../src/lib/poGroupCsv.ts";

test("parseCsv handles quotes, commas, line breaks and CRLF", () => {
  const text =
    String.fromCharCode(0xfeff) +
    'email,description\r\n"a@x.com","Boots, safety"\r\nb@x.com,"Say ""hi""\nthere"\r\n';
  assert.deepEqual(parseCsv(text), [
    ["email", "description"],
    ["a@x.com", "Boots, safety"],
    ["b@x.com", 'Say "hi"\nthere'],
  ]);
});

test("parseCsv keeps a last row without a trailing newline", () => {
  assert.deepEqual(parseCsv("email\na@x.com"), [["email"], ["a@x.com"]]);
});

test("poGroupColumns puts email first and keeps a fixed order", () => {
  assert.deepEqual(poGroupColumns(["branch", "total"]), ["email", "total", "branch"]);
  assert.deepEqual(poGroupColumns([]), ["email"]);
});

test("readPoGroupCsv matches columns loosely and points rows at file lines", () => {
  const text = '\n Email ,TOTAL,Name\na@x.com, 250 ,Ann\n\nb@x.com,"$1,000.00",Bob\n';
  const csv = readPoGroupCsv(text, ["total"]);
  assert.deepEqual(csv.missingColumns, []);
  assert.deepEqual(csv.rows, [
    { email: "a@x.com", total: "250" },
    { email: "b@x.com", total: "$1,000.00" },
  ]);
  assert.deepEqual(csv.lines, [3, 5]);
});

test("readPoGroupCsv reports columns the upload needs but lacks", () => {
  const csv = readPoGroupCsv("email\na@x.com\n", ["total", "division"]);
  assert.deepEqual(csv.missingColumns, ["total", "division"]);
  assert.deepEqual(csv.rows, [{ email: "a@x.com", total: "", division: "" }]);
});

test("readPoGroupCsv treats an empty file as missing every column", () => {
  const csv = readPoGroupCsv(" \n", ["total"]);
  assert.deepEqual(csv.rows, []);
  assert.deepEqual(csv.missingColumns, ["email", "total"]);
});

test("newPoGroupId makes 15 lowercase letters or digits", () => {
  assert.match(newPoGroupId(), /^[a-z0-9]{15}$/);
});
