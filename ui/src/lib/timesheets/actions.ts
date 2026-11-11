export interface TimesheetActionTarget {
  uid?: string | null;
  approver?: string | null;
  approved?: string | null;
  rejected?: string | null;
  committed?: string | null;
  submitted?: boolean | null;
}

function normalized(value?: string | null): string {
  return value ?? "";
}

function isSubmitted(target: TimesheetActionTarget): boolean {
  return target.submitted ?? true;
}

export function canRecallTimesheet(target: TimesheetActionTarget, viewerId: string): boolean {
  return (
    normalized(target.uid) === viewerId &&
    normalized(target.committed) === "" &&
    ((isSubmitted(target) &&
      normalized(target.approved) === "" &&
      normalized(target.rejected) === "") ||
      normalized(target.rejected) !== "")
  );
}

export function canApproveTimesheet(target: TimesheetActionTarget, viewerId: string): boolean {
  return (
    normalized(target.approver) === viewerId &&
    isSubmitted(target) &&
    normalized(target.approved) === "" &&
    normalized(target.rejected) === "" &&
    normalized(target.committed) === ""
  );
}

export function canRejectTimesheet(
  target: TimesheetActionTarget,
  viewerId: string,
  hasCommitAccess: boolean,
): boolean {
  if (
    !isSubmitted(target) ||
    normalized(target.committed) !== "" ||
    normalized(target.rejected) !== ""
  ) {
    return false;
  }

  if (normalized(target.approver) === viewerId) {
    return true;
  }

  return hasCommitAccess && normalized(target.approved) !== "";
}

// Sharing lets another manager read a time sheet. Work on a person's time sheet
// may have been done under another manager. The person's own manager approves
// the time sheet and shares it so that the other manager can see the hours
// worked for them. The other manager can read it but can't approve or reject it.
//
// Only the approver of the time sheet can share it. The create and delete rules
// of time_sheet_reviewers allow only the approver to add or remove a viewer.
// Those rules have no state condition, so a committed or rejected time sheet can
// still be shared. A viewer can open only a submitted time sheet, so the share
// control needs one.
export function canShareTimesheet(target: TimesheetActionTarget, viewerId: string): boolean {
  return viewerId !== "" && normalized(target.approver) === viewerId && isSubmitted(target);
}

export function canCommitTimesheet(
  target: TimesheetActionTarget,
  hasCommitAccess: boolean,
): boolean {
  return (
    hasCommitAccess &&
    isSubmitted(target) &&
    normalized(target.approved) !== "" &&
    normalized(target.rejected) === "" &&
    normalized(target.committed) === ""
  );
}
