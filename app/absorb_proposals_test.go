package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"tybalt/internal/testutils"
	"tybalt/routes"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

// Fixtures (absorb_proposals.csv, absorb_proposal_records.csv):
//   - absorbprop00001: pending, BN Rail into Mobilia SA, by time@test.com
//   - absorbprop00002: pending, Sepulchi & Co into City of Ipperville, by time@test.com
//   - absorbprop00003: pending, Stuff Unlimited into Big Vendor Industries, by time@test.com
//   - absorbprop00004: declined vendor proposal by fakemanager@fakesite.xyz
//
// book@keeper.com holds the absorb claim; time@test.com does not.
const (
	proposalRequesterEmail = "time@test.com"
	proposalRequesterID    = "rzr98oadsp9qc11"
	proposalHolderEmail    = "book@keeper.com"
	proposalHolderID       = "tqqf7q0f3378rvp"
	proposalOtherEmail     = "fakemanager@fakesite.xyz"

	proposalMobiliaClients  = "absorbprop00001"
	proposalSepulchiClients = "absorbprop00002"
	proposalVendors         = "absorbprop00003"
	proposalDeclined        = "absorbprop00004"

	clientMobilia     = "lb0fnenkeyitsny"
	clientBNRail      = "eldtxi3i4h00k8r"
	clientSepulchi    = "pqpd90fqd5ohjcs"
	clientIpperville  = "iq7ge8osf6kuj3l"
	clientAtlantic    = "ee3xvodl583b61o"
	contactJoe        = "nh5u9z3cyknjclv"
	contactFranklin   = "nqy9f0oojeyzasa"
	contactHappy      = "contactdelok001"
	contactDavidBrown = "7lalt1egn7wi8bn"
	vendorBig         = "2zqxtsmymf670ha"
	vendorStuff       = "z66xe6vqhwtokt4"
	vendorRicky       = "mmgxrnn144767x7"
	vendorDicks       = "yxhycv2ycpvsbt4"
)

type proposalTestRow struct {
	ID         string `json:"id"`
	Status     string `json:"status"`
	TargetName string `json:"target_name"`
	ClientName string `json:"client_name"`
	InProgress bool   `json:"in_progress"`
	Records    []struct {
		RecordID   string `json:"record_id"`
		RecordName string `json:"record_name"`
	} `json:"records"`
}

func proposalJSON(t testing.TB, body any) *strings.Reader {
	t.Helper()
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	return strings.NewReader(string(data))
}

func proposalErrorCode(t testing.TB, body []byte) string {
	t.Helper()
	var failure struct {
		Code string `json:"code"`
	}
	if err := json.Unmarshal(body, &failure); err != nil {
		t.Fatalf("not a JSON error: %s", body)
	}
	return failure.Code
}

func proposalStatus(t testing.TB, app core.App, id string) string {
	t.Helper()
	proposal, err := app.FindRecordById("absorb_proposals", id)
	if err != nil {
		t.Fatal(err)
	}
	return proposal.GetString("status")
}

func listProposals(t testing.TB, app *tests.TestApp, email string, query string) []proposalTestRow {
	t.Helper()
	res := performTestAPIRequest(t, app, http.MethodGet, "/api/absorb_proposals"+query, nil, uncommitToken(t, email))
	mustStatus(t, res, http.StatusOK)
	rows := []proposalTestRow{}
	if err := json.Unmarshal(res.Body.Bytes(), &rows); err != nil {
		t.Fatal(err)
	}
	return rows
}

func proposalIDs(rows []proposalTestRow) []string {
	ids := []string{}
	for _, row := range rows {
		ids = append(ids, row.ID)
	}
	return ids
}

