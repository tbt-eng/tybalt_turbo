import assert from "node:assert/strict";
import test from "node:test";
import {
  safeWorkflowReturn,
  keepClientWorkflowDraft,
  mergeClientWorkflowDraft,
  readClientWorkflowDraft,
  clearClientWorkflowDraft,
  selectClientWorkflowProfile,
  workflowDraftSnapshot,
  workflowFormName,
} from "../src/lib/clientWorkflowNavigation.ts";

// Mimic Storage: stored entries are its enumerable keys, as Object.keys(sessionStorage) expects.
function storageMock(storage: Map<string, string>, blocked = () => false) {
  const methods = {
    setItem: (key: string, value: string) => {
      if (blocked()) throw new Error("Storage is full");
      storage.set(key, value);
    },
    getItem: (key: string) => storage.get(key) ?? null,
    removeItem: (key: string) => storage.delete(key),
  };
  return new Proxy(methods, {
    ownKeys: () => [...storage.keys()],
    getOwnPropertyDescriptor: (_, key) =>
      storage.has(String(key))
        ? { configurable: true, enumerable: true, value: storage.get(String(key)) }
        : undefined,
  });
}

test("workflow returns accept only local project and request forms", () => {
  for (const path of [
    "/job-requests/new?parent=abc",
    "/job-requests/request1/edit",
    "/jobs/add",
    "/jobs/add/parent1",
    "/jobs/add/from/proposal1",
    "/jobs/job1/edit",
  ])
    assert.equal(safeWorkflowReturn(path), path);
  for (const path of [
    "https://example.com/jobs/add",
    "//example.com/jobs/add",
    "/\\example.com/jobs/add",
    "/clients/client1/details",
    "/jobs/job1/details#value",
    "/jobs/value-approvals",
    "/jobs",
    "/jobs/job1/delete",
    "/job-requests/request1/print",
    // A request's own page is read only, so a draft has no form there.
    "/job-requests/request1",
    "/job-requests/request1/edit/more",
    "/jobs/%2fexample.com/edit",
    "/jobs/add\n",
    null,
    undefined,
  ])
    assert.equal(safeWorkflowReturn(path), null);
});

test("workflow returns name the originating form", () => {
  assert.equal(workflowFormName("/job-requests/new?parent=abc"), "project request");
  assert.equal(workflowFormName("/job-requests/request1/edit"), "project request");
  assert.equal(workflowFormName("/jobs/add/parent1"), "job form");
  assert.equal(workflowFormName("/jobs/job1/edit"), "job form");
});

