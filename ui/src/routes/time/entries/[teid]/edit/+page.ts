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
    work_record_subject_id: "",
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
    if (item.work_record_subject_id) {
      const subject = await pb
        .collection("work_records_subjects")
        .getOne<any>(item.work_record_subject_id);
      const workRecord = (await pb.send(
        `/api/work_records/${encodeURIComponent(subject.work_record)}`,
        {
          method: "GET",
        },
      )) as WorkRecordDetailsResponse;
      const detailSubject =
        workRecord.subjects?.find((row) => row.id === item.work_record_subject_id) ?? subject;
      linkedWorkRecord = {
        id: workRecord.id,
        subject_id: item.work_record_subject_id,
        work_record: workRecord.work_record,
        date: workRecord.date,
        job_id: workRecord.job_id,
        job_number: workRecord.job_number,
        type_id: workRecord.type_id,
        type_name: workRecord.type_name,
        location: workRecord.location,
        work_description: workRecord.work_description,
        hours_on_site: detailSubject.hours_on_site,
        hours_travel_time: detailSubject.hours_travel_time,
        approved: detailSubject.approved,
        linked_time_entry_id: detailSubject.linked_time_entry_id,
        linked_timesheet_id: detailSubject.linked_timesheet_id,
        linked_time_entry_date: detailSubject.linked_time_entry_date,
        linked_time_entry_description: detailSubject.linked_time_entry_description,
        linked_time_sheet_week_ending: detailSubject.linked_timesheet_week_ending,
        can_create_time_entry: detailSubject.can_create_time_entry,
      };
    }
    return { item, editing: true, id: params.teid, linked_work_record: linkedWorkRecord };
  } catch (error) {
    console.error(`error loading data, returning default item: ${error}`);
    return { item: { ...defaultItem } as TimeEntriesRecord, editing: false, id: null };
  }
};
