package routes

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"tybalt/errs"
	apphooks "tybalt/hooks"
	"tybalt/utilities"

	"github.com/pocketbase/dbx"
	validation "github.com/pocketbase/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

const (
	workRecordsCollectionIdentifier                 = "work_records"
	workRecordNotesCollectionIdentifier             = "work_records_notes"
	workRecordConsumableEntriesCollectionIdentifier = "work_records_consumable_entries"
)

type WorkRecordListRow struct {
	ID                   string                      `db:"id" json:"id"`
	WorkRecord           string                      `db:"work_record" json:"work_record"`
	Date                 string                      `db:"date" json:"date"`
	Creator              string                      `db:"creator" json:"creator"`
	CreatorName          string                      `db:"creator_name" json:"creator_name"`
	JobID                string                      `db:"job_id" json:"job_id"`
	JobNumber            string                      `db:"job_number" json:"job_number"`
	TypeID               string                      `db:"type_id" json:"type_id"`
	TypeName             string                      `db:"type_name" json:"type_name"`
	Location             string                      `db:"location" json:"location"`
	WorkDescription      string                      `db:"work_description" json:"work_description"`
	Attachment           string                      `db:"attachment" json:"attachment"`
	Approved             bool                        `db:"approved" json:"approved"`
	SearchText           string                      `json:"search_text"`
	CanCreateTimeEntry   bool                        `json:"can_create_time_entry"`
	SubjectCount         int                         `db:"subject_count" json:"subject_count"`
	ApprovedSubjectCount int                         `db:"approved_subject_count" json:"approved_subject_count"`
	Subjects             []WorkRecordSubjectResponse `json:"subjects"`
}

type WorkRecordNoteResponse struct {
	ID         string `db:"id" json:"id"`
	UID        string `db:"uid" json:"uid"`
	AuthorName string `db:"author_name" json:"author_name"`
	Note       string `db:"note" json:"note"`
	Created    string `db:"created" json:"created"`
}

type WorkRecordConsumableResponse struct {
	EntryID         string  `db:"entry_id" json:"entry_id"`
	ConsumableID    string  `db:"consumable_id" json:"consumable_id"`
	ConsumableName  string  `db:"consumable_name" json:"consumable_name"`
	InputKind       string  `db:"input_kind" json:"input_kind"`
	UnitLabel       string  `db:"unit_label" json:"unit_label"`
	SortOrder       int     `db:"sort_order" json:"sort_order"`
	QuantityNumber  float64 `db:"quantity_number" json:"quantity_number"`
	BooleanSelected bool    `db:"boolean_selected" json:"boolean_selected"`
}

type WorkRecordSubjectResponse struct {
	ID                         string  `db:"id" json:"id"`
	WorkRecordID               string  `db:"work_record" json:"work_record_id"`
	UID                        string  `db:"uid" json:"uid"`
	WorkerName                 string  `db:"worker_name" json:"worker_name"`
	HoursOnSite                float64 `db:"hours_on_site" json:"hours_on_site"`
	HoursTravelTime            float64 `db:"hours_travel_time" json:"hours_travel_time"`
	DistanceTravelledKM        float64 `db:"distance_travelled_km" json:"distance_travelled_km"`
	VehicleType                string  `db:"vehicle_type" json:"vehicle_type"`
	IsPassenger                bool    `db:"is_passenger" json:"is_passenger"`
	CompanyVehicleUnitNumber   string  `db:"company_vehicle_unit_number" json:"company_vehicle_unit_number"`
	Approved                   bool    `db:"approved" json:"approved"`
	LinkedTimeEntryID          string  `db:"linked_time_entry_id" json:"linked_time_entry_id"`
	LinkedTimeSheetID          string  `db:"linked_timesheet_id" json:"linked_timesheet_id"`
	LinkedTimeEntryDate        string  `db:"linked_time_entry_date" json:"linked_time_entry_date"`
	LinkedTimeEntryDescription string  `db:"linked_time_entry_description" json:"linked_time_entry_description"`
	LinkedTimeSheetWeekEnding  string  `db:"linked_timesheet_week_ending" json:"linked_timesheet_week_ending"`
	CanCreateTimeEntry         bool    `json:"can_create_time_entry"`
	CanEdit                    bool    `json:"can_edit"`
}

type WorkRecordDetailResponse struct {
	ID                    string                         `json:"id"`
	WorkRecord            string                         `json:"work_record"`
	Date                  string                         `json:"date"`
	Creator               string                         `json:"creator"`
	CreatorName           string                         `json:"creator_name"`
	JobID                 string                         `json:"job_id"`
	JobNumber             string                         `json:"job_number"`
	ClientName            string                         `json:"client_name"`
	JobDescription        string                         `json:"job_description"`
	ClientReferenceNumber string                         `json:"client_reference_number"`
	ClientPO              string                         `json:"client_po"`
	TypeID                string                         `json:"type_id"`
	TypeName              string                         `json:"type_name"`
	Location              string                         `json:"location"`
	WorkDescription       string                         `json:"work_description"`
	ReportTo              string                         `json:"report_to"`
	LaneKms               float64                        `json:"lane_kms"`
	FieldBookNumber       string                         `json:"field_book_number"`
	FieldBookPageNumber   int                            `json:"field_book_page_number"`
	SubContractor         string                         `json:"sub_contractor"`
	Equipment             string                         `json:"equipment"`
	Supplies              string                         `json:"supplies"`
	Attachment            string                         `json:"attachment"`
	Approved              bool                           `json:"approved"`
	ParentWorkRecordID    string                         `json:"parent_work_record_id"`
	ParentWorkRecord      string                         `json:"parent_work_record"`
	RootWorkRecordID      string                         `json:"root_work_record_id"`
	RootWorkRecord        string                         `json:"root_work_record"`
	CanCopyToTomorrow     bool                           `json:"can_copy_to_tomorrow"`
	CanEdit               bool                           `json:"can_edit"`
	CanEditShared         bool                           `json:"can_edit_shared"`
	CanAddNote            bool                           `json:"can_add_note"`
	Subjects              []WorkRecordSubjectResponse    `json:"subjects"`
	Notes                 []WorkRecordNoteResponse       `json:"notes"`
	Consumables           []WorkRecordConsumableResponse `json:"consumables"`
}

type WorkRecordSaveConsumableRequest struct {
	EntryID         string  `json:"entry_id"`
	ConsumableID    string  `json:"consumable_id"`
	InputKind       string  `json:"input_kind"`
	QuantityNumber  float64 `json:"quantity_number"`
	BooleanSelected bool    `json:"boolean_selected"`
}

type WorkRecordSubjectSaveRequest struct {
	ID                       string  `json:"id"`
	UID                      string  `json:"uid"`
	HoursOnSite              float64 `json:"hours_on_site"`
	HoursTravelTime          float64 `json:"hours_travel_time"`
	DistanceTravelledKM      float64 `json:"distance_travelled_km"`
	VehicleType              string  `json:"vehicle_type"`
	IsPassenger              bool    `json:"is_passenger"`
	CompanyVehicleUnitNumber string  `json:"company_vehicle_unit_number"`
}

type WorkRecordCreateRequest struct {
	Date                 string                            `json:"date"`
	Job                  string                            `json:"job"`
	Type                 string                            `json:"type"`
	Location             string                            `json:"location"`
	WorkDescription      string                            `json:"work_description"`
	ReportTo             string                            `json:"report_to"`
	LaneKms              float64                           `json:"lane_kms"`
	FieldBookNumber      string                            `json:"field_book_number"`
	FieldBookPageNumber  any                               `json:"field_book_page_number"`
	SubContractor        string                            `json:"sub_contractor"`
	Equipment            string                            `json:"equipment"`
	Supplies             string                            `json:"supplies"`
	Subjects             []WorkRecordSubjectSaveRequest    `json:"subjects"`
	Consumables          []WorkRecordSaveConsumableRequest `json:"consumables"`
	CopyFromWorkRecordID string                            `json:"copy_from_work_record_id"`
}

type WorkRecordSaveRequest struct {
	Date                string                            `json:"date"`
	Job                 string                            `json:"job"`
	Type                string                            `json:"type"`
	Location            string                            `json:"location"`
	WorkDescription     string                            `json:"work_description"`
	ReportTo            string                            `json:"report_to"`
	LaneKms             float64                           `json:"lane_kms"`
	FieldBookNumber     string                            `json:"field_book_number"`
	FieldBookPageNumber any                               `json:"field_book_page_number"`
	SubContractor       string                            `json:"sub_contractor"`
	Equipment           string                            `json:"equipment"`
	Supplies            string                            `json:"supplies"`
	Subjects            []WorkRecordSubjectSaveRequest    `json:"subjects"`
	Consumables         []WorkRecordSaveConsumableRequest `json:"consumables"`
}

type WorkRecordSaveResponse struct {
	ID string `json:"id"`
}

type WorkRecordAttachmentResponse struct {
	ID         string `json:"id"`
	Attachment string `json:"attachment"`
}

type WorkRecordNoteCreateRequest struct {
	Note string `json:"note"`
}

type WorkRecordCopyDraftResponse struct {
	Date                 string                         `json:"date"`
	Job                  string                         `json:"job"`
	Type                 string                         `json:"type"`
	Location             string                         `json:"location"`
	WorkDescription      string                         `json:"work_description"`
	ReportTo             string                         `json:"report_to"`
	LaneKms              float64                        `json:"lane_kms"`
	FieldBookNumber      string                         `json:"field_book_number"`
	FieldBookPageNumber  int                            `json:"field_book_page_number"`
	SubContractor        string                         `json:"sub_contractor"`
	Equipment            string                         `json:"equipment"`
	Supplies             string                         `json:"supplies"`
	Subjects             []WorkRecordSubjectResponse    `json:"subjects"`
	Consumables          []WorkRecordConsumableResponse `json:"consumables"`
	CopyFromWorkRecordID string                         `json:"copy_from_work_record_id"`
	SourceWorkRecord     string                         `json:"source_work_record"`
	RootWorkRecordID     string                         `json:"root_work_record_id"`
	RootWorkRecord       string                         `json:"root_work_record"`
}

type WorkRecordOption struct {
	ID         string `db:"id" json:"id"`
	WorkRecord string `db:"work_record" json:"work_record"`
}

type workRecordDetailRow struct {
	ID                     string  `db:"id"`
	WorkRecord             string  `db:"work_record"`
	Date                   string  `db:"date"`
	Creator                string  `db:"creator"`
	CreatorName            string  `db:"creator_name"`
	JobID                  string  `db:"job_id"`
	JobNumber              string  `db:"job_number"`
	ClientName             string  `db:"client_name"`
	JobDescription         string  `db:"job_description"`
	ClientReferenceNumber  string  `db:"client_reference_number"`
	ClientPO               string  `db:"client_po"`
	TypeID                 string  `db:"type_id"`
	TypeName               string  `db:"type_name"`
	Location               string  `db:"location"`
	WorkDescription        string  `db:"work_description"`
	ReportTo               string  `db:"report_to"`
	LaneKms                float64 `db:"lane_kms"`
	FieldBookNumber        string  `db:"field_book_number"`
	FieldBookPageNumber    int     `db:"field_book_page_number"`
	SubContractor          string  `db:"sub_contractor"`
	Equipment              string  `db:"equipment"`
	Supplies               string  `db:"supplies"`
	Attachment             string  `db:"attachment"`
	Approved               bool    `db:"approved"`
	ParentWorkRecordID     string  `db:"parent_work_record_id"`
	ParentWorkRecordNumber string  `db:"parent_work_record_number"`
}

