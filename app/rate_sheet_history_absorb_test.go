package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"tybalt/absorb"
	"tybalt/routes"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

// Rate sheet absorb and the project rate sheet history (job_rate_sheets).
// Time is priced from the history, so a period is a use of its sheet, moves
// with an absorb and is restored by undo. The fixtures are rows of the seed
// CSV files (app/test_seed_data/data). Rates H are the 2025 Standard Rates
// sheet (c41ofep525bcacj) with 7000 added to each rate, so that they match no
// other seeded sheet.
//
// Sheets with rates H:
//   - Absorb History Target (rsabhisttgt0001): rev 0, 2024-01-01, active
//   - Absorb History Source (rsabhistsrc0001): rev 0, 2024-03-01, active;
//     only the 2024-04-01 period of jobabhist000001 uses it
//   - Absorb History Current (rsabhistcur0001): rev 0, 2024-03-01, active;
//     the current sheet of jobabhist000002
//   - Absorb History Revised rev 0 (rsabhistrev0r01): 2024-02-01, inactive;
//     only the opening period of jobabhist000003 uses it
//   - Absorb History Revised rev 1 (rsabhistrev1r01): 2024-05-01, active;
//     nothing uses it
//
// Projects:
//   - jobabhist000001, on the standard sheet: opening period on the standard
//     sheet, from 2024-04-01 on Absorb History Source, from 2025-02-01 on the
//     standard sheet. One time entry in each period, role 0l1llckilwz1mqz
//     (160 on the standard sheet, 7160 on rates H): 2 h, 3 h and 1 h.
//   - jobabhist000002, on Absorb History Current: opening period on the
//     target, from 2024-04-01 on Absorb History Current.
//   - jobabhist000003, on the standard sheet: opening period on Absorb
//     History Revised rev 0, from 2025-02-01 on the standard sheet.
//
// history.revise@example.com is active and holds job, rate_sheet_revise and
// absorb. One test changes fixture data other than through the API, and
// says so where it does.

const (
	rsHistTarget  = "rsabhisttgt0001"
	rsHistSource  = "rsabhistsrc0001"
	rsHistCurrent = "rsabhistcur0001"
	rsHistRev0    = "rsabhistrev0r01"
	rsHistRev1    = "rsabhistrev1r01"
	histJobOld    = "jobabhist000001"
	histJobNow    = "jobabhist000002"
	histJobRev    = "jobabhist000003"
	histPeriodSrc = "jrsabhist01src0"
	histPeriodCur = "jrsabhist02cur0"
	histRevise    = "history.revise@example.com"
)

// periodSnapshot returns every field of each period of a project, in date
// order, so that undo can be checked to restore the history exactly.
func periodSnapshot(t testing.TB, app core.App, job string) []string {
	t.Helper()
	rows, err := app.FindRecordsByFilter("job_rate_sheets", "job = {:job}", "effective_date", 0, 0, dbx.Params{"job": job})
	if err != nil {
		t.Fatal(err)
	}
	result := []string{}
	for _, row := range rows {
		raw, _ := json.Marshal(recordAllFields(row))
		result = append(result, string(raw))
	}
	return result
}

func recordAllFields(record *core.Record) map[string]any {
	fields := map[string]any{"id": record.Id}
	for _, field := range record.Collection().Fields {
		fields[field.GetName()] = record.Get(field.GetName())
	}
	return fields
}

// jobWritebackState is what the jobs writeback reads to decide whether to
// send a job.
func jobWritebackState(t testing.TB, app core.App, job string) string {
	t.Helper()
	record, err := app.FindRecordById("jobs", job)
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := json.Marshal([]any{record.GetString("rate_sheet"), record.GetBool("_imported"), record.GetString("updated")})
	return string(raw)
}

