package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

const (
	deployGroupsCollectionID           = "pbc_deploy_groups"
	deployGroupMembershipsCollectionID = "pbc_deploy_group_memberships"
)

// Add deploy groups, which give a feature to some users before the rest.
//
// deploy_groups names a group. deploy_group_memberships puts a user in a
// group. A feature's app_config row can name one group by its id, in the
// property deploy_group. Then the feature is on for a user when the feature's
// own switch is on and the user is in that group. With no deploy_group, the
// switch alone decides, as before. See descriptions/app_config.md.
//
// Only a PocketBase superuser writes these rows: neither collection has a
// create, update or delete rule. A user can list and view their own
// memberships, so that the UI can show or hide a feature. deploy_groups has no
// rule at all, because the UI compares ids and never needs a group's name.
func init() {
	m.Register(func(app core.App) error {
		users, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}

		groups := core.NewBaseCollection("deploy_groups")
		groups.Id = deployGroupsCollectionID
		groups.Fields.Add(
			&core.TextField{Name: "name", Required: true, Max: 100},
			&core.AutodateField{Name: "created", OnCreate: true},
			&core.AutodateField{Name: "updated", OnCreate: true, OnUpdate: true},
		)
		groups.AddIndex("idx_deploy_groups_name", true, "`name`", "")
		if err := app.Save(groups); err != nil {
			return err
		}

		memberships := core.NewBaseCollection("deploy_group_memberships")
		memberships.Id = deployGroupMembershipsCollectionID
		memberships.ListRule = types.Pointer("@request.auth.id != '' && user = @request.auth.id")
		memberships.ViewRule = types.Pointer("@request.auth.id != '' && user = @request.auth.id")
		memberships.Fields.Add(
			&core.RelationField{Name: "deploy_group", CollectionId: groups.Id, MaxSelect: 1, Required: true, CascadeDelete: true},
			&core.RelationField{Name: "user", CollectionId: users.Id, MaxSelect: 1, Required: true, CascadeDelete: true},
			&core.AutodateField{Name: "created", OnCreate: true},
			&core.AutodateField{Name: "updated", OnCreate: true, OnUpdate: true},
		)
		// One row for a user in a group. The user comes first, so the index also
		// serves the list of one user's memberships.
		memberships.AddIndex("idx_deploy_group_memberships_user_group", true, "`user`, `deploy_group`", "")
		return app.Save(memberships)
	}, func(app core.App) error {
		for _, id := range []string{deployGroupMembershipsCollectionID, deployGroupsCollectionID} {
			collection, err := app.FindCollectionByNameOrId(id)
			if err != nil {
				continue
			}
			if err := app.Delete(collection); err != nil {
				return err
			}
		}
		return nil
	})
}
