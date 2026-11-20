package routes

import (
	"encoding/json"
	"strings"
	"testing"
)

// The cases of the legacy utilities.test.ts, plus JavaScript white space.
func TestMakeSlug(t *testing.T) {
	for _, tc := range []struct {
		serial, mfg, want string
	}{
		{"kQDS 9/452,At", "Dell, Inc. Ltd  ", "kQDS9452At,dell"},
		{"LENISJ/PFU ", "Lenovo Computer, Inc. Ltd  ", "LENISJPFU,lenovo_computer"},
		{"FE/KJ", "Dell, Inc. Ltd", "FEKJ,dell"},
		{"FE/KJ", "JR Inc.", "FEKJ,jr"},
		{"AB CD\tEF", "　Hewlett-Packard ", "ABCDEF,hewlett-packard"},
		{"SN123", "manufac", "SN123,manufac"},
	} {
		got, err := makeSlug(tc.serial, tc.mfg)
		if err != nil || got != tc.want {
			t.Errorf("makeSlug(%q, %q) = %q, %v; want %q", tc.serial, tc.mfg, got, err, tc.want)
		}
	}
	for _, tc := range []struct{ serial, mfg string }{
		{"FE/K", "Dell, Inc. Ltd"},
		{"FE/KJ", "J Inc."},
	} {
		if got, err := makeSlug(tc.serial, tc.mfg); err == nil {
			t.Errorf("makeSlug(%q, %q) = %q; want an error", tc.serial, tc.mfg, got)
		}
	}
}

func TestCoerceToJSInteger(t *testing.T) {
	for _, tc := range []struct {
		raw  any
		want float64
		ok   bool
	}{
		{json.Number("5"), 5, true},
		{json.Number("5.0"), 5, true},
		{json.Number("5.5"), 0, false},
		{" 5.", 5, true},
		{"48", 48, true},
		{"1e3", 1000, true},
		{"0x10", 16, true},
		{"   ", 0, true},
		{"", 0, false},
		{"4.5", 0, false},
		{"abc", 0, false},
		{"Infinity", 0, false},
		{"NaN", 0, false},
		{"inf", 0, false},
		{true, 1, true},
		{false, 0, true},
		{nil, 0, true},
		{map[string]any{}, 0, false},
	} {
		got, ok := coerceToJSInteger(tc.raw)
		if ok != tc.ok || (ok && got != tc.want) {
			t.Errorf("coerceToJSInteger(%#v) = %v, %v; want %v, %v", tc.raw, got, ok, tc.want, tc.ok)
		}
	}
}

func TestCoerceToJSString(t *testing.T) {
	for _, tc := range []struct {
		raw  any
		want string
		ok   bool
	}{
		{"abc", "abc", true},
		{json.Number("12345678"), "12345678", true},
		{json.Number("1.50"), "1.5", true},
		{json.Number("1e21"), "1e+21", true},
		{json.Number("0.0000001"), "1e-7", true},
		{json.Number("0.000001"), "0.000001", true},
		{true, "true", true},
		{nil, "", true},
		{[]any{"a"}, "", false},
	} {
		got, ok := coerceToJSString(tc.raw)
		if ok != tc.ok || got != tc.want {
			t.Errorf("coerceToJSString(%#v) = %q, %v; want %q, %v", tc.raw, got, ok, tc.want, tc.ok)
		}
	}
}

func TestValidateRadiatorReport(t *testing.T) {
	valid := map[string]any{
		"serial":           "SN123",
		"mfg":              "manufac",
		"upn":              "TTesterson@testco.co",
		"email":            "",
		"userSourceAnchor": "F25D2A25F25D2A25F25D2A25F25D2A25",
		"networkConfig":    map[string]any{"DC:4A:3E:E0:45:00": map[string]any{}},
		"radiatorVersion":  json.Number("7"),
		"systemType":       " 5.",
		"osSku":            "48",
		"computerName":     "Tromsø",
		"userGivenName":    "Testy",
		"userSurname":      "Testerson",
		"extra":            "dropped",
	}
	report, problems := validateRadiatorReport(valid)
	if len(problems) != 0 {
		t.Fatalf("problems = %v", problems)
	}
	if _, ok := report["email"]; ok {
		t.Error("an invalid email must be dropped")
	}
	if _, ok := report["extra"]; ok {
		t.Error("an unknown property must be dropped")
	}
	if report["systemType"] != float64(5) || report["osSku"] != float64(48) {
		t.Errorf("coercion: systemType %v, osSku %v", report["systemType"], report["osSku"])
	}

	invalid := map[string]any{
		"serial":           "SN123",
		"mfg":              "m",
		"upn":              "no-at-sign",
		"userSourceAnchor": nil,
		"networkConfig":    map[string]any{"not-a-mac": true},
		"radiatorVersion":  json.Number("9"),
		"userGivenName":    "Testy",
		"userSurname":      "Testerson",
		"bootDrive":        "CC:",
	}
	_, problems = validateRadiatorReport(invalid)
	joined := strings.Join(problems, "\n")
	for _, want := range []string{
		"mfg: must be at least 2 characters",
		"upn: must match the pattern",
		"userSourceAnchor: must match the pattern",
		`networkConfig: the property name "not-a-mac" must be a MAC address`,
		"bootDrive: must match the pattern",
	} {
		if !strings.Contains(joined, want) {
			t.Errorf("problems %q do not contain %q", joined, want)
		}
	}
}
