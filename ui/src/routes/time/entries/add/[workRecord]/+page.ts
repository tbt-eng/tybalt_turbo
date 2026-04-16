import type { TimeEntriesRecord } from "$lib/pocketbase-types";
import { pb } from "$lib/pocketbase";
import type { PageLoad } from "./$types";
import type { TimeEntriesPageData, WorkRecordDetailsResponse } from "$lib/svelte-types";
import { error, redirect } from "@sveltejs/kit";

export const load: PageLoad<TimeEntriesPageData> = async ({ params }) => {
  const defaultItem = {
    uid: "",
    date: new Date().toISOString().split("T")[0],
    time_type: "sdyfl3q7j7ap849",
    division: "",
    description: "",
    job: "",
    work_record_id: "",
    work_record: "",
    hours: 0,
    meals_hours: 0,
    payout_request_amount: 0,
    category: "",
    week_ending: "2006-01-02",
  };

  try {
    const workRecord = (await pb.send(`/api/work_records/${encodeURIComponent(params.workRecord)}`, {
      method: "GET",
    })) as WorkRecordDetailsResponse;

    if (!workRecord.can_create_time_entry && workRecord.linked_time_entry_id) {
      if (workRecord.linked_timesheet_id) {
        throw redirect(
          303,
          `/time/sheets/${workRecord.linked_timesheet_id}/details?entry=${workRecord.linked_time_entry_id}`,
        );
      }
      throw redirect(303, `/time/entries/${workRecord.linked_time_entry_id}/details`);
    }

    return {
      item: {
        ...defaultItem,
        uid: workRecord.uid,
        date: workRecord.date,
        description: workRecord.work_description,
        job: workRecord.job_id,
        work_record_id: workRecord.id,
        hours: workRecord.hours_on_site + workRecord.hours_travel_time,
      } as TimeEntriesRecord,
      editing: false,
      id: null,
      linked_work_record: {
        id: workRecord.id,
        work_record: workRecord.work_record,
        date: workRecord.date,
        job_id: workRecord.job_id,
        job_number: workRecord.job_number,
        type_id: workRecord.type_id,
        type_name: workRecord.type_name,
        location: workRecord.location,
        work_description: workRecord.work_description,
        hours_on_site: workRecord.hours_on_site,
        hours_travel_time: workRecord.hours_travel_time,
        approved: workRecord.approved,
        linked_time_entry_id: workRecord.linked_time_entry_id,
        linked_timesheet_id: workRecord.linked_timesheet_id,
        linked_time_entry_date: workRecord.linked_time_entry_date,
        linked_time_entry_description: workRecord.linked_time_entry_description,
        linked_time_sheet_week_ending: workRecord.linked_time_sheet_week_ending,
        can_create_time_entry: workRecord.can_create_time_entry,
      },
    };
  } catch (err) {
    if (err && typeof err === "object" && "status" in err) {
      throw err;
    }
    console.error(`loading linked work record: ${err}`);
    throw error(404, "Failed to load work record");
  }
};