func TestAbsorbProposalCreate(t *testing.T) {
	app := testutils.SetupTestApp(t)

	create := func(body map[string]any) (int, []byte) {
		res := performTestAPIRequest(t, app, http.MethodPost, "/api/absorb_proposals", proposalJSON(t, body), uncommitToken(t, proposalRequesterEmail))
		return res.Code, res.Body.Bytes()
	}

	refusals := []struct {
		name   string
		body   map[string]any
		status int
		code   string
	}{
		{"rate sheets", map[string]any{"collection_name": "rate_sheets", "target_id": "a", "record_ids": []string{"b"}, "reason": "x"}, 400, "invalid_collection"},
		{"no records", map[string]any{"collection_name": "client_contacts", "target_id": contactJoe, "record_ids": []string{}, "reason": "x"}, 400, "no_records"},
		{"target in list", map[string]any{"collection_name": "client_contacts", "target_id": contactJoe, "record_ids": []string{contactJoe}, "reason": "x"}, 400, "target_in_records"},
		{"duplicate ids", map[string]any{"collection_name": "client_contacts", "target_id": contactJoe, "record_ids": []string{contactFranklin, contactFranklin}, "reason": "x"}, 400, "duplicate_records"},
		{"missing record", map[string]any{"collection_name": "client_contacts", "target_id": contactJoe, "record_ids": []string{"doesnotexist000"}, "reason": "x"}, 404, "record_not_found"},
		{"contacts of two clients", map[string]any{"collection_name": "client_contacts", "target_id": contactJoe, "record_ids": []string{contactDavidBrown}, "reason": "x"}, 400, "different_parent"},
		{"blank reason", map[string]any{"collection_name": "client_contacts", "target_id": contactJoe, "record_ids": []string{contactFranklin}, "reason": "  "}, 400, "reason_required"},
		{"reason too long", map[string]any{"collection_name": "client_contacts", "target_id": contactJoe, "record_ids": []string{contactFranklin}, "reason": strings.Repeat("x", 1001)}, 400, "reason_too_long"},
		// absorbprop00001 already absorbs BN Rail.
		{"item of another open proposal", map[string]any{"collection_name": "clients", "target_id": clientAtlantic, "record_ids": []string{clientBNRail}, "reason": "x"}, 409, "record_in_open_proposal"},
		{"target is an item of another open proposal", map[string]any{"collection_name": "clients", "target_id": clientBNRail, "record_ids": []string{clientAtlantic}, "reason": "x"}, 409, "record_in_open_proposal"},
		// absorbprop00001 keeps Mobilia, so no proposal may absorb it.
		{"item kept by another open proposal", map[string]any{"collection_name": "clients", "target_id": clientAtlantic, "record_ids": []string{clientMobilia}, "reason": "x"}, 409, "record_kept_by_open_proposal"},
	}
	for _, refusal := range refusals {
		t.Run(refusal.name, func(t *testing.T) {
			status, body := create(refusal.body)
			if status != refusal.status || proposalErrorCode(t, body) != refusal.code {
				t.Fatalf("got %d %s, want %d %s", status, body, refusal.status, refusal.code)
			}
		})
	}

	// Two proposals may keep the same record.
	status, body := create(map[string]any{"collection_name": "clients", "target_id": clientMobilia, "record_ids": []string{clientAtlantic}, "reason": "Atlantic is Mobilia"})
	if status != http.StatusOK {
		t.Fatalf("second proposal keeping Mobilia: %d %s", status, body)
	}

	status, body = create(map[string]any{"collection_name": "client_contacts", "target_id": contactJoe, "record_ids": []string{contactFranklin, contactHappy}, "reason": " Same person "})
	if status != http.StatusOK {
		t.Fatalf("create: %d %s", status, body)
	}
	var created struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(body, &created); err != nil {
		t.Fatal(err)
	}
	proposal, err := app.FindRecordById("absorb_proposals", created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if proposal.GetString("status") != "pending" || proposal.GetString("requester") != proposalRequesterID ||
		proposal.GetString("reason") != "Same person" || proposal.GetString("target_name") != "Joe Muchico (joe@example2.com)" ||
		proposal.GetString("client_id") != clientMobilia || proposal.GetString("client_name") != "Mobilia SA" {
		t.Fatalf("unexpected proposal: %v", proposal.FieldsData())
	}
	rows := listProposals(t, app, proposalRequesterEmail, "")
	var found *proposalTestRow
	for i := range rows {
		if rows[i].ID == created.ID {
			found = &rows[i]
		}
	}
	if found == nil || len(found.Records) != 2 || found.Records[0].RecordName != "Franklin Costanza (fcostanza@seinfeld.com)" {
		t.Fatalf("created proposal not listed with its records: %+v", rows)
	}

	// The open index refuses a record in two open proposals even without the
	// route's check.
	item, err := app.FindFirstRecordByData("absorb_proposal_records", "record_id", contactFranklin)
	if err != nil {
		t.Fatal(err)
	}
	collection, _ := app.FindCollectionByNameOrId("absorb_proposal_records")
	duplicate := core.NewRecord(collection)
	duplicate.Set("proposal", proposalVendors)
	duplicate.Set("collection_name", "client_contacts")
	duplicate.Set("record_id", item.GetString("record_id"))
	duplicate.Set("open", true)
	if err := app.Save(duplicate); err == nil {
		t.Fatal("expected the open index to refuse a second open row for the record")
	}
}

func TestAbsorbProposalListAndView(t *testing.T) {
	app := testutils.SetupTestApp(t)

	// A holder sees every pending proposal, oldest first.
	got := proposalIDs(listProposals(t, app, proposalHolderEmail, ""))
	want := []string{proposalMobiliaClients, proposalSepulchiClients, proposalVendors}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("holder pending: got %v, want %v", got, want)
	}
	if got := proposalIDs(listProposals(t, app, proposalHolderEmail, "?status=decided")); strings.Join(got, ",") != proposalDeclined {
		t.Fatalf("holder decided: got %v", got)
	}
	// Anyone else sees only their own.
	if got := proposalIDs(listProposals(t, app, proposalRequesterEmail, "?status=decided")); len(got) != 0 {
		t.Fatalf("requester decided: got %v", got)
	}
	if got := proposalIDs(listProposals(t, app, proposalOtherEmail, "?status=decided")); strings.Join(got, ",") != proposalDeclined {
		t.Fatalf("other decided: got %v", got)
	}
	if got := proposalIDs(listProposals(t, app, proposalOtherEmail, "")); len(got) != 0 {
		t.Fatalf("other pending: got %v", got)
	}

	res := performTestAPIRequest(t, app, http.MethodGet, "/api/absorb_proposals/"+proposalVendors, nil, uncommitToken(t, proposalOtherEmail))
	mustStatus(t, res, http.StatusNotFound)
	res = performTestAPIRequest(t, app, http.MethodGet, "/api/absorb_proposals/"+proposalVendors, nil, uncommitToken(t, proposalHolderEmail))
	mustStatus(t, res, http.StatusOK)
	var row proposalTestRow
	if err := json.Unmarshal(res.Body.Bytes(), &row); err != nil {
		t.Fatal(err)
	}
	if row.TargetName != "Big Vendor Industries (BVI)" || len(row.Records) != 1 || row.Records[0].RecordID != vendorStuff {
		t.Fatalf("unexpected proposal: %+v", row)
	}

	res = performTestAPIRequest(t, app, http.MethodGet, "/api/absorb_proposals", nil, nil)
	mustStatus(t, res, http.StatusUnauthorized)

	// The menu badge counts pending proposals for holders only.
	badges := func(email string) map[string]int {
		res := performTestAPIRequest(t, app, http.MethodGet, "/api/nav/badges", nil, uncommitToken(t, email))
		mustStatus(t, res, http.StatusOK)
		counts := map[string]int{}
		if err := json.Unmarshal(res.Body.Bytes(), &counts); err != nil {
			t.Fatal(err)
		}
		return counts
	}
	if count := badges(proposalHolderEmail)["/absorb/proposals"]; count != 3 {
		t.Fatalf("holder badge: got %d, want 3", count)
	}
	if _, ok := badges(proposalRequesterEmail)["/absorb/proposals"]; ok {
		t.Fatal("a user without the absorb claim got the proposals badge")
	}
}

