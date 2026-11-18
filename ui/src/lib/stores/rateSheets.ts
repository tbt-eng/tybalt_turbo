import { createCollectionStore } from "./collectionStore";
import { pb } from "$lib/pocketbase";
import type { BaseSystemFields } from "$lib/pocketbase-types";
import { rateSheetSearchFields } from "$lib/rateSheetLabel";

export interface RateSheetResponse extends BaseSystemFields {
  // The sheet's label and identity, shared by its revisions.
  code: string;
  // An optional description.
  name: string;
  effective_date: string;
  revision: number;
  active: boolean;
  job_count: number;
  // The job numbers and the client names of every revision of the code, now
  // or in a job's rate sheet history, comma separated, for the search.
  job_numbers: string;
  client_names: string;
}

// Fetch from the augmented view which includes job_count
const fetchAllRateSheets = async (): Promise<RateSheetResponse[]> =>
  pb.collection("rate_sheets_augmented").getFullList({ sort: "-effective_date" });

export const rateSheets = createCollectionStore<RateSheetResponse>(
  "rate_sheets", // Subscribe to the base collection for real-time updates
  { requestKey: "rate_sheets" },
  {
    // A search finds a sheet by its code, its name, its effective date, or a
    // job or client that used the code.
    fields: rateSheetSearchFields,
    storeFields: ["id", "code", "name", "effective_date", "revision", "active", "job_count"],
  },
  // onCreate handler - refetch from augmented view
  async (item) => {
    const record = await pb.collection("rate_sheets_augmented").getOne<RateSheetResponse>(item.id);
    rateSheets.update((s) => ({
      ...s,
      items: [...s.items, record],
      index: s.index?.add(record) || s.index,
    }));
  },
  // onUpdate handler - refetch from augmented view. A request's proposed
  // rates are hidden until job creation activates them, so their first
  // event is an update of a sheet the store does not have yet.
  async (item) => {
    const record = await pb.collection("rate_sheets_augmented").getOne<RateSheetResponse>(item.id);
    rateSheets.update((state) => {
      const known = state.items.some((i) => i.id === item.id);
      if (!known) {
        state.index?.add(record);
        return { ...state, items: [...state.items, record] };
      }
      return {
        ...state,
        items: state.items.map((i) => (i.id === item.id ? record : i)),
        index: state.index?.replace(record) || state.index,
      };
    });
  },
  "rate_sheets",
  fetchAllRateSheets,
  true, // refresh after a rate sheet absorb or its undo
);