type workRecordCopySourceRow struct {
	ID                  string  `db:"id"`
	WorkRecord          string  `db:"work_record"`
	Creator             string  `db:"creator"`
	Date                string  `db:"date"`
	Job                 string  `db:"job"`
	Type                string  `db:"type"`
	Location            string  `db:"location"`
	WorkDescription     string  `db:"work_description"`
	ReportTo            string  `db:"report_to"`
	LaneKms             float64 `db:"lane_kms"`
	FieldBookNumber     string  `db:"field_book_number"`
	FieldBookPageNumber int     `db:"field_book_page_number"`
	SubContractor       string  `db:"sub_contractor"`
	Equipment           string  `db:"equipment"`
	Supplies            string  `db:"supplies"`
	ParentWorkRecordID  string  `db:"parent_work_record_id"`
	ParentWorkRecord    string  `db:"parent_work_record"`
}

type workRecordConsumableEntryState struct {
	ID           string `db:"id"`
	ConsumableID string `db:"consumable_id"`
}

type workRecordSaveStateRow struct {
	ID                  string  `db:"id"`
	Number              string  `db:"number"`
	Date                string  `db:"date"`
	Creator             string  `db:"creator"`
	Job                 string  `db:"job"`
	Type                string  `db:"type"`
	Location            string  `db:"location"`
	WorkDescription     string  `db:"work_description"`
	ReportTo            string  `db:"report_to"`
	LaneKms             float64 `db:"lane_kms"`
	FieldBookNumber     string  `db:"field_book_number"`
	FieldBookPageNumber int     `db:"field_book_page_number"`
	SubContractor       string  `db:"sub_contractor"`
	Equipment           string  `db:"equipment"`
	Supplies            string  `db:"supplies"`
	ParentWorkRecord    string  `db:"parent_work_record"`
}

type workRecordConsumableEntryPersistenceRow struct {
	ID             string  `db:"id"`
	WorkRecordID   string  `db:"work_record"`
	ConsumableID   string  `db:"consumable"`
	QuantityNumber float64 `db:"quantity_number"`
}

func findWorkRecordCollection(app core.App) (*core.Collection, error) {
	return app.FindCollectionByNameOrId(workRecordsCollectionIdentifier)
}

func findWorkRecordConsumableEntriesCollection(app core.App) (*core.Collection, error) {
	return app.FindCollectionByNameOrId(workRecordConsumableEntriesCollectionIdentifier)
}

func findWorkRecordSubjectsCollection(app core.App) (*core.Collection, error) {
	return app.FindCollectionByNameOrId("work_records_subjects")
}

func findWorkRecordSubjectByID(app core.App, subjectID string) (*core.Record, error) {
	subjectID = strings.TrimSpace(subjectID)
	if subjectID == "" {
		return nil, sql.ErrNoRows
	}

	row := WorkRecordSubjectResponse{}
	if err := app.DB().NewQuery(`
		SELECT
			id,
			COALESCE(work_record, '') AS work_record,
			COALESCE(uid, '') AS uid,
			COALESCE(hours_on_site, 0) AS hours_on_site,
			COALESCE(hours_travel_time, 0) AS hours_travel_time,
			COALESCE(distance_travelled_km, 0) AS distance_travelled_km,
			COALESCE(vehicle_type, '') AS vehicle_type,
			COALESCE(is_passenger, 0) AS is_passenger,
			COALESCE(company_vehicle_unit_number, '') AS company_vehicle_unit_number,
			COALESCE(approved, 0) AS approved
		FROM work_records_subjects
		WHERE id = {:id}
		LIMIT 1
	`).Bind(dbx.Params{
		"id": subjectID,
	}).One(&row); err != nil {
		return nil, err
	}

	collection, err := findWorkRecordSubjectsCollection(app)
	if err != nil {
		return nil, err
	}
	record := core.NewRecord(collection)
	record.Load(map[string]any{
		"id":                          row.ID,
		"work_record":                 row.WorkRecordID,
		"uid":                         row.UID,
		"hours_on_site":               row.HoursOnSite,
		"hours_travel_time":           row.HoursTravelTime,
		"distance_travelled_km":       row.DistanceTravelledKM,
		"vehicle_type":                row.VehicleType,
		"is_passenger":                row.IsPassenger,
		"company_vehicle_unit_number": row.CompanyVehicleUnitNumber,
		"approved":                    row.Approved,
	})
	if err := record.PostScan(); err != nil {
		return nil, err
	}
	return record, nil
}

func findWorkRecordConsumableEntryByID(app core.App, entryID string) (*core.Record, error) {
	row := workRecordConsumableEntryPersistenceRow{}
	if err := app.DB().NewQuery(`
		SELECT
			id,
			COALESCE(work_record, '') AS work_record,
			COALESCE(consumable, '') AS consumable,
			COALESCE(quantity_number, 0) AS quantity_number
		FROM work_records_consumable_entries
		WHERE id = {:id}
		LIMIT 1
	`).Bind(dbx.Params{
		"id": entryID,
	}).One(&row); err != nil {
		return nil, err
	}

	collection, err := findWorkRecordConsumableEntriesCollection(app)
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(collection)
	record.Load(map[string]any{
		"id":              row.ID,
		"work_record":     row.WorkRecordID,
		"consumable":      row.ConsumableID,
		"quantity_number": row.QuantityNumber,
	})
	if err := record.PostScan(); err != nil {
		return nil, err
	}

	return record, nil
}

func canViewAllWorkRecords(app core.App, auth *core.Record) (bool, error) {
	if auth == nil {
		return false, nil
	}

	return utilities.HasClaim(app, auth, "report")
}

func parseWorkRecordListScope(e *core.RequestEvent) (string, error) {
	scope := strings.TrimSpace(e.Request.URL.Query().Get("scope"))
	if scope == "" {
		return "mine", nil
	}
	if scope != "mine" && scope != "all" {
		return "", workRecordSearchFieldError("scope", "invalid_scope", "scope must be mine or all")
	}
	return scope, nil
}

func loadVisibleWorkRecordRows(app core.App, auth *core.Record, scope string) ([]WorkRecordListRow, int, error) {
	filterSQL := "1 = 1"
	params := dbx.Params{}
	if scope == "mine" {
		filterSQL += " AND (wr.creator = {:authId} OR EXISTS (SELECT 1 FROM work_records_subjects wrs_visible WHERE wrs_visible.work_record = wr.id AND wrs_visible.uid = {:authId}))"
		params["authId"] = auth.Id
	}

	rows := []WorkRecordListRow{}
	if err := app.DB().NewQuery(`
		SELECT
			wr.id AS id,
			wr.number AS work_record,
			wr.date AS date,
			wr.creator AS creator,
			COALESCE(creator_profile.given_name || ' ' || creator_profile.surname, '') AS creator_name,
			COALESCE(wr.job, '') AS job_id,
			COALESCE(j.number, '') AS job_number,
			COALESCE(wr.type, '') AS type_id,
			COALESCE(wrt.name, '') AS type_name,
			COALESCE(wr.location, '') AS location,
			COALESCE(wr.work_description, '') AS work_description,
			COALESCE(wr.attachment, '') AS attachment,
			CASE WHEN COALESCE(subject_counts.approved_subject_count, 0) > 0 THEN 1 ELSE 0 END AS approved,
			COALESCE(subject_counts.subject_count, 0) AS subject_count,
			COALESCE(subject_counts.approved_subject_count, 0) AS approved_subject_count
		FROM work_records wr
		LEFT JOIN profiles creator_profile
		  ON creator_profile.uid = wr.creator
		LEFT JOIN (
			SELECT
				work_record,
				COUNT(*) AS subject_count,
				SUM(CASE WHEN COALESCE(approved, 0) = 1 THEN 1 ELSE 0 END) AS approved_subject_count
			FROM work_records_subjects
			GROUP BY work_record
		) subject_counts
		  ON subject_counts.work_record = wr.id
		LEFT JOIN jobs j
		  ON j.id = wr.job
		LEFT JOIN work_records_types wrt
		  ON wrt.id = wr.type
		WHERE ` + filterSQL + `
		ORDER BY wr.date DESC, wr.number DESC
	`).Bind(params).All(&rows); err != nil {
		return nil, 0, err
	}

	for i := range rows {
		subjects, err := loadWorkRecordSubjects(app, rows[i].ID, auth, rows[i].Creator)
		if err != nil {
			return nil, 0, err
		}
		rows[i].Subjects = subjects

		subjectSearchText := strings.Builder{}
		for _, subject := range subjects {
			subjectSearchText.WriteString(" ")
			subjectSearchText.WriteString(subject.WorkerName)
			if subject.UID == auth.Id {
				rows[i].CanCreateTimeEntry = subject.CanCreateTimeEntry
			}
		}

		rows[i].SearchText = buildWorkRecordSearchText(rows[i].WorkRecord + " " + rows[i].CreatorName + " " + rows[i].JobNumber + " " + rows[i].TypeName + " " + rows[i].Location + " " + rows[i].WorkDescription + subjectSearchText.String())
	}

	return rows, len(rows), nil
}

func canViewWorkRecord(app core.App, auth *core.Record, record *core.Record) (bool, error) {
	if auth == nil || record == nil {
		return false, nil
	}

	if record.GetString("creator") == auth.Id {
		return true, nil
	}

	var row struct {
		Count int `db:"count"`
	}
	if err := app.DB().NewQuery(`
		SELECT COUNT(*) AS count
		FROM work_records_subjects
		WHERE work_record = {:workRecordId}
		  AND uid = {:uid}
	`).Bind(dbx.Params{
		"workRecordId": record.Id,
		"uid":          auth.Id,
	}).One(&row); err != nil {
		return false, err
	}
	if row.Count > 0 {
		return true, nil
	}

	return canViewAllWorkRecords(app, auth)
}

func requireWorkRecordDetailViewer(app core.App, auth *core.Record, record *core.Record) error {
	allowed, err := canViewWorkRecord(app, auth, record)
	if err != nil {
		return err
	}
	if allowed {
		return nil
	}

	return &errs.HookError{
		Status:  http.StatusForbidden,
		Message: "you are not authorized to view this work record",
		Data: map[string]errs.CodeError{
			"global": {
				Code:    "unauthorized",
				Message: "you are not authorized to view this work record",
			},
		},
	}
}

func loadWorkRecordNotes(app core.App, workRecordID string) ([]WorkRecordNoteResponse, error) {
	rows := []WorkRecordNoteResponse{}
	if err := app.DB().NewQuery(`
		SELECT
			wrn.id AS id,
			wrn.uid AS uid,
			COALESCE(p.given_name || ' ' || p.surname, '') AS author_name,
			COALESCE(wrn.note, '') AS note,
			COALESCE(wrn.created, '') AS created
		FROM work_records_notes wrn
		LEFT JOIN profiles p
		  ON p.uid = wrn.uid
		WHERE wrn.work_record = {:workRecordId}
		ORDER BY wrn.created DESC
	`).Bind(dbx.Params{
		"workRecordId": workRecordID,
	}).All(&rows); err != nil {
		return nil, err
	}
	return rows, nil
}