func TestAbsorbProposalWithdrawAndDecline(t *testing.T) {
	app := testutils.SetupTestApp(t)
	post := func(email, path string, body any) (int, []byte) {
		res := performTestAPIRequest(t, app, http.MethodPost, path, proposalJSON(t, body), uncommitToken(t, email))
		return res.Code, res.Body.Bytes()
	}

	// Only the requester withdraws.
	if status, body := post(proposalHolderEmail, "/api/absorb_proposals/"+proposalVendors+"/withdraw", map[string]any{}); status != http.StatusForbidden {
		t.Fatalf("holder withdraw: %d %s", status, body)
	}
	// Only a holder declines, and only with a comment.
	if status, body := post(proposalRequesterEmail, "/api/absorb_proposals/"+proposalMobiliaClients+"/decline", map[string]any{"comment": "no"}); status != http.StatusForbidden {
		t.Fatalf("requester decline: %d %s", status, body)
	}
	if status, body := post(proposalHolderEmail, "/api/absorb_proposals/"+proposalMobiliaClients+"/decline", map[string]any{"comment": " "}); status != http.StatusBadRequest || proposalErrorCode(t, body) != "comment_required" {
		t.Fatalf("blank comment: %d %s", status, body)
	}

	if status, body := post(proposalRequesterEmail, "/api/absorb_proposals/"+proposalVendors+"/withdraw", map[string]any{}); status != http.StatusOK {
		t.Fatalf("withdraw: %d %s", status, body)
	}
	if got := proposalStatus(t, app, proposalVendors); got != "withdrawn" {
		t.Fatalf("status after withdraw: %s", got)
	}
	// The withdrawn proposal frees its record for a new one.
	if status, body := post(proposalRequesterEmail, "/api/absorb_proposals", map[string]any{"collection_name": "vendors", "target_id": vendorRicky, "record_ids": []string{vendorStuff}, "reason": "again"}); status != http.StatusOK {
		t.Fatalf("new proposal for a freed record: %d %s", status, body)
	}
	if status, body := post(proposalRequesterEmail, "/api/absorb_proposals/"+proposalVendors+"/withdraw", map[string]any{}); status != http.StatusBadRequest || proposalErrorCode(t, body) != "proposal_not_pending" {
		t.Fatalf("second withdraw: %d %s", status, body)
	}

	if status, body := post(proposalHolderEmail, "/api/absorb_proposals/"+proposalMobiliaClients+"/decline", map[string]any{"comment": "Different companies"}); status != http.StatusOK {
		t.Fatalf("decline: %d %s", status, body)
	}
	proposal, err := app.FindRecordById("absorb_proposals", proposalMobiliaClients)
	if err != nil {
		t.Fatal(err)
	}
	if proposal.GetString("status") != "declined" || proposal.GetString("decided_by") != proposalHolderID ||
		proposal.GetString("decision_comment") != "Different companies" || proposal.GetDateTime("decided_at").IsZero() {
		t.Fatalf("unexpected declined proposal: %v", proposal.FieldsData())
	}
	if n := countNotifications(t, app, "absorb_proposal_declined", proposalRequesterID); n != 1 {
		t.Fatalf("declined notifications: got %d, want 1", n)
	}
	// Absorbing from a declined proposal is refused.
	res := performTestAPIRequest(t, app, http.MethodPost, "/api/clients/"+clientMobilia+"/absorb",
		proposalJSON(t, map[string]any{"ids_to_absorb": []string{clientBNRail}, "proposal": proposalMobiliaClients}), uncommitToken(t, proposalHolderEmail))
	if res.Code != http.StatusBadRequest || proposalErrorCode(t, res.Body.Bytes()) != "proposal_not_pending" {
		t.Fatalf("absorb from a declined proposal: %d %s", res.Code, res.Body.String())
	}
}

