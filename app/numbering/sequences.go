package numbering

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type sequenceValueRow struct {
	Value int `db:"value"`
}

var errSequenceExhausted = errors.New("sequence exhausted")

// reserveSequence atomically initializes and increments a named counter row.
//
// The caller is expected to pass a transactional app when the reserved value
// must commit together with some related record save.
func reserveSequence(app core.App, scope string, initialValueQuery string, params dbx.Params, maxValue int) (int, error) {
	if maxValue < 1 {
		return 0, fmt.Errorf("invalid max sequence value %d", maxValue)
	}

	insertParams := cloneParams(params)
	insertParams["scope"] = scope

	if _, err := app.DB().NewQuery(`
		INSERT INTO number_sequences (scope, value)
		VALUES ({:scope}, COALESCE((` + initialValueQuery + `), 0))
		ON CONFLICT(scope) DO NOTHING
	`).Bind(insertParams).Execute(); err != nil {
		return 0, fmt.Errorf("initialize sequence %s: %w", scope, err)
	}

	updateParams := cloneParams(insertParams)
	updateParams["maxValue"] = maxValue

	var row sequenceValueRow
	if err := app.DB().NewQuery(`
		UPDATE number_sequences
		SET value = value + 1
		WHERE scope = {:scope}
		  AND value < {:maxValue}
		RETURNING value
	`).Bind(updateParams).One(&row); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("%w: %s", errSequenceExhausted, scope)
		}
		return 0, fmt.Errorf("reserve next value for sequence %s: %w", scope, err)
	}

	return row.Value, nil
}

func ReserveWorkRecordNumber(app core.App, now time.Time) (string, error) {
	prefix := fmt.Sprintf("W%02d%02d-", now.Year()%100, int(now.Month()))

	value, err := reserveSequence(
		app,
		"work_records:monthly:"+prefix,
		`SELECT MAX(CAST(SUBSTR(number, LENGTH({:prefix}) + 1, 4) AS INTEGER))
		FROM work_records
		WHERE number GLOB {:glob}`,
		dbx.Params{
			"prefix": prefix,
			"glob":   prefix + "[0-9][0-9][0-9][0-9]",
		},
		9999,
	)
	if err != nil {
		if errors.Is(err, errSequenceExhausted) {
			return "", fmt.Errorf("unable to generate a unique work record number")
		}
		return "", err
	}

	return fmt.Sprintf("%s%04d", prefix, value), nil
}

func ReserveParentPONumber(app core.App, now time.Time) (string, error) {
	prefix := fmt.Sprintf("%02d%02d-", now.Year()%100, int(now.Month()))
	upperBound := fmt.Sprintf("%s%04d", prefix, 5000)

	value, err := reserveSequence(
		app,
		"purchase_orders:monthly_parent:"+prefix,
		`SELECT MAX(CAST(
			CASE
				WHEN INSTR(SUBSTR(po_number, LENGTH({:prefix}) + 1), '-') > 0
					THEN SUBSTR(
						SUBSTR(po_number, LENGTH({:prefix}) + 1),
						1,
						INSTR(SUBSTR(po_number, LENGTH({:prefix}) + 1), '-') - 1
					)
				ELSE SUBSTR(po_number, LENGTH({:prefix}) + 1)
			END
		AS INTEGER))
		FROM purchase_orders
		WHERE parent_po = ''
		  AND po_number LIKE {:like}
		  AND po_number < {:upperBound}`,
		dbx.Params{
			"prefix":     prefix,
			"like":       prefix + "%",
			"upperBound": upperBound,
		},
		4999,
	)
	if err != nil {
		if errors.Is(err, errSequenceExhausted) {
			return "", fmt.Errorf("unable to generate a unique PO number")
		}
		return "", err
	}

	return fmt.Sprintf("%s%04d", prefix, value), nil
}

func cloneParams(params dbx.Params) dbx.Params {
	cloned := dbx.Params{}
	for k, v := range params {
		cloned[k] = v
	}
	return cloned
}
