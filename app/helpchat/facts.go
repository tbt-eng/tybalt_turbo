package helpchat

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"

	"tybalt/helpcorpus"
	"tybalt/utilities"
)

// The corpus says how a limit applies. It cannot say what this user's limit
// is, or what an admin has set a switch to, because those are in the database.
// loadFacts reads them for one request and writes them for the prompt.
//
// Every value here is the user's own, or is the same for all users. Most are
// values that the user can already see in Turbo: their own profile, the
// defaults that /api/user_defaults returns, their own row of
// user_po_approver_profile, the app_config collection and the
// expenditure_kinds collection. A project approver's limit is shown to every
// user who picks an approver for a request.
//
// The values of adminProfileFacts are the exception. Only an admin, HR and a
// time off manager can open admin profiles in Turbo, and the owner of Turbo
// decided that Help can tell every user these values of their own admin
// profile. Never add a value of another user, and ask the owner before you add
// a value that the user cannot see in Turbo.
//
// The user's settings go after the cached part of the prompt, so they cost
// input tokens on every call. Keep them short, and write a part only for a user
// it applies to. Turbo's settings are the same for every user and are cached.

// maxFactDivisions limits the list of divisions in the prompt.
const maxFactDivisions = 60

// loadUserFacts returns the section "This user's settings". A part that cannot
// be read is left out and logged: the chat still answers, and says that it
// cannot see the value.
func loadUserFacts(app core.App, uid string, held helpcorpus.Claims, kinds []kind) string {
	var lines []string
	add := func(part string, load func() ([]string, error)) {
		more, err := load()
		if err != nil {
			app.Logger().Error("help chat can't read a part of the user's settings", "part", part, "error", err)
			return
		}
		lines = append(lines, more...)
	}
	add("profile", func() ([]string, error) { return profileFacts(app, uid) })
	add("admin profile", func() ([]string, error) { return adminProfileFacts(app, uid) })
	if held["po_approver"] {
		add("po limits", func() ([]string, error) { return poLimitFacts(app, uid, kinds) })
	}
	if held["project_authorization"] {
		add("project limit", func() ([]string, error) { return projectLimitFacts(app, uid) })
	}
	if len(lines) == 0 {
		return ""
	}
	var out strings.Builder
	out.WriteString("# This user's settings\n\nTurbo read these from its database for this request. They are this user's own.\n")
	writeLines(&out, lines)
	return out.String()
}

// loadTurboFacts returns the section "Turbo's settings today". It is the same
// text for every user, and it changes only when an admin changes a setting,
// so it goes in the cached part of the prompt. Keep its lines in a fixed order.
func loadTurboFacts(app core.App, kinds []kind) string {
	lines, err := turboFacts(app, kinds)
	if err != nil {
		app.Logger().Error("help chat can't read Turbo's settings", "error", err)
	}
	if len(lines) == 0 {
		return ""
	}
	var out strings.Builder
	out.WriteString("# Turbo's settings today\n\nThe same for every user. The name in brackets is the name of the setting in the code.\n")
	writeLines(&out, lines)
	return out.String()
}

func writeLines(out *strings.Builder, lines []string) {
	for _, line := range lines {
		out.WriteString("- " + line + "\n")
	}
}

// kind is one kind of PO, with the amount above which a PO of that kind needs
// a second approval.
type kind struct {
	Name      string  `db:"name"`
	Label     string  `db:"label"`
	Threshold float64 `db:"threshold"`
}

// loadKinds returns the kinds of PO in the order of the screen. With an error
// it returns none, and the lines that need the kinds are left out.
func loadKinds(app core.App) []kind {
	var kinds []kind
	err := app.DB().NewQuery(`
		SELECT name, COALESCE(NULLIF(en_ui_label, ''), name) AS label, COALESCE(second_approval_threshold, 0) AS threshold
		FROM expenditure_kinds
		ORDER BY ui_order, name
	`).All(&kinds)
	if err != nil {
		app.Logger().Error("help chat can't read the kinds of PO", "error", err)
		return nil
	}
	return kinds
}

