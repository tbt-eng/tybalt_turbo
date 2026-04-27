import type {
  CurrenciesResponse,
  TimeEntriesRecord,
  PurchaseOrdersRecord,
  PurchaseOrdersResponse,
  PoApproversResponse,
  ExpensesRecord,
  ExpensesResponse,
  JobsRecord,
  CategoriesResponse,
  ClientsRecord,
  ClientContactsResponse,
  VendorsRecord,
  VendorsResponse,
  TimeAmendmentsRecord,
  AdminProfilesRecord,
  AdminProfilesAugmentedResponse,
  ClientDetails,
  ClientNotesResponse,
  ExpensesAugmentedResponse,
} from "$lib/pocketbase-types";
import type { VisiblePurchaseOrderResponse } from "$lib/poVisibility";

export interface PageData<T> {
  item: T;
  editing: boolean;
  id: string | null;
}

export type LinkedWorkRecordSummary = {
  id: string;
  subject_id: string;
  work_record: string;
  date: string;
  job_id: string;
  job_number: string;
  type_id?: string;
  type_name?: string;
  location: string;
  work_description: string;
  hours_on_site: number;
  hours_travel_time: number;
  approved: boolean;
  linked_time_entry_id?: string;
  linked_timesheet_id?: string;
  linked_time_entry_date?: string;
  linked_time_entry_description?: string;
  linked_time_sheet_week_ending?: string;
  can_create_time_entry?: boolean;
};

export type TimeEntriesPageData = PageData<TimeEntriesRecord> & {
  linked_work_record?: LinkedWorkRecordSummary;
};
export type TimeAmendmentsPageData = PageData<TimeAmendmentsRecord>;
export type PurchaseOrdersPageData = PageData<PurchaseOrdersRecord | PurchaseOrdersResponse> & {
  parent_po_number?: string;
  parent_currency?: string;
  loadError?: string;
};
export type SecondApproverStatus =
  "not_required" | "requester_qualifies" | "candidates_available" | "required_no_candidates";
export type SecondApproversResponse = {
  approvers: PoApproversResponse[];
  meta: {
    second_approval_required: boolean;
    requester_qualifies: boolean;
    status: SecondApproverStatus;
    reason_code: string;
    reason_message: string;
    evaluated_amount: number;
    occurrences: number;
    second_approval_threshold: number;
    limit_column: string;
    second_stage_timeout_hours: number;
  };
};
export type PurchaseOrderDetailsPageData = {
  po: VisiblePurchaseOrderResponse;
  expenses: ExpensesAugmentedResponse[];
  secondApproverDiagnostics: SecondApproversResponse | null;
  canApproveOrReject: boolean;
};
export type LinkedPurchaseOrderSummary = {
  id: string;
  uid: string;
  po_number: string;
  currency: string;
  type: VisiblePurchaseOrderResponse["type"];
  payment_type: VisiblePurchaseOrderResponse["payment_type"];
  status: VisiblePurchaseOrderResponse["status"];
  recurring_expected_occurrences?: number;
  recurring_remaining_occurrences?: number;
  remaining_amount?: number;
};
export type ExpensesPageData = PageData<ExpensesRecord | ExpensesResponse> & {
  linked_purchase_order?: LinkedPurchaseOrderSummary;
  source_expense?: string;
};
export type JobsPageData = PageData<JobsRecord> & {
  categories: CategoriesResponse[];
  // When creating a project from a proposal, this contains any existing projects that already reference the proposal
  existingReferencingProjects?: { id: string; number: string }[];
};
export type ClientsPageData = PageData<ClientsRecord> & {
  client_contacts: ClientContactsResponse[];
};
export type ClientDetailsPageData = {
  client: ClientDetails;
  jobs: unknown[];
  notes: ClientNotesResponse[];
  noteJobs: JobsRecord[];
  tab: string;
  page: number;
  totalPages: number;
  projectsPage: number;
  proposalsPage: number;
  ownerPage: number;
  counts: {
    projects: number;
    proposals: number;
    owner: number;
  };
};
export type VendorsPageData = PageData<VendorsRecord | VendorsResponse>;
export type AdminProfilesPageData = PageData<AdminProfilesRecord | AdminProfilesAugmentedResponse>;
export type AdminProfilesEditPageData = PageData<AdminProfilesAugmentedResponse>;

