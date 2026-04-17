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
)

const (
	workRecordsCollectionIdentifier                 = "work_records"
	workRecordConsumableEntriesCollectionIdentifier = "work_records_consumable_entries"
)

type WorkRecordListRow struct {
	ID                         string `db:"id" json:"id"`
	WorkRecord                 string `db:"work_record" json:"work_record"`
	Date                       string `db:"date" json:"date"`
	UID                        string `db:"uid" json:"uid"`
	OwnerName                  string `db:"owner_name" json:"owner_name"`
	Creator                    string `db:"creator" json:"creator"`
	JobID                      string `db:"job_id" json:"job_id"`
	JobNumber                  string `db:"job_number" json:"job_number"`
	TypeID                     string `db:"type_id" json:"type_id"`
	TypeName                   string `db:"type_name" json:"type_name"`
	Location                   string `db:"location" json:"location"`
	WorkDescription            string `db:"work_description" json:"work_description"`
	Approved                   bool   `db:"approved" json:"approved"`
	LinkedTimeEntryID          string `db:"linked_time_entry_id" json:"linked_time_entry_id"`
	LinkedTimeSheetID          string `db:"linked_timesheet_id" json:"linked_timesheet_id"`
	LinkedTimeEntryDate        string `db:"linked_time_entry_date" json:"linked_time_entry_date"`
	LinkedTimeEntryDescription string `db:"linked_time_entry_description" json:"linked_time_entry_description"`
	SearchText                 string `json:"search_text"`
	CanCreateTimeEntry         bool   `json:"can_create_time_entry"`
}

