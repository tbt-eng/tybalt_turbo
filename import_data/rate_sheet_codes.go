package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"time"

	"imports/load"
)

// Rate sheet codes (issue 207). Every rate sheet in Turbo has a code such as
// 2026-11-00004: a month and a five digit serial that starts at 00001 in
// each month. All revisions of a sheet share the code. Turbo's hooks issue
// codes, but the import writes SQLite directly, so it gives imported sheets
// their codes itself:
//
//   - A sheet that already had a code keeps it. --clear_rates deletes the
//     rate tables first, so the codes are read before the delete and matched
//     to the reimported rows by id, which the export keeps the same.
//   - A sheet without one takes the code of its other revisions (the legacy
//     system links revisions by name), or else a new code in the month of the
//     import.
//   - rate_sheet_codes records every code issued and is never cleared, so a
//     new code is one more than the highest issued in the month, and no code
//     is issued twice.

var importRateSheetCodePattern = regexp.MustCompile(`^(\d{4}-(0[1-9]|1[0-2]))-(\d{5})$`)

// rateSheetCodeSnapshot is the codes of the rate sheets in the database
// before an import clears them.
type rateSheetCodeSnapshot struct {
	byID   map[string]string
	byName map[string]string
}

// snapshotRateSheetCodes reads the code of every rate sheet. It returns an
// empty snapshot for a database whose rate_sheets has no code column yet.
func snapshotRateSheetCodes(dbPath string) (rateSheetCodeSnapshot, error) {
	snapshot := rateSheetCodeSnapshot{byID: map[string]string{}, byName: map[string]string{}}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return snapshot, err
	}
	defer db.Close()
	var hasCode int
	if err := db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('rate_sheets') WHERE name = 'code'`).Scan(&hasCode); err != nil {
		return snapshot, err
	}
	if hasCode == 0 {
		return snapshot, nil
	}
	rows, err := db.Query(`SELECT id, COALESCE(name, ''), COALESCE(code, '') FROM rate_sheets ORDER BY revision`)
	if err != nil {
		return snapshot, err
	}
	defer rows.Close()
	for rows.Next() {
		var id, name, code string
		if err := rows.Scan(&id, &name, &code); err != nil {
			return snapshot, err
		}
		if code == "" {
			continue
		}
		snapshot.byID[id] = code
		if _, ok := snapshot.byName[name]; !ok && name != "" {
			snapshot.byName[name] = code
		}
	}
	return snapshot, rows.Err()
}

// issuedRateSheetSerial returns the highest serial issued in month, from
// rate_sheet_codes and from the codes in snapshot.
func issuedRateSheetSerial(dbPath string, month string, snapshot rateSheetCodeSnapshot) (int, error) {
	highest := 0
	for _, code := range snapshot.byID {
		if m, serial, ok := parseImportRateSheetCode(code); ok && m == month && serial > highest {
			highest = serial
		}
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return 0, err
	}
	defer db.Close()
	var issued int
	if err := db.QueryRow(`SELECT COALESCE(MAX(serial), 0) FROM rate_sheet_codes WHERE month = ?`, month).Scan(&issued); err != nil {
		return 0, err
	}
	return max(highest, issued), nil
}

func parseImportRateSheetCode(code string) (string, int, bool) {
	match := importRateSheetCodePattern.FindStringSubmatch(code)
	if match == nil {
		return "", 0, false
	}
	serial, err := strconv.Atoi(match[3])
	if err != nil || serial < 1 {
		return "", 0, false
	}
	return match[1], serial, true
}

// planImportedRateSheetCodes gives each imported sheet its code. A sheet
// keeps the code it had (by id). Otherwise it takes the code of a revision
// with the same name, from this import or from before it. A sheet of a new
// name gets the next code of month, after issued. New codes are given in
// order of name, then id, so that a run gives the same result.
func planImportedRateSheetCodes(items []load.RateSheet, snapshot rateSheetCodeSnapshot, month string, issued int) (map[string]string, error) {
	codes := map[string]string{}
	byName := map[string]string{}
	for _, item := range items {
		if code, ok := snapshot.byID[item.Id]; ok {
			codes[item.Id] = code
			if _, taken := byName[item.Name]; !taken {
				byName[item.Name] = code
			}
		}
	}
	rest := []load.RateSheet{}
	for _, item := range items {
		if _, ok := codes[item.Id]; !ok {
			rest = append(rest, item)
		}
	}
	sort.SliceStable(rest, func(i, j int) bool {
		if rest[i].Name != rest[j].Name {
			return rest[i].Name < rest[j].Name
		}
		return rest[i].Id < rest[j].Id
	})
	serial := issued
	for _, item := range rest {
		code, ok := byName[item.Name]
		if !ok && item.Name != "" {
			code, ok = snapshot.byName[item.Name]
		}
		if !ok {
			serial++
			if serial > 99999 {
				return nil, fmt.Errorf("no rate sheet code is left in %s", month)
			}
			code = fmt.Sprintf("%s-%05d", month, serial)
		}
		codes[item.Id] = code
		if item.Name != "" {
			byName[item.Name] = code
		}
	}
	return codes, nil
}

// planRateSheetCodesForImport reads the rate sheets of the import and plans
// their codes in the month of now (UTC).
func planRateSheetCodesForImport(dbPath string, items []load.RateSheet, snapshot rateSheetCodeSnapshot, now time.Time) (map[string]string, error) {
	month := now.UTC().Format("2006-01")
	issued, err := issuedRateSheetSerial(dbPath, month, snapshot)
	if err != nil {
		return nil, err
	}
	return planImportedRateSheetCodes(items, snapshot, month, issued)
}

// recordIssuedRateSheetCodes adds each code to rate_sheet_codes, unless it is
// there already.
func recordIssuedRateSheetCodes(dbPath string, codes map[string]string) error {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	seen := map[string]bool{}
	created := time.Now().UTC().Format("2006-01-02 15:04:05.000Z")
	for _, code := range codes {
		if seen[code] {
			continue
		}
		seen[code] = true
		month, serial, ok := parseImportRateSheetCode(code)
		if !ok {
			return fmt.Errorf("invalid rate sheet code %q", code)
		}
		var exists int
		if err := tx.QueryRow(`SELECT COUNT(*) FROM rate_sheet_codes WHERE code = ?`, code).Scan(&exists); err != nil {
			return err
		}
		if exists > 0 {
			continue
		}
		if _, err := tx.Exec(`INSERT INTO rate_sheet_codes (id, code, month, serial, created) VALUES (?, ?, ?, ?, ?)`,
			randomRecordID(), code, month, serial, created); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// randomRecordID returns a PocketBase style id: 15 characters from a-z and
// 0-9.
func randomRecordID() string {
	const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 15)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	for i := range b {
		b[i] = alphabet[int(b[i])%len(alphabet)]
	}
	return string(b)
}
