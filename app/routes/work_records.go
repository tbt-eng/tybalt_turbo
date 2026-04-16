package routes

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"
	"tybalt/errs"
	"tybalt/numbering"
	"tybalt/utilities"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/security"
)

type WorkRecordListRow struct {
	ID                         string `db:"id" json:"id"`
	WorkRecord                 string `db:"work_record" json:"work_record"`
	Date                       string `db:"date" json:"date"`
	UID                        string `db:"uid" json:"uid"`
	Creator                    string `db:"creator" json:"creator"`
	JobID                      string `db:"job_id" json:"job_id"`
	JobNumber                  string `db:"job_number" json:"job_number"`
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

type WorkRecordDetailResponse struct {
	ID                         string  `json:"id"`
	WorkRecord                 string  `json:"work_record"`
	Date                       string  `json:"date"`
	UID                        string  `json:"uid"`
	Creator                    string  `json:"creator"`
	JobID                      string  `json:"job_id"`
	JobNumber                  string  `json:"job_number"`
	TypeID                     string  `json:"type_id"`
	TypeName                   string  `json:"type_name"`
	Location                   string  `json:"location"`
	WorkDescription            string  `json:"work_description"`
	HoursOnSite                float64 `json:"hours_on_site"`
	HoursTravelTime            float64 `json:"hours_travel_time"`
	Approved                   bool    `json:"approved"`
	LinkedTimeEntryID          string  `json:"linked_time_entry_id"`
	LinkedTimeSheetID          string  `json:"linked_timesheet_id"`
	LinkedTimeEntryDate        string  `json:"linked_time_entry_date"`
	LinkedTimeEntryDescription string  `json:"linked_time_entry_description"`
	LinkedTimeSheetWeekEnding  string  `json:"linked_time_sheet_week_ending"`
	CanCreateTimeEntry         bool    `json:"can_create_time_entry"`
	CanCopyToTomorrow          bool    `json:"can_copy_to_tomorrow"`
}

type WorkRecordCopyResponse struct {
	ID               string `json:"id"`
	WorkRecord       string `json:"work_record"`
	Date             string `json:"date"`
	ParentWorkRecord string `json:"parent_work_record"`
}

type WorkRecordOption struct {
	ID         string `db:"id" json:"id"`
	WorkRecord string `db:"work_record" json:"work_record"`
}

type workRecordDetailRow struct {
	ID              string  `db:"id"`
	WorkRecord      string  `db:"work_record"`
	Date            string  `db:"date"`
	UID             string  `db:"uid"`
	Creator         string  `db:"creator"`
	JobID           string  `db:"job_id"`
	JobNumber       string  `db:"job_number"`
	TypeID          string  `db:"type_id"`
	TypeName        string  `db:"type_name"`
	Location        string  `db:"location"`
	WorkDescription string  `db:"work_description"`
	HoursOnSite     float64 `db:"hours_on_site"`
	HoursTravelTime float64 `db:"hours_travel_time"`
	Approved        bool    `db:"approved"`
}

type workRecordCopySourceRow struct {
	ID                       string  `db:"id"`
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
	ParentWorkRecord         string  `db:"parent_work_record"`
}

func canViewAllWorkRecords(app core.App, auth *core.Record) (bool, error) {
	if auth == nil {
		return false, nil
	}

	return utilities.HasClaim(app, auth, "report")
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

func findWorkRecordByToken(app core.App, token string) (*core.Record, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, sql.ErrNoRows
	}

	var row struct {
		ID string `db:"id"`
	}
	if err := app.DB().NewQuery(`
		SELECT id
		FROM work_records
		WHERE id = {:token} OR number = {:token}
		LIMIT 1
	`).Bind(dbx.Params{
		"token": token,
	}).One(&row); err != nil {
		return nil, err
	}

	return app.FindRecordById("work_records", row.ID)
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
			COALESCE(wr.parent_work_record, '') AS parent_work_record
		FROM work_records wr
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

func createCopyWorkRecordToTomorrowHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		auth := e.Auth
		if auth == nil {
			return e.Error(http.StatusUnauthorized, "unauthorized", nil)
		}

		workRecordToken := strings.TrimSpace(e.Request.PathValue("workRecord"))
		if workRecordToken == "" {
			return e.Error(http.StatusBadRequest, "workRecord is required", nil)
		}

		var (
			httpResponseStatusCode int
			response               WorkRecordCopyResponse
		)

		err := app.RunInTransaction(func(txApp core.App) error {
			source, err := findWorkRecordCopySourceByToken(txApp, workRecordToken)
			if err != nil {
				httpResponseStatusCode = http.StatusNotFound
				return &CodeError{
					Code:    "work_record_not_found",
					Message: fmt.Sprintf("work record not found: %v", err),
				}
			}

			if source.UID != auth.Id && source.Creator != auth.Id {
				canViewAll, err := canViewAllWorkRecords(txApp, auth)
				if err != nil {
					httpResponseStatusCode = http.StatusInternalServerError
					return err
				}
				if !canViewAll {
					httpResponseStatusCode = http.StatusForbidden
					return &CodeError{
						Code:    "unauthorized",
						Message: "you are not authorized to view this work record",
					}
				}
			}

			canCreate, err := canCreateWorkRecordForUID(txApp, auth, source.UID)
			if err != nil {
				httpResponseStatusCode = http.StatusInternalServerError
				return err
			}
			if !canCreate {
				httpResponseStatusCode = http.StatusForbidden
				return &CodeError{
					Code:    "unauthorized",
					Message: "you are not authorized to copy this work record",
				}
			}

			sourceDate, err := time.Parse("2006-01-02", source.Date)
			if err != nil {
				httpResponseStatusCode = http.StatusInternalServerError
				return &CodeError{
					Code:    "invalid_source_date",
					Message: fmt.Sprintf("work record has invalid date: %v", err),
				}
			}
			targetDateTime := sourceDate.AddDate(0, 0, 1)
			targetDate := targetDateTime.Format("2006-01-02")

			var duplicate struct {
				ID string `db:"id"`
			}
			err = txApp.DB().NewQuery(`
				SELECT id
				FROM work_records
				WHERE uid = {:uid}
				  AND job = {:job}
				  AND date = {:date}
				LIMIT 1
			`).Bind(dbx.Params{
				"uid":  source.UID,
				"job":  source.Job,
				"date": targetDate,
			}).One(&duplicate)
			if err == nil && duplicate.ID != "" {
				httpResponseStatusCode = http.StatusBadRequest
				return &CodeError{
					Code:    "work_record_for_date_exists",
					Message: "a work record for tomorrow already exists for this worker and job",
				}
			}
			if err != nil && err != sql.ErrNoRows {
				httpResponseStatusCode = http.StatusInternalServerError
				return err
			}

			rootParentID := strings.TrimSpace(source.ParentWorkRecord)
			if rootParentID == "" {
				rootParentID = source.ID
			}

			number, err := numbering.ReserveWorkRecordNumber(txApp, targetDateTime)
			if err != nil {
				httpResponseStatusCode = http.StatusInternalServerError
				return err
			}

			now := time.Now().UTC().Format("2006-01-02 15:04:05.000Z")
			newID := security.PseudorandomString(15)
			if _, err := txApp.DB().NewQuery(`
				INSERT INTO work_records (
					id, number, date, uid, creator, job, type, location,
					work_description, report_to, hours_on_site, hours_travel_time,
					distance_travelled_km, field_book_number, field_book_page_number,
					sub_contractor, equipment, supplies, vehicle_type, is_passenger,
					company_vehicle_unit_number, approved, created, updated, parent_work_record
				) VALUES (
					{:id}, {:number}, {:date}, {:uid}, {:creator}, {:job}, {:type}, {:location},
					{:workDescription}, {:reportTo}, {:hoursOnSite}, {:hoursTravelTime},
					{:distanceTravelledKM}, {:fieldBookNumber}, {:fieldBookPageNumber},
					{:subContractor}, {:equipment}, {:supplies}, {:vehicleType}, {:isPassenger},
					{:companyVehicleUnitNumber}, 0, {:created}, {:updated}, {:parentWorkRecord}
				)
			`).Bind(dbx.Params{
				"id":                       newID,
				"number":                   number,
				"date":                     targetDate,
				"uid":                      source.UID,
				"creator":                  auth.Id,
				"job":                      source.Job,
				"type":                     source.Type,
				"location":                 source.Location,
				"workDescription":          source.WorkDescription,
				"reportTo":                 source.ReportTo,
				"hoursOnSite":              source.HoursOnSite,
				"hoursTravelTime":          source.HoursTravelTime,
				"distanceTravelledKM":      source.DistanceTravelledKM,
				"fieldBookNumber":          source.FieldBookNumber,
				"fieldBookPageNumber":      source.FieldBookPageNumber,
				"subContractor":            source.SubContractor,
				"equipment":                source.Equipment,
				"supplies":                 source.Supplies,
				"vehicleType":              source.VehicleType,
				"isPassenger":              source.IsPassenger,
				"companyVehicleUnitNumber": source.CompanyVehicleUnitNumber,
				"created":                  now,
				"updated":                  now,
				"parentWorkRecord":         rootParentID,
			}).Execute(); err != nil {
				httpResponseStatusCode = http.StatusInternalServerError
				return err
			}

			response = WorkRecordCopyResponse{
				ID:               newID,
				WorkRecord:       number,
				Date:             targetDate,
				ParentWorkRecord: rootParentID,
			}

			return nil
		})

		if err != nil {
			switch typed := err.(type) {
			case *CodeError:
				return e.JSON(httpResponseStatusCode, map[string]any{
					"error": typed.Message,
					"code":  typed.Code,
				})
			default:
				return e.JSON(httpResponseStatusCode, map[string]string{"error": err.Error()})
			}
		}

		return e.JSON(http.StatusOK, response)
	}
}

func createGetWorkRecordsHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		auth := e.Auth
		if auth == nil {
			return e.Error(http.StatusUnauthorized, "unauthorized", nil)
		}

		canViewAll, err := canViewAllWorkRecords(app, auth)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to check work record permissions", err)
		}

		var rows []WorkRecordListRow
		if err := app.DB().NewQuery(`
			SELECT
				wr.id AS id,
				wr.number AS work_record,
				wr.date AS date,
				wr.uid AS uid,
				wr.creator AS creator,
				COALESCE(wr.job, '') AS job_id,
				COALESCE(j.number, '') AS job_number,
				COALESCE(wrt.name, '') AS type_name,
				COALESCE(wr.location, '') AS location,
				COALESCE(wr.work_description, '') AS work_description,
				COALESCE(wr.approved, 0) AS approved,
				COALESCE(te.id, '') AS linked_time_entry_id,
				COALESCE(te.tsid, '') AS linked_timesheet_id,
				COALESCE(te.date, '') AS linked_time_entry_date,
				COALESCE(te.description, '') AS linked_time_entry_description
			FROM work_records wr
			LEFT JOIN jobs j
			  ON wr.job = j.id
			LEFT JOIN work_records_types wrt
			  ON wr.type = wrt.id
			LEFT JOIN time_entries te
			  ON te.work_record_id = wr.id
			WHERE {:canViewAll} = 1
			   OR wr.uid = {:authId}
			   OR wr.creator = {:authId}
			ORDER BY wr.date DESC, wr.number DESC
		`).Bind(dbx.Params{
			"canViewAll": boolToSQLInt(canViewAll),
			"authId":     auth.Id,
		}).All(&rows); err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load work records", err)
		}

		for i := range rows {
			rows[i].SearchText = buildWorkRecordSearchText(rows[i].WorkRecord + " " + rows[i].JobNumber + " " + rows[i].TypeName + " " + rows[i].Location + " " + rows[i].WorkDescription)
			rows[i].CanCreateTimeEntry = rows[i].UID == auth.Id && rows[i].LinkedTimeEntryID == ""
		}

		return e.JSON(http.StatusOK, rows)
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
				wr.creator AS creator,
				COALESCE(wr.job, '') AS job_id,
				COALESCE(j.number, '') AS job_number,
				COALESCE(wr.type, '') AS type_id,
				COALESCE(wrt.name, '') AS type_name,
				COALESCE(wr.location, '') AS location,
				COALESCE(wr.work_description, '') AS work_description,
				COALESCE(wr.hours_on_site, 0) AS hours_on_site,
				COALESCE(wr.hours_travel_time, 0) AS hours_travel_time,
				COALESCE(wr.approved, 0) AS approved
			FROM work_records wr
			LEFT JOIN jobs j
			  ON wr.job = j.id
			LEFT JOIN work_records_types wrt
			  ON wr.type = wrt.id
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

		canCopyToTomorrow, err := canCreateWorkRecordForUID(app, auth, row.UID)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to check copy permissions", err)
		}

		return e.JSON(http.StatusOK, WorkRecordDetailResponse{
			ID:                         row.ID,
			WorkRecord:                 row.WorkRecord,
			Date:                       row.Date,
			UID:                        row.UID,
			Creator:                    row.Creator,
			JobID:                      row.JobID,
			JobNumber:                  row.JobNumber,
			TypeID:                     row.TypeID,
			TypeName:                   row.TypeName,
			Location:                   row.Location,
			WorkDescription:            row.WorkDescription,
			HoursOnSite:                row.HoursOnSite,
			HoursTravelTime:            row.HoursTravelTime,
			Approved:                   row.Approved,
			LinkedTimeEntryID:          linkedEntry.LinkedTimeEntryID,
			LinkedTimeSheetID:          linkedEntry.LinkedTimeSheetID,
			LinkedTimeEntryDate:        linkedEntry.LinkedTimeEntryDate,
			LinkedTimeEntryDescription: linkedEntry.LinkedTimeEntryDescription,
			LinkedTimeSheetWeekEnding:  linkedEntry.LinkedTimeSheetWeekEnding,
			CanCreateTimeEntry:         row.UID == auth.Id && linkedEntry.LinkedTimeEntryID == "",
			CanCopyToTomorrow:          canCopyToTomorrow,
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

func buildWorkRecordSearchText(workRecord string) string {
	workRecord = strings.ToUpper(strings.TrimSpace(workRecord))
	if workRecord == "" {
		return ""
	}

	return strings.Join(strings.Fields(strings.ReplaceAll(workRecord, "-", " ")), " ")
}
