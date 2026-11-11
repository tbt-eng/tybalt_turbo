import test from "node:test";
import assert from "node:assert/strict";
import {
  cacheReadShare,
  configValueText,
  helpChatNotices,
  type HelpChatStatus,
} from "../src/lib/helpChatStatus.ts";

function status(change: Partial<HelpChatStatus> = {}): HelpChatStatus {
  return {
    settings_error: "",
    has_row: true,
    settings: {
      enabled: true,
      deploy_group: "",
      provider: "anthropic",
      model: "claude-haiku-5-5",
      effort: "low",
      max_output_tokens: 8000,
      max_question_chars: 2000,
      questions_per_hour: 30,
      max_tool_calls: 6,
      history_turns: 10,
      turn_timeout_seconds: 60,
    },
    keys: [
      { provider: "anthropic", secret: "ANTHROPIC_API_KEY", set: true, in_use: true },
      { provider: "openai", secret: "OPENAI_API_KEY", set: false, in_use: false },
    ],
    deploy_group: null,
    deploy_groups: [],
    corpus: { state: "not_loaded", chunks: 0, hash: "", error: "" },
    last_day: {
      questions: 0,
      users: 0,
      errors: 0,
      tokens_in: 0,
      tokens_cache_read: 0,
      tokens_cache_write: 0,
      tokens_out: 0,
    },
    last_answered: "",
    last_error: null,
    ...change,
  };
}

test("a chat that can run gets one success notice that says who has it", () => {
  let notices = helpChatNotices(status());
  assert.equal(notices.length, 1);
  assert.equal(notices[0].tone, "success");
  assert.match(notices[0].text, /every user/);

  notices = helpChatNotices(
    status({ deploy_group: { id: "g1", name: "testers", found: true, members: 3 } }),
  );
  assert.equal(notices[0].tone, "success");
  assert.match(notices[0].text, /the 3 members of testers/);
});

test("a chat that is off says why", () => {
  assert.match(helpChatNotices(status({ has_row: false }))[0].text, /no help_chat row/);
  const off = status();
  off.settings.enabled = false;
  assert.equal(helpChatNotices(off)[0].tone, "info");
  const bad = helpChatNotices(status({ settings_error: 'unknown field "deploy_grp"' }));
  assert.equal(bad[0].tone, "danger");
  assert.match(bad[0].text, /deploy_grp/);
});

test("each problem of a chat that is on gets its own notice", () => {
  const broken = status({
    keys: [{ provider: "anthropic", secret: "ANTHROPIC_API_KEY", set: false, in_use: true }],
    deploy_group: { id: "nosuchgroup", name: "", found: false, members: 0 },
    corpus: { state: "failed", chunks: 0, hash: "", error: "no corpus" },
  });
  const notices = helpChatNotices(broken);
  assert.equal(notices.length, 3);
  assert.match(notices[0].text, /ANTHROPIC_API_KEY/);
  assert.match(notices[0].text, /fly secrets set/);
  assert.match(notices[0].text, /Do not use the Fly dashboard/);
  assert.match(notices[1].text, /nosuchgroup/);
  assert.equal(notices[2].text, "no corpus");
  assert.ok(notices.every((notice) => notice.tone !== "success"));

  const empty = helpChatNotices(
    status({ deploy_group: { id: "g1", name: "testers", found: true, members: 0 } }),
  );
  assert.equal(empty[0].tone, "warning");
  assert.match(empty[0].text, /testers/);
});

test("the cache share counts read, written and plain input tokens", () => {
  const day = status().last_day;
  assert.equal(cacheReadShare(day), null);
  assert.equal(
    cacheReadShare({ ...day, tokens_in: 100, tokens_cache_read: 800, tokens_cache_write: 100 }),
    80,
  );
});

test("a config value shows as text", () => {
  assert.equal(configValueText("anthropic"), "anthropic");
  assert.equal(configValueText(""), '""');
  assert.equal(configValueText(false), "false");
  assert.equal(configValueText(0), "0");
  assert.equal(configValueText({ mode: "percent", value: 5 }), '{"mode":"percent","value":5}');
});
