package main

import (
	"encoding/json"
	"net/http"
	"testing"
	"tybalt/internal/testutils"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

type workRecordListRow struct {
	ID                 string `json:"id"`
	WorkRecord         string `json:"work_record"`
	Date               string `json:"date"`
	JobNumber          string `json:"job_number"`
	Approved           bool   `json:"approved"`
	LinkedTimeEntryID  string `json:"linked_time_entry_id"`
	SearchText         string `json:"search_text"`
	CanCreateTimeEntry bool   `json:"can_create_time_entry"`
}

type workRecordDetailResponse struct {
	ID                         string  `json:"id"`
	WorkRecord                 string  `json:"work_record"`
	Date                       string  `json:"date"`
	UID                        string  `json:"uid"`
	JobID                      string  `json:"job_id"`
	JobNumber                  string  `json:"job_number"`
	TypeName                   string  `json:"type_name"`
	Location                   string  `json:"location"`
	WorkDescription            string  `json:"work_description"`
	HoursOnSite                float64 `json:"hours_on_site"`
	HoursTravelTime            float64 `json:"hours_travel_time"`
	LinkedTimeEntryID          string  `json:"linked_time_entry_id"`
	LinkedTimeSheetID          string  `json:"linked_timesheet_id"`
	LinkedTimeEntryDate        string  `json:"linked_time_entry_date"`
	LinkedTimeEntryDescription string  `json:"linked_time_entry_description"`
	CanCreateTimeEntry         bool    `json:"can_create_time_entry"`
}

type workRecordCopyResponse struct {
	ID               string `json:"id"`
	WorkRecord       string `json:"work_record"`
	Date             string `json:"date"`
	ParentWorkRecord string `json:"parent_work_record"`
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
			Name:   "report holder can fetch first class work records list",
			Method: http.MethodGet,
			URL:    "/api/work_records",
			Headers: map[string]string{
				"Authorization": reportToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"work_record":"W3002-0050"`,
				`"work_record":"W2604-0010"`,
				`"work_record":"W2604-0003"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, _ *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var rows []workRecordListRow
				if err := json.NewDecoder(res.Body).Decode(&rows); err != nil {
					t.Fatalf("failed decoding work records list response: %v", err)
				}

				if len(rows) != 3 {
					t.Fatalf("expected 3 work records, got %d: %#v", len(rows), rows)
				}
				if rows[0].WorkRecord != "W3002-0050" || rows[1].WorkRecord != "W2604-0010" || rows[2].WorkRecord != "W2604-0003" {
					t.Fatalf("unexpected row order: %#v", rows)
				}
				if rows[0].CanCreateTimeEntry {
					t.Fatalf("report holder should not see create action on %#v", rows[0])
				}
				if rows[0].SearchText == "" {
					t.Fatalf("expected populated search_text for %#v", rows[0])
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
				`"can_create_time_entry":true`,
			},
			TestAppFactory: testutils.SetupTestApp,
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
				`[]`,
			},
			TestAppFactory: testutils.SetupTestApp,
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
				if row.LinkedTimeEntryID != "" || row.LinkedTimeSheetID != "" {
					t.Fatalf("expected seeded work record to be unlinked: %#v", row)
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
			Name:   "work record holder can copy a root work record to tomorrow",
			Method: http.MethodPost,
			URL:    "/api/work_records/wr260400000003a/copy_to_tomorrow",
			Headers: map[string]string{
				"Authorization": workRecordToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"date":"2026-04-04"`,
				`"parent_work_record":"wr260400000003a"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, _ *core.ServeEvent) {
				if err := app.RunAllMigrations(); err != nil {
					t.Fatalf("failed to apply migrations in test app: %v", err)
				}
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var row workRecordCopyResponse
				if err := json.NewDecoder(res.Body).Decode(&row); err != nil {
					t.Fatalf("failed decoding work record copy response: %v", err)
				}

				if row.ID == "" || row.WorkRecord == "" {
					t.Fatalf("expected copied work record identifiers, got %#v", row)
				}
				if row.ParentWorkRecord != "wr260400000003a" || row.Date != "2026-04-04" {
					t.Fatalf("unexpected copy response: %#v", row)
				}

				var record struct {
					ParentWorkRecord string `db:"parent_work_record"`
				}
				if err := app.DB().NewQuery(`
					SELECT COALESCE(parent_work_record, '') AS parent_work_record
					FROM work_records
					WHERE id = {:id}
				`).Bind(map[string]any{
					"id": row.ID,
				}).One(&record); err != nil {
					t.Fatalf("expected copied work record to be saved: %v", err)
				}
				if record.ParentWorkRecord != "wr260400000003a" {
					t.Fatalf("expected root parent to point at source, got %q", record.ParentWorkRecord)
				}
			},
		},
		{
			Name:   "copying a child work record keeps the root parent",
			Method: http.MethodPost,
			URL:    "/api/work_records/wrchild0000001a/copy_to_tomorrow",
			Headers: map[string]string{
				"Authorization": workRecordToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"date":"2026-04-05"`,
				`"parent_work_record":"wr260400000003a"`,
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

				var row workRecordCopyResponse
				if err := json.NewDecoder(res.Body).Decode(&row); err != nil {
					t.Fatalf("failed decoding work record copy response: %v", err)
				}

				var record struct {
					ParentWorkRecord string `db:"parent_work_record"`
				}
				if err := app.DB().NewQuery(`
					SELECT COALESCE(parent_work_record, '') AS parent_work_record
					FROM work_records
					WHERE id = {:id}
				`).Bind(map[string]any{
					"id": row.ID,
				}).One(&record); err != nil {
					t.Fatalf("expected copied child work record to be saved: %v", err)
				}
				if record.ParentWorkRecord != "wr260400000003a" {
					t.Fatalf("expected copied child to keep root parent, got %q", record.ParentWorkRecord)
				}
			},
		},
		{
			Name:   "copy to tomorrow uses target date month for numbering",
			Method: http.MethodPost,
			URL:    "/api/work_records/wrjan310000001a/copy_to_tomorrow",
			Headers: map[string]string{
				"Authorization": workRecordToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"date":"2026-02-01"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, _ *core.ServeEvent) {
				if err := app.RunAllMigrations(); err != nil {
					t.Fatalf("failed to apply migrations in test app: %v", err)
				}

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
			AfterTestFunc: func(t testing.TB, _ *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var row workRecordCopyResponse
				if err := json.NewDecoder(res.Body).Decode(&row); err != nil {
					t.Fatalf("failed decoding work record copy response: %v", err)
				}

				if row.Date != "2026-02-01" {
					t.Fatalf("expected copied date to cross month boundary, got %#v", row)
				}
				if len(row.WorkRecord) < len("W2602-") || row.WorkRecord[:len("W2602-")] != "W2602-" {
					t.Fatalf("expected February-prefixed work record number, got %q", row.WorkRecord)
				}
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
