import type { PurchaseOrdersRecord } from "$lib/pocketbase-types";
import {
  PurchaseOrdersFrequencyOptions,
  PurchaseOrdersPaymentTypeOptions,
  PurchaseOrdersStatusOptions,
  PurchaseOrdersTypeOptions,
} from "$lib/pocketbase-types";
import type { PageLoad } from "./$types";
import type { PurchaseOrdersPageData } from "$lib/svelte-types";
import { localIsoDate } from "$lib/localDate";

export const load: PageLoad<PurchaseOrdersPageData> = async () => {
  const defaultItem: Partial<PurchaseOrdersRecord> = {
    _imported: false,
    legacy_manual_entry: true,
    po_number: "",
    status: PurchaseOrdersStatusOptions.Active,
    uid: "",
    type: PurchaseOrdersTypeOptions["One-Time"],
    date: localIsoDate(),
    end_date: "",
    frequency: PurchaseOrdersFrequencyOptions.Weekly,
    division: "",
    description: "",
    total: 0,
    covered_within_project_budget: false,
    payment_type: PurchaseOrdersPaymentTypeOptions.OnAccount,
    vendor: "",
    job: "",
    category: "",
    kind: "",
    branch: "",
    approver: "",
    priority_second_approver: "",
  };

  return {
    item: { ...defaultItem } as PurchaseOrdersRecord,
    editing: false,
    id: null,
  };
};
