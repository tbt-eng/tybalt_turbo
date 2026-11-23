package migrations

import (
	"database/sql"
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// Retire the business development lead of a client, without losing the data.
//
// The business does not want to record a lead for now, so the UI no longer
// offers the field. This migration does two things.
//
//  1. It makes clients.business_development_lead optional. A schema `required`
//     flag is checked on every save, so it would reject any save of a client
//     that has no lead, which the UI can no longer set. Existing values stay.
//     An optional relation also lets PocketBase clear the reference when the
//     referenced user is deleted, where a required one refused the delete.
//
//  2. It deletes the busdev claim and takes it from everyone who holds it. The
//     claim only marked who could be picked as a lead. If the lead comes back
//     it will not use this claim.
//
// The down migration restores the flag and the claim, but not who held the
// claim. PocketBase doesn't recheck existing rows when a collection's schema
// is saved, so it succeeds even with blank leads present, but those clients
// can no longer be saved afterwards.
const (
	busdevClaimID          = "v748dp232p1onmh"
	busdevClaimName        = "busdev"
	busdevClaimDescription = "Marks a person as a business development lead, so that they can be picked as the lead of a client. It gives the holder no other access."
)

func init() {
	m.Register(func(app core.App) error {
		if err := setClientsBusdevLeadRequired(app, false); err != nil {
			return err
		}
		return deleteBusdevClaim(app)
	}, func(app core.App) error {
		if err := restoreBusdevClaim(app); err != nil {
			return err
		}
		return setClientsBusdevLeadRequired(app, true)
	})
}

func setClientsBusdevLeadRequired(app core.App, required bool) error {
	clients, err := app.FindCollectionByNameOrId("clients")
	if err != nil {
		return err
	}
	clients.Fields.GetByName("business_development_lead").(*core.RelationField).Required = required
	return app.Save(clients)
}

// deleteBusdevClaim removes the claim and every user_claims record that gives
// it. A database without the claim is left alone.
func deleteBusdevClaim(app core.App) error {
	claim, err := app.FindFirstRecordByFilter("claims", "name={:name}", dbx.Params{"name": busdevClaimName})
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}
	holders, err := app.FindAllRecords("user_claims", dbx.HashExp{"cid": claim.Id})
	if err != nil {
		return err
	}
	for _, holder := range holders {
		if err := app.Delete(holder); err != nil {
			return err
		}
	}
	return app.Delete(claim)
}

func restoreBusdevClaim(app core.App) error {
	_, err := app.FindFirstRecordByFilter("claims", "name={:name}", dbx.Params{"name": busdevClaimName})
	if err == nil {
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	collection, err := app.FindCollectionByNameOrId("claims")
	if err != nil {
		return err
	}
	record := core.NewRecord(collection)
	record.Set("id", busdevClaimID)
	record.Set("name", busdevClaimName)
	record.Set("description", busdevClaimDescription)
	return app.Save(record)
}
