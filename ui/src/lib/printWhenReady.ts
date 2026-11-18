import { tick } from "svelte";

// Opens the browser print dialog once the page has rendered, its fonts have
// loaded and its images have decoded, so the dialog never shows a half-drawn
// page. Call it from onMount and return its result, which cancels the print
// when the page is left first.
export function printWhenReady(): () => void {
  let cancelled = false;
  (async () => {
    await tick();
    await document.fonts.ready;
    await Promise.all(Array.from(document.images).map((img) => img.decode().catch(() => {})));
    if (!cancelled) window.print();
  })();
  return () => {
    cancelled = true;
  };
}
