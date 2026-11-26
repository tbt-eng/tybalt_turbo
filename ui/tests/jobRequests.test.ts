import test from "node:test";
import assert from "node:assert/strict";
import {
  canViewJobCreateQueue,
  copiedAllocations,
  hasAllocations,
  printMissing,
  contactEmail,
  billingAddress,
  wholeCAD,
  emptyChecklist,
  proposalAcceptsRequest,
  requestNextStep,
  type RequestDetail,
} from "../src/lib/jobRequests.ts";
const complete = (): RequestDetail => ({
  checklist: { version: 1, questions: [] },
  request: {
    id: "request1",
    version: 2,
    description: "Field inspection\nPrepare a report",
    location: "86VHQPRM+5C",
    project_value: 120000,
    time_and_materials: false,
    time_and_materials_set: true,
  },
  allocations: [],
  review_token: "",
  permissions: {},
  related: {
    client: {
      id: "client1",
      name: "Example Client",
      address: "1 Main Street",
      city: "Thunder Bay",
      province_state: "ON",
    },
    contact: { given_name: "Alex", surname: "Client", email: "alex@example.com" },
    billing_contact: { given_name: "Pat", surname: "Billing", email: "pat@example.com" },
    rate_sheet: {
      id: "rate1",
      active: true,
      name: "Standard",
      revision: 3,
      effective_date: "2026-01-01",
    },
  },
});
test("complete unsigned draft is ready before internal approval, PDF and number", () => {
  const d = complete();
  assert.deepEqual(printMissing(d), []);
  assert.equal(d.attachment, undefined);
  assert.equal(d.request.number, undefined);
});
test("optional contacts, owner, references and addresses do not block signature readiness", () => {
  const d = complete();
  delete d.related.client.address;
  delete d.related.client.city;
  delete d.related.client.province_state;
  assert.deepEqual(printMissing(d), []);
  d.related.billing_contact = d.related.contact;
  assert.deepEqual(printMissing(d), []);
});
test("each required client value is independently checked", () => {
  const mutations: ((d: RequestDetail) => void)[] = [
    (d) => (d.related.client.name = ""),
    (d) => (d.request.description = " "),
    (d) => (d.request.location = "bad"),
    (d) => (d.related.contact.given_name = d.related.contact.surname = ""),
    (d) => (d.related.contact.email = "invalid"),
    (d) => (d.related.billing_contact.given_name = d.related.billing_contact.surname = ""),
    (d) => (d.related.billing_contact.email = "invalid"),
    (d) => (d.request.project_value = 0),
    (d) => (d.request.project_value = 1.5),
    (d) => (d.request.time_and_materials_set = false),
    (d) => delete d.related.rate_sheet,
  ];
  for (const mutate of mutations) {
    const d = complete();
    mutate(d);
    assert.equal(printMissing(d).length, 1);
  }
});
test("a contact marked as having no email address is ready for signature", () => {
  const d = complete();
  d.related.contact.email = "";
  assert.deepEqual(printMissing(d), ["Project contact email"]);
  d.related.contact.no_email = true;
  assert.deepEqual(printMissing(d), []);
  assert.equal(contactEmail(d.related.contact), "No email address");
  assert.equal(contactEmail(d.related.billing_contact), d.related.billing_contact.email);
});
test("T&M requires matching available billing rates", () => {
  const d = complete();
  d.request.time_and_materials = true;
  assert.deepEqual(printMissing(d), ["Client billing rates"]);
  d.related.rate_sheet_entries = [{ role_name: "Engineer", rate: 180, overtime_rate: 270 }];
  assert.deepEqual(printMissing(d), []);
  d.related.rate_sheet_entries[0].rate = undefined;
  assert.deepEqual(printMissing(d), ["Client billing rates"]);
});
test("complete proposed rates are ready for signature while the sheet is inactive", () => {
  const d = complete();
  d.related.rate_sheet = { id: "draft1", active: false, code: "2026-11-00025", revision: 0 };
  d.related.rate_sheet_proposed = true;
  d.related.rate_sheet_missing_roles = [];
  assert.deepEqual(printMissing(d), []);
  d.related.rate_sheet_missing_roles = ["Survey Drone Pilot"];
  assert.deepEqual(printMissing(d), ["Rate sheet"]);
  delete d.related.rate_sheet_missing_roles;
  assert.deepEqual(printMissing(d), ["Rate sheet"]);
});
test("a chosen rate sheet must be active for signature", () => {
  const d = complete();
  d.related.rate_sheet.active = false;
  assert.deepEqual(printMissing(d), ["Rate sheet"]);
  // A missing-roles list does not make an inactive chosen sheet ready.
  d.related.rate_sheet_missing_roles = [];
  assert.deepEqual(printMissing(d), ["Rate sheet"]);
});
test("billing address selects one complete source without mixing partial addresses", () => {
  const client = { address: "1 Main Street", city: "Thunder Bay", province_state: "ON" };
  assert.equal(billingAddress({}, client), "1 Main Street, Thunder Bay, ON");
  assert.equal(billingAddress({ city: "Ottawa" }, client), "Ottawa");
});
test("unknown and fractional values are not printed as zero or rounded", () => {
  for (const value of [undefined, null, "120000", 0, -1, 10.2])
    assert.equal(wholeCAD(value), "Not provided");
  assert.equal(wholeCAD(120000), "$120,000");
});
test("every checklist question starts unchecked on each review", () => {
  const questions = [
    { key: "pdf_sufficient", text: "The PDF is sufficient." },
    { key: "rate_sheet_correct", text: "The rate sheet is correct." },
  ];
  const first = emptyChecklist(questions);
  assert.deepEqual(first, { pdf_sufficient: false, rate_sheet_correct: false });
  first.pdf_sufficient = true;
  assert.equal(emptyChecklist(questions).pdf_sufficient, false);
});

