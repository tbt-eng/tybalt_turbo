import { pb } from "$lib/pocketbase";
import type { PageLoad } from "./$types";
import type { WorkRecordDetailsResponse } from "$lib/svelte-types";
import { error } from "@sveltejs/kit";

export const load: PageLoad = async ({ params }) => {
  try {
    const detail = (await pb.send(`/api/work_records/${encodeURIComponent(params.workRecord)}`, {
      method: "GET",
    })) as WorkRecordDetailsResponse;

    return {
      item: {
        id: detail.id,
        uid: detail.uid,
        creator: detail.creator,
        date: detail.date,
        job: detail.job_id,
        type: detail.type_id,
        location: detail.location,
        work_description: detail.work_description,
        report_to: detail.report_to,
        hours_on_site: detail.hours_on_site,
        hours_travel_time: detail.hours_travel_time,
        distance_travelled_km: detail.distance_travelled_km,
        field_book_number: detail.field_book_number,
        field_book_page_number: detail.field_book_page_number === 0 ? "" : detail.field_book_page_number,
        sub_contractor: detail.sub_contractor,
        equipment: detail.equipment,
        supplies: detail.supplies,
        vehicle_type: detail.vehicle_type,
        is_passenger: detail.is_passenger,
        company_vehicle_unit_number: detail.company_vehicle_unit_number,
      },
      editing: true,
      id: detail.id,
      detail,
    };
  } catch (err) {
    console.error(`loading work record for edit: ${err}`);
    throw error(404, "Failed to load work record");
  }
};
