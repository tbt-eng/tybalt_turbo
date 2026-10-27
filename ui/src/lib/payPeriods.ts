import { localIsoDate } from "$lib/localDate";

function parseIsoDate(dateString: string) {
  return new Date(`${dateString}T12:00:00`);
}

// Pay periods are two weeks long and end on a Saturday, so they start on the
// Sunday 13 days earlier.
export function payPeriodStart(payPeriodEnding: string) {
  const date = parseIsoDate(payPeriodEnding);
  date.setDate(date.getDate() - 13);
  return localIsoDate(date);
}

export function isUpcomingPayDate(payDate: string) {
  return payDate > localIsoDate();
}