func TestAbsorbProposalAcceptedOnCommit(t *testing.T) {
	app := testutils.SetupTestApp(t)
	holder := uncommitToken(t, proposalHolderEmail)
	absorbFrom := func(email, target string, ids []string, proposal string) (int, []byte) {
		res := performTestAPIRequest(t, app, http.MethodPost, "/api/clients/"+target+"/absorb",
			proposalJSON(t, map[string]any{"ids_to_absorb": ids, "proposal": proposal}), uncommitToken(t, email))
		return res.Code, res.Body.Bytes()
	}

	// Only a holder absorbs, proposal or not.
	if status, body := absorbFrom(proposalRequesterEmail, clientMobilia, []string{clientBNRail}, proposalMobiliaClients); status != http.StatusForbidden {
		t.Fatalf("requester absorb: %d %s", status, body)
	}
	// The records must come from the proposal.
	if status, body := absorbFrom(proposalHolderEmail, clientMobilia, []string{clientBNRail, clientAtlantic}, proposalMobiliaClients); status != http.StatusBadRequest || proposalErrorCode(t, body) != "record_not_in_proposal" {
		t.Fatalf("extra record: %d %s", status, body)
	}
	// A proposal of one collection can't carry out an absorb of another.
	res := performTestAPIRequest(t, app, http.MethodPost, "/api/vendors/"+vendorBig+"/absorb",
		proposalJSON(t, map[string]any{"ids_to_absorb": []string{vendorStuff}, "proposal": proposalMobiliaClients}), holder)
	if res.Code != http.StatusBadRequest || proposalErrorCode(t, res.Body.Bytes()) != "proposal_mismatch" {
		t.Fatalf("mismatched collection: %d %s", res.Code, res.Body.String())
	}

	// The reviewer keeps BN Rail instead of Mobilia: a swap.
	if status, body := absorbFrom(proposalHolderEmail, clientBNRail, []string{clientMobilia}, proposalMobiliaClients); status != http.StatusOK {
		t.Fatalf("absorb from proposal: %d %s", status, body)
	}
	action, err := app.FindFirstRecordByData("absorb_actions", "collection_name", "clients")
	if err != nil {
		t.Fatal(err)
	}
	if action.GetString("proposal") != proposalMobiliaClients {
		t.Fatalf("action proposal: got %q", action.GetString("proposal"))
	}
	// Absorbed, not committed: the proposal is still pending, and can be
	// neither withdrawn nor declined.
	if got := proposalStatus(t, app, proposalMobiliaClients); got != "pending" {
		t.Fatalf("status after absorb: %s", got)
	}
	rows := listProposals(t, app, proposalHolderEmail, "")
	if rows[0].ID != proposalMobiliaClients || !rows[0].InProgress {
		t.Fatalf("expected the proposal to show as in progress: %+v", rows[0])
	}
	res = performTestAPIRequest(t, app, http.MethodPost, "/api/absorb_proposals/"+proposalMobiliaClients+"/withdraw", nil, uncommitToken(t, proposalRequesterEmail))
	if res.Code != http.StatusConflict {
		t.Fatalf("withdraw while absorbed: %d %s", res.Code, res.Body.String())
	}
	if n := countNotifications(t, app, "absorb_proposal_accepted", proposalRequesterID); n != 0 {
		t.Fatalf("accepted notification before commit: %d", n)
	}

	// Commit accepts it and tells the requester.
	res = performTestAPIRequest(t, app, http.MethodDelete, "/api/collections/absorb_actions/records/"+action.Id, nil, holder)
	mustStatus(t, res, http.StatusNoContent)
	proposal, err := app.FindRecordById("absorb_proposals", proposalMobiliaClients)
	if err != nil {
		t.Fatal(err)
	}
	if proposal.GetString("status") != "accepted" || proposal.GetString("decided_by") != proposalHolderID || proposal.GetDateTime("decided_at").IsZero() {
		t.Fatalf("unexpected accepted proposal: %v", proposal.FieldsData())
	}
	if n := countNotifications(t, app, "absorb_proposal_accepted", proposalRequesterID); n != 1 {
		t.Fatalf("accepted notifications: got %d, want 1", n)
	}
	// Its records are free again.
	item, err := app.FindFirstRecordByData("absorb_proposal_records", "proposal", proposalMobiliaClients)
	if err != nil {
		t.Fatal(err)
	}
	if item.GetBool("open") {
		t.Fatal("the accepted proposal still holds its record")
	}
	// The other pending proposals name no absorbed record.
	if got := proposalStatus(t, app, proposalSepulchiClients); got != "pending" {
		t.Fatalf("unrelated proposal: %s", got)
	}
}

