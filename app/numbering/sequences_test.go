// Package numbering tests the shared backend sequence reservation logic used by
// work records and purchase orders.
//
// This file exists because numbering bugs are easy to miss in normal request
// tests but expensive in production: duplicate numbers, stale seeds after
// manual inserts, exhaustion handling, or races under concurrent writers can all
// corrupt user-visible identifiers. The tests here focus on those failure modes
// directly so the sequence layer can evolve without silently breaking numbering
// guarantees higher up the stack.
//
// The test strategy uses seeded PocketBase apps so reservations run against
// realistic fixture data, then exercises both single-threaded and concurrent
// paths. That gives confidence that the number_sequences-backed implementation
// is correct both for legacy-seeded data and for new transactional reservations.
package numbering

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"testing"
	"time"
	"tybalt/internal/testseed"

	"github.com/pocketbase/pocketbase/core"
)

func TestReserveSequence_ExhaustionReturnsSentinel(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	t.Cleanup(app.Cleanup)

	_, err := reserveSequence(app, "test:sequence", "SELECT 0", nil, 1)
	if err == nil {
		t.Fatal("expected exhaustion error, got nil")
	}
	if !errors.Is(err, errSequenceExhausted) {
		t.Fatalf("expected errSequenceExhausted, got %v", err)
	}
}

func TestReserveWorkRecordNumber_SeedCorrectness(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	t.Cleanup(app.Cleanup)

	when := time.Date(2026, time.April, 16, 0, 0, 0, 0, time.UTC)
	got, err := ReserveWorkRecordNumber(app, when)
	if err != nil {
		t.Fatalf("unexpected error reserving work record number: %v", err)
	}
	if got != "W2604-0011" {
		t.Fatalf("expected W2604-0011, got %s", got)
	}
}

func TestReserveWorkRecordNumber_Exhaustion(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	t.Cleanup(app.Cleanup)

	when := time.Date(2030, time.March, 16, 0, 0, 0, 0, time.UTC)

	_, err := ReserveWorkRecordNumber(app, when)
	if err == nil {
		t.Fatal("expected exhaustion error, got nil")
	}
	if err.Error() != "unable to generate a unique work record number" {
		t.Fatalf("expected public exhaustion error, got %q", err.Error())
	}
}

func TestReserveWorkRecordNumber_ConcurrentReservations(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	t.Cleanup(app.Cleanup)

	when := time.Date(2030, time.January, 16, 0, 0, 0, 0, time.UTC)
	numbers := reserveConcurrently(t, app, 20, func(txApp core.App) (string, error) {
		return ReserveWorkRecordNumber(txApp, when)
	})

	assertContiguousSequence(t, numbers, "W3001-", 20)
}

func TestReserveWorkRecordNumber_ManualInsertAfterSequenceSeed(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	t.Cleanup(app.Cleanup)

	when := time.Date(2030, time.February, 16, 0, 0, 0, 0, time.UTC)

	got, err := ReserveWorkRecordNumber(app, when)
	if err != nil {
		t.Fatalf("unexpected error reserving work record number: %v", err)
	}
	if got != "W3002-0011" {
		t.Fatalf("expected W3002-0011 after stale manual insert, got %s", got)
	}
}

func TestReserveParentPONumber_SeedCorrectnessAndCompatibility(t *testing.T) {
	tests := []struct {
		name string
		when time.Time
		want string
	}{
		{
			name: "existing parent max seeds next sequential value",
			when: time.Date(2024, time.January, 16, 0, 0, 0, 0, time.UTC),
			want: "2401-0010",
		},
		{
			name: "manual 5000 range is ignored on initial seed",
			when: time.Date(2024, time.February, 16, 0, 0, 0, 0, time.UTC),
			want: "2402-0001",
		},
		{
			name: "dashed historical parent max is parsed defensively",
			when: time.Date(2024, time.April, 16, 0, 0, 0, 0, time.UTC),
			want: "2404-0011",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := testseed.NewSeededTestApp(t)
			t.Cleanup(app.Cleanup)

			got, err := ReserveParentPONumber(app, tt.when)
			if err != nil {
				t.Fatalf("unexpected error reserving parent PO number: %v", err)
			}
			if got != tt.want {
				t.Fatalf("expected %s, got %s", tt.want, got)
			}
		})
	}
}

func TestReserveParentPONumber_Exhaustion(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	t.Cleanup(app.Cleanup)

	when := time.Date(2024, time.June, 16, 0, 0, 0, 0, time.UTC)

	_, err := ReserveParentPONumber(app, when)
	if err == nil {
		t.Fatal("expected exhaustion error, got nil")
	}
	if err.Error() != "unable to generate a unique PO number" {
		t.Fatalf("expected public exhaustion error, got %q", err.Error())
	}
}

func TestReserveParentPONumber_ConcurrentReservations(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	t.Cleanup(app.Cleanup)

	when := time.Date(2030, time.April, 16, 0, 0, 0, 0, time.UTC)
	numbers := reserveConcurrently(t, app, 20, func(txApp core.App) (string, error) {
		return ReserveParentPONumber(txApp, when)
	})

	assertContiguousSequence(t, numbers, "3004-", 20)
}

func TestReserveParentPONumber_ManualInsertAfterSequenceSeed(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	t.Cleanup(app.Cleanup)

	when := time.Date(2030, time.May, 16, 0, 0, 0, 0, time.UTC)

	got, err := ReserveParentPONumber(app, when)
	if err != nil {
		t.Fatalf("unexpected error reserving parent PO number: %v", err)
	}
	if got != "3005-0011" {
		t.Fatalf("expected 3005-0011 after stale manual insert, got %s", got)
	}
}

func reserveConcurrently(t *testing.T, app core.App, count int, reserve func(core.App) (string, error)) []string {
	t.Helper()

	results := make([]string, 0, count)
	errCh := make(chan error, count)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < count; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			var number string
			err := app.RunInTransaction(func(txApp core.App) error {
				var err error
				number, err = reserve(txApp)
				return err
			})
			if err != nil {
				errCh <- err
				return
			}

			mu.Lock()
			results = append(results, number)
			mu.Unlock()
		}()
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Fatalf("unexpected concurrent reservation error: %v", err)
	}

	return results
}

func assertContiguousSequence(t *testing.T, numbers []string, prefix string, count int) {
	t.Helper()

	if len(numbers) != count {
		t.Fatalf("expected %d numbers, got %d", count, len(numbers))
	}

	sorted := append([]string(nil), numbers...)
	sort.Strings(sorted)

	for i := 1; i <= count; i++ {
		want := fmt.Sprintf("%s%04d", prefix, i)
		if sorted[i-1] != want {
			t.Fatalf("expected sorted number %d to be %s, got %s", i, want, sorted[i-1])
		}
	}
}
