import type { WorkRecordsTypesResponse } from "$lib/pocketbase-types";
import { createCollectionStore } from "./collectionStore";

export const workRecordTypes = createCollectionStore<WorkRecordsTypesResponse>(
  "work_records_types",
  {
    sort: "sort_order,name",
    filter: "active = true",
    requestKey: "work_record_types",
  },
  {
    idField: "id",
    fields: ["name"],
    storeFields: ["id", "name", "sort_order", "active"],
    extractField: (doc, field) => (doc as WorkRecordsTypesResponse)[field as keyof WorkRecordsTypesResponse],
  },
);
