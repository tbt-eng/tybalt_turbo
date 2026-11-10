// The state of the Help chat: the turns of the conversation that is open, and
// the calls to /api/help. The server keeps the conversation. This store holds
// only what the screen shows.
import { derived, get, writable } from "svelte/store";
import { pb } from "$lib/pocketbase";
import { renderAnswer, type ChatSource } from "$lib/helpChatFormat";
import { HelpRequestError, helpUnavailable, streamHelp } from "$lib/helpChatStream";

export type TurnStatus = "running" | "answered" | "gap" | "refused" | "cancelled" | "error";

export interface ChatTurn {
  // A number for the screen. The server's id arrives with the first event.
  key: number;
  id: string;
  question: string;
  text: string;
  // The answer as safe HTML. See renderAnswer.
  html: string;
  status: TurnStatus;
  // What the server is doing, while the turn runs and has no text yet.
  activity: string;
  // What the user reads when there is no answer.
  message: string;
  sources: ChatSource[];
  rating: "" | "up" | "down";
  noteSent: boolean;
}

interface HelpChatState {
  // The conversation. It stays while the user reads the docs.
  turns: ChatTurn[];
  // True while the conversation shows, and false while the docs show.
  open: boolean;
}

const finalStatuses: TurnStatus[] = ["answered", "gap", "refused", "cancelled", "error"];

const store = writable<HelpChatState>({ turns: [], open: false });

let chatId = "";
let nextKey = 1;
let abort: AbortController | null = null;

function turnByKey(key: number): ChatTurn | undefined {
  return get(store).turns.find((turn) => turn.key === key);
}

function change(key: number, values: Partial<ChatTurn>) {
  store.update((state) => ({
    ...state,
    turns: state.turns.map((turn) => (turn.key === key ? { ...turn, ...values } : turn)),
  }));
}

function isBusy(): boolean {
  return get(store).turns.at(-1)?.status === "running";
}

// Asks a question. A question in an open conversation continues it.
async function ask(question: string, route: string): Promise<void> {
  if (isBusy()) return;
  const key = nextKey++;
  store.update((state) => ({
    open: true,
    turns: [
      ...state.turns,
      {
        key,
        id: "",
        question,
        text: "",
        html: "",
        status: "running",
        activity: "",
        message: "",
        sources: [],
        rating: "",
        noteSent: false,
      },
    ],
  }));
  const controller = new AbortController();
  abort = controller;
  let text = "";
  let finished = false;

  try {
    await streamHelp(
      fetch,
      `${pb.baseURL}/api/help/ask`,
      pb.authStore.token,
      { question, chat: chatId, route, surface: "help" },
      controller.signal,
      ({ event, data }) => {
        const value = (data ?? {}) as Record<string, unknown>;
        if (event === "meta") {
          chatId = String(value.chat ?? "");
          change(key, { id: String(value.turn ?? "") });
        } else if (event === "status") {
          change(key, { activity: value.state === "searching" ? "Searching Turbo…" : "Reading…" });
        } else if (event === "text") {
          text += String(value.delta ?? "");
          change(key, { text, html: renderAnswer(text, null) });
        } else if (event === "sources" && Array.isArray(data)) {
          const sources = data as ChatSource[];
          change(key, { sources, html: renderAnswer(text, sources) });
        } else if (event === "done") {
          const status = finalStatuses.find((candidate) => candidate === value.status);
          finished = true;
          change(key, {
            status: status ?? "error",
            message: String(value.message ?? "") || (status ? "" : helpUnavailable),
          });
        }
      },
    );
    // A stream that ends with no result was cut off.
    if (!finished) change(key, { status: "error", message: helpUnavailable });
  } catch (error) {
    if (controller.signal.aborted) {
      change(key, { status: "cancelled", message: "You stopped this answer." });
    } else {
      change(key, {
        status: "error",
        message: error instanceof HelpRequestError ? error.message : helpUnavailable,
      });
    }
  } finally {
    if (abort === controller) abort = null;
  }
}

// Stops the answer that is running. The server stops the model too.
function stop() {
  abort?.abort();
}

// Shows the conversation or the docs. The conversation shows only if there is
// one.
function show(conversation: boolean) {
  store.update((state) => ({ ...state, open: conversation && state.turns.length > 0 }));
}

// Ends the conversation and shows the docs page again.
function close() {
  stop();
  chatId = "";
  store.set({ turns: [], open: false });
}

// Feedback is a signal for the people who improve Help. A failure to save it
// is not the user's problem, so it shows nothing.
function feedback(key: number, body: Record<string, unknown>) {
  const id = turnByKey(key)?.id;
  if (!id) return;
  pb.send(`/api/help/turns/${id}/feedback`, { method: "POST", body, requestKey: null }).catch(
    () => {},
  );
}

function rate(key: number, rating: "up" | "down") {
  // A second press on the same thumb takes the rating back.
  const next = turnByKey(key)?.rating === rating ? "" : rating;
  change(key, { rating: next });
  feedback(key, { rating: next });
}

function sendNote(key: number, note: string) {
  if (!note.trim()) return;
  change(key, { noteSent: true });
  feedback(key, { note });
}

// Asks the server to get the model ready, so the first answer starts faster.
function warm() {
  pb.send("/api/help/warm", { method: "POST", requestKey: null }).catch(() => {});
}

export const helpChat = {
  subscribe: store.subscribe,
  ask,
  stop,
  show,
  close,
  rate,
  sendNote,
  warm,
  // The user opened a source: a chunk's ref or a docs page's url.
  opened: (key: number, source: string) => feedback(key, { opened: source }),
  // The user copied the answer. Returns its text.
  copied: (key: number): string => {
    feedback(key, { copied: true });
    return turnByKey(key)?.text ?? "";
  },
};

export const helpChatBusy = derived(store, ($store) => $store.turns.at(-1)?.status === "running");
