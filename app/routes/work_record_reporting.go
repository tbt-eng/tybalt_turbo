package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"tybalt/errs"
	"tybalt/utilities"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type WorkRecordSearchResponse struct {
	Items []WorkRecordListRow `json:"items"`
	Total int                 `json:"total"`
}

type PaginatedJobWorkRecordsResponse struct {
	Data       []WorkRecordListRow `json:"data"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	Total      int                 `json:"total"`
	TotalPages int                 `json:"total_pages"`
}

type WorkRecordConsumablesReportResponse struct {
	Items []WorkRecordConsumableReportRow `json:"items"`
}

type WorkRecordConsumableReportRow struct {
	WorkRecordID    string  `db:"work_record_id" json:"work_record_id"`
	WorkRecord      string  `db:"work_record" json:"work_record"`
	Date            string  `db:"date" json:"date"`
	UID             string  `db:"uid" json:"uid"`
	OwnerName       string  `db:"owner_name" json:"owner_name"`
	JobID           string  `db:"job_id" json:"job_id"`
	JobNumber       string  `db:"job_number" json:"job_number"`
	TypeID          string  `db:"type_id" json:"type_id"`
	TypeName        string  `db:"type_name" json:"type_name"`
	ConsumableID    string  `db:"consumable_id" json:"consumable_id"`
	ConsumableName  string  `db:"consumable_name" json:"consumable_name"`
	InputKind       string  `db:"input_kind" json:"input_kind"`
	UnitLabel       string  `db:"unit_label" json:"unit_label"`
	QuantityNumber  float64 `db:"quantity_number" json:"quantity_number"`
	BooleanSelected bool    `db:"boolean_selected" json:"boolean_selected"`
}

type workRecordSearchFilters struct {
	UID      string
	JobID    string
	TypeID   string
	DateFrom string
	DateTo   string
	Linked   *bool
	Approved *bool
}

func workRecordSearchFieldError(field string, code string, message string) *errs.HookError {
	return &errs.HookError{
		Status:  http.StatusBadRequest,
		Message: "Validation error",
		Data: map[string]errs.CodeError{
			field: {
				Code:    code,
				Message: message,
			},
		},
	}
}

func requireWorkRecordReportViewer(app core.App, auth *core.Record) error {
	if auth == nil {
		return &errs.HookError{
			Status:  http.StatusUnauthorized,
			Message: "authentication is required",
			Data: map[string]errs.CodeError{
				"global": {Code: "authentication_required", Message: "authentication is required"},
			},
		}
	}

	hasReport, err := utilities.HasClaim(app, auth, "report")
	if err != nil {
		return err
	}
	if hasReport {
		return nil
	}

	return &errs.HookError{
		Status:  http.StatusForbidden,
		Message: "you are not authorized to view work record reporting",
		Data: map[string]errs.CodeError{
			"global": {Code: "unauthorized", Message: "you are not authorized to view work record reporting"},
		},
	}
}

func parseWorkRecordSearchFilters(e *core.RequestEvent) (workRecordSearchFilters, error) {
	query := e.Request.URL.Query()
	filters := workRecordSearchFilters{
		UID:      strings.TrimSpace(query.Get("uid")),
		JobID:    strings.TrimSpace(query.Get("job")),
		TypeID:   strings.TrimSpace(query.Get("type")),
		DateFrom: strings.TrimSpace(query.Get("date_from")),
		DateTo:   strings.TrimSpace(query.Get("date_to")),
	}

	if filters.DateFrom != "" {
		if _, err := time.Parse("2006-01-02", filters.DateFrom); err != nil {
			return filters, workRecordSearchFieldError("date_from", "invalid_date", "date_from must use YYYY-MM-DD")
		}
	}
	if filters.DateTo != "" {
		if _, err := time.Parse("2006-01-02", filters.DateTo); err != nil {
			return filters, workRecordSearchFieldError("date_to", "invalid_date", "date_to must use YYYY-MM-DD")
		}
	}

	parseOptionalBool := func(name string) (*bool, error) {
		value := strings.TrimSpace(query.Get(name))
		if value == "" {
			return nil, nil
		}
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return nil, workRecordSearchFieldError(name, "invalid_boolean", fmt.Sprintf("%s must be true or false", name))
		}
		return &parsed, nil
	}

	var err error
	filters.Linked, err = parseOptionalBool("linked")
	if err != nil {
		return filters, err
	}
	filters.Approved, err = parseOptionalBool("approved")
	if err != nil {
		return filters, err
	}

	return filters, nil
}

func buildWorkRecordFilterSQL(alias string, filters workRecordSearchFilters) (string, dbx.Params) {
	conditions := []string{"1 = 1"}
	params := dbx.Params{}

	if filters.UID != "" {
		conditions = append(conditions, "EXISTS (SELECT 1 FROM work_records_subjects wrs_uid WHERE wrs_uid.work_record = "+alias+".id AND wrs_uid.uid = {:uid})")
		params["uid"] = filters.UID
	}
	if filters.JobID != "" {
		conditions = append(conditions, alias+".job = {:job}")
		params["job"] = filters.JobID
	}
	if filters.TypeID != "" {
		conditions = append(conditions, alias+".type = {:type}")
		params["type"] = filters.TypeID
	}
	if filters.DateFrom != "" {
		conditions = append(conditions, alias+".date >= {:dateFrom}")
		params["dateFrom"] = filters.DateFrom
	}
	if filters.DateTo != "" {
		conditions = append(conditions, alias+".date <= {:dateTo}")
		params["dateTo"] = filters.DateTo
	}
	if filters.Linked != nil {
		if *filters.Linked {
			conditions = append(conditions, "EXISTS (SELECT 1 FROM work_records_subjects wrs_filter JOIN time_entries te_filter ON te_filter.work_record_subject_id = wrs_filter.id WHERE wrs_filter.work_record = "+alias+".id)")
		} else {
			conditions = append(conditions, "EXISTS (SELECT 1 FROM work_records_subjects wrs_filter WHERE wrs_filter.work_record = "+alias+".id AND NOT EXISTS (SELECT 1 FROM time_entries te_filter WHERE te_filter.work_record_subject_id = wrs_filter.id))")
		}
	}
	if filters.Approved != nil {
		if *filters.Approved {
			conditions = append(conditions, "EXISTS (SELECT 1 FROM work_records_subjects wrs_approved WHERE wrs_approved.work_record = "+alias+".id AND COALESCE(wrs_approved.approved, 0) = 1)")
		} else {
			conditions = append(conditions, "EXISTS (SELECT 1 FROM work_records_subjects wrs_approved WHERE wrs_approved.work_record = "+alias+".id AND COALESCE(wrs_approved.approved, 0) = 0)")
		}
	}

	return strings.Join(conditions, " AND "), params
}

func loadFilteredWorkRecordRows(app core.App, filters workRecordSearchFilters, limit int, offset int) ([]WorkRecordListRow, int, error) {
	filterSQL, params := buildWorkRecordFilterSQL("wr", filters)
	params["limit"] = limit
	params["offset"] = offset

	var totalRow struct {
		Count int `db:"count"`
	}
	if err := app.DB().NewQuery(`
		SELECT COUNT(*) AS count
		FROM work_records wr
		WHERE ` + filterSQL).Bind(params).One(&totalRow); err != nil {
		return nil, 0, err
	}
	params["displayUid"] = filters.UID

	rows := []WorkRecordListRow{}
	if err := app.DB().NewQuery(`
		SELECT
			wr.id AS id,
			wr.number AS work_record,
			wr.date AS date,
			COALESCE(display_subject.uid, '') AS uid,
			COALESCE(p.given_name || ' ' || p.surname, '') AS owner_name,
			wr.creator AS creator,
			COALESCE(wr.job, '') AS job_id,
			COALESCE(j.number, '') AS job_number,
			COALESCE(wr.type, '') AS type_id,
			COALESCE(wrt.name, '') AS type_name,
			COALESCE(wr.location, '') AS location,
			COALESCE(wr.work_description, '') AS work_description,
			CASE WHEN COALESCE(subject_counts.approved_subject_count, 0) > 0 THEN 1 ELSE 0 END AS approved,
			COALESCE(te.id, '') AS linked_time_entry_id,
			COALESCE(te.tsid, '') AS linked_timesheet_id,
			COALESCE(te.date, '') AS linked_time_entry_date,
			COALESCE(te.description, '') AS linked_time_entry_description,
			COALESCE(subject_counts.subject_count, 0) AS subject_count,
			COALESCE(subject_counts.approved_subject_count, 0) AS approved_subject_count
		FROM work_records wr
		LEFT JOIN (
			SELECT work_record, MIN(id) AS subject_id
			FROM work_records_subjects
			WHERE {:displayUid} = '' OR uid = {:displayUid}
			GROUP BY work_record
		) display_subject_id
		  ON display_subject_id.work_record = wr.id
		LEFT JOIN work_records_subjects display_subject
		  ON display_subject.id = display_subject_id.subject_id
		LEFT JOIN profiles p
		  ON p.uid = display_subject.uid
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
		LEFT JOIN time_entries te
		  ON te.work_record_subject_id = display_subject.id
		WHERE ` + filterSQL + `
		ORDER BY wr.date DESC, wr.number DESC
		LIMIT {:limit} OFFSET {:offset}
	`).Bind(params).All(&rows); err != nil {
		return nil, 0, err
	}

	for i := range rows {
		rows[i].SearchText = buildWorkRecordSearchText(rows[i].WorkRecord + " " + rows[i].OwnerName + " " + rows[i].JobNumber + " " + rows[i].TypeName + " " + rows[i].Location + " " + rows[i].WorkDescription)
		rows[i].CanCreateTimeEntry = false
	}

	return rows, totalRow.Count, nil
}

func loadAllFilteredWorkRecordRows(app core.App, filters workRecordSearchFilters) ([]WorkRecordListRow, int, error) {
	filterSQL, params := buildWorkRecordFilterSQL("wr", filters)
	params["displayUid"] = filters.UID

	rows := []WorkRecordListRow{}
	if err := app.DB().NewQuery(`
		SELECT
			wr.id AS id,
			wr.number AS work_record,
			wr.date AS date,
			COALESCE(display_subject.uid, '') AS uid,
			COALESCE(p.given_name || ' ' || p.surname, '') AS owner_name,
			wr.creator AS creator,
			COALESCE(wr.job, '') AS job_id,
			COALESCE(j.number, '') AS job_number,
			COALESCE(wr.type, '') AS type_id,
			COALESCE(wrt.name, '') AS type_name,
			COALESCE(wr.location, '') AS location,
			COALESCE(wr.work_description, '') AS work_description,
			CASE WHEN COALESCE(subject_counts.approved_subject_count, 0) > 0 THEN 1 ELSE 0 END AS approved,
			COALESCE(te.id, '') AS linked_time_entry_id,
			COALESCE(te.tsid, '') AS linked_timesheet_id,
			COALESCE(te.date, '') AS linked_time_entry_date,
			COALESCE(te.description, '') AS linked_time_entry_description,
			COALESCE(subject_counts.subject_count, 0) AS subject_count,
			COALESCE(subject_counts.approved_subject_count, 0) AS approved_subject_count
		FROM work_records wr
		LEFT JOIN (
			SELECT work_record, MIN(id) AS subject_id
			FROM work_records_subjects
			WHERE {:displayUid} = '' OR uid = {:displayUid}
			GROUP BY work_record
		) display_subject_id
		  ON display_subject_id.work_record = wr.id
		LEFT JOIN work_records_subjects display_subject
		  ON display_subject.id = display_subject_id.subject_id
		LEFT JOIN profiles p
		  ON p.uid = display_subject.uid
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
		LEFT JOIN time_entries te
		  ON te.work_record_subject_id = display_subject.id
		WHERE ` + filterSQL + `
		ORDER BY wr.date DESC, wr.number DESC
	`).Bind(params).All(&rows); err != nil {
		return nil, 0, err
	}

	for i := range rows {
		rows[i].SearchText = buildWorkRecordSearchText(rows[i].WorkRecord + " " + rows[i].OwnerName + " " + rows[i].JobNumber + " " + rows[i].TypeName + " " + rows[i].Location + " " + rows[i].WorkDescription)
		rows[i].CanCreateTimeEntry = false
	}

	return rows, len(rows), nil
}

func createGetWorkRecordSearchHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		if err := requireWorkRecordReportViewer(app, e.Auth); err != nil {
			return writeHookError(e, err)
		}

		filters, err := parseWorkRecordSearchFilters(e)
		if err != nil {
			return writeHookError(e, err)
		}

		rows, total, err := loadAllFilteredWorkRecordRows(app, filters)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load work record search results", err)
		}

		return e.JSON(http.StatusOK, WorkRecordSearchResponse{
			Items: rows,
			Total: total,
		})
	}
}

func createGetJobWorkRecordsHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		if e.Auth == nil {
			return writeHookError(e, &errs.HookError{
				Status:  http.StatusUnauthorized,
				Message: "authentication is required",
				Data: map[string]errs.CodeError{
					"global": {Code: "authentication_required", Message: "authentication is required"},
				},
			})
		}

		jobID := strings.TrimSpace(e.Request.PathValue("id"))
		if jobID == "" {
			return e.Error(http.StatusBadRequest, "id is required", nil)
		}

		page := 1
		if raw := strings.TrimSpace(e.Request.URL.Query().Get("page")); raw != "" {
			parsed, err := strconv.Atoi(raw)
			if err != nil || parsed <= 0 {
				return writeHookError(e, workRecordSearchFieldError("page", "invalid_page", "page must be a positive integer"))
			}
			page = parsed
		}

		limit := 50
		if raw := strings.TrimSpace(e.Request.URL.Query().Get("limit")); raw != "" {
			parsed, err := strconv.Atoi(raw)
			if err != nil || parsed <= 0 {
				return writeHookError(e, workRecordSearchFieldError("limit", "invalid_limit", "limit must be a positive integer"))
			}
			if parsed > 200 {
				parsed = 200
			}
			limit = parsed
		}

		offset := (page - 1) * limit
		rows, total, err := loadFilteredWorkRecordRows(app, workRecordSearchFilters{
			JobID:    jobID,
			Approved: func() *bool { v := true; return &v }(),
		}, limit, offset)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load job work records", err)
		}

		totalPages := 1
		if total > 0 {
			totalPages = (total + limit - 1) / limit
		}

		return e.JSON(http.StatusOK, PaginatedJobWorkRecordsResponse{
			Data:       rows,
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		})
	}
}

func createGetUnlinkedWorkRecordsReportHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		if err := requireWorkRecordReportViewer(app, e.Auth); err != nil {
			return writeHookError(e, err)
		}
		filters, err := parseWorkRecordSearchFilters(e)
		if err != nil {
			return writeHookError(e, err)
		}
		linked := false
		filters.Linked = &linked
		rows, total, err := loadAllFilteredWorkRecordRows(app, filters)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load unlinked work records", err)
		}
		return e.JSON(http.StatusOK, WorkRecordSearchResponse{Items: rows, Total: total})
	}
}

func createGetPendingApprovalWorkRecordsReportHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		if err := requireWorkRecordReportViewer(app, e.Auth); err != nil {
			return writeHookError(e, err)
		}
		filters, err := parseWorkRecordSearchFilters(e)
		if err != nil {
			return writeHookError(e, err)
		}
		approved := false
		filters.Approved = &approved
		rows, total, err := loadAllFilteredWorkRecordRows(app, filters)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load pending-approval work records", err)
		}
		return e.JSON(http.StatusOK, WorkRecordSearchResponse{Items: rows, Total: total})
	}
}

func createGetWorkRecordConsumablesReportHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		if err := requireWorkRecordReportViewer(app, e.Auth); err != nil {
			return writeHookError(e, err)
		}

		filters, err := parseWorkRecordSearchFilters(e)
		if err != nil {
			return writeHookError(e, err)
		}

		filterSQL, params := buildWorkRecordFilterSQL("wr", filters)
		params["displayUid"] = filters.UID

		rows := []WorkRecordConsumableReportRow{}
		if err := app.DB().NewQuery(`
			SELECT
				wr.id AS work_record_id,
				wr.number AS work_record,
				wr.date AS date,
				COALESCE(display_subject.uid, '') AS uid,
				COALESCE(p.given_name || ' ' || p.surname, '') AS owner_name,
				COALESCE(wr.job, '') AS job_id,
				COALESCE(j.number, '') AS job_number,
				COALESCE(wr.type, '') AS type_id,
				COALESCE(wrt.name, '') AS type_name,
				COALESCE(wrc.id, '') AS consumable_id,
				COALESCE(wrc.name, '') AS consumable_name,
				COALESCE(wrc.input_kind, '') AS input_kind,
				COALESCE(wrc.unit_label, '') AS unit_label,
				COALESCE(wrce.quantity_number, 0) AS quantity_number,
				CASE WHEN wrce.id != '' THEN 1 ELSE 0 END AS boolean_selected
			FROM work_records wr
			JOIN work_records_consumable_entries wrce
			  ON wrce.work_record = wr.id
			JOIN work_records_consumables wrc
			  ON wrc.id = wrce.consumable
			LEFT JOIN (
				SELECT work_record, MIN(id) AS subject_id
				FROM work_records_subjects
				WHERE {:displayUid} = '' OR uid = {:displayUid}
				GROUP BY work_record
			) display_subject_id
			  ON display_subject_id.work_record = wr.id
			LEFT JOIN work_records_subjects display_subject
			  ON display_subject.id = display_subject_id.subject_id
			LEFT JOIN profiles p
			  ON p.uid = display_subject.uid
			LEFT JOIN jobs j
			  ON j.id = wr.job
			LEFT JOIN work_records_types wrt
			  ON wrt.id = wr.type
			WHERE ` + filterSQL + `
			ORDER BY wr.date DESC, wr.number DESC, wrc.name ASC
		`).Bind(params).All(&rows); err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load work record consumables report", err)
		}

		return e.JSON(http.StatusOK, WorkRecordConsumablesReportResponse{Items: rows})
	}
}
