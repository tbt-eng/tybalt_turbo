package main

import (
	"encoding/json"
	"io"
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
	Attachment         string `json:"attachment"`
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

type workRecordSubjectDetailRow struct {
	ID                       string  `json:"id"`
	UID                      string  `json:"uid"`
	HoursOnSite              float64 `json:"hours_on_site"`
	HoursTravelTime          float64 `json:"hours_travel_time"`
	DistanceTravelledKM      float64 `json:"distance_travelled_km"`
	VehicleType              string  `json:"vehicle_type"`
	IsPassenger              bool    `json:"is_passenger"`
	CompanyVehicleUnitNumber string  `json:"company_vehicle_unit_number"`
	Approved                 bool    `json:"approved"`
}

type workRecordDetailResponse struct {
	ID                         string                       `json:"id"`
	WorkRecord                 string                       `json:"work_record"`
	Date                       string                       `json:"date"`
	UID                        string                       `json:"uid"`
	OwnerName                  string                       `json:"owner_name"`
	Creator                    string                       `json:"creator"`
	JobID                      string                       `json:"job_id"`
	JobNumber                  string                       `json:"job_number"`
	ClientName                 string                       `json:"client_name"`
	JobDescription             string                       `json:"job_description"`
	ClientReferenceNumber      string                       `json:"client_reference_number"`
	ClientPO                   string                       `json:"client_po"`
	TypeID                     string                       `json:"type_id"`
	TypeName                   string                       `json:"type_name"`
	Location                   string                       `json:"location"`
	WorkDescription            string                       `json:"work_description"`
	HoursOnSite                float64                      `json:"hours_on_site"`
	HoursTravelTime            float64                      `json:"hours_travel_time"`
	LaneKms                    float64                      `json:"lane_kms"`
	ReportTo                   string                       `json:"report_to"`
	DistanceTravelledKM        float64                      `json:"distance_travelled_km"`
	FieldBookNumber            string                       `json:"field_book_number"`
	FieldBookPageNumber        int                          `json:"field_book_page_number"`
	SubContractor              string                       `json:"sub_contractor"`
	Equipment                  string                       `json:"equipment"`
	Supplies                   string                       `json:"supplies"`
	VehicleType                string                       `json:"vehicle_type"`
	IsPassenger                bool                         `json:"is_passenger"`
	CompanyVehicleUnitNumber   string                       `json:"company_vehicle_unit_number"`
	Attachment                 string                       `json:"attachment"`
	Approved                   bool                         `json:"approved"`
	ParentWorkRecordID         string                       `json:"parent_work_record_id"`
	ParentWorkRecord           string                       `json:"parent_work_record"`
	RootWorkRecordID           string                       `json:"root_work_record_id"`
	RootWorkRecord             string                       `json:"root_work_record"`
	LinkedTimeEntryID          string                       `json:"linked_time_entry_id"`
	LinkedTimeSheetID          string                       `json:"linked_timesheet_id"`
	LinkedTimeEntryDate        string                       `json:"linked_time_entry_date"`
	LinkedTimeEntryDescription string                       `json:"linked_time_entry_description"`
	LinkedTimeSheetWeekEnding  string                       `json:"linked_time_sheet_week_ending"`
	CanCreateTimeEntry         bool                         `json:"can_create_time_entry"`
	CanCopyToTomorrow          bool                         `json:"can_copy_to_tomorrow"`
	CanEdit                    bool                         `json:"can_edit"`
	CanEditShared              bool                         `json:"can_edit_shared"`
	CanAddNote                 bool                         `json:"can_add_note"`
	Subjects                   []workRecordSubjectDetailRow `json:"subjects"`
	Notes                      []workRecordNoteRow          `json:"notes"`
	Consumables                []workRecordConsumableRow    `json:"consumables"`
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

func countWorkRecordSubjectsForTest(t testing.TB, app *tests.TestApp, workRecordID string) int {
	t.Helper()

	var row struct {
		Count int `db:"count"`
	}
	if err := app.DB().NewQuery(`
		SELECT COUNT(*) AS count
		FROM work_records_subjects
		WHERE work_record = {:id}
	`).Bind(map[string]any{
		"id": workRecordID,
	}).One(&row); err != nil {
		t.Fatalf("failed counting work record workers: %v", err)
	}

	return row.Count
}

func setWorkRecordSubjectsApprovedForTest(t testing.TB, app *tests.TestApp, workRecordID string, approved bool) {
	t.Helper()

	value := 0
	if approved {
		value = 1
	}
	if _, err := app.DB().NewQuery(`
		UPDATE work_records_subjects
		SET approved = {:approved}
		WHERE work_record = {:id}
	`).Bind(map[string]any{
		"approved": value,
		"id":       workRecordID,
	}).Execute(); err != nil {
		t.Fatalf("failed updating work record subject approval: %v", err)
	}
}

func linkTimeEntryToWorkRecordSubjectForTest(t testing.TB, app *tests.TestApp, timeEntryID string, subjectID string) {
	t.Helper()

	if _, err := app.DB().NewQuery(`
		UPDATE time_entries
		SET work_record_subject_id = {:subjectId},
		    work_record = ''
		WHERE id = {:timeEntryId}
	`).Bind(map[string]any{
		"subjectId":   subjectID,
		"timeEntryId": timeEntryID,
	}).Execute(); err != nil {
		t.Fatalf("failed linking time entry to work record worker: %v", err)
	}
}

func insertWorkRecordSubjectForTest(t testing.TB, app *tests.TestApp, params map[string]any) {
	t.Helper()

	isPassenger, _ := params["isPassenger"].(bool)
	approved, _ := params["approved"].(bool)
	if _, err := app.DB().NewQuery(`
		INSERT INTO work_records_subjects (
			id,
			work_record,
			uid,
			hours_on_site,
			hours_travel_time,
			distance_travelled_km,
			vehicle_type,
			is_passenger,
			company_vehicle_unit_number,
			approved
		) VALUES (
			{:subjectId},
			{:workRecordId},
			{:uid},
			{:hoursOnSite},
			{:hoursTravelTime},
			{:distanceTravelledKM},
			{:vehicleType},
			{:isPassenger},
			{:companyVehicleUnitNumber},
			{:approved}
		)
	`).Bind(map[string]any{
		"subjectId":                params["subjectId"],
		"workRecordId":             params["workRecordId"],
		"uid":                      params["uid"],
		"hoursOnSite":              params["hoursOnSite"],
		"hoursTravelTime":          params["hoursTravelTime"],
		"distanceTravelledKM":      params["distanceTravelledKM"],
		"vehicleType":              params["vehicleType"],
		"isPassenger":              isPassenger,
		"companyVehicleUnitNumber": testStringParam(params, "companyVehicleUnitNumber", "7"),
		"approved":                 approved,
	}).Execute(); err != nil {
		t.Fatalf("failed creating work record worker: %v", err)
	}
}

func insertWorkRecordWithSubjectForTest(t testing.TB, app *tests.TestApp, params map[string]any) {
	t.Helper()

	if _, err := app.DB().NewQuery(`
		INSERT INTO work_records (
			id,
			number,
			date,
			creator,
			job,
			type,
			location,
			work_description,
			created,
			updated,
			parent_work_record
		) VALUES (
			{:id},
			{:number},
			{:date},
			{:creator},
			{:job},
			{:type},
			{:location},
			{:workDescription},
			{:created},
			{:updated},
			{:parentWorkRecord}
		)
	`).Bind(map[string]any{
		"id":               params["id"],
		"number":           params["number"],
		"date":             params["date"],
		"creator":          params["creator"],
		"job":              params["job"],
		"type":             params["type"],
		"location":         params["location"],
		"workDescription":  params["workDescription"],
		"created":          firstNonEmptyTestString(params["created"], "2026-04-24 12:00:00.000Z"),
		"updated":          firstNonEmptyTestString(params["updated"], "2026-04-24 12:00:00.000Z"),
		"parentWorkRecord": firstNonEmptyTestString(params["parentWorkRecord"], ""),
	}).Execute(); err != nil {
		t.Fatalf("failed creating work record: %v", err)
	}

	insertWorkRecordSubjectForTest(t, app, map[string]any{
		"subjectId":                params["subjectId"],
		"workRecordId":             params["id"],
		"uid":                      params["uid"],
		"hoursOnSite":              params["hoursOnSite"],
		"hoursTravelTime":          params["hoursTravelTime"],
		"distanceTravelledKM":      params["distanceTravelledKM"],
		"vehicleType":              params["vehicleType"],
		"companyVehicleUnitNumber": firstNonEmptyTestString(params["companyVehicleUnitNumber"], "7"),
		"approved":                 params["approved"],
	})
}

func firstNonEmptyTestString(values ...any) string {
	for _, value := range values {
		if str, ok := value.(string); ok && strings.TrimSpace(str) != "" {
			return str
		}
	}
	return ""
}

func testStringParam(values map[string]any, key string, fallback string) string {
	value, ok := values[key]
	if !ok {
		return fallback
	}
	str, ok := value.(string)
	if !ok {
		return fallback
	}
	return str
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

func loadWorkRecordLaneKmsForTest(t testing.TB, app *tests.TestApp, workRecordID string) float64 {
	t.Helper()

	var row struct {
		LaneKms float64 `db:"lane_kms"`
	}
	if err := app.DB().NewQuery(`
		SELECT COALESCE(lane_kms, 0) AS lane_kms
		FROM work_records
		WHERE id = {:id}
	`).Bind(map[string]any{
		"id": workRecordID,
	}).One(&row); err != nil {
		t.Fatalf("failed loading work-record lane kms: %v", err)
	}

	return row.LaneKms
}

func loadWorkRecordSubjectRowsForTest(t testing.TB, app *tests.TestApp, workRecordID string) []workRecordSubjectDetailRow {
	t.Helper()

	rows := []workRecordSubjectDetailRow{}
	if err := app.DB().NewQuery(`
		SELECT
			id,
			COALESCE(uid, '') AS uid,
			COALESCE(hours_on_site, 0) AS hours_on_site,
			COALESCE(hours_travel_time, 0) AS hours_travel_time,
			COALESCE(distance_travelled_km, 0) AS distance_travelled_km,
			COALESCE(vehicle_type, '') AS vehicle_type,
			COALESCE(is_passenger, 0) AS is_passenger,
			COALESCE(company_vehicle_unit_number, '') AS company_vehicle_unit_number,
			COALESCE(approved, 0) AS approved
		FROM work_records_subjects
		WHERE work_record = {:id}
		ORDER BY id
	`).Bind(map[string]any{
		"id": workRecordID,
	}).All(&rows); err != nil {
		t.Fatalf("failed loading work-record workers: %v", err)
	}

	return rows
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

func workRecordRequestBody(raw string) io.Reader {
	var payload map[string]any
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		panic(err)
	}

	if _, hasSubjects := payload["subjects"]; !hasSubjects {
		if _, hasWorkerFields := payload["hours_on_site"]; hasWorkerFields {
			uid, _ := payload["uid"].(string)
			if uid == "" {
				uid = "f2j5a8vk006baub"
			}
			payload["subjects"] = []map[string]any{{
				"uid":                         uid,
				"hours_on_site":               payload["hours_on_site"],
				"hours_travel_time":           payload["hours_travel_time"],
				"distance_travelled_km":       payload["distance_travelled_km"],
				"vehicle_type":                payload["vehicle_type"],
				"is_passenger":                payload["is_passenger"],
				"company_vehicle_unit_number": payload["company_vehicle_unit_number"],
			}}
		}
	}

	for _, key := range []string{
		"uid",
		"hours_on_site",
		"hours_travel_time",
		"distance_travelled_km",
		"vehicle_type",
		"is_passenger",
		"company_vehicle_unit_number",
	} {
		delete(payload, key)
	}

	body, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	return strings.NewReader(string(body))
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
			Body: workRecordRequestBody(`{
				"date":"2026-04-10",
				"job":"cjf0kt0defhq480",
				"type":"rqayp9yoz1qe6ui",
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
				if typeID != "rqayp9yoz1qe6ui" || location != "Updated location" {
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
			Body: workRecordRequestBody(`{
				"date":"2026-04-10",
				"job":"cjf0kt0defhq480",
				"type":"rqayp9yoz1qe6ui",
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
				if typeID != "oask5u50nka0km8" {
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
			Body: workRecordRequestBody(`{
				"date":"2026-04-10",
				"job":"cjf0kt0defhq480",
				"type":"rqayp9yoz1qe6ui",
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
				if typeID != "oask5u50nka0km8" {
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
			Body: workRecordRequestBody(`{
				"date":"2026-04-10",
				"job":"cjf0kt0defhq480",
				"type":"rqayp9yoz1qe6ui",
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
				setWorkRecordSubjectsApprovedForTest(t, app, "wr260400000010a", true)
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, _ *http.Response) {
				typeID, _ := loadWorkRecordTypeAndLocationForTest(t, app, "wr260400000010a")
				if typeID != "oask5u50nka0km8" {
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
			Body: workRecordRequestBody(`{
				"date":"2026-04-10",
				"job":"cjf0kt0defhq480",
				"type":"oask5u50nka0km8",
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
				if typeID != "oask5u50nka0km8" || location != "Regression location" {
					t.Fatalf("expected non-type save to update location only, got type=%q location=%q", typeID, location)
				}
				if got := countWorkRecordConsumableEntriesForTest(t, app, "wr260400000010a"); got != 2 {
					t.Fatalf("expected unchanged consumables to remain present, got %d", got)
				}
			},
		},
		{
			Name:   "self create rejects lane kms when type disallows them",
			Method: http.MethodPost,
			URL:    "/api/work_records/save",
			Headers: map[string]string{
				"Authorization": regularUserToken,
				"Content-Type":  "application/json",
			},
			Body: workRecordRequestBody(`{
				"uid":"rzr98oadsp9qc11",
				"date":"2030-02-07",
				"job":"cjf0kt0defhq480",
				"type":"oask5u50nka0km8",
				"location":"Lane kms disallowed",
				"work_description":"Lane kms disallowed",
				"report_to":"",
				"lane_kms":12.5,
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
				`"code":"lane_kms_not_allowed"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, _ *http.Response) {
				if got := countWorkRecordsForTest(t, app); got != 5 {
					t.Fatalf("expected failed lane-kms create to avoid inserts, got %d work records", got)
				}
			},
		},
		{
			Name:   "self create saves lane kms when type allows them",
			Method: http.MethodPost,
			URL:    "/api/work_records/save",
			Headers: map[string]string{
				"Authorization": regularUserToken,
				"Content-Type":  "application/json",
			},
			Body: workRecordRequestBody(`{
				"uid":"rzr98oadsp9qc11",
				"date":"2030-02-07",
				"job":"cjf0kt0defhq480",
				"type":"lanekmstype0001",
				"location":"Lane kms allowed",
				"work_description":"Lane kms allowed",
				"report_to":"",
				"lane_kms":12.5,
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
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"id":"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var row workRecordSaveResponse
				if err := json.NewDecoder(res.Body).Decode(&row); err != nil {
					t.Fatalf("failed decoding lane-kms create response: %v", err)
				}
				if row.ID == "" {
					t.Fatalf("expected created work record id, got %#v", row)
				}
				if got := loadWorkRecordLaneKmsForTest(t, app, row.ID); got != 12.5 {
					t.Fatalf("expected lane kms to persist, got %v", got)
				}
				if got := countWorkRecordSubjectsForTest(t, app, row.ID); got != 1 {
					t.Fatalf("expected lane-kms create to save one worker, got %d", got)
				}
			},
		},
		{
			Name:   "creator can save multiple workers up to type max with separate fields",
			Method: http.MethodPost,
			URL:    "/api/work_records/wrmulti0000001a/save",
			Headers: map[string]string{
				"Authorization": workRecordToken,
				"Content-Type":  "application/json",
			},
			Body: workRecordRequestBody(`{
				"date":"2026-04-13",
				"job":"cjf0kt0defhq480",
				"type":"rqayp9yoz1qe6ui",
				"location":"Multi worker update",
				"work_description":"Multi worker update",
				"report_to":"",
				"lane_kms":0,
				"field_book_number":"",
				"field_book_page_number":"",
				"sub_contractor":"",
				"equipment":"",
				"supplies":"",
				"subjects":[
					{
						"id":"wrsubmulti0001a",
						"uid":"f2j5a8vk006baub",
						"hours_on_site":7.5,
						"hours_travel_time":0.5,
						"distance_travelled_km":11,
						"vehicle_type":"company",
						"is_passenger":false,
						"company_vehicle_unit_number":"7"
					},
					{
						"id":"",
						"uid":"rzr98oadsp9qc11",
						"hours_on_site":6,
						"hours_travel_time":1.5,
						"distance_travelled_km":22.2,
						"vehicle_type":"personal",
						"is_passenger":true,
						"company_vehicle_unit_number":""
					}
				],
				"consumables":[]
			}`),
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"id":"wrmulti0000001a"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, _ *core.ServeEvent) {
				insertWorkRecordWithSubjectForTest(t, app, map[string]any{
					"id":                  "wrmulti0000001a",
					"subjectId":           "wrsubmulti0001a",
					"number":              "W2604-0013",
					"date":                "2026-04-13",
					"uid":                 "f2j5a8vk006baub",
					"creator":             "f2j5a8vk006baub",
					"job":                 "cjf0kt0defhq480",
					"type":                "rqayp9yoz1qe6ui",
					"location":            "Multi worker seeded",
					"workDescription":     "Multi worker seeded",
					"hoursOnSite":         8,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 10,
					"vehicleType":         "company",
				})
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, _ *http.Response) {
				rows := loadWorkRecordSubjectRowsForTest(t, app, "wrmulti0000001a")
				if len(rows) != 2 {
					t.Fatalf("expected two workers after creator save, got %#v", rows)
				}
				var second *workRecordSubjectDetailRow
				for i := range rows {
					if rows[i].UID == "rzr98oadsp9qc11" {
						second = &rows[i]
						break
					}
				}
				if second == nil || second.HoursOnSite != 6 || second.HoursTravelTime != 1.5 || second.DistanceTravelledKM != 22.2 || !second.IsPassenger || second.VehicleType != "personal" {
					t.Fatalf("expected second worker fields to persist independently, got %#v", rows)
				}
			},
		},
		{
			Name:   "creator cannot save more workers than type max",
			Method: http.MethodPost,
			URL:    "/api/work_records/wrmaxroute0001a/save",
			Headers: map[string]string{
				"Authorization": workRecordToken,
				"Content-Type":  "application/json",
			},
			Body: workRecordRequestBody(`{
				"date":"2026-04-14",
				"job":"cjf0kt0defhq480",
				"type":"rqayp9yoz1qe6ui",
				"location":"Max worker update",
				"work_description":"Max worker update",
				"report_to":"",
				"lane_kms":0,
				"field_book_number":"",
				"field_book_page_number":"",
				"sub_contractor":"",
				"equipment":"",
				"supplies":"",
				"subjects":[
					{
						"id":"wrsubmaxrt0001a",
						"uid":"f2j5a8vk006baub",
						"hours_on_site":8,
						"hours_travel_time":1,
						"distance_travelled_km":10,
						"vehicle_type":"company",
						"is_passenger":false,
						"company_vehicle_unit_number":"7"
					},
					{
						"id":"",
						"uid":"rzr98oadsp9qc11",
						"hours_on_site":6,
						"hours_travel_time":1,
						"distance_travelled_km":12,
						"vehicle_type":"company",
						"is_passenger":false,
						"company_vehicle_unit_number":"8"
					},
					{
						"id":"",
						"uid":"u_mileage_valid",
						"hours_on_site":5,
						"hours_travel_time":1,
						"distance_travelled_km":14,
						"vehicle_type":"company",
						"is_passenger":false,
						"company_vehicle_unit_number":"9"
					}
				],
				"consumables":[]
			}`),
			ExpectedStatus: http.StatusBadRequest,
			ExpectedContent: []string{
				`"code":"max_subjects_exceeded"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, _ *core.ServeEvent) {
				insertWorkRecordWithSubjectForTest(t, app, map[string]any{
					"id":                  "wrmaxroute0001a",
					"subjectId":           "wrsubmaxrt0001a",
					"number":              "W2604-0014",
					"date":                "2026-04-14",
					"uid":                 "f2j5a8vk006baub",
					"creator":             "f2j5a8vk006baub",
					"job":                 "cjf0kt0defhq480",
					"type":                "rqayp9yoz1qe6ui",
					"location":            "Max worker seeded",
					"workDescription":     "Max worker seeded",
					"hoursOnSite":         8,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 10,
					"vehicleType":         "company",
				})
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, _ *http.Response) {
				if got := countWorkRecordSubjectsForTest(t, app, "wrmaxroute0001a"); got != 1 {
					t.Fatalf("expected rejected max-workers save to roll back worker changes, got %d workers", got)
				}
			},
		},
		{
			Name:   "creator can lower type max while removing extra workers",
			Method: http.MethodPost,
			URL:    "/api/work_records/wrmaxshrink001a/save",
			Headers: map[string]string{
				"Authorization": workRecordToken,
				"Content-Type":  "application/json",
			},
			Body: workRecordRequestBody(`{
					"date":"2026-04-15",
					"job":"cjf0kt0defhq480",
					"type":"oask5u50nka0km8",
					"location":"Max worker shrink update",
					"work_description":"Max worker shrink update",
					"report_to":"",
					"lane_kms":0,
					"field_book_number":"",
					"field_book_page_number":"",
					"sub_contractor":"",
					"equipment":"",
					"supplies":"",
					"subjects":[
						{
							"id":"wrsubshrink001a",
							"uid":"f2j5a8vk006baub",
							"hours_on_site":8,
							"hours_travel_time":1,
							"distance_travelled_km":10,
							"vehicle_type":"company",
							"is_passenger":false,
							"company_vehicle_unit_number":"7"
						}
					],
					"consumables":[]
				}`),
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"id":"wrmaxshrink001a"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, _ *core.ServeEvent) {
				insertWorkRecordWithSubjectForTest(t, app, map[string]any{
					"id":                  "wrmaxshrink001a",
					"subjectId":           "wrsubshrink001a",
					"number":              "W2604-0015",
					"date":                "2026-04-15",
					"uid":                 "f2j5a8vk006baub",
					"creator":             "f2j5a8vk006baub",
					"job":                 "cjf0kt0defhq480",
					"type":                "rqayp9yoz1qe6ui",
					"location":            "Max worker shrink seeded",
					"workDescription":     "Max worker shrink seeded",
					"hoursOnSite":         8,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 10,
					"vehicleType":         "company",
				})
				insertWorkRecordSubjectForTest(t, app, map[string]any{
					"subjectId":           "wrsubshrink002a",
					"workRecordId":        "wrmaxshrink001a",
					"uid":                 "rzr98oadsp9qc11",
					"hoursOnSite":         6,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 12,
					"vehicleType":         "company",
				})
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, _ *http.Response) {
				typeID, _ := loadWorkRecordTypeAndLocationForTest(t, app, "wrmaxshrink001a")
				if typeID != "oask5u50nka0km8" {
					t.Fatalf("expected type downgrade to persist, got %q", typeID)
				}
				rows := loadWorkRecordSubjectRowsForTest(t, app, "wrmaxshrink001a")
				if len(rows) != 1 || rows[0].UID != "f2j5a8vk006baub" {
					t.Fatalf("expected type downgrade to keep only submitted worker, got %#v", rows)
				}
			},
		},
		{
			Name:   "failed validation rolls back workers pruned before max check",
			Method: http.MethodPost,
			URL:    "/api/work_records/wrmaxroll00001a/save",
			Headers: map[string]string{
				"Authorization": workRecordToken,
				"Content-Type":  "application/json",
			},
			Body: workRecordRequestBody(`{
						"date":"2026-04-18",
						"job":"cjf0kt0defhq480",
						"type":"oask5u50nka0km8",
						"location":"Max worker rollback update",
						"work_description":"Max worker rollback update",
						"report_to":"",
						"lane_kms":0,
						"field_book_number":"",
						"field_book_page_number":"",
						"sub_contractor":"",
						"equipment":"",
						"supplies":"",
						"subjects":[
							{
								"id":"wrsubroll0001a",
								"uid":"f2j5a8vk006baub",
								"hours_on_site":8,
								"hours_travel_time":1,
								"distance_travelled_km":10,
								"vehicle_type":"company",
								"is_passenger":false,
								"company_vehicle_unit_number":""
							}
						],
						"consumables":[]
					}`),
			ExpectedStatus: http.StatusBadRequest,
			ExpectedContent: []string{
				`"code":"required_for_company_vehicle"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, _ *core.ServeEvent) {
				insertWorkRecordWithSubjectForTest(t, app, map[string]any{
					"id":                  "wrmaxroll00001a",
					"subjectId":           "wrsubroll0001a",
					"number":              "W2604-0018",
					"date":                "2026-04-18",
					"uid":                 "f2j5a8vk006baub",
					"creator":             "f2j5a8vk006baub",
					"job":                 "cjf0kt0defhq480",
					"type":                "rqayp9yoz1qe6ui",
					"location":            "Max worker rollback seeded",
					"workDescription":     "Max worker rollback seeded",
					"hoursOnSite":         8,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 10,
					"vehicleType":         "company",
				})
				insertWorkRecordSubjectForTest(t, app, map[string]any{
					"subjectId":           "wrsubroll0002a",
					"workRecordId":        "wrmaxroll00001a",
					"uid":                 "rzr98oadsp9qc11",
					"hoursOnSite":         6,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 12,
					"vehicleType":         "company",
				})
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, _ *http.Response) {
				typeID, location := loadWorkRecordTypeAndLocationForTest(t, app, "wrmaxroll00001a")
				if typeID != "rqayp9yoz1qe6ui" || location != "Max worker rollback seeded" {
					t.Fatalf("expected failed save to roll back parent fields, got type=%q location=%q", typeID, location)
				}
				rows := loadWorkRecordSubjectRowsForTest(t, app, "wrmaxroll00001a")
				if len(rows) != 2 {
					t.Fatalf("expected failed save to restore pruned worker roster, got %#v", rows)
				}
			},
		},
		{
			Name:   "creator can remove workers before approval",
			Method: http.MethodPost,
			URL:    "/api/work_records/wrremove000001a/save",
			Headers: map[string]string{
				"Authorization": workRecordToken,
				"Content-Type":  "application/json",
			},
			Body: workRecordRequestBody(`{
				"date":"2026-04-16",
				"job":"cjf0kt0defhq480",
				"type":"rqayp9yoz1qe6ui",
				"location":"Remove worker update",
				"work_description":"Remove worker update",
				"report_to":"",
				"lane_kms":0,
				"field_book_number":"",
				"field_book_page_number":"",
				"sub_contractor":"",
				"equipment":"",
				"supplies":"",
				"subjects":[
					{
						"id":"wrsubremove001a",
						"uid":"f2j5a8vk006baub",
						"hours_on_site":8,
						"hours_travel_time":1,
						"distance_travelled_km":10,
						"vehicle_type":"company",
						"is_passenger":false,
						"company_vehicle_unit_number":"7"
					}
				],
				"consumables":[]
			}`),
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"id":"wrremove000001a"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, _ *core.ServeEvent) {
				insertWorkRecordWithSubjectForTest(t, app, map[string]any{
					"id":                  "wrremove000001a",
					"subjectId":           "wrsubremove001a",
					"number":              "W2604-0016",
					"date":                "2026-04-16",
					"uid":                 "f2j5a8vk006baub",
					"creator":             "f2j5a8vk006baub",
					"job":                 "cjf0kt0defhq480",
					"type":                "rqayp9yoz1qe6ui",
					"location":            "Remove worker seeded",
					"workDescription":     "Remove worker seeded",
					"hoursOnSite":         8,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 10,
					"vehicleType":         "company",
				})
				insertWorkRecordSubjectForTest(t, app, map[string]any{
					"subjectId":           "wrsubremove002a",
					"workRecordId":        "wrremove000001a",
					"uid":                 "rzr98oadsp9qc11",
					"hoursOnSite":         6,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 12,
					"vehicleType":         "company",
				})
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, _ *http.Response) {
				rows := loadWorkRecordSubjectRowsForTest(t, app, "wrremove000001a")
				if len(rows) != 1 || rows[0].UID != "f2j5a8vk006baub" {
					t.Fatalf("expected creator save to remove omitted worker, got %#v", rows)
				}
			},
		},
		{
			Name:   "non-creator can edit own worker row without mutating shared record",
			Method: http.MethodPost,
			URL:    "/api/work_records/wrownedit00001a/save",
			Headers: map[string]string{
				"Authorization": regularUserToken,
				"Content-Type":  "application/json",
			},
			Body: workRecordRequestBody(`{
				"date":"2026-04-15",
				"job":"cjf0kt0defhq480",
				"type":"4k4nxilem8k10cu",
				"location":"Own worker edit",
				"work_description":"Own worker edit",
				"report_to":"",
				"lane_kms":0,
				"field_book_number":"",
				"field_book_page_number":"",
				"sub_contractor":"",
				"equipment":"",
				"supplies":"",
				"subjects":[
					{
						"id":"wrsubown000002a",
						"uid":"rzr98oadsp9qc11",
						"hours_on_site":5.25,
						"hours_travel_time":1.25,
						"distance_travelled_km":18,
						"vehicle_type":"personal",
						"is_passenger":true,
						"company_vehicle_unit_number":""
					}
				],
				"consumables":[
					{
						"entry_id":"wrcentryown0001",
						"consumable_id":"wrconsmbl000003",
						"input_kind":"number",
						"quantity_number":9,
						"boolean_selected":false
					}
				]
			}`),
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"id":"wrownedit00001a"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, _ *core.ServeEvent) {
				insertWorkRecordWithSubjectForTest(t, app, map[string]any{
					"id":                  "wrownedit00001a",
					"subjectId":           "wrsubown000001a",
					"number":              "W2604-0015",
					"date":                "2026-04-15",
					"uid":                 "f2j5a8vk006baub",
					"creator":             "f2j5a8vk006baub",
					"job":                 "cjf0kt0defhq480",
					"type":                "4k4nxilem8k10cu",
					"location":            "Own worker seeded",
					"workDescription":     "Own worker seeded",
					"hoursOnSite":         8,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 10,
					"vehicleType":         "company",
				})
				if _, err := app.DB().NewQuery(`
					INSERT INTO work_records_consumable_entries (
						id,
						work_record,
						consumable,
						quantity_number,
						created,
						updated
					) VALUES (
						'wrcentryown0001',
						'wrownedit00001a',
						'wrconsmbl000003',
						4,
						'2026-04-15 12:00:00.000Z',
						'2026-04-15 12:00:00.000Z'
					)
				`).Execute(); err != nil {
					t.Fatalf("failed creating non-creator consumable fixture: %v", err)
				}
				insertWorkRecordSubjectForTest(t, app, map[string]any{
					"subjectId":           "wrsubown000002a",
					"workRecordId":        "wrownedit00001a",
					"uid":                 "rzr98oadsp9qc11",
					"hoursOnSite":         6,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 12,
					"vehicleType":         "company",
				})
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, _ *http.Response) {
				_, location := loadWorkRecordTypeAndLocationForTest(t, app, "wrownedit00001a")
				if location != "Own worker seeded" {
					t.Fatalf("expected non-creator save to preserve shared work-record fields, got location %q", location)
				}
				if got := loadWorkRecordConsumableQuantityForTest(t, app, "wrownedit00001a", "wrconsmbl000003"); got != 4 {
					t.Fatalf("expected non-creator save to preserve consumables, got quantity %v", got)
				}

				rows := loadWorkRecordSubjectRowsForTest(t, app, "wrownedit00001a")
				if len(rows) != 2 {
					t.Fatalf("expected non-creator save to preserve other workers, got %#v", rows)
				}
				var own *workRecordSubjectDetailRow
				for i := range rows {
					if rows[i].UID == "rzr98oadsp9qc11" {
						own = &rows[i]
						break
					}
				}
				if own == nil || own.HoursOnSite != 5.25 || own.HoursTravelTime != 1.25 || own.DistanceTravelledKM != 18 || !own.IsPassenger || own.VehicleType != "personal" {
					t.Fatalf("expected non-creator to update only own worker fields, got %#v", rows)
				}
			},
		},
		{
			Name:   "non-creator can edit own unapproved worker row after another worker is approved",
			Method: http.MethodPost,
			URL:    "/api/work_records/wrpartialedit01/save",
			Headers: map[string]string{
				"Authorization": regularUserToken,
				"Content-Type":  "application/json",
			},
			Body: workRecordRequestBody(`{
				"date":"2026-04-19",
				"job":"cjf0kt0defhq480",
				"type":"rqayp9yoz1qe6ui",
				"location":"Partial locked seeded",
				"work_description":"Partial locked seeded",
				"report_to":"",
				"lane_kms":0,
				"field_book_number":"",
				"field_book_page_number":"",
				"sub_contractor":"",
				"equipment":"",
				"supplies":"",
				"subjects":[
					{
						"id":"wrsubpartialed2",
						"uid":"rzr98oadsp9qc11",
						"hours_on_site":6.5,
						"hours_travel_time":1.5,
						"distance_travelled_km":22,
						"vehicle_type":"personal",
						"is_passenger":true,
						"company_vehicle_unit_number":""
					}
				],
				"consumables":[]
			}`),
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"id":"wrpartialedit01"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, _ *core.ServeEvent) {
				insertWorkRecordWithSubjectForTest(t, app, map[string]any{
					"id":                  "wrpartialedit01",
					"subjectId":           "wrsubpartialed1",
					"number":              "W2604-0019",
					"date":                "2026-04-19",
					"uid":                 "f2j5a8vk006baub",
					"creator":             "f2j5a8vk006baub",
					"job":                 "cjf0kt0defhq480",
					"type":                "rqayp9yoz1qe6ui",
					"location":            "Partial locked seeded",
					"workDescription":     "Partial locked seeded",
					"hoursOnSite":         8,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 10,
					"vehicleType":         "company",
					"approved":            true,
				})
				insertWorkRecordSubjectForTest(t, app, map[string]any{
					"subjectId":           "wrsubpartialed2",
					"workRecordId":        "wrpartialedit01",
					"uid":                 "rzr98oadsp9qc11",
					"hoursOnSite":         6,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 12,
					"vehicleType":         "company",
				})
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, _ *http.Response) {
				_, location := loadWorkRecordTypeAndLocationForTest(t, app, "wrpartialedit01")
				if location != "Partial locked seeded" {
					t.Fatalf("expected locked shared fields to remain unchanged, got location %q", location)
				}
				rows := loadWorkRecordSubjectRowsForTest(t, app, "wrpartialedit01")
				for _, row := range rows {
					if row.ID == "wrsubpartialed1" && !row.Approved {
						t.Fatalf("expected approved sibling to remain approved, got %#v", rows)
					}
					if row.ID == "wrsubpartialed2" && (row.HoursOnSite != 6.5 || row.HoursTravelTime != 1.5 || !row.IsPassenger || row.VehicleType != "personal") {
						t.Fatalf("expected unapproved worker row to update, got %#v", rows)
					}
				}
			},
		},
		{
			Name:   "creator can edit an unapproved worker row after another worker is approved",
			Method: http.MethodPost,
			URL:    "/api/work_records/wrpartialedit02/save",
			Headers: map[string]string{
				"Authorization": workRecordToken,
				"Content-Type":  "application/json",
			},
			Body: workRecordRequestBody(`{
				"date":"2026-04-20",
				"job":"cjf0kt0defhq480",
				"type":"rqayp9yoz1qe6ui",
				"location":"Creator partial locked seeded",
				"work_description":"Creator partial locked seeded",
				"report_to":"",
				"lane_kms":0,
				"field_book_number":"",
				"field_book_page_number":"",
				"sub_contractor":"",
				"equipment":"",
				"supplies":"",
				"subjects":[
					{
						"id":"wrsubpartialcr2",
						"uid":"rzr98oadsp9qc11",
						"hours_on_site":7,
						"hours_travel_time":1.5,
						"distance_travelled_km":30,
						"vehicle_type":"personal",
						"is_passenger":false,
						"company_vehicle_unit_number":""
					}
				],
				"consumables":[]
			}`),
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"id":"wrpartialedit02"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, _ *core.ServeEvent) {
				insertWorkRecordWithSubjectForTest(t, app, map[string]any{
					"id":                  "wrpartialedit02",
					"subjectId":           "wrsubpartialcr1",
					"number":              "W2604-0020",
					"date":                "2026-04-20",
					"uid":                 "f2j5a8vk006baub",
					"creator":             "f2j5a8vk006baub",
					"job":                 "cjf0kt0defhq480",
					"type":                "rqayp9yoz1qe6ui",
					"location":            "Creator partial locked seeded",
					"workDescription":     "Creator partial locked seeded",
					"hoursOnSite":         8,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 10,
					"vehicleType":         "company",
					"approved":            true,
				})
				insertWorkRecordSubjectForTest(t, app, map[string]any{
					"subjectId":           "wrsubpartialcr2",
					"workRecordId":        "wrpartialedit02",
					"uid":                 "rzr98oadsp9qc11",
					"hoursOnSite":         6,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 12,
					"vehicleType":         "company",
				})
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, _ *http.Response) {
				rows := loadWorkRecordSubjectRowsForTest(t, app, "wrpartialedit02")
				for _, row := range rows {
					if row.ID == "wrsubpartialcr2" && (row.HoursOnSite != 7 || row.HoursTravelTime != 1.5 || row.DistanceTravelledKM != 30 || row.VehicleType != "personal") {
						t.Fatalf("expected creator to update unapproved worker row, got %#v", rows)
					}
				}
			},
		},
		{
			Name:   "locked parent rejects shared field changes while preserving worker rows",
			Method: http.MethodPost,
			URL:    "/api/work_records/wrpartialedit03/save",
			Headers: map[string]string{
				"Authorization": regularUserToken,
				"Content-Type":  "application/json",
			},
			Body: workRecordRequestBody(`{
				"date":"2026-04-21",
				"job":"cjf0kt0defhq480",
				"type":"rqayp9yoz1qe6ui",
				"location":"Changed locked location",
				"work_description":"Shared lock seeded",
				"report_to":"",
				"lane_kms":0,
				"field_book_number":"",
				"field_book_page_number":"",
				"sub_contractor":"",
				"equipment":"",
				"supplies":"",
				"subjects":[
					{
						"id":"wrsubpartialsh2",
						"uid":"rzr98oadsp9qc11",
						"hours_on_site":6.5,
						"hours_travel_time":1.5,
						"distance_travelled_km":22,
						"vehicle_type":"personal",
						"is_passenger":true,
						"company_vehicle_unit_number":""
					}
				],
				"consumables":[]
			}`),
			ExpectedStatus: http.StatusBadRequest,
			ExpectedContent: []string{
				`"code":"approved_work_record_locked"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, _ *core.ServeEvent) {
				insertWorkRecordWithSubjectForTest(t, app, map[string]any{
					"id":                  "wrpartialedit03",
					"subjectId":           "wrsubpartialsh1",
					"number":              "W2604-0021",
					"date":                "2026-04-21",
					"uid":                 "f2j5a8vk006baub",
					"creator":             "f2j5a8vk006baub",
					"job":                 "cjf0kt0defhq480",
					"type":                "rqayp9yoz1qe6ui",
					"location":            "Shared lock seeded",
					"workDescription":     "Shared lock seeded",
					"hoursOnSite":         8,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 10,
					"vehicleType":         "company",
					"approved":            true,
				})
				insertWorkRecordSubjectForTest(t, app, map[string]any{
					"subjectId":           "wrsubpartialsh2",
					"workRecordId":        "wrpartialedit03",
					"uid":                 "rzr98oadsp9qc11",
					"hoursOnSite":         6,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 12,
					"vehicleType":         "company",
				})
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, _ *http.Response) {
				_, location := loadWorkRecordTypeAndLocationForTest(t, app, "wrpartialedit03")
				if location != "Shared lock seeded" {
					t.Fatalf("expected locked shared field rejection to preserve location, got %q", location)
				}
				rows := loadWorkRecordSubjectRowsForTest(t, app, "wrpartialedit03")
				for _, row := range rows {
					if row.ID == "wrsubpartialsh2" && row.HoursOnSite != 6 {
						t.Fatalf("expected rejected save to preserve worker row, got %#v", rows)
					}
				}
			},
		},
		{
			Name:   "locked parent rejects consumable changes while preserving worker rows",
			Method: http.MethodPost,
			URL:    "/api/work_records/wrpartialedit04/save",
			Headers: map[string]string{
				"Authorization": regularUserToken,
				"Content-Type":  "application/json",
			},
			Body: workRecordRequestBody(`{
				"date":"2026-04-22",
				"job":"cjf0kt0defhq480",
				"type":"rqayp9yoz1qe6ui",
				"location":"Consumable lock seeded",
				"work_description":"Consumable lock seeded",
				"report_to":"",
				"lane_kms":0,
				"field_book_number":"",
				"field_book_page_number":"",
				"sub_contractor":"",
				"equipment":"",
				"supplies":"",
				"subjects":[
					{
						"id":"wrsubpartialcn2",
						"uid":"rzr98oadsp9qc11",
						"hours_on_site":6.5,
						"hours_travel_time":1.5,
						"distance_travelled_km":22,
						"vehicle_type":"personal",
						"is_passenger":true,
						"company_vehicle_unit_number":""
					}
				],
				"consumables":[]
			}`),
			ExpectedStatus: http.StatusBadRequest,
			ExpectedContent: []string{
				`"code":"approved_work_record_locked"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, _ *core.ServeEvent) {
				insertWorkRecordWithSubjectForTest(t, app, map[string]any{
					"id":                  "wrpartialedit04",
					"subjectId":           "wrsubpartialcn1",
					"number":              "W2604-0022",
					"date":                "2026-04-22",
					"uid":                 "f2j5a8vk006baub",
					"creator":             "f2j5a8vk006baub",
					"job":                 "cjf0kt0defhq480",
					"type":                "rqayp9yoz1qe6ui",
					"location":            "Consumable lock seeded",
					"workDescription":     "Consumable lock seeded",
					"hoursOnSite":         8,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 10,
					"vehicleType":         "company",
					"approved":            true,
				})
				insertWorkRecordSubjectForTest(t, app, map[string]any{
					"subjectId":           "wrsubpartialcn2",
					"workRecordId":        "wrpartialedit04",
					"uid":                 "rzr98oadsp9qc11",
					"hoursOnSite":         6,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 12,
					"vehicleType":         "company",
				})
				if _, err := app.DB().NewQuery(`
					INSERT INTO work_records_consumable_entries (
						id,
						work_record,
						consumable,
						quantity_number,
						created,
						updated
					) VALUES (
						'wrcentrylock001',
						'wrpartialedit04',
						'wrconsmbl000003',
						4,
						'2026-04-22 12:00:00.000Z',
						'2026-04-22 12:00:00.000Z'
					)
				`).Execute(); err != nil {
					t.Fatalf("failed creating locked consumable fixture: %v", err)
				}
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, _ *http.Response) {
				if got := loadWorkRecordConsumableQuantityForTest(t, app, "wrpartialedit04", "wrconsmbl000003"); got != 4 {
					t.Fatalf("expected locked consumable rejection to preserve quantity, got %v", got)
				}
				rows := loadWorkRecordSubjectRowsForTest(t, app, "wrpartialedit04")
				for _, row := range rows {
					if row.ID == "wrsubpartialcn2" && row.HoursOnSite != 6 {
						t.Fatalf("expected rejected save to preserve worker row, got %#v", rows)
					}
				}
			},
		},
		{
			Name:   "non-creator cannot add workers",
			Method: http.MethodPost,
			URL:    "/api/work_records/wrnoadd0000001a/save",
			Headers: map[string]string{
				"Authorization": regularUserToken,
				"Content-Type":  "application/json",
			},
			Body: workRecordRequestBody(`{
				"date":"2026-04-17",
				"job":"cjf0kt0defhq480",
				"type":"rqayp9yoz1qe6ui",
				"location":"Non creator add",
				"work_description":"Non creator add",
				"report_to":"",
				"lane_kms":0,
				"field_book_number":"",
				"field_book_page_number":"",
				"sub_contractor":"",
				"equipment":"",
				"supplies":"",
				"subjects":[
					{
						"id":"wrsubnoadd0002a",
						"uid":"rzr98oadsp9qc11",
						"hours_on_site":6,
						"hours_travel_time":1,
						"distance_travelled_km":12,
						"vehicle_type":"company",
						"is_passenger":false,
						"company_vehicle_unit_number":"8"
					},
					{
						"id":"",
						"uid":"dkv192wxprcqmho",
						"hours_on_site":5,
						"hours_travel_time":1,
						"distance_travelled_km":14,
						"vehicle_type":"company",
						"is_passenger":false,
						"company_vehicle_unit_number":"9"
					}
				],
				"consumables":[]
			}`),
			ExpectedStatus: http.StatusBadRequest,
			ExpectedContent: []string{
				`"code":"creator_required"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, _ *core.ServeEvent) {
				insertWorkRecordWithSubjectForTest(t, app, map[string]any{
					"id":                  "wrnoadd0000001a",
					"subjectId":           "wrsubnoadd0001a",
					"number":              "W2604-0017",
					"date":                "2026-04-17",
					"uid":                 "f2j5a8vk006baub",
					"creator":             "f2j5a8vk006baub",
					"job":                 "cjf0kt0defhq480",
					"type":                "rqayp9yoz1qe6ui",
					"location":            "Non creator add seeded",
					"workDescription":     "Non creator add seeded",
					"hoursOnSite":         8,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 10,
					"vehicleType":         "company",
				})
				insertWorkRecordSubjectForTest(t, app, map[string]any{
					"subjectId":                "wrsubnoadd0002a",
					"workRecordId":             "wrnoadd0000001a",
					"uid":                      "rzr98oadsp9qc11",
					"hoursOnSite":              6,
					"hoursTravelTime":          1,
					"distanceTravelledKM":      12,
					"vehicleType":              "company",
					"companyVehicleUnitNumber": "8",
				})
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, _ *http.Response) {
				if got := countWorkRecordSubjectsForTest(t, app, "wrnoadd0000001a"); got != 2 {
					t.Fatalf("expected non-creator add rejection to preserve worker roster, got %d workers", got)
				}
			},
		},
		{
			Name:   "time-only creator can add another time-eligible worker",
			Method: http.MethodPost,
			URL:    "/api/work_records/wrcreatornoadd1/save",
			Headers: map[string]string{
				"Authorization": regularUserToken,
				"Content-Type":  "application/json",
			},
			Body: workRecordRequestBody(`{
				"date":"2026-04-18",
				"job":"cjf0kt0defhq480",
				"type":"rqayp9yoz1qe6ui",
				"location":"Creator no claim add",
				"work_description":"Creator no claim add",
				"report_to":"",
				"lane_kms":0,
				"field_book_number":"",
				"field_book_page_number":"",
				"sub_contractor":"",
				"equipment":"",
				"supplies":"",
				"subjects":[
					{
						"id":"wrsubcrnoadd01",
						"uid":"rzr98oadsp9qc11",
						"hours_on_site":8,
						"hours_travel_time":1,
						"distance_travelled_km":12,
						"vehicle_type":"company",
						"is_passenger":false,
						"company_vehicle_unit_number":"8"
					},
					{
						"id":"",
						"uid":"u_mileage_valid",
						"hours_on_site":5,
						"hours_travel_time":1,
						"distance_travelled_km":14,
						"vehicle_type":"company",
						"is_passenger":false,
						"company_vehicle_unit_number":"9"
					}
				],
				"consumables":[]
			}`),
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"id":"wrcreatornoadd1"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, _ *core.ServeEvent) {
				insertWorkRecordWithSubjectForTest(t, app, map[string]any{
					"id":                  "wrcreatornoadd1",
					"subjectId":           "wrsubcrnoadd01",
					"number":              "W2604-0018",
					"date":                "2026-04-18",
					"uid":                 "rzr98oadsp9qc11",
					"creator":             "rzr98oadsp9qc11",
					"job":                 "cjf0kt0defhq480",
					"type":                "rqayp9yoz1qe6ui",
					"location":            "Creator no claim seeded",
					"workDescription":     "Creator no claim seeded",
					"hoursOnSite":         8,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 12,
					"vehicleType":         "company",
				})
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, _ *http.Response) {
				if got := countWorkRecordSubjectsForTest(t, app, "wrcreatornoadd1"); got != 2 {
					t.Fatalf("expected creator add to save worker roster, got %d workers", got)
				}
				rows := loadWorkRecordSubjectRowsForTest(t, app, "wrcreatornoadd1")
				found := false
				for _, row := range rows {
					if row.UID == "u_mileage_valid" {
						found = true
						break
					}
				}
				if !found {
					t.Fatalf("expected added worker to be saved, got %#v", rows)
				}
			},
		},
		{
			Name:   "work record creator cannot add worker without time claim",
			Method: http.MethodPost,
			URL:    "/api/work_records/wrtargetnotime1/save",
			Headers: map[string]string{
				"Authorization": workRecordToken,
				"Content-Type":  "application/json",
			},
			Body: workRecordRequestBody(`{
				"date":"2026-04-19",
				"job":"cjf0kt0defhq480",
				"type":"rqayp9yoz1qe6ui",
				"location":"Target no time add",
				"work_description":"Target no time add",
				"report_to":"",
				"lane_kms":0,
				"field_book_number":"",
				"field_book_page_number":"",
				"sub_contractor":"",
				"equipment":"",
				"supplies":"",
				"subjects":[
					{
						"id":"wrsubtargetnt1",
						"uid":"f2j5a8vk006baub",
						"hours_on_site":8,
						"hours_travel_time":1,
						"distance_travelled_km":12,
						"vehicle_type":"company",
						"is_passenger":false,
						"company_vehicle_unit_number":"8"
					},
					{
						"id":"",
						"uid":"u_no_claims",
						"hours_on_site":5,
						"hours_travel_time":1,
						"distance_travelled_km":14,
						"vehicle_type":"company",
						"is_passenger":false,
						"company_vehicle_unit_number":"9"
					}
				],
				"consumables":[]
			}`),
			ExpectedStatus: http.StatusForbidden,
			ExpectedContent: []string{
				`"code":"unauthorized"`,
				`"you are not authorized to add this worker"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, _ *core.ServeEvent) {
				insertWorkRecordWithSubjectForTest(t, app, map[string]any{
					"id":                  "wrtargetnotime1",
					"subjectId":           "wrsubtargetnt1",
					"number":              "W2604-0019",
					"date":                "2026-04-19",
					"uid":                 "f2j5a8vk006baub",
					"creator":             "f2j5a8vk006baub",
					"job":                 "cjf0kt0defhq480",
					"type":                "rqayp9yoz1qe6ui",
					"location":            "Target no time seeded",
					"workDescription":     "Target no time seeded",
					"hoursOnSite":         8,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 12,
					"vehicleType":         "company",
				})
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, _ *http.Response) {
				if got := countWorkRecordSubjectsForTest(t, app, "wrtargetnotime1"); got != 1 {
					t.Fatalf("expected ineligible target add to roll back worker roster, got %d workers", got)
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
			Name:   "time holder can fetch own work records list",
			Method: http.MethodGet,
			URL:    "/api/work_records",
			Headers: map[string]string{
				"Authorization": regularUserToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"total":2`,
				`"work_record":"W3002-1051"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, _ *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var payload workRecordSearchResponse
				if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
					t.Fatalf("failed decoding work records list response: %v", err)
				}

				if payload.Total != 2 || len(payload.Items) != 2 {
					t.Fatalf("expected 2 visible self-owned work records, got %#v", payload)
				}
				if payload.Items[0].WorkRecord != "W3002-1051" || payload.Items[1].WorkRecord != "W3002-1050" {
					t.Fatalf("unexpected self-owned row order: %#v", payload.Items)
				}
				if !payload.Items[0].CanCreateTimeEntry || payload.Items[0].LinkedTimeEntryID != "" {
					t.Fatalf("expected top self-owned row to allow time-entry creation: %#v", payload.Items[0])
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
				insertWorkRecordWithSubjectForTest(t, app, map[string]any{
					"id":                  "wrcreator000001a",
					"subjectId":           "wrsubcreator01",
					"number":              "W2604-0012",
					"date":                "2026-04-12",
					"uid":                 "rzr98oadsp9qc11",
					"creator":             "f2j5a8vk006baub",
					"job":                 "cjf0kt0defhq480",
					"type":                "oask5u50nka0km8",
					"location":            "Creator visibility test",
					"workDescription":     "Created for another worker",
					"hoursOnSite":         8,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 5,
					"vehicleType":         "personal",
				})
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
				`"total":5`,
				`"work_record":"W3002-1051"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, _ *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var payload workRecordSearchResponse
				if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
					t.Fatalf("failed decoding work record list response: %v", err)
				}
				if payload.Total != 5 || len(payload.Items) != 5 {
					t.Fatalf("expected 5 organization work records without endpoint-level filtering, got %#v", payload)
				}
				if payload.Items[0].CanCreateTimeEntry {
					t.Fatalf("report all-scope should not expose create action for another user's row: %#v", payload.Items[0])
				}
			},
		},
		{
			Name:   "report holder can fetch work record details by id",
			Method: http.MethodGet,
			URL:    "/api/work_records/wr260400000010a",
			Headers: map[string]string{
				"Authorization": reportToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"work_record":"W2604-0010"`,
				`"job_number":"24-321"`,
				`"client_name":"Mobilia SA"`,
				`"job_description":"Clarke St Boiler Replacement"`,
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
				if row.ClientName != "Mobilia SA" || row.JobDescription != "Clarke St Boiler Replacement" {
					t.Fatalf("expected job-derived client context in detail response: %#v", row)
				}
				if row.ClientPO != "" || row.ClientReferenceNumber != "" {
					t.Fatalf("expected blank optional client fields to serialize as empty strings: %#v", row)
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
			Name:   "details response includes all workers with separate fields",
			Method: http.MethodGet,
			URL:    "/api/work_records/wrdetail000001a",
			Headers: map[string]string{
				"Authorization": reportToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"subjects":[`,
				`"uid":"rzr98oadsp9qc11"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, _ *core.ServeEvent) {
				insertWorkRecordWithSubjectForTest(t, app, map[string]any{
					"id":                  "wrdetail000001a",
					"subjectId":           "wrsubdetail001a",
					"number":              "W2604-0018",
					"date":                "2026-04-18",
					"uid":                 "f2j5a8vk006baub",
					"creator":             "f2j5a8vk006baub",
					"job":                 "cjf0kt0defhq480",
					"type":                "rqayp9yoz1qe6ui",
					"location":            "Detail workers seeded",
					"workDescription":     "Detail workers seeded",
					"hoursOnSite":         8,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 10,
					"vehicleType":         "company",
				})
				insertWorkRecordSubjectForTest(t, app, map[string]any{
					"subjectId":                "wrsubdetail002a",
					"workRecordId":             "wrdetail000001a",
					"uid":                      "rzr98oadsp9qc11",
					"hoursOnSite":              6.5,
					"hoursTravelTime":          1.25,
					"distanceTravelledKM":      18.75,
					"vehicleType":              "personal",
					"isPassenger":              true,
					"companyVehicleUnitNumber": "",
				})
			},
			AfterTestFunc: func(t testing.TB, _ *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var row workRecordDetailResponse
				if err := json.NewDecoder(res.Body).Decode(&row); err != nil {
					t.Fatalf("failed decoding multi-worker detail response: %v", err)
				}
				if len(row.Subjects) != 2 {
					t.Fatalf("expected detail response to include both workers, got %#v", row.Subjects)
				}
				var second *workRecordSubjectDetailRow
				for i := range row.Subjects {
					if row.Subjects[i].UID == "rzr98oadsp9qc11" {
						second = &row.Subjects[i]
						break
					}
				}
				if second == nil || second.HoursOnSite != 6.5 || second.HoursTravelTime != 1.25 || second.DistanceTravelledKM != 18.75 || !second.IsPassenger || second.VehicleType != "personal" {
					t.Fatalf("expected detail response to preserve separate worker fields, got %#v", row.Subjects)
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
			Body: workRecordRequestBody(`{
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
			Body: workRecordRequestBody(`{
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
			Body: workRecordRequestBody(`{
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
			Name:   "regular user cannot fetch another users work record details by id",
			Method: http.MethodGet,
			URL:    "/api/work_records/wr260400000010a",
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
				`"total":2`,
			},
			TestAppFactory: testutils.SetupTestApp,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, _ *core.ServeEvent) {
				setWorkRecordSubjectsApprovedForTest(t, app, "wr260400000010a", true)
			},
			AfterTestFunc: func(t testing.TB, _ *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var payload paginatedJobWorkRecordsResponse
				if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
					t.Fatalf("failed decoding job work records response: %v", err)
				}
				if payload.Total != 2 || payload.Page != 1 || payload.TotalPages != 1 {
					t.Fatalf("unexpected pagination metadata: %#v", payload)
				}
				if len(payload.Data) != 2 {
					t.Fatalf("expected two approved job work records, got %#v", payload)
				}
				if payload.Data[0].WorkRecord != "W3002-0050" || !payload.Data[0].Approved {
					t.Fatalf("expected seeded approved work record first, got %#v", payload)
				}
				if payload.Data[1].WorkRecord != "W2604-0010" || !payload.Data[1].Approved {
					t.Fatalf("expected requested approved work record second, got %#v", payload)
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
				`"total":2`,
			},
			TestAppFactory: testutils.SetupTestApp,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, _ *core.ServeEvent) {
				setWorkRecordSubjectsApprovedForTest(t, app, "wr260400000010a", true)
			},
			AfterTestFunc: func(t testing.TB, _ *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var payload paginatedJobWorkRecordsResponse
				if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
					t.Fatalf("failed decoding job work records response: %v", err)
				}
				if payload.Total != 2 || len(payload.Data) != 2 {
					t.Fatalf("expected two approved job work records, got %#v", payload)
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
				`"total_pages":3`,
			},
			TestAppFactory: testutils.SetupTestApp,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, _ *core.ServeEvent) {
				setWorkRecordSubjectsApprovedForTest(t, app, "wr260400000010a", true)
				setWorkRecordSubjectsApprovedForTest(t, app, "wr260400000003a", true)
			},
			AfterTestFunc: func(t testing.TB, _ *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var payload paginatedJobWorkRecordsResponse
				if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
					t.Fatalf("failed decoding paginated job work records response: %v", err)
				}
				if payload.Page != 2 || payload.Limit != 1 || payload.Total != 3 || payload.TotalPages != 3 {
					t.Fatalf("unexpected pagination metadata: %#v", payload)
				}
				if len(payload.Data) != 1 || payload.Data[0].WorkRecord != "W2604-0010" {
					t.Fatalf("expected second page to contain the middle approved work record, got %#v", payload)
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
				`"total":5`,
			},
			TestAppFactory: testutils.SetupTestApp,
		},
		{
			Name:   "unlinked report includes partially linked multi-worker records",
			Method: http.MethodGet,
			URL:    "/api/reports/work_records/unlinked",
			Headers: map[string]string{
				"Authorization": reportToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"work_record":"W3101-0001"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, _ *core.ServeEvent) {
				insertWorkRecordWithSubjectForTest(t, app, map[string]any{
					"id":                  "wrpartiallink1a",
					"subjectId":           "wrsubpartialln1",
					"number":              "W3101-0001",
					"date":                "2031-01-02",
					"uid":                 "f2j5a8vk006baub",
					"creator":             "f2j5a8vk006baub",
					"job":                 "cjf0kt0defhq480",
					"type":                "rqayp9yoz1qe6ui",
					"location":            "Partial linked report",
					"workDescription":     "Partial linked report",
					"hoursOnSite":         8,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 10,
					"vehicleType":         "company",
				})
				insertWorkRecordSubjectForTest(t, app, map[string]any{
					"subjectId":           "wrsubpartialln2",
					"workRecordId":        "wrpartiallink1a",
					"uid":                 "rzr98oadsp9qc11",
					"hoursOnSite":         7,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 12,
					"vehicleType":         "company",
				})
				linkTimeEntryToWorkRecordSubjectForTest(t, app, "1v2bciifdg8uv5c", "wrsubpartialln1")
			},
		},
		{
			Name:   "report holder can fetch filtered work records report",
			Method: http.MethodGet,
			URL:    "/api/reports/work_records/search?job=cjf0kt0defhq480&type=oask5u50nka0km8&date_from=2026-04-01&date_to=2026-04-30",
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
			Name:   "employee-filtered report displays the matching worker row",
			Method: http.MethodGet,
			URL:    "/api/reports/work_records/search?uid=rzr98oadsp9qc11&job=cjf0kt0defhq480&date_from=2031-02-01&date_to=2031-02-28",
			Headers: map[string]string{
				"Authorization": reportToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"total":1`,
				`"work_record":"W3102-0001"`,
				`"uid":"rzr98oadsp9qc11"`,
				`"owner_name":"Tester Time"`,
				`"linked_time_entry_id":"1v2bciifdg8uv5c"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, _ *core.ServeEvent) {
				insertWorkRecordWithSubjectForTest(t, app, map[string]any{
					"id":                  "wrfilteruid001",
					"subjectId":           "aaafilter00001",
					"number":              "W3102-0001",
					"date":                "2031-02-03",
					"uid":                 "f2j5a8vk006baub",
					"creator":             "f2j5a8vk006baub",
					"job":                 "cjf0kt0defhq480",
					"type":                "oask5u50nka0km8",
					"location":            "Filtered worker seeded",
					"workDescription":     "Filtered worker seeded",
					"hoursOnSite":         8,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 10,
					"vehicleType":         "company",
				})
				insertWorkRecordSubjectForTest(t, app, map[string]any{
					"subjectId":                "zzzfilter00001",
					"workRecordId":             "wrfilteruid001",
					"uid":                      "rzr98oadsp9qc11",
					"hoursOnSite":              6,
					"hoursTravelTime":          1,
					"distanceTravelledKM":      12,
					"vehicleType":              "personal",
					"companyVehicleUnitNumber": "",
				})
				linkTimeEntryToWorkRecordSubjectForTest(t, app, "1v2bciifdg8uv5c", "zzzfilter00001")
			},
			AfterTestFunc: func(t testing.TB, _ *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var payload workRecordSearchResponse
				if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
					t.Fatalf("failed decoding employee-filtered work record report response: %v", err)
				}
				if payload.Total != 1 || len(payload.Items) != 1 {
					t.Fatalf("expected one employee-filtered report row, got %#v", payload)
				}
				row := payload.Items[0]
				if row.UID != "rzr98oadsp9qc11" || row.OwnerName != "Tester Time" || row.LinkedTimeEntryID != "1v2bciifdg8uv5c" {
					t.Fatalf("expected filtered row to display matching worker fields, got %#v", row)
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
				`"total":4`,
			},
			TestAppFactory: testutils.SetupTestApp,
		},
		{
			Name:   "pending approval report includes partially approved multi-worker records",
			Method: http.MethodGet,
			URL:    "/api/reports/work_records/pending_approval",
			Headers: map[string]string{
				"Authorization": reportToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"work_record":"W3101-0002"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, _ *core.ServeEvent) {
				insertWorkRecordWithSubjectForTest(t, app, map[string]any{
					"id":                  "wrpartialappr1a",
					"subjectId":           "wrsubpartialap1",
					"number":              "W3101-0002",
					"date":                "2031-01-03",
					"uid":                 "f2j5a8vk006baub",
					"creator":             "f2j5a8vk006baub",
					"job":                 "cjf0kt0defhq480",
					"type":                "rqayp9yoz1qe6ui",
					"location":            "Partial approved report",
					"workDescription":     "Partial approved report",
					"hoursOnSite":         8,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 10,
					"vehicleType":         "company",
				})
				setWorkRecordSubjectsApprovedForTest(t, app, "wrpartialappr1a", true)
				insertWorkRecordSubjectForTest(t, app, map[string]any{
					"subjectId":           "wrsubpartialap2",
					"workRecordId":        "wrpartialappr1a",
					"uid":                 "rzr98oadsp9qc11",
					"hoursOnSite":         7,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 12,
					"vehicleType":         "company",
				})
			},
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
				if len(payload.Items) != 4 {
					t.Fatalf("expected 4 consumable report rows, got %#v", payload)
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
				if len(payload.Items) != 4 {
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
				if len(payload.Items) != 4 {
					t.Fatalf("expected full consumables response when offset is supplied, got %#v", payload)
				}
			},
		},
		{
			Name:   "time holder can load a self copy draft without creating a record",
			Method: http.MethodGet,
			URL:    "/api/work_records/wrtimecopy3002a/copy_to_tomorrow",
			Headers: map[string]string{
				"Authorization": regularUserToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"date":"2030-02-06"`,
				`"source_work_record":"W3002-1050"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var row workRecordCopyDraftResponse
				if err := json.NewDecoder(res.Body).Decode(&row); err != nil {
					t.Fatalf("failed decoding work record copy draft response: %v", err)
				}

				if row.CopyFromWorkRecordID != "wrtimecopy3002a" || row.RootWorkRecordID != "wrtimecopy3002a" {
					t.Fatalf("expected root copy draft metadata, got %#v", row)
				}
				if row.RootWorkRecord != "W3002-1050" || row.Date != "2030-02-06" {
					t.Fatalf("unexpected copy draft response: %#v", row)
				}
				if len(row.Consumables) != 2 {
					t.Fatalf("expected copied draft consumables, got %#v", row.Consumables)
				}
				if got := countWorkRecordsForTest(t, app); got != 5 {
					t.Fatalf("expected copy draft load to avoid inserts, got %d work records", got)
				}
			},
		},
		{
			Name:   "time holder can save a copied self draft with copied consumables",
			Method: http.MethodPost,
			URL:    "/api/work_records/save",
			Headers: map[string]string{
				"Authorization": regularUserToken,
				"Content-Type":  "application/json",
			},
			Body: workRecordRequestBody(`{
				"uid":"rzr98oadsp9qc11",
				"date":"2030-02-07",
				"job":"cjf0kt0defhq480",
				"type":"oask5u50nka0km8",
				"location":"Time seeded location",
				"work_description":"Time seeded work description",
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
				"copy_from_work_record_id":"wrtimecopy3002a",
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
							COALESCE(wr.number, '') AS number,
							COALESCE(wr.date, '') AS date,
							COALESCE(wrs.uid, '') AS uid,
							COALESCE(wr.job, '') AS job,
							COALESCE(wr.parent_work_record, '') AS parent_work_record
						FROM work_records wr
						LEFT JOIN work_records_subjects wrs
						  ON wrs.work_record = wr.id
						WHERE wr.id = {:id}
						ORDER BY wrs.id
						LIMIT 1
					`).Bind(map[string]any{
					"id": row.ID,
				}).One(&record); err != nil {
					t.Fatalf("expected copied work record to be saved: %v", err)
				}

				if record.ParentWorkRecord != "wrtimecopy3002a" || record.Date != "2030-02-07" {
					t.Fatalf("unexpected copied work record lineage/date: %#v", record)
				}
				if record.UID != "rzr98oadsp9qc11" || record.Job != "cjf0kt0defhq480" {
					t.Fatalf("expected copied work record to keep source worker/job, got %#v", record)
				}
				if record.Number != "W3002-1052" {
					t.Fatalf("expected copied work record number to continue seeded February sequence, got %q", record.Number)
				}
				if got := countWorkRecordConsumableEntriesForTest(t, app, row.ID); got != 2 {
					t.Fatalf("expected copied consumable entries, got %d", got)
				}
			},
		},
		{
			Name:   "time holder cannot tamper copied draft worker",
			Method: http.MethodPost,
			URL:    "/api/work_records/save",
			Headers: map[string]string{
				"Authorization": regularUserToken,
				"Content-Type":  "application/json",
			},
			Body: workRecordRequestBody(`{
				"uid":"rzr98oadsp9qc11",
				"date":"2030-02-07",
				"job":"cjf0kt0defhq480",
				"type":"oask5u50nka0km8",
				"location":"Time seeded location",
				"work_description":"Time seeded work description",
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
				"subjects":[
					{
						"id":"",
						"uid":"dkv192wxprcqmho",
						"hours_on_site":8,
						"hours_travel_time":1,
						"distance_travelled_km":10,
						"vehicle_type":"company",
						"is_passenger":false,
						"company_vehicle_unit_number":"7"
					}
				],
				"copy_from_work_record_id":"wrtimecopy3002a",
				"consumables":[]
			}`),
			ExpectedStatus: http.StatusForbidden,
			ExpectedContent: []string{
				`"message":"you are not authorized to create this work record"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, _ *http.Response) {
				if got := countWorkRecordsForTest(t, app); got != 5 {
					t.Fatalf("expected rejected copied worker tampering to avoid inserts, got %d work records", got)
				}
			},
		},
		{
			Name:   "work record holder without time cannot create a self-owned work record",
			Method: http.MethodPost,
			URL:    "/api/work_records/save",
			Headers: map[string]string{
				"Authorization": workRecordToken,
				"Content-Type":  "application/json",
			},
			Body: workRecordRequestBody(`{
				"uid":"f2j5a8vk006baub",
				"date":"2026-04-11",
				"job":"cjf0kt0defhq480",
				"type":"oask5u50nka0km8",
				"location":"Unauthorized self create",
				"work_description":"Unauthorized self create",
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
			ExpectedStatus: http.StatusForbidden,
			ExpectedContent: []string{
				`"you are not authorized to create this work record"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, _ *http.Response) {
				if got := countWorkRecordsForTest(t, app); got != 5 {
					t.Fatalf("expected unauthorized self create to avoid inserts, got %d work records", got)
				}
			},
		},
		{
			Name:   "self create rejects schema-required work record fields when payload is otherwise empty",
			Method: http.MethodPost,
			URL:    "/api/work_records/save",
			Headers: map[string]string{
				"Authorization": regularUserToken,
				"Content-Type":  "application/json",
			},
			Body: workRecordRequestBody(`{
				"uid":"rzr98oadsp9qc11",
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

				requiredFields := []string{"date", "job", "type", "location", "work_description"}
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

				if got := countWorkRecordsForTest(t, app); got != 5 {
					t.Fatalf("expected failed empty create to avoid inserts, got %d work records", got)
				}
			},
		},
		{
			Name:   "self create rejects invalid field book number format through custom save route",
			Method: http.MethodPost,
			URL:    "/api/work_records/save",
			Headers: map[string]string{
				"Authorization": regularUserToken,
				"Content-Type":  "application/json",
			},
			Body: workRecordRequestBody(`{
				"uid":"rzr98oadsp9qc11",
				"date":"2026-04-11",
				"job":"cjf0kt0defhq480",
				"type":"oask5u50nka0km8",
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
				if got := countWorkRecordsForTest(t, app); got != 5 {
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
			Body: workRecordRequestBody(`{
				"date":"2026-04-10",
				"job":"cjf0kt0defhq480",
				"type":"oask5u50nka0km8",
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
				if got := countWorkRecordsForTest(t, app); got != 5 {
					t.Fatalf("expected failed edit to avoid inserts, got %d work records", got)
				}
			},
		},
		{
			Name:   "self create rejects field book page number below schema minimum through custom save route",
			Method: http.MethodPost,
			URL:    "/api/work_records/save",
			Headers: map[string]string{
				"Authorization": regularUserToken,
				"Content-Type":  "application/json",
			},
			Body: workRecordRequestBody(`{
				"uid":"rzr98oadsp9qc11",
				"date":"2026-04-11",
				"job":"cjf0kt0defhq480",
				"type":"oask5u50nka0km8",
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
				if got := countWorkRecordsForTest(t, app); got != 5 {
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
			Body: workRecordRequestBody(`{
				"date":"2026-04-10",
				"job":"cjf0kt0defhq480",
				"type":"oask5u50nka0km8",
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
				if got := countWorkRecordsForTest(t, app); got != 5 {
					t.Fatalf("expected failed edit to avoid inserts, got %d work records", got)
				}
			},
		},
		{
			Name:   "self create rejects company vehicle unit values outside schema pattern through custom save route",
			Method: http.MethodPost,
			URL:    "/api/work_records/save",
			Headers: map[string]string{
				"Authorization": regularUserToken,
				"Content-Type":  "application/json",
			},
			Body: workRecordRequestBody(`{
				"uid":"rzr98oadsp9qc11",
				"date":"2026-04-11",
				"job":"cjf0kt0defhq480",
				"type":"oask5u50nka0km8",
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
				if got := countWorkRecordsForTest(t, app); got != 5 {
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
			Body: workRecordRequestBody(`{
				"date":"2026-04-10",
				"job":"cjf0kt0defhq480",
				"type":"oask5u50nka0km8",
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
				if got := countWorkRecordsForTest(t, app); got != 5 {
					t.Fatalf("expected failed edit to avoid inserts, got %d work records", got)
				}
			},
		},
		{
			Name:   "self create rejects invalid job relation through custom save route",
			Method: http.MethodPost,
			URL:    "/api/work_records/save",
			Headers: map[string]string{
				"Authorization": regularUserToken,
				"Content-Type":  "application/json",
			},
			Body: workRecordRequestBody(`{
				"uid":"rzr98oadsp9qc11",
				"date":"2026-04-11",
				"job":"notarealjobid",
				"type":"oask5u50nka0km8",
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
				if got := countWorkRecordsForTest(t, app); got != 5 {
					t.Fatalf("expected failed create to avoid inserts, got %d work records", got)
				}
			},
		},
		{
			Name:   "self create rejects invalid type relation through custom save route",
			Method: http.MethodPost,
			URL:    "/api/work_records/save",
			Headers: map[string]string{
				"Authorization": regularUserToken,
				"Content-Type":  "application/json",
			},
			Body: workRecordRequestBody(`{
				"uid":"rzr98oadsp9qc11",
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
				if got := countWorkRecordsForTest(t, app); got != 5 {
					t.Fatalf("expected failed create to avoid inserts, got %d work records", got)
				}
			},
		},
		{
			Name:   "self create ignores forged boolean input kind when numeric quantity is provided",
			Method: http.MethodPost,
			URL:    "/api/work_records/save",
			Headers: map[string]string{
				"Authorization": regularUserToken,
				"Content-Type":  "application/json",
			},
			Body: workRecordRequestBody(`{
					"uid":"rzr98oadsp9qc11",
					"date":"2026-04-11",
					"job":"cjf0kt0defhq480",
					"type":"oask5u50nka0km8",
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
			Name:   "self create ignores forged numeric input kind for boolean consumables",
			Method: http.MethodPost,
			URL:    "/api/work_records/save",
			Headers: map[string]string{
				"Authorization": regularUserToken,
				"Content-Type":  "application/json",
			},
			Body: workRecordRequestBody(`{
					"uid":"rzr98oadsp9qc11",
					"date":"2026-04-11",
					"job":"cjf0kt0defhq480",
					"type":"oask5u50nka0km8",
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
			Body: workRecordRequestBody(`{
					"date":"2026-04-10",
					"job":"cjf0kt0defhq480",
					"type":"oask5u50nka0km8",
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
			Name:   "self copied draft duplicate-date conflicts are raised at save time",
			Method: http.MethodPost,
			URL:    "/api/work_records/save",
			Headers: map[string]string{
				"Authorization": regularUserToken,
				"Content-Type":  "application/json",
			},
			Body: workRecordRequestBody(`{
				"uid":"rzr98oadsp9qc11",
				"date":"2030-02-06",
				"job":"cjf0kt0defhq480",
				"type":"oask5u50nka0km8",
				"location":"Time seeded location",
				"work_description":"Time seeded work description",
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
				"copy_from_work_record_id":"wrtimecopy3002a",
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
				"Authorization": regularUserToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"date":"2026-04-05"`,
				`"root_work_record_id":"wr260400000003a"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, _ *core.ServeEvent) {
				insertWorkRecordWithSubjectForTest(t, app, map[string]any{
					"id":                  "wrchild0000001a",
					"subjectId":           "wrsubchild001",
					"number":              "W9901-0001",
					"date":                "2026-04-04",
					"uid":                 "rzr98oadsp9qc11",
					"creator":             "rzr98oadsp9qc11",
					"job":                 "cjf0kt0defhq480",
					"type":                "oask5u50nka0km8",
					"location":            "Seeded child location",
					"workDescription":     "Seeded child work description",
					"hoursOnSite":         8,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 10,
					"vehicleType":         "company",
					"created":             "2026-04-04 12:00:00.000Z",
					"updated":             "2026-04-04 12:00:00.000Z",
					"parentWorkRecord":    "wr260400000003a",
				})
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
				if got := countWorkRecordsForTest(t, app); got != 6 {
					t.Fatalf("expected copy draft to avoid inserts after seeding child, got %d work records", got)
				}
			},
		},
		{
			Name:   "copied draft save uses the edited date for numbering",
			Method: http.MethodPost,
			URL:    "/api/work_records/save",
			Headers: map[string]string{
				"Authorization": regularUserToken,
				"Content-Type":  "application/json",
			},
			Body: workRecordRequestBody(`{
				"uid":"rzr98oadsp9qc11",
				"date":"2026-03-01",
				"job":"cjf0kt0defhq480",
				"type":"oask5u50nka0km8",
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
				insertWorkRecordWithSubjectForTest(t, app, map[string]any{
					"id":                  "wrjan310000001a",
					"subjectId":           "wrsubjan31001",
					"number":              "W2601-0001",
					"date":                "2026-01-31",
					"uid":                 "rzr98oadsp9qc11",
					"creator":             "rzr98oadsp9qc11",
					"job":                 "cjf0kt0defhq480",
					"type":                "oask5u50nka0km8",
					"location":            "January location",
					"workDescription":     "January work description",
					"hoursOnSite":         8,
					"hoursTravelTime":     1,
					"distanceTravelledKM": 10,
					"vehicleType":         "company",
					"created":             "2026-01-31 12:00:00.000Z",
					"updated":             "2026-01-31 12:00:00.000Z",
				})
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
			Body: workRecordRequestBody(`{
				"uid":"f2j5a8vk006baub",
				"date":"2026-04-11",
				"job":"cjf0kt0defhq480",
				"type":"oask5u50nka0km8",
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
