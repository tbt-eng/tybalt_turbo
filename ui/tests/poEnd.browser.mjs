import assert from "node:assert/strict";
import { runBrowserHarness } from "./browserHarness.mjs";

// Closing and cancelling a purchase order are final. The dialog says so, sends
// nothing until the person confirms, and shows the server's refusal. The owner
// of a purchase order and payables staff get the buttons that open it.
await runBrowserHarness(async (page, origin) => {
  const errors = [];
  page.on("pageerror", (error) => errors.push(error.message));
  const sent = [];
  let refusal = null;
  await page.route("**/api/**", (route) => {
    const request = route.request();
    const path = new URL(request.url()).pathname;
    const send = (body, status = 200) =>
      route.fulfill({ status, contentType: "application/json", body: JSON.stringify(body) });
    if (path.startsWith("/api/purchase_orders/po1/") && request.method() === "POST") {
      sent.push(path);
      return refusal ? send(refusal, 400) : send({ id: "po1" });
    }
    return send({});
  });
  const dialog = () => page.getByRole("dialog");

  // Cancel: the dialog explains it, and "Keep it" sends nothing.
  await page.goto(`${origin}/tests/fixtures/poEnd.html?action=cancel`);
  await dialog().getByText("Cancel PO 2611-0379?").waitFor();
  await dialog().getByText("This is permanent").waitFor();
  await dialog().getByText("You can't undo this.", { exact: false }).waitFor();
  await dialog().getByText("its number can't be used again", { exact: false }).waitFor();
  await dialog().getByRole("button", { name: "Keep it", exact: true }).click();
  await page.waitForFunction(() => !document.querySelector("dialog[open]"));
  assert.deepEqual(sent, []);
  assert.equal(await page.evaluate(() => window.poEnded), undefined);

  // Confirming sends the cancel, then tells the page.
  await page.goto(`${origin}/tests/fixtures/poEnd.html?action=cancel`);
  await dialog().getByRole("button", { name: "Cancel Purchase Order", exact: true }).click();
  await page.waitForFunction(() => window.poEnded === true);
  assert.deepEqual(sent, ["/api/purchase_orders/po1/cancel"]);

  // Close: a refusal stays in the dialog, which stays open.
  refusal = {
    code: "pending_expenses",
    message: "This PO has one or more approved expenses awaiting commitment.",
  };
  await page.goto(`${origin}/tests/fixtures/poEnd.html?action=close`);
  await dialog().getByText("Close PO 2611-0379?").waitFor();
  await dialog().getByText("This is permanent").waitFor();
  await dialog().getByRole("button", { name: "Close Purchase Order", exact: true }).click();
  await dialog().getByText(refusal.message).waitFor();
  assert.equal(sent.at(-1), "/api/purchase_orders/po1/close");
  assert.equal(await page.evaluate(() => window.poEnded), undefined);

  // Who gets the buttons: the owner or payables staff, on an Active PO.
  const rules = await page.evaluate(() => {
    const { poMayBeCancelledBy, poMayBeClosedBy } = window.poRules;
    const po = (over) => ({
      status: "Active",
      type: "Cumulative",
      uid: "owner",
      committed_expenses_count: 0,
      ...over,
    });
    return {
      ownerCancels: poMayBeCancelledBy(po(), [], "owner"),
      payablesCancels: poMayBeCancelledBy(po(), ["payables_admin"], "other"),
      otherCancels: poMayBeCancelledBy(po(), ["time"], "other"),
      cancelWithExpense: poMayBeCancelledBy(po({ committed_expenses_count: 1 }), [], "owner"),
      cancelClosed: poMayBeCancelledBy(po({ status: "Closed" }), [], "owner"),
      ownerCloses: poMayBeClosedBy(po({ committed_expenses_count: 1 }), [], "owner"),
      otherCloses: poMayBeClosedBy(po({ committed_expenses_count: 1 }), [], "other"),
      closeNoExpense: poMayBeClosedBy(po(), [], "owner"),
      closeOneTime: poMayBeClosedBy(
        po({ type: "One-Time", committed_expenses_count: 1 }),
        ["payables_admin"],
        "owner",
      ),
    };
  });
  assert.deepEqual(rules, {
    ownerCancels: true,
    payablesCancels: true,
    otherCancels: false,
    cancelWithExpense: false,
    cancelClosed: false,
    ownerCloses: true,
    otherCloses: false,
    closeNoExpense: false,
    closeOneTime: false,
  });
  assert.deepEqual(errors, []);
});
console.log("Close and cancel purchase order browser checks passed.");
