package main

import (
	"encoding/json"
	"io"
	"net/http"
	"sort"
	"testing"
	"tybalt/internal/testutils"
	"tybalt/routes"

	"github.com/pocketbase/pocketbase/tests"
)

// Fixtures for absorb suggestions (clients.csv, vendors.csv,
// absorb_dismissals.csv), with the proposal fixtures described in
// absorb_proposals_test.go:
//   - sugclient000001 BN Rail Corp. duplicates BN Rail Corporation, which
//     open proposal absorbprop00001 absorbs and 36 jobs use.
//   - sugclient000002/3 Kestrel Surveying (Ltd): dismissed together
//     (absorbdism00001), so the bucket is hidden.
//   - sugclient000004/5/6 Pinecrest Paving (Inc.), Pine Crest Paving: 4 and 5
//     are dismissed (absorbdism00002), but 6 joined them, so it shows.
//   - sugvendor000001 Stuff Unlimited Inc. duplicates Stuff Unlimited, which
//     open proposal absorbprop00003 absorbs.
//   - sugvendor000002 Big Vendor Industries Ltd duplicates Big Vendor
//     Industries, which absorbprop00003 keeps.
const (
	sugBNRailCorp       = "sugclient000001"
	sugKestrel          = "sugclient000002"
	sugKestrelLtd       = "sugclient000003"
	sugPinecrest        = "sugclient000004"
	sugPinecrestInc     = "sugclient000005"
	sugPineCrest        = "sugclient000006"
	sugStuffUnlimitedCo = "sugvendor000001"
	sugBigVendorLtd     = "sugvendor000002"
)

type suggestionTestMember struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	References int     `json:"references"`
	Score      float64 `json:"score"`
	Proposal   string  `json:"proposal"`
}

type suggestionTestBucket struct {
	Key      string                 `json:"key"`
	TargetID string                 `json:"target_id"`
	Members  []suggestionTestMember `json:"members"`
}

type suggestionTestResponse struct {
	Buckets       []suggestionTestBucket `json:"buckets"`
	Dismissed     int                    `json:"dismissed"`
	PendingAction *struct {
		TargetID   string   `json:"target_id"`
		TargetName string   `json:"target_name"`
		Absorbed   []string `json:"absorbed"`
		Created    string   `json:"created"`
	} `json:"pending_action"`
}

func getSuggestions(t testing.TB, app *tests.TestApp, collection string, query string) suggestionTestResponse {
	t.Helper()
	res := performTestAPIRequest(t, app, http.MethodGet, "/api/absorb_suggestions/"+collection+query, nil, uncommitToken(t, proposalHolderEmail))
	mustStatus(t, res, http.StatusOK)
	var out suggestionTestResponse
	if err := json.Unmarshal(res.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	return out
}

// bucketWith returns the bucket that holds id, or nil.
func bucketWith(response suggestionTestResponse, id string) *suggestionTestBucket {
	for i, bucket := range response.Buckets {
		for _, member := range bucket.Members {
			if member.ID == id {
				return &response.Buckets[i]
			}
		}
	}
	return nil
}

func memberIDs(bucket *suggestionTestBucket) []string {
	ids := []string{}
	for _, member := range bucket.Members {
		ids = append(ids, member.ID)
	}
	sort.Strings(ids)
	return ids
}

func TestAbsorbSuggestionsAccess(t *testing.T) {
	app := testutils.SetupTestApp(t)
	cases := []struct {
		name   string
		method string
		url    string
		email  string
		body   any
		status int
		code   string
	}{
		{"list without the absorb claim", http.MethodGet, "/api/absorb_suggestions/clients", proposalRequesterEmail, nil, 403, "unauthorized"},
		{"dismiss without the absorb claim", http.MethodPost, "/api/absorb_suggestions/vendors/dismiss", proposalRequesterEmail, map[string]any{"record_ids": []string{vendorStuff, sugStuffUnlimitedCo}}, 403, "unauthorized"},
		{"contacts", http.MethodGet, "/api/absorb_suggestions/client_contacts", proposalHolderEmail, nil, 400, "invalid_collection"},
		{"rate sheets", http.MethodGet, "/api/absorb_suggestions/rate_sheets", proposalHolderEmail, nil, 400, "invalid_collection"},
		{"unknown level", http.MethodGet, "/api/absorb_suggestions/clients?level=tight", proposalHolderEmail, nil, 400, "invalid_level"},
		{"one record", http.MethodPost, "/api/absorb_suggestions/vendors/dismiss", proposalHolderEmail, map[string]any{"record_ids": []string{vendorStuff}}, 400, "too_few_records"},
		{"a record twice", http.MethodPost, "/api/absorb_suggestions/vendors/dismiss", proposalHolderEmail, map[string]any{"record_ids": []string{vendorStuff, vendorStuff}}, 400, "duplicate_records"},
		{"a missing record", http.MethodPost, "/api/absorb_suggestions/vendors/dismiss", proposalHolderEmail, map[string]any{"record_ids": []string{vendorStuff, "doesnotexist000"}}, 404, "record_not_found"},
		{"a record of another collection", http.MethodPost, "/api/absorb_suggestions/vendors/dismiss", proposalHolderEmail, map[string]any{"record_ids": []string{vendorStuff, clientMobilia}}, 404, "record_not_found"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var body io.Reader
			if c.body != nil {
				body = proposalJSON(t, c.body)
			}
			res := performTestAPIRequest(t, app, c.method, c.url, body, uncommitToken(t, c.email))
			if res.Code != c.status || proposalErrorCode(t, res.Body.Bytes()) != c.code {
				t.Fatalf("got %d %s, want %d %s", res.Code, res.Body.String(), c.status, c.code)
			}
		})
	}
	res := performTestAPIRequest(t, app, http.MethodGet, "/api/absorb_suggestions/clients", nil, nil)
	if res.Code != http.StatusUnauthorized {
		t.Fatalf("signed out: got %d", res.Code)
	}
}

