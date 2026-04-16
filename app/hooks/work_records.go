// Package hooks owns PocketBase record hooks that enforce backend-side business
// rules during record persistence.
//
// This file exists to keep work record numbering server-assigned and consistent
// regardless of which UI or internal caller creates the record. Work records
// are intended to use a generated W-number sequence, and the transition in this
// branch makes that number part of the user-visible workflow because linked time
// entries and compatibility reads resolve through the first-class work_records
// table.
//
// The hook deliberately does one small job: if a work record arrives without a
// number, it reserves the next valid W-number and writes it before save. If the
// number is already present, the hook preserves it so controlled import/backfill
// flows remain possible. Keeping the logic here means the numbering policy is
// enforced even when callers bypass the normal UI and talk directly to the
// backend.
package hooks

import (
	"fmt"
	"strings"
	"time"
	"tybalt/numbering"

	validation "github.com/pocketbase/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
)

// GenerateWorkRecordNumber reserves the next work-record W-number using the
// shared numbering package.
func GenerateWorkRecordNumber(app core.App) (string, error) {
	return numbering.ReserveWorkRecordNumber(app, time.Now())
}

// ProcessWorkRecord applies work-record-specific save behavior before the
// record is persisted.
//
// Today that means assigning a generated number when the caller did not provide
// one. The hook does not overwrite an existing number because some internal
// import and repair flows may need to preserve seeded values explicitly.
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
