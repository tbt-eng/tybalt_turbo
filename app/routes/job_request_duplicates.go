package routes

import (
	_ "embed" // for go:embed
	"net/http"
	"slices"
	"sort"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"tybalt/utilities"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// The duplicate-work review on a project request. Procurement scans the list
// before creating the job, so it holds what can plausibly be the same work:
//
//   - exact matches, whatever their age or status: the same proposal, client
//     PO or client reference, the request's parent, or a sibling subjob;
//   - the client's other projects that are open (not Closed, not Cancelled)
//     and recent (an award date or a time entry within
//     JOB_REQUEST_DUPLICATE_RECENT_DAYS).
//
// Location and description never list a project by themselves; they are
// extra reasons on a listed row. The response counts the client's projects
// that the rules leave out, so the page can say what the list does not hold.

//go:embed job_request_duplicates.sql
var jobRequestDuplicatesQuery string

// Tiers order the list: exact matches first, then the client's projects with
// a similar description, then the client's other recent projects.
const (
	duplicateTierExact   = "exact"
	duplicateTierSimilar = "similar"
	duplicateTierClient  = "client"
)

var duplicateTierRank = map[string]int{duplicateTierExact: 0, duplicateTierSimilar: 1, duplicateTierClient: 2}

const jobRequestDuplicatesNowStoreKey = "tybalt.routes.jobRequestDuplicatesNow"

func jobRequestDuplicatesNow(app core.App) time.Time {
	if override, ok := app.Store().Get(jobRequestDuplicatesNowStoreKey).(func() time.Time); ok && override != nil {
		return override().UTC()
	}
	return time.Now().UTC()
}

// SetJobRequestDuplicatesNowForTest fixes the reference date of the
// duplicate-work review for one app instance, so tests do not depend on how
// far fixture dates are from today. nil restores the real clock.
func SetJobRequestDuplicatesNowForTest(app core.App, now func() time.Time) {
	if now == nil {
		app.Store().Remove(jobRequestDuplicatesNowStoreKey)
		return
	}
	app.Store().Set(jobRequestDuplicatesNowStoreKey, now)
}

// duplicateSubject holds the fields that the review compares, for the request
// under review and for each job or other request.
type duplicateSubject struct {
	ID              string
	Client          string
	Proposal        string
	Parent          string
	ClientPO        string
	ClientReference string
	Location        string
	Description     string
}

func requestDuplicateSubject(r *core.Record) duplicateSubject {
	return duplicateSubject{
		ID:              r.Id,
		Client:          r.GetString("client"),
		Proposal:        r.GetString("proposal"),
		Parent:          r.GetString("parent"),
		ClientPO:        r.GetString("client_po"),
		ClientReference: r.GetString("client_reference_number"),
		Location:        r.GetString("location"),
		Description:     r.GetString("description"),
	}
}

type duplicateJobRow struct {
	ID              string `db:"id" json:"id"`
	Number          string `db:"number" json:"number"`
	Description     string `db:"description" json:"description"`
	Client          string `db:"client" json:"client"`
	ClientName      string `db:"client_name" json:"client_name"`
	Status          string `db:"status" json:"status"`
	Manager         string `db:"manager" json:"manager"`
	ManagerName     string `db:"manager_name" json:"manager_name"`
	Parent          string `db:"parent" json:"parent"`
	ParentNumber    string `db:"parent_number" json:"parent_number"`
	Proposal        string `db:"proposal" json:"-"`
	Location        string `db:"location" json:"-"`
	ClientPO        string `db:"client_po" json:"-"`
	ClientReference string `db:"client_reference_number" json:"-"`
	AwardDate       string `db:"project_award_date" json:"-"`
	LastTime        string `db:"last_time" json:"-"`
	// Filled by classifyDuplicateJobs.
	Reasons      []string `db:"-" json:"reasons"`
	Tier         string   `db:"-" json:"tier"`
	LastActivity string   `db:"-" json:"last_activity"`
}

func (j duplicateJobRow) subject() duplicateSubject {
	return duplicateSubject{ID: j.ID, Client: j.Client, Proposal: j.Proposal, Parent: j.Parent, ClientPO: j.ClientPO, ClientReference: j.ClientReference, Location: j.Location, Description: j.Description}
}

// sameCode compares client POs and client references: blank never matches,
// and case and surrounding spaces do not count.
func sameCode(a, b string) bool {
	a = strings.ToLower(strings.TrimSpace(a))
	return a != "" && a == strings.ToLower(strings.TrimSpace(b))
}

func sameID(a, b string) bool {
	return a != "" && a == b
}

// duplicateReasons compares the request under review with one job or other
// request. exact is true for a match that lists a row whatever its age or
// status: the same proposal, client PO or client reference, the request's
// parent, or the same parent. The description is compared only for the same
// client, and location is only ever a supporting reason.
func duplicateReasons(r, other duplicateSubject) (reasons []string, exact, sameClient bool) {
	reasons = []string{}
	sameClient = sameID(r.Client, other.Client)
	if sameClient {
		reasons = append(reasons, "Same client")
	}
	if sameID(r.Proposal, other.Proposal) {
		reasons = append(reasons, "Same proposal")
		exact = true
	}
	if sameID(r.Parent, other.Parent) {
		reasons = append(reasons, "Same parent")
		exact = true
	}
	if sameID(strings.TrimSpace(r.Location), strings.TrimSpace(other.Location)) {
		reasons = append(reasons, "Same location")
	}
	if sameCode(r.ClientPO, other.ClientPO) {
		reasons = append(reasons, "Same client PO")
		exact = true
	}
	if sameCode(r.ClientReference, other.ClientReference) {
		reasons = append(reasons, "Same client reference")
		exact = true
	}
	if sameClient && similarDescriptions(r.Description, other.Description) {
		reasons = append(reasons, "Similar description")
	}
	if sameID(r.Parent, other.ID) {
		reasons = append(reasons, "Parent project")
		exact = true
	}
	return reasons, exact, sameClient
}

// descriptionStopwords are words of three or more letters that say nothing
// about which work a description means. Shorter words are dropped anyway.
var descriptionStopwords = map[string]bool{
	"and": true, "the": true, "for": true, "with": true, "from": true, "into": true,
	"onto": true, "over": true, "under": true, "per": true, "via": true, "its": true,
	"all": true, "new": true, "project": true, "projects": true, "phase": true,
	"work": true, "works": true, "service": true, "services": true,
}

// descriptionWords returns the distinct meaningful words of a description:
// lower case, split on anything but letters and digits, without stopwords and
// without words of fewer than three characters. A two-character word with a
// digit, such as "17" in "Highway 17", is kept.
func descriptionWords(description string) map[string]bool {
	words := map[string]bool{}
	fields := strings.FieldsFunc(strings.ToLower(description), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	for _, word := range fields {
		length := utf8.RuneCountInString(word)
		if length < 2 || (length == 2 && !strings.ContainsFunc(word, unicode.IsDigit)) {
			continue
		}
		if descriptionStopwords[word] {
			continue
		}
		words[word] = true
	}
	return words
}

// similarDescriptions is true when the two descriptions share at least half
// of the meaningful words of the shorter one, and at least two words. When
// the shorter one has a single meaningful word, that word must be shared.
func similarDescriptions(a, b string) bool {
	small, large := descriptionWords(a), descriptionWords(b)
	if len(small) > len(large) {
		small, large = large, small
	}
	if len(small) == 0 {
		return false
	}
	shared := 0
	for word := range small {
		if large[word] {
			shared++
		}
	}
	if len(small) == 1 {
		return shared == 1
	}
	return shared >= 2 && shared*2 >= len(small)
}

// isoDay returns the YYYY-MM-DD part of a date string, or "" if it has none.
func isoDay(value string) string {
	value = strings.TrimSpace(value)
	if len(value) < 10 {
		return ""
	}
	if _, err := time.Parse("2006-01-02", value[:10]); err != nil {
		return ""
	}
	return value[:10]
}

// lastActivity is the later of a project's award date and its latest time
// entry date, or "" when it has neither.
func lastActivity(awardDate, lastTime string) string {
	award, entry := isoDay(awardDate), isoDay(lastTime)
	if award > entry {
		return award
	}
	return entry
}

// recentProject is true when the project's award date or its latest time
// entry is on or after the start of the recent window that ends at now. An
// award date in the future counts as recent.
func recentProject(awardDate, lastTime string, now time.Time) bool {
	activity := lastActivity(awardDate, lastTime)
	if activity == "" {
		return false
	}
	cutoff := now.UTC().AddDate(0, 0, -utilities.JOB_REQUEST_DUPLICATE_RECENT_DAYS).Format("2006-01-02")
	return activity >= cutoff
}

// openProject is true unless the project is Closed or Cancelled.
func openProject(status string) bool {
	return status != "Closed" && status != "Cancelled"
}

// classifyDuplicateJobs decides which candidate jobs the review lists, in
// which tier, and in which order. hidden counts the client's projects that
// are left out because they are closed, cancelled or not recent.
func classifyDuplicateJobs(r duplicateSubject, candidates []duplicateJobRow, now time.Time) (listed []duplicateJobRow, hidden int) {
	listed = []duplicateJobRow{}
	for _, job := range candidates {
		if strings.HasPrefix(job.Number, "P") {
			continue
		}
		reasons, exact, sameClient := duplicateReasons(r, job.subject())
		switch {
		case exact:
			job.Tier = duplicateTierExact
		case sameClient && openProject(job.Status) && recentProject(job.AwardDate, job.LastTime, now):
			job.Tier = duplicateTierClient
			if slices.Contains(reasons, "Similar description") {
				job.Tier = duplicateTierSimilar
			}
		case sameClient:
			hidden++
			continue
		default:
			continue
		}
		job.Reasons = reasons
		job.LastActivity = lastActivity(job.AwardDate, job.LastTime)
		listed = append(listed, job)
	}
	sort.SliceStable(listed, func(i, k int) bool {
		a, b := listed[i], listed[k]
		if duplicateTierRank[a.Tier] != duplicateTierRank[b.Tier] {
			return duplicateTierRank[a.Tier] < duplicateTierRank[b.Tier]
		}
		if a.LastActivity != b.LastActivity {
			return a.LastActivity > b.LastActivity
		}
		return a.Number > b.Number
	})
	return listed, hidden
}

// duplicateRequestTier places another pending request. It is listed for the
// same client or an exact match; pending requests are all current, so no date
// rule applies. The second result is false when it is not listed.
func duplicateRequestTier(reasons []string, exact, sameClient bool) (string, bool) {
	switch {
	case exact:
		return duplicateTierExact, true
	case sameClient && slices.Contains(reasons, "Similar description"):
		return duplicateTierSimilar, true
	case sameClient:
		return duplicateTierClient, true
	}
	return "", false
}

func requestDuplicatesHandler(app core.App) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		r, err := loadVisibleRequest(app, e.Request.PathValue("id"), e.Auth)
		if err != nil {
			return err
		}
		subject := requestDuplicateSubject(r)
		candidates := []duplicateJobRow{}
		if err := app.DB().NewQuery(jobRequestDuplicatesQuery).Bind(dbx.Params{
			"client":           subject.Client,
			"proposal":         subject.Proposal,
			"parent":           subject.Parent,
			"client_po":        strings.ToLower(strings.TrimSpace(subject.ClientPO)),
			"client_reference": strings.ToLower(strings.TrimSpace(subject.ClientReference)),
		}).All(&candidates); err != nil {
			return err
		}
		items, hidden := classifyDuplicateJobs(subject, candidates, jobRequestDuplicatesNow(app))

		pending := []map[string]any{}
		hasJob, err := utilities.HasClaim(app, e.Auth, "job")
		if err != nil {
			return err
		}
		if hasJob {
			requests, err := app.FindAllRecords("job_requests")
			if err != nil {
				return err
			}
			// Newest edit first within a tier.
			sort.SliceStable(requests, func(i, k int) bool {
				return requests[i].GetDateTime("last_edited").Time().After(requests[k].GetDateTime("last_edited").Time())
			})
			for _, other := range requests {
				if other.Id == r.Id {
					continue
				}
				reasons, exact, sameClient := duplicateReasons(subject, requestDuplicateSubject(other))
				tier, ok := duplicateRequestTier(reasons, exact, sameClient)
				if !ok {
					continue
				}
				row := recordPublic(other)
				row["reasons"] = reasons
				row["tier"] = tier
				pending = append(pending, row)
			}
			sort.SliceStable(pending, func(i, k int) bool {
				return duplicateTierRank[pending[i]["tier"].(string)] < duplicateTierRank[pending[k]["tier"].(string)]
			})
		}
		return e.JSON(http.StatusOK, map[string]any{"jobs": items, "hidden_client_jobs": hidden, "requests": pending})
	}
}
