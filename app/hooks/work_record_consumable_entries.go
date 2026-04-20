package hooks

import (
	"database/sql"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type workRecordConsumableContext struct {
	Approved         bool   `db:"approved"`
	WorkRecordType   string `db:"work_record_type"`
	InputKind        string `db:"input_kind"`
	ConsumableActive bool   `db:"consumable_active"`
}

func loadWorkRecordConsumableContext(app core.App, workRecordID string, consumableID string) (*workRecordConsumableContext, error) {
	var row workRecordConsumableContext
	if err := app.DB().NewQuery(`
		SELECT
			COALESCE(wr.approved, 0) AS approved,
			COALESCE(wr.type, '') AS work_record_type,
			COALESCE(wrc.input_kind, '') AS input_kind,
			COALESCE(wrc.active, 0) AS consumable_active
		FROM work_records wr
		JOIN work_records_consumables wrc
		  ON wrc.id = {:consumableId}
		WHERE wr.id = {:workRecordId}
		LIMIT 1
	`).Bind(dbx.Params{
		"workRecordId": workRecordID,
		"consumableId": consumableID,
	}).One(&row); err != nil {
		return nil, err
	}
	return &row, nil
}

func validateWorkRecordConsumableAllowlist(app core.App, workRecordType string, consumableID string) error {
	var row struct {
		Count int `db:"count"`
	}
	if err := app.DB().NewQuery(`
		SELECT COUNT(*) AS count
		FROM work_records_type_consumables
		WHERE type = {:type}
		  AND consumable = {:consumable}
		  AND COALESCE(active, 0) = 1
	`).Bind(dbx.Params{
		"type":       workRecordType,
		"consumable": consumableID,
	}).One(&row); err != nil {
		return err
	}
	if row.Count == 0 {
		return validationFieldError(
			"consumable",
			"consumable_not_allowed_for_type",
			"This consumable is not allowed for the selected work record type",
		)
	}
	return nil
}

func validateWorkRecordConsumableUniqueness(app core.App, record *core.Record) error {
	var row struct {
		ID string `db:"id"`
	}
	err := app.DB().NewQuery(`
		SELECT id
		FROM work_records_consumable_entries
		WHERE work_record = {:workRecordId}
		  AND consumable = {:consumableId}
		  AND id != {:id}
		LIMIT 1
	`).Bind(dbx.Params{
		"workRecordId": record.GetString("work_record"),
		"consumableId": record.GetString("consumable"),
		"id":           record.Id,
	}).One(&row)
	if err == nil && row.ID != "" {
		return validationFieldError(
			"consumable",
			"duplicate_consumable_entry",
			"This consumable already has an entry for the work record",
		)
	}
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

func ValidateWorkRecordConsumableEntry(app core.App, record *core.Record) error {
	original := record.Original()
	if original != nil && original.Id != "" {
		if original.GetString("work_record") != record.GetString("work_record") {
			return validationFieldError(
				"work_record",
				"immutable_field",
				"Work record cannot be changed for an existing consumable entry",
			)
		}
		if original.GetString("consumable") != record.GetString("consumable") {
			return validationFieldError(
				"consumable",
				"immutable_field",
				"Consumable cannot be changed for an existing consumable entry",
			)
		}
	}

	workRecordID := strings.TrimSpace(record.GetString("work_record"))
	consumableID := strings.TrimSpace(record.GetString("consumable"))
	if workRecordID == "" {
		return validationFieldError("work_record", "required", "Work record is required")
	}
	if consumableID == "" {
		return validationFieldError("consumable", "required", "Consumable is required")
	}

	ctx, err := loadWorkRecordConsumableContext(app, workRecordID, consumableID)
	if err != nil {
		return validationFieldError(
			"work_record",
			"invalid_consumable_context",
			"Referenced work record or consumable was not found",
		)
	}

	if ctx.Approved {
		return validationFieldError(
			"global",
			"approved_work_record_locked",
			"Consumables cannot be edited on an approved work record",
		)
	}

	if !ctx.ConsumableActive {
		return validationFieldError(
			"consumable",
			"inactive_consumable",
			"Consumable is inactive and cannot be used",
		)
	}

	if err := validateWorkRecordConsumableAllowlist(app, ctx.WorkRecordType, consumableID); err != nil {
		return err
	}

	if err := validateWorkRecordConsumableUniqueness(app, record); err != nil {
		return err
	}

	switch ctx.InputKind {
	case "number":
		if record.GetFloat("quantity_number") <= 0 {
			return validationFieldError(
				"quantity_number",
				"quantity_required",
				"Numeric consumables require a quantity greater than 0",
			)
		}
	case "boolean":
		if record.GetFloat("quantity_number") != 0 {
			return validationFieldError(
				"quantity_number",
				"boolean_consumable_quantity_forbidden",
				"Boolean consumables must not store a numeric quantity",
			)
		}
	default:
		return validationFieldError(
			"consumable",
			"invalid_input_kind",
			"Consumable input kind is invalid",
		)
	}

	return nil
}

func ValidateWorkRecordConsumableEntryDelete(app core.App, record *core.Record) error {
	workRecordID := strings.TrimSpace(record.GetString("work_record"))
	consumableID := strings.TrimSpace(record.GetString("consumable"))
	ctx, err := loadWorkRecordConsumableContext(app, workRecordID, consumableID)
	if err != nil {
		return validationFieldError(
			"work_record",
			"invalid_consumable_context",
			"Referenced work record or consumable was not found",
		)
	}
	if ctx.Approved {
		return validationFieldError(
			"global",
			"approved_work_record_locked",
			"Consumables cannot be edited on an approved work record",
		)
	}
	return nil
}
