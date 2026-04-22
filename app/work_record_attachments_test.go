package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"
	"tybalt/internal/testutils"

	"github.com/pocketbase/pocketbase/tests"
)

type workRecordAttachmentRouteResponse struct {
	ID         string `json:"id"`
	Attachment string `json:"attachment"`
}

func loadWorkRecordAttachmentStateForTest(t testing.TB, app *tests.TestApp, workRecordID string) (string, string) {
	t.Helper()

	var row struct {
		Attachment     string `db:"attachment"`
		AttachmentHash string `db:"attachment_hash"`
	}
	if err := app.DB().NewQuery(`
		SELECT
			COALESCE(attachment, '') AS attachment,
			COALESCE(attachment_hash, '') AS attachment_hash
		FROM work_records
		WHERE id = {:id}
		LIMIT 1
	`).Bind(map[string]any{
		"id": workRecordID,
	}).One(&row); err != nil {
		t.Fatalf("failed loading work-record attachment state: %v", err)
	}

	return row.Attachment, row.AttachmentHash
}

func makeWorkRecordPDFMultipart(t testing.TB, fileName string, fileContent []byte) (*bytes.Buffer, string) {
	t.Helper()

	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)
	fw, err := w.CreateFormFile("attachment", fileName)
	if err != nil {
		t.Fatalf("failed creating attachment form file: %v", err)
	}
	if _, err := fw.Write(fileContent); err != nil {
		t.Fatalf("failed writing attachment content: %v", err)
	}
	contentType := w.FormDataContentType()
	if err := w.Close(); err != nil {
		t.Fatalf("failed closing multipart writer: %v", err)
	}

	return buf, contentType
}