func TestAbsorbProposalPendingAfterUndo(t *testing.T) {
	app := testutils.SetupTestApp(t)
	holder := uncommitToken(t, proposalHolderEmail)

	res := performTestAPIRequest(t, app, http.MethodPost, "/api/vendors/"+vendorBig+"/absorb",
		proposalJSON(t, map[string]any{"ids_to_absorb": []string{vendorStuff}, "proposal": proposalVendors}), holder)
	mustStatus(t, res, http.StatusOK)
	res = performTestAPIRequest(t, app, http.MethodPost, "/api/vendors/undo_absorb", nil, holder)
	mustStatus(t, res, http.StatusOK)

	if got := proposalStatus(t, app, proposalVendors); got != "pending" {
		t.Fatalf("status after undo: %s", got)
	}
	if _, err := app.FindRecordById("vendors", vendorStuff); err != nil {
		t.Fatalf("undo did not restore the vendor: %v", err)
	}
	// The proposal can be carried out again, and then withdrawn is refused.
	res = performTestAPIRequest(t, app, http.MethodPost, "/api/vendors/"+vendorBig+"/absorb",
		proposalJSON(t, map[string]any{"ids_to_absorb": []string{vendorStuff}, "proposal": proposalVendors}), holder)
	mustStatus(t, res, http.StatusOK)
	if n := countNotifications(t, app, "absorb_proposal_accepted", proposalRequesterID); n != 0 {
		t.Fatalf("accepted notifications without a commit: %d", n)
	}
}

