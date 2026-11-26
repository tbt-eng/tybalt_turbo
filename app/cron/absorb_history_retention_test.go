package cron

import (
	"testing"
	"time"

	"tybalt/internal/testutils"
)

// The fixture rows of absorb_history.csv were committed at 2026-08-20 10:00
// (absorbhist00001), 2026-11-20 09:00 (absorbhist00002) and 2026-11-22 09:00
// (absorbhist00003). The tests pass a fixed now, so they add no rows.
func TestDeleteExpiredAbsorbHistory(t *testing.T) {
	tests := []struct {
		name string
		now  time.Time
		gone []string
		kept []string
	}{
		{
			name: "nothing is old enough",
			// Cutoff 2026-08-20 10:00:00.000, the commit time of absorbhist00001.
			now:  time.Date(2026, 11, 18, 10, 0, 0, 0, time.UTC),
			kept: []string{"absorbhist00001", "absorbhist00002", "absorbhist00003"},
		},
		{
			name: "a row committed a millisecond before the cutoff goes",
			now:  time.Date(2026, 11, 18, 10, 0, 0, int(time.Millisecond), time.UTC),
			gone: []string{"absorbhist00001"},
			kept: []string{"absorbhist00002", "absorbhist00003"},
		},
		{
			name: "every row goes a year later",
			now:  time.Date(2027, 12, 1, 0, 0, 0, 0, time.UTC),
			gone: []string{"absorbhist00001", "absorbhist00002", "absorbhist00003"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := testutils.SetupTestApp(t)
			defer app.Cleanup()

			deleted, err := deleteExpiredAbsorbHistory(app, tt.now)
			if err != nil {
				t.Fatalf("deleteExpiredAbsorbHistory: %v", err)
			}
			if deleted != int64(len(tt.gone)) {
				t.Errorf("deleted %d rows, want %d", deleted, len(tt.gone))
			}
			assertRowCount(t, app, "absorb_history", len(tt.kept))
			for _, id := range tt.gone {
				assertRowExists(t, app, "absorb_history", id, false)
			}
			for _, id := range tt.kept {
				assertRowExists(t, app, "absorb_history", id, true)
			}
		})
	}
}
