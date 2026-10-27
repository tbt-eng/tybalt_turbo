// Calendar dates (YYYY-MM-DD) are stored without a timezone. "Today" is the
// date on the user's own device, not UTC: toISOString() rolls over to the next
// day in the evening anywhere west of UTC.
export function localIsoDate(date = new Date()) {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  return `${year}-${month}-${day}`;
}