func TestAbsorbProposalObsoleteOnCommit(t *testing.T) {
	app := testutils.SetupTestApp(t)

	// A plain absorb of Sepulchi, which absorbprop00002 proposed to absorb
	// into another client. Then a plain absorb of the record that
	// absorbprop00001 keeps would make it obsolete too; it is checked after.
	if err := routes.AbsorbRecords(app, "clients", clientAtlantic, []string{clientSepulchi}); err != nil {
		t.Fatal(err)
	}
	if got := proposalStatus(t, app, proposalSepulchiClients); got != "pending" {
		t.Fatalf("status before commit: %s", got)
	}
	action, err := app.FindFirstRecordByData("absorb_actions", "collection_name", "clients")
	if err != nil {
		t.Fatal(err)
	}
	res := performTestAPIRequest(t, app, http.MethodDelete, "/api/collections/absorb_actions/records/"+action.Id, nil, uncommitToken(t, proposalHolderEmail))
	mustStatus(t, res, http.StatusNoContent)

	proposal, err := app.FindRecordById("absorb_proposals", proposalSepulchiClients)
	if err != nil {
		t.Fatal(err)
	}
	if proposal.GetString("status") != "obsolete" || !strings.Contains(proposal.GetString("decision_comment"), "Sepulchi & Co") {
		t.Fatalf("unexpected obsolete proposal: %v", proposal.FieldsData())
	}
	if n := countNotifications(t, app, "absorb_proposal_obsolete", proposalRequesterID); n != 1 {
		t.Fatalf("obsolete notifications: got %d, want 1", n)
	}
	if got := proposalStatus(t, app, proposalMobiliaClients); got != "pending" {
		t.Fatalf("unrelated proposal: %s", got)
	}

	// A commit that absorbs the target of a proposal closes it too.
	if err := routes.AbsorbRecords(app, "clients", clientAtlantic, []string{clientMobilia}); err != nil {
		t.Fatal(err)
	}
	action, err = app.FindFirstRecordByData("absorb_actions", "collection_name", "clients")
	if err != nil {
		t.Fatal(err)
	}
	res = performTestAPIRequest(t, app, http.MethodDelete, "/api/collections/absorb_actions/records/"+action.Id, nil, uncommitToken(t, proposalHolderEmail))
	mustStatus(t, res, http.StatusNoContent)
	if got := proposalStatus(t, app, proposalMobiliaClients); got != "obsolete" {
		t.Fatalf("proposal whose target was absorbed: %s", got)
	}
}