func loadWorkRecordSubjects(app core.App, workRecordID string, auth *core.Record, creatorIDs ...string) ([]WorkRecordSubjectResponse, error) {
	rows := []WorkRecordSubjectResponse{}
	if err := app.DB().NewQuery(`
		SELECT
			wrs.id AS id,
			wrs.work_record AS work_record,
			wrs.uid AS uid,
			COALESCE(p.given_name || ' ' || p.surname, '') AS worker_name,
			COALESCE(wrs.hours_on_site, 0) AS hours_on_site,
			COALESCE(wrs.hours_travel_time, 0) AS hours_travel_time,
			COALESCE(wrs.distance_travelled_km, 0) AS distance_travelled_km,
			COALESCE(wrs.vehicle_type, '') AS vehicle_type,
			COALESCE(wrs.is_passenger, 0) AS is_passenger,
			COALESCE(wrs.company_vehicle_unit_number, '') AS company_vehicle_unit_number,
			COALESCE(wrs.approved, 0) AS approved,
			COALESCE(te.id, '') AS linked_time_entry_id,
			COALESCE(te.tsid, '') AS linked_timesheet_id,
			COALESCE(te.date, '') AS linked_time_entry_date,
			COALESCE(te.description, '') AS linked_time_entry_description,
			COALESCE(ts.week_ending, '') AS linked_timesheet_week_ending
		FROM work_records_subjects wrs
		LEFT JOIN profiles p
		  ON p.uid = wrs.uid
		LEFT JOIN time_entries te
		  ON te.work_record_subject_id = wrs.id
		LEFT JOIN time_sheets ts
		  ON ts.id = te.tsid
		WHERE wrs.work_record = {:workRecordId}
		ORDER BY p.surname, p.given_name, wrs.created, wrs.id
	`).Bind(dbx.Params{
		"workRecordId": workRecordID,
	}).All(&rows); err != nil {
		return nil, err
	}
	authID := ""
	if auth != nil {
		authID = auth.Id
	}
	creatorID := ""
	if len(creatorIDs) > 0 {
		creatorID = creatorIDs[0]
	}
	for i := range rows {
		rows[i].CanCreateTimeEntry = rows[i].UID == authID && rows[i].LinkedTimeEntryID == ""
		rows[i].CanEdit = !rows[i].Approved && (rows[i].UID == authID || creatorID == authID)
	}
	return rows, nil
}

func loadWorkRecordNote(app core.App, noteID string) (WorkRecordNoteResponse, error) {
	var row WorkRecordNoteResponse
	err := app.DB().NewQuery(`
		SELECT
			wrn.id AS id,
			wrn.uid AS uid,
			COALESCE(p.given_name || ' ' || p.surname, '') AS author_name,
			COALESCE(wrn.note, '') AS note,
			COALESCE(wrn.created, '') AS created
		FROM work_records_notes wrn
		LEFT JOIN profiles p
		  ON p.uid = wrn.uid
		WHERE wrn.id = {:noteId}
		LIMIT 1
	`).Bind(dbx.Params{
		"noteId": noteID,
	}).One(&row)
	if err != nil {
		return WorkRecordNoteResponse{}, err
	}

	return row, nil
}

func loadWorkRecordConsumables(app core.App, workRecordID string, typeID string) ([]WorkRecordConsumableResponse, error) {
	rows := []WorkRecordConsumableResponse{}
	if err := app.DB().NewQuery(`
		SELECT
			COALESCE(wrce.id, '') AS entry_id,
			COALESCE(wrc.id, '') AS consumable_id,
			COALESCE(wrc.name, '') AS consumable_name,
			COALESCE(wrc.input_kind, '') AS input_kind,
			COALESCE(wrc.unit_label, '') AS unit_label,
			COALESCE(wrtc.sort_order, 0) AS sort_order,
			COALESCE(wrce.quantity_number, 0) AS quantity_number,
			CASE WHEN wrce.id != '' THEN 1 ELSE 0 END AS boolean_selected
		FROM work_records_type_consumables wrtc
		JOIN work_records_consumables wrc
		  ON wrc.id = wrtc.consumable
		LEFT JOIN work_records_consumable_entries wrce
		  ON wrce.work_record = {:workRecordId}
		 AND wrce.consumable = wrtc.consumable
		WHERE wrtc.type = {:typeId}
		  AND COALESCE(wrtc.active, 0) = 1
		  AND COALESCE(wrc.active, 0) = 1
		ORDER BY COALESCE(wrtc.sort_order, 0), wrc.name
	`).Bind(dbx.Params{
		"workRecordId": workRecordID,
		"typeId":       typeID,
	}).All(&rows); err != nil {
		return nil, err
	}
	return rows, nil
}

func findWorkRecordByID(app core.App, workRecordID string) (*core.Record, error) {
	workRecordID = strings.TrimSpace(workRecordID)
	if workRecordID == "" {
		return nil, sql.ErrNoRows
	}

	collection, err := findWorkRecordCollection(app)
	if err != nil {
		return nil, err
	}

	row := workRecordSaveStateRow{}
	if err := app.DB().NewQuery(`
		SELECT
			wr.id AS id,
			COALESCE(wr.number, '') AS number,
			COALESCE(wr.date, '') AS date,
			COALESCE(wr.creator, '') AS creator,
			COALESCE(wr.job, '') AS job,
			COALESCE(wr.type, '') AS type,
			COALESCE(wr.location, '') AS location,
			COALESCE(wr.work_description, '') AS work_description,
			COALESCE(wr.report_to, '') AS report_to,
			COALESCE(wr.lane_kms, 0) AS lane_kms,
			COALESCE(wr.field_book_number, '') AS field_book_number,
			COALESCE(wr.field_book_page_number, 0) AS field_book_page_number,
			COALESCE(wr.sub_contractor, '') AS sub_contractor,
			COALESCE(wr.equipment, '') AS equipment,
			COALESCE(wr.supplies, '') AS supplies,
			COALESCE(wr.parent_work_record, '') AS parent_work_record
		FROM work_records wr
		WHERE wr.id = {:id}
		LIMIT 1
	`).Bind(dbx.Params{
		"id": workRecordID,
	}).One(&row); err != nil {
		return nil, err
	}

	record := core.NewRecord(collection)
	record.Load(map[string]any{
		"id":                     row.ID,
		"number":                 row.Number,
		"date":                   row.Date,
		"creator":                row.Creator,
		"job":                    row.Job,
		"type":                   row.Type,
		"location":               row.Location,
		"work_description":       row.WorkDescription,
		"report_to":              row.ReportTo,
		"lane_kms":               row.LaneKms,
		"field_book_number":      row.FieldBookNumber,
		"field_book_page_number": row.FieldBookPageNumber,
		"sub_contractor":         row.SubContractor,
		"equipment":              row.Equipment,
		"supplies":               row.Supplies,
		"parent_work_record":     row.ParentWorkRecord,
	})
	if err := record.PostScan(); err != nil {
		return nil, err
	}

	return record, nil
}

func findWorkRecordCopySourceByID(app core.App, workRecordID string) (workRecordCopySourceRow, error) {
	workRecordID = strings.TrimSpace(workRecordID)
	if workRecordID == "" {
		return workRecordCopySourceRow{}, sql.ErrNoRows
	}

	var row workRecordCopySourceRow
	err := app.DB().NewQuery(`
		SELECT
			wr.id AS id,
			wr.number AS work_record,
			wr.creator AS creator,
			wr.date AS date,
			COALESCE(wr.job, '') AS job,
			COALESCE(wr.type, '') AS type,
			COALESCE(wr.location, '') AS location,
			COALESCE(wr.work_description, '') AS work_description,
			COALESCE(wr.report_to, '') AS report_to,
			COALESCE(wr.lane_kms, 0) AS lane_kms,
			COALESCE(wr.field_book_number, '') AS field_book_number,
			COALESCE(wr.field_book_page_number, 0) AS field_book_page_number,
			COALESCE(wr.sub_contractor, '') AS sub_contractor,
			COALESCE(wr.equipment, '') AS equipment,
			COALESCE(wr.supplies, '') AS supplies,
			COALESCE(wr.parent_work_record, '') AS parent_work_record_id,
			COALESCE(parent_wr.number, '') AS parent_work_record
		FROM work_records wr
		LEFT JOIN work_records parent_wr
		  ON parent_wr.id = wr.parent_work_record
		WHERE wr.id = {:id}
		LIMIT 1
	`).Bind(dbx.Params{
		"id": workRecordID,
	}).One(&row)
	if err != nil {
		return workRecordCopySourceRow{}, err
	}

	return row, nil
}

func toWorkRecordSaveRequest(req WorkRecordCreateRequest) WorkRecordSaveRequest {
	return WorkRecordSaveRequest{
		Date:                req.Date,
		Job:                 req.Job,
		Type:                req.Type,
		Location:            req.Location,
		WorkDescription:     req.WorkDescription,
		ReportTo:            req.ReportTo,
		LaneKms:             req.LaneKms,
		FieldBookNumber:     req.FieldBookNumber,
		FieldBookPageNumber: req.FieldBookPageNumber,
		SubContractor:       req.SubContractor,
		Equipment:           req.Equipment,
		Supplies:            req.Supplies,
		Subjects:            req.Subjects,
		Consumables:         req.Consumables,
	}
}

func canCreateWorkRecordForRoster(app core.App, auth *core.Record, targetUIDs []string) (bool, error) {
	if auth == nil {
		return false, nil
	}

	seen := map[string]struct{}{}
	includesAuth := false
	for _, targetUID := range targetUIDs {
		targetUID = strings.TrimSpace(targetUID)
		if targetUID == "" {
			return false, nil
		}
		if _, ok := seen[targetUID]; ok {
			continue
		}
		seen[targetUID] = struct{}{}

		if targetUID == auth.Id {
			includesAuth = true
		}

		canBeWorker, err := canBeWorkRecordWorker(app, targetUID)
		if err != nil || !canBeWorker {
			return false, err
		}
	}
	if len(seen) == 0 {
		return false, nil
	}

	if includesAuth {
		return utilities.HasClaim(app, auth, "time")
	}

	return utilities.HasClaim(app, auth, "work_record")
}

func canBeWorkRecordWorker(app core.App, targetUID string) (bool, error) {
	targetUID = strings.TrimSpace(targetUID)
	if targetUID == "" {
		return false, nil
	}

	targetUser, err := app.FindRecordById("users", targetUID)
	if err != nil {
		return false, err
	}

	return utilities.HasClaim(app, targetUser, "time")
}

func resolveWorkRecordRootParent(source workRecordCopySourceRow) (string, string) {
	return firstNonEmpty(source.ParentWorkRecordID, source.ID), firstNonEmpty(source.ParentWorkRecord, source.WorkRecord)
}

func requireCopySourceAccess(app core.App, auth *core.Record, source workRecordCopySourceRow) error {
	if auth == nil {
		return workRecordRouteError(http.StatusUnauthorized, "unauthorized", "unauthorized")
	}

	subjects, err := loadWorkRecordSubjects(app, source.ID, auth)
	if err != nil {
		return err
	}
	isSubject := false
	for _, subject := range subjects {
		if subject.UID == auth.Id {
			isSubject = true
			break
		}
	}

	if !isSubject && source.Creator != auth.Id {
		canViewAll, err := canViewAllWorkRecords(app, auth)
		if err != nil {
			return err
		}
		if !canViewAll {
			return workRecordRouteError(http.StatusForbidden, "unauthorized", "you are not authorized to view this work record")
		}
	}

	targetUIDs := make([]string, 0, len(subjects))
	for _, subject := range subjects {
		targetUIDs = append(targetUIDs, subject.UID)
	}
	if err := requireWorkRecordRosterCreateAccess(app, auth, targetUIDs, "copy this work record"); err != nil {
		return err
	}

	return nil
}

