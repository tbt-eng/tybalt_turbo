package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"tybalt/internal/testutils"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

type workRecordListRow struct {
	ID                 string `json:"id"`
	WorkRecord         string `json:"work_record"`
	Date               string `json:"date"`
	UID                string `json:"uid"`
	OwnerName          string `json:"owner_name"`
	JobNumber          string `json:"job_number"`
	Approved           bool   `json:"approved"`
	LinkedTimeEntryID  string `json:"linked_time_entry_id"`
	SearchText         string `json:"search_text"`
	CanCreateTimeEntry bool   `json:"can_create_time_entry"`
}

type workRecordSearchResponse struct {
	Items []workRecordListRow `json:"items"`
	Total int                 `json:"total"`
}

type paginatedJobWorkRecordsResponse struct {
	Data       []workRecordListRow `json:"data"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	Total      int                 `json:"total"`
	TotalPages int                 `json:"total_pages"`
}

type workRecordNoteRow struct {
	ID         string `json:"id"`
	UID        string `json:"uid"`
	AuthorName string `json:"author_name"`
	Note       string `json:"note"`
	Created    string `json:"created"`
}

type workRecordConsumableRow struct {
	ConsumableID    string  `json:"consumable_id"`
	ConsumableName  string  `json:"consumable_name"`
	InputKind       string  `json:"input_kind"`
	UnitLabel       string  `json:"unit_label"`
	QuantityNumber  float64 `json:"quantity_number"`
	BooleanSelected bool    `json:"boolean_selected"`
}

type workRecordDetailResponse struct {
	ID                         string                    `json:"id"`
	WorkRecord                 string                    `json:"work_record"`
	Date                       string                    `json:"date"`
	UID                        string                    `json:"uid"`
	OwnerName                  string                    `json:"owner_name"`
	JobID                      string                    `json:"job_id"`
	JobNumber                  string                    `json:"job_number"`
	TypeID                     string                    `json:"type_id"`
	TypeName                   string                    `json:"type_name"`
	Location                   string                    `json:"location"`
	WorkDescription            string                    `json:"work_description"`
	HoursOnSite                float64                   `json:"hours_on_site"`
	HoursTravelTime            float64                   `json:"hours_travel_time"`
	ReportTo                   string                    `json:"report_to"`
	DistanceTravelledKM        float64                   `json:"distance_travelled_km"`
	FieldBookNumber            string                    `json:"field_book_number"`
	FieldBookPageNumber        int                       `json:"field_book_page_number"`
	SubContractor              string                    `json:"sub_contractor"`
	Equipment                  string                    `json:"equipment"`
	Supplies                   string                    `json:"supplies"`
	VehicleType                string                    `json:"vehicle_type"`
	IsPassenger                bool                      `json:"is_passenger"`
	CompanyVehicleUnitNumber   string                    `json:"company_vehicle_unit_number"`
	Approved                   bool                      `json:"approved"`
	ParentWorkRecordID         string                    `json:"parent_work_record_id"`
	ParentWorkRecord           string                    `json:"parent_work_record"`
	RootWorkRecordID           string                    `json:"root_work_record_id"`
	RootWorkRecord             string                    `json:"root_work_record"`
	LinkedTimeEntryID          string                    `json:"linked_time_entry_id"`
	LinkedTimeSheetID          string                    `json:"linked_timesheet_id"`
	LinkedTimeEntryDate        string                    `json:"linked_time_entry_date"`
	LinkedTimeEntryDescription string                    `json:"linked_time_entry_description"`
	LinkedTimeSheetWeekEnding  string                    `json:"linked_time_sheet_week_ending"`
	CanCreateTimeEntry         bool                      `json:"can_create_time_entry"`
	CanCopyToTomorrow          bool                      `json:"can_copy_to_tomorrow"`
	CanEdit                    bool                      `json:"can_edit"`
	CanAddNote                 bool                      `json:"can_add_note"`
	Notes                      []workRecordNoteRow       `json:"notes"`
	Consumables                []workRecordConsumableRow `json:"consumables"`
}

type workRecordConsumableReportRow struct {
	WorkRecord      string  `json:"work_record"`
	ConsumableID    string  `json:"consumable_id"`
	ConsumableName  string  `json:"consumable_name"`
	InputKind       string  `json:"input_kind"`
	UnitLabel       string  `json:"unit_label"`
	QuantityNumber  float64 `json:"quantity_number"`
	BooleanSelected bool    `json:"boolean_selected"`
}

type workRecordConsumablesReportResponse struct {
	Items []workRecordConsumableReportRow `json:"items"`
}

type workRecordCopyDraftResponse struct {
	UID                      string                    `json:"uid"`
	Date                     string                    `json:"date"`
	Job                      string                    `json:"job"`
	Type                     string                    `json:"type"`
	Location                 string                    `json:"location"`
	WorkDescription          string                    `json:"work_description"`
	ReportTo                 string                    `json:"report_to"`
	HoursOnSite              float64                   `json:"hours_on_site"`
	HoursTravelTime          float64                   `json:"hours_travel_time"`
	DistanceTravelledKM      float64                   `json:"distance_travelled_km"`
	FieldBookNumber          string                    `json:"field_book_number"`
	FieldBookPageNumber      int                       `json:"field_book_page_number"`
	SubContractor            string                    `json:"sub_contractor"`
	Equipment                string                    `json:"equipment"`
	Supplies                 string                    `json:"supplies"`
	VehicleType              string                    `json:"vehicle_type"`
	IsPassenger              bool                      `json:"is_passenger"`
	CompanyVehicleUnitNumber string                    `json:"company_vehicle_unit_number"`
	Consumables              []workRecordConsumableRow `json:"consumables"`
	CopyFromWorkRecordID     string                    `json:"copy_from_work_record_id"`
	SourceWorkRecord         string                    `json:"source_work_record"`
	RootWorkRecordID         string                    `json:"root_work_record_id"`
	RootWorkRecord           string                    `json:"root_work_record"`
}

type workRecordSaveResponse struct {
	ID string `json:"id"`
}

type workRecordAPIFieldError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type workRecordAPIErrorResponse struct {
	Message string                             `json:"message"`
	Data    map[string]workRecordAPIFieldError `json:"data"`
}

func countWorkRecordConsumableEntriesForTest(t testing.TB, app *tests.TestApp, workRecordID string) int {
	t.Helper()

	var row struct {
		Count int `db:"count"`
	}
	if err := app.DB().NewQuery(`
		SELECT COUNT(*) AS count
		FROM work_records_consumable_entries
		WHERE work_record = {:id}
	`).Bind(map[string]any{
		"id": workRecordID,
	}).One(&row); err != nil {
		t.Fatalf("failed counting work-record consumables: %v", err)
	}

	return row.Count
}

func countWorkRecordsForTest(t testing.TB, app *tests.TestApp) int {
	t.Helper()

	var row struct {
		Count int `db:"count"`
	}
	if err := app.DB().NewQuery(`SELECT COUNT(*) AS count FROM work_records`).One(&row); err != nil {
		t.Fatalf("failed counting work records: %v", err)
	}

	return row.Count
}

func countWorkRecordNotesForTest(t testing.TB, app *tests.TestApp, workRecordID string) int {
	t.Helper()

	var row struct {
		Count int `db:"count"`
	}
	if err := app.DB().NewQuery(`
		SELECT COUNT(*) AS count
		FROM work_records_notes
		WHERE work_record = {:id}
	`).Bind(map[string]any{
		"id": workRecordID,
	}).One(&row); err != nil {
		t.Fatalf("failed counting work record notes: %v", err)
	}

	return row.Count
}

func loadWorkRecordNoteForTest(t testing.TB, app *tests.TestApp, noteID string) workRecordNoteRow {
	t.Helper()

	var row workRecordNoteRow
	if err := app.DB().NewQuery(`
		SELECT
			id,
			COALESCE(uid, '') AS uid,
			COALESCE(note, '') AS note,
			COALESCE(created, '') AS created
		FROM work_records_notes
		WHERE id = {:id}
		LIMIT 1
	`).Bind(map[string]any{
		"id": noteID,
	}).One(&row); err != nil {
		t.Fatalf("failed loading work record note: %v", err)
	}

	return row
}

func loadWorkRecordTypeAndLocationForTest(t testing.TB, app *tests.TestApp, workRecordID string) (string, string) {
	t.Helper()

	var row struct {
		TypeID   string `db:"type"`
		Location string `db:"location"`
	}
	if err := app.DB().NewQuery(`
		SELECT
			COALESCE(type, '') AS type,
			COALESCE(location, '') AS location
		FROM work_records
		WHERE id = {:id}
	`).Bind(map[string]any{
		"id": workRecordID,
	}).One(&row); err != nil {
		t.Fatalf("failed loading work-record state: %v", err)
	}

	return row.TypeID, row.Location
}

func loadWorkRecordConsumableQuantityForTest(t testing.TB, app *tests.TestApp, workRecordID string, consumableID string) float64 {
	t.Helper()

	var row struct {
		QuantityNumber float64 `db:"quantity_number"`
	}
	if err := app.DB().NewQuery(`
		SELECT COALESCE(quantity_number, 0) AS quantity_number
		FROM work_records_consumable_entries
		WHERE work_record = {:workRecordId}
		  AND consumable = {:consumableId}
		LIMIT 1
	`).Bind(map[string]any{
		"workRecordId": workRecordID,
		"consumableId": consumableID,
	}).One(&row); err != nil {
		t.Fatalf("failed loading consumable quantity: %v", err)
	}

	return row.QuantityNumber
}

func decodeWorkRecordAPIErrorResponse(t testing.TB, res *http.Response) workRecordAPIErrorResponse {
	t.Helper()

	defer res.Body.Close()

	var payload workRecordAPIErrorResponse
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("failed decoding work record error response: %v", err)
	}

	return payload
}

func TestWorkRecordsRoutes(t *testing.T) {
	reportToken, err := testutils.GenerateRecordToken("users", "fatt@mac.com")
	if err != nil {
		t.Fatal(err)
	}
	workRecordToken, err := testutils.GenerateRecordToken("users", "author@soup.com")
	if err != nil {
		t.Fatal(err)
	}
	regularUserToken, err := testutils.GenerateRecordToken("users", "time@test.com")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []tests.ApiScenario{
		{
			Name:   "work record holder can save type change after clearing incompatible consumables",
			Method: http.MethodPost,
			URL:    "/api/work_records/wr260400000010a/save",
			Headers: map[string]string{
				"Authorization": workRecordToken,
				"Content-Type":  "application/json",
			},
			Body: strings.NewReader(`{
				"date":"2026-04-10",
				"job":"cjf0kt0defhq480",
				"type":"wrtype000000002",
				"location":"Updated location",
				"work_description":"Updated description",
				"report_to":"",
				"hours_on_site":8,
				"hours_travel_time":1,
				"distance_travelled_km":10,
				"field_book_number":"",
				"field_book_page_number":"",
				"sub_contractor":"",
				"equipment":"",
				"supplies":"",
				"vehicle_type":"company",
				"is_passenger":false,
				"company_vehicle_unit_number":"7",
				"consumables":[
					{
						"entry_id":"wrcentry0000001",
						"consumable_id":"wrconsmbl000001",
						"input_kind":"number",
						"quantity_number":0,
						"boolean_selected":false
					},
					{
						"entry_id":"wrcentry0000002",
						"consumable_id":"wrconsmbl000002",
						"input_kind":"boolean",
						"quantity_number":0,
						"boolean_selected":false
					}
				]
			}`),
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"id":"wr260400000010a"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var payload workRecordSaveResponse
				if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
					t.Fatalf("failed decoding work record save response: %v", err)
				}
				if payload.ID != "wr260400000010a" {
					t.Fatalf("unexpected save payload: %#v", payload)
				}

				typeID, location := loadWorkRecordTypeAndLocationForTest(t, app, "wr260400000010a")
				if typeID != "wrtype000000002" || location != "Updated location" {
					t.Fatalf("expected updated type/location, got type=%q location=%q", typeID, location)
				}
				if got := countWorkRecordConsumableEntriesForTest(t, app, "wr260400000010a"); got != 0 {
					t.Fatalf("expected incompatible consumables to be removed, got %d", got)
				}
			},
		},
		{
			Name:   "type change still fails when incompatible consumables remain selected",
			Method: http.MethodPost,
			URL:    "/api/work_records/wr260400000010a/save",
			Headers: map[string]string{
				"Authorization": workRecordToken,
				"Content-Type":  "application/json",
			},
			Body: strings.NewReader(`{
				"date":"2026-04-10",
				"job":"cjf0kt0defhq480",
				"type":"wrtype000000002",
				"location":"Seeded location",
				"work_description":"Seeded work description",
				"report_to":"",
				"hours_on_site":8,
				"hours_travel_time":1,
				"distance_travelled_km":10,
				"field_book_number":"",
				"field_book_page_number":"",
				"sub_contractor":"",
				"equipment":"",
				"supplies":"",
				"vehicle_type":"company",
				"is_passenger":false,
				"company_vehicle_unit_number":"7",
				"consumables":[
					{
						"entry_id":"wrcentry0000001",
						"consumable_id":"wrconsmbl000001",
						"input_kind":"number",
						"quantity_number":12,
						"boolean_selected":false
					},
					{
						"entry_id":"wrcentry0000002",
						"consumable_id":"wrconsmbl000002",
						"input_kind":"boolean",
						"quantity_number":0,
						"boolean_selected":true
					}
				]
			}`),
			ExpectedStatus: http.StatusBadRequest,
			ExpectedContent: []string{
				`"code":"incompatible_consumable_entries"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, _ *http.Response) {
				typeID, _ := loadWorkRecordTypeAndLocationForTest(t, app, "wr260400000010a")
				if typeID != "wrtype000000001" {
					t.Fatalf("expected type to remain unchanged after rejection, got %q", typeID)
				}
				if got := countWorkRecordConsumableEntriesForTest(t, app, "wr260400000010a"); got != 2 {
					t.Fatalf("expected incompatible save rejection to preserve consumables, got %d", got)
				}
			},
		},
		{
			Name:   "failed work record validation rolls back prior consumable deletions",
			Method: http.MethodPost,
			URL:    "/api/work_records/wr260400000010a/save",
			Headers: map[string]string{
				"Authorization": workRecordToken,
				"Content-Type":  "application/json",
			},
			Body: strings.NewReader(`{
				"date":"2026-04-10",
				"job":"cjf0kt0defhq480",
				"type":"wrtype000000002",
				"location":"Seeded location",
				"work_description":"Seeded work description",
				"report_to":"",
				"hours_on_site":8,
				"hours_travel_time":1,
				"distance_travelled_km":10,
				"field_book_number":"",
				"field_book_page_number":"",
				"sub_contractor":"",
				"equipment":"",
				"supplies":"",
				"vehicle_type":"company",
				"is_passenger":false,
				"company_vehicle_unit_number":"",
				"consumables":[
					{
						"entry_id":"wrcentry0000001",
						"consumable_id":"wrconsmbl000001",
						"input_kind":"number",
						"quantity_number":0,
						"boolean_selected":false
					},
					{
						"entry_id":"wrcentry0000002",
						"consumable_id":"wrconsmbl000002",
						"input_kind":"boolean",
						"quantity_number":0,
						"boolean_selected":false
					}
				]
			}`),
			ExpectedStatus: http.StatusBadRequest,
			ExpectedContent: []string{
				`"code":"required_for_company_vehicle"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, _ *http.Response) {
				typeID, _ := loadWorkRecordTypeAndLocationForTest(t, app, "wr260400000010a")
				if typeID != "wrtype000000001" {
					t.Fatalf("expected failed save to leave type unchanged, got %q", typeID)
				}
				if got := countWorkRecordConsumableEntriesForTest(t, app, "wr260400000010a"); got != 2 {
					t.Fatalf("expected rollback to restore original consumables, got %d", got)
				}
			},
		},
		{
			Name:   "approved work records reject save before mutating consumables",
			Method: http.MethodPost,
			URL:    "/api/work_records/wr260400000010a/save",
			Headers: map[string]string{
				"Authorization": workRecordToken,
				"Content-Type":  "application/json",
			},
			Body: strings.NewReader(`{
				"date":"2026-04-10",
				"job":"cjf0kt0defhq480",
				"type":"wrtype000000002",
				"location":"Seeded location",
				"work_description":"Seeded work description",
				"report_to":"",
				"hours_on_site":8,
				"hours_travel_time":1,
				"distance_travelled_km":10,
				"field_book_number":"",
				"field_book_page_number":"",
				"sub_contractor":"",
				"equipment":"",
				"supplies":"",
				"vehicle_type":"company",
				"is_passenger":false,
				"company_vehicle_unit_number":"7",
				"consumables":[
					{
						"entry_id":"wrcentry0000001",
						"consumable_id":"wrconsmbl000001",
						"input_kind":"number",
						"quantity_number":0,
						"boolean_selected":false
					},
					{
						"entry_id":"wrcentry0000002",
						"consumable_id":"wrconsmbl000002",
						"input_kind":"boolean",
						"quantity_number":0,
						"boolean_selected":false
					}
				]
			}`),
			ExpectedStatus: http.StatusBadRequest,
			ExpectedContent: []string{
				`"code":"approved_work_record_locked"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, _ *core.ServeEvent) {
				if _, err := app.DB().NewQuery(`
					UPDATE work_records
					SET approved = 1
					WHERE id = {:id}
				`).Bind(map[string]any{
					"id": "wr260400000010a",
				}).Execute(); err != nil {
					t.Fatalf("failed approving seeded work record: %v", err)
				}
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, _ *http.Response) {
				typeID, _ := loadWorkRecordTypeAndLocationForTest(t, app, "wr260400000010a")
				if typeID != "wrtype000000001" {
					t.Fatalf("expected approved save rejection to leave type unchanged, got %q", typeID)
				}
				if got := countWorkRecordConsumableEntriesForTest(t, app, "wr260400000010a"); got != 2 {
					t.Fatalf("expected approved save rejection to leave consumables untouched, got %d", got)
				}
			},
		},
		{
			Name:   "non-type edits keep unchanged consumables through save endpoint",
			Method: http.MethodPost,
			URL:    "/api/work_records/wr260400000010a/save",
			Headers: map[string]string{
				"Authorization": workRecordToken,
				"Content-Type":  "application/json",
			},
			Body: strings.NewReader(`{
				"date":"2026-04-10",
				"job":"cjf0kt0defhq480",
				"type":"wrtype000000001",
				"location":"Regression location",
				"work_description":"Regression description",
				"report_to":"",
				"hours_on_site":8,
				"hours_travel_time":1,
				"distance_travelled_km":10,
				"field_book_number":"",
				"field_book_page_number":"",
				"sub_contractor":"",
				"equipment":"",
				"supplies":"",
				"vehicle_type":"company",
				"is_passenger":false,
				"company_vehicle_unit_number":"7",
				"consumables":[
					{
						"entry_id":"wrcentry0000001",
						"consumable_id":"wrconsmbl000001",
						"input_kind":"number",
						"quantity_number":12,
						"boolean_selected":false
					},
					{
						"entry_id":"wrcentry0000002",
						"consumable_id":"wrconsmbl000002",
						"input_kind":"boolean",
						"quantity_number":0,
						"boolean_selected":true
					}
				]
			}`),
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"id":"wr260400000010a"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, _ *http.Response) {
				typeID, location := loadWorkRecordTypeAndLocationForTest(t, app, "wr260400000010a")
				if typeID != "wrtype000000001" || location != "Regression location" {
					t.Fatalf("expected non-type save to update location only, got type=%q location=%q", typeID, location)
				}
				if got := countWorkRecordConsumableEntriesForTest(t, app, "wr260400000010a"); got != 2 {
					t.Fatalf("expected unchanged consumables to remain present, got %d", got)
				}
			},
		},
		{
			Name:   "report holder gets empty current user work records list",
			Method: http.MethodGet,
			URL:    "/api/work_records",
			Headers: map[string]string{
				"Authorization": reportToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"total":0`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, _ *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var payload workRecordSearchResponse
				if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
					t.Fatalf("failed decoding work records list response: %v", err)
				}

				if payload.Total != 0 || len(payload.Items) != 0 {
					t.Fatalf("expected report holder current-user list to be empty, got %#v", payload)
				}
			},
		},
		{
			Name:   "work record holder can fetch own work records list with create actions",
			Method: http.MethodGet,
			URL:    "/api/work_records",
			Headers: map[string]string{
				"Authorization": workRecordToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"work_record":"W2604-0010"`,
				`"total":3`,
				`"can_create_time_entry":true`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, _ *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var payload workRecordSearchResponse
				if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
					t.Fatalf("failed decoding work records list response: %v", err)
				}

				rows := payload.Items
				if payload.Total != 3 || len(rows) != 3 {
					t.Fatalf("expected 3 visible work records, got %#v", payload)
				}
				if rows[0].WorkRecord != "W3002-0050" || rows[1].WorkRecord != "W2604-0010" || rows[2].WorkRecord != "W2604-0003" {
					t.Fatalf("unexpected row order: %#v", rows)
				}
				if rows[0].OwnerName == "" || rows[0].SearchText == "" {
					t.Fatalf("expected owner/search text on owned row: %#v", rows[0])
				}
				if !rows[0].CanCreateTimeEntry || rows[0].LinkedTimeEntryID != "" {
					t.Fatalf("expected top row to allow time-entry creation: %#v", rows[0])
				}
			},
		},
		{
			Name:   "regular user gets empty work records list",
			Method: http.MethodGet,
			URL:    "/api/work_records",
			Headers: map[string]string{
				"Authorization": regularUserToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"total":0`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, _ *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var payload workRecordSearchResponse
				if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
					t.Fatalf("failed decoding work records list response: %v", err)
				}

				if payload.Total != 0 || len(payload.Items) != 0 {
					t.Fatalf("expected empty work record list, got %#v", payload)
				}
			},
		},
		{
			Name:   "work record holder cannot request all work records",
			Method: http.MethodGet,
			URL:    "/api/work_records?scope=all",
			Headers: map[string]string{
				"Authorization": workRecordToken,
			},
			ExpectedStatus: http.StatusForbidden,
			ExpectedContent: []string{
				`"you are not authorized to view all work records"`,
			},
			TestAppFactory: testutils.SetupTestApp,
		},
		{
			Name:   "work record creator can see records they created for someone else",
			Method: http.MethodGet,
			URL:    "/api/work_records?scope=mine",
			Headers: map[string]string{
				"Authorization": workRecordToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"total":4`,
				`"uid":"rzr98oadsp9qc11"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, _ *core.ServeEvent) {
				if _, err := app.DB().NewQuery(`
					INSERT INTO work_records (
						id,
						number,
						date,
						uid,
						creator,
						job,
						type,
						location,
						work_description,
						hours_on_site,
						hours_travel_time,
						distance_travelled_km,
						vehicle_type,
						approved
					) VALUES (
						{:id},
						{:number},
						{:date},
						{:uid},
						{:creator},
						{:job},
						{:type},
						{:location},
						{:workDescription},
						{:hoursOnSite},
						{:hoursTravelTime},
						{:distanceTravelledKM},
						{:vehicleType},
						0
					)
				`).Bind(map[string]any{
					"id":                  "wrcreator000001a",
					"number":              "W2604-0012",
					"date":                "2026-04-12",
					"uid":                 "rzr98oadsp9qc11",
					"creator":             "f2j5a8vk006baub",
					"job":                 "cjf0kt0defhq480",
					"type":                "wrtype000000001",
					"location":            "Creator visibility test",
					"workDescription":     "Created for another worker",
					"hoursOnSite":         8,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 5,
					"vehicleType":         "personal",
				}).Execute(); err != nil {
					t.Fatalf("failed creating creator-visible work record: %v", err)
				}
			},
			AfterTestFunc: func(t testing.TB, _ *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var payload workRecordSearchResponse
				if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
					t.Fatalf("failed decoding creator-visible work record response: %v", err)
				}
				if payload.Total != 4 || len(payload.Items) != 4 {
					t.Fatalf("expected 4 visible work records including the created-for-another-worker row, got %#v", payload)
				}
				var creatorVisible *workRecordListRow
				for i := range payload.Items {
					if payload.Items[i].UID == "rzr98oadsp9qc11" {
						creatorVisible = &payload.Items[i]
						break
					}
				}
				if creatorVisible == nil {
					t.Fatalf("expected creator-visible row for another worker, got %#v", payload.Items)
				}
				if creatorVisible.CanCreateTimeEntry {
					t.Fatalf("creator-visible row for another worker should not expose create action: %#v", creatorVisible)
				}
			},
		},
		{
			Name:   "report holder can view all organization work records through unified endpoint",
			Method: http.MethodGet,
			URL:    "/api/work_records?scope=all",
			Headers: map[string]string{
				"Authorization": reportToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"total":3`,
				`"work_record":"W3002-0050"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, _ *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var payload workRecordSearchResponse
				if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
					t.Fatalf("failed decoding work record list response: %v", err)
				}
				if payload.Total != 3 || len(payload.Items) != 3 {
					t.Fatalf("expected 3 organization work records without endpoint-level filtering, got %#v", payload)
				}
				if payload.Items[0].CanCreateTimeEntry {
					t.Fatalf("report all-scope should not expose create action for another user's row: %#v", payload.Items[0])
				}
			},
		},
		{
			Name:   "report holder can fetch work record details by number",
			Method: http.MethodGet,
			URL:    "/api/work_records/W2604-0010",
			Headers: map[string]string{
				"Authorization": reportToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"work_record":"W2604-0010"`,
				`"job_number":"24-321"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, _ *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var row workRecordDetailResponse
				if err := json.NewDecoder(res.Body).Decode(&row); err != nil {
					t.Fatalf("failed decoding work record details response: %v", err)
				}

				if row.WorkRecord != "W2604-0010" || row.Date != "2026-04-10" {
					t.Fatalf("unexpected detail response: %#v", row)
				}
				if row.JobID == "" || row.UID == "" {
					t.Fatalf("expected linked identifiers in detail response: %#v", row)
				}
				if row.CanCreateTimeEntry {
					t.Fatalf("report holder should not see create action in detail response: %#v", row)
				}
				if row.CanEdit || row.CanAddNote {
					t.Fatalf("report holder should not get edit/note actions in detail response: %#v", row)
				}
				if row.LinkedTimeEntryID != "" || row.LinkedTimeSheetID != "" {
					t.Fatalf("expected seeded work record to be unlinked: %#v", row)
				}
				if len(row.Notes) != 1 || row.Notes[0].Note == "" {
					t.Fatalf("expected seeded notes in detail response: %#v", row.Notes)
				}
				if len(row.Consumables) != 2 {
					t.Fatalf("expected seeded consumables in detail response: %#v", row.Consumables)
				}
				if row.Consumables[0].ConsumableName == "" || row.Consumables[0].InputKind == "" {
					t.Fatalf("expected enriched consumable details: %#v", row.Consumables)
				}
			},
		},
		{
			Name:   "work record holder can fetch work record details by id",
			Method: http.MethodGet,
			URL:    "/api/work_records/wr260400000003a",
			Headers: map[string]string{
				"Authorization": workRecordToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"work_record":"W2604-0003"`,
				`"can_create_time_entry":true`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, _ *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var row workRecordDetailResponse
				if err := json.NewDecoder(res.Body).Decode(&row); err != nil {
					t.Fatalf("failed decoding work record details response: %v", err)
				}
				if !row.CanEdit || !row.CanAddNote {
					t.Fatalf("expected owner to be able to edit and add notes: %#v", row)
				}
				if row.RootWorkRecordID != row.ID || row.RootWorkRecord != row.WorkRecord {
					t.Fatalf("expected root lineage to point at self for root record: %#v", row)
				}
			},
		},
		{
			Name:   "work record holder can add a note through the custom work record route",
			Method: http.MethodPost,
			URL:    "/api/work_records/wr260400000010a/notes",
			Headers: map[string]string{
				"Authorization": workRecordToken,
				"Content-Type":  "application/json",
			},
			Body: strings.NewReader(`{
				"note":"Created from the unified work record note route"
			}`),
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"note":"Created from the unified work record note route"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var payload workRecordNoteRow
				if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
					t.Fatalf("failed decoding work record note response: %v", err)
				}
				if payload.ID == "" {
					t.Fatalf("expected created note id, got %#v", payload)
				}
				if payload.UID != "f2j5a8vk006baub" {
					t.Fatalf("expected note uid to match authenticated worker, got %#v", payload)
				}

				saved := loadWorkRecordNoteForTest(t, app, payload.ID)
				if saved.Note != "Created from the unified work record note route" {
					t.Fatalf("unexpected saved note: %#v", saved)
				}
				if got := countWorkRecordNotesForTest(t, app, "wr260400000010a"); got != 2 {
					t.Fatalf("expected note count to increment, got %d", got)
				}
			},
		},
		{
			Name:   "work record holder cannot add an empty note through the custom work record route",
			Method: http.MethodPost,
			URL:    "/api/work_records/wr260400000010a/notes",
			Headers: map[string]string{
				"Authorization": workRecordToken,
				"Content-Type":  "application/json",
			},
			Body: strings.NewReader(`{
				"note":"   "
			}`),
			ExpectedStatus: http.StatusBadRequest,
			ExpectedContent: []string{
				`"note"`,
				`"validation_required"`,
				`"Note is required"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var payload struct {
					Message string                    `json:"message"`
					Data    map[string]map[string]any `json:"data"`
				}
				if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
					t.Fatalf("failed decoding empty-note rejection: %v", err)
				}

				fieldErr, ok := payload.Data["note"]
				if !ok {
					t.Fatalf("expected note validation error, got %#v", payload.Data)
				}
				if code, _ := fieldErr["code"].(string); code != "validation_required" {
					t.Fatalf("expected note required code, got %#v", fieldErr)
				}
				if message, _ := fieldErr["message"].(string); message != "Note is required" {
					t.Fatalf("expected note required message, got %#v", fieldErr)
				}
				if got := countWorkRecordNotesForTest(t, app, "wr260400000010a"); got != 1 {
					t.Fatalf("expected empty note rejection to avoid inserts, got %d notes", got)
				}
			},
		},
		{
			Name:   "report holder cannot add a note through the custom work record route",
			Method: http.MethodPost,
			URL:    "/api/work_records/wr260400000010a/notes",
			Headers: map[string]string{
				"Authorization": reportToken,
				"Content-Type":  "application/json",
			},
			Body: strings.NewReader(`{
				"note":"Unauthorized note attempt"
			}`),
			ExpectedStatus: http.StatusForbidden,
			ExpectedContent: []string{
				`"you are not authorized to add a note to this work record"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, _ *http.Response) {
				if got := countWorkRecordNotesForTest(t, app, "wr260400000010a"); got != 1 {
					t.Fatalf("expected unauthorized note attempt to avoid inserts, got %d notes", got)
				}
			},
		},
		{
			Name:   "regular user cannot fetch another users work record details",
			Method: http.MethodGet,
			URL:    "/api/work_records/W2604-0010",
			Headers: map[string]string{
				"Authorization": regularUserToken,
			},
			ExpectedStatus: http.StatusForbidden,
			ExpectedContent: []string{
				`"message":"you are not authorized to view this work record"`,
			},
			TestAppFactory: testutils.SetupTestApp,
		},
		{
			Name:   "authenticated non-report user can fetch approved work records for a job",
			Method: http.MethodGet,
			URL:    "/api/jobs/cjf0kt0defhq480/work-records?limit=10",
			Headers: map[string]string{
				"Authorization": regularUserToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"work_record":"W2604-0010"`,
				`"total":1`,
			},
			TestAppFactory: testutils.SetupTestApp,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, _ *core.ServeEvent) {
				if _, err := app.DB().NewQuery(`
					UPDATE work_records
					SET approved = 1
					WHERE id = {:id}
				`).Bind(map[string]any{
					"id": "wr260400000010a",
				}).Execute(); err != nil {
					t.Fatalf("failed to approve seeded work record: %v", err)
				}
			},
			AfterTestFunc: func(t testing.TB, _ *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var payload paginatedJobWorkRecordsResponse
				if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
					t.Fatalf("failed decoding job work records response: %v", err)
				}
				if payload.Total != 1 || payload.Page != 1 || payload.TotalPages != 1 {
					t.Fatalf("unexpected pagination metadata: %#v", payload)
				}
				if len(payload.Data) != 1 || payload.Data[0].WorkRecord != "W2604-0010" || !payload.Data[0].Approved {
					t.Fatalf("expected one approved job work record, got %#v", payload)
				}
			},
		},
		{
			Name:   "report holder can still fetch approved work records for a job",
			Method: http.MethodGet,
			URL:    "/api/jobs/cjf0kt0defhq480/work-records?limit=10",
			Headers: map[string]string{
				"Authorization": reportToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"work_record":"W2604-0010"`,
				`"total":1`,
			},
			TestAppFactory: testutils.SetupTestApp,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, _ *core.ServeEvent) {
				if _, err := app.DB().NewQuery(`
					UPDATE work_records
					SET approved = 1
					WHERE id = {:id}
				`).Bind(map[string]any{
					"id": "wr260400000010a",
				}).Execute(); err != nil {
					t.Fatalf("failed to approve seeded work record: %v", err)
				}
			},
			AfterTestFunc: func(t testing.TB, _ *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var payload paginatedJobWorkRecordsResponse
				if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
					t.Fatalf("failed decoding job work records response: %v", err)
				}
				if payload.Total != 1 || len(payload.Data) != 1 {
					t.Fatalf("expected one approved job work record, got %#v", payload)
				}
			},
		},
		{
			Name:   "job approved work records endpoint paginates",
			Method: http.MethodGet,
			URL:    "/api/jobs/cjf0kt0defhq480/work-records?page=2&limit=1",
			Headers: map[string]string{
				"Authorization": regularUserToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"page":2`,
				`"total_pages":2`,
			},
			TestAppFactory: testutils.SetupTestApp,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, _ *core.ServeEvent) {
				if _, err := app.DB().NewQuery(`
					UPDATE work_records
					SET approved = 1
					WHERE id IN ({:firstId}, {:secondId})
				`).Bind(map[string]any{
					"firstId":  "wr260400000010a",
					"secondId": "wr260400000003a",
				}).Execute(); err != nil {
					t.Fatalf("failed to approve seeded work records: %v", err)
				}
			},
			AfterTestFunc: func(t testing.TB, _ *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var payload paginatedJobWorkRecordsResponse
				if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
					t.Fatalf("failed decoding paginated job work records response: %v", err)
				}
				if payload.Page != 2 || payload.Limit != 1 || payload.Total != 2 || payload.TotalPages != 2 {
					t.Fatalf("unexpected pagination metadata: %#v", payload)
				}
				if len(payload.Data) != 1 || payload.Data[0].WorkRecord != "W2604-0003" {
					t.Fatalf("expected second page to contain the older approved work record, got %#v", payload)
				}
			},
		},
		{
			Name:   "report holder can fetch unlinked work records report",
			Method: http.MethodGet,
			URL:    "/api/reports/work_records/unlinked",
			Headers: map[string]string{
				"Authorization": reportToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"total":3`,
			},
			TestAppFactory: testutils.SetupTestApp,
		},
		{
			Name:   "report holder can fetch filtered work records report",
			Method: http.MethodGet,
			URL:    "/api/reports/work_records/search?job=cjf0kt0defhq480&type=wrtype000000001&date_from=2026-04-01&date_to=2026-04-30",
			Headers: map[string]string{
				"Authorization": reportToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"total":2`,
				`"work_record":"W2604-0010"`,
				`"work_record":"W2604-0003"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, _ *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var payload workRecordSearchResponse
				if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
					t.Fatalf("failed decoding filtered work record report response: %v", err)
				}
				if payload.Total != 2 || len(payload.Items) != 2 {
					t.Fatalf("expected 2 filtered work records in report response, got %#v", payload)
				}
				if payload.Items[0].WorkRecord != "W2604-0010" || payload.Items[1].WorkRecord != "W2604-0003" {
					t.Fatalf("unexpected filtered report row order: %#v", payload.Items)
				}
			},
		},
		{
			Name:   "report holder can fetch pending approval work records report",
			Method: http.MethodGet,
			URL:    "/api/reports/work_records/pending_approval",
			Headers: map[string]string{
				"Authorization": reportToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"total":3`,
			},
			TestAppFactory: testutils.SetupTestApp,
		},
		{
			Name:   "report holder can fetch work record consumables report",
			Method: http.MethodGet,
			URL:    "/api/reports/work_records/consumables",
			Headers: map[string]string{
				"Authorization": reportToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"consumable_name":"Survey Stakes"`,
				`"work_record":"W2604-0010"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, _ *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var payload workRecordConsumablesReportResponse
				if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
					t.Fatalf("failed decoding work record consumables report response: %v", err)
				}
				if len(payload.Items) != 2 {
					t.Fatalf("expected 2 consumable report rows, got %#v", payload)
				}
				if payload.Items[0].ConsumableName == "" || payload.Items[0].InputKind == "" {
					t.Fatalf("expected enriched consumable report rows, got %#v", payload)
				}
			},
		},
		{
			Name:   "consumables report still returns the full filtered set when extra paging params are sent",
			Method: http.MethodGet,
			URL:    "/api/reports/work_records/consumables?limit=1",
			Headers: map[string]string{
				"Authorization": reportToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"items":[`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, _ *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var payload workRecordConsumablesReportResponse
				if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
					t.Fatalf("failed decoding limited consumables report response: %v", err)
				}
				if len(payload.Items) != 2 {
					t.Fatalf("expected full consumables response when limit is supplied, got %#v", payload)
				}
			},
		},
		{
			Name:   "consumables report still returns the full filtered set when extra offset params are sent",
			Method: http.MethodGet,
			URL:    "/api/reports/work_records/consumables?limit=1&offset=1",
			Headers: map[string]string{
				"Authorization": reportToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"items":[`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, _ *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var payload workRecordConsumablesReportResponse
				if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
					t.Fatalf("failed decoding final-page consumables report response: %v", err)
				}
				if len(payload.Items) != 2 {
					t.Fatalf("expected full consumables response when offset is supplied, got %#v", payload)
				}
			},
		},
		{
			Name:   "work record holder can load a copy draft without creating a record",
			Method: http.MethodGet,
			URL:    "/api/work_records/wr260400000010a/copy_to_tomorrow",
			Headers: map[string]string{
				"Authorization": workRecordToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"date":"2026-04-11"`,
				`"source_work_record":"W2604-0010"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var row workRecordCopyDraftResponse
				if err := json.NewDecoder(res.Body).Decode(&row); err != nil {
					t.Fatalf("failed decoding work record copy draft response: %v", err)
				}

				if row.CopyFromWorkRecordID != "wr260400000010a" || row.RootWorkRecordID != "wr260400000010a" {
					t.Fatalf("expected root copy draft metadata, got %#v", row)
				}
				if row.RootWorkRecord != "W2604-0010" || row.Date != "2026-04-11" {
					t.Fatalf("unexpected copy draft response: %#v", row)
				}
				if len(row.Consumables) != 2 {
					t.Fatalf("expected copied draft consumables, got %#v", row.Consumables)
				}
				if got := countWorkRecordsForTest(t, app); got != 3 {
					t.Fatalf("expected copy draft load to avoid inserts, got %d work records", got)
				}
			},
		},
		{
			Name:   "saving a copied draft creates one linked work record with copied consumables",
			Method: http.MethodPost,
			URL:    "/api/work_records/save",
			Headers: map[string]string{
				"Authorization": workRecordToken,
				"Content-Type":  "application/json",
			},
			Body: strings.NewReader(`{
				"uid":"f2j5a8vk006baub",
				"date":"2026-04-11",
				"job":"cjf0kt0defhq480",
				"type":"wrtype000000001",
				"location":"Seeded location",
				"work_description":"Seeded work description",
				"report_to":"",
				"hours_on_site":8,
				"hours_travel_time":1,
				"distance_travelled_km":10,
				"field_book_number":"",
				"field_book_page_number":"",
				"sub_contractor":"",
				"equipment":"",
				"supplies":"",
				"vehicle_type":"company",
				"is_passenger":false,
				"company_vehicle_unit_number":"7",
				"copy_from_work_record_id":"wr260400000010a",
				"consumables":[
					{
						"entry_id":"",
						"consumable_id":"wrconsmbl000001",
						"input_kind":"number",
						"quantity_number":12,
						"boolean_selected":false
					},
					{
						"entry_id":"",
						"consumable_id":"wrconsmbl000002",
						"input_kind":"boolean",
						"quantity_number":0,
						"boolean_selected":true
					}
				]
			}`),
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"id":"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var row workRecordSaveResponse
				if err := json.NewDecoder(res.Body).Decode(&row); err != nil {
					t.Fatalf("failed decoding work record create response: %v", err)
				}
				if row.ID == "" {
					t.Fatalf("expected created work record id, got %#v", row)
				}

				var record struct {
					Number           string `db:"number"`
					Date             string `db:"date"`
					UID              string `db:"uid"`
					Job              string `db:"job"`
					ParentWorkRecord string `db:"parent_work_record"`
				}
				if err := app.DB().NewQuery(`
					SELECT
						COALESCE(number, '') AS number,
						COALESCE(date, '') AS date,
						COALESCE(uid, '') AS uid,
						COALESCE(job, '') AS job,
						COALESCE(parent_work_record, '') AS parent_work_record
					FROM work_records
					WHERE id = {:id}
				`).Bind(map[string]any{
					"id": row.ID,
				}).One(&record); err != nil {
					t.Fatalf("expected copied work record to be saved: %v", err)
				}

				if record.ParentWorkRecord != "wr260400000010a" || record.Date != "2026-04-11" {
					t.Fatalf("unexpected copied work record lineage/date: %#v", record)
				}
				if record.UID != "f2j5a8vk006baub" || record.Job != "cjf0kt0defhq480" {
					t.Fatalf("expected copied work record to keep source worker/job, got %#v", record)
				}
				if len(record.Number) < len("W2604-") || record.Number[:len("W2604-")] != "W2604-" {
					t.Fatalf("expected April-prefixed work record number, got %q", record.Number)
				}
				if got := countWorkRecordConsumableEntriesForTest(t, app, row.ID); got != 2 {
					t.Fatalf("expected copied consumable entries, got %d", got)
				}
			},
		},
		{
			Name:   "create rejects schema-required work record fields when payload is otherwise empty",
			Method: http.MethodPost,
			URL:    "/api/work_records/save",
			Headers: map[string]string{
				"Authorization": workRecordToken,
				"Content-Type":  "application/json",
			},
			Body: strings.NewReader(`{
				"uid":"f2j5a8vk006baub",
				"date":"",
				"job":"",
				"type":"",
				"location":"",
				"work_description":"",
				"report_to":"",
				"hours_on_site":0,
				"hours_travel_time":0,
				"distance_travelled_km":0,
				"field_book_number":"",
				"field_book_page_number":"",
				"sub_contractor":"",
				"equipment":"",
				"supplies":"",
				"vehicle_type":"",
				"is_passenger":false,
				"company_vehicle_unit_number":"",
				"consumables":[]
			}`),
			ExpectedStatus: http.StatusBadRequest,
			ExpectedContent: []string{
				`"code":"required"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var payload struct {
					Message string                    `json:"message"`
					Data    map[string]map[string]any `json:"data"`
				}
				if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
					t.Fatalf("failed decoding empty-create rejection: %v", err)
				}

				if strings.TrimSpace(payload.Message) == "" {
					t.Fatalf("expected validation failure message, got %#v", payload)
				}

				requiredFields := []string{"date", "job", "type", "location", "work_description", "vehicle_type"}
				for _, field := range requiredFields {
					fieldErr, ok := payload.Data[field]
					if !ok {
						t.Fatalf("expected %q validation error, got %#v", field, payload.Data)
					}
					if code, _ := fieldErr["code"].(string); code != "required" {
						t.Fatalf("expected %q required code, got %#v", field, fieldErr)
					}
					if message, _ := fieldErr["message"].(string); strings.TrimSpace(message) == "" {
						t.Fatalf("expected %q validation message, got %#v", field, fieldErr)
					}
				}

				if got := countWorkRecordsForTest(t, app); got != 3 {
					t.Fatalf("expected failed empty create to avoid inserts, got %d work records", got)
				}
			},
		},
		{
			Name:   "create rejects invalid field book number format through custom save route",
			Method: http.MethodPost,
			URL:    "/api/work_records/save",
			Headers: map[string]string{
				"Authorization": workRecordToken,
				"Content-Type":  "application/json",
			},
			Body: strings.NewReader(`{
				"uid":"f2j5a8vk006baub",
				"date":"2026-04-11",
				"job":"cjf0kt0defhq480",
				"type":"wrtype000000001",
				"location":"Invalid field book number create",
				"work_description":"Invalid field book number create",
				"report_to":"",
				"hours_on_site":8,
				"hours_travel_time":1,
				"distance_travelled_km":10,
				"field_book_number":"2026-001",
				"field_book_page_number":"",
				"sub_contractor":"",
				"equipment":"",
				"supplies":"",
				"vehicle_type":"company",
				"is_passenger":false,
				"company_vehicle_unit_number":"7",
				"consumables":[]
			}`),
			ExpectedStatus: http.StatusBadRequest,
			ExpectedContent: []string{
				`"field_book_number"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				payload := decodeWorkRecordAPIErrorResponse(t, res)
				if _, ok := payload.Data["field_book_number"]; !ok {
					t.Fatalf("expected field_book_number validation error, got %#v", payload)
				}
				if got := countWorkRecordsForTest(t, app); got != 3 {
					t.Fatalf("expected failed create to avoid inserts, got %d work records", got)
				}
			},
		},
		{
			Name:   "edit rejects invalid field book number format through custom save route",
			Method: http.MethodPost,
			URL:    "/api/work_records/wr260400000010a/save",
			Headers: map[string]string{
				"Authorization": workRecordToken,
				"Content-Type":  "application/json",
			},
			Body: strings.NewReader(`{
				"date":"2026-04-10",
				"job":"cjf0kt0defhq480",
				"type":"wrtype000000001",
				"location":"Seeded location",
				"work_description":"Seeded work description",
				"report_to":"",
				"hours_on_site":8,
				"hours_travel_time":1,
				"distance_travelled_km":10,
				"field_book_number":"2026-001",
				"field_book_page_number":"",
				"sub_contractor":"",
				"equipment":"",
				"supplies":"",
				"vehicle_type":"company",
				"is_passenger":false,
				"company_vehicle_unit_number":"7",
				"consumables":[
					{
						"entry_id":"wrcentry0000001",
						"consumable_id":"wrconsmbl000001",
						"input_kind":"number",
						"quantity_number":12,
						"boolean_selected":false
					},
					{
						"entry_id":"wrcentry0000002",
						"consumable_id":"wrconsmbl000002",
						"input_kind":"boolean",
						"quantity_number":0,
						"boolean_selected":true
					}
				]
			}`),
			ExpectedStatus: http.StatusBadRequest,
			ExpectedContent: []string{
				`"field_book_number"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				payload := decodeWorkRecordAPIErrorResponse(t, res)
				if _, ok := payload.Data["field_book_number"]; !ok {
					t.Fatalf("expected field_book_number validation error, got %#v", payload)
				}
				if got := countWorkRecordsForTest(t, app); got != 3 {
					t.Fatalf("expected failed edit to avoid inserts, got %d work records", got)
				}
			},
		},
		{
			Name:   "create rejects field book page number below schema minimum through custom save route",
			Method: http.MethodPost,
			URL:    "/api/work_records/save",
			Headers: map[string]string{
				"Authorization": workRecordToken,
				"Content-Type":  "application/json",
			},
			Body: strings.NewReader(`{
				"uid":"f2j5a8vk006baub",
				"date":"2026-04-11",
				"job":"cjf0kt0defhq480",
				"type":"wrtype000000001",
				"location":"Invalid page number create",
				"work_description":"Invalid page number create",
				"report_to":"",
				"hours_on_site":8,
				"hours_travel_time":1,
				"distance_travelled_km":10,
				"field_book_number":"",
				"field_book_page_number":0,
				"sub_contractor":"",
				"equipment":"",
				"supplies":"",
				"vehicle_type":"company",
				"is_passenger":false,
				"company_vehicle_unit_number":"7",
				"consumables":[]
			}`),
			ExpectedStatus: http.StatusBadRequest,
			ExpectedContent: []string{
				`"field_book_page_number"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				payload := decodeWorkRecordAPIErrorResponse(t, res)
				if _, ok := payload.Data["field_book_page_number"]; !ok {
					t.Fatalf("expected field_book_page_number validation error, got %#v", payload)
				}
				if got := countWorkRecordsForTest(t, app); got != 3 {
					t.Fatalf("expected failed create to avoid inserts, got %d work records", got)
				}
			},
		},
		{
			Name:   "edit rejects field book page number below schema minimum through custom save route",
			Method: http.MethodPost,
			URL:    "/api/work_records/wr260400000010a/save",
			Headers: map[string]string{
				"Authorization": workRecordToken,
				"Content-Type":  "application/json",
			},
			Body: strings.NewReader(`{
				"date":"2026-04-10",
				"job":"cjf0kt0defhq480",
				"type":"wrtype000000001",
				"location":"Seeded location",
				"work_description":"Seeded work description",
				"report_to":"",
				"hours_on_site":8,
				"hours_travel_time":1,
				"distance_travelled_km":10,
				"field_book_number":"",
				"field_book_page_number":0,
				"sub_contractor":"",
				"equipment":"",
				"supplies":"",
				"vehicle_type":"company",
				"is_passenger":false,
				"company_vehicle_unit_number":"7",
				"consumables":[
					{
						"entry_id":"wrcentry0000001",
						"consumable_id":"wrconsmbl000001",
						"input_kind":"number",
						"quantity_number":12,
						"boolean_selected":false
					},
					{
						"entry_id":"wrcentry0000002",
						"consumable_id":"wrconsmbl000002",
						"input_kind":"boolean",
						"quantity_number":0,
						"boolean_selected":true
					}
				]
			}`),
			ExpectedStatus: http.StatusBadRequest,
			ExpectedContent: []string{
				`"field_book_page_number"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				payload := decodeWorkRecordAPIErrorResponse(t, res)
				if _, ok := payload.Data["field_book_page_number"]; !ok {
					t.Fatalf("expected field_book_page_number validation error, got %#v", payload)
				}
				if got := countWorkRecordsForTest(t, app); got != 3 {
					t.Fatalf("expected failed edit to avoid inserts, got %d work records", got)
				}
			},
		},
		{
			Name:   "create rejects company vehicle unit values outside schema pattern through custom save route",
			Method: http.MethodPost,
			URL:    "/api/work_records/save",
			Headers: map[string]string{
				"Authorization": workRecordToken,
				"Content-Type":  "application/json",
			},
			Body: strings.NewReader(`{
				"uid":"f2j5a8vk006baub",
				"date":"2026-04-11",
				"job":"cjf0kt0defhq480",
				"type":"wrtype000000001",
				"location":"Invalid vehicle unit create",
				"work_description":"Invalid vehicle unit create",
				"report_to":"",
				"hours_on_site":8,
				"hours_travel_time":1,
				"distance_travelled_km":10,
				"field_book_number":"",
				"field_book_page_number":"",
				"sub_contractor":"",
				"equipment":"",
				"supplies":"",
				"vehicle_type":"company",
				"is_passenger":false,
				"company_vehicle_unit_number":"1000",
				"consumables":[]
			}`),
			ExpectedStatus: http.StatusBadRequest,
			ExpectedContent: []string{
				`"company_vehicle_unit_number"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				payload := decodeWorkRecordAPIErrorResponse(t, res)
				if _, ok := payload.Data["company_vehicle_unit_number"]; !ok {
					t.Fatalf("expected company_vehicle_unit_number validation error, got %#v", payload)
				}
				if got := countWorkRecordsForTest(t, app); got != 3 {
					t.Fatalf("expected failed create to avoid inserts, got %d work records", got)
				}
			},
		},
		{
			Name:   "edit rejects company vehicle unit values outside schema pattern through custom save route",
			Method: http.MethodPost,
			URL:    "/api/work_records/wr260400000010a/save",
			Headers: map[string]string{
				"Authorization": workRecordToken,
				"Content-Type":  "application/json",
			},
			Body: strings.NewReader(`{
				"date":"2026-04-10",
				"job":"cjf0kt0defhq480",
				"type":"wrtype000000001",
				"location":"Seeded location",
				"work_description":"Seeded work description",
				"report_to":"",
				"hours_on_site":8,
				"hours_travel_time":1,
				"distance_travelled_km":10,
				"field_book_number":"",
				"field_book_page_number":"",
				"sub_contractor":"",
				"equipment":"",
				"supplies":"",
				"vehicle_type":"company",
				"is_passenger":false,
				"company_vehicle_unit_number":"1000",
				"consumables":[
					{
						"entry_id":"wrcentry0000001",
						"consumable_id":"wrconsmbl000001",
						"input_kind":"number",
						"quantity_number":12,
						"boolean_selected":false
					},
					{
						"entry_id":"wrcentry0000002",
						"consumable_id":"wrconsmbl000002",
						"input_kind":"boolean",
						"quantity_number":0,
						"boolean_selected":true
					}
				]
			}`),
			ExpectedStatus: http.StatusBadRequest,
			ExpectedContent: []string{
				`"company_vehicle_unit_number"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				payload := decodeWorkRecordAPIErrorResponse(t, res)
				if _, ok := payload.Data["company_vehicle_unit_number"]; !ok {
					t.Fatalf("expected company_vehicle_unit_number validation error, got %#v", payload)
				}
				if got := countWorkRecordsForTest(t, app); got != 3 {
					t.Fatalf("expected failed edit to avoid inserts, got %d work records", got)
				}
			},
		},
		{
			Name:   "create rejects invalid job relation through custom save route",
			Method: http.MethodPost,
			URL:    "/api/work_records/save",
			Headers: map[string]string{
				"Authorization": workRecordToken,
				"Content-Type":  "application/json",
			},
			Body: strings.NewReader(`{
				"uid":"f2j5a8vk006baub",
				"date":"2026-04-11",
				"job":"notarealjobid",
				"type":"wrtype000000001",
				"location":"Invalid job create",
				"work_description":"Invalid job create",
				"report_to":"",
				"hours_on_site":8,
				"hours_travel_time":1,
				"distance_travelled_km":10,
				"field_book_number":"",
				"field_book_page_number":"",
				"sub_contractor":"",
				"equipment":"",
				"supplies":"",
				"vehicle_type":"company",
				"is_passenger":false,
				"company_vehicle_unit_number":"7",
				"consumables":[]
			}`),
			ExpectedStatus: http.StatusBadRequest,
			ExpectedContent: []string{
				`"job"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				payload := decodeWorkRecordAPIErrorResponse(t, res)
				if _, ok := payload.Data["job"]; !ok {
					t.Fatalf("expected job validation error, got %#v", payload)
				}
				if got := countWorkRecordsForTest(t, app); got != 3 {
					t.Fatalf("expected failed create to avoid inserts, got %d work records", got)
				}
			},
		},
		{
			Name:   "create rejects invalid type relation through custom save route",
			Method: http.MethodPost,
			URL:    "/api/work_records/save",
			Headers: map[string]string{
				"Authorization": workRecordToken,
				"Content-Type":  "application/json",
			},
			Body: strings.NewReader(`{
				"uid":"f2j5a8vk006baub",
				"date":"2026-04-11",
				"job":"cjf0kt0defhq480",
				"type":"notarealtypeid",
				"location":"Invalid type create",
				"work_description":"Invalid type create",
				"report_to":"",
				"hours_on_site":8,
				"hours_travel_time":1,
				"distance_travelled_km":10,
				"field_book_number":"",
				"field_book_page_number":"",
				"sub_contractor":"",
				"equipment":"",
				"supplies":"",
				"vehicle_type":"company",
				"is_passenger":false,
				"company_vehicle_unit_number":"7",
				"consumables":[]
			}`),
			ExpectedStatus: http.StatusBadRequest,
			ExpectedContent: []string{
				`"type"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				payload := decodeWorkRecordAPIErrorResponse(t, res)
				if _, ok := payload.Data["type"]; !ok {
					t.Fatalf("expected type validation error, got %#v", payload)
				}
				if got := countWorkRecordsForTest(t, app); got != 3 {
					t.Fatalf("expected failed create to avoid inserts, got %d work records", got)
				}
			},
		},
		{
			Name:   "create ignores forged boolean input kind when numeric quantity is provided",
			Method: http.MethodPost,
			URL:    "/api/work_records/save",
			Headers: map[string]string{
				"Authorization": workRecordToken,
				"Content-Type":  "application/json",
			},
			Body: strings.NewReader(`{
					"uid":"f2j5a8vk006baub",
					"date":"2026-04-11",
					"job":"cjf0kt0defhq480",
					"type":"wrtype000000001",
					"location":"Forged input kind create",
					"work_description":"Forged input kind create",
					"report_to":"",
					"hours_on_site":8,
					"hours_travel_time":1,
					"distance_travelled_km":10,
					"field_book_number":"",
					"field_book_page_number":"",
					"sub_contractor":"",
					"equipment":"",
					"supplies":"",
					"vehicle_type":"company",
					"is_passenger":false,
					"company_vehicle_unit_number":"7",
					"consumables":[
						{
							"entry_id":"",
							"consumable_id":"wrconsmbl000001",
							"input_kind":"boolean",
							"quantity_number":7,
							"boolean_selected":false
						}
					]
					}`),
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"id":"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var row workRecordSaveResponse
				if err := json.NewDecoder(res.Body).Decode(&row); err != nil {
					t.Fatalf("failed decoding forged-input create response: %v", err)
				}
				if got := loadWorkRecordConsumableQuantityForTest(t, app, row.ID, "wrconsmbl000001"); got != 7 {
					t.Fatalf("expected numeric consumable quantity to be persisted from server kind, got %v", got)
				}
			},
		},
		{
			Name:   "create ignores forged numeric input kind for boolean consumables",
			Method: http.MethodPost,
			URL:    "/api/work_records/save",
			Headers: map[string]string{
				"Authorization": workRecordToken,
				"Content-Type":  "application/json",
			},
			Body: strings.NewReader(`{
					"uid":"f2j5a8vk006baub",
					"date":"2026-04-11",
					"job":"cjf0kt0defhq480",
					"type":"wrtype000000001",
					"location":"Forged boolean quantity",
					"work_description":"Forged boolean quantity",
					"report_to":"",
					"hours_on_site":8,
					"hours_travel_time":1,
					"distance_travelled_km":10,
					"field_book_number":"",
					"field_book_page_number":"",
					"sub_contractor":"",
					"equipment":"",
					"supplies":"",
					"vehicle_type":"company",
					"is_passenger":false,
					"company_vehicle_unit_number":"7",
					"consumables":[
						{
							"entry_id":"",
							"consumable_id":"wrconsmbl000002",
							"input_kind":"number",
							"quantity_number":99,
							"boolean_selected":true
						}
					]
					}`),
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"id":"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var row workRecordSaveResponse
				if err := json.NewDecoder(res.Body).Decode(&row); err != nil {
					t.Fatalf("failed decoding forged-boolean create response: %v", err)
				}
				if got := loadWorkRecordConsumableQuantityForTest(t, app, row.ID, "wrconsmbl000002"); got != 0 {
					t.Fatalf("expected boolean consumable quantity to be forced to zero, got %v", got)
				}
			},
		},
		{
			Name:   "edit ignores forged boolean input kind when numeric consumable has quantity",
			Method: http.MethodPost,
			URL:    "/api/work_records/wr260400000010a/save",
			Headers: map[string]string{
				"Authorization": workRecordToken,
				"Content-Type":  "application/json",
			},
			Body: strings.NewReader(`{
					"date":"2026-04-10",
					"job":"cjf0kt0defhq480",
					"type":"wrtype000000001",
					"location":"Seeded location",
					"work_description":"Seeded work description",
					"report_to":"",
					"hours_on_site":8,
					"hours_travel_time":1,
					"distance_travelled_km":10,
					"field_book_number":"",
					"field_book_page_number":"",
					"sub_contractor":"",
					"equipment":"",
					"supplies":"",
					"vehicle_type":"company",
					"is_passenger":false,
					"company_vehicle_unit_number":"7",
					"consumables":[
						{
							"entry_id":"wrcentry0000001",
							"consumable_id":"wrconsmbl000001",
							"input_kind":"boolean",
							"quantity_number":9,
							"boolean_selected":false
						},
						{
							"entry_id":"wrcentry0000002",
							"consumable_id":"wrconsmbl000002",
							"input_kind":"boolean",
							"quantity_number":0,
							"boolean_selected":true
						}
					]
					}`),
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"id":"wr260400000010a"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				defer res.Body.Close()
				if got := loadWorkRecordConsumableQuantityForTest(t, app, "wr260400000010a", "wrconsmbl000001"); got != 9 {
					t.Fatalf("expected edit flow to respect server numeric kind, got %v", got)
				}
			},
		},
		{
			Name:   "copied draft duplicate-date conflicts are raised at save time",
			Method: http.MethodPost,
			URL:    "/api/work_records/save",
			Headers: map[string]string{
				"Authorization": workRecordToken,
				"Content-Type":  "application/json",
			},
			Body: strings.NewReader(`{
				"uid":"f2j5a8vk006baub",
				"date":"2026-04-10",
				"job":"cjf0kt0defhq480",
				"type":"wrtype000000001",
				"location":"Seeded location",
				"work_description":"Seeded work description",
				"report_to":"",
				"hours_on_site":8,
				"hours_travel_time":1,
				"distance_travelled_km":10,
				"field_book_number":"",
				"field_book_page_number":"",
				"sub_contractor":"",
				"equipment":"",
				"supplies":"",
				"vehicle_type":"company",
				"is_passenger":false,
				"company_vehicle_unit_number":"7",
				"copy_from_work_record_id":"wr260400000010a",
				"consumables":[
					{
						"entry_id":"",
						"consumable_id":"wrconsmbl000001",
						"input_kind":"number",
						"quantity_number":12,
						"boolean_selected":false
					},
					{
						"entry_id":"",
						"consumable_id":"wrconsmbl000002",
						"input_kind":"boolean",
						"quantity_number":0,
						"boolean_selected":true
					}
				]
			}`),
			ExpectedStatus: http.StatusBadRequest,
			ExpectedContent: []string{
				`"code":"duplicate_work_record_day"`,
			},
			TestAppFactory: testutils.SetupTestApp,
		},
		{
			Name:   "copying a child work record keeps the root parent in draft metadata",
			Method: http.MethodGet,
			URL:    "/api/work_records/wrchild0000001a/copy_to_tomorrow",
			Headers: map[string]string{
				"Authorization": workRecordToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"date":"2026-04-05"`,
				`"root_work_record_id":"wr260400000003a"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, _ *core.ServeEvent) {
				if _, err := app.DB().NewQuery(`
					INSERT INTO work_records (
						id, number, date, uid, creator, job, type, location,
						work_description, hours_on_site, hours_travel_time,
						distance_travelled_km, vehicle_type, approved, created,
						updated, parent_work_record
					) VALUES (
						{:id}, {:number}, {:date}, {:uid}, {:creator}, {:job}, {:type}, {:location},
						{:workDescription}, {:hoursOnSite}, {:hoursTravelTime},
						{:distanceTravelledKM}, {:vehicleType}, 0, {:created},
						{:updated}, {:parentWorkRecord}
					)
				`).Bind(map[string]any{
					"id":                  "wrchild0000001a",
					"number":              "W9901-0001",
					"date":                "2026-04-04",
					"uid":                 "f2j5a8vk006baub",
					"creator":             "f2j5a8vk006baub",
					"job":                 "cjf0kt0defhq480",
					"type":                "wrtype000000001",
					"location":            "Seeded child location",
					"workDescription":     "Seeded child work description",
					"hoursOnSite":         8,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 10,
					"vehicleType":         "company",
					"created":             "2026-04-04 12:00:00.000Z",
					"updated":             "2026-04-04 12:00:00.000Z",
					"parentWorkRecord":    "wr260400000003a",
				}).Execute(); err != nil {
					t.Fatalf("failed creating child work record: %v", err)
				}
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var row workRecordCopyDraftResponse
				if err := json.NewDecoder(res.Body).Decode(&row); err != nil {
					t.Fatalf("failed decoding work record copy draft response: %v", err)
				}

				if row.CopyFromWorkRecordID != "wrchild0000001a" {
					t.Fatalf("expected child draft to point at source record, got %#v", row)
				}
				if row.RootWorkRecordID != "wr260400000003a" || row.RootWorkRecord != "W2604-0003" {
					t.Fatalf("expected child draft to keep root metadata, got %#v", row)
				}
				if got := countWorkRecordsForTest(t, app); got != 4 {
					t.Fatalf("expected copy draft to avoid inserts after seeding child, got %d work records", got)
				}
			},
		},
		{
			Name:   "copied draft save uses the edited date for numbering",
			Method: http.MethodPost,
			URL:    "/api/work_records/save",
			Headers: map[string]string{
				"Authorization": workRecordToken,
				"Content-Type":  "application/json",
			},
			Body: strings.NewReader(`{
				"uid":"f2j5a8vk006baub",
				"date":"2026-03-01",
				"job":"cjf0kt0defhq480",
				"type":"wrtype000000001",
				"location":"January location",
				"work_description":"January work description",
				"report_to":"",
				"hours_on_site":8,
				"hours_travel_time":1,
				"distance_travelled_km":10,
				"field_book_number":"",
				"field_book_page_number":"",
				"sub_contractor":"",
				"equipment":"",
				"supplies":"",
				"vehicle_type":"company",
				"is_passenger":false,
				"company_vehicle_unit_number":"7",
				"copy_from_work_record_id":"wrjan310000001a",
				"consumables":[]
			}`),
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"id":"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, _ *core.ServeEvent) {
				if _, err := app.DB().NewQuery(`
					INSERT INTO work_records (
						id, number, date, uid, creator, job, type, location,
						work_description, hours_on_site, hours_travel_time,
						distance_travelled_km, vehicle_type, approved, created,
						updated
					) VALUES (
						{:id}, {:number}, {:date}, {:uid}, {:creator}, {:job}, {:type}, {:location},
						{:workDescription}, {:hoursOnSite}, {:hoursTravelTime},
						{:distanceTravelledKM}, {:vehicleType}, 0, {:created},
						{:updated}
					)
				`).Bind(map[string]any{
					"id":                  "wrjan310000001a",
					"number":              "W2601-0001",
					"date":                "2026-01-31",
					"uid":                 "f2j5a8vk006baub",
					"creator":             "f2j5a8vk006baub",
					"job":                 "cjf0kt0defhq480",
					"type":                "wrtype000000001",
					"location":            "January location",
					"workDescription":     "January work description",
					"hoursOnSite":         8,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 10,
					"vehicleType":         "company",
					"created":             "2026-01-31 12:00:00.000Z",
					"updated":             "2026-01-31 12:00:00.000Z",
				}).Execute(); err != nil {
					t.Fatalf("failed creating january work record: %v", err)
				}
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var row workRecordSaveResponse
				if err := json.NewDecoder(res.Body).Decode(&row); err != nil {
					t.Fatalf("failed decoding work record save response: %v", err)
				}

				var record struct {
					Number           string `db:"number"`
					Date             string `db:"date"`
					ParentWorkRecord string `db:"parent_work_record"`
				}
				if err := app.DB().NewQuery(`
					SELECT
						COALESCE(number, '') AS number,
						COALESCE(date, '') AS date,
						COALESCE(parent_work_record, '') AS parent_work_record
					FROM work_records
					WHERE id = {:id}
				`).Bind(map[string]any{
					"id": row.ID,
				}).One(&record); err != nil {
					t.Fatalf("failed loading copied march work record: %v", err)
				}
				if record.Date != "2026-03-01" || record.ParentWorkRecord != "wrjan310000001a" {
					t.Fatalf("unexpected copied march work record state: %#v", record)
				}
				if len(record.Number) < len("W2603-") || record.Number[:len("W2603-")] != "W2603-" {
					t.Fatalf("expected March-prefixed work record number, got %q", record.Number)
				}
			},
		},
		{
			Name:   "report holder cannot open a copy draft they cannot create",
			Method: http.MethodGet,
			URL:    "/api/work_records/wr260400000010a/copy_to_tomorrow",
			Headers: map[string]string{
				"Authorization": reportToken,
			},
			ExpectedStatus: http.StatusForbidden,
			ExpectedContent: []string{
				`"message":"you are not authorized to copy this work record"`,
			},
			TestAppFactory: testutils.SetupTestApp,
		},
		{
			Name:   "report holder cannot save a copied work record they cannot create",
			Method: http.MethodPost,
			URL:    "/api/work_records/save",
			Headers: map[string]string{
				"Authorization": reportToken,
				"Content-Type":  "application/json",
			},
			Body: strings.NewReader(`{
				"uid":"f2j5a8vk006baub",
				"date":"2026-04-11",
				"job":"cjf0kt0defhq480",
				"type":"wrtype000000001",
				"location":"Seeded location",
				"work_description":"Seeded work description",
				"report_to":"",
				"hours_on_site":8,
				"hours_travel_time":1,
				"distance_travelled_km":10,
				"field_book_number":"",
				"field_book_page_number":"",
				"sub_contractor":"",
				"equipment":"",
				"supplies":"",
				"vehicle_type":"company",
				"is_passenger":false,
				"company_vehicle_unit_number":"7",
				"copy_from_work_record_id":"wr260400000010a",
				"consumables":[]
			}`),
			ExpectedStatus: http.StatusForbidden,
			ExpectedContent: []string{
				`"message":"you are not authorized to copy this work record"`,
			},
			TestAppFactory: testutils.SetupTestApp,
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