test("copied profiles change only the matching unfinished form and preserve its draft", () => {
  const storage = new Map<string, string>();
  const originalWindow = Object.getOwnPropertyDescriptor(globalThis, "window");
  const originalStorage = Object.getOwnPropertyDescriptor(globalThis, "sessionStorage");
  const location = { href: "https://app.example/jobs/job1/edit" };
  let storageBlocked = false;
  Object.defineProperty(globalThis, "window", { configurable: true, value: { location } });
  Object.defineProperty(globalThis, "sessionStorage", {
    configurable: true,
    value: storageMock(storage, () => storageBlocked),
  });
  try {
    for (const [source, kind, field] of [
      ["/job-requests/new?parent=parent1", "request:new", "item"],
      ["/job-requests/request1/edit", "request:request1", "item"],
      ["/jobs/add", "job:new", "item"],
      ["/jobs/add/parent1", "job:new", "item"],
      ["/jobs/add/from/proposal1", "job:new", "item"],
      ["/jobs/job1/edit", "job:job1", "item"],
    ]) {
      location.href = "https://app.example" + source;
      const documentFile = new File(["%PDF-1.4"], "authorization.pdf", { type: "application/pdf" });
      const item = {
        client: "client1",
        invoicing_information: "original",
        description: "Unsaved scope",
        version: 4,
      };
      const draft = {
        [field]: item,
        allocations: [{ division: "division1", hours: 12 }],
        editing: { id: "request1", version: 7 },
        originalAllocations: "original allocations",
        saved: "original saved baseline",
        documentFile,
        fee: "fixed",
      };
      const workspace = new URL(
        keepClientWorkflowDraft(kind, "user1", draft, "client1", "contact1"),
        "https://app.example",
      );
      const destination = workspace.searchParams.get("workflow_return")!;
      location.href = "https://app.example/clients/client1/invoicing/profile1/edit";
      assert.equal(selectClientWorkflowProfile(destination, "user2", "client1", "copy"), null);
      assert.equal(selectClientWorkflowProfile(destination, "user1", "client2", "copy"), null);
      assert.equal(selectClientWorkflowProfile(destination, "", "client1", "copy"), null);
      assert.equal(selectClientWorkflowProfile(destination, "user1", "client1", ""), null);
      const otherPath = new URL(destination, "https://app.example");
      otherPath.pathname = "/job-requests/other/edit";
      assert.equal(
        selectClientWorkflowProfile(
          otherPath.pathname + otherPath.search,
          "user1",
          "client1",
          "copy",
        ),
        null,
      );
      assert.equal(
        selectClientWorkflowProfile(
          "https://other.example" + destination,
          "user1",
          "client1",
          "copy",
        ),
        null,
      );
      assert.equal(
        selectClientWorkflowProfile("/job-requests/new", "user1", "client1", "copy"),
        null,
      );
      const oldNow = Date.now;
      Date.now = () => oldNow() + 5 * 60 * 60 * 1000;
      try {
        assert.equal(selectClientWorkflowProfile(destination, "user1", "client1", "copy"), null);
      } finally {
        Date.now = oldNow;
      }
      storageBlocked = true;
      assert.equal(selectClientWorkflowProfile(destination, "user1", "client1", "copy"), null);
      storageBlocked = false;
      location.href = "https://app.example" + destination;
      assert.deepEqual(
        readClientWorkflowDraft(kind, "user1"),
        draft,
        "failed selection leaves all draft data intact",
      );
      assert.equal(
        selectClientWorkflowProfile(destination, "user1", "client1", "copy"),
        destination,
      );
      const restored = readClientWorkflowDraft<typeof draft>(kind, "user1")!;
      assert.deepEqual(restored, {
        ...draft,
        [field]: { ...item, invoicing_information: "copy" },
        selectedProfile: "copy",
      });
      assert.equal(restored.documentFile, documentFile, "normal navigation keeps the selected PDF");
      assert.equal(item.invoicing_information, "original", "the saved baseline remains unchanged");
      const token = new URL(destination, "https://app.example").searchParams.get("client_draft")!;
      const reloadCopy = JSON.parse(storage.get("client-workspace-draft:" + token)!);
      assert.equal(reloadCopy.payload[field].invoicing_information, "copy");
      assert.equal(reloadCopy.payload.editing.version, 7);
      clearClientWorkflowDraft();
    }

    // A reloaded tab reads the JSON copy without requiring an in-memory draft.
    const token = crypto.randomUUID();
    const destination = `/job-requests/new?client_draft=${token}`;
    const envelope = {
      kind: "request:new",
      user: "user1",
      path: "/job-requests/new",
      created: Date.now(),
      payload: {
        item: { client: "client1", invoicing_information: "original" },
        saved: "baseline",
      },
    };
    storage.set(
      "client-workspace-draft:" + token,
      JSON.stringify({ ...envelope, kind: "job:new" }),
    );
    assert.equal(selectClientWorkflowProfile(destination, "user1", "client1", "copy"), null);
    storage.set("client-workspace-draft:" + token, JSON.stringify(envelope));
    assert.equal(selectClientWorkflowProfile(destination, "user1", "client1", "copy"), destination);
    location.href = "https://app.example" + destination;
    assert.deepEqual(readClientWorkflowDraft("request:new", "user1"), {
      item: { client: "client1", invoicing_information: "copy" },
      saved: "baseline",
      selectedProfile: "copy",
    });
    clearClientWorkflowDraft();
  } finally {
    if (originalWindow) Object.defineProperty(globalThis, "window", originalWindow);
    else Reflect.deleteProperty(globalThis, "window");
    if (originalStorage) Object.defineProperty(globalThis, "sessionStorage", originalStorage);
    else Reflect.deleteProperty(globalThis, "sessionStorage");
  }
});