func requireWorkRecordRosterCreateAccess(app core.App, auth *core.Record, targetUIDs []string, action string) error {
	canCreate, err := canCreateWorkRecordForRoster(app, auth, targetUIDs)
	if err != nil {
		return err
	}
	if !canCreate {
		return workRecordRouteError(http.StatusForbidden, "unauthorized", "you are not authorized to "+action)
	}
	return nil
}

func requireWorkRecordWorkerEligibility(app core.App, targetUID string, action string) error {
	canBeWorker, err := canBeWorkRecordWorker(app, targetUID)
	if err != nil {
		return err
	}
	if !canBeWorker {
		return workRecordRouteError(http.StatusForbidden, "unauthorized", "you are not authorized to "+action)
	}
	return nil
}

func loadWorkRecordCopyDraft(app core.App, auth *core.Record, workRecordID string) (WorkRecordCopyDraftResponse, error) {
	source, err := findWorkRecordCopySourceByID(app, workRecordID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return WorkRecordCopyDraftResponse{}, workRecordRouteError(http.StatusNotFound, "work_record_not_found", "work record not found")
		}
		return WorkRecordCopyDraftResponse{}, err
	}

	if err := requireCopySourceAccess(app, auth, source); err != nil {
		return WorkRecordCopyDraftResponse{}, err
	}

	sourceDate, err := time.Parse("2006-01-02", source.Date)
	if err != nil {
		return WorkRecordCopyDraftResponse{}, workRecordRouteError(http.StatusBadRequest, "invalid_source_date", fmt.Sprintf("work record has invalid date: %v", err))
	}

	consumables, err := loadWorkRecordConsumables(app, source.ID, source.Type)
	if err != nil {
		return WorkRecordCopyDraftResponse{}, err
	}
	for i := range consumables {
		consumables[i].EntryID = ""
	}
	subjects, err := loadWorkRecordSubjects(app, source.ID, auth)
	if err != nil {
		return WorkRecordCopyDraftResponse{}, err
	}
	for i := range subjects {
		subjects[i].ID = ""
		subjects[i].Approved = false
		subjects[i].LinkedTimeEntryID = ""
		subjects[i].LinkedTimeSheetID = ""
		subjects[i].LinkedTimeEntryDate = ""
		subjects[i].LinkedTimeEntryDescription = ""
		subjects[i].LinkedTimeSheetWeekEnding = ""
	}

	rootWorkRecordID, rootWorkRecord := resolveWorkRecordRootParent(source)

	return WorkRecordCopyDraftResponse{
		Date:                 sourceDate.AddDate(0, 0, 1).Format("2006-01-02"),
		Job:                  source.Job,
		Type:                 source.Type,
		Location:             source.Location,
		WorkDescription:      source.WorkDescription,
		ReportTo:             source.ReportTo,
		LaneKms:              source.LaneKms,
		FieldBookNumber:      source.FieldBookNumber,
		FieldBookPageNumber:  source.FieldBookPageNumber,
		SubContractor:        source.SubContractor,
		Equipment:            source.Equipment,
		Supplies:             source.Supplies,
		Subjects:             subjects,
		Consumables:          consumables,
		CopyFromWorkRecordID: source.ID,
		SourceWorkRecord:     source.WorkRecord,
		RootWorkRecordID:     rootWorkRecordID,
		RootWorkRecord:       rootWorkRecord,
	}, nil
}

func workRecordRouteError(status int, code string, message string) *errs.HookError {
	return &errs.HookError{
		Status:  status,
		Message: message,
		Data: map[string]errs.CodeError{
			"global": {
				Code:    code,
				Message: message,
			},
		},
	}
}

func workRecordFieldError(field string, code string, message string) *errs.HookError {
	return &errs.HookError{
		Status:  http.StatusBadRequest,
		Message: "validation failed",
		Data: map[string]errs.CodeError{
			field: {
				Code:    code,
				Message: message,
			},
		},
	}
}

func findWorkRecordRecordByID(app core.App, workRecordID string) (*core.Record, error) {
	workRecordID = strings.TrimSpace(workRecordID)
	if workRecordID == "" {
		return nil, sql.ErrNoRows
	}

	return app.FindRecordById(workRecordsCollectionIdentifier, workRecordID)
}

func requireWorkRecordEditAccess(app core.App, auth *core.Record, record *core.Record) error {
	if auth == nil {
		return workRecordRouteError(http.StatusUnauthorized, "unauthorized", "unauthorized")
	}

	allowed := record.GetString("creator") == auth.Id
	if !allowed {
		var row struct {
			Count int `db:"count"`
		}
		if err := app.DB().NewQuery(`
			SELECT COUNT(*) AS count
			FROM work_records_subjects
			WHERE work_record = {:workRecordId}
			  AND uid = {:uid}
		`).Bind(dbx.Params{
			"workRecordId": record.Id,
			"uid":          auth.Id,
		}).One(&row); err != nil {
			return err
		}
		allowed = row.Count > 0
	}
	if !allowed {
		return workRecordRouteError(http.StatusForbidden, "unauthorized", "you are not authorized to edit this work record")
	}

	return nil
}

func workRecordHasApprovedSubjectForRoute(app core.App, workRecordID string) (bool, error) {
	var lockRow struct {
		Approved bool `db:"approved"`
	}
	if err := app.DB().NewQuery(`
		SELECT EXISTS(
			SELECT 1
			FROM work_records_subjects
			WHERE work_record = {:workRecordId}
			  AND COALESCE(approved, 0) = 1
		) AS approved
	`).Bind(dbx.Params{
		"workRecordId": workRecordID,
	}).One(&lockRow); err != nil {
		return false, err
	}
	return lockRow.Approved, nil
}

func requireUnlockedSharedWorkRecordEditAccess(app core.App, auth *core.Record, record *core.Record) error {
	if err := requireWorkRecordEditAccess(app, auth, record); err != nil {
		return err
	}

	locked, err := workRecordHasApprovedSubjectForRoute(app, record.Id)
	if err != nil {
		return err
	}
	if locked {
		return workRecordRouteError(http.StatusBadRequest, "approved_work_record_locked", "Approved work records cannot be edited")
	}

	return nil
}

func loadWorkRecordConsumableEntryStates(app core.App, workRecordID string) ([]workRecordConsumableEntryState, error) {
	rows := []workRecordConsumableEntryState{}
	if err := app.DB().NewQuery(`
		SELECT
			id,
			COALESCE(consumable, '') AS consumable_id
		FROM work_records_consumable_entries
		WHERE work_record = {:workRecordId}
	`).Bind(dbx.Params{
		"workRecordId": workRecordID,
	}).All(&rows); err != nil {
		return nil, err
	}
	return rows, nil
}

func canAddWorkRecordNote(app core.App, auth *core.Record, workRecordID string, creatorUID string) bool {
	if auth == nil {
		return false
	}

	if creatorUID == auth.Id {
		return true
	}
	var row struct {
		Count int `db:"count"`
	}
	if err := app.DB().NewQuery(`
		SELECT COUNT(*) AS count
		FROM work_records_subjects
		WHERE work_record = {:workRecordId}
		  AND uid = {:uid}
	`).Bind(dbx.Params{
		"workRecordId": workRecordID,
		"uid":          auth.Id,
	}).One(&row); err != nil {
		return false
	}
	return row.Count > 0
}

func validateWorkRecordNoteCreateRequest(req WorkRecordNoteCreateRequest) error {
	return validation.Errors{
		"note": validation.Validate(
			strings.TrimSpace(req.Note),
			validation.Required.Error("Note is required"),
		),
	}.Filter()
}

func normalizeOptionalFieldBookPageNumber(raw any) (any, int, error) {
	switch value := raw.(type) {
	case nil:
		return "", 0, nil
	case string:
		if strings.TrimSpace(value) == "" {
			return "", 0, nil
		}
	case float64:
		return value, int(value), nil
	case float32:
		return value, int(value), nil
	case int:
		return value, value, nil
	case int32:
		return value, int(value), nil
	case int64:
		return value, int(value), nil
	}
	return nil, 0, workRecordFieldError("field_book_page_number", "invalid_field_value", "Field book page number must be a number or empty")
}

func applyWorkRecordSaveRequest(record *core.Record, req WorkRecordSaveRequest) error {
	record.Set("date", req.Date)
	record.Set("job", req.Job)
	record.Set("type", req.Type)
	record.Set("location", req.Location)
	record.Set("work_description", req.WorkDescription)
	record.Set("report_to", req.ReportTo)
	record.Set("lane_kms", req.LaneKms)
	record.Set("field_book_number", req.FieldBookNumber)

	pageNumberValue, _, err := normalizeOptionalFieldBookPageNumber(req.FieldBookPageNumber)
	if err != nil {
		return err
	}
	record.Set("field_book_page_number", pageNumberValue)

	record.Set("sub_contractor", req.SubContractor)
	record.Set("equipment", req.Equipment)
	record.Set("supplies", req.Supplies)

	return nil
}

func subjectsFromSaveRequest(req WorkRecordSaveRequest) ([]WorkRecordSubjectSaveRequest, error) {
	if req.Subjects != nil {
		return req.Subjects, nil
	}
	return nil, workRecordFieldError("subjects", "required", "At least one worker is required")
}

func validateSubjectPayload(subjects []WorkRecordSubjectSaveRequest) error {
	if len(subjects) == 0 {
		return &errs.HookError{
			Status:  http.StatusBadRequest,
			Message: "validation failed",
			Data: map[string]errs.CodeError{
				"subjects": {Code: "required", Message: "At least one worker is required"},
			},
		}
	}
	seen := map[string]struct{}{}
	for _, subject := range subjects {
		uid := strings.TrimSpace(subject.UID)
		if uid == "" {
			return workRecordFieldError("subjects", "subject_required", "Every worker must include a user")
		}
		if _, exists := seen[uid]; exists {
			return workRecordFieldError("subjects", "duplicate_subject", "Each worker may appear only once")
		}
		seen[uid] = struct{}{}
	}
	return nil
}

func applyWorkRecordSubjectSaveRequest(record *core.Record, workRecordID string, req WorkRecordSubjectSaveRequest) {
	record.Set("work_record", workRecordID)
	record.Set("uid", strings.TrimSpace(req.UID))
	record.Set("hours_on_site", req.HoursOnSite)
	record.Set("hours_travel_time", req.HoursTravelTime)
	record.Set("distance_travelled_km", req.DistanceTravelledKM)
	record.Set("vehicle_type", req.VehicleType)
	record.Set("is_passenger", req.IsPassenger)
	record.Set("company_vehicle_unit_number", req.CompanyVehicleUnitNumber)
}

func workRecordSubjectErrorField(field string) (string, bool) {
	switch field {
	case "uid":
		return "uid", true
	case "hours_on_site", "hours_on_site_min":
		return "hours_on_site", true
	case "hours_travel_time", "travel_time_min":
		return "hours_travel_time", true
	case "distance_travelled_km", "distance_km_min":
		return "distance_travelled_km", true
	case "vehicle_type":
		return "vehicle_type", true
	case "company_vehicle_unit_number":
		return "company_vehicle_unit_number", true
	default:
		return "", false
	}
}