func profileFacts(app core.App, uid string) ([]string, error) {
	var row struct {
		HasProfile    bool   `db:"has_profile"`
		HasAdmin      bool   `db:"has_admin"`
		Manager       string `db:"manager"`
		Alternate     string `db:"alternate"`
		DivisionCode  string `db:"division_code"`
		DivisionName  string `db:"division_name"`
		Branch        string `db:"branch"`
		Reimbursement bool   `db:"reimbursement"`
	}
	err := app.DB().NewQuery(`
		SELECT
			p.id IS NOT NULL AS has_profile,
			a.id IS NOT NULL AS has_admin,
			TRIM(COALESCE(m.given_name, '') || ' ' || COALESCE(m.surname, '')) AS manager,
			TRIM(COALESCE(am.given_name, '') || ' ' || COALESCE(am.surname, '')) AS alternate,
			COALESCE(d.code, '') AS division_code,
			COALESCE(d.name, '') AS division_name,
			COALESCE(b.name, '') AS branch,
			COALESCE(a.allow_personal_reimbursement, 0) AS reimbursement
		FROM users u
		LEFT JOIN profiles p ON p.uid = u.id
		LEFT JOIN profiles m ON m.uid = p.manager AND p.manager != ''
		LEFT JOIN profiles am ON am.uid = p.alternate_manager AND p.alternate_manager != ''
		LEFT JOIN divisions d ON d.id = p.default_division
		LEFT JOIN admin_profiles a ON a.uid = u.id
		LEFT JOIN branches b ON b.id = a.default_branch
		WHERE u.id = {:uid}
		LIMIT 1
	`).Bind(dbx.Params{"uid": uid}).One(&row)
	if err != nil {
		return nil, err
	}
	var lines []string
	if row.HasProfile {
		lines = append(lines,
			"Manager: "+orNone(row.Manager),
			"Alternate manager: "+orNone(row.Alternate),
			"Default division: "+orNone(strings.TrimSpace(row.DivisionCode+" "+row.DivisionName)),
		)
	}
	if row.HasAdmin {
		lines = append(lines, "Default branch: "+orNone(row.Branch))
		if row.Reimbursement {
			lines = append(lines, "Personal Reimbursement expenses: allowed for this user")
		} else {
			lines = append(lines, "Personal Reimbursement expenses: not allowed for this user")
		}
	}
	return lines, nil
}

// adminProfileFacts gives the values of the user's own admin profile that
// decide what Turbo lets them do, with the labels of the Admin Profiles screen.
// It leaves out the payroll ID and the charge-out rate, which decide nothing
// that a user asks Help about.
func adminProfileFacts(app core.App, uid string) ([]string, error) {
	var rows []struct {
		Salary        bool    `db:"salary"`
		WorkWeek      float64 `db:"work_week_hours"`
		Expected      bool    `db:"time_sheet_expected"`
		OffRotation   bool    `db:"off_rotation_permitted"`
		Untracked     bool    `db:"untracked_time_off"`
		SkipMinTime   string  `db:"skip_min_time_check"`
		OpeningDate   string  `db:"opening_date"`
		OpeningOP     float64 `db:"opening_op"`
		OpeningOV     float64 `db:"opening_ov"`
		InsuranceDate string  `db:"insurance"`
	}
	err := app.DB().NewQuery(`
		SELECT
			COALESCE(salary, 0) AS salary,
			COALESCE(work_week_hours, 0) AS work_week_hours,
			COALESCE(time_sheet_expected, 0) AS time_sheet_expected,
			COALESCE(off_rotation_permitted, 0) AS off_rotation_permitted,
			COALESCE(untracked_time_off, 0) AS untracked_time_off,
			COALESCE(skip_min_time_check, '') AS skip_min_time_check,
			COALESCE(opening_date, '') AS opening_date,
			COALESCE(opening_op, 0) AS opening_op,
			COALESCE(opening_ov, 0) AS opening_ov,
			COALESCE(personal_vehicle_insurance_expiry, '') AS insurance
		FROM admin_profiles
		WHERE uid = {:uid}
		LIMIT 1
	`).Bind(dbx.Params{"uid": uid}).All(&rows)
	if err != nil || len(rows) == 0 {
		return nil, err
	}
	row := rows[0]
	yes := func(value bool) string {
		if value {
			return "Yes"
		}
		return "No"
	}
	return []string{"The user's own admin profile: " + strings.Join([]string{
		"Salary " + yes(row.Salary),
		"Work Week Hours " + number(row.WorkWeek),
		"Time Sheet Expected " + yes(row.Expected),
		"Off Rotation Permitted " + yes(row.OffRotation),
		"Untracked Time Off " + yes(row.Untracked),
		"Skip Min Time Check " + orNone(row.SkipMinTime),
		"Opening Date " + orNone(row.OpeningDate),
		"Opening OP " + number(row.OpeningOP),
		"Opening OV " + number(row.OpeningOV),
		"Personal Vehicle Insurance Expiry " + orNone(row.InsuranceDate),
	}, "; ")}, nil
}