test("drafts preserve values, require the same user and source, and clear after save", () => {
  const storage = new Map<string, string>();
  const originalWindow = Object.getOwnPropertyDescriptor(globalThis, "window");
  const originalStorage = Object.getOwnPropertyDescriptor(globalThis, "sessionStorage");
  const location = { href: "https://app.example/job-requests/new?parent=parent1" };
  Object.defineProperty(globalThis, "window", { configurable: true, value: { location } });
  Object.defineProperty(globalThis, "sessionStorage", {
    configurable: true,
    value: storageMock(storage),
  });
  try {
    const draft = {
      item: { description: "Unsaved scope", version: 4 },
      allocations: [{ hours: 12 }],
    };
    let sourceEntry = "";
    const href = keepClientWorkflowDraft(
      "request:new",
      "user1",
      draft,
      "client1",
      "contact1",
      "workspace",
      (url) => {
        sourceEntry = url;
      },
    );
    const workspace = new URL(href, "https://app.example");
    assert.equal(workspace.pathname, "/clients/client1/details");
    assert.equal(workspace.searchParams.get("project_contact"), "contact1");
    const returnUrl = workspace.searchParams.get("workflow_return")!;
    assert.equal(
      sourceEntry,
      returnUrl,
      "browser Back and the explicit return use the same draft token",
    );
    assert.ok(safeWorkflowReturn(returnUrl));
    assert.match(returnUrl, /parent=parent1/);
    location.href = "https://app.example" + returnUrl;
    assert.deepEqual(readClientWorkflowDraft("request:new", "user1"), draft);
    assert.equal(readClientWorkflowDraft("request:new", "user2"), null);
    assert.equal(readClientWorkflowDraft("job:job1", "user1"), null);
    location.href = location.href.replace("/job-requests/new", "/job-requests/other");
    assert.equal(readClientWorkflowDraft("request:new", "user1"), null);
    location.href = "https://app.example" + returnUrl;
    const oldNow = Date.now;
    Date.now = () => oldNow() + 5 * 60 * 60 * 1000;
    try {
      assert.equal(readClientWorkflowDraft("request:new", "user1"), null);
    } finally {
      Date.now = oldNow;
    }
    const second = new URL(
      keepClientWorkflowDraft("request:new", "user1", draft, "client1", ""),
      "https://app.example",
    );
    assert.equal(
      readClientWorkflowDraft("request:new", "user1"),
      null,
      "the earlier history entry cannot restore a superseded draft",
    );
    location.href = "https://app.example" + second.searchParams.get("workflow_return");
    assert.deepEqual(readClientWorkflowDraft("request:new", "user1"), draft);
    clearClientWorkflowDraft();
    assert.equal(readClientWorkflowDraft("request:new", "user1"), null);
    assert.equal(storage.size, 0);

    const addUrl = new URL(
      keepClientWorkflowDraft("request:new", "user1", draft, "client1", "", "contact"),
      "https://app.example",
    );
    assert.equal(addUrl.pathname, "/clients/client1/contacts/add");
    assert.match(addUrl.searchParams.get("return_to")!, /^\/clients\/client1\/details\?/);
    const profileUrl = new URL(
      keepClientWorkflowDraft("request:new", "user1", draft, "client1", "contact1", "invoicing"),
      "https://app.example",
    );
    assert.equal(profileUrl.pathname, "/clients/client1/invoicing/add");
    assert.equal(profileUrl.searchParams.get("project_contact"), "contact1");
    const profileReturn = new URL(profileUrl.searchParams.get("return_to")!, "https://app.example");
    assert.equal(profileReturn.searchParams.get("project_contact"), "contact1");
    location.href = "https://app.example" + profileReturn.searchParams.get("workflow_return");
    assert.deepEqual(readClientWorkflowDraft("request:new", "user1"), draft);
  } finally {
    if (originalWindow) Object.defineProperty(globalThis, "window", originalWindow);
    else Reflect.deleteProperty(globalThis, "window");
    if (originalStorage) Object.defineProperty(globalThis, "sessionStorage", originalStorage);
    else Reflect.deleteProperty(globalThis, "sessionStorage");
  }
});

