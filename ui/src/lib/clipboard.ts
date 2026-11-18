// Copy text to the clipboard. It resolves true when the browser took the
// text, and false when there is no clipboard (an insecure page, an old
// browser) or the browser refused. It never throws, so a caller only shows
// the result. DSCopyButton uses it.
export async function copyText(
  text: string,
  clipboard: Pick<Clipboard, "writeText"> | undefined = globalThis.navigator?.clipboard,
): Promise<boolean> {
  if (typeof clipboard?.writeText !== "function") return false;
  try {
    await clipboard.writeText(text);
    return true;
  } catch {
    return false;
  }
}
