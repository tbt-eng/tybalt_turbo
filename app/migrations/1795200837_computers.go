package migrations

import (
	"database/sql"
	"errors"
	"slices"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

// Add the collections that radiator login reports fill (issue 211).
//
// Radiator is a script that runs on each Windows computer at sign-in. It
// POSTs one report about the computer and the user to /api/radiator, which
// checks a machine secret with the new radiator role.
//
//   - computers: one row per computer, keyed by a slug of the serial and the
//     manufacturer, updated in place by each valid report. assigned_to,
//     assigned_time and retired are not written by reports.
//   - computer_logins: one row per valid report.
//   - computer_raw_logins: reports that failed validation, with the errors. A
//     valid report from the same computer name deletes them.
//
// admin_profiles gets user_source_anchor, the AD mS-DS-ConsistencyGuid of the
// person as 32 lowercase hex characters. The endpoint matches a report to a
// user on it first, then on email. The unique index skips empty values.
//
// Only server code writes the three collections: they have no create, update
// or delete rule. A holder of the new computers claim or of admin can list and
// view them.
//
// See descriptions/computers.md.

const (
	computersCollectionID         = "pbc_computers"
	computerLoginsCollectionID    = "pbc_comp_logins"
	computerRawLoginsCollectionID = "pbc_comp_raw_logins"

	computersClaimID          = "computersclaim1"
	computersClaimName        = "computers"
	computersClaimDescription = "Can see the computers that report at sign-in, the log of their sign-ins and the reports that were rejected."

	radiatorMachineSecretRole = "radiator"

	adminProfilesUserSourceAnchorIndex = "idx_admin_profiles_user_source_anchor"

	computersReadRule = "@request.auth.id != '' && (\n  @request.auth.user_claims_via_uid.cid.name ?= 'computers' ||\n  @request.auth.user_claims_via_uid.cid.name ?= 'admin'\n)"
)

func init() {
	m.Register(func(app core.App) error {
		if err := ensureComputersClaim(app); err != nil {
			return err
		}
		if err := setRadiatorMachineSecretRole(app, true); err != nil {
			return err
		}
		if err := addAdminProfilesUserSourceAnchor(app); err != nil {
			return err
		}

		users, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}

		computers := core.NewBaseCollection("computers")
		computers.Id = computersCollectionID
		computers.ListRule = types.Pointer(computersReadRule)
		computers.ViewRule = types.Pointer(computersReadRule)
		computers.Fields.Add(
			// The key: serial and manufacturer, cleaned (see makeSlug in
			// routes/radiator.go).
			&core.TextField{Name: "slug", Required: true},

			// The computer, as last reported.
			&core.TextField{Name: "serial", Required: true},
			&core.TextField{Name: "mfg", Required: true},
			&core.TextField{Name: "model"},
			&core.TextField{Name: "computer_name"},
			&core.NumberField{Name: "system_type", OnlyInt: true},
			&core.NumberField{Name: "ram", OnlyInt: true},
			&core.TextField{Name: "os_version"},
			&core.NumberField{Name: "os_sku", OnlyInt: true},
			&core.TextField{Name: "os_arch"},
			&core.TextField{Name: "boot_drive"},
			&core.TextField{Name: "boot_drive_fs"},
			&core.NumberField{Name: "boot_drive_cap", OnlyInt: true},
			&core.NumberField{Name: "boot_drive_free", OnlyInt: true},
			&core.NumberField{Name: "radiator_version", OnlyInt: true},
			&core.JSONField{Name: "network_config"},

			// The user of the last report. last_user is the matched user, if any.
			&core.TextField{Name: "user_given_name"},
			&core.TextField{Name: "user_surname"},
			&core.TextField{Name: "email"},
			&core.TextField{Name: "upn"},
			&core.TextField{Name: "user_source_anchor"},
			&core.RelationField{Name: "last_user", CollectionId: users.Id, MaxSelect: 1},
			&core.DateField{Name: "last_login"},

			// Set by people, never by a report.
			&core.RelationField{Name: "assigned_to", CollectionId: users.Id, MaxSelect: 1},
			&core.DateField{Name: "assigned_time"},
			&core.DateField{Name: "retired"},

			&core.AutodateField{Name: "created", OnCreate: true},
			&core.AutodateField{Name: "updated", OnCreate: true, OnUpdate: true},
		)
		computers.AddIndex("idx_computers_slug", true, "`slug`", "")
		computers.AddIndex("idx_computers_computer_name", false, "`computer_name`", "")
		if err := app.Save(computers); err != nil {
			return err
		}

		logins := core.NewBaseCollection("computer_logins")
		logins.Id = computerLoginsCollectionID
		logins.ListRule = types.Pointer(computersReadRule)
		logins.ViewRule = types.Pointer(computersReadRule)
		logins.Fields.Add(
			&core.RelationField{Name: "computer", CollectionId: computers.Id, MaxSelect: 1, Required: true, CascadeDelete: true},
			&core.TextField{Name: "given_name"},
			&core.TextField{Name: "surname"},
			&core.TextField{Name: "user_source_anchor"},
			&core.RelationField{Name: "user", CollectionId: users.Id, MaxSelect: 1},
			&core.AutodateField{Name: "created", OnCreate: true},
		)
		logins.AddIndex("idx_computer_logins_created", false, "`created`", "")
		logins.AddIndex("idx_computer_logins_computer_created", false, "`computer`, `created`", "")
		if err := app.Save(logins); err != nil {
			return err
		}

		rawLogins := core.NewBaseCollection("computer_raw_logins")
		rawLogins.Id = computerRawLoginsCollectionID
		rawLogins.ListRule = types.Pointer(computersReadRule)
		rawLogins.ViewRule = types.Pointer(computersReadRule)
		rawLogins.Fields.Add(
			&core.JSONField{Name: "report"},
			&core.TextField{Name: "computer_name"},
			&core.TextField{Name: "error"},
			&core.AutodateField{Name: "created", OnCreate: true},
		)
		rawLogins.AddIndex("idx_computer_raw_logins_computer_name", false, "`computer_name`", "")
		rawLogins.AddIndex("idx_computer_raw_logins_created", false, "`created`", "")
		return app.Save(rawLogins)
	}, func(app core.App) error {
		for _, id := range []string{computerRawLoginsCollectionID, computerLoginsCollectionID, computersCollectionID} {
			collection, err := app.FindCollectionByNameOrId(id)
			if err != nil {
				continue
			}
			if err := app.Delete(collection); err != nil {
				return err
			}
		}
		if err := removeAdminProfilesUserSourceAnchor(app); err != nil {
			return err
		}
		if err := setRadiatorMachineSecretRole(app, false); err != nil {
			return err
		}
		record, err := app.FindRecordById("claims", computersClaimID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}
			return err
		}
		return app.Delete(record)
	})
}