type linkedWorkRecordEntryRow struct {
	LinkedTimeEntryID          string `db:"linked_time_entry_id" json:"linked_time_entry_id"`
	LinkedTimeSheetID          string `db:"linked_timesheet_id" json:"linked_timesheet_id"`
	LinkedTimeEntryDate        string `db:"linked_time_entry_date" json:"linked_time_entry_date"`
	LinkedTimeEntryDescription string `db:"linked_time_entry_description" json:"linked_time_entry_description"`
	LinkedTimeSheetWeekEnding  string `db:"linked_timesheet_week_ending" json:"linked_timesheet_week_ending"`
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

type WorkRecordDetailResponse struct {
	ID                         string                         `json:"id"`
	WorkRecord                 string                         `json:"work_record"`
	Date                       string                         `json:"date"`
	UID                        string                         `json:"uid"`
	OwnerName                  string                         `json:"owner_name"`
	Creator                    string                         `json:"creator"`
	JobID                      string                         `json:"job_id"`
	JobNumber                  string                         `json:"job_number"`
	TypeID                     string                         `json:"type_id"`
	TypeName                   string                         `json:"type_name"`
	Location                   string                         `json:"location"`
	WorkDescription            string                         `json:"work_description"`
	ReportTo                   string                         `json:"report_to"`
	HoursOnSite                float64                        `json:"hours_on_site"`
	HoursTravelTime            float64                        `json:"hours_travel_time"`
	DistanceTravelledKM        float64                        `json:"distance_travelled_km"`
	FieldBookNumber            string                         `json:"field_book_number"`
	FieldBookPageNumber        int                            `json:"field_book_page_number"`
	SubContractor              string                         `json:"sub_contractor"`
	Equipment                  string                         `json:"equipment"`
	Supplies                   string                         `json:"supplies"`
	VehicleType                string                         `json:"vehicle_type"`
	IsPassenger                bool                           `json:"is_passenger"`
	CompanyVehicleUnitNumber   string                         `json:"company_vehicle_unit_number"`
	Approved                   bool                           `json:"approved"`
	ParentWorkRecordID         string                         `json:"parent_work_record_id"`
	ParentWorkRecord           string                         `json:"parent_work_record"`
	RootWorkRecordID           string                         `json:"root_work_record_id"`
	RootWorkRecord             string                         `json:"root_work_record"`
	LinkedTimeEntryID          string                         `json:"linked_time_entry_id"`
	LinkedTimeSheetID          string                         `json:"linked_timesheet_id"`
	LinkedTimeEntryDate        string                         `json:"linked_time_entry_date"`
	LinkedTimeEntryDescription string                         `json:"linked_time_entry_description"`
	LinkedTimeSheetWeekEnding  string                         `json:"linked_time_sheet_week_ending"`
	CanCreateTimeEntry         bool                           `json:"can_create_time_entry"`
	CanCopyToTomorrow          bool                           `json:"can_copy_to_tomorrow"`
	CanEdit                    bool                           `json:"can_edit"`
	CanAddNote                 bool                           `json:"can_add_note"`
	Notes                      []WorkRecordNoteResponse       `json:"notes"`
	Consumables                []WorkRecordConsumableResponse `json:"consumables"`
}

type WorkRecordSaveConsumableRequest struct {
	EntryID         string  `json:"entry_id"`
	ConsumableID    string  `json:"consumable_id"`
	InputKind       string  `json:"input_kind"`
	QuantityNumber  float64 `json:"quantity_number"`
	BooleanSelected bool    `json:"boolean_selected"`
}

type WorkRecordCreateRequest struct {
	UID                      string                            `json:"uid"`
	Date                     string                            `json:"date"`
	Job                      string                            `json:"job"`
	Type                     string                            `json:"type"`
	Location                 string                            `json:"location"`
	WorkDescription          string                            `json:"work_description"`
	ReportTo                 string                            `json:"report_to"`
	HoursOnSite              float64                           `json:"hours_on_site"`
	HoursTravelTime          float64                           `json:"hours_travel_time"`
	DistanceTravelledKM      float64                           `json:"distance_travelled_km"`
	FieldBookNumber          string                            `json:"field_book_number"`
	FieldBookPageNumber      any                               `json:"field_book_page_number"`
	SubContractor            string                            `json:"sub_contractor"`
	Equipment                string                            `json:"equipment"`
	Supplies                 string                            `json:"supplies"`
	VehicleType              string                            `json:"vehicle_type"`
	IsPassenger              bool                              `json:"is_passenger"`
	CompanyVehicleUnitNumber string                            `json:"company_vehicle_unit_number"`
	Consumables              []WorkRecordSaveConsumableRequest `json:"consumables"`
	CopyFromWorkRecordID     string                            `json:"copy_from_work_record_id"`
}

type WorkRecordSaveRequest struct {
	Date                     string                            `json:"date"`
	Job                      string                            `json:"job"`
	Type                     string                            `json:"type"`
	Location                 string                            `json:"location"`
	WorkDescription          string                            `json:"work_description"`
	ReportTo                 string                            `json:"report_to"`
	HoursOnSite              float64                           `json:"hours_on_site"`
	HoursTravelTime          float64                           `json:"hours_travel_time"`
	DistanceTravelledKM      float64                           `json:"distance_travelled_km"`
	FieldBookNumber          string                            `json:"field_book_number"`
	FieldBookPageNumber      any                               `json:"field_book_page_number"`
	SubContractor            string                            `json:"sub_contractor"`
	Equipment                string                            `json:"equipment"`
	Supplies                 string                            `json:"supplies"`
	VehicleType              string                            `json:"vehicle_type"`
	IsPassenger              bool                              `json:"is_passenger"`
	CompanyVehicleUnitNumber string                            `json:"company_vehicle_unit_number"`
	Consumables              []WorkRecordSaveConsumableRequest `json:"consumables"`
}

type WorkRecordSaveResponse struct {
	ID string `json:"id"`
}

type WorkRecordCopyDraftResponse struct {
	UID                      string                         `json:"uid"`
	Date                     string                         `json:"date"`
	Job                      string                         `json:"job"`
	Type                     string                         `json:"type"`
	Location                 string                         `json:"location"`
	WorkDescription          string                         `json:"work_description"`
	ReportTo                 string                         `json:"report_to"`
	HoursOnSite              float64                        `json:"hours_on_site"`
	HoursTravelTime          float64                        `json:"hours_travel_time"`
	DistanceTravelledKM      float64                        `json:"distance_travelled_km"`
	FieldBookNumber          string                         `json:"field_book_number"`
	FieldBookPageNumber      int                            `json:"field_book_page_number"`
	SubContractor            string                         `json:"sub_contractor"`
	Equipment                string                         `json:"equipment"`
	Supplies                 string                         `json:"supplies"`
	VehicleType              string                         `json:"vehicle_type"`
	IsPassenger              bool                           `json:"is_passenger"`
	CompanyVehicleUnitNumber string                         `json:"company_vehicle_unit_number"`
	Consumables              []WorkRecordConsumableResponse `json:"consumables"`
	CopyFromWorkRecordID     string                         `json:"copy_from_work_record_id"`
	SourceWorkRecord         string                         `json:"source_work_record"`
	RootWorkRecordID         string                         `json:"root_work_record_id"`
	RootWorkRecord           string                         `json:"root_work_record"`
}

type WorkRecordOption struct {
	ID         string `db:"id" json:"id"`
	WorkRecord string `db:"work_record" json:"work_record"`
}

type workRecordDetailRow struct {
	ID                       string  `db:"id"`
	WorkRecord               string  `db:"work_record"`
	Date                     string  `db:"date"`
	UID                      string  `db:"uid"`
	OwnerName                string  `db:"owner_name"`
	Creator                  string  `db:"creator"`
	JobID                    string  `db:"job_id"`
	JobNumber                string  `db:"job_number"`
	TypeID                   string  `db:"type_id"`
	TypeName                 string  `db:"type_name"`
	Location                 string  `db:"location"`
	WorkDescription          string  `db:"work_description"`
	ReportTo                 string  `db:"report_to"`
	HoursOnSite              float64 `db:"hours_on_site"`
	HoursTravelTime          float64 `db:"hours_travel_time"`
	DistanceTravelledKM      float64 `db:"distance_travelled_km"`
	FieldBookNumber          string  `db:"field_book_number"`
	FieldBookPageNumber      int     `db:"field_book_page_number"`
	SubContractor            string  `db:"sub_contractor"`
	Equipment                string  `db:"equipment"`
	Supplies                 string  `db:"supplies"`
	VehicleType              string  `db:"vehicle_type"`
	IsPassenger              bool    `db:"is_passenger"`
	CompanyVehicleUnitNumber string  `db:"company_vehicle_unit_number"`
	Approved                 bool    `db:"approved"`
	ParentWorkRecordID       string  `db:"parent_work_record_id"`
	ParentWorkRecordNumber   string  `db:"parent_work_record_number"`
}

type workRecordCopySourceRow struct {
	ID                       string  `db:"id"`
	WorkRecord               string  `db:"work_record"`
	UID                      string  `db:"uid"`
	Creator                  string  `db:"creator"`
	Date                     string  `db:"date"`
	Job                      string  `db:"job"`
	Type                     string  `db:"type"`
	Location                 string  `db:"location"`
	WorkDescription          string  `db:"work_description"`
	ReportTo                 string  `db:"report_to"`
	HoursOnSite              float64 `db:"hours_on_site"`
	HoursTravelTime          float64 `db:"hours_travel_time"`
	DistanceTravelledKM      float64 `db:"distance_travelled_km"`
	FieldBookNumber          string  `db:"field_book_number"`
	FieldBookPageNumber      int     `db:"field_book_page_number"`
	SubContractor            string  `db:"sub_contractor"`
	Equipment                string  `db:"equipment"`
	Supplies                 string  `db:"supplies"`
	VehicleType              string  `db:"vehicle_type"`
	IsPassenger              bool    `db:"is_passenger"`
	CompanyVehicleUnitNumber string  `db:"company_vehicle_unit_number"`
	ParentWorkRecordID       string  `db:"parent_work_record_id"`
	ParentWorkRecord         string  `db:"parent_work_record"`
}

type workRecordConsumableEntryState struct {
	ID           string `db:"id"`
	ConsumableID string `db:"consumable_id"`
}

type workRecordSaveStateRow struct {
	ID                       string  `db:"id"`
	Number                   string  `db:"number"`
	Date                     string  `db:"date"`
	UID                      string  `db:"uid"`
	Creator                  string  `db:"creator"`
	Job                      string  `db:"job"`
	Type                     string  `db:"type"`
	Location                 string  `db:"location"`
	WorkDescription          string  `db:"work_description"`
	ReportTo                 string  `db:"report_to"`
	HoursOnSite              float64 `db:"hours_on_site"`
	HoursTravelTime          float64 `db:"hours_travel_time"`
	DistanceTravelledKM      float64 `db:"distance_travelled_km"`
	FieldBookNumber          string  `db:"field_book_number"`
	FieldBookPageNumber      int     `db:"field_book_page_number"`
	SubContractor            string  `db:"sub_contractor"`
	Equipment                string  `db:"equipment"`
	Supplies                 string  `db:"supplies"`
	VehicleType              string  `db:"vehicle_type"`
	IsPassenger              bool    `db:"is_passenger"`
	CompanyVehicleUnitNumber string  `db:"company_vehicle_unit_number"`
	Approved                 bool    `db:"approved"`
	ParentWorkRecord         string  `db:"parent_work_record"`
}

type workRecordConsumableEntryPersistenceRow struct {
	ID             string  `db:"id"`
	WorkRecordID   string  `db:"work_record"`
	ConsumableID   string  `db:"consumable"`
	QuantityNumber float64 `db:"quantity_number"`
}

func newWorkRecordPersistenceCollection() *core.Collection {
	collection := core.NewBaseCollection(workRecordsCollectionIdentifier)
	collection.Fields.Add(&core.TextField{Name: "number"})
	collection.Fields.Add(&core.TextField{Name: "date"})
	collection.Fields.Add(&core.TextField{Name: "uid"})
	collection.Fields.Add(&core.TextField{Name: "creator"})
	collection.Fields.Add(&core.TextField{Name: "job"})
	collection.Fields.Add(&core.TextField{Name: "type"})
	collection.Fields.Add(&core.TextField{Name: "location"})
	collection.Fields.Add(&core.TextField{Name: "work_description"})
	collection.Fields.Add(&core.TextField{Name: "report_to"})
	collection.Fields.Add(&core.NumberField{Name: "hours_on_site"})
	collection.Fields.Add(&core.NumberField{Name: "hours_travel_time"})
	collection.Fields.Add(&core.NumberField{Name: "distance_travelled_km"})
	collection.Fields.Add(&core.TextField{Name: "field_book_number"})
	collection.Fields.Add(&core.NumberField{Name: "field_book_page_number"})
	collection.Fields.Add(&core.TextField{Name: "sub_contractor"})
	collection.Fields.Add(&core.TextField{Name: "equipment"})
	collection.Fields.Add(&core.TextField{Name: "supplies"})
	collection.Fields.Add(&core.TextField{Name: "vehicle_type"})
	collection.Fields.Add(&core.BoolField{Name: "is_passenger"})
	collection.Fields.Add(&core.TextField{Name: "company_vehicle_unit_number"})
	collection.Fields.Add(&core.BoolField{Name: "approved"})
	collection.Fields.Add(&core.TextField{Name: "parent_work_record"})
	return collection
}

func newWorkRecordConsumableEntryPersistenceCollection() *core.Collection {
	collection := core.NewBaseCollection(workRecordConsumableEntriesCollectionIdentifier)
	collection.Fields.Add(&core.TextField{Name: "work_record"})
	collection.Fields.Add(&core.TextField{Name: "consumable"})
	collection.Fields.Add(&core.NumberField{Name: "quantity_number"})
	return collection
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

	record := core.NewRecord(newWorkRecordConsumableEntryPersistenceCollection())
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
		filterSQL += " AND (wr.uid = {:authId} OR wr.creator = {:authId})"
		params["authId"] = auth.Id
	}

	rows := []WorkRecordListRow{}
	if err := app.DB().NewQuery(`
		SELECT
			wr.id AS id,
			wr.number AS work_record,
			wr.date AS date,
			wr.uid AS uid,
			COALESCE(p.given_name || ' ' || p.surname, '') AS owner_name,
			wr.creator AS creator,
			COALESCE(wr.job, '') AS job_id,
			COALESCE(j.number, '') AS job_number,
			COALESCE(wr.type, '') AS type_id,
			COALESCE(wrt.name, '') AS type_name,
			COALESCE(wr.location, '') AS location,
			COALESCE(wr.work_description, '') AS work_description,
			COALESCE(wr.approved, 0) AS approved,
			COALESCE(te.id, '') AS linked_time_entry_id,
			COALESCE(te.tsid, '') AS linked_timesheet_id,
			COALESCE(te.date, '') AS linked_time_entry_date,
			COALESCE(te.description, '') AS linked_time_entry_description
		FROM work_records wr
		LEFT JOIN profiles p
		  ON p.uid = wr.uid
		LEFT JOIN jobs j
		  ON j.id = wr.job
		LEFT JOIN work_records_types wrt
		  ON wrt.id = wr.type
		LEFT JOIN time_entries te
		  ON te.work_record_id = wr.id
		WHERE ` + filterSQL + `
		ORDER BY wr.date DESC, wr.number DESC
	`).Bind(params).All(&rows); err != nil {
		return nil, 0, err
	}

	for i := range rows {
		rows[i].SearchText = buildWorkRecordSearchText(rows[i].WorkRecord + " " + rows[i].OwnerName + " " + rows[i].JobNumber + " " + rows[i].TypeName + " " + rows[i].Location + " " + rows[i].WorkDescription)
		rows[i].CanCreateTimeEntry = rows[i].UID == auth.Id && rows[i].LinkedTimeEntryID == ""
	}

	return rows, len(rows), nil
}

