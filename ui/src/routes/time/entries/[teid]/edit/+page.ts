import type { TimeEntriesRecord } from "$lib/pocketbase-types";
import { pb } from "$lib/pocketbase";
import type { PageLoad } from "./$types";
import type { TimeEntriesPageData, WorkRecordDetailsResponse } from "$lib/svelte-types";
import { localIsoDate } from "$lib/localDate";

export const load: PageLoad<TimeEntriesPageData> = async ({ params }) => {
  const defaultItem = {
    uid: "",
    // date in YYYY-MM-DD format
    date: localIsoDate(),
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
    const item = await pb.collection("time_entries").getOne<TimeEntriesRecord>(params.teid);
    let linkedWorkRecord: TimeEntriesPageData["linked_work_record"] | undefined = undefined;
    if (item.work_record_id) {
      const workRecord = (await pb.send(`/api/work_records/${encodeURIComponent(item.work_record_id)}`, {
        method: "GET",
      })) as WorkRecordDetailsResponse;
      linkedWorkRecord = {
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
      };
    }
    return { item, editing: true, id: params.teid, linked_work_record: linkedWorkRecord };
  } catch (error) {
    console.error(`error loading data, returning default item: ${error}`);
    return { item: { ...defaultItem } as TimeEntriesRecord, editing: false, id: null };
  }
};