func histWIP(t *testing.T, app *tests.TestApp, job string) routes.JobWIP {
	t.Helper()
	status, body := absorbRequest(t, app, http.MethodGet, "/api/jobs/"+job+"/wip", "valuation1@example.com", nil)
	if status != http.StatusOK {
		t.Fatalf("WIP status %d: %s", status, body)
	}
	var result routes.JobWIP
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		t.Fatal(err)
	}
	return result
}

func commitRateSheetAbsorb(t *testing.T, app *tests.TestApp) {
	t.Helper()
	path := "/api/collections/absorb_actions/records/" + rateSheetAction(t, app).Id
	if status, body := absorbRequest(t, app, http.MethodDelete, path, absorbAdminEmail, nil); status != http.StatusNoContent {
		t.Fatalf("commit status %d: %s", status, body)
	}
}

// A single revision sheet that only an older period uses is absorbed. The
// period moves to the target and the time keeps its value. The job itself
// does not change, so it is not sent to the writeback. Undo restores the
// history exactly, and commit keeps the period on the target.
func TestRateSheetAbsorbMovesHistoryPeriod(t *testing.T) {
	app := newRateSheetAbsorbApp(t)
	defer app.Cleanup()

	before := histWIP(t, app, histJobOld)
	// 2 h and 1 h at 160 on the standard sheet, 3 h at 7160 on the source.
	if before.Hours != 6 || before.TimeValue != 21960 {
		t.Fatalf("WIP before = %+v; want 6 h worth 21960", before)
	}
	periods := periodSnapshot(t, app, histJobOld)
	job := jobWritebackState(t, app, histJobOld)

	if status, body := absorbRequest(t, app, http.MethodPost, "/api/rate_sheets/"+rsHistTarget+"/absorb", absorbAdminEmail, absorbBody(rsHistSource)); status != http.StatusOK {
		t.Fatalf("absorb status %d: %s", status, body)
	}
	if rateSheetExists(t, app, rsHistSource) {
		t.Fatal("the source sheet was not absorbed")
	}
	if got := fieldOf(t, app, "job_rate_sheets", histPeriodSrc, "rate_sheet"); got != rsHistTarget {
		t.Fatalf("the period is on %s; want the target", got)
	}
	if got := jobWritebackState(t, app, histJobOld); got != job {
		t.Fatalf("the job changed from %s to %s; only its period moved", job, got)
	}
	var updates map[string]map[string]map[string]string
	if err := rateSheetAction(t, app).UnmarshalJSONField("updated_references", &updates); err != nil {
		t.Fatal(err)
	}
	if updates["job_rate_sheets"]["rate_sheet"][histPeriodSrc] != rsHistSource {
		t.Fatalf("the undo data does not hold the period: %v", updates)
	}
	if _, tracked := updates["jobs"]["rate_sheet"][histJobOld]; tracked {
		t.Fatal("the undo data holds the job, which did not move")
	}
	if after := histWIP(t, app, histJobOld); after != before {
		t.Fatalf("WIP after absorb = %+v; want %+v", after, before)
	}

	if status, body := absorbRequest(t, app, http.MethodPost, "/api/rate_sheets/undo_absorb", absorbAdminEmail, nil); status != http.StatusOK {
		t.Fatalf("undo status %d: %s", status, body)
	}
	if got := periodSnapshot(t, app, histJobOld); strings.Join(got, "\n") != strings.Join(periods, "\n") {
		t.Fatalf("undo restored the history\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(periods, "\n"))
	}
	if got := jobWritebackState(t, app, histJobOld); got != job {
		t.Fatalf("undo changed the job from %s to %s", job, got)
	}
	if undone := histWIP(t, app, histJobOld); undone != before {
		t.Fatalf("WIP after undo = %+v; want %+v", undone, before)
	}

	if err := routes.AbsorbRecords(app, "rate_sheets", rsHistTarget, []string{rsHistSource}); err != nil {
		t.Fatal(err)
	}
	commitRateSheetAbsorb(t, app)
	if got := fieldOf(t, app, "job_rate_sheets", histPeriodSrc, "rate_sheet"); got != rsHistTarget {
		t.Fatalf("after commit the period is on %s", got)
	}
	if rateSheetExists(t, app, rsHistSource) {
		t.Fatal("commit restored the source sheet")
	}
	if committed := histWIP(t, app, histJobOld); committed != before {
		t.Fatalf("WIP after commit = %+v; want %+v", committed, before)
	}
}

// When the job's current sheet is absorbed, the job and its latest period move
// together, so jobs.rate_sheet stays the sheet of the latest period. The
// period before it is already on the target: the two neighbours then share a
// sheet. While the absorb is pending, the project's history can't change.
// After commit, that pair does not block other changes, and a change can't
// make a new pair.
func TestRateSheetAbsorbCurrentSheetWithHistory(t *testing.T) {
	app := newRateSheetAbsorbApp(t)
	defer app.Cleanup()

	periods := periodSnapshot(t, app, histJobNow)
	if status, body := absorbRequest(t, app, http.MethodPost, "/api/rate_sheets/"+rsHistTarget+"/absorb", absorbAdminEmail, absorbBody(rsHistCurrent)); status != http.StatusOK {
		t.Fatalf("absorb status %d: %s", status, body)
	}
	if got := fieldOf(t, app, "jobs", histJobNow, "rate_sheet"); got != rsHistTarget {
		t.Fatalf("the job is on %s", got)
	}
	if got := fieldOf(t, app, "job_rate_sheets", histPeriodCur, "rate_sheet"); got != rsHistTarget {
		t.Fatalf("the latest period is on %s", got)
	}
	job, err := app.FindRecordById("jobs", histJobNow)
	if err != nil {
		t.Fatal(err)
	}
	if job.GetBool("_imported") {
		t.Fatal("the moved job was not marked for writeback")
	}

	// The project's history is frozen until the absorb is committed or undone.
	newPeriod := map[string]any{"rate_sheet": rsAbStandardRate, "effective_date": "2025-06-02", "reason": "New contract"}
	status, body := absorbRequest(t, app, http.MethodPost, "/api/jobs/"+histJobNow+"/rate_sheets", histRevise, newPeriod)
	if status != http.StatusConflict || !strings.Contains(body, "waiting to be committed or undone") {
		t.Fatalf("history change while the absorb is pending: status %d: %s", status, body)
	}
	status, body = absorbRequest(t, app, http.MethodPatch, "/api/jobs/"+histJobNow+"/rate_sheets/"+histPeriodCur, histRevise, map[string]any{"effective_date": "2024-04-02", "reason": "Wrong date"})
	if status != http.StatusConflict {
		t.Fatalf("period edit while the absorb is pending: status %d: %s", status, body)
	}
	// A project that the absorb did not move can still change.
	status, body = absorbRequest(t, app, http.MethodPatch, "/api/jobs/"+histJobRev+"/rate_sheets/jrsabhist03std0", histRevise, map[string]any{"effective_date": "2025-02-03", "reason": "Wrong date"})
	if status != http.StatusOK {
		t.Fatalf("edit of another project while the absorb is pending: status %d: %s", status, body)
	}

	// Test-time change to fixture data: the moved job's updated is set to a
	// past time, as if the absorb was done earlier and the writeback has sent
	// the job with the target's id since. Undo must mark the job to be sent
	// again, with the restored sheet.
	const sentAt = "2026-01-01 00:00:00.000Z"
	rateSheetAbsorbSQL(t, app, "UPDATE jobs SET updated = '"+sentAt+"' WHERE id = {:job}", dbx.Params{"job": histJobNow})

	if status, body := absorbRequest(t, app, http.MethodPost, "/api/rate_sheets/undo_absorb", absorbAdminEmail, nil); status != http.StatusOK {
		t.Fatalf("undo status %d: %s", status, body)
	}
	if job, err = app.FindRecordById("jobs", histJobNow); err != nil {
		t.Fatal(err)
	}
	if job.GetBool("_imported") || job.GetString("updated") <= sentAt {
		t.Fatalf("undo did not mark the job for writeback: _imported %v, updated %s", job.GetBool("_imported"), job.GetString("updated"))
	}
	if got := periodSnapshot(t, app, histJobNow); strings.Join(got, "\n") != strings.Join(periods, "\n") {
		t.Fatalf("undo restored the history\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(periods, "\n"))
	}
	if got := fieldOf(t, app, "jobs", histJobNow, "rate_sheet"); got != rsHistCurrent {
		t.Fatalf("undo left the job on %s", got)
	}

	if err := routes.AbsorbRecords(app, "rate_sheets", rsHistTarget, []string{rsHistCurrent}); err != nil {
		t.Fatal(err)
	}
	commitRateSheetAbsorb(t, app)

	// Both periods are on the target now. A new period is still allowed.
	status, body = absorbRequest(t, app, http.MethodPost, "/api/jobs/"+histJobNow+"/rate_sheets", histRevise, newPeriod)
	if status != http.StatusOK {
		t.Fatalf("new period after the commit: status %d: %s", status, body)
	}
	if got := fieldOf(t, app, "jobs", histJobNow, "rate_sheet"); got != rsAbStandardRate {
		t.Fatalf("the job is on %s; want the new period's sheet", got)
	}
	added, err := app.FindFirstRecordByFilter("job_rate_sheets", "job = {:job} && effective_date = '2025-06-02'", dbx.Params{"job": histJobNow})
	if err != nil {
		t.Fatal(err)
	}
	// A change that makes a new pair of neighbours on one sheet is refused.
	status, body = absorbRequest(t, app, http.MethodPatch, "/api/jobs/"+histJobNow+"/rate_sheets/"+added.Id, histRevise, map[string]any{"rate_sheet": rsHistTarget, "reason": "Back to the target"})
	if status != http.StatusBadRequest || !strings.Contains(body, "same_as_previous") {
		t.Fatalf("a new pair on one sheet: status %d: %s", status, body)
	}
	// The later period of the pair can be deleted.
	status, body = absorbRequest(t, app, http.MethodDelete, "/api/jobs/"+histJobNow+"/rate_sheets/"+histPeriodCur, histRevise, nil)
	if status != http.StatusOK {
		t.Fatalf("delete of the later period of the pair: status %d: %s", status, body)
	}
	if n := len(periodSnapshot(t, app, histJobNow)); n != 2 {
		t.Fatalf("the history has %d periods; want 2", n)
	}
}

// A name with revisions can't be absorbed when only a period uses one of its
// revisions. The Duplicates view and the candidates say why, and count the
// project as a job that uses the sheet.
func TestRateSheetAbsorbNameUsedByHistory(t *testing.T) {
	app := newRateSheetAbsorbApp(t)
	defer app.Cleanup()

	const reason = "a job or project request uses one of them"
	for _, id := range []string{rsHistRev0, rsHistRev1} {
		err := routes.AbsorbRecords(app, "rate_sheets", rsHistTarget, []string{id})
		var refused *absorb.RefusedError
		if !errors.As(err, &refused) || !strings.Contains(refused.Message, reason) {
			t.Fatalf("absorb of %s: %v; want a refusal that says %q", id, err, reason)
		}
		if !rateSheetExists(t, app, id) {
			t.Fatalf("a refused absorb deleted %s", id)
		}
	}

	status, body := absorbRequest(t, app, http.MethodGet, "/api/rate_sheets/duplicates", reviseOnlyEmail, nil)
	if status != http.StatusOK {
		t.Fatalf("duplicates status %d: %s", status, body)
	}
	var duplicates struct {
		Groups []struct {
			Sheets []struct {
				ID       string `json:"id"`
				Target   bool   `json:"target"`
				Reason   string `json:"reason"`
				Absorbs  int    `json:"absorbs"`
				JobCount int    `json:"job_count"`
			} `json:"sheets"`
		} `json:"groups"`
	}
	if err := json.Unmarshal([]byte(body), &duplicates); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, group := range duplicates.Groups {
		if group.Sheets[0].ID != rsHistTarget {
			continue
		}
		found = true
		if !group.Sheets[0].Target {
			t.Fatal("the earliest sheet is not marked as the target")
		}
		want := map[string]struct {
			reason string
			jobs   int
		}{
			rsHistSource:  {"", 1},
			rsHistCurrent: {"", 1},
			rsHistRev0:    {reason, 1},
			rsHistRev1:    {reason, 0},
		}
		for _, sheet := range group.Sheets[1:] {
			w, ok := want[sheet.ID]
			if !ok {
				t.Fatalf("unexpected sheet %s in the group", sheet.ID)
			}
			if (w.reason == "") != (sheet.Reason == "") || !strings.Contains(sheet.Reason, w.reason) {
				t.Fatalf("%s reason = %q; want %q", sheet.ID, sheet.Reason, w.reason)
			}
			if sheet.JobCount != w.jobs {
				t.Fatalf("%s job count = %d; want %d", sheet.ID, sheet.JobCount, w.jobs)
			}
			delete(want, sheet.ID)
		}
		if len(want) != 0 {
			t.Fatalf("missing sheets %v", want)
		}
		// The target's only use is the opening period of jobabhist000002.
		if group.Sheets[0].JobCount != 1 {
			t.Fatalf("target job count = %d; want 1", group.Sheets[0].JobCount)
		}
	}
	if !found {
		t.Fatal("the rates H group is missing")
	}

	status, body = absorbRequest(t, app, http.MethodGet, "/api/rate_sheets/"+rsHistTarget+"/absorb_candidates", reviseOnlyEmail, nil)
	if status != http.StatusOK {
		t.Fatalf("candidates status %d: %s", status, body)
	}
	var candidates struct {
		Candidates []struct {
			ID       string `json:"id"`
			Reason   string `json:"reason"`
			Absorbs  int    `json:"absorbs"`
			JobCount int    `json:"job_count"`
		} `json:"candidates"`
	}
	if err := json.Unmarshal([]byte(body), &candidates); err != nil {
		t.Fatal(err)
	}
	got := map[string]string{}
	for _, c := range candidates.Candidates {
		got[c.ID] = c.Reason
		if c.ID == rsHistSource && (c.Absorbs != 1 || c.JobCount != 1) {
			t.Fatalf("source candidate absorbs %d, job count %d; want 1 and 1", c.Absorbs, c.JobCount)
		}
	}
	if r, ok := got[rsHistSource]; !ok || r != "" {
		t.Fatalf("source candidate reason = %q (listed %v); want none", r, ok)
	}
	for _, id := range []string{rsHistRev0, rsHistRev1} {
		if !strings.Contains(got[id], reason) {
			t.Fatalf("%s candidate reason = %q; want %q", id, got[id], reason)
		}
	}
}

// The Rate Sheets list (rate_sheets_augmented) counts a job that uses a sheet
// now or in any period of its history, once, as the Duplicates view does.
func TestRateSheetListCountsHistoryJobs(t *testing.T) {
	app := newRateSheetAbsorbApp(t)
	defer app.Cleanup()

	for id, want := range map[string]int{
		rsHistSource:  1, // only the 2024-04-01 period of jobabhist000001
		rsHistCurrent: 1, // jobabhist000002: its current sheet and its latest period
		rsHistTarget:  1, // only the opening period of jobabhist000002
		rsHistRev0:    1, // only the opening period of jobabhist000003
		rsHistRev1:    0,
	} {
		var row struct {
			JobCount int `db:"job_count"`
		}
		if err := app.DB().NewQuery("SELECT job_count FROM rate_sheets_augmented WHERE id = {:id}").Bind(dbx.Params{"id": id}).One(&row); err != nil {
			t.Fatal(err)
		}
		if row.JobCount != want {
			t.Errorf("%s job_count = %d; want %d", id, row.JobCount, want)
		}
	}
}