func scopeWorkRecordSubjectValidationError(err error, index int) error {
	var hookErr *errs.HookError
	if !errors.As(err, &hookErr) || len(hookErr.Data) == 0 {
		return err
	}

	scoped := &errs.HookError{
		Status:  hookErr.Status,
		Message: hookErr.Message,
		Data:    make(map[string]errs.CodeError, len(hookErr.Data)+1),
	}
	for field, codeErr := range hookErr.Data {
		scoped.Data[field] = codeErr
		if subjectField, ok := workRecordSubjectErrorField(field); ok {
			scoped.Data[fmt.Sprintf("subjects_%d_%s", index, subjectField)] = codeErr
		}
	}
	return scoped
}

func loadExistingSubjectRecords(app core.App, workRecordID string) (map[string]*core.Record, map[string]*core.Record, error) {
	subjects, err := loadWorkRecordSubjects(app, workRecordID, nil)
	if err != nil {
		return nil, nil, err
	}
	byID := map[string]*core.Record{}
	byUID := map[string]*core.Record{}
	for _, subject := range subjects {
		record, err := findWorkRecordSubjectByID(app, subject.ID)
		if err != nil {
			return nil, nil, err
		}
		record.Set("work_record", workRecordID)
		byID[subject.ID] = record
		byUID[subject.UID] = record
	}
	return byID, byUID, nil
}

func workRecordSubjectRequestMatchesRecord(record *core.Record, req WorkRecordSubjectSaveRequest) bool {
	return record.GetString("uid") == strings.TrimSpace(req.UID) &&
		record.GetFloat("hours_on_site") == req.HoursOnSite &&
		record.GetFloat("hours_travel_time") == req.HoursTravelTime &&
		record.GetFloat("distance_travelled_km") == req.DistanceTravelledKM &&
		record.GetString("vehicle_type") == req.VehicleType &&
		record.GetBool("is_passenger") == req.IsPassenger &&
		record.GetString("company_vehicle_unit_number") == req.CompanyVehicleUnitNumber
}

func saveWorkRecordSubjects(txApp core.App, auth *core.Record, workRecordID string, creatorID string, subjects []WorkRecordSubjectSaveRequest, allowRosterChanges bool, rosterLocked bool) error {
	if err := validateSubjectPayload(subjects); err != nil {
		return err
	}

	collection, err := findWorkRecordSubjectsCollection(txApp)
	if err != nil {
		return err
	}
	existingByID, existingByUID, err := loadExistingSubjectRecords(txApp, workRecordID)
	if err != nil {
		return err
	}

	seenExistingIDs := map[string]struct{}{}
	if allowRosterChanges && auth != nil && auth.Id == creatorID {
		for _, subject := range subjects {
			subject.ID = strings.TrimSpace(subject.ID)
			subject.UID = strings.TrimSpace(subject.UID)
			if subject.ID != "" {
				if existingByID[subject.ID] != nil {
					seenExistingIDs[subject.ID] = struct{}{}
				}
				continue
			}
			if existing := existingByUID[subject.UID]; existing != nil {
				seenExistingIDs[existing.Id] = struct{}{}
			}
		}
		for id, record := range existingByID {
			if _, seen := seenExistingIDs[id]; seen {
				continue
			}
			if err := normalizeRouteValidationError(apphooks.ValidateWorkRecordSubjectDelete(txApp, record)); err != nil {
				return err
			}
			if err := normalizeRouteValidationError(txApp.Delete(record)); err != nil {
				return err
			}
		}
	}

	seenExistingIDs = map[string]struct{}{}
	for index, subject := range subjects {
		subject.ID = strings.TrimSpace(subject.ID)
		subject.UID = strings.TrimSpace(subject.UID)

		var record *core.Record
		if subject.ID != "" {
			existing, ok := existingByID[subject.ID]
			if !ok {
				return workRecordFieldError("subjects", "invalid_subject", "Worker was not found for this work record")
			}
			record = existing
		} else if existing := existingByUID[subject.UID]; existing != nil {
			record = existing
		} else {
			if rosterLocked {
				return workRecordFieldError("subjects", "approved_roster_locked", "Workers cannot be added or reassigned after approval")
			}
			if !allowRosterChanges || auth == nil || auth.Id != creatorID {
				return workRecordFieldError("subjects", "creator_required", "Only the work record creator can add workers")
			}
			if err := requireWorkRecordWorkerEligibility(txApp, subject.UID, "add this worker"); err != nil {
				return err
			}
			record = core.NewRecord(collection)
		}

		if auth == nil {
			return workRecordRouteError(http.StatusUnauthorized, "unauthorized", "unauthorized")
		}
		if auth.Id != creatorID && record.GetString("uid") != auth.Id {
			return workRecordFieldError("subjects", "unauthorized", "You may only edit your own worker row")
		}

		if record.GetBool("approved") && workRecordSubjectRequestMatchesRecord(record, subject) {
			seenExistingIDs[record.Id] = struct{}{}
			continue
		}

		applyWorkRecordSubjectSaveRequest(record, workRecordID, subject)
		if err := normalizeRouteValidationError(apphooks.ValidateWorkRecordSubject(txApp, record)); err != nil {
			return scopeWorkRecordSubjectValidationError(err, index)
		}
		if err := normalizeRouteValidationError(txApp.Save(record)); err != nil {
			return scopeWorkRecordSubjectValidationError(err, index)
		}
		seenExistingIDs[record.Id] = struct{}{}
	}

	return nil
}

func validateOptionalNumberMinimum(collection *core.Collection, fieldName string, raw any) error {
	numberField, ok := collection.Fields.GetByName(fieldName).(*core.NumberField)
	if !ok || numberField == nil || numberField.Required || numberField.Min == nil || *numberField.Min <= 0 {
		return nil
	}

	switch value := raw.(type) {
	case nil:
		return nil
	case string:
		if strings.TrimSpace(value) == "" {
			return nil
		}
		return nil
	case float64:
		if value != 0 {
			return nil
		}
	case float32:
		if value != 0 {
			return nil
		}
	case int:
		if value != 0 {
			return nil
		}
	case int32:
		if value != 0 {
			return nil
		}
	case int64:
		if value != 0 {
			return nil
		}
	default:
		return nil
	}

	// PocketBase treats 0 as "empty" for optional number fields, so a
	// migrated min constraint like field_book_page_number >= 1 won't fire on
	// the custom save route unless we patch that gap explicitly.
	return &errs.HookError{
		Status:  http.StatusBadRequest,
		Message: "validation failed",
		Data: map[string]errs.CodeError{
			fieldName: {
				Code:    "validation_min_number_constraint",
				Message: fmt.Sprintf("Must be larger than %f", *numberField.Min),
			},
		},
	}
}

func hasRequestedWorkRecordConsumableValue(item WorkRecordSaveConsumableRequest) bool {
	if item.InputKind == "boolean" {
		return item.BooleanSelected
	}
	return item.QuantityNumber > 0
}

func consumableFieldError(code string, message string) *errs.HookError {
	return &errs.HookError{
		Status:  http.StatusBadRequest,
		Message: "validation failed",
		Data: map[string]errs.CodeError{
			"consumables": {
				Code:    code,
				Message: message,
			},
		},
	}
}

func loadWorkRecordConsumableInputKinds(app core.App, consumableIDs []string) (map[string]string, error) {
	kinds := map[string]string{}
	if len(consumableIDs) == 0 {
		return kinds, nil
	}

	placeholders := make([]string, 0, len(consumableIDs))
	params := dbx.Params{}
	for i, consumableID := range consumableIDs {
		key := fmt.Sprintf("consumableId%d", i)
		placeholders = append(placeholders, "{:"+key+"}")
		params[key] = consumableID
	}

	rows := []struct {
		ID        string `db:"id"`
		InputKind string `db:"input_kind"`
	}{}
	if err := app.DB().NewQuery(`
		SELECT
			id,
			COALESCE(input_kind, '') AS input_kind
		FROM work_records_consumables
		WHERE id IN (` + strings.Join(placeholders, ", ") + `)
	`).Bind(params).All(&rows); err != nil {
		return nil, err
	}
	for _, row := range rows {
		kinds[row.ID] = row.InputKind
	}
	return kinds, nil
}

func resolveWorkRecordConsumableInputKinds(app core.App, items []WorkRecordSaveConsumableRequest) (map[string]string, error) {
	uniqueIDs := make([]string, 0, len(items))
	seen := map[string]struct{}{}
	for _, item := range items {
		consumableID := strings.TrimSpace(item.ConsumableID)
		if consumableID == "" {
			continue
		}
		if _, exists := seen[consumableID]; exists {
			continue
		}
		seen[consumableID] = struct{}{}
		uniqueIDs = append(uniqueIDs, consumableID)
	}

	kinds, err := loadWorkRecordConsumableInputKinds(app, uniqueIDs)
	if err != nil {
		return nil, err
	}
	for _, consumableID := range uniqueIDs {
		if strings.TrimSpace(kinds[consumableID]) == "" {
			return nil, consumableFieldError("invalid_consumable_entry", "Consumable was not found")
		}
	}
	return kinds, nil
}

func hasResolvedWorkRecordConsumableValue(inputKind string, item WorkRecordSaveConsumableRequest) bool {
	switch inputKind {
	case "boolean":
		return item.BooleanSelected
	case "number":
		return item.QuantityNumber > 0
	default:
		return hasRequestedWorkRecordConsumableValue(item)
	}
}

func setWorkRecordConsumableQuantity(entryRecord *core.Record, inputKind string, item WorkRecordSaveConsumableRequest) {
	if inputKind == "boolean" {
		entryRecord.Set("quantity_number", 0)
		return
	}
	entryRecord.Set("quantity_number", item.QuantityNumber)
}

func lockedWorkRecordEditError() *errs.HookError {
	return workRecordFieldError("global", "approved_work_record_locked", "Approved work records cannot be edited")
}

func validateLockedWorkRecordSharedPayload(record *core.Record, req WorkRecordSaveRequest) error {
	_, pageNumber, err := normalizeOptionalFieldBookPageNumber(req.FieldBookPageNumber)
	if err != nil {
		return err
	}

	if record.GetString("date") != req.Date ||
		record.GetString("job") != req.Job ||
		record.GetString("type") != req.Type ||
		record.GetString("location") != req.Location ||
		record.GetString("work_description") != req.WorkDescription ||
		record.GetString("report_to") != req.ReportTo ||
		record.GetFloat("lane_kms") != req.LaneKms ||
		record.GetString("field_book_number") != req.FieldBookNumber ||
		record.GetInt("field_book_page_number") != pageNumber ||
		record.GetString("sub_contractor") != req.SubContractor ||
		record.GetString("equipment") != req.Equipment ||
		record.GetString("supplies") != req.Supplies {
		return lockedWorkRecordEditError()
	}

	return nil
}