func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func TestWorkRecordAttachmentRoutes(t *testing.T) {
	ownerToken, err := testutils.GenerateRecordToken("users", "author@soup.com")
	if err != nil {
		t.Fatal(err)
	}
	reportToken, err := testutils.GenerateRecordToken("users", "fatt@mac.com")
	if err != nil {
		t.Fatal(err)
	}

	uploadContent := []byte("%PDF-1.4 upload test\n")
	replaceContent := []byte("%PDF-1.4 replace test\n")
	duplicateContent := []byte("%PDF-1.4 duplicate test\n")

	scenarios := []tests.ApiScenario{
		func() tests.ApiScenario {
			body, contentType := makeWorkRecordPDFMultipart(t, "paper-copy.pdf", uploadContent)
			return tests.ApiScenario{
				Name:   "owner can upload work record paper copy pdf",
				Method: http.MethodPost,
				URL:    "/api/work_records/wr260400000010a/attachment",
				Headers: map[string]string{
					"Authorization": ownerToken,
					"Content-Type":  contentType,
				},
				Body:           body,
				ExpectedStatus: http.StatusOK,
				ExpectedContent: []string{
					`"id":"wr260400000010a"`,
					`"attachment":`,
				},
				TestAppFactory: testutils.SetupTestApp,
				AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
					defer res.Body.Close()

					var payload workRecordAttachmentRouteResponse
					if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
						t.Fatalf("failed decoding attachment upload response: %v", err)
					}
					if payload.ID != "wr260400000010a" {
						t.Fatalf("expected uploaded response id, got %#v", payload)
					}
					if payload.Attachment == "" {
						t.Fatal("expected uploaded attachment filename in response")
					}

					attachment, attachmentHash := loadWorkRecordAttachmentStateForTest(t, app, "wr260400000010a")
					if attachment == "" {
						t.Fatal("expected uploaded attachment filename to be stored")
					}
					if attachmentHash != sha256Hex(uploadContent) {
						t.Fatalf("expected uploaded hash %q, got %q", sha256Hex(uploadContent), attachmentHash)
					}
				},
			}
		}(),
		func() tests.ApiScenario {
			body, contentType := makeWorkRecordPDFMultipart(t, "replacement.pdf", replaceContent)
			return tests.ApiScenario{
				Name:   "owner can replace work record paper copy pdf",
				Method: http.MethodPost,
				URL:    "/api/work_records/wr260400000003a/attachment",
				Headers: map[string]string{
					"Authorization": ownerToken,
					"Content-Type":  contentType,
				},
				Body:           body,
				ExpectedStatus: http.StatusOK,
				ExpectedContent: []string{
					`"attachment":`,
				},
				TestAppFactory: testutils.SetupTestApp,
				AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
					defer res.Body.Close()

					var payload workRecordAttachmentRouteResponse
					if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
						t.Fatalf("failed decoding attachment replace response: %v", err)
					}
					if payload.Attachment == "" || payload.Attachment == "paper_copy.pdf" {
						t.Fatalf("expected replacement attachment filename in response, got %#v", payload)
					}

					attachment, attachmentHash := loadWorkRecordAttachmentStateForTest(t, app, "wr260400000003a")
					if attachment == "" || attachment == "paper_copy.pdf" {
						t.Fatalf("expected replacement attachment filename, got %q", attachment)
					}
					if attachmentHash != sha256Hex(replaceContent) {
						t.Fatalf("expected replacement hash %q, got %q", sha256Hex(replaceContent), attachmentHash)
					}
				},
			}
		}(),
		{
			Name:   "owner can remove work record paper copy pdf",
			Method: http.MethodDelete,
			URL:    "/api/work_records/wr260400000003a/attachment",
			Headers: map[string]string{
				"Authorization": ownerToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"id":"wr260400000003a"`,
				`"attachment":""`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var payload workRecordAttachmentRouteResponse
				if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
					t.Fatalf("failed decoding attachment delete response: %v", err)
				}
				if payload.Attachment != "" {
					t.Fatalf("expected empty attachment response after delete, got %#v", payload)
				}

				attachment, attachmentHash := loadWorkRecordAttachmentStateForTest(t, app, "wr260400000003a")
				if attachment != "" || attachmentHash != "" {
					t.Fatalf("expected attachment removal to clear state, got attachment=%q hash=%q", attachment, attachmentHash)
				}
			},
		},
		func() tests.ApiScenario {
			body, contentType := makeWorkRecordPDFMultipart(t, "duplicate.pdf", duplicateContent)
			return tests.ApiScenario{
				Name:   "duplicate work record paper copy pdf is rejected",
				Method: http.MethodPost,
				URL:    "/api/work_records/wr260400000010a/attachment",
				Headers: map[string]string{
					"Authorization": ownerToken,
					"Content-Type":  contentType,
				},
				Body:           body,
				ExpectedStatus: http.StatusBadRequest,
				ExpectedContent: []string{
					`"attachment":{"code":"duplicate_file"`,
					`"message":"This file has already been uploaded to another work record"`,
				},
				TestAppFactory: testutils.SetupTestApp,
			}
		}(),
		func() tests.ApiScenario {
			body, contentType := makeWorkRecordPDFMultipart(t, "approved.pdf", uploadContent)
			return tests.ApiScenario{
				Name:   "approved work record attachment mutation is rejected",
				Method: http.MethodPost,
				URL:    "/api/work_records/wr300200000050a/attachment",
				Headers: map[string]string{
					"Authorization": ownerToken,
					"Content-Type":  contentType,
				},
				Body:           body,
				ExpectedStatus: http.StatusBadRequest,
				ExpectedContent: []string{
					`"code":"approved_work_record_locked"`,
				},
				TestAppFactory: testutils.SetupTestApp,
			}
		}(),
		func() tests.ApiScenario {
			body, contentType := makeWorkRecordPDFMultipart(t, "report-viewer.pdf", uploadContent)
			return tests.ApiScenario{
				Name:   "report viewer cannot upload work record attachment",
				Method: http.MethodPost,
				URL:    "/api/work_records/wr260400000010a/attachment",
				Headers: map[string]string{
					"Authorization": reportToken,
					"Content-Type":  contentType,
				},
				Body:           body,
				ExpectedStatus: http.StatusForbidden,
				ExpectedContent: []string{
					`"code":"unauthorized"`,
					`"message":"you are not authorized to edit this work record"`,
				},
				TestAppFactory: testutils.SetupTestApp,
			}
		}(),
		{
			Name:   "report viewer cannot remove work record attachment",
			Method: http.MethodDelete,
			URL:    "/api/work_records/wr260400000003a/attachment",
			Headers: map[string]string{
				"Authorization": reportToken,
			},
			ExpectedStatus: http.StatusForbidden,
			ExpectedContent: []string{
				`"code":"unauthorized"`,
				`"message":"you are not authorized to edit this work record"`,
			},
			TestAppFactory: testutils.SetupTestApp,
		},
		func() tests.ApiScenario {
			body, contentType := makeWorkRecordPDFMultipart(t, "disabled.pdf", uploadContent)
			return tests.ApiScenario{
				Name:   "time editing disabled blocks work record attachment upload",
				Method: http.MethodPost,
				URL:    "/api/work_records/wr260400000010a/attachment",
				Headers: map[string]string{
					"Authorization": ownerToken,
					"Content-Type":  contentType,
				},
				Body:           body,
				ExpectedStatus: http.StatusForbidden,
				ExpectedContent: []string{
					`"Time editing is currently disabled."`,
				},
				TestAppFactory: setupTimeEditingDisabledApp,
			}
		}(),
		{
			Name:   "direct collection patch remains unavailable to regular users",
			Method: http.MethodPatch,
			URL:    "/api/collections/work_records/records/wr260400000010a",
			Body:   strings.NewReader(`{"location":"bypass attempt"}`),
			Headers: map[string]string{
				"Authorization": ownerToken,
				"Content-Type":  "application/json",
			},
			ExpectedStatus: http.StatusForbidden,
			ExpectedContent: []string{
				`"message":"Only superusers can perform this action."`,
			},
			TestAppFactory: testutils.SetupTestApp,
		},
		{
			Name:   "work record detail response includes attachment",
			Method: http.MethodGet,
			URL:    "/api/work_records/wr260400000003a",
			Headers: map[string]string{
				"Authorization": ownerToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"attachment":"paper_copy.pdf"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, _ *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var payload workRecordDetailResponse
				if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
					t.Fatalf("failed decoding work record detail response: %v", err)
				}
				if payload.Attachment != "paper_copy.pdf" {
					t.Fatalf("expected attachment in detail response, got %#v", payload)
				}
			},
		},
		{
			Name:   "work record list response includes attachment",
			Method: http.MethodGet,
			URL:    "/api/work_records",
			Headers: map[string]string{
				"Authorization": ownerToken,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"attachment":"paper_copy.pdf"`,
			},
			TestAppFactory: testutils.SetupTestApp,
			AfterTestFunc: func(t testing.TB, _ *tests.TestApp, res *http.Response) {
				defer res.Body.Close()

				var payload workRecordSearchResponse
				if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
					t.Fatalf("failed decoding work record list response: %v", err)
				}
				for _, item := range payload.Items {
					if item.ID != "wr260400000003a" {
						continue
					}
					if item.Attachment != "paper_copy.pdf" {
						t.Fatalf("expected attachment in list response, got %#v", item)
					}
					return
				}
				t.Fatal("expected seeded work record in list response")
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