func orNone(value string) string {
	if value == "" {
		return "none"
	}
	return value
}

// poLimitFacts gives a PO approver's limit for each kind of PO, and the
// divisions of the POs that they can approve.
func poLimitFacts(app core.App, uid string, kinds []kind) ([]string, error) {
	var rows []struct {
		Capital        float64 `db:"capital"`
		Project        float64 `db:"project"`
		Sponsorship    float64 `db:"sponsorship"`
		StaffAndSocial float64 `db:"staff_and_social"`
		MediaAndEvent  float64 `db:"media_and_event"`
		Computer       float64 `db:"computer"`
		Divisions      string  `db:"divisions"`
	}
	err := app.DB().NewQuery(`
		SELECT
			COALESCE(pap.max_amount, 0) AS capital,
			COALESCE(pap.project_max, 0) AS project,
			COALESCE(pap.sponsorship_max, 0) AS sponsorship,
			COALESCE(pap.staff_and_social_max, 0) AS staff_and_social,
			COALESCE(pap.media_and_event_max, 0) AS media_and_event,
			COALESCE(pap.computer_max, 0) AS computer,
			COALESCE(pap.divisions, '[]') AS divisions
		FROM user_claims uc
		JOIN claims c ON c.id = uc.cid AND c.name = 'po_approver'
		LEFT JOIN po_approver_props pap ON pap.user_claim = uc.id
		WHERE uc.uid = {:uid}
		LIMIT 1
	`).Bind(dbx.Params{"uid": uid}).All(&rows)
	if err != nil || len(rows) == 0 {
		return nil, err
	}
	row := rows[0]
	limits := map[string]float64{
		utilities.ExpenditureKindNameCapital:        row.Capital,
		utilities.ExpenditureKindNameProject:        row.Project,
		utilities.ExpenditureKindNameSponsorship:    row.Sponsorship,
		utilities.ExpenditureKindNameStaffAndSocial: row.StaffAndSocial,
		utilities.ExpenditureKindNameMediaAndEvent:  row.MediaAndEvent,
		utilities.ExpenditureKindNameComputer:       row.Computer,
	}
	var parts []string
	for _, k := range kinds {
		if limit, ok := limits[k.Name]; ok {
			parts = append(parts, k.Label+" "+dollars(limit))
		}
	}
	if len(parts) == 0 {
		return nil, nil
	}
	lines := []string{"PO approval limit for each kind of PO: " + strings.Join(parts, "; ") +
		". A limit is the largest PO amount for which this user's approval is final. With a limit of $0, this user cannot approve a PO of that kind." +
		" For a larger PO that needs a second approval, this user can give the first approval and not the final one."}

	var ids []string
	if err := json.Unmarshal([]byte(row.Divisions), &ids); err != nil {
		return lines, nil
	}
	if len(ids) == 0 {
		return append(lines, "Divisions of the POs that this user can approve: all divisions"), nil
	}
	var codes []string
	err = app.DB().NewQuery(`
		SELECT code FROM divisions WHERE id IN (SELECT value FROM json_each({:ids})) ORDER BY code
	`).Bind(dbx.Params{"ids": row.Divisions}).Column(&codes)
	if err != nil {
		return lines, err
	}
	list := strings.Join(codes, ", ")
	if len(codes) > maxFactDivisions {
		list = fmt.Sprintf("%s and %d more", strings.Join(codes[:maxFactDivisions], ", "), len(codes)-maxFactDivisions)
	}
	return append(lines, fmt.Sprintf("Divisions of the POs that this user can approve (%d): %s", len(codes), list)), nil
}