func ensureComputersClaim(app core.App) error {
	existing, err := app.FindFirstRecordByFilter("claims", "name={:name}", dbx.Params{"name": computersClaimName})
	if err == nil && existing != nil {
		existing.Set("description", computersClaimDescription)
		return app.Save(existing)
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	collection, err := app.FindCollectionByNameOrId("claims")
	if err != nil {
		return err
	}

	record := core.NewRecord(collection)
	record.Set("id", computersClaimID)
	record.Set("name", computersClaimName)
	record.Set("description", computersClaimDescription)
	return app.Save(record)
}

// setRadiatorMachineSecretRole adds radiator to, or removes it from, the
// values of machine_secrets.role.
func setRadiatorMachineSecretRole(app core.App, present bool) error {
	secrets, err := app.FindCollectionByNameOrId("machine_secrets")
	if err != nil {
		return err
	}
	role, ok := secrets.Fields.GetByName("role").(*core.SelectField)
	if !ok {
		return errors.New("machine_secrets.role is not a select field")
	}
	has := slices.Contains(role.Values, radiatorMachineSecretRole)
	switch {
	case present && !has:
		role.Values = append(role.Values, radiatorMachineSecretRole)
	case !present && has:
		role.Values = slices.DeleteFunc(role.Values, func(v string) bool { return v == radiatorMachineSecretRole })
	default:
		return nil
	}
	return app.Save(secrets)
}

func addAdminProfilesUserSourceAnchor(app core.App) error {
	profiles, err := app.FindCollectionByNameOrId("admin_profiles")
	if err != nil {
		return err
	}
	if profiles.Fields.GetByName("user_source_anchor") == nil {
		profiles.Fields.Add(&core.TextField{Name: "user_source_anchor", Max: 32, Pattern: "^[0-9a-f]{32}$"})
	}
	profiles.AddIndex(adminProfilesUserSourceAnchorIndex, true, "`user_source_anchor`", "`user_source_anchor` != ''")
	return app.Save(profiles)
}

func removeAdminProfilesUserSourceAnchor(app core.App) error {
	profiles, err := app.FindCollectionByNameOrId("admin_profiles")
	if err != nil {
		return err
	}
	profiles.RemoveIndex(adminProfilesUserSourceAnchorIndex)
	profiles.Fields.RemoveByName("user_source_anchor")
	return app.Save(profiles)
}
