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
          uid: draft.uid,
          creator: currentUserId,
          date: draft.date,
          job: draft.job,
          type: draft.type,
          location: draft.location,
          work_description: draft.work_description,
          report_to: draft.report_to,
          hours_on_site: draft.hours_on_site,
          hours_travel_time: draft.hours_travel_time,
          distance_travelled_km: draft.distance_travelled_km,
          field_book_number: draft.field_book_number,
          field_book_page_number: draft.field_book_page_number === 0 ? "" : draft.field_book_page_number,
          sub_contractor: draft.sub_contractor,
          equipment: draft.equipment,
          supplies: draft.supplies,
          vehicle_type: draft.vehicle_type,
          is_passenger: draft.is_passenger,
          company_vehicle_unit_number: draft.company_vehicle_unit_number,
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
      uid: currentUserId,
      creator: currentUserId,
      date: localToday(),
      job: "",
      type: defaultType,
      location: "",
      work_description: "",
      report_to: "",
      hours_on_site: 0,
      hours_travel_time: 0,
      distance_travelled_km: 0,
      field_book_number: "",
      field_book_page_number: "",
      sub_contractor: "",
      equipment: "",
      supplies: "",
      vehicle_type: "",
      is_passenger: false,
      company_vehicle_unit_number: "",
    },
    editing: false,
    id: null,
    detail: null,
    copyContext: null,
    cancelAction: "/time/work-records",
  };
};
