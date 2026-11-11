package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// An admin reads the description of a claim on the Claims screen before they
// give the claim to a person, and Help chat reads it to say who can do a task.
// Many descriptions were written when the claim did one thing, and the code
// now checks the claim in many places: admin said only "Can delete jobs".
// Each new description says what the code and the collection rules let a
// holder do today. A claim that is not in the database is left alone.
var claimDescriptions = []struct{ name, before, after string }{
	{"absorb",
		"Can absorb (merge) many existing records in a collection to a single record to handle duplicates or dirty data",
		"Can merge duplicate clients, client contacts and vendors into one record, undo the latest merge of each type, and see the Absorb Actions list."},
	{"accounting",
		"Can cancel purchase_orders records",
		"Not in use. Turbo does not check this claim."},
	{"admin",
		"Can delete jobs",
		"System administrator: can edit admin profiles, assign claims and approval limits, see committed time sheets and expenses, renumber or delete jobs, and manage currencies, rate sheet rates, machine secrets and attachment audits."},
	{"book_keeper",
		"Can create eligible PO-linked expenses on behalf of the purchase order owner",
		"Can enter On Account or corporate card expenses against another person's Active PO for that person, becomes the approver of those expenses, and can reuse a receipt that another expense already has."},
	{"busdev",
		"Business Development Leads. These are users that are branch managers or executives who can be responsible for business development with a client. Users with this claim can be assigned to the business_development_lead field of the clients collection",
		"Marks a person as a business development lead, so that they can be picked as the lead of a client. It gives the holder no other access."},
	{"commit",
		"Can commit records",
		"Can commit approved time sheets, expenses and time amendments, reject approved ones, see the commit queue and the tracking screens, and ignore missing time sheets."},
	{"corporate_branch",
		"Allows assigning and using the `Corporate` branch for time, expense, and po entries",
		"Lets the holder choose the Corporate branch on their time entries, expenses and POs, and have Corporate as their default branch."},
	{"hr",
		"Can edit specific fields on admin_profiles records",
		"Can view all admin profiles and edit their HR fields: active status, job title, phone, default branch, charge-out rate, salary, payroll ID, time off opening balances and time sheet settings."},
	{"it",
		"Can manage identity repair fields and authorized providers for users",
		"Can view each user's sign-in identity, set the legacy user ID, remove a linked sign-in provider, and make users who are not admins active or inactive."},
	{"job",
		"Create and update jobs",
		"Procurement: can create and edit jobs, clients, contacts and rate sheets, close projects, change the client of a project, create jobs from approved project requests, and propose project value changes."},
	{"kpi",
		"Can view management KPI reports",
		"Can view the KPI reports, the Branch and Division WIP reports, and Branch Hours for every branch."},
	{"legacy_po_create_update",
		"Can create and update legacy manually-entered purchase orders while the transition flag is enabled",
		"Can create and edit legacy POs that are entered by hand, and see them in any status, while the legacy PO setting is on."},
	{"payables_admin",
		"Can manage purchase orders (cancel active POs with no expenses, close active Cumulative/Recurring POs with committed expenses, convert One-Time POs to Cumulative), create OnAccount expenses above the no-PO limit without a PO, create/update/delete vendors, and edit the outstanding_balance field on jobs",
		"Payables staff: can cancel an Active PO that has no expenses, close an Active Cumulative or Recurring PO, convert a One-Time PO to Cumulative, settle foreign-currency expenses, create, edit and delete vendors, and enter On Account expenses over the no-PO limit."},
	{"po_approver",
		"Holders can approve purchase_orders records. payload contains a JSON object which will restrict the type of purchase_orders record that can be approved. Properties are \"max_amount\" (required, number) and \"divisions\" (not required, list of strings that are division ids)",
		"Can be chosen as a PO approver and approve POs up to the limits that an admin sets for them, for each kind of PO and for their divisions. With no limits, the claim does nothing."},
	{"project_authorization",
		"Approve full project values within the assigned limit",
		"Can be chosen to approve project requests and project value changes when the full project value is within the holder's limit. An admin must set that limit, or the claim does nothing."},
	{"rate_sheet_revise",
		"Holders can create revisions of a rate_sheets record (revision > 0). Holders may also change the rate_sheet value of a jobs record that has a non-blank rate_sheet value.",
		"Can create a new revision of a rate sheet, and change the rate sheet of a job that already has one. To edit the job, the holder also needs the job claim."},
	{"report",
		"Broad read-access to data for reporting",
		"Read-only access for reports: the payroll, weekly and payables reports, all committed time sheets and expenses, the tracking screens, everyone's time off, time amendments, work records, PO approver limits, and Closed or Cancelled POs."},
	{"super_po",
		"Create purchase orders for other staff from a template and a CSV upload",
		"Can create a group of POs for other staff from one template PO and a CSV upload, and see every PO of a group in any status."},
	{"tame",
		"Can CRUD uncommitted time_amendments records and read committed ones",
		"Can create, edit and delete time amendments for any employee until they are committed, and view all amendments. To commit one, the holder needs the commit claim."},
	{"tapr",
		"Approve Time Sheets (manager)",
		"Marks a person as a manager: they can be set as a person's manager or alternate manager, and as manager they approve that person's time sheets and expenses."},
	{"time",
		"CRUD for time_entries, bundling time_sheets where records belong to the user",
		"Basic staff access, given at the first sign-in: enter your own time and bundle it into time sheets, create and submit project requests, and add vendors and client invoicing profiles."},
	{"time_off_manager",
		"Can edit opening_date opening_op and opening_ov on admin_profiles records",
		"Can view all admin profiles and edit each person's time off opening date, opening balances and payroll ID."},
	{"tt",
		"Create Time Types",
		"Can add time types. The Time Types screen is in the menu only for an admin, so give this claim together with admin."},
	{"work_record",
		"Can audit time entries by work record across the organization",
		"Can search work records and see every employee's time entries on each one."},
}

func init() {
	m.Register(func(app core.App) error {
		for _, claim := range claimDescriptions {
			if _, err := app.DB().NewQuery(`
				UPDATE claims
				SET description = {:description}, updated = strftime('%Y-%m-%d %H:%M:%fZ', 'now')
				WHERE name = {:name}
			`).Bind(dbx.Params{"name": claim.name, "description": claim.after}).Execute(); err != nil {
				return err
			}
		}
		return nil
	}, func(app core.App) error {
		// Put back a description only if nobody changed it after this migration.
		for _, claim := range claimDescriptions {
			if _, err := app.DB().NewQuery(`
				UPDATE claims
				SET description = {:before}, updated = strftime('%Y-%m-%d %H:%M:%fZ', 'now')
				WHERE name = {:name} AND description = {:after}
			`).Bind(dbx.Params{"name": claim.name, "before": claim.before, "after": claim.after}).Execute(); err != nil {
				return err
			}
		}
		return nil
	})
}