// projectLimitFacts gives the largest project value that a project approver
// can approve.
func projectLimitFacts(app core.App, uid string) ([]string, error) {
	var limits []float64
	err := app.DB().NewQuery(`
		SELECT COALESCE(pap.max_value, 0)
		FROM user_claims uc
		JOIN claims c ON c.id = uc.cid AND c.name = 'project_authorization'
		LEFT JOIN project_authorization_props pap ON pap.user_claim = uc.id
		WHERE uc.uid = {:uid}
		LIMIT 1
	`).Bind(dbx.Params{"uid": uid}).Column(&limits)
	if err != nil || len(limits) == 0 {
		return nil, err
	}
	if limits[0] <= 0 {
		return []string{"Limit for project requests and project value changes: none is set, so this user cannot approve a project value"}, nil
	}
	return []string{"Limit for project requests and project value changes, as the largest project value that this user can approve: " +
		dollars(limits[0]) + ". This is not a limit for POs."}, nil
}

// turboFacts gives the switches and amounts of app_config, and the second
// approval amount of each kind of PO. The functions of utilities apply the
// defaults, so the model gets the value that Turbo uses.
func turboFacts(app core.App, kinds []kind) ([]string, error) {
	var lines []string
	// A switch that is on is the usual state, so only one that is off is written.
	for _, f := range []struct {
		label, name string
		read        func(core.App) (bool, error)
	}{
		{"creating and editing time entries", "time.create_edit", utilities.IsTimeEditingEnabled},
		{"creating and editing expenses, POs and vendors", "expenses.create_edit_absorb", utilities.IsExpensesEditingEnabled},
		{"creating and editing jobs", "jobs.create_edit_absorb", utilities.IsJobsEditingEnabled},
	} {
		on, err := f.read(app)
		if err != nil {
			return lines, err
		}
		if !on {
			lines = append(lines, fmt.Sprintf("Turned off for everyone now: %s (%s)", f.label, f.name))
		}
	}
	required, err := utilities.IsJobRequestApprovalRequired(app)
	if err != nil {
		return lines, err
	}
	if required {
		lines = append(lines, "A new project needs an approved project request (jobs.require_job_request_approval): yes")
	} else {
		lines = append(lines, "A new project needs an approved project request (jobs.require_job_request_approval): no")
	}

	lines = append(lines, fmt.Sprintf("An expense of this amount or more needs a PO (no_po_expense_limit): %s", dollars(utilities.GetNoPOExpenseLimit(app))))
	excess := utilities.GetPOExpenseExcessConfig(app)
	which := "the lesser"
	if excess.Mode == "greater_of" {
		which = "the greater"
	}
	lines = append(lines, fmt.Sprintf("The expenses of a PO can go over its total by %s of %s%% and %s CAD (po_expense_allowed_excess)",
		which, number(math.Round(excess.Percent*10000)/100), dollars(excess.Value)))
	lines = append(lines, fmt.Sprintf("Hours that a PO waits for its second approval before the timeout (second_stage_timeout_hours): %s",
		number(utilities.GetPurchaseOrderSecondStageTimeoutHours(app))))

	var need, none []string
	for _, k := range kinds {
		if k.Threshold > 0 {
			need = append(need, k.Label+" "+dollars(k.Threshold))
		} else {
			none = append(none, k.Label)
		}
	}
	if len(need) > 0 {
		line := "A PO needs a second approval when its amount is above (second_approval_threshold): " + strings.Join(need, "; ")
		if len(none) > 0 {
			line += ". No second approval for: " + strings.Join(none, ", ")
		}
		lines = append(lines, line)
	}
	return lines, nil
}

// dollars writes an amount as $1,234 or $1,234.50.
func dollars(amount float64) string {
	text := strconv.FormatFloat(amount, 'f', 2, 64)
	whole, cents, _ := strings.Cut(text, ".")
	negative := strings.HasPrefix(whole, "-")
	whole = strings.TrimPrefix(whole, "-")
	var grouped []byte
	for i := range len(whole) {
		if i > 0 && (len(whole)-i)%3 == 0 {
			grouped = append(grouped, ',')
		}
		grouped = append(grouped, whole[i])
	}
	out := "$" + string(grouped)
	if cents != "00" {
		out += "." + cents
	}
	if negative {
		out = "-" + out
	}
	return out
}

// number writes 24 as "24" and 2.5 as "2.5".
func number(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}
