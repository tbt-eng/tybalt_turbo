// Deploy groups give a feature to some users before the rest.
//
// A feature's app_config row can name one group by its id, in the property
// deploy_group. The feature is then on for a user when the feature's own
// switch is on AND the user is in that group. With no deploy_group, the switch
// alone decides.
//
// The server does the same check and is the gate (utilities.InDeployGroup in
// app/utilities/config.go). This file only decides what the UI shows. It has
// no Vite or PocketBase imports, so a test can load it.

// Says whether the deploy group of a feature lets the user in. value is the
// value of the feature's app_config row, and groups holds the ids of the
// user's deploy groups.
export function deployGroupAllows(value: unknown, groups: readonly string[]): boolean {
  if (value === null || typeof value !== "object") return true;
  const group = (value as Record<string, unknown>).deploy_group;
  if (group === undefined || group === null || group === "") return true;
  // An id that matches no group has no members, so it lets nobody in.
  if (typeof group !== "string") return false;
  return groups.includes(group);
}
