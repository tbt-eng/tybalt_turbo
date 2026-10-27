import type { ExpensesRecord } from "$lib/pocketbase-types";
import { ExpensesPaymentTypeOptions } from "$lib/pocketbase-types";
import type { PageLoad } from "./$types";
import type { ExpensesPageData } from "$lib/svelte-types";
import { localIsoDate } from "$lib/localDate";

export const load: PageLoad<ExpensesPageData> = async () => {
  const defaultItem: Partial<ExpensesRecord> = {
    uid: "",
    // date in YYYY-MM-DD format
    date: localIsoDate(),
    division: "",
    description: "",
    payment_type: ExpensesPaymentTypeOptions.OnAccount,
    purchase_order: "",
    currency: "",
    vendor: "",
    job: "",
    category: "",
    kind: "",
    settled_total: 0,
    settled: "",
    settler: "",
    allowance_types: [],
    // approver is configured as not required in pocketbase so we do not have to
    // set it here, but is set by the server side hook
  };
  return { item: { ...defaultItem } as ExpensesRecord, editing: false, id: null };
};
