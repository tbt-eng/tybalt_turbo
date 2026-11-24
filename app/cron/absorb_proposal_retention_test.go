package cron

import (
	"testing"
	"time"

	"tybalt/internal/testutils"
)

// The fixture rows of absorb_proposals.csv and absorb_proposal_records.csv.
// absorbprop00004 was declined at 2026-11-19 18:00 (one record row) and
// absorbprop00005 withdrawn at 2026-11-21 10:00 (two record rows). The other
// three are pending. The tests pass a fixed now, so they add no rows.
const (
	retentionFixtureProposals       = 5
	retentionFixtureProposalRecords = 6
)

func TestDeleteExpiredAbsorbProposals(t *testing.T) {
	tests := []struct {
		name        string
		now         time.Time
		gone        []string
		goneRecords []string
		kept        []string
	}{
		{
			name: "nothing is old enough",
			// Cutoff 2026-11-19 18:00:00.000, the decided time of absorbprop00004.
			now:  time.Date(2026, 12, 19, 18, 0, 0, 0, time.UTC),
			kept: []string{"absorbprop00004", "absorbprop00005"},
		},
		{
			name:        "a proposal decided a millisecond before the cutoff goes",
			now:         time.Date(2026, 12, 19, 18, 0, 0, int(time.Millisecond), time.UTC),
			gone:        []string{"absorbprop00004"},
			goneRecords: []string{"absorbproprec04"},
			kept:        []string{"absorbprop00005"},
		},
		{
			name: "every decided proposal goes and the pending ones stay",
			// A year later. The pending proposals are older still, and stay.
			now:         time.Date(2027, 12, 1, 0, 0, 0, 0, time.UTC),
			gone:        []string{"absorbprop00004", "absorbprop00005"},
			goneRecords: []string{"absorbproprec04", "absorbproprec05", "absorbproprec06"},
			kept:        []string{"absorbprop00001", "absorbprop00002", "absorbprop00003"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := testutils.SetupTestApp(t)
			defer app.Cleanup()

			deleted, err := deleteExpiredAbsorbProposals(app, tt.now)
			if err != nil {
				t.Fatalf("deleteExpiredAbsorbProposals: %v", err)
			}
			if deleted != int64(len(tt.gone)) {
				t.Errorf("deleted %d proposals, want %d", deleted, len(tt.gone))
			}

			assertRowCount(t, app, "absorb_proposals", retentionFixtureProposals-len(tt.gone))
			assertRowCount(t, app, "absorb_proposal_records", retentionFixtureProposalRecords-len(tt.goneRecords))
			for _, id := range tt.gone {
				assertRowExists(t, app, "absorb_proposals", id, false)
			}
			for _, id := range tt.goneRecords {
				assertRowExists(t, app, "absorb_proposal_records", id, false)
			}
			for _, id := range tt.kept {
				assertRowExists(t, app, "absorb_proposals", id, true)
			}
		})
	}
}
