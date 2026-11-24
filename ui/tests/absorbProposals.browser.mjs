import assert from "node:assert/strict";
import { runBrowserHarness } from "./browserHarness.mjs";

// Without the absorb permission the absorb editor proposes: it needs a
// reason and sends the proposal. The queue lets a requester withdraw their own
// proposal, and lets an absorb holder review or decline. Review opens the
// editor with the proposal's records, where the holder may keep another
// record before absorbing. On the Decided tab a requester deletes their
// withdrawn proposal after a confirmation.
await runBrowserHarness(
  async (page, origin) => {
    const errors = [];
    page.on("pageerror", (error) => errors.push(error.message));
    const sent = [];
    // The next DELETE answers with this refusal, once.
    let refuseDelete = null;
    await page.route("**/api/**", (route) => {
      const request = route.request();
      const path = new URL(request.url()).pathname;
      const send = (body, status = 200) =>
        route.fulfill({ status, contentType: "application/json", body: JSON.stringify(body) });
      if (request.method() === "DELETE") {
        sent.push({ path, method: "DELETE" });
        const refusal = refuseDelete;
        refuseDelete = null;
        return refusal ? send(refusal, 400) : send({ message: "absorb proposal deleted" });
      }
      if (request.method() === "POST") {
        sent.push({ path, body: request.postDataJSON() });
        return send({ id: "new" });
      }
      if (path.startsWith("/api/collections/absorb_actions"))
        return send({ items: [], page: 1, perPage: 1, totalItems: 0, totalPages: 0 });
      return send({});
    });

    // Propose: no absorb action is read, and Propose waits for a reason.
    await page.goto(`${origin}/tests/fixtures/absorbProposals.html?view=propose`);
    await page.getByText("Target Record", { exact: true }).waitFor();
    assert.equal(await page.getByRole("button", { name: "Absorb", exact: true }).count(), 0);
    await page.getByLabel("Select Record").selectOption("clientb00000001");
    const propose = page.getByRole("button", { name: "Propose", exact: true });
    assert.equal(await propose.isDisabled(), true);
    await page.getByLabel("Reason").fill("Acme was entered twice");
    await propose.click();
    await page.waitForURL("**/absorb/proposals");
    assert.deepEqual(sent, [
      {
        path: "/api/absorb_proposals",
        body: {
          collection_name: "clients",
          target_id: "clienta00000001",
          record_ids: ["clientb00000001"],
          reason: "Acme was entered twice",
        },
      },
    ]);

    // Queue for a requester: Withdraw on their own proposal only.
    sent.length = 0;
    await page.goto(`${origin}/tests/fixtures/absorbProposals.html?view=queue`);
    await page.getByText("Robin Requester", { exact: false }).first().waitFor();
    assert.equal(await page.getByTitle("Withdraw", { exact: true }).count(), 1);
    assert.equal(await page.getByTitle("Review", { exact: true }).count(), 0);
    assert.equal(await page.getByTitle("Decline", { exact: true }).count(), 0);
    await page.getByTitle("Withdraw", { exact: true }).click();
    await page.waitForFunction(() => window.reloads === 1);
    assert.deepEqual(
      sent.map((s) => s.path),
      ["/api/absorb_proposals/proposal0000001/withdraw"],
    );

    // Queue for a holder: Review links to the absorb page; Decline needs a reason.
    sent.length = 0;
    await page.goto(`${origin}/tests/fixtures/absorbProposals.html?view=queue&claims=absorb`);
    await page.getByText("Sam Other", { exact: false }).first().waitFor();
    assert.equal(await page.getByTitle("Review", { exact: true }).count(), 2);
    await page.getByTitle("Decline", { exact: true }).nth(1).click();
    const dialog = page.getByRole("dialog");
    const submit = dialog.getByRole("button", { name: "Decline", exact: true });
    assert.equal(await submit.isDisabled(), true);
    await dialog.getByLabel("Why").fill("Different companies");
    await submit.click();
    await page.waitForFunction(() => window.reloads === 1);
    assert.deepEqual(sent, [
      {
        path: "/api/absorb_proposals/proposal0000002/decline",
        body: { comment: "Different companies" },
      },
    ]);
    await page.getByTitle("Review", { exact: true }).first().click();
    await page.waitForURL("**/clients/clienta00000001/absorb?proposal=proposal0000001");

    // Review: the proposal's records are filled in, with no search. Keeping
    // Acme Limited swaps it with the target.
    sent.length = 0;
    await page.goto(`${origin}/tests/fixtures/absorbProposals.html?view=review`);
    await page.getByText("Proposed by Robin Requester", { exact: true }).waitFor();
    assert.equal(await page.getByLabel("Select Record").count(), 0);
    await page.getByTitle("Keep this record instead", { exact: true }).click();
    await page.getByRole("button", { name: "Absorb", exact: true }).click();
    await page.waitForURL("**/absorb/actions");
    assert.deepEqual(sent, [
      {
        path: "/api/clients/clientb00000001/absorb",
        body: { ids_to_absorb: ["clienta00000001"], proposal: "proposal0000001" },
      },
    ]);

    // Decided tab: Delete only on the viewer's own withdrawn proposal, for a
    // holder too. Nothing is sent until the person confirms.
    for (const claims of ["", "&claims=absorb"]) {
      await page.goto(`${origin}/tests/fixtures/absorbProposals.html?view=decided${claims}`);
      await page.getByText("Sam Other", { exact: false }).first().waitFor();
      assert.equal(await page.getByTitle("Delete", { exact: true }).count(), 1);
      assert.equal(await page.getByTitle("Withdraw", { exact: true }).count(), 0);
    }
    sent.length = 0;
    const deleteButton = page.getByTitle("Delete", { exact: true });
    await deleteButton.click();
    const confirm = page.getByRole("dialog");
    await confirm.getByText("You can't undo this.").waitFor();
    await confirm.getByRole("button", { name: "Keep it", exact: true }).click();
    await confirm.waitFor({ state: "hidden" });
    assert.deepEqual(sent, []);

    // A refusal shows in the dialog, and the row stays.
    refuseDelete = {
      code: "proposal_not_withdrawn",
      message: "only a withdrawn absorb proposal can be deleted",
    };
    await deleteButton.click();
    await confirm.getByRole("button", { name: "Delete", exact: true }).click();
    await confirm.getByText("only a withdrawn absorb proposal can be deleted").waitFor();
    assert.equal(await deleteButton.count(), 1);
    await confirm.getByRole("button", { name: "Delete", exact: true }).click();
    await confirm.waitFor({ state: "hidden" });
    await page.waitForFunction(() => window.reloads === 1);
    assert.deepEqual(sent, [
      { path: "/api/absorb_proposals/proposal0000003", method: "DELETE" },
      { path: "/api/absorb_proposals/proposal0000003", method: "DELETE" },
    ]);
    // The row leaves the list; the others stay.
    assert.equal(await deleteButton.count(), 0);
    assert.equal(await page.getByText("Different companies", { exact: false }).count(), 1);
    assert.equal((await page.getByText("Sam Other", { exact: false }).count()) > 0, true);

    assert.deepEqual(errors, []);
  },
  {
    services: {
      "$app/stores":
        "import { readable } from 'svelte/store'; export const page = readable({ params: {}, url: new URL(location.href) });",
    },
  },
);
console.log("Absorb proposal browser checks passed.");
