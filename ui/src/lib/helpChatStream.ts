// Reads the server-sent events of /api/help/ask. This file has no Vite or
// Svelte imports, so tests can run it directly.

export interface HelpEvent {
  event: string;
  data: unknown;
}

// Returns a function that takes the stream's text piece by piece and calls
// onEvent for each complete event. An event can arrive split across pieces.
export function createEventParser(onEvent: (event: HelpEvent) => void): (text: string) => void {
  let buffer = "";
  return (text) => {
    buffer += text.replaceAll("\r\n", "\n");
    let end = buffer.indexOf("\n\n");
    while (end >= 0) {
      const block = buffer.slice(0, end);
      buffer = buffer.slice(end + 2);
      end = buffer.indexOf("\n\n");
      let event = "message";
      const data: string[] = [];
      for (const line of block.split("\n")) {
        if (line.startsWith("event:")) event = line.slice("event:".length).trim();
        else if (line.startsWith("data:")) data.push(line.slice("data:".length).trimStart());
      }
      if (data.length === 0) continue;
      try {
        onEvent({ event, data: JSON.parse(data.join("\n")) });
      } catch {
        // A block that is not JSON is not one of the server's events.
      }
    }
  };
}

// HelpRequestError is a refusal from the server before the answer starts, such
// as "Help chat is off." The message is safe to show.
export class HelpRequestError extends Error {
  status: number;
  constructor(status: number, message: string) {
    super(message);
    this.status = status;
  }
}

export const helpUnavailable = "Help can't answer now. Try again in a minute.";

// Posts a question and calls onEvent for each event of the answer. It returns
// when the stream ends.
export async function streamHelp(
  request: typeof fetch,
  url: string,
  token: string,
  body: unknown,
  signal: AbortSignal,
  onEvent: (event: HelpEvent) => void,
): Promise<void> {
  const response = await request(url, {
    method: "POST",
    headers: { Authorization: token, "Content-Type": "application/json" },
    body: JSON.stringify(body),
    signal,
  });
  if (!response.ok || !response.body) {
    let message = helpUnavailable;
    try {
      const parsed = (await response.json()) as { message?: unknown };
      if (typeof parsed.message === "string" && parsed.message) message = parsed.message;
    } catch {
      // Keep the general message.
    }
    throw new HelpRequestError(response.status, message);
  }
  const parse = createEventParser(onEvent);
  const reader = response.body.getReader();
  const decoder = new TextDecoder();
  for (;;) {
    const { done, value } = await reader.read();
    if (done) break;
    parse(decoder.decode(value, { stream: true }));
  }
  parse(decoder.decode() + "\n\n");
}
