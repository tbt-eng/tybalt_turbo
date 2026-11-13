// Package absorb contains the configuration for record absorption/merging.
// This is a separate package to avoid import cycles between routes and hooks.
package absorb

import (
	"fmt"
	"sort"

	"github.com/pocketbase/pocketbase/core"

	"tybalt/utilities"
)

// Config defines the complete configuration for absorbing a collection.
// All absorb-related relationships are defined here in one place.
type Config struct {
	// RefConfigs defines which tables/columns get updated when absorbing this collection.
	// When records are absorbed, references in these tables are updated to point to the target.
	RefConfigs []RefConfig

	// ParentConstraint enforces that records being absorbed must share the same parent.
	// For example, client_contacts being absorbed must all belong to the same client.
	ParentConstraint *ParentConstraint

	// DependsOn specifies the parent collection in the absorb dependency graph.
	// If set:
	//   - This collection's COMMIT is blocked when the parent has a pending absorb
	//   - The parent's UNDO is blocked when this collection has a pending absorb
	// This prevents data inconsistencies where undoing a parent absorb would try to
	// restore references to records that were deleted by a committed child absorb.
	DependsOn string

	// The fields below are optional. A collection that leaves them unset
	// behaves as before: the absorb claim allows absorb, undo and commit, any
	// record can absorb any other, and only the absorbed records are deleted.

	// ClaimSets replaces the absorb claim check for absorb, undo and commit.
	// A user passes when they hold every claim of at least one set. The
	// absorb_actions collection rules still apply to commit and to the list.
	ClaimSets [][]string

	// OwnedRows are child rows that belong to an absorbed record, such as the
	// entries of a rate sheet. Column is the child's relation to the absorbed
	// record. Absorb uses plain SQL, so a cascadeDelete does not run: the rows
	// are saved in the action with their record, deleted with it and restored
	// on undo. They are not moved to the target.
	OwnedRows []RefConfig

	// Prepare checks a request inside the absorb transaction, before anything
	// changes, and returns the ids to absorb. It may add ids, for example the
	// other revisions of a rate sheet. A *RefusedError tells the user why the
	// absorb is refused.
	Prepare func(app core.App, target *core.Record, ids []string) ([]string, error)

	// OnCommit runs in the commit transaction, before the action is deleted.
	// auth is the user who commits.
	OnCommit func(app core.App, action *core.Record, auth *core.Record) error

	// ChangesJobs marks a collection whose absorb moves references in jobs, so
	// the absorb is refused while job editing is turned off.
	ChangesJobs bool
}

// RefusedError is a rule that refuses an absorb. Message is for the user.
type RefusedError struct {
	Message string
}

func (e *RefusedError) Error() string {
	return e.Message
}

// Refused returns a *RefusedError with a formatted message.
func Refused(format string, args ...any) error {
	return &RefusedError{Message: fmt.Sprintf(format, args...)}
}

// OwnedRowsKey is the key under which an absorbed record keeps its owned
// rows in absorbed_records, as map[table][]row. A field name can't contain
// "@", so the key never collides with a field.
const OwnedRowsKey = "@owned_rows"

type RefConfig struct {
	Table  string
	Column string
}

// ParentConstraint enforces data integrity during record absorption by ensuring
// that records being absorbed together must share the same parent relationship.
// This prevents merging records that don't logically belong together.
type ParentConstraint struct {
	Collection string // The parent collection that must match (e.g., "clients")
	Field      string // The field that must have matching values (e.g., "client")
}

