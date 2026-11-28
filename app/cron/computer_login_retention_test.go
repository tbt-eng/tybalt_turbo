package cron

import (
	"reflect"
	"testing"
	"time"

	"tybalt/internal/testutils"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// The fixture rows of computer_logins.csv and computer_raw_logins.csv. The
// tests pass a fixed now, so they use these rows and add none.
//
// The raw logins: cmprawlogin0001 to 0003 and 0005 to 0007 have a computer
// name and never expire. cmprawlogin0004 (2026-10-18 13:00) and
// cmprawlogin0008 (2026-10-20 07:00) have none and expire after 30 days.
const (
	retentionFixtureLogins            = 26
	retentionFixtureRawLogins         = 8
	retentionFixtureNamelessRawLogins = 2
)

// The raw logins with a computer name, which the job must never delete.
var namedFixtureRawLogins = []string{
	"cmprawlogin0001", // FAKE-PC-0001, 2026-10-15 13:00
	"cmprawlogin0002", // FAKE-PC-0099, 2026-10-16 13:00
	"cmprawlogin0003", // FAKE-PC-0098, 2026-10-17 13:00
	"cmprawlogin0005", // FAKE-PC-0097, 2026-10-19 13:00
	"cmprawlogin0006", // FAKE-PC-0097, 2026-10-21 13:00
	"cmprawlogin0007", // FAKE-PC-0097, 2026-10-21 13:00
}

func TestDeleteExpiredComputerLogins(t *testing.T) {
	tests := []struct {
		name string
		now  time.Time
		// The rows that must be gone. Every other fixture row must stay.
		goneLogins    []string
		goneRawLogins []string
		// Rows that must stay, checked by id as well as by count.
		keptLogins    []string
		keptRawLogins []string
	}{
		{
			name: "nothing is old enough",
			// Cutoff 2026-10-14 13:00, a day before the oldest raw login.
			now:           time.Date(2026, 11, 13, 13, 0, 0, 0, time.UTC),
			keptLogins:    []string{"cmplogin0000026"},
			keptRawLogins: []string{"cmprawlogin0001"},
		},
		{
			name: "a login created at the cutoff stays",
			// Cutoff 2026-10-20 13:00:00.000, the created time of cmplogin0000026.
			now:           time.Date(2026, 11, 19, 13, 0, 0, 0, time.UTC),
			goneRawLogins: []string{"cmprawlogin0004", "cmprawlogin0008"},
			keptLogins:    []string{"cmplogin0000026"},
			keptRawLogins: namedFixtureRawLogins,
		},
		{
			name:          "a login created a millisecond before the cutoff goes",
			now:           time.Date(2026, 11, 19, 13, 0, 0, int(time.Millisecond), time.UTC),
			goneLogins:    []string{"cmplogin0000026"},
			goneRawLogins: []string{"cmprawlogin0004", "cmprawlogin0008"},
			keptLogins:    []string{"cmplogin0000001"},
			keptRawLogins: namedFixtureRawLogins,
		},
		{
			name: "a raw login without a name created at the cutoff stays",
			// Cutoff 2026-10-20 07:00:00.000, the created time of cmprawlogin0008.
			// The named raw logins are older still and stay too.
			now:           time.Date(2026, 11, 19, 7, 0, 0, 0, time.UTC),
			goneRawLogins: []string{"cmprawlogin0004"},
			keptLogins:    []string{"cmplogin0000026"},
			keptRawLogins: append([]string{"cmprawlogin0008"}, namedFixtureRawLogins...),
		},
		{
			name: "old logins go and newer ones stay",
			// Cutoff 2026-11-06 13:00, the created time of cmplogin0000007.
			now: time.Date(2026, 12, 6, 13, 0, 0, 0, time.UTC),
			goneLogins: []string{
				"cmplogin0000026", // 2026-10-20 13:00
				"cmplogin0000001", // 2026-11-01 13:00
				"cmplogin0000002", // 2026-11-05 07:00
				"cmplogin0000003", // 2026-11-05 13:00
				"cmplogin0000004", // 2026-11-05 19:00
				"cmplogin0000005", // 2026-11-06 01:00
				"cmplogin0000006", // 2026-11-06 07:00
			},
			goneRawLogins: []string{"cmprawlogin0004", "cmprawlogin0008"},
			keptLogins:    []string{"cmplogin0000007", "cmplogin0000024", "cmplogin0000025"},
			keptRawLogins: namedFixtureRawLogins,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := testutils.SetupTestApp(t)
			defer app.Cleanup()

			logins, rawLogins, err := deleteExpiredComputerLogins(app, tt.now)
			if err != nil {
				t.Fatalf("deleteExpiredComputerLogins: %v", err)
			}
			if logins != int64(len(tt.goneLogins)) {
				t.Errorf("deleted %d logins, want %d", logins, len(tt.goneLogins))
			}
			if rawLogins != int64(len(tt.goneRawLogins)) {
				t.Errorf("deleted %d raw logins, want %d", rawLogins, len(tt.goneRawLogins))
			}

			assertRowCount(t, app, "computer_logins", retentionFixtureLogins-len(tt.goneLogins))
			assertRowCount(t, app, "computer_raw_logins", retentionFixtureRawLogins-len(tt.goneRawLogins))
			for _, id := range tt.goneLogins {
				assertRowExists(t, app, "computer_logins", id, false)
			}
			for _, id := range tt.goneRawLogins {
				assertRowExists(t, app, "computer_raw_logins", id, false)
			}
			for _, id := range tt.keptLogins {
				assertRowExists(t, app, "computer_logins", id, true)
			}
			for _, id := range tt.keptRawLogins {
				assertRowExists(t, app, "computer_raw_logins", id, true)
			}
		})
	}
}