func requestedConsumableValueMap(app core.App, items []WorkRecordSaveConsumableRequest) (map[string]float64, error) {
	values := map[string]float64{}
	inputKindsByConsumable, err := resolveWorkRecordConsumableInputKinds(app, items)
	if err != nil {
		return nil, err
	}

	seenConsumables := map[string]struct{}{}
	for _, item := range items {
		consumableID := strings.TrimSpace(item.ConsumableID)
		if consumableID == "" {
			return nil, consumableFieldError("consumable_required", "Consumable id is required for all submitted consumables")
		}
		if _, exists := seenConsumables[consumableID]; exists {
			return nil, consumableFieldError("duplicate_consumable_payload", "Each consumable may be submitted only once")
		}
		seenConsumables[consumableID] = struct{}{}

		inputKind := inputKindsByConsumable[consumableID]
		if !hasResolvedWorkRecordConsumableValue(inputKind, item) {
			continue
		}
		if inputKind == "boolean" {
			values[consumableID] = 0
		} else {
			values[consumableID] = item.QuantityNumber
		}
	}

	return values, nil
}

func storedConsumableValueMap(app core.App, workRecordID string) (map[string]float64, error) {
	rows := []struct {
		ConsumableID   string  `db:"consumable"`
		QuantityNumber float64 `db:"quantity_number"`
	}{}
	if err := app.DB().NewQuery(`
		SELECT
			COALESCE(consumable, '') AS consumable,
			COALESCE(quantity_number, 0) AS quantity_number
		FROM work_records_consumable_entries
		WHERE work_record = {:workRecordId}
	`).Bind(dbx.Params{
		"workRecordId": workRecordID,
	}).All(&rows); err != nil {
		return nil, err
	}

	values := map[string]float64{}
	for _, row := range rows {
		values[row.ConsumableID] = row.QuantityNumber
	}
	return values, nil
}

func validateLockedWorkRecordConsumablePayload(app core.App, workRecordID string, items []WorkRecordSaveConsumableRequest) error {
	requested, err := requestedConsumableValueMap(app, items)
	if err != nil {
		return err
	}
	stored, err := storedConsumableValueMap(app, workRecordID)
	if err != nil {
		return err
	}
	if len(requested) != len(stored) {
		return lockedWorkRecordEditError()
	}
	for consumableID, quantity := range stored {
		if requested[consumableID] != quantity {
			return lockedWorkRecordEditError()
		}
	}
	return nil
}

func saveWorkRecordEdit(txApp core.App, auth *core.Record, record *core.Record, req WorkRecordSaveRequest) (WorkRecordSaveResponse, error) {
	currentRecord, err := findWorkRecordByID(txApp, record.Id)
	if err != nil {
		return WorkRecordSaveResponse{}, err
	}

	if err := requireWorkRecordEditAccess(txApp, auth, currentRecord); err != nil {
		return WorkRecordSaveResponse{}, err
	}

	locked, err := workRecordHasApprovedSubjectForRoute(txApp, currentRecord.Id)
	if err != nil {
		return WorkRecordSaveResponse{}, err
	}
	if locked {
		if err := validateLockedWorkRecordSharedPayload(currentRecord, req); err != nil {
			return WorkRecordSaveResponse{}, err
		}
		if err := validateLockedWorkRecordConsumablePayload(txApp, currentRecord.Id, req.Consumables); err != nil {
			return WorkRecordSaveResponse{}, err
		}
		subjects, err := subjectsFromSaveRequest(req)
		if err != nil {
			return WorkRecordSaveResponse{}, err
		}
		if err := saveWorkRecordSubjects(txApp, auth, currentRecord.Id, currentRecord.GetString("creator"), subjects, false, true); err != nil {
			return WorkRecordSaveResponse{}, err
		}
		return WorkRecordSaveResponse{ID: currentRecord.Id}, nil
	}

	subjects, err := subjectsFromSaveRequest(req)
	if err != nil {
		return WorkRecordSaveResponse{}, err
	}

	if auth == nil || auth.Id != currentRecord.GetString("creator") {
		if err := saveWorkRecordSubjects(txApp, auth, currentRecord.Id, currentRecord.GetString("creator"), subjects, false, false); err != nil {
			return WorkRecordSaveResponse{}, err
		}
		return WorkRecordSaveResponse{ID: currentRecord.Id}, nil
	}

	existingEntries, err := loadWorkRecordConsumableEntryStates(txApp, currentRecord.Id)
	if err != nil {
		return WorkRecordSaveResponse{}, err
	}

	existingByID := make(map[string]workRecordConsumableEntryState, len(existingEntries))
	existingByConsumable := make(map[string]workRecordConsumableEntryState, len(existingEntries))
	for _, row := range existingEntries {
		existingByID[row.ID] = row
		existingByConsumable[row.ConsumableID] = row
	}

	inputKindsByConsumable, err := resolveWorkRecordConsumableInputKinds(txApp, req.Consumables)
	if err != nil {
		return WorkRecordSaveResponse{}, err
	}

	seenEntryIDs := map[string]struct{}{}
	seenConsumables := map[string]struct{}{}
	deleteEntryIDs := []string{}
	upserts := []WorkRecordSaveConsumableRequest{}

	for _, item := range req.Consumables {
		item.EntryID = strings.TrimSpace(item.EntryID)
		item.ConsumableID = strings.TrimSpace(item.ConsumableID)

		if item.ConsumableID == "" {
			return WorkRecordSaveResponse{}, consumableFieldError("consumable_required", "Consumable id is required for all submitted consumables")
		}

		if _, exists := seenConsumables[item.ConsumableID]; exists {
			return WorkRecordSaveResponse{}, consumableFieldError("duplicate_consumable_payload", "Each consumable may be submitted only once")
		}
		seenConsumables[item.ConsumableID] = struct{}{}

		existingEntry, foundExistingEntry := existingByConsumable[item.ConsumableID]
		if item.EntryID != "" {
			if _, exists := seenEntryIDs[item.EntryID]; exists {
				return WorkRecordSaveResponse{}, consumableFieldError("duplicate_consumable_payload", "Each consumable entry may be submitted only once")
			}
			seenEntryIDs[item.EntryID] = struct{}{}

			entryByID, ok := existingByID[item.EntryID]
			if !ok {
				return WorkRecordSaveResponse{}, consumableFieldError("invalid_consumable_entry", "Consumable entry was not found for this work record")
			}
			if entryByID.ConsumableID != item.ConsumableID {
				return WorkRecordSaveResponse{}, consumableFieldError("immutable_consumable_entry", "Consumable entry ids must match their original consumable")
			}
			existingEntry = entryByID
			foundExistingEntry = true
		}

		inputKind := inputKindsByConsumable[item.ConsumableID]
		if !hasResolvedWorkRecordConsumableValue(inputKind, item) {
			if foundExistingEntry {
				deleteEntryIDs = append(deleteEntryIDs, existingEntry.ID)
			}
			continue
		}

		if foundExistingEntry {
			item.EntryID = existingEntry.ID
		}
		upserts = append(upserts, item)
	}

	for _, entryID := range deleteEntryIDs {
		entryRecord, err := findWorkRecordConsumableEntryByID(txApp, entryID)
		if err != nil {
			return WorkRecordSaveResponse{}, err
		}
		if err := normalizeRouteValidationError(apphooks.ValidateWorkRecordConsumableEntryDelete(txApp, entryRecord)); err != nil {
			return WorkRecordSaveResponse{}, err
		}
		if err := normalizeRouteValidationError(txApp.Delete(entryRecord)); err != nil {
			return WorkRecordSaveResponse{}, err
		}
	}

	// When the type changes, the UI only sends consumables valid for the new
	// type. Delete any existing DB entries not covered by the request so the
	// hook-level compatibility check passes.
	if currentRecord.GetString("type") != req.Type {
		for _, existing := range existingEntries {
			if _, covered := seenConsumables[existing.ConsumableID]; covered {
				continue
			}
			entryRecord, err := findWorkRecordConsumableEntryByID(txApp, existing.ID)
			if err != nil {
				return WorkRecordSaveResponse{}, err
			}
			if err := normalizeRouteValidationError(apphooks.ValidateWorkRecordConsumableEntryDelete(txApp, entryRecord)); err != nil {
				return WorkRecordSaveResponse{}, err
			}
			if err := normalizeRouteValidationError(txApp.Delete(entryRecord)); err != nil {
				return WorkRecordSaveResponse{}, err
			}
		}
	}

	if err := applyWorkRecordSaveRequest(currentRecord, req); err != nil {
		return WorkRecordSaveResponse{}, err
	}
	if err := validateOptionalNumberMinimum(currentRecord.Collection(), "field_book_page_number", req.FieldBookPageNumber); err != nil {
		return WorkRecordSaveResponse{}, err
	}
	if err := normalizeRouteValidationError(apphooks.ValidateWorkRecord(txApp, currentRecord)); err != nil {
		return WorkRecordSaveResponse{}, err
	}
	if err := normalizeRouteValidationError(txApp.Save(currentRecord)); err != nil {
		return WorkRecordSaveResponse{}, err
	}
	if err := saveWorkRecordSubjects(txApp, auth, currentRecord.Id, currentRecord.GetString("creator"), subjects, auth != nil && auth.Id == currentRecord.GetString("creator"), false); err != nil {
		return WorkRecordSaveResponse{}, err
	}

	for _, item := range upserts {
		var entryRecord *core.Record
		if item.EntryID != "" {
			entryRecord, err = findWorkRecordConsumableEntryByID(txApp, item.EntryID)
			if err != nil {
				return WorkRecordSaveResponse{}, err
			}
		} else {
			entriesCollection, err := findWorkRecordConsumableEntriesCollection(txApp)
			if err != nil {
				return WorkRecordSaveResponse{}, err
			}
			entryRecord = core.NewRecord(entriesCollection)
		}

		entryRecord.Set("work_record", currentRecord.Id)
		entryRecord.Set("consumable", item.ConsumableID)
		setWorkRecordConsumableQuantity(entryRecord, inputKindsByConsumable[item.ConsumableID], item)

		if err := normalizeRouteValidationError(apphooks.ValidateWorkRecordConsumableEntry(txApp, entryRecord)); err != nil {
			return WorkRecordSaveResponse{}, err
		}
		if err := normalizeRouteValidationError(txApp.Save(entryRecord)); err != nil {
			return WorkRecordSaveResponse{}, err
		}
	}

	return WorkRecordSaveResponse{ID: currentRecord.Id}, nil
}