// Configs is the single source of truth for all absorb configurations.
// Add new absorbable collections here.
var Configs = map[string]Config{
	"clients": {
		RefConfigs: []RefConfig{
			{"client_contacts", "client"},
			{"client_notes", "client"},
			{"jobs", "client"},
			{"jobs", "job_owner"},
			{"job_requests", "client"},
			{"job_requests", "job_owner"},
			{"client_invoicing_information", "client"},
		},
		ChangesJobs: true,
		// clients has no parent - it's at the top of the hierarchy
	},
	"client_contacts": {
		RefConfigs: []RefConfig{
			{"jobs", "contact"},
			{"job_requests", "contact"},
			{"client_invoicing_information", "contact"},
		},
		ParentConstraint: &ParentConstraint{
			Collection: "clients",
			Field:      "client",
		},
		ChangesJobs: true,
		// client_contacts depends on clients for absorb ordering:
		// - Can't commit client_contacts while clients absorb exists
		// - Can't undo clients while client_contacts absorb exists
		DependsOn: "clients",
	},
	"vendors": {
		RefConfigs: []RefConfig{
			{"purchase_orders", "vendor"},
			{"expenses", "vendor"},
		},
		// vendors is independent - no parent dependency
	},
	// rate_sheets is added by Register in the hooks package
	// (hooks/rate_sheet_absorb.go), because its Prepare and OnCommit use the
	// rate sheet comparison there, and hooks imports this package.
}

// Register adds a collection to Configs. Call it from an init function only.
func Register(collectionName string, config Config) {
	if _, exists := Configs[collectionName]; exists {
		panic(fmt.Sprintf("absorb config for %s is already registered", collectionName))
	}
	Configs[collectionName] = config
}

// RequiredClaimSets returns the claim sets that allow a user to absorb, undo
// and commit in a collection. A collection without ClaimSets, or one that is
// not configured, needs the absorb claim.
func RequiredClaimSets(collectionName string) [][]string {
	if config, ok := Configs[collectionName]; ok && len(config.ClaimSets) > 0 {
		return config.ClaimSets
	}
	return [][]string{{"absorb"}}
}

// HasRequiredClaims reports whether auth may absorb, undo and commit in the
// collection, by RequiredClaimSets.
func HasRequiredClaims(app core.App, auth *core.Record, collectionName string) (bool, error) {
	for _, set := range RequiredClaimSets(collectionName) {
		hasAll := true
		for _, claim := range set {
			has, err := utilities.HasClaim(app, auth, claim)
			if err != nil {
				return false, err
			}
			if !has {
				hasAll = false
				break
			}
		}
		if hasAll {
			return true, nil
		}
	}
	return false, nil
}

// CollectionsReferencing returns, sorted, the collections whose absorb moves
// references in table.
func CollectionsReferencing(table string) []string {
	names := []string{}
	for name, config := range Configs {
		for _, ref := range config.RefConfigs {
			if ref.Table == table {
				names = append(names, name)
				break
			}
		}
	}
	sort.Strings(names)
	return names
}

// GetConfig returns the absorb configuration for a collection.
func GetConfig(collectionName string) (Config, error) {
	config, ok := Configs[collectionName]
	if !ok {
		return Config{}, fmt.Errorf("unknown collection: %s", collectionName)
	}
	return config, nil
}

// GetRefConfigs returns the reference configs and parent constraint for a collection.
func GetRefConfigs(collectionName string) ([]RefConfig, *ParentConstraint, error) {
	config, err := GetConfig(collectionName)
	if err != nil {
		return nil, nil, err
	}
	return config.RefConfigs, config.ParentConstraint, nil
}

// GetUndoBlockers returns collections that block undoing this collection's absorb.
// These are collections that have DependsOn pointing to this collection.
func GetUndoBlockers(collectionName string) []string {
	var blocks []string
	for name, config := range Configs {
		if config.DependsOn == collectionName {
			blocks = append(blocks, name)
		}
	}
	return blocks
}

// GetCommitBlockers returns collections that block committing this collection's absorb.
// This is simply the DependsOn field - if this collection depends on another, that other
// collection's absorb must be resolved before this one can be committed.
func GetCommitBlockers(collectionName string) []string {
	config, ok := Configs[collectionName]
	if !ok || config.DependsOn == "" {
		return nil
	}
	return []string{config.DependsOn}
}