func TestDeleteExpiredComputerLogins_LeavesComputersUntouched(t *testing.T) {
	app := testutils.SetupTestApp(t)
	defer app.Cleanup()

	before := computerRows(t, app)
	if len(before) == 0 {
		t.Fatal("expected fixture computers")
	}

	// Far enough ahead that every fixture login and raw login is old. Every
	// login and every raw login without a name goes; the named raw logins stay.
	now := time.Date(2027, 6, 1, 0, 0, 0, 0, time.UTC)
	logins, rawLogins, err := deleteExpiredComputerLogins(app, now)
	if err != nil {
		t.Fatalf("deleteExpiredComputerLogins: %v", err)
	}
	if logins != retentionFixtureLogins || rawLogins != retentionFixtureNamelessRawLogins {
		t.Fatalf("deleted %d logins and %d raw logins, want %d and %d",
			logins, rawLogins, retentionFixtureLogins, retentionFixtureNamelessRawLogins)
	}
	assertRowCount(t, app, "computer_raw_logins", len(namedFixtureRawLogins))
	for _, id := range namedFixtureRawLogins {
		assertRowExists(t, app, "computer_raw_logins", id, true)
	}

	if after := computerRows(t, app); !reflect.DeepEqual(before, after) {
		t.Fatalf("computers changed:\nbefore %v\nafter  %v", before, after)
	}

	// A second run finds nothing more to delete.
	logins, rawLogins, err = deleteExpiredComputerLogins(app, now)
	if err != nil {
		t.Fatalf("second deleteExpiredComputerLogins: %v", err)
	}
	if logins != 0 || rawLogins != 0 {
		t.Fatalf("second run deleted %d logins and %d raw logins, want 0 and 0", logins, rawLogins)
	}
}

func assertRowCount(t *testing.T, app core.App, table string, want int) {
	t.Helper()

	var got int
	if err := app.DB().NewQuery("SELECT COUNT(*) FROM " + table).Row(&got); err != nil {
		t.Fatalf("count %s: %v", table, err)
	}
	if got != want {
		t.Errorf("%s has %d rows, want %d", table, got, want)
	}
}

func assertRowExists(t *testing.T, app core.App, table string, id string, want bool) {
	t.Helper()

	var got int
	if err := app.DB().NewQuery("SELECT COUNT(*) FROM " + table + " WHERE id = {:id}").
		Bind(dbx.Params{"id": id}).Row(&got); err != nil {
		t.Fatalf("find %s %s: %v", table, id, err)
	}
	if (got == 1) != want {
		t.Errorf("%s %s exists = %v, want %v", table, id, got == 1, want)
	}
}

// computerRows returns every column of every computer, so a test can check
// that nothing in computers changed.
func computerRows(t *testing.T, app core.App) []dbx.NullStringMap {
	t.Helper()

	var rows []dbx.NullStringMap
	if err := app.DB().NewQuery("SELECT * FROM computers ORDER BY id").All(&rows); err != nil {
		t.Fatalf("read computers: %v", err)
	}
	return rows
}
