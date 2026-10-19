package hooks

import "testing"

func TestVendorInvoiceNumberKey(t *testing.T) {
	for number, want := range map[string]string{
		"INV-0042":   "inv0042",
		" inv 0042 ": "inv0042",
		"INV0042":    "inv0042",
		"0042":       "0042",
		"A/B.12#3":   "ab123",
		"Façture-7":  "façture7",
		"--":         "",
		"":           "",
	} {
		if got := VendorInvoiceNumberKey(number); got != want {
			t.Errorf("VendorInvoiceNumberKey(%q) = %q; want %q", number, got, want)
		}
	}
}