// Expenses list pages use API endpoints returning augmented rows plus pagination metadata
export type ExpensesListData = {
  items: ExpensesAugmentedResponse[];
  createdItemIsVisible?: (record: ExpensesResponse) => boolean;
  totalPages?: number;
  limit?: number;
};

export type CurrencyListRow = CurrenciesResponse & {
  used_by_purchase_orders: boolean;
  used_by_expenses: boolean;
};

export type CurrencyInitStatus = {
  home_currency_id: string;
  home_currency_exists: boolean;
  home_currency_ready: boolean;
  blank_purchase_orders: number;
  blank_expenses: number;
};

export type PurchaseOrdersListData = {
  items?: VisiblePurchaseOrderResponse[];
  createdItemIsVisible?: (record: PurchaseOrdersResponse) => boolean;
  realtime_source?: "visible" | "pending" | "none";
};

export type ExpenseCommitQueueRow = {
  id: string;
  given_name: string;
  surname: string;
  submitted: boolean;
  approved: string;
  rejected: string;
  committed?: string;
  approver_name: string;
  committer_name: string;
  rejector_name: string;
  settler_name: string;
  date: string;
  description: string;
  attachment: string;
  allowance_str: string;
  job_number: string;
  job_description: string;
  client_name: string;
  total: number;
  settled_total: number;
  currency: string;
  currency_code: string;
  currency_symbol: string;
  currency_icon: string;
  currency_rate: number;
  currency_rate_date: string;
  po_number: string;
};

export type ExpenseSettlementRow = {
  id: string;
  uid: string;
  creator_name: string;
  date: string;
  description: string;
  payment_type: string;
  cc_last_4_digits: string;
  total: number;
  settled_total: number;
  settled: string;
  settler: string;
  settler_name: string;
  po_number: string;
  vendor_name: string;
  currency: string;
  currency_code: string;
  currency_symbol: string;
  currency_icon: string;
  currency_rate: number;
  currency_rate_date: string;
  indicative_cad_total: number;
  age_days: number;
};

export interface WorkRecordAttachmentResponse {
  id: string;
  attachment: string;
}

// Claims types returned by custom API endpoints
export interface ClaimListItem {
  id: string;
  name: string;
  description: string;
  holder_count: number;
}

export interface ClaimHolder {
  id: string;
  user_claim_id: string;
  uid: string;
  admin_profile_id: string;
  given_name: string;
  surname: string;
}

export interface ClaimAssignableUser {
  id: string;
  admin_profile_id: string;
  given_name: string;
  surname: string;
  name: string;
  username: string;
  email: string;
}

export interface ClaimDetails {
  id: string;
  name: string;
  description: string;
  holders: ClaimHolder[];
}

export interface ClaimAssignableUsers {
  id: string;
  name: string;
  description: string;
  assignable_users: ClaimAssignableUser[];
}

export interface WorkRecordListRow {
  id: string;
  work_record: string;
  date: string;
  creator: string;
  creator_name: string;
  job_number: string;
  job_id: string;
  type_id: string;
  type_name: string;
  location: string;
  work_description: string;
  attachment: string;
  approved: boolean;
  search_text: string;
  can_create_time_entry: boolean;
  subject_count: number;
  approved_subject_count: number;
  subjects: WorkRecordSubjectResponse[];
}

export interface WorkRecordDetailsResponse {
  id: string;
  work_record: string;
  date: string;
  creator: string;
  creator_name: string;
  job_id: string;
  job_number: string;
  client_name: string;
  job_description: string;
  client_reference_number: string;
  client_po: string;
  type_id: string;
  type_name: string;
  location: string;
  work_description: string;
  report_to: string;
  lane_kms: number;
  field_book_number: string;
  field_book_page_number: number;
  sub_contractor: string;
  equipment: string;
  supplies: string;
  attachment: string;
  approved: boolean;
  parent_work_record_id: string;
  parent_work_record: string;
  root_work_record_id: string;
  root_work_record: string;
  can_copy_to_tomorrow: boolean;
  can_edit: boolean;
  can_edit_shared: boolean;
  can_add_note: boolean;
  subjects: WorkRecordSubjectResponse[];
  notes: WorkRecordNoteResponse[];
  consumables: WorkRecordConsumableResponse[];
}