func TestAbsorbSuggestionsBuckets(t *testing.T) {
	app := testutils.SetupTestApp(t)

	clients := getSuggestions(t, app, "clients", "")
	bnRail := bucketWith(clients, sugBNRailCorp)
	if bnRail == nil {
		t.Fatalf("no bucket for BN Rail Corp.: %+v", clients.Buckets)
	}
	// The suggested target is the member with the most references: BN Rail
	// Corporation, with its 36 jobs. It is in an open proposal.
	target := bnRail.Members[0]
	if bnRail.TargetID != clientBNRail || bnRail.Key != clientBNRail || target.ID != clientBNRail || target.Score != 1 {
		t.Fatalf("BN Rail target = %+v", bnRail)
	}
	if target.References != 36 {
		t.Errorf("BN Rail Corporation references = %d, want 36", target.References)
	}
	if target.Proposal != proposalMobiliaClients {
		t.Errorf("BN Rail Corporation proposal = %q, want %s", target.Proposal, proposalMobiliaClients)
	}
	if other := bnRail.Members[1]; other.ID != sugBNRailCorp || other.Proposal != "" || other.References != 0 {
		t.Errorf("BN Rail Corp. = %+v", other)
	}

	// A dismissed bucket is hidden; one that a new record joined shows.
	if bucketWith(clients, sugKestrel) != nil || bucketWith(clients, sugKestrelLtd) != nil {
		t.Error("dismissed Kestrel bucket is shown")
	}
	if clients.Dismissed < 1 {
		t.Errorf("dismissed = %d, want at least 1", clients.Dismissed)
	}
	pinecrest := bucketWith(clients, sugPineCrest)
	if pinecrest == nil || len(pinecrest.Members) != 3 {
		t.Fatalf("Pinecrest bucket = %+v, want three members", pinecrest)
	}
	if clients.PendingAction != nil {
		t.Errorf("pending action = %+v, want none", clients.PendingAction)
	}

	vendors := getSuggestions(t, app, "vendors", "?level=normal")
	stuff := bucketWith(vendors, sugStuffUnlimitedCo)
	if stuff == nil || stuff.TargetID != vendorStuff || stuff.Members[0].References != 113 {
		t.Fatalf("Stuff Unlimited bucket = %+v", stuff)
	}
	// Stuff Unlimited is absorbed by an open proposal, and Big Vendor
	// Industries is kept by it: both are flagged.
	if stuff.Members[0].Proposal != proposalVendors {
		t.Errorf("Stuff Unlimited proposal = %q", stuff.Members[0].Proposal)
	}
	big := bucketWith(vendors, sugBigVendorLtd)
	if big == nil || big.TargetID != vendorBig || big.Members[0].Proposal != proposalVendors || big.Members[0].References != 52 {
		t.Fatalf("Big Vendor bucket = %+v", big)
	}

	// The same request gives the same answer.
	again := getSuggestions(t, app, "vendors", "")
	first, _ := json.Marshal(vendors)
	second, _ := json.Marshal(again)
	if string(first) != string(second) {
		t.Errorf("two requests differ:\n%s\n%s", first, second)
	}

	// The level is read from the query.
	if strict := getSuggestions(t, app, "clients", "?level=strict"); bucketWith(strict, sugBNRailCorp) == nil {
		t.Error("strict lost the BN Rail bucket")
	}
}