func TestAbsorbProposalContactsFollowClientCommit(t *testing.T) {
	app := testutils.SetupTestApp(t)

	// A contact proposal under Mobilia, then Mobilia is absorbed into
	// Atlantic. absorbprop00001 keeps Mobilia, so it becomes obsolete, and
	// the contact proposal names Atlantic as its client.
	res := performTestAPIRequest(t, app, http.MethodPost, "/api/absorb_proposals",
		proposalJSON(t, map[string]any{"collection_name": "client_contacts", "target_id": contactJoe, "record_ids": []string{contactFranklin}, "reason": "Same person"}),
		uncommitToken(t, proposalRequesterEmail))
	mustStatus(t, res, http.StatusOK)
	var created struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}

	if err := routes.AbsorbRecords(app, "clients", clientAtlantic, []string{clientMobilia}); err != nil {
		t.Fatal(err)
	}
	action, err := app.FindFirstRecordByData("absorb_actions", "collection_name", "clients")
	if err != nil {
		t.Fatal(err)
	}
	res = performTestAPIRequest(t, app, http.MethodDelete, "/api/collections/absorb_actions/records/"+action.Id, nil, uncommitToken(t, proposalHolderEmail))
	mustStatus(t, res, http.StatusNoContent)

	proposal, err := app.FindRecordById("absorb_proposals", created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if proposal.GetString("status") != "pending" || proposal.GetString("client_id") != clientAtlantic || proposal.GetString("client_name") != "Atlantic Power Generation" {
		t.Fatalf("unexpected contact proposal: %v", proposal.FieldsData())
	}
}

func TestAbsorbProposalRateSheetsRefused(t *testing.T) {
	app := testutils.SetupTestApp(t)
	// A rate sheet absorb can't carry out a proposal, and no rate sheet
	// proposal can be made (TestAbsorbProposalCreate).
	res := performTestAPIRequest(t, app, http.MethodPost, "/api/rate_sheets/"+rsAbTarget+"/absorb",
		proposalJSON(t, map[string]any{"ids_to_absorb": []string{rsAbSingle}, "proposal": proposalMobiliaClients}), uncommitToken(t, absorbAdminEmail))
	if res.Code != http.StatusBadRequest || proposalErrorCode(t, res.Body.Bytes()) != "proposal_mismatch" {
		t.Fatalf("rate sheet absorb with a proposal: %d %s", res.Code, res.Body.String())
	}
	if _, err := app.FindFirstRecordByData("absorb_actions", "collection_name", "rate_sheets"); err == nil {
		t.Fatal("a rate sheet absorb ran with a proposal")
	}
}
