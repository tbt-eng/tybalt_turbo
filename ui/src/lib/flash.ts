// A one-line message for the page navigated to, such as "Approved X" after an
// action moves the user on. Pass it as goto(url, { state: { flash } }); the
// layout's MessageBar shows it. It lives in the history entry, so Back and
// Forward keep it and a reload or fresh visit drops it.
export type Flash = {
  tone: "success" | "info";
  message: string;
  link?: { href: string; label: string };
};