func TestAbsorbSuggestionsDismiss(t *testing.T) {
	app := testutils.SetupTestApp(t)
	before := getSuggestions(t, app, "vendors", "")
	if bucketWith(before, sugStuffUnlimitedCo) == nil {
		t.Fatal("no Stuff Unlimited bucket")
	}

	// This writes an absorb_dismissals row at test time: the dismissal under
	// test.
	res := performTestAPIRequest(t, app, http.MethodPost, "/api/absorb_suggestions/vendors/dismiss",
		proposalJSON(t, map[string]any{"record_ids": []string{sugStuffUnlimitedCo, vendorStuff}}), uncommitToken(t, proposalHolderEmail))
	mustStatus(t, res, http.StatusOK)

	after := getSuggestions(t, app, "vendors", "")
	if bucketWith(after, sugStuffUnlimitedCo) != nil {
		t.Error("dismissed bucket is still shown")
	}
	if after.Dismissed != before.Dismissed+1 {
		t.Errorf("dismissed = %d, want %d", after.Dismissed, before.Dismissed+1)
	}
	if bucketWith(after, sugBigVendorLtd) == nil {
		t.Error("another bucket was hidden")
	}
	// The dismissal is shared: it keeps who and when.
	row, err := app.FindFirstRecordByData("absorb_dismissals", "collection_name", "vendors")
	if err != nil {
		t.Fatal(err)
	}
	if row.GetString("dismissed_by") != proposalHolderID || row.GetDateTime("created").IsZero() {
		t.Errorf("dismissal = %v", row.FieldsData())
	}
	ids := []string{}
	if err := row.UnmarshalJSONField("record_ids", &ids); err != nil || len(ids) != 2 {
		t.Errorf("record_ids = %v %v", ids, err)
	}
}

func TestAbsorbSuggestionsPendingAction(t *testing.T) {
	app := testutils.SetupTestApp(t)
	// This absorbs BN Rail Corp. into BN Rail Corporation at test time, to
	// leave a pending absorb action for clients.
	if err := routes.AbsorbRecords(app, "clients", clientBNRail, []string{sugBNRailCorp}); err != nil {
		t.Fatal(err)
	}
	clients := getSuggestions(t, app, "clients", "")
	action := clients.PendingAction
	if action == nil {
		t.Fatal("no pending action")
	}
	if action.TargetID != clientBNRail || action.TargetName != "BN Rail Corporation" || action.Created == "" {
		t.Errorf("pending action = %+v", action)
	}
	if len(action.Absorbed) != 1 || action.Absorbed[0] != "BN Rail Corp." {
		t.Errorf("absorbed = %v", action.Absorbed)
	}
	if bucketWith(clients, clientBNRail) != nil {
		t.Error("the absorbed record still makes a bucket")
	}
	if vendors := getSuggestions(t, app, "vendors", ""); vendors.PendingAction != nil {
		t.Errorf("vendors pending action = %+v", vendors.PendingAction)
	}
	// The pinecrest bucket is unchanged by the absorb.
	if got := memberIDs(bucketWith(clients, sugPineCrest)); len(got) != 3 {
		t.Errorf("pinecrest = %v", got)
	}
}