test("creation queue navigation requires current server eligibility even when count is zero", () => {
  const counts = { "/job-requests?queue=creation": 0 };
  assert.equal(canViewJobCreateQueue(["job"], true, counts), true);
  assert.equal(canViewJobCreateQueue(["job"], true, {}), false);
  assert.equal(canViewJobCreateQueue(["job"], false, counts), false);
  assert.equal(canViewJobCreateQueue(["admin"], true, counts), false);
});

test("only an In Progress, Submitted or Awarded proposal accepts a request", () => {
  for (const status of ["In Progress", "Submitted", "Awarded"])
    assert.equal(proposalAcceptsRequest(status), true, status);
  for (const status of ["Not Awarded", "No Bid", "Cancelled", "Active", "", undefined])
    assert.equal(proposalAcceptsRequest(status), false, String(status));
});

test("the request page says who acts next, by name or as you", () => {
  const next = (status: string, viewer: Partial<Parameters<typeof requestNextStep>[0]> = {}) =>
    requestNextStep({
      status,
      creatorName: "Casey Creator",
      approverName: "Avery Approver",
      viewerIsCreator: false,
      viewerIsApprover: false,
      canCreate: false,
      ...viewer,
    });
  const creator = { viewerIsCreator: true };
  assert.equal(next("Draft"), "Casey Creator has not submitted this request.");
  assert.equal(next("Draft", creator), "You have not submitted this request.");
  assert.equal(next("Recalled"), "Recalled. Casey Creator has not submitted it again.");
  assert.equal(next("Recalled", creator), "Recalled. You have not submitted it again.");
  assert.equal(next("Awaiting approval"), "Waiting for approval by Avery Approver.");
  assert.equal(next("Awaiting approval", creator), "Waiting for approval by Avery Approver.");
  assert.equal(next("Awaiting approval", { viewerIsApprover: true }), "Waiting for your approval.");
  assert.equal(next("Awaiting creation"), "Approved. Waiting for procurement to create the job.");
  assert.equal(
    next("Awaiting creation", { canCreate: true }),
    "Approved. Waiting for you or another procurement user to create the job.",
  );
  assert.equal(
    next("Awaiting creation", { unawardedProposal: "P25-123" }),
    "Approved. Waiting for procurement to award proposal P25-123 and create the job.",
  );
  assert.equal(
    next("Awaiting creation", { canCreate: true, unawardedProposal: "P25-123" }),
    "Approved. Award proposal P25-123, then create the job.",
  );
  // The proposal matters only while the request awaits creation.
  assert.equal(
    next("Awaiting approval", { unawardedProposal: "P25-123" }),
    "Waiting for approval by Avery Approver.",
  );
  assert.equal(next("Rejected"), "Rejected. Waiting for Casey Creator to correct it.");
  assert.equal(next("Rejected", creator), "Rejected. Correct the request, then submit it again.");
  assert.equal(next("Unknown"), "");
});

test("a missing name is left out, never replaced by a role", () => {
  const next = (status: string) =>
    requestNextStep({
      status,
      creatorName: " ",
      approverName: "",
      viewerIsCreator: false,
      viewerIsApprover: false,
      canCreate: false,
    });
  assert.equal(next("Draft"), "Not submitted yet.");
  assert.equal(next("Recalled"), "Recalled. Not submitted again yet.");
  assert.equal(next("Awaiting approval"), "Waiting for approval.");
  assert.equal(next("Rejected"), "Rejected. Waiting for a correction.");
});

test("a request copies the source job's allocations with active divisions only", () => {
  const divisions = [
    { id: "eng", active: true },
    { id: "env", active: true },
    { id: "old", active: false },
  ];
  const source = [
    { id: "a1", division: { id: "eng", code: "ENG", name: "Engineering" }, hours: 40 },
    { id: "a2", division: { id: "old", code: "OLD", name: "Retired" }, hours: 10 },
    { id: "a3", division: { id: "gone", code: "GNE", name: "Unlisted" }, hours: 5 },
    { id: "a4", division: { id: "env", code: "ENV", name: "Environmental" }, hours: 0 },
  ];
  assert.deepEqual(copiedAllocations(source, divisions), [
    { division: "eng", hours: 40 },
    { division: "env", hours: 0 },
  ]);
  assert.deepEqual(copiedAllocations(undefined, divisions), []);
  assert.deepEqual(copiedAllocations(source, undefined), []);
});
test("rows without a division leave the allocations empty", () => {
  assert.equal(hasAllocations([]), false);
  assert.equal(hasAllocations([{ division: "", hours: 0 }]), false);
  assert.equal(
    hasAllocations([
      { division: "", hours: 0 },
      { division: "eng", hours: 0 },
    ]),
    true,
  );
});
