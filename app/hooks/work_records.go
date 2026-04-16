package hooks

import (
	"fmt"
	"strings"
	"time"
	"tybalt/numbering"

	validation "github.com/pocketbase/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
)

func GenerateWorkRecordNumber(app core.App) (string, error) {
	return numbering.ReserveWorkRecordNumber(app, time.Now())
}

func ProcessWorkRecord(app core.App, e *core.RecordEvent) error {
	record := e.Record
	// Status quo:
	// - the collection createRule already forbids API clients from submitting a
	//   non-empty number;
	// - internal/non-request callers may pre-seed a number and this hook
	//   preserves it, which leaves room for controlled import/backfill flows.
	//
	// Future consideration:
	// - if we later want all non-import creates to be strictly backend-numbered,
	//   tighten this branch to require an explicit import/backfill signal.
	if strings.TrimSpace(record.GetString("number")) != "" {
		return nil
	}

	number, err := GenerateWorkRecordNumber(app)
	if err != nil {
		return validation.Errors{
			"number": validation.NewError(
				"error_generating_work_record_number",
				fmt.Sprintf("error generating work record number: %v", err),
			),
		}
	}

	record.Set("number", number)

	return nil
}
