// The status of Help chat, for the App Config page. The types match
// app/helpchat/status.go. This file has no Vite or PocketBase imports, so a
// test can load it.

export interface HelpChatSettings {
  enabled: boolean;
  deploy_group: string;
  provider: string;
  model: string;
  effort: string;
  max_output_tokens: number;
  max_question_chars: number;
  questions_per_hour: number;
  max_tool_calls: number;
  history_turns: number;
  turn_timeout_seconds: number;
}

export interface HelpChatKey {
  provider: string;
  secret: string;
  set: boolean;
  in_use: boolean;
}

export interface DeployGroupStatus {
  id: string;
  name: string;
  found: boolean;
  members: number;
}

export interface HelpChatStatus {
  settings_error: string;
  has_row: boolean;
  settings: HelpChatSettings;
  keys: HelpChatKey[];
  deploy_group: DeployGroupStatus | null;
  deploy_groups: DeployGroupStatus[];
  corpus: {
    state: "not_loaded" | "loaded" | "failed";
    chunks: number;
    hash: string;
    error: string;
  };
  last_day: {
    questions: number;
    users: number;
    errors: number;
    tokens_in: number;
    tokens_cache_read: number;
    tokens_cache_write: number;
    tokens_out: number;
  };
  last_answered: string;
  last_error: { at: string; error: string } | null;
}

export interface HelpChatKeyTest {
  ok: boolean;
  provider: string;
  model: string;
  message: string;
}

export interface StatusNotice {
  tone: "info" | "success" | "warning" | "danger";
  title: string;
  text: string;
}

function people(count: number): string {
  return count === 1 ? "1 member" : `${count} members`;
}

// Says in plain words whether Help chat can run, and each reason why not. The
// result has one success notice, or one notice for each problem.
export function helpChatNotices(status: HelpChatStatus): StatusNotice[] {
  if (status.settings_error) {
    return [
      {
        tone: "danger",
        title: "The help_chat setting is not valid",
        text: `Help chat is off for everyone. ${status.settings_error}`,
      },
    ];
  }
  if (!status.has_row) {
    return [
      {
        tone: "info",
        title: "Help chat is off",
        text: "app_config has no help_chat row.",
      },
    ];
  }
  if (!status.settings.enabled) {
    return [
      {
        tone: "info",
        title: "Help chat is off",
        text: "The help_chat row has enabled set to false.",
      },
    ];
  }

  const problems: StatusNotice[] = [];
  const key = status.keys.find((entry) => entry.in_use);
  if (key && !key.set) {
    problems.push({
      tone: "danger",
      title: "The server has no key for the provider",
      // A secret set in the Fly web dashboard did not reach the app, though
      // Fly listed it as deployed. See descriptions/help_chat.md, step 2.
      text: `Set the secret ${key.secret} with the command "fly secrets set", which restarts the app. Do not use the Fly dashboard: a secret set there did not reach the app. Until then, a question gets "Help chat is not set up."`,
    });
  }
  const group = status.deploy_group;
  if (group && !group.found) {
    problems.push({
      tone: "warning",
      title: "The deploy group does not exist",
      text: `No deploy group has the id ${group.id}, so Help chat is off for everyone.`,
    });
  } else if (group && group.members === 0) {
    problems.push({
      tone: "warning",
      title: "The deploy group has no members",
      text: `Help chat is off for everyone until ${group.name} has a member.`,
    });
  }
  if (status.corpus.state === "failed") {
    problems.push({
      tone: "danger",
      title: "The search index did not load",
      text: status.corpus.error,
    });
  }
  if (problems.length > 0) return problems;

  const who = group ? `the ${people(group.members)} of ${group.name}` : "every user";
  return [
    {
      tone: "success",
      title: "Help chat is on",
      text: `It is on for ${who}. Use Test the key to check that the provider accepts the key.`,
    },
  ];
}

// The share of the input tokens that the model read from its cache, as a
// whole percent. Null when there were no input tokens.
export function cacheReadShare(day: HelpChatStatus["last_day"]): number | null {
  const total = day.tokens_in + day.tokens_cache_read + day.tokens_cache_write;
  if (total <= 0) return null;
  return Math.round((day.tokens_cache_read / total) * 100);
}

// Shows one property of an app_config row as text.
export function configValueText(value: unknown): string {
  if (typeof value === "string") return value === "" ? '""' : value;
  if (value === undefined) return "";
  return JSON.stringify(value);
}
