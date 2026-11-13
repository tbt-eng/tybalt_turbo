package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// A rate sheet absorb may be done by an admin who does not hold the absorb
// claim (issue #191). The absorb_actions rules let only absorb holders list,
// view and delete (commit) an action, so such an admin could start a rate
// sheet absorb and then neither see nor commit it. Admit an admin to the
// rate sheet actions only. Actions of the other collections still need the
// absorb claim. The delete hook checks the rest of the rate sheet claims.
const (
	absorbActionsOldRule = "@request.auth.user_claims_via_uid.cid.name ?= 'absorb'"
	absorbActionsNewRule = "@request.auth.user_claims_via_uid.cid.name ?= 'absorb' || (collection_name = 'rate_sheets' && @request.auth.user_claims_via_uid.cid.name ?= 'admin')"
)

var absorbActionsRuleNames = []string{"listRule", "viewRule", "deleteRule"}

func init() {
	m.Register(func(app core.App) error {
		for _, rule := range absorbActionsRuleNames {
			if err := updateRule(app, "absorb_actions", rule, absorbActionsNewRule); err != nil {
				return err
			}
		}
		return nil
	}, func(app core.App) error {
		for _, rule := range absorbActionsRuleNames {
			if err := updateRule(app, "absorb_actions", rule, absorbActionsOldRule); err != nil {
				return err
			}
		}
		return nil
	})
}
