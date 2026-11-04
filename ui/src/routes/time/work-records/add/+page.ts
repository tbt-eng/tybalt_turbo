import { pb } from "$lib/pocketbase";
import type { PageLoad } from "./$types";
import { error } from "@sveltejs/kit";
import type { WorkRecordCopyDraftResponse } from "$lib/svelte-types";

function localToday() {
  const now = new Date();
  const year = now.getFullYear();
  const month = `${now.getMonth() + 1}`.padStart(2, "0");
  const day = `${now.getDate()}`.padStart(2, "0");
  return `${year}-${month}-${day}`;
}

export const load: PageLoad = async ({ url }) => {
  const currentUserId = pb.authStore.record?.id ?? "";
  const copyFromWorkRecordId = url.searchParams.get("copy_from") ?? "";

  if (copyFromWorkRecordId !== "") {
    try {
      const draft = (await pb.send(
        `/api/work_records/${encodeURIComponent(copyFromWorkRecordId)}/copy_to_tomorrow`,
        {
          method: "GET",
          requestKey: null,
        },
      )) as WorkRecordCopyDraftResponse;

      return {
        item: {
          creator: currentUserId,
          date: draft.date,
          job: draft.job,
          type: draft.type,
          location: draft.location,
          work_description: draft.work_description,
          report_to: draft.report_to,
          lane_kms: draft.lane_kms,
          field_book_number: draft.field_book_number,
          field_book_page_number:
            draft.field_book_page_number === 0 ? "" : draft.field_book_page_number,
          sub_contractor: draft.sub_contractor,
          equipment: draft.equipment,
          supplies: draft.supplies,
        },
        editing: false,
        id: null,
        detail: null,
        copyContext: draft,
        cancelAction: `/time/work-records/${draft.copy_from_work_record_id}`,
      };
    } catch (err: any) {
      console.error(`loading work record copy draft: ${err}`);
      throw error(err?.status === 403 ? 403 : 404, "Failed to load work record copy draft");
    }
  }

  let defaultType = "";

  try {
    const defaults = (await pb.send("/api/users/defaults", {
      method: "GET",
      requestKey: null,
    })) as { default_work_record_type?: string };
    defaultType = defaults.default_work_record_type ?? "";
  } catch (error) {
    console.error(`loading work record defaults: ${error}`);
  }

  return {
    item: {
      creator: currentUserId,
      date: localToday(),
      job: "",
      type: defaultType,
      location: "",
      work_description: "",
      report_to: "",
      lane_kms: 0,
      field_book_number: "",
      field_book_page_number: "",
      sub_contractor: "",
      equipment: "",
      supplies: "",
    },
    editing: false,
    id: null,
    detail: null,
    copyContext: null,
    cancelAction: "/time/work-records",
  };
};