export interface WorkRecordOption {
  id: string;
  work_record: string;
}

export interface WorkRecordNoteResponse {
  id: string;
  uid: string;
  author_name: string;
  note: string;
  created: string;
}

export interface WorkRecordConsumableResponse {
  entry_id: string;
  consumable_id: string;
  consumable_name: string;
  input_kind: string;
  unit_label: string;
  sort_order: number;
  quantity_number: number;
  boolean_selected: boolean;
}

export interface WorkRecordSubjectResponse {
  id: string;
  work_record_id: string;
  uid: string;
  worker_name: string;
  hours_on_site: number;
  hours_travel_time: number;
  distance_travelled_km: number;
  vehicle_type: string;
  is_passenger: boolean;
  company_vehicle_unit_number: string;
  approved: boolean;
  linked_time_entry_id: string;
  linked_timesheet_id: string;
  linked_time_entry_date: string;
  linked_time_entry_description: string;
  linked_timesheet_week_ending: string;
  can_create_time_entry: boolean;
  can_edit: boolean;
}

export interface WorkRecordSubjectSaveRequest {
  id?: string;
  uid: string;
  hours_on_site: number;
  hours_travel_time: number;
  distance_travelled_km: number;
  vehicle_type: string;
  is_passenger: boolean;
  company_vehicle_unit_number: string;
}

export interface WorkRecordSaveConsumableRequest {
  entry_id: string;
  consumable_id: string;
  input_kind: string;
  quantity_number: number;
  boolean_selected: boolean;
}

export interface WorkRecordConsumablesReportResponse {
  items: WorkRecordConsumableReportRow[];
}

export interface WorkRecordCreateRequest {
  date: string;
  job: string;
  type: string;
  location: string;
  work_description: string;
  report_to: string;
  lane_kms?: number;
  field_book_number: string;
  field_book_page_number: number | "";
  sub_contractor: string;
  equipment: string;
  supplies: string;
  subjects: WorkRecordSubjectSaveRequest[];
  consumables: WorkRecordSaveConsumableRequest[];
  copy_from_work_record_id?: string;
}

export interface WorkRecordSaveRequest {
  date: string;
  job: string;
  type: string;
  location: string;
  work_description: string;
  report_to: string;
  lane_kms?: number;
  field_book_number: string;
  field_book_page_number: number | "";
  sub_contractor: string;
  equipment: string;
  supplies: string;
  subjects: WorkRecordSubjectSaveRequest[];
  consumables: WorkRecordSaveConsumableRequest[];
}

export interface WorkRecordSaveResponse {
  id: string;
}

export interface WorkRecordCopyDraftResponse {
  date: string;
  job: string;
  type: string;
  location: string;
  work_description: string;
  report_to: string;
  lane_kms: number;
  field_book_number: string;
  field_book_page_number: number;
  sub_contractor: string;
  equipment: string;
  supplies: string;
  subjects: WorkRecordSubjectResponse[];
  consumables: WorkRecordConsumableResponse[];
  copy_from_work_record_id: string;
  source_work_record: string;
  root_work_record_id: string;
  root_work_record: string;
}

export interface WorkRecordSearchResponse {
  items: WorkRecordListRow[];
  total: number;
}

export interface PaginatedJobWorkRecordsResponse {
  data: WorkRecordListRow[];
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface WorkRecordConsumableReportRow {
  id: string;
  work_record_id: string;
  work_record: string;
  date: string;
  creator: string;
  creator_name: string;
  job_id: string;
  job_number: string;
  type_id: string;
  type_name: string;
  consumable_id: string;
  consumable_name: string;
  input_kind: string;
  unit_label: string;
  quantity_number: number;
  boolean_selected: boolean;
}