func createWorkRecord(txApp core.App, auth *core.Record, req WorkRecordCreateRequest) (WorkRecordSaveResponse, error) {
	subjects := req.Subjects
	if err := validateSubjectPayload(subjects); err != nil {
		return WorkRecordSaveResponse{}, err
	}

	rootParentID := ""
	copyFromWorkRecordID := strings.TrimSpace(req.CopyFromWorkRecordID)
	if copyFromWorkRecordID != "" {
		source, err := findWorkRecordCopySourceByID(txApp, copyFromWorkRecordID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return WorkRecordSaveResponse{}, workRecordRouteError(http.StatusNotFound, "work_record_not_found", "work record not found")
			}
			return WorkRecordSaveResponse{}, err
		}
		if err := requireCopySourceAccess(txApp, auth, source); err != nil {
			return WorkRecordSaveResponse{}, err
		}
		if strings.TrimSpace(req.Job) != source.Job {
			return WorkRecordSaveResponse{}, &errs.HookError{
				Status:  http.StatusBadRequest,
				Message: "validation failed",
				Data: map[string]errs.CodeError{
					"job": {
						Code:    "copy_job_mismatch",
						Message: "Copied work records must keep the same job as the source record",
					},
				},
			}
		}

		rootParentID, _ = resolveWorkRecordRootParent(source)
	}
	targetUIDs := make([]string, 0, len(subjects))
	for _, subject := range subjects {
		targetUIDs = append(targetUIDs, subject.UID)
	}
	if err := requireWorkRecordRosterCreateAccess(txApp, auth, targetUIDs, "create this work record"); err != nil {
		return WorkRecordSaveResponse{}, err
	}

	collection, err := findWorkRecordCollection(txApp)
	if err != nil {
		return WorkRecordSaveResponse{}, err
	}

	record := core.NewRecord(collection)
	record.Set("creator", auth.Id)
	record.Set("parent_work_record", rootParentID)

	saveReq := toWorkRecordSaveRequest(req)
	if err := applyWorkRecordSaveRequest(record, saveReq); err != nil {
		return WorkRecordSaveResponse{}, err
	}
	if err := validateOptionalNumberMinimum(collection, "field_book_page_number", saveReq.FieldBookPageNumber); err != nil {
		return WorkRecordSaveResponse{}, err
	}
	if err := normalizeRouteValidationError(apphooks.ValidateWorkRecord(txApp, record)); err != nil {
		return WorkRecordSaveResponse{}, err
	}
	if err := normalizeRouteValidationError(apphooks.EnsureWorkRecordNumber(txApp, record)); err != nil {
		return WorkRecordSaveResponse{}, err
	}
	if err := normalizeRouteValidationError(txApp.Save(record)); err != nil {
		return WorkRecordSaveResponse{}, err
	}
	if err := saveWorkRecordSubjects(txApp, auth, record.Id, auth.Id, subjects, true, false); err != nil {
		return WorkRecordSaveResponse{}, err
	}

	inputKindsByConsumable, err := resolveWorkRecordConsumableInputKinds(txApp, req.Consumables)
	if err != nil {
		return WorkRecordSaveResponse{}, err
	}

	seenConsumables := map[string]struct{}{}
	for _, item := range req.Consumables {
		item.EntryID = strings.TrimSpace(item.EntryID)
		item.ConsumableID = strings.TrimSpace(item.ConsumableID)

		if item.EntryID != "" {
			return WorkRecordSaveResponse{}, consumableFieldError("invalid_consumable_entry", "New work records cannot reference existing consumable entries")
		}
		if item.ConsumableID == "" {
			return WorkRecordSaveResponse{}, consumableFieldError("consumable_required", "Consumable id is required for all submitted consumables")
		}
		if _, exists := seenConsumables[item.ConsumableID]; exists {
			return WorkRecordSaveResponse{}, consumableFieldError("duplicate_consumable_payload", "Each consumable may be submitted only once")
		}
		seenConsumables[item.ConsumableID] = struct{}{}

		if !hasResolvedWorkRecordConsumableValue(inputKindsByConsumable[item.ConsumableID], item) {
			continue
		}

		entriesCollection, err := findWorkRecordConsumableEntriesCollection(txApp)
		if err != nil {
			return WorkRecordSaveResponse{}, err
		}

		entryRecord := core.NewRecord(entriesCollection)
		entryRecord.Set("work_record", record.Id)
		entryRecord.Set("consumable", item.ConsumableID)
		setWorkRecordConsumableQuantity(entryRecord, inputKindsByConsumable[item.ConsumableID], item)

		if err := normalizeRouteValidationError(apphooks.ValidateWorkRecordConsumableEntry(txApp, entryRecord)); err != nil {
			return WorkRecordSaveResponse{}, err
		}
		if err := normalizeRouteValidationError(txApp.Save(entryRecord)); err != nil {
			return WorkRecordSaveResponse{}, err
		}
	}

	return WorkRecordSaveResponse{ID: record.Id}, nil
}

func saveWorkRecordAttachment(
	txApp core.App,
	auth *core.Record,
	workRecordID string,
	uploadedFiles []*filesystem.File,
	clearAttachment bool,
) (WorkRecordAttachmentResponse, error) {
	// Attachment writes intentionally use a dedicated route instead of enabling a
	// collection-level updateRule on work_records.
	//
	// That narrower route keeps all regular field edits on the existing save
	// endpoint, preserves the time-editing gate and worker/creator checks in one
	// place, and avoids exposing direct PATCH /api/collections/work_records/... as
	// a second general-purpose mutation path.
	//
	// The flow here is:
	//   1. load the real work-record row;
	//   2. enforce the same edit permissions used by normal work-record saves;
	//   3. apply the attachment add/replace/remove mutation only;
	//   4. run standard validation plus attachment-specific hash/dedupe syncing;
	//   5. persist and return a minimal response for the UI.
	record, err := findWorkRecordRecordByID(txApp, workRecordID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return WorkRecordAttachmentResponse{}, workRecordRouteError(http.StatusNotFound, "work_record_not_found", "work record not found")
		}
		return WorkRecordAttachmentResponse{}, err
	}

	if err := requireUnlockedSharedWorkRecordEditAccess(txApp, auth, record); err != nil {
		return WorkRecordAttachmentResponse{}, err
	}

	if clearAttachment {
		record.Set("attachment", "")
	} else {
		record.Set("attachment", uploadedFiles)
	}

	if err := normalizeRouteValidationError(apphooks.ValidateWorkRecord(txApp, record)); err != nil {
		return WorkRecordAttachmentResponse{}, err
	}
	if err := normalizeRouteValidationError(apphooks.SyncWorkRecordAttachmentFields(txApp, record)); err != nil {
		return WorkRecordAttachmentResponse{}, err
	}
	if err := normalizeRouteValidationError(txApp.Save(record)); err != nil {
		return WorkRecordAttachmentResponse{}, err
	}

	return WorkRecordAttachmentResponse{
		ID:         record.Id,
		Attachment: record.GetString("attachment"),
	}, nil
}

func createCreateWorkRecordHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		requestApp := e.App
		if err := requireTimeEditing(requestApp); err != nil {
			return err
		}

		auth := e.Auth
		if auth == nil {
			return e.Error(http.StatusUnauthorized, "unauthorized", nil)
		}

		var req WorkRecordCreateRequest
		if err := e.BindBody(&req); err != nil {
			return writeHookError(e, workRecordRouteError(http.StatusBadRequest, "invalid_request_body", "invalid request body"))
		}

		response := WorkRecordSaveResponse{}
		err := requestApp.RunInTransaction(func(txApp core.App) error {
			var createErr error
			response, createErr = createWorkRecord(txApp, auth, req)
			return createErr
		})
		if err != nil {
			return writeHookError(e, err)
		}

		return e.JSON(http.StatusOK, response)
	}
}

func createSaveWorkRecordHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		requestApp := e.App
		if err := requireTimeEditing(requestApp); err != nil {
			return err
		}

		auth := e.Auth
		if auth == nil {
			return e.Error(http.StatusUnauthorized, "unauthorized", nil)
		}

		workRecordID := strings.TrimSpace(e.Request.PathValue("workRecord"))
		if workRecordID == "" {
			return writeHookError(e, workRecordRouteError(http.StatusBadRequest, "work_record_required", "workRecord is required"))
		}

		var req WorkRecordSaveRequest
		if err := e.BindBody(&req); err != nil {
			return writeHookError(e, workRecordRouteError(http.StatusBadRequest, "invalid_request_body", "invalid request body"))
		}

		record, err := findWorkRecordByID(requestApp, workRecordID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return writeHookError(e, workRecordRouteError(http.StatusNotFound, "work_record_not_found", "work record not found"))
			}
			return e.Error(http.StatusInternalServerError, "failed to load work record", err)
		}

		response := WorkRecordSaveResponse{}
		err = requestApp.RunInTransaction(func(txApp core.App) error {
			var saveErr error
			response, saveErr = saveWorkRecordEdit(txApp, auth, record, req)
			return saveErr
		})
		if err != nil {
			return writeHookError(e, err)
		}

		return e.JSON(http.StatusOK, response)
	}
}

func createUploadWorkRecordAttachmentHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		requestApp := e.App
		if err := requireTimeEditing(requestApp); err != nil {
			return err
		}

		auth := e.Auth
		if auth == nil {
			return e.Error(http.StatusUnauthorized, "unauthorized", nil)
		}

		workRecordID := strings.TrimSpace(e.Request.PathValue("workRecord"))
		if workRecordID == "" {
			return writeHookError(e, workRecordRouteError(http.StatusBadRequest, "work_record_required", "workRecord is required"))
		}

		// The attachment route accepts multipart/form-data with a single
		// "attachment" file field. The file field itself already restricts allowed
		// MIME types at the collection level, while this route supplies the edit
		// permissions and narrower mutation scope.
		uploadedFiles, err := e.FindUploadedFiles("attachment")
		if errors.Is(err, http.ErrMissingFile) {
			return writeHookError(e, workRecordFieldError("attachment", "required", "Paper copy PDF is required"))
		}
		if err != nil {
			return writeHookError(e, workRecordRouteError(http.StatusBadRequest, "invalid_request_body", "invalid request body"))
		}

		response := WorkRecordAttachmentResponse{}
		err = requestApp.RunInTransaction(func(txApp core.App) error {
			var saveErr error
			response, saveErr = saveWorkRecordAttachment(txApp, auth, workRecordID, uploadedFiles, false)
			return saveErr
		})
		if err != nil {
			return writeHookError(e, err)
		}

		return e.JSON(http.StatusOK, response)
	}
}

func createDeleteWorkRecordAttachmentHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		requestApp := e.App
		if err := requireTimeEditing(requestApp); err != nil {
			return err
		}

		auth := e.Auth
		if auth == nil {
			return e.Error(http.StatusUnauthorized, "unauthorized", nil)
		}

		workRecordID := strings.TrimSpace(e.Request.PathValue("workRecord"))
		if workRecordID == "" {
			return writeHookError(e, workRecordRouteError(http.StatusBadRequest, "work_record_required", "workRecord is required"))
		}

		// Removal is its own endpoint so the UI does not need to send a synthetic
		// empty file payload and so the backend can explicitly clear both the file
		// field and its derived attachment_hash.
		response := WorkRecordAttachmentResponse{}
		err := requestApp.RunInTransaction(func(txApp core.App) error {
			var saveErr error
			response, saveErr = saveWorkRecordAttachment(txApp, auth, workRecordID, nil, true)
			return saveErr
		})
		if err != nil {
			return writeHookError(e, err)
		}

		return e.JSON(http.StatusOK, response)
	}
}

func createCreateWorkRecordNoteHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		requestApp := e.App
		if err := requireTimeEditing(requestApp); err != nil {
			return err
		}

		auth := e.Auth
		if auth == nil {
			return e.Error(http.StatusUnauthorized, "unauthorized", nil)
		}

		workRecordID := strings.TrimSpace(e.Request.PathValue("workRecord"))
		if workRecordID == "" {
			return writeHookError(e, workRecordRouteError(http.StatusBadRequest, "work_record_required", "workRecord is required"))
		}

		var req WorkRecordNoteCreateRequest
		if err := e.BindBody(&req); err != nil {
			return writeHookError(e, workRecordRouteError(http.StatusBadRequest, "invalid_request_body", "invalid request body"))
		}
		if err := normalizeRouteValidationError(validateWorkRecordNoteCreateRequest(req)); err != nil {
			return writeHookError(e, err)
		}

		workRecord, err := findWorkRecordByID(requestApp, workRecordID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return writeHookError(e, workRecordRouteError(http.StatusNotFound, "work_record_not_found", "work record not found"))
			}
			return e.Error(http.StatusInternalServerError, "failed to load work record", err)
		}

		if !canAddWorkRecordNote(requestApp, auth, workRecord.Id, workRecord.GetString("creator")) {
			return writeHookError(e, workRecordRouteError(http.StatusForbidden, "unauthorized", "you are not authorized to add a note to this work record"))
		}

		response := WorkRecordNoteResponse{}
		err = requestApp.RunInTransaction(func(txApp core.App) error {
			notesCollection, err := txApp.FindCollectionByNameOrId(workRecordNotesCollectionIdentifier)
			if err != nil {
				return err
			}

			record := core.NewRecord(notesCollection)
			record.Set("work_record", workRecord.Id)
			record.Set("uid", auth.Id)
			record.Set("note", strings.TrimSpace(req.Note))

			if err := normalizeRouteValidationError(txApp.Save(record)); err != nil {
				return err
			}

			response, err = loadWorkRecordNote(txApp, record.Id)
			return err
		})
		if err != nil {
			return writeHookError(e, err)
		}

		return e.JSON(http.StatusOK, response)
	}
}

func createGetWorkRecordCopyDraftHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		auth := e.Auth
		if auth == nil {
			return e.Error(http.StatusUnauthorized, "unauthorized", nil)
		}

		workRecordID := strings.TrimSpace(e.Request.PathValue("workRecord"))
		if workRecordID == "" {
			return writeHookError(e, workRecordRouteError(http.StatusBadRequest, "work_record_required", "workRecord is required"))
		}

		response, err := loadWorkRecordCopyDraft(app, auth, workRecordID)
		if err != nil {
			return writeHookError(e, err)
		}

		return e.JSON(http.StatusOK, response)
	}
}

func createGetCompatibleWorkRecordsHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		auth := e.Auth
		if auth == nil {
			return e.Error(http.StatusUnauthorized, "unauthorized", nil)
		}

		jobID := strings.TrimSpace(e.Request.URL.Query().Get("job"))
		date := strings.TrimSpace(e.Request.URL.Query().Get("date"))
		excludeEntryID := strings.TrimSpace(e.Request.URL.Query().Get("exclude_entry"))
		if jobID == "" || date == "" {
			return e.Error(http.StatusBadRequest, "job and date are required", nil)
		}

		var rows []WorkRecordOption
		if err := app.DB().NewQuery(`
			SELECT
				wrs.id AS id,
				wr.number || ' - ' || COALESCE(p.given_name || ' ' || p.surname, wrs.uid) AS work_record
			FROM work_records wr
			JOIN work_records_subjects wrs
			  ON wrs.work_record = wr.id
			LEFT JOIN profiles p
			  ON p.uid = wrs.uid
			WHERE wrs.uid = {:uid}
			  AND wr.job = {:job}
			  AND wr.date = {:date}
			  AND NOT EXISTS (
					SELECT 1
					FROM time_entries te
					WHERE te.work_record_subject_id = wrs.id
					  AND ({:excludeEntryId} = '' OR te.id != {:excludeEntryId})
			  )
			ORDER BY wr.number ASC
			LIMIT 100
		`).Bind(dbx.Params{
			"uid":            auth.Id,
			"job":            jobID,
			"date":           date,
			"excludeEntryId": excludeEntryID,
		}).All(&rows); err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load compatible work records", err)
		}

		return e.JSON(http.StatusOK, rows)
	}
}

func createGetWorkRecordsHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		auth := e.Auth
		if auth == nil {
			return e.Error(http.StatusUnauthorized, "unauthorized", nil)
		}

		scope, err := parseWorkRecordListScope(e)
		if err != nil {
			return writeHookError(e, err)
		}

		if scope == "all" {
			canViewAll, err := canViewAllWorkRecords(app, auth)
			if err != nil {
				return e.Error(http.StatusInternalServerError, "failed to check work record permissions", err)
			}
			if !canViewAll {
				return writeHookError(e, &errs.HookError{
					Status:  http.StatusForbidden,
					Message: "you are not authorized to view all work records",
					Data: map[string]errs.CodeError{
						"global": {
							Code:    "unauthorized",
							Message: "you are not authorized to view all work records",
						},
					},
				})
			}
		}

		rows, total, err := loadVisibleWorkRecordRows(app, auth, scope)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load work records", err)
		}

		return e.JSON(http.StatusOK, WorkRecordSearchResponse{
			Items: rows,
			Total: total,
		})
	}
}

func createGetWorkRecordDetailsHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		auth := e.Auth
		if auth == nil {
			return e.Error(http.StatusUnauthorized, "unauthorized", nil)
		}

		workRecordID := strings.TrimSpace(e.Request.PathValue("workRecord"))
		if workRecordID == "" {
			return e.Error(http.StatusBadRequest, "workRecord is required", nil)
		}

		var row workRecordDetailRow
		// Keep detail responses attachment-aware for the same reason as the list
		// response: the editor and details page read record state from this custom
		// route, while the paper-copy PDF itself is created or removed through the
		// dedicated /api/work_records/{id}/attachment endpoints.
		err := app.DB().NewQuery(`
			SELECT
				wr.id AS id,
				wr.number AS work_record,
				wr.date AS date,
				wr.creator AS creator,
				COALESCE(creator_profile.given_name || ' ' || creator_profile.surname, '') AS creator_name,
				COALESCE(wr.job, '') AS job_id,
				COALESCE(j.number, '') AS job_number,
				COALESCE(c.name, '') AS client_name,
				COALESCE(j.description, '') AS job_description,
				COALESCE(j.client_reference_number, '') AS client_reference_number,
				COALESCE(j.client_po, '') AS client_po,
				COALESCE(wr.type, '') AS type_id,
				COALESCE(wrt.name, '') AS type_name,
				COALESCE(wr.location, '') AS location,
				COALESCE(wr.work_description, '') AS work_description,
				COALESCE(wr.report_to, '') AS report_to,
				COALESCE(wr.lane_kms, 0) AS lane_kms,
				COALESCE(wr.field_book_number, '') AS field_book_number,
				COALESCE(wr.field_book_page_number, 0) AS field_book_page_number,
				COALESCE(wr.sub_contractor, '') AS sub_contractor,
				COALESCE(wr.equipment, '') AS equipment,
				COALESCE(wr.supplies, '') AS supplies,
				COALESCE(wr.attachment, '') AS attachment,
				EXISTS(
					SELECT 1
					FROM work_records_subjects wrs_lock
					WHERE wrs_lock.work_record = wr.id
					  AND COALESCE(wrs_lock.approved, 0) = 1
				) AS approved,
				COALESCE(wr.parent_work_record, '') AS parent_work_record_id,
				COALESCE(parent_wr.number, '') AS parent_work_record_number
			FROM work_records wr
			LEFT JOIN profiles creator_profile
			  ON creator_profile.uid = wr.creator
			LEFT JOIN jobs j
			  ON wr.job = j.id
			LEFT JOIN clients c
			  ON j.client = c.id
			LEFT JOIN work_records_types wrt
			  ON wr.type = wrt.id
			LEFT JOIN work_records parent_wr
			  ON parent_wr.id = wr.parent_work_record
			WHERE wr.id = {:id}
			LIMIT 1
		`).Bind(dbx.Params{
			"id": workRecordID,
		}).One(&row)
		if err != nil {
			return e.Error(http.StatusNotFound, "work record not found", err)
		}

		subjects, err := loadWorkRecordSubjects(app, row.ID, auth, row.Creator)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load work record workers", err)
		}
		isSubject := false
		for _, subject := range subjects {
			if subject.UID == auth.Id {
				isSubject = true
				break
			}
		}

		if !isSubject && row.Creator != auth.Id {
			canViewAll, err := canViewAllWorkRecords(app, auth)
			if err != nil {
				return e.Error(http.StatusInternalServerError, "failed to check work record permissions", err)
			}
			if !canViewAll {
				return writeHookError(e, &errs.HookError{
					Status:  http.StatusForbidden,
					Message: "you are not authorized to view this work record",
					Data: map[string]errs.CodeError{
						"global": {
							Code:    "unauthorized",
							Message: "you are not authorized to view this work record",
						},
					},
				})
			}
		}

		notes, err := loadWorkRecordNotes(app, row.ID)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load work record notes", err)
		}
		consumables, err := loadWorkRecordConsumables(app, row.ID, row.TypeID)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load work record consumables", err)
		}

		canCopyToTomorrow := row.Creator == auth.Id
		if !canCopyToTomorrow {
			targetUIDs := make([]string, 0, len(subjects))
			for _, subject := range subjects {
				targetUIDs = append(targetUIDs, subject.UID)
			}
			ok, err := canCreateWorkRecordForRoster(app, auth, targetUIDs)
			if err != nil {
				return e.Error(http.StatusInternalServerError, "failed to check copy permissions", err)
			}
			canCopyToTomorrow = ok
		}
		canMutateShared := (isSubject || row.Creator == auth.Id) && !row.Approved
		canMutateRecord := canMutateShared
		if !canMutateRecord {
			for _, subject := range subjects {
				if subject.CanEdit {
					canMutateRecord = true
					break
				}
			}
		}
		return e.JSON(http.StatusOK, WorkRecordDetailResponse{
			ID:                    row.ID,
			WorkRecord:            row.WorkRecord,
			Date:                  row.Date,
			Creator:               row.Creator,
			CreatorName:           row.CreatorName,
			JobID:                 row.JobID,
			JobNumber:             row.JobNumber,
			ClientName:            row.ClientName,
			JobDescription:        row.JobDescription,
			ClientReferenceNumber: row.ClientReferenceNumber,
			ClientPO:              row.ClientPO,
			TypeID:                row.TypeID,
			TypeName:              row.TypeName,
			Location:              row.Location,
			WorkDescription:       row.WorkDescription,
			ReportTo:              row.ReportTo,
			LaneKms:               row.LaneKms,
			FieldBookNumber:       row.FieldBookNumber,
			FieldBookPageNumber:   row.FieldBookPageNumber,
			SubContractor:         row.SubContractor,
			Equipment:             row.Equipment,
			Supplies:              row.Supplies,
			Attachment:            row.Attachment,
			Approved:              row.Approved,
			ParentWorkRecordID:    row.ParentWorkRecordID,
			ParentWorkRecord:      row.ParentWorkRecordNumber,
			RootWorkRecordID:      firstNonEmpty(row.ParentWorkRecordID, row.ID),
			RootWorkRecord:        firstNonEmpty(row.ParentWorkRecordNumber, row.WorkRecord),
			CanCopyToTomorrow:     canCopyToTomorrow,
			CanEdit:               canMutateRecord,
			CanEditShared:         canMutateShared,
			CanAddNote:            canAddWorkRecordNote(app, auth, row.ID, row.Creator),
			Subjects:              subjects,
			Notes:                 notes,
			Consumables:           consumables,
		})
	}
}

func valueOrEmpty(record *core.Record, field string) string {
	if record == nil {
		return ""
	}
	return record.GetString(field)
}

func boolToSQLInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}

func buildWorkRecordSearchText(workRecord string) string {
	workRecord = strings.ToUpper(strings.TrimSpace(workRecord))
	if workRecord == "" {
		return ""
	}

	return strings.Join(strings.Fields(strings.ReplaceAll(workRecord, "-", " ")), " ")
}