test("a saved request's draft leaves from its edit page and returns there", () => {
  const storage = new Map<string, string>();
  const originalWindow = Object.getOwnPropertyDescriptor(globalThis, "window");
  const originalStorage = Object.getOwnPropertyDescriptor(globalThis, "sessionStorage");
  const location = { href: "https://app.example/job-requests/request1/edit" };
  Object.defineProperty(globalThis, "window", { configurable: true, value: { location } });
  Object.defineProperty(globalThis, "sessionStorage", {
    configurable: true,
    value: storageMock(storage),
  });
  try {
    const draft = {
      item: { id: "request1", client: "client1", description: "Unsaved scope", version: 4 },
      allocations: [{ division: "division1", hours: 12 }],
      saved: "baseline",
    };
    const profileUrl = new URL(
      keepClientWorkflowDraft("request:request1", "user1", draft, "client1", "", "invoicing"),
      "https://app.example",
    );
    const workspace = new URL(profileUrl.searchParams.get("return_to")!, "https://app.example");
    const workflowReturn = workspace.searchParams.get("workflow_return")!;
    const back = new URL(workflowReturn, "https://app.example");
    assert.equal(back.pathname, "/job-requests/request1/edit");
    assert.equal(safeWorkflowReturn(workflowReturn), workflowReturn);
    assert.equal(workflowFormName(workflowReturn), "project request");
    // A profile created in the workspace goes into the request on its edit page.
    location.href = "https://app.example/clients/client1/invoicing/add";
    assert.equal(
      selectClientWorkflowProfile(workflowReturn, "user1", "client1", "profile1"),
      workflowReturn,
    );
    location.href = "https://app.example" + workflowReturn;
    assert.deepEqual(readClientWorkflowDraft("request:request1", "user1"), {
      ...draft,
      item: { ...draft.item, invoicing_information: "profile1" },
      selectedProfile: "profile1",
    });
    // The request's read-only page does not restore the draft.
    location.href = `https://app.example/job-requests/request1${back.search}`;
    assert.equal(readClientWorkflowDraft("request:request1", "user1"), null);
    // Nor can a form there keep one.
    assert.throws(() => keepClientWorkflowDraft("request:request1", "user1", draft, "client1", ""));
  } finally {
    if (originalWindow) Object.defineProperty(globalThis, "window", originalWindow);
    else Reflect.deleteProperty(globalThis, "window");
    if (originalStorage) Object.defineProperty(globalThis, "sessionStorage", originalStorage);
    else Reflect.deleteProperty(globalThis, "sessionStorage");
  }
});

// A request form as kept on leaving for the client workspace, and as loaded
// again on return. Empty proposed-rate rows are not compared, as in the editor.
type RequestForm = {
  item: Record<string, unknown>;
  allocations: { division: string; hours: number }[];
  fee: string;
  rateMode: string;
  proposal: { code: string; rows: { role: string; rate: number | null }[] };
};
const ratedRows = (part: string, value: unknown) =>
  part === "proposal"
    ? {
        ...(value as RequestForm["proposal"]),
        rows: (value as RequestForm["proposal"]).rows.filter((row) => row.rate !== null),
      }
    : value;
const savedForm = (): RequestForm => ({
  item: {
    id: "request1",
    version: 3,
    description: "Bridge survey",
    client: "client1",
    contact: "contact1",
    notes: "",
  },
  allocations: [{ division: "division1", hours: 10 }],
  fee: "fixed",
  rateMode: "propose",
  proposal: { code: "R1", rows: [{ role: "role1", rate: 100 }] },
});
const mergeForm = (draft: RequestForm, baseline: RequestForm, loaded: RequestForm) =>
  mergeClientWorkflowDraft(
    draft,
    workflowDraftSnapshot(baseline, ratedRows),
    loaded,
    ["item"],
    ratedRows,
  );

test("a kept request draft takes what an absorb changed and keeps the user's changes", () => {
  const baseline = savedForm();
  const draft = savedForm();
  draft.item.description = "Bridge survey, phase 2";
  const loaded = savedForm();
  // An absorb rewrote the contact and bumped the version meanwhile.
  loaded.item.contact = "contact2";
  loaded.item.version = 4;
  loaded.item.notes = "Saved elsewhere";
  const merged = mergeForm(draft, baseline, loaded);
  assert.equal(merged.item.version, 4);
  assert.equal(merged.item.contact, "contact2");
  assert.equal(merged.item.notes, "Saved elsewhere");
  assert.equal(merged.item.description, "Bridge survey, phase 2");
});

