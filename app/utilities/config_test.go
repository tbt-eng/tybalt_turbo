package utilities

import (
	"testing"

	"tybalt/constants"
	"tybalt/internal/testseed"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

func upsertPurchaseOrdersConfig(t *testing.T, app *tests.TestApp, rawValue string) {
	t.Helper()

	collection, err := app.FindCollectionByNameOrId("app_config")
	if err != nil {
		t.Fatalf("failed to find app_config collection: %v", err)
	}

	record, err := app.FindFirstRecordByData("app_config", "key", "purchase_orders")
	if err != nil || record == nil {
		record = core.NewRecord(collection)
		record.Set("key", "purchase_orders")
	}

	record.Set("value", rawValue)
	if err := app.Save(record); err != nil {
		t.Fatalf("failed to save purchase_orders config: %v", err)
	}
}

func deletePurchaseOrdersConfig(t *testing.T, app *tests.TestApp) {
	t.Helper()

	record, err := app.FindFirstRecordByData("app_config", "key", "purchase_orders")
	if err != nil || record == nil {
		return
	}
	if err := app.Delete(record); err != nil {
		t.Fatalf("failed to delete purchase_orders config: %v", err)
	}
}

func upsertNotificationsConfig(t *testing.T, app *tests.TestApp, rawValue string) {
	t.Helper()

	collection, err := app.FindCollectionByNameOrId("app_config")
	if err != nil {
		t.Fatalf("failed to find app_config collection: %v", err)
	}

	record, err := app.FindFirstRecordByData("app_config", "key", "notifications")
	if err != nil || record == nil {
		record = core.NewRecord(collection)
		record.Set("key", "notifications")
	}

	record.Set("value", rawValue)
	if err := app.Save(record); err != nil {
		t.Fatalf("failed to save notifications config: %v", err)
	}
}

func deleteNotificationsConfig(t *testing.T, app *tests.TestApp) {
	t.Helper()

	record, err := app.FindFirstRecordByData("app_config", "key", "notifications")
	if err != nil || record == nil {
		return
	}
	if err := app.Delete(record); err != nil {
		t.Fatalf("failed to delete notifications config: %v", err)
	}
}

func TestGetPurchaseOrderSecondStageTimeoutHours(t *testing.T) {
	testsTable := []struct {
		name  string
		setup func(t *testing.T, app *tests.TestApp)
		want  float64
	}{
		{
			name: "returns default when purchase_orders config is missing",
			setup: func(t *testing.T, app *tests.TestApp) {
				deletePurchaseOrdersConfig(t, app)
			},
			want: 24.0,
		},
		{
			name: "returns default when timeout key is missing",
			setup: func(t *testing.T, app *tests.TestApp) {
				upsertPurchaseOrdersConfig(t, app, `{"some_other_key":1}`)
			},
			want: 24.0,
		},
		{
			name: "returns default when timeout is zero",
			setup: func(t *testing.T, app *tests.TestApp) {
				upsertPurchaseOrdersConfig(t, app, `{"second_stage_timeout_hours":0}`)
			},
			want: 24.0,
		},
		{
			name: "returns default when timeout is negative",
			setup: func(t *testing.T, app *tests.TestApp) {
				upsertPurchaseOrdersConfig(t, app, `{"second_stage_timeout_hours":-3}`)
			},
			want: 24.0,
		},
		{
			name: "returns configured positive timeout",
			setup: func(t *testing.T, app *tests.TestApp) {
				upsertPurchaseOrdersConfig(t, app, `{"second_stage_timeout_hours":36.5}`)
			},
			want: 36.5,
		},
	}

	for _, tc := range testsTable {
		t.Run(tc.name, func(t *testing.T) {
			app := testseed.NewSeededTestApp(t)
			defer app.Cleanup()

			tc.setup(t, app)

			got := GetPurchaseOrderSecondStageTimeoutHours(app)
			if got != tc.want {
				t.Fatalf("GetPurchaseOrderSecondStageTimeoutHours() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestIsLegacyPOCreateUpdateEnabled(t *testing.T) {
	testsTable := []struct {
		name  string
		setup func(t *testing.T, app *tests.TestApp)
		want  bool
	}{
		{
			name: "returns default false when purchase_orders config is missing",
			setup: func(t *testing.T, app *tests.TestApp) {
				deletePurchaseOrdersConfig(t, app)
			},
			want: false,
		},
		{
			name: "returns default false when legacy key is missing",
			setup: func(t *testing.T, app *tests.TestApp) {
				upsertPurchaseOrdersConfig(t, app, `{"second_stage_timeout_hours":24}`)
			},
			want: false,
		},
		{
			name: "returns configured true value",
			setup: func(t *testing.T, app *tests.TestApp) {
				upsertPurchaseOrdersConfig(t, app, `{"enable_legacy_po_create_update":true}`)
			},
			want: true,
		},
		{
			name: "returns configured false value",
			setup: func(t *testing.T, app *tests.TestApp) {
				upsertPurchaseOrdersConfig(t, app, `{"enable_legacy_po_create_update":false}`)
			},
			want: false,
		},
	}

	for _, tc := range testsTable {
		t.Run(tc.name, func(t *testing.T) {
			app := testseed.NewSeededTestApp(t)
			defer app.Cleanup()

			tc.setup(t, app)

			got, err := IsLegacyPOCreateUpdateEnabled(app)
			if err != nil {
				t.Fatalf("IsLegacyPOCreateUpdateEnabled() returned error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("IsLegacyPOCreateUpdateEnabled() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestIsNotificationFeatureEnabled(t *testing.T) {
	testsTable := []struct {
		name        string
		setup       func(t *testing.T, app *tests.TestApp)
		template    string
		wantEnabled bool
	}{
		{
			name: "returns default false when notifications config is missing",
			setup: func(t *testing.T, app *tests.TestApp) {
				deleteNotificationsConfig(t, app)
			},
			template:    "po_active",
			wantEnabled: false,
		},
		{
			name: "returns default false when template key is missing",
			setup: func(t *testing.T, app *tests.TestApp) {
				upsertNotificationsConfig(t, app, `{"expense_rejected":false}`)
			},
			template:    "po_active",
			wantEnabled: false,
		},
		{
			name: "returns false when template key is explicitly disabled",
			setup: func(t *testing.T, app *tests.TestApp) {
				upsertNotificationsConfig(t, app, `{"po_active":false}`)
			},
			template:    "po_active",
			wantEnabled: false,
		},
		{
			name: "returns true when template key is explicitly enabled",
			setup: func(t *testing.T, app *tests.TestApp) {
				upsertNotificationsConfig(t, app, `{"po_active":true}`)
			},
			template:    "po_active",
			wantEnabled: true,
		},
	}

	for _, tc := range testsTable {
		t.Run(tc.name, func(t *testing.T) {
			app := testseed.NewSeededTestApp(t)
			defer app.Cleanup()

			tc.setup(t, app)

			enabled, err := IsNotificationFeatureEnabled(app, tc.template)
			if err != nil {
				t.Fatalf("IsNotificationFeatureEnabled() returned error: %v", err)
			}
			if enabled != tc.wantEnabled {
				t.Fatalf("IsNotificationFeatureEnabled() = %v, want %v", enabled, tc.wantEnabled)
			}
		})
	}
}

func upsertExpensesConfig(t *testing.T, app *tests.TestApp, rawValue string) {
	t.Helper()

	collection, err := app.FindCollectionByNameOrId("app_config")
	if err != nil {
		t.Fatalf("failed to find app_config collection: %v", err)
	}

	record, err := app.FindFirstRecordByData("app_config", "key", "expenses")
	if err != nil || record == nil {
		record = core.NewRecord(collection)
		record.Set("key", "expenses")
	}

	record.Set("value", rawValue)
	if err := app.Save(record); err != nil {
		t.Fatalf("failed to save expenses config: %v", err)
	}
}

func deleteExpensesConfig(t *testing.T, app *tests.TestApp) {
	t.Helper()

	record, err := app.FindFirstRecordByData("app_config", "key", "expenses")
	if err != nil || record == nil {
		return
	}
	if err := app.Delete(record); err != nil {
		t.Fatalf("failed to delete expenses config: %v", err)
	}
}

func upsertTimeConfig(t *testing.T, app *tests.TestApp, rawValue string) {
	t.Helper()

	collection, err := app.FindCollectionByNameOrId("app_config")
	if err != nil {
		t.Fatalf("failed to find app_config collection: %v", err)
	}

	record, err := app.FindFirstRecordByData("app_config", "key", "time")
	if err != nil || record == nil {
		record = core.NewRecord(collection)
		record.Set("key", "time")
	}

	record.Set("value", rawValue)
	if err := app.Save(record); err != nil {
		t.Fatalf("failed to save time config: %v", err)
	}
}

func deleteTimeConfig(t *testing.T, app *tests.TestApp) {
	t.Helper()

	record, err := app.FindFirstRecordByData("app_config", "key", "time")
	if err != nil || record == nil {
		return
	}
	if err := app.Delete(record); err != nil {
		t.Fatalf("failed to delete time config: %v", err)
	}
}

func upsertLinkedWorkRecordsConfig(t *testing.T, app *tests.TestApp, rawValue string) {
	t.Helper()

	collection, err := app.FindCollectionByNameOrId("app_config")
	if err != nil {
		t.Fatalf("failed to find app_config collection: %v", err)
	}

	record, err := app.FindFirstRecordByData("app_config", "key", "linked_work_records")
	if err != nil || record == nil {
		record = core.NewRecord(collection)
		record.Set("key", "linked_work_records")
	}

	record.Set("value", rawValue)
	if err := app.Save(record); err != nil {
		t.Fatalf("failed to save linked_work_records config: %v", err)
	}
}

func deleteLinkedWorkRecordsConfig(t *testing.T, app *tests.TestApp) {
	t.Helper()

	record, err := app.FindFirstRecordByData("app_config", "key", "linked_work_records")
	if err != nil || record == nil {
		return
	}
	if err := app.Delete(record); err != nil {
		t.Fatalf("failed to delete linked_work_records config: %v", err)
	}
}

func TestIsTimeEditingEnabled(t *testing.T) {
	testsTable := []struct {
		name  string
		setup func(t *testing.T, app *tests.TestApp)
		want  bool
	}{
		{
			name: "returns default true when time config is missing",
			setup: func(t *testing.T, app *tests.TestApp) {
				deleteTimeConfig(t, app)
			},
			want: true,
		},
		{
			name: "returns default true when create_edit key is missing",
			setup: func(t *testing.T, app *tests.TestApp) {
				upsertTimeConfig(t, app, `{"some_other_key":true}`)
			},
			want: true,
		},
		{
			name: "returns configured true value",
			setup: func(t *testing.T, app *tests.TestApp) {
				upsertTimeConfig(t, app, `{"create_edit":true}`)
			},
			want: true,
		},
		{
			name: "returns configured false value",
			setup: func(t *testing.T, app *tests.TestApp) {
				upsertTimeConfig(t, app, `{"create_edit":false}`)
			},
			want: false,
		},
	}

	for _, tc := range testsTable {
		t.Run(tc.name, func(t *testing.T) {
			app := testseed.NewSeededTestApp(t)
			defer app.Cleanup()

			tc.setup(t, app)

			got, err := IsTimeEditingEnabled(app)
			if err != nil {
				t.Fatalf("IsTimeEditingEnabled() returned error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("IsTimeEditingEnabled() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestIsLinkedWorkRecordsEnabled(t *testing.T) {
	testsTable := []struct {
		name  string
		setup func(t *testing.T, app *tests.TestApp)
		want  bool
	}{
		{
			name: "returns default false when linked_work_records config is missing",
			setup: func(t *testing.T, app *tests.TestApp) {
				deleteLinkedWorkRecordsConfig(t, app)
			},
			want: false,
		},
		{
			name: "returns default false when enabled key is missing",
			setup: func(t *testing.T, app *tests.TestApp) {
				upsertLinkedWorkRecordsConfig(t, app, `{"some_other_key":true}`)
			},
			want: false,
		},
		{
			name: "returns configured true value",
			setup: func(t *testing.T, app *tests.TestApp) {
				upsertLinkedWorkRecordsConfig(t, app, `{"enabled":true}`)
			},
			want: true,
		},
		{
			name: "returns configured false value",
			setup: func(t *testing.T, app *tests.TestApp) {
				upsertLinkedWorkRecordsConfig(t, app, `{"enabled":false}`)
			},
			want: false,
		},
	}

	for _, tc := range testsTable {
		t.Run(tc.name, func(t *testing.T) {
			app := testseed.NewSeededTestApp(t)
			defer app.Cleanup()

			tc.setup(t, app)

			got, err := IsLinkedWorkRecordsEnabled(app)
			if err != nil {
				t.Fatalf("IsLinkedWorkRecordsEnabled() returned error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("IsLinkedWorkRecordsEnabled() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestGetPOExpenseExcessConfig(t *testing.T) {
	testsTable := []struct {
		name        string
		setup       func(t *testing.T, app *tests.TestApp)
		wantPercent float64
		wantValue   float64
		wantMode    string
	}{
		{
			name: "returns defaults when expenses config is missing",
			setup: func(t *testing.T, app *tests.TestApp) {
				deleteExpensesConfig(t, app)
			},
			wantPercent: constants.MAX_PURCHASE_ORDER_EXCESS_PERCENT,
			wantValue:   constants.MAX_PURCHASE_ORDER_EXCESS_VALUE,
			wantMode:    "lesser_of",
		},
		{
			name: "returns defaults when po_expense_allowed_excess key is missing",
			setup: func(t *testing.T, app *tests.TestApp) {
				upsertExpensesConfig(t, app, `{"create_edit_absorb":true}`)
			},
			wantPercent: constants.MAX_PURCHASE_ORDER_EXCESS_PERCENT,
			wantValue:   constants.MAX_PURCHASE_ORDER_EXCESS_VALUE,
			wantMode:    "lesser_of",
		},
		{
			name: "returns configured values when all fields present",
			setup: func(t *testing.T, app *tests.TestApp) {
				upsertExpensesConfig(t, app, `{"po_expense_allowed_excess":{"percent":10,"value":200.0,"mode":"greater_of"}}`)
			},
			wantPercent: 0.10,
			wantValue:   200.0,
			wantMode:    "greater_of",
		},
		{
			name: "returns partial override with defaults for missing fields",
			setup: func(t *testing.T, app *tests.TestApp) {
				upsertExpensesConfig(t, app, `{"po_expense_allowed_excess":{"percent":3}}`)
			},
			wantPercent: 0.03,
			wantValue:   constants.MAX_PURCHASE_ORDER_EXCESS_VALUE,
			wantMode:    "lesser_of",
		},
		{
			name: "returns default mode for invalid mode string",
			setup: func(t *testing.T, app *tests.TestApp) {
				upsertExpensesConfig(t, app, `{"po_expense_allowed_excess":{"percent":5,"value":100.0,"mode":"invalid"}}`)
			},
			wantPercent: 0.05,
			wantValue:   100.0,
			wantMode:    "lesser_of",
		},
		{
			name: "accepts zero percent",
			setup: func(t *testing.T, app *tests.TestApp) {
				upsertExpensesConfig(t, app, `{"po_expense_allowed_excess":{"percent":0,"value":50.0,"mode":"lesser_of"}}`)
			},
			wantPercent: 0,
			wantValue:   50.0,
			wantMode:    "lesser_of",
		},
		{
			name: "rejects negative percent and uses default",
			setup: func(t *testing.T, app *tests.TestApp) {
				upsertExpensesConfig(t, app, `{"po_expense_allowed_excess":{"percent":-5,"value":100.0,"mode":"lesser_of"}}`)
			},
			wantPercent: constants.MAX_PURCHASE_ORDER_EXCESS_PERCENT,
			wantValue:   100.0,
			wantMode:    "lesser_of",
		},
		{
			name: "rejects percent above 100 and uses default",
			setup: func(t *testing.T, app *tests.TestApp) {
				upsertExpensesConfig(t, app, `{"po_expense_allowed_excess":{"percent":150,"value":100.0,"mode":"lesser_of"}}`)
			},
			wantPercent: constants.MAX_PURCHASE_ORDER_EXCESS_PERCENT,
			wantValue:   100.0,
			wantMode:    "lesser_of",
		},
	}

	for _, tc := range testsTable {
		t.Run(tc.name, func(t *testing.T) {
			app := testseed.NewSeededTestApp(t)
			defer app.Cleanup()

			tc.setup(t, app)

			got := GetPOExpenseExcessConfig(app)
			if got.Percent != tc.wantPercent {
				t.Errorf("Percent = %v, want %v", got.Percent, tc.wantPercent)
			}
			if got.Value != tc.wantValue {
				t.Errorf("Value = %v, want %v", got.Value, tc.wantValue)
			}
			if got.Mode != tc.wantMode {
				t.Errorf("Mode = %q, want %q", got.Mode, tc.wantMode)
			}
		})
	}
}

func TestCalculatePOExpenseTotalLimit(t *testing.T) {
	tests := []struct {
		name                string
		poTotal             float64
		cfg                 POExpenseExcessConfig
		wantTotalLimitCents int64
		wantExcessText      string
	}{
		{
			name:    "lesser_of: percent is lesser (small PO)",
			poTotal: 500.0,
			cfg:     POExpenseExcessConfig{Percent: 0.05, Value: 100.0, Mode: "lesser_of"},
			// 5% of 500 = 25, which is < 100, so percent applies
			wantTotalLimitCents: 52500,
			wantExcessText:      "5.00%",
		},
		{
			name:    "lesser_of: value is lesser (large PO)",
			poTotal: 5000.0,
			cfg:     POExpenseExcessConfig{Percent: 0.05, Value: 100.0, Mode: "lesser_of"},
			// 5% of 5000 = 250, which is > 100, so value applies
			wantTotalLimitCents: 510000,
			wantExcessText:      "$100.00",
		},
		{
			name:    "lesser_of: equal at breakpoint",
			poTotal: 2000.0,
			cfg:     POExpenseExcessConfig{Percent: 0.05, Value: 100.0, Mode: "lesser_of"},
			// 5% of 2000 = 100, equals value, so percent branch applies (value < percent is false)
			wantTotalLimitCents: 210000,
			wantExcessText:      "5.00%",
		},
		{
			name:    "greater_of: value is greater (small PO)",
			poTotal: 500.0,
			cfg:     POExpenseExcessConfig{Percent: 0.05, Value: 100.0, Mode: "greater_of"},
			// 5% of 500 = 25, value 100 >= 25, so value applies
			wantTotalLimitCents: 60000,
			wantExcessText:      "$100.00",
		},
		{
			name:    "greater_of: percent is greater (large PO)",
			poTotal: 5000.0,
			cfg:     POExpenseExcessConfig{Percent: 0.05, Value: 100.0, Mode: "greater_of"},
			// 5% of 5000 = 250, value 100 < 250, so percent applies
			wantTotalLimitCents: 525000,
			wantExcessText:      "5.00%",
		},
		{
			name:    "zero PO total",
			poTotal: 0,
			cfg:     POExpenseExcessConfig{Percent: 0.05, Value: 100.0, Mode: "lesser_of"},
			// 5% of 0 = 0, which is < 100, so percent applies; total = 0
			wantTotalLimitCents: 0,
			wantExcessText:      "5.00%",
		},
		{
			name:    "percent limit landing exactly on a cent is not floored below it",
			poTotal: 3.80,
			cfg:     POExpenseExcessConfig{Percent: 0.05, Value: 100.0, Mode: "lesser_of"},
			// 3.80 * 1.05 is 3.9899999999999998 in floats; the exact limit is 3.99
			wantTotalLimitCents: 399,
			wantExcessText:      "5.00%",
		},
		{
			name:    "value limit landing exactly on a cent is not floored below it",
			poTotal: 2000.03,
			cfg:     POExpenseExcessConfig{Percent: 0.05, Value: 100.0, Mode: "lesser_of"},
			// 2000.03 + 100 is 2100.0299999999997 in floats; the exact limit is 2100.03
			wantTotalLimitCents: 210003,
			wantExcessText:      "$100.00",
		},
		{
			name:    "limit between cents is floored to the cent below",
			poTotal: 16.27,
			cfg:     POExpenseExcessConfig{Percent: 0.05, Value: 100.0, Mode: "lesser_of"},
			// 16.27 * 1.05 = 17.0835, so 17.08 is the largest allowed total
			wantTotalLimitCents: 1708,
			wantExcessText:      "5.00%",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := CalculatePOExpenseTotalLimit(tc.poTotal, POExpenseCurrency{}, tc.cfg)
			if got.TotalLimitCents != tc.wantTotalLimitCents {
				t.Errorf("TotalLimitCents = %v, want %v", got.TotalLimitCents, tc.wantTotalLimitCents)
			}
			if got.ExcessText != tc.wantExcessText {
				t.Errorf("ExcessText = %q, want %q", got.ExcessText, tc.wantExcessText)
			}
		})
	}
}

// The fixed amount is in the home currency (issue #183). On a foreign PO it is
// divided by the PO's rate (home units per unit of the PO's currency) before
// the mode compares it with the percent and before it is added.
func TestCalculatePOExpenseTotalLimitForeignCurrency(t *testing.T) {
	usd := POExpenseCurrency{Code: "USD", HomeRate: 1.35} // $100 is 74.074... USD
	lesser := POExpenseExcessConfig{Percent: 0.05, Value: 100.0, Mode: "lesser_of"}
	greater := POExpenseExcessConfig{Percent: 0.05, Value: 100.0, Mode: "greater_of"}

	tests := []struct {
		name                string
		poTotal             float64
		currency            POExpenseCurrency
		cfg                 POExpenseExcessConfig
		wantTotalLimitCents int64
		wantExcessText      string
		wantPOTotalText     string
	}{
		{
			name: "lesser_of: fixed amount is lesser after conversion (old code chose the percent)",
			// 5% of 1500 is 75.00 USD; $100 is 74.07 USD. The old code compared
			// 100 with 75 and allowed 1575.00.
			poTotal: 1500, currency: usd, cfg: lesser,
			wantTotalLimitCents: 157407,
			wantExcessText:      "100.00 CAD (74.07 USD)",
			wantPOTotalText:     "1,500.00 USD",
		},
		{
			name:    "lesser_of: percent is lesser",
			poTotal: 1000, currency: usd, cfg: lesser,
			// 5% of 1000 is 50.00 USD < 74.07 USD
			wantTotalLimitCents: 105000,
			wantExcessText:      "5.00%",
			wantPOTotalText:     "1,000.00 USD",
		},
		{
			name:    "lesser_of: large PO uses the converted fixed amount",
			poTotal: 4000, currency: usd, cfg: lesser,
			// The old code allowed 4100.00
			wantTotalLimitCents: 407407,
			wantExcessText:      "100.00 CAD (74.07 USD)",
			wantPOTotalText:     "4,000.00 USD",
		},
		{
			name: "greater_of: percent is greater after conversion (old code chose the fixed amount)",
			// 75.00 USD > 74.07 USD. The old code compared 100 with 75 and
			// allowed 1600.00.
			poTotal: 1500, currency: usd, cfg: greater,
			wantTotalLimitCents: 157500,
			wantExcessText:      "5.00%",
			wantPOTotalText:     "1,500.00 USD",
		},
		{
			name:    "greater_of: fixed amount is greater",
			poTotal: 1000, currency: usd, cfg: greater,
			// 74.07 USD > 50.00 USD; the old code allowed 1100.00
			wantTotalLimitCents: 107407,
			wantExcessText:      "100.00 CAD (74.07 USD)",
			wantPOTotalText:     "1,000.00 USD",
		},
		{
			name: "a currency with a small unit gets a large fixed amount",
			// 1 JPY is 0.0091 CAD, so $100 is 10989.01 JPY, less than 5% of
			// 1,000,000 JPY (50000). The old code allowed 1,000,100 JPY.
			poTotal: 1000000, currency: POExpenseCurrency{Code: "JPY", HomeRate: 0.0091}, cfg: lesser,
			wantTotalLimitCents: 101098901,
			wantExcessText:      "100.00 CAD (10,989.01 JPY)",
			wantPOTotalText:     "1,000,000.00 JPY",
		},
		{
			name: "a rate of 0 leaves the fixed amount unconverted",
			// GetPOExpenseCurrency never returns 0; this guards a direct caller.
			poTotal: 4000, currency: POExpenseCurrency{Code: "USD"}, cfg: lesser,
			wantTotalLimitCents: 410000,
			wantExcessText:      "100.00 CAD (100.00 USD)",
			wantPOTotalText:     "4,000.00 USD",
		},
		{
			name: "the home currency is never converted, even with a rate that is not 1",
			// A stored home total can differ from the total by rounding.
			poTotal: 4000, currency: POExpenseCurrency{Code: "CAD", HomeRate: 1.0001}, cfg: lesser,
			wantTotalLimitCents: 410000,
			wantExcessText:      "$100.00",
			wantPOTotalText:     "$4000.00",
		},
		{
			name:    "home currency (blank code) keeps the old text",
			poTotal: 1000, currency: POExpenseCurrency{}, cfg: greater,
			wantTotalLimitCents: 110000,
			wantExcessText:      "$100.00",
			wantPOTotalText:     "$1000.00",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := CalculatePOExpenseTotalLimit(tc.poTotal, tc.currency, tc.cfg)
			if got.TotalLimitCents != tc.wantTotalLimitCents {
				t.Errorf("TotalLimitCents = %v, want %v", got.TotalLimitCents, tc.wantTotalLimitCents)
			}
			if got.ExcessText != tc.wantExcessText {
				t.Errorf("ExcessText = %q, want %q", got.ExcessText, tc.wantExcessText)
			}
			if got.POTotalText != tc.wantPOTotalText {
				t.Errorf("POTotalText = %q, want %q", got.POTotalText, tc.wantPOTotalText)
			}
		})
	}
}

// GetPOExpenseCurrency takes the rate from the PO's stored home total, and
// falls back to the currency's rate today, then to 1. The POs are built in
// memory; the currencies are fixtures (USD rate 1.35, CHF wipcurr00000001 has
// no rate).
func TestGetPOExpenseCurrency(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	collection, err := app.FindCollectionByNameOrId("purchase_orders")
	if err != nil {
		t.Fatalf("failed to load purchase_orders collection: %v", err)
	}

	tests := []struct {
		name string
		po   map[string]any
		want POExpenseCurrency
	}{
		{
			name: "blank currency is the home currency",
			po:   map[string]any{"approval_total": 100.0, "approval_total_home": 100.0},
			want: POExpenseCurrency{},
		},
		{
			name: "CAD is never converted, even when the stored home total differs by rounding",
			po:   map[string]any{"currency": "cadcurr00000001", "approval_total": 100.0, "approval_total_home": 100.004},
			want: POExpenseCurrency{},
		},
		{
			name: "USD uses the rate of the stored home total, not the rate today",
			po:   map[string]any{"currency": "usdcurr00000001", "approval_total": 1000.0, "approval_total_home": 1400.0},
			want: POExpenseCurrency{Code: "USD", HomeRate: 1.4},
		},
		{
			name: "USD without a stored home total falls back to the rate today",
			po:   map[string]any{"currency": "usdcurr00000001", "approval_total": 1000.0},
			want: POExpenseCurrency{Code: "USD", HomeRate: 1.35},
		},
		{
			name: "USD with a zero approval total falls back to the rate today",
			po:   map[string]any{"currency": "usdcurr00000001", "approval_total_home": 1350.0},
			want: POExpenseCurrency{Code: "USD", HomeRate: 1.35},
		},
		{
			name: "a currency with no rate and no stored home total falls back to 1",
			po:   map[string]any{"currency": "wipcurr00000001", "approval_total": 1000.0},
			want: POExpenseCurrency{Code: "CHF", HomeRate: 1},
		},
		{
			name: "a currency that can't be found is treated as the home currency",
			po:   map[string]any{"currency": "nosuchcurrency1", "approval_total": 1000.0, "approval_total_home": 1350.0},
			want: POExpenseCurrency{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			po := core.NewRecord(collection)
			for k, v := range tc.po {
				po.Set(k, v)
			}
			if got := GetPOExpenseCurrency(app, po); got != tc.want {
				t.Errorf("GetPOExpenseCurrency = %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestGetNoPOExpenseLimit(t *testing.T) {
	testsTable := []struct {
		name  string
		setup func(t *testing.T, app *tests.TestApp)
		want  float64
	}{
		{
			name: "returns default when expenses config is missing",
			setup: func(t *testing.T, app *tests.TestApp) {
				deleteExpensesConfig(t, app)
			},
			want: constants.NO_PO_EXPENSE_LIMIT,
		},
		{
			name: "returns default when no_po_expense_limit key is missing",
			setup: func(t *testing.T, app *tests.TestApp) {
				upsertExpensesConfig(t, app, `{"create_edit_absorb":true}`)
			},
			want: constants.NO_PO_EXPENSE_LIMIT,
		},
		{
			name: "returns configured positive value",
			setup: func(t *testing.T, app *tests.TestApp) {
				upsertExpensesConfig(t, app, `{"no_po_expense_limit":250.0}`)
			},
			want: 250.0,
		},
		{
			name: "returns zero when value is zero",
			setup: func(t *testing.T, app *tests.TestApp) {
				upsertExpensesConfig(t, app, `{"no_po_expense_limit":0}`)
			},
			want: 0,
		},
		{
			name: "returns default when value is negative",
			setup: func(t *testing.T, app *tests.TestApp) {
				upsertExpensesConfig(t, app, `{"no_po_expense_limit":-50}`)
			},
			want: constants.NO_PO_EXPENSE_LIMIT,
		},
		{
			name: "returns default when value is non-numeric",
			setup: func(t *testing.T, app *tests.TestApp) {
				upsertExpensesConfig(t, app, `{"no_po_expense_limit":"abc"}`)
			},
			want: constants.NO_PO_EXPENSE_LIMIT,
		},
	}

	for _, tc := range testsTable {
		t.Run(tc.name, func(t *testing.T) {
			app := testseed.NewSeededTestApp(t)
			defer app.Cleanup()

			tc.setup(t, app)

			got := GetNoPOExpenseLimit(app)
			if got != tc.want {
				t.Fatalf("GetNoPOExpenseLimit() = %v, want %v", got, tc.want)
			}
		})
	}
}

// addDeployGroup makes a deploy group with the users as members and returns
// its id.
func addDeployGroup(t *testing.T, app *tests.TestApp, name string, uids ...string) string {
	t.Helper()
	groups, err := app.FindCollectionByNameOrId("deploy_groups")
	if err != nil {
		t.Fatal(err)
	}
	group := core.NewRecord(groups)
	group.Set("name", name)
	if err := app.Save(group); err != nil {
		t.Fatal(err)
	}
	memberships, err := app.FindCollectionByNameOrId("deploy_group_memberships")
	if err != nil {
		t.Fatal(err)
	}
	for _, uid := range uids {
		membership := core.NewRecord(memberships)
		membership.Set("deploy_group", group.Id)
		membership.Set("user", uid)
		if err := app.Save(membership); err != nil {
			t.Fatal(err)
		}
	}
	return group.Id
}

func setConfigRow(t *testing.T, app *tests.TestApp, key, rawValue string) {
	t.Helper()
	collection, err := app.FindCollectionByNameOrId("app_config")
	if err != nil {
		t.Fatal(err)
	}
	record, err := app.FindFirstRecordByData("app_config", "key", key)
	if err != nil || record == nil {
		record = core.NewRecord(collection)
		record.Set("key", key)
	}
	record.Set("value", rawValue)
	if err := app.Save(record); err != nil {
		t.Fatal(err)
	}
}

func twoUserIDs(t *testing.T, app *tests.TestApp) (string, string) {
	t.Helper()
	users, err := app.FindRecordsByFilter("users", "id != ''", "id", 2, 0)
	if err != nil || len(users) < 2 {
		t.Fatalf("the seed needs two users: %v", err)
	}
	return users[0].Id, users[1].Id
}

func TestInDeployGroup(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()
	member, outsider := twoUserIDs(t, app)
	group := addDeployGroup(t, app, "testers", member)
	empty := addDeployGroup(t, app, "nobody yet")

	cases := []struct {
		name  string
		group string
		uid   string
		want  bool
	}{
		{"no group is every user", "", outsider, true},
		{"no group is every user, also with no user", "", "", true},
		{"a member", group, member, true},
		{"a user outside the group", group, outsider, false},
		{"a group with no members", empty, member, false},
		{"an id that matches no group", "nosuchgroup00000", member, false},
		{"no user", group, "", false},
	}
	for _, c := range cases {
		got, err := InDeployGroup(app, c.group, c.uid)
		if err != nil {
			t.Errorf("%s: %v", c.name, err)
		}
		if got != c.want {
			t.Errorf("%s: got %v, want %v", c.name, got, c.want)
		}
	}
}

// A second membership row for the same user and group is refused, so a group
// can't count a user twice.
func TestADeployGroupHoldsAUserOnce(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()
	member, _ := twoUserIDs(t, app)
	group := addDeployGroup(t, app, "testers", member)

	memberships, err := app.FindCollectionByNameOrId("deploy_group_memberships")
	if err != nil {
		t.Fatal(err)
	}
	again := core.NewRecord(memberships)
	again.Set("deploy_group", group)
	again.Set("user", member)
	if err := app.Save(again); err == nil {
		t.Fatal("the same user was added to the group twice")
	}
}

func TestDeployGroupAllows(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()
	member, outsider := twoUserIDs(t, app)
	group := addDeployGroup(t, app, "testers", member)
	const key = "some_feature"

	check := func(name, uid string, want, wantErr bool) {
		t.Helper()
		got, err := DeployGroupAllows(app, key, uid)
		if (err != nil) != wantErr {
			t.Errorf("%s: error %v, want error %v", name, err, wantErr)
		}
		if got != want {
			t.Errorf("%s: got %v, want %v", name, got, want)
		}
	}

	check("no row", outsider, true, false)

	setConfigRow(t, app, key, `{"enabled": true}`)
	check("no deploy_group", outsider, true, false)

	setConfigRow(t, app, key, `{"enabled": true, "deploy_group": ""}`)
	check("an empty deploy_group", outsider, true, false)

	setConfigRow(t, app, key, `{"enabled": true, "deploy_group": "`+group+`"}`)
	check("a member", member, true, false)
	check("a user outside the group", outsider, false, false)

	setConfigRow(t, app, key, `{"enabled": true, "deploy_group": "nosuchgroup00000"}`)
	check("an id that matches no group", member, false, false)

	setConfigRow(t, app, key, `{"enabled": true, "deploy_group": 7}`)
	check("a deploy_group that is not text", member, false, true)
}
