namespace go deadliner.v1

struct Session {
    1: string account_uid
    2: string access_token
    3: string refresh_token
    4: string expires_at
}

struct RegisterRequest {
    1: string email
    2: string password
    3: string display_name
    4: string device_uid
    5: string device_name
    6: string platform
}

struct RegisterResponse {
    1: optional Session session
}

struct LoginRequest {
    1: string email
    2: string password
    3: string device_uid
    4: string device_name
    5: string platform
}

struct LoginResponse {
    1: optional Session session
}

struct RefreshSessionRequest {
    1: string refresh_token
    2: string device_uid
}

struct RefreshSessionResponse {
    1: optional Session session
}

struct ClientVersion {
    1: string ts
    2: i32 ctr
    3: string dev
}

struct ServerVersion {
    1: i64 change_id
    2: string committed_at
}

struct SubTask {
    1: string id
    2: string content
    3: bool is_completed
    4: i32 sort_order
    5: string created_at
    6: string updated_at
}

struct DeadlineDocument {
    1: string uid
    2: i64 legacy_id
    3: string name
    4: string start_time
    5: string end_time
    6: string state
    7: string complete_time
    8: string note
    9: bool is_stared
    10: string type
    11: i32 habit_count
    12: i32 habit_total_count
    13: i64 calendar_event
    14: string timestamp
    15: list<SubTask> sub_tasks
}

struct HabitConfig {
    1: string name
    2: string description
    3: i32 color
    4: string icon_key
    5: string period
    6: i32 times_per_period
    7: string goal_type
    8: i32 total_target
    9: string created_at
    10: string updated_at
    11: string status
    12: i32 sort_order
    13: string alarm_time
}

struct HabitRecord {
    1: string date
    2: i32 count
    3: string status
    4: string created_at
}

struct HabitDocument {
    1: string ddl_uid
    2: HabitConfig habit
    3: list<HabitRecord> records
}

struct DeadlineMutation {
    1: bool deleted
    2: optional DeadlineDocument doc
}

struct HabitMutation {
    1: bool deleted
    2: optional HabitDocument doc
}

union MutationPayload {
    1: DeadlineMutation deadline
    2: HabitMutation habit
}

struct Mutation {
    1: string mutation_id
    2: string device_uid
    3: string entity_uid
    4: ClientVersion client_version
    5: optional i64 base_change_id
    6: MutationPayload payload
}

struct MutationResult {
    1: string mutation_id
    2: string entity_uid
    3: bool accepted
    4: string rejection_reason
    5: optional ServerVersion server_version
    6: bool replayed
    7: optional string status
}

struct DeadlineChange {
    1: string entity_uid
    2: bool deleted
    3: ServerVersion server_version
    4: optional DeadlineDocument doc
}

struct HabitChange {
    1: string entity_uid
    2: bool deleted
    3: ServerVersion server_version
    4: optional HabitDocument doc
}

struct PullChangesRequest {
    1: string device_uid
    2: string cursor
    3: i32 limit
    4: bool include_deleted
}

struct PullChangesResponse {
    1: list<DeadlineChange> deadline_changes
    2: list<HabitChange> habit_changes
    3: string next_cursor
    4: bool has_more
}

struct PushChangesRequest {
    1: string device_uid
    2: string base_cursor
    3: list<Mutation> mutations
}

struct PushChangesResponse {
    1: list<MutationResult> results
    2: list<DeadlineChange> deadline_changes
    3: list<HabitChange> habit_changes
    4: string next_cursor
}

service DeadlinerService {
    RegisterResponse Register(1: RegisterRequest req)
    LoginResponse Login(1: LoginRequest req)
    RefreshSessionResponse RefreshSession(1: RefreshSessionRequest req)
    PullChangesResponse PullChanges(1: PullChangesRequest req)
    PushChangesResponse PushChanges(1: PushChangesRequest req)
}
