package routes

import (
	"reflect"
	"testing"

	"tybalt/internal/testseed"

	"github.com/pocketbase/dbx"
)

// TestVendorMergeKeepsInvoiceNumbersUnique checks that a vendor merge finds
// invoice numbers both vendors have recorded, and otherwise moves the invoices
// to the kept vendor.
func TestVendorMergeKeepsInvoiceNumbersUnique(t *testing.T) {
	app := testseed.NewSeededTestApp(t)
	defer app.Cleanup()

	const kept, absorbed = "z66xe6vqhwtokt4", "2zqxtsmymf670ha"
	insert := func(id, vendor, number, key string) {
		t.Helper()
		_, err := app.DB().Insert("expense_documents", dbx.Params{
			"id": id, "attachment": "invoice.png", "attachment_hash": id, "uploaded_by": "",
			"vendor": vendor, "vendor_invoice_number": number, "vendor_invoice_number_key": key,
			"vendor_invoice_date": "2024-08-30", "vendor_invoice_total": 10, "vendor_invoice_currency": "",
			"created": "2026-10-19 00:00:00.000Z", "updated": "2026-10-19 00:00:00.000Z",
		}).Execute()
		if err != nil {
			t.Fatal(err)
		}
	}
	insert("vidockept000001", kept, "INV-1", "inv1")
	insert("vidocabsorbed01", absorbed, "inv 1", "inv1")
	insert("vidocabsorbed02", absorbed, "INV-2", "inv2")

	clash, err := clashingVendorInvoiceNumbers(app, []string{kept, absorbed})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(clash, []string{"INV-1"}) {
		t.Fatalf("clash = %v; want [INV-1]", clash)
	}

	// With the clash corrected, the merge moves the absorbed vendor's invoices.
	if _, err := app.DB().Delete("expense_documents", dbx.HashExp{"id": "vidocabsorbed01"}).Execute(); err != nil {
		t.Fatal(err)
	}
	if clash, _ := clashingVendorInvoiceNumbers(app, []string{kept, absorbed}); len(clash) != 0 {
		t.Fatalf("clash = %v; want none", clash)
	}
	if err := AbsorbRecords(app, "vendors", kept, []string{absorbed}); err != nil {
		t.Fatal(err)
	}
	document, err := app.FindRecordById("expense_documents", "vidocabsorbed02")
	if err != nil {
		t.Fatal(err)
	}
	if got := document.GetString("vendor"); got != kept {
		t.Fatalf("vendor = %q; want the kept vendor", got)
	}
}