func canViewWorkRecord(app core.App, auth *core.Record, record *core.Record) (bool, error) {
	if auth == nil || record == nil {
		return false, nil
	}

	if record.GetString("uid") == auth.Id || record.GetString("creator") == auth.Id {
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

func loadLinkedTimeEntry(app core.App, workRecordID string) (linkedWorkRecordEntryRow, error) {
	var row linkedWorkRecordEntryRow
	err := app.DB().NewQuery(`
		SELECT
			COALESCE(te.id, '') AS linked_time_entry_id,
			COALESCE(te.tsid, '') AS linked_timesheet_id,
			COALESCE(te.date, '') AS linked_time_entry_date,
			COALESCE(te.description, '') AS linked_time_entry_description,
			COALESCE(ts.week_ending, '') AS linked_timesheet_week_ending
		FROM time_entries te
		LEFT JOIN time_sheets ts
		  ON ts.id = te.tsid
		WHERE te.work_record_id = {:workRecordId}
		LIMIT 1
	`).Bind(dbx.Params{
		"workRecordId": workRecordID,
	}).One(&row)
	if err != nil {
		if err == sql.ErrNoRows {
			return linkedWorkRecordEntryRow{}, nil
		}
		return linkedWorkRecordEntryRow{}, err
	}

	return row, nil
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

func loadWorkRecordConsumables(app core.App, workRecordID string, typeID string) ([]WorkRecordConsumableResponse, error) {
	rows := []WorkRecordConsumableResponse{}
	if err := app.DB().NewQuery(`
		SELECT
			COALESCE(wrce.id, '') AS entry_id,
			COALESCE(wrc.id, '') AS consumable_id,
			COALESCE(wrc.name, '') AS consumable_name,
			COALESCE(wrc.input_kind, '') AS input_kind,
			COALESCE(wrtc.unit_label, '') AS unit_label,
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

func findWorkRecordByToken(app core.App, token string) (*core.Record, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, sql.ErrNoRows
	}

	collection := newWorkRecordPersistenceCollection()
	row := workRecordSaveStateRow{}
	if err := app.DB().NewQuery(`
		SELECT
			wr.id AS id,
			COALESCE(wr.number, '') AS number,
			COALESCE(wr.date, '') AS date,
			COALESCE(wr.uid, '') AS uid,
			COALESCE(wr.creator, '') AS creator,
			COALESCE(wr.job, '') AS job,
			COALESCE(wr.type, '') AS type,
			COALESCE(wr.location, '') AS location,
			COALESCE(wr.work_description, '') AS work_description,
			COALESCE(wr.report_to, '') AS report_to,
			COALESCE(wr.hours_on_site, 0) AS hours_on_site,
			COALESCE(wr.hours_travel_time, 0) AS hours_travel_time,
			COALESCE(wr.distance_travelled_km, 0) AS distance_travelled_km,
			COALESCE(wr.field_book_number, '') AS field_book_number,
			COALESCE(wr.field_book_page_number, 0) AS field_book_page_number,
			COALESCE(wr.sub_contractor, '') AS sub_contractor,
			COALESCE(wr.equipment, '') AS equipment,
			COALESCE(wr.supplies, '') AS supplies,
			COALESCE(wr.vehicle_type, '') AS vehicle_type,
			COALESCE(wr.is_passenger, 0) AS is_passenger,
			COALESCE(wr.company_vehicle_unit_number, '') AS company_vehicle_unit_number,
			COALESCE(wr.approved, 0) AS approved,
			COALESCE(wr.parent_work_record, '') AS parent_work_record
		FROM work_records wr
		WHERE wr.id = {:token} OR wr.number = {:token}
		LIMIT 1
	`).Bind(dbx.Params{
		"token": token,
	}).One(&row); err != nil {
		return nil, err
	}

	record := core.NewRecord(collection)
	record.Load(map[string]any{
		"id":                          row.ID,
		"number":                      row.Number,
		"date":                        row.Date,
		"uid":                         row.UID,
		"creator":                     row.Creator,
		"job":                         row.Job,
		"type":                        row.Type,
		"location":                    row.Location,
		"work_description":            row.WorkDescription,
		"report_to":                   row.ReportTo,
		"hours_on_site":               row.HoursOnSite,
		"hours_travel_time":           row.HoursTravelTime,
		"distance_travelled_km":       row.DistanceTravelledKM,
		"field_book_number":           row.FieldBookNumber,
		"field_book_page_number":      row.FieldBookPageNumber,
		"sub_contractor":              row.SubContractor,
		"equipment":                   row.Equipment,
		"supplies":                    row.Supplies,
		"vehicle_type":                row.VehicleType,
		"is_passenger":                row.IsPassenger,
		"company_vehicle_unit_number": row.CompanyVehicleUnitNumber,
		"approved":                    row.Approved,
		"parent_work_record":          row.ParentWorkRecord,
	})
	if err := record.PostScan(); err != nil {
		return nil, err
	}

	return record, nil
}

func findWorkRecordCopySourceByToken(app core.App, token string) (workRecordCopySourceRow, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return workRecordCopySourceRow{}, sql.ErrNoRows
	}

	var row workRecordCopySourceRow
	err := app.DB().NewQuery(`
		SELECT
			wr.id AS id,
			wr.number AS work_record,
			wr.uid AS uid,
			wr.creator AS creator,
			wr.date AS date,
			COALESCE(wr.job, '') AS job,
			COALESCE(wr.type, '') AS type,
			COALESCE(wr.location, '') AS location,
			COALESCE(wr.work_description, '') AS work_description,
			COALESCE(wr.report_to, '') AS report_to,
			COALESCE(wr.hours_on_site, 0) AS hours_on_site,
			COALESCE(wr.hours_travel_time, 0) AS hours_travel_time,
			COALESCE(wr.distance_travelled_km, 0) AS distance_travelled_km,
			COALESCE(wr.field_book_number, '') AS field_book_number,
			COALESCE(wr.field_book_page_number, 0) AS field_book_page_number,
			COALESCE(wr.sub_contractor, '') AS sub_contractor,
			COALESCE(wr.equipment, '') AS equipment,
			COALESCE(wr.supplies, '') AS supplies,
			COALESCE(wr.vehicle_type, '') AS vehicle_type,
			COALESCE(wr.is_passenger, 0) AS is_passenger,
			COALESCE(wr.company_vehicle_unit_number, '') AS company_vehicle_unit_number,
			COALESCE(wr.parent_work_record, '') AS parent_work_record_id,
			COALESCE(parent_wr.number, '') AS parent_work_record
		FROM work_records wr
		LEFT JOIN work_records parent_wr
		  ON parent_wr.id = wr.parent_work_record
		WHERE wr.id = {:token} OR wr.number = {:token}
		LIMIT 1
	`).Bind(dbx.Params{
		"token": token,
	}).One(&row)
	if err != nil {
		return workRecordCopySourceRow{}, err
	}

	return row, nil
}

func toWorkRecordSaveRequest(req WorkRecordCreateRequest) WorkRecordSaveRequest {
	return WorkRecordSaveRequest{
		Date:                     req.Date,
		Job:                      req.Job,
		Type:                     req.Type,
		Location:                 req.Location,
		WorkDescription:          req.WorkDescription,
		ReportTo:                 req.ReportTo,
		HoursOnSite:              req.HoursOnSite,
		HoursTravelTime:          req.HoursTravelTime,
		DistanceTravelledKM:      req.DistanceTravelledKM,
		FieldBookNumber:          req.FieldBookNumber,
		FieldBookPageNumber:      req.FieldBookPageNumber,
		SubContractor:            req.SubContractor,
		Equipment:                req.Equipment,
		Supplies:                 req.Supplies,
		VehicleType:              req.VehicleType,
		IsPassenger:              req.IsPassenger,
		CompanyVehicleUnitNumber: req.CompanyVehicleUnitNumber,
		Consumables:              req.Consumables,
	}
}

func canCreateWorkRecordForUID(app core.App, auth *core.Record, targetUID string) (bool, error) {
	if auth == nil {
		return false, nil
	}

	targetUID = strings.TrimSpace(targetUID)
	if targetUID == "" {
		return false, nil
	}

	if targetUID == auth.Id {
		hasTimeClaim, err := utilities.HasClaim(app, auth, "time")
		if err != nil {
			return false, err
		}
		if hasTimeClaim {
			return true, nil
		}

		return utilities.HasClaim(app, auth, "work_record")
	}

	hasWorkRecordClaim, err := utilities.HasClaim(app, auth, "work_record")
	if err != nil || !hasWorkRecordClaim {
		return false, err
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

	if source.UID != auth.Id && source.Creator != auth.Id {
		canViewAll, err := canViewAllWorkRecords(app, auth)
		if err != nil {
			return err
		}
		if !canViewAll {
			return workRecordRouteError(http.StatusForbidden, "unauthorized", "you are not authorized to view this work record")
		}
	}

	if err := requireWorkRecordCreateAccess(app, auth, source.UID, "copy this work record"); err != nil {
		return err
	}

	return nil
}

func requireWorkRecordCreateAccess(app core.App, auth *core.Record, targetUID string, action string) error {
	canCreate, err := canCreateWorkRecordForUID(app, auth, targetUID)
	if err != nil {
		return err
	}
	if !canCreate {
		return workRecordRouteError(http.StatusForbidden, "unauthorized", "you are not authorized to "+action)
	}
	return nil
}

func loadWorkRecordCopyDraft(app core.App, auth *core.Record, workRecordToken string) (WorkRecordCopyDraftResponse, error) {
	source, err := findWorkRecordCopySourceByToken(app, workRecordToken)
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

	rootWorkRecordID, rootWorkRecord := resolveWorkRecordRootParent(source)

	return WorkRecordCopyDraftResponse{
		UID:                      source.UID,
		Date:                     sourceDate.AddDate(0, 0, 1).Format("2006-01-02"),
		Job:                      source.Job,
		Type:                     source.Type,
		Location:                 source.Location,
		WorkDescription:          source.WorkDescription,
		ReportTo:                 source.ReportTo,
		HoursOnSite:              source.HoursOnSite,
		HoursTravelTime:          source.HoursTravelTime,
		DistanceTravelledKM:      source.DistanceTravelledKM,
		FieldBookNumber:          source.FieldBookNumber,
		FieldBookPageNumber:      source.FieldBookPageNumber,
		SubContractor:            source.SubContractor,
		Equipment:                source.Equipment,
		Supplies:                 source.Supplies,
		VehicleType:              source.VehicleType,
		IsPassenger:              source.IsPassenger,
		CompanyVehicleUnitNumber: source.CompanyVehicleUnitNumber,
		Consumables:              consumables,
		CopyFromWorkRecordID:     source.ID,
		SourceWorkRecord:         source.WorkRecord,
		RootWorkRecordID:         rootWorkRecordID,
		RootWorkRecord:           rootWorkRecord,
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

func normalizeWorkRecordValidationError(err error) error {
	if err == nil {
		return nil
	}

	var hookErr *errs.HookError
	if errors.As(err, &hookErr) {
		return hookErr
	}

	var validationErrs validation.Errors
	if errors.As(err, &validationErrs) {
		fieldErrors := make(map[string]errs.CodeError, len(validationErrs))
		for field, fieldErr := range validationErrs {
			fieldErrors[field] = errs.CodeError{
				Code:    "validation_error",
				Message: fieldErr.Error(),
			}
		}
		return &errs.HookError{
			Status:  http.StatusBadRequest,
			Message: "validation failed",
			Data:    fieldErrors,
		}
	}

	return err
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

func applyWorkRecordSaveRequest(record *core.Record, req WorkRecordSaveRequest) error {
	record.Set("date", req.Date)
	record.Set("job", req.Job)
	record.Set("type", req.Type)
	record.Set("location", req.Location)
	record.Set("work_description", req.WorkDescription)
	record.Set("report_to", req.ReportTo)
	record.Set("hours_on_site", req.HoursOnSite)
	record.Set("hours_travel_time", req.HoursTravelTime)
	record.Set("distance_travelled_km", req.DistanceTravelledKM)
	record.Set("field_book_number", req.FieldBookNumber)

	switch value := req.FieldBookPageNumber.(type) {
	case nil:
		record.Set("field_book_page_number", "")
	case string:
		if strings.TrimSpace(value) == "" {
			record.Set("field_book_page_number", "")
		} else {
			return &errs.HookError{
				Status:  http.StatusBadRequest,
				Message: "validation failed",
				Data: map[string]errs.CodeError{
					"field_book_page_number": {
						Code:    "invalid_field_value",
						Message: "Field book page number must be a number or empty",
					},
				},
			}
		}
	case float64, float32, int, int32, int64:
		record.Set("field_book_page_number", value)
	default:
		return &errs.HookError{
			Status:  http.StatusBadRequest,
			Message: "validation failed",
			Data: map[string]errs.CodeError{
				"field_book_page_number": {
					Code:    "invalid_field_value",
					Message: "Field book page number must be a number or empty",
				},
			},
		}
	}

	record.Set("sub_contractor", req.SubContractor)
	record.Set("equipment", req.Equipment)
	record.Set("supplies", req.Supplies)
	record.Set("vehicle_type", req.VehicleType)
	record.Set("is_passenger", req.IsPassenger)
	record.Set("company_vehicle_unit_number", req.CompanyVehicleUnitNumber)

	return nil
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

func saveWorkRecordEdit(txApp core.App, auth *core.Record, record *core.Record, req WorkRecordSaveRequest) (WorkRecordSaveResponse, error) {
	currentRecord, err := findWorkRecordByToken(txApp, record.Id)
	if err != nil {
		return WorkRecordSaveResponse{}, err
	}

	if currentRecord.GetString("uid") != auth.Id && currentRecord.GetString("creator") != auth.Id {
		return WorkRecordSaveResponse{}, workRecordRouteError(http.StatusForbidden, "unauthorized", "you are not authorized to edit this work record")
	}

	if currentRecord.GetBool("approved") {
		return WorkRecordSaveResponse{}, workRecordRouteError(http.StatusBadRequest, "approved_work_record_locked", "Approved work records cannot be edited")
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
		if err := normalizeWorkRecordValidationError(apphooks.ValidateWorkRecordConsumableEntryDelete(txApp, entryRecord)); err != nil {
			return WorkRecordSaveResponse{}, err
		}
		if err := normalizeWorkRecordValidationError(txApp.Delete(entryRecord)); err != nil {
			return WorkRecordSaveResponse{}, err
		}
	}

	if err := applyWorkRecordSaveRequest(currentRecord, req); err != nil {
		return WorkRecordSaveResponse{}, err
	}
	if err := normalizeWorkRecordValidationError(apphooks.ValidateWorkRecord(txApp, currentRecord)); err != nil {
		return WorkRecordSaveResponse{}, err
	}
	if err := normalizeWorkRecordValidationError(txApp.Save(currentRecord)); err != nil {
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
			entryRecord = core.NewRecord(newWorkRecordConsumableEntryPersistenceCollection())
		}

		entryRecord.Set("work_record", currentRecord.Id)
		entryRecord.Set("consumable", item.ConsumableID)
		setWorkRecordConsumableQuantity(entryRecord, inputKindsByConsumable[item.ConsumableID], item)

		if err := normalizeWorkRecordValidationError(apphooks.ValidateWorkRecordConsumableEntry(txApp, entryRecord)); err != nil {
			return WorkRecordSaveResponse{}, err
		}
		if err := normalizeWorkRecordValidationError(txApp.Save(entryRecord)); err != nil {
			return WorkRecordSaveResponse{}, err
		}
	}

	return WorkRecordSaveResponse{ID: currentRecord.Id}, nil
}

func createWorkRecord(txApp core.App, auth *core.Record, req WorkRecordCreateRequest) (WorkRecordSaveResponse, error) {
	targetUID := strings.TrimSpace(req.UID)
	if targetUID == "" {
		return WorkRecordSaveResponse{}, &errs.HookError{
			Status:  http.StatusBadRequest,
			Message: "validation failed",
			Data: map[string]errs.CodeError{
				"uid": {
					Code:    "required",
					Message: "Work record owner is required",
				},
			},
		}
	}

	rootParentID := ""
	copyFromWorkRecordID := strings.TrimSpace(req.CopyFromWorkRecordID)
	if copyFromWorkRecordID == "" {
		if err := requireWorkRecordCreateAccess(txApp, auth, targetUID, "create this work record"); err != nil {
			return WorkRecordSaveResponse{}, err
		}
	} else {
		source, err := findWorkRecordCopySourceByToken(txApp, copyFromWorkRecordID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return WorkRecordSaveResponse{}, workRecordRouteError(http.StatusNotFound, "work_record_not_found", "work record not found")
			}
			return WorkRecordSaveResponse{}, err
		}
		if err := requireCopySourceAccess(txApp, auth, source); err != nil {
			return WorkRecordSaveResponse{}, err
		}
		if targetUID != source.UID {
			return WorkRecordSaveResponse{}, &errs.HookError{
				Status:  http.StatusBadRequest,
				Message: "validation failed",
				Data: map[string]errs.CodeError{
					"uid": {
						Code:    "copy_uid_mismatch",
						Message: "Copied work records must keep the same worker as the source record",
					},
				},
			}
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

	record := core.NewRecord(newWorkRecordPersistenceCollection())
	record.Set("uid", targetUID)
	record.Set("creator", auth.Id)
	record.Set("parent_work_record", rootParentID)

	if err := applyWorkRecordSaveRequest(record, toWorkRecordSaveRequest(req)); err != nil {
		return WorkRecordSaveResponse{}, err
	}
	if err := normalizeWorkRecordValidationError(apphooks.ValidateWorkRecord(txApp, record)); err != nil {
		return WorkRecordSaveResponse{}, err
	}
	if err := normalizeWorkRecordValidationError(apphooks.EnsureWorkRecordNumber(txApp, record)); err != nil {
		return WorkRecordSaveResponse{}, err
	}
	if err := normalizeWorkRecordValidationError(txApp.Save(record)); err != nil {
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

		entryRecord := core.NewRecord(newWorkRecordConsumableEntryPersistenceCollection())
		entryRecord.Set("work_record", record.Id)
		entryRecord.Set("consumable", item.ConsumableID)
		setWorkRecordConsumableQuantity(entryRecord, inputKindsByConsumable[item.ConsumableID], item)

		if err := normalizeWorkRecordValidationError(apphooks.ValidateWorkRecordConsumableEntry(txApp, entryRecord)); err != nil {
			return WorkRecordSaveResponse{}, err
		}
		if err := normalizeWorkRecordValidationError(txApp.Save(entryRecord)); err != nil {
			return WorkRecordSaveResponse{}, err
		}
	}

	return WorkRecordSaveResponse{ID: record.Id}, nil
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

		workRecordToken := strings.TrimSpace(e.Request.PathValue("workRecord"))
		if workRecordToken == "" {
			return writeHookError(e, workRecordRouteError(http.StatusBadRequest, "work_record_required", "workRecord is required"))
		}

		var req WorkRecordSaveRequest
		if err := e.BindBody(&req); err != nil {
			return writeHookError(e, workRecordRouteError(http.StatusBadRequest, "invalid_request_body", "invalid request body"))
		}

		record, err := findWorkRecordByToken(requestApp, workRecordToken)
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

func createGetWorkRecordCopyDraftHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		auth := e.Auth
		if auth == nil {
			return e.Error(http.StatusUnauthorized, "unauthorized", nil)
		}

		workRecordToken := strings.TrimSpace(e.Request.PathValue("workRecord"))
		if workRecordToken == "" {
			return writeHookError(e, workRecordRouteError(http.StatusBadRequest, "work_record_required", "workRecord is required"))
		}

		response, err := loadWorkRecordCopyDraft(app, auth, workRecordToken)
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
				wr.id AS id,
				wr.number AS work_record
			FROM work_records wr
			WHERE wr.uid = {:uid}
			  AND wr.job = {:job}
			  AND wr.date = {:date}
			  AND NOT EXISTS (
					SELECT 1
					FROM time_entries te
					WHERE te.work_record_id = wr.id
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

		workRecordToken := strings.TrimSpace(e.Request.PathValue("workRecord"))
		if workRecordToken == "" {
			pathParts := strings.Split(strings.Trim(e.Request.URL.Path, "/"), "/")
			if len(pathParts) > 0 {
				workRecordToken = strings.TrimSpace(pathParts[len(pathParts)-1])
			}
		}
		if workRecordToken == "" {
			return e.Error(http.StatusBadRequest, "workRecord is required", nil)
		}

		var row workRecordDetailRow
		err := app.DB().NewQuery(`
			SELECT
				wr.id AS id,
				wr.number AS work_record,
				wr.date AS date,
				wr.uid AS uid,
				COALESCE(p.given_name || ' ' || p.surname, '') AS owner_name,
				wr.creator AS creator,
				COALESCE(wr.job, '') AS job_id,
				COALESCE(j.number, '') AS job_number,
				COALESCE(wr.type, '') AS type_id,
				COALESCE(wrt.name, '') AS type_name,
				COALESCE(wr.location, '') AS location,
				COALESCE(wr.work_description, '') AS work_description,
				COALESCE(wr.report_to, '') AS report_to,
				COALESCE(wr.hours_on_site, 0) AS hours_on_site,
				COALESCE(wr.hours_travel_time, 0) AS hours_travel_time,
				COALESCE(wr.distance_travelled_km, 0) AS distance_travelled_km,
				COALESCE(wr.field_book_number, '') AS field_book_number,
				COALESCE(wr.field_book_page_number, 0) AS field_book_page_number,
				COALESCE(wr.sub_contractor, '') AS sub_contractor,
				COALESCE(wr.equipment, '') AS equipment,
				COALESCE(wr.supplies, '') AS supplies,
				COALESCE(wr.vehicle_type, '') AS vehicle_type,
				COALESCE(wr.is_passenger, 0) AS is_passenger,
				COALESCE(wr.company_vehicle_unit_number, '') AS company_vehicle_unit_number,
				COALESCE(wr.approved, 0) AS approved,
				COALESCE(wr.parent_work_record, '') AS parent_work_record_id,
				COALESCE(parent_wr.number, '') AS parent_work_record_number
			FROM work_records wr
			LEFT JOIN profiles p
			  ON p.uid = wr.uid
			LEFT JOIN jobs j
			  ON wr.job = j.id
			LEFT JOIN work_records_types wrt
			  ON wr.type = wrt.id
			LEFT JOIN work_records parent_wr
			  ON parent_wr.id = wr.parent_work_record
			WHERE wr.id = {:token} OR wr.number = {:token}
			LIMIT 1
		`).Bind(dbx.Params{
			"token": workRecordToken,
		}).One(&row)
		if err != nil {
			return e.Error(http.StatusNotFound, "work record not found", err)
		}

		if row.UID != auth.Id && row.Creator != auth.Id {
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

		linkedEntry, err := loadLinkedTimeEntry(app, row.ID)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load linked time entry", err)
		}
		notes, err := loadWorkRecordNotes(app, row.ID)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load work record notes", err)
		}
		consumables, err := loadWorkRecordConsumables(app, row.ID, row.TypeID)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load work record consumables", err)
		}

		canCopyToTomorrow, err := canCreateWorkRecordForUID(app, auth, row.UID)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to check copy permissions", err)
		}
		canMutateRecord := (row.UID == auth.Id || row.Creator == auth.Id) && !row.Approved

		return e.JSON(http.StatusOK, WorkRecordDetailResponse{
			ID:                         row.ID,
			WorkRecord:                 row.WorkRecord,
			Date:                       row.Date,
			UID:                        row.UID,
			OwnerName:                  row.OwnerName,
			Creator:                    row.Creator,
			JobID:                      row.JobID,
			JobNumber:                  row.JobNumber,
			TypeID:                     row.TypeID,
			TypeName:                   row.TypeName,
			Location:                   row.Location,
			WorkDescription:            row.WorkDescription,
			ReportTo:                   row.ReportTo,
			HoursOnSite:                row.HoursOnSite,
			HoursTravelTime:            row.HoursTravelTime,
			DistanceTravelledKM:        row.DistanceTravelledKM,
			FieldBookNumber:            row.FieldBookNumber,
			FieldBookPageNumber:        row.FieldBookPageNumber,
			SubContractor:              row.SubContractor,
			Equipment:                  row.Equipment,
			Supplies:                   row.Supplies,
			VehicleType:                row.VehicleType,
			IsPassenger:                row.IsPassenger,
			CompanyVehicleUnitNumber:   row.CompanyVehicleUnitNumber,
			Approved:                   row.Approved,
			ParentWorkRecordID:         row.ParentWorkRecordID,
			ParentWorkRecord:           row.ParentWorkRecordNumber,
			RootWorkRecordID:           firstNonEmpty(row.ParentWorkRecordID, row.ID),
			RootWorkRecord:             firstNonEmpty(row.ParentWorkRecordNumber, row.WorkRecord),
			LinkedTimeEntryID:          linkedEntry.LinkedTimeEntryID,
			LinkedTimeSheetID:          linkedEntry.LinkedTimeSheetID,
			LinkedTimeEntryDate:        linkedEntry.LinkedTimeEntryDate,
			LinkedTimeEntryDescription: linkedEntry.LinkedTimeEntryDescription,
			LinkedTimeSheetWeekEnding:  linkedEntry.LinkedTimeSheetWeekEnding,
			CanCreateTimeEntry:         row.UID == auth.Id && linkedEntry.LinkedTimeEntryID == "",
			CanCopyToTomorrow:          canCopyToTomorrow,
			CanEdit:                    canMutateRecord,
			CanAddNote:                 row.UID == auth.Id || row.Creator == auth.Id,
			Notes:                      notes,
			Consumables:                consumables,
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
