// Why a printed PO can show less than the approved amount.
//
// The approved amount of a PO is a ceiling, not the price. A person often asks
// for approval of more than the vendor's quote, so that a price rise, tax or
// shipping does not force a second PO and a second approval.
//
// The vendor must not see that ceiling. A printed PO that shows the full
// approved amount tells the vendor how much room there is. So the person can
// print a lower amount, usually the quoted price. The print page shows only
// that amount: never the print limit, the allowance or the warning. The print
// link holds only an opaque key, so a browser print footer cannot show the
// limit either (poPrintStorageKey).
//
// Plus tax and Plus shipping let the person authorize the quoted price and tell
// the vendor that tax or shipping come on top. Those extras still come out of
// the ceiling. Their allowance is the print limit less the printed amount, and
// tax and shipping share it. So a plus box needs a printed amount below the
// limit (validatePOPrintOptions), and the person must acknowledge the warning
// that states the allowance (poPrintWarning). Printing does not change the PO
// or its approved amount, so an extra cost above the allowance needs a new
// approved PO.
//
// The print limit (print_max_amount) comes from the server: the PO total (for a
// Recurring PO, the amount per period), or for a Cumulative PO the total less
// committed expenses. See app/routes/po_visibility_base.sql.
//
// Keep the popup and print-page checks consistent. Amounts use integer cents.
export type PrintablePO = {
  status: string;
  type: string;
  currency_code: string;
  frequency?: string;
  print_max_amount: number;
};

export type POPrintOptions = {
  amount: string;
  plusTax: boolean;
  plusShipping: boolean;
};

export function amountInCents(amount: string): number | null {
  if (!/^\d+(?:\.\d{1,2})?$/.test(amount)) return null;
  const cents = Math.round(Number(amount) * 100);
  return Number.isSafeInteger(cents) ? cents : null;
}

export function printLimitInCents(po: PrintablePO): number | null {
  if (!Number.isFinite(po.print_max_amount)) return null;
  return amountInCents(po.print_max_amount.toFixed(2));
}

export function defaultPOPrintOptions(po: PrintablePO): POPrintOptions {
  const limit = printLimitInCents(po);
  return {
    amount: limit !== null && limit > 0 ? (limit / 100).toFixed(2) : "",
    plusTax: false,
    plusShipping: false,
  };
}

export function validatePOPrintOptions(po: PrintablePO, options: POPrintOptions): string | null {
  if (po.status !== "Active") return "Only Active purchase orders can be printed.";
  const limit = printLimitInCents(po);
  if (limit === null || limit <= 0) return "No approved balance remains for printing.";
  const amount = amountInCents(options.amount);
  if (amount === null || amount <= 0)
    return "Enter an amount greater than zero with up to two decimal places.";
  if (amount > limit) return "The amount to print exceeds the available approved amount.";
  if ((options.plusTax || options.plusShipping) && amount === limit) {
    return "Enter a lower amount to allow for the additional charges.";
  }
  return null;
}

export function printMoney(amount: number, currency: string): string {
  return new Intl.NumberFormat("en-CA", {
    style: "currency",
    currency: currency || "CAD",
    currencyDisplay: "code",
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(amount);
}

export function authorizedAmountText(po: PrintablePO, options: POPrintOptions): string {
  const amount = amountInCents(options.amount);
  if (amount === null || amount <= 0) return "Enter a valid amount.";
  const charges =
    options.plusTax && options.plusShipping
      ? " plus tax and shipping"
      : options.plusTax
        ? " plus tax"
        : options.plusShipping
          ? " plus shipping"
          : "";
  const period = po.type === "Recurring" ? ` / ${po.frequency?.trim() || "period"}` : "";
  return `${printMoney(amount / 100, po.currency_code)}${charges}${period}`;
}

export function poPrintWarning(po: PrintablePO, options: POPrintOptions): string | null {
  if (!(options.plusTax || options.plusShipping) || validatePOPrintOptions(po, options))
    return null;
  const allowance = (printLimitInCents(po)! - amountInCents(options.amount)!) / 100;
  const charges =
    options.plusTax && options.plusShipping
      ? "Tax and shipping combined"
      : options.plusTax
        ? "Tax"
        : "Shipping";
  const period = po.type === "Recurring" ? " per period" : "";
  return `${charges} must not exceed ${printMoney(allowance, po.currency_code)}${period}. If these charges exceed this amount, you must have a new PO approved before authorizing the extra cost.`;
}

// Bind acknowledgement to the displayed amount, terms, currency, and limit.
export function poPrintAcknowledgementKey(po: PrintablePO, options: POPrintOptions): string {
  return JSON.stringify([
    options.amount,
    options.plusTax,
    options.plusShipping,
    po.print_max_amount,
    po.currency_code,
    po.type,
    po.frequency,
  ]);
}

export function poPrintSearchParams(po: PrintablePO, options: POPrintOptions): URLSearchParams {
  return new URLSearchParams({
    amount: options.amount,
    plusTax: options.plusTax ? "1" : "0",
    plusShipping: options.plusShipping ? "1" : "0",
    // The print page checks this limit against a fresh PO read.
    limit: (printLimitInCents(po)! / 100).toFixed(2),
    acknowledged: options.plusTax || options.plusShipping ? "1" : "0",
  });
}

export function readPOPrintOptions(po: PrintablePO, params: URLSearchParams): POPrintOptions {
  const keys = ["amount", "plusTax", "plusShipping", "limit", "acknowledged"];
  let options = defaultPOPrintOptions(po);
  if (keys.some((key) => params.has(key))) {
    if (
      keys.some((key) => params.getAll(key).length !== 1) ||
      ["plusTax", "plusShipping", "acknowledged"].some((key) => !/^[01]$/.test(params.get(key)!))
    ) {
      throw new Error("Invalid print options. Open Print options again.");
    }
    options = {
      amount: params.get("amount")!,
      plusTax: params.get("plusTax") === "1",
      plusShipping: params.get("plusShipping") === "1",
    };
    const limit = amountInCents(params.get("limit")!);
    if (limit === null) throw new Error("Invalid print limit. Open Print options again.");
    if (options.plusTax || options.plusShipping) {
      if (params.get("acknowledged") !== "1")
        throw new Error("Acknowledge the additional-charge warning in Print options.");
      if (limit !== printLimitInCents(po))
        throw new Error(
          "The available amount has changed. Open Print options again and review the warning.",
        );
    }
  }
  const error = validatePOPrintOptions(po, options);
  if (error) throw new Error(error);
  return options;
}

// The URL contains only an opaque key. Browser print footers must not expose the limit.
export function poPrintStorageKey(id: string, token: string): string {
  return `po-print:${id}:${token}`;
}

export function readPOPrintRequest(
  po: PrintablePO,
  id: string,
  params: URLSearchParams,
  storage: Pick<Storage, "getItem">,
): POPrintOptions {
  if (
    ["amount", "plusTax", "plusShipping", "limit", "acknowledged"].some((key) => params.has(key))
  ) {
    throw new Error("Invalid print link. Open Print options again.");
  }
  if (!params.has("print")) return readPOPrintOptions(po, new URLSearchParams());
  const token = params.get("print")!;
  if (params.getAll("print").length !== 1 || !/^[0-9a-f-]{36}$/.test(token)) {
    throw new Error("Invalid print link. Open Print options again.");
  }
  const saved = storage.getItem(poPrintStorageKey(id, token));
  if (!saved)
    throw new Error("Print options are not available in this tab. Open Print options again.");
  return readPOPrintOptions(po, new URLSearchParams(saved));
}