test("a field both the user and the server changed keeps the user's value", () => {
  const draft = savedForm();
  draft.item.contact = "contact3";
  draft.item.version = 3;
  const loaded = savedForm();
  loaded.item.contact = "contact2";
  loaded.item.version = 4;
  const merged = mergeForm(draft, savedForm(), loaded);
  assert.equal(merged.item.contact, "contact3");
  assert.equal(merged.item.version, 4);
});

test("a profile chosen in the client workspace survives the merge", () => {
  const draft = savedForm();
  // selectClientWorkflowProfile writes the choice into the kept item.
  draft.item.invoicing_information = "profile1";
  const merged = mergeForm(draft, savedForm(), savedForm());
  assert.equal(merged.item.invoicing_information, "profile1");
});

test("allocations, proposed rates, fee and rate mode merge as whole units", () => {
  const draft = savedForm();
  draft.allocations = [
    { division: "division1", hours: 10 },
    { division: "division2", hours: 5 },
  ];
  // An empty row the rates editor added is not a change.
  draft.proposal.rows.push({ role: "role2", rate: null });
  const loaded = savedForm();
  loaded.allocations = [{ division: "division1", hours: 20 }];
  loaded.proposal = { code: "R2", rows: [{ role: "role1", rate: 120 }] };
  loaded.fee = "tm";
  const merged = mergeForm(draft, savedForm(), loaded);
  assert.deepEqual(merged.allocations, draft.allocations);
  assert.deepEqual(merged.proposal, loaded.proposal);
  assert.equal(merged.fee, "tm");
  assert.equal(merged.rateMode, "propose");

  const changedRates = savedForm();
  changedRates.proposal.rows[0].rate = 110;
  changedRates.fee = "tm";
  changedRates.rateMode = "choose";
  const kept = mergeForm(changedRates, savedForm(), loaded);
  assert.deepEqual(kept.proposal, changedRates.proposal);
  assert.deepEqual(kept.allocations, loaded.allocations);
  assert.equal(kept.fee, "tm");
  assert.equal(kept.rateMode, "choose");
});

test("a merged request is unsaved only when it differs from the server copy", () => {
  const loaded = savedForm();
  loaded.item.contact = "contact2";
  loaded.item.version = 4;
  const baselineAfter = workflowDraftSnapshot(loaded, ratedRows);
  const untouched = mergeForm(savedForm(), savedForm(), loaded);
  assert.equal(workflowDraftSnapshot(untouched, ratedRows), baselineAfter);

  const draft = savedForm();
  draft.item.notes = "Call first";
  const changed = mergeForm(draft, savedForm(), loaded);
  assert.notEqual(workflowDraftSnapshot(changed, ratedRows), baselineAfter);
  changed.item.notes = "";
  assert.equal(workflowDraftSnapshot(changed, ratedRows), baselineAfter);
});

test("a new request's draft restores as kept", () => {
  const blank: RequestForm = {
    item: { description: "", client: "", contact: "" },
    allocations: [],
    fee: "",
    rateMode: "choose",
    proposal: { code: "", rows: [] },
  };
  const draft: RequestForm = {
    item: { description: "New survey", client: "client1", contact: "contact1" },
    allocations: [{ division: "division1", hours: 8 }],
    fee: "tm",
    rateMode: "choose",
    proposal: { code: "", rows: [] },
  };
  assert.deepEqual(mergeForm(draft, blank, structuredClone(blank)), draft);
});

test("a draft with an unreadable baseline keeps all of its values but the version", () => {
  const draft = savedForm();
  draft.item.contact = "contact1";
  const loaded = savedForm();
  loaded.item.contact = "contact2";
  loaded.item.version = 4;
  const merged = mergeClientWorkflowDraft(draft, "not json", loaded, ["item"], ratedRows);
  assert.equal(merged.item.contact, "contact1");
  assert.equal(merged.item.version, 4);
});
