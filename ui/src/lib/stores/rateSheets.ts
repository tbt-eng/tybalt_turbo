import { createCollectionStore } from "./collectionStore";
import { pb } from "$lib/pocketbase";
import type { BaseSystemFields } from "$lib/pocketbase-types";

export interface RateSheetResponse extends BaseSystemFields {
  name: string;
  effective_date: string;
  revision: number;
  active: boolean;
  job_count: number;
}

// Fetch from the augmented view which includes job_count
const fetchAllRateSheets = async (): Promise<RateSheetResponse[]> =>
  pb.collection("rate_sheets_augmented").getFullList({ sort: "-effective_date" });

export const rateSheets = createCollectionStore<RateSheetResponse>(
  "rate_sheets", // Subscribe to the base collection for real-time updates
  { requestKey: "rate_sheets" },
  {
    fields: ["name", "effective_date"],
    storeFields: ["id", "name", "effective_date", "revision", "active", "job_count"],
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
