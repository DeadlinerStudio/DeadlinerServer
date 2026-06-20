CREATE TABLE accounts (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    account_uid VARCHAR(64) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    UNIQUE KEY uk_accounts_account_uid (account_uid),
    UNIQUE KEY uk_accounts_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE devices (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    device_uid VARCHAR(64) NOT NULL,
    account_id BIGINT NOT NULL,
    platform VARCHAR(64) NOT NULL,
    device_name VARCHAR(255) NOT NULL,
    last_seen_at DATETIME(3) NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    UNIQUE KEY uk_devices_device_uid (device_uid),
    KEY idx_devices_account_id (account_id),
    CONSTRAINT fk_devices_account_id FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE sessions (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    session_uid VARCHAR(64) NOT NULL,
    account_id BIGINT NOT NULL,
    refresh_token_hash VARCHAR(255) NOT NULL,
    expires_at DATETIME(3) NOT NULL,
    revoked_at DATETIME(3) NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    UNIQUE KEY uk_sessions_session_uid (session_uid),
    KEY idx_sessions_account_id (account_id),
    CONSTRAINT fk_sessions_account_id FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE deadline_items (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    account_id BIGINT NOT NULL,
    uid VARCHAR(128) NOT NULL,
    legacy_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    start_time VARCHAR(64) NOT NULL,
    end_time VARCHAR(64) NOT NULL,
    state VARCHAR(64) NOT NULL,
    complete_time VARCHAR(64) NOT NULL DEFAULT '',
    note TEXT NOT NULL,
    is_stared TINYINT(1) NOT NULL DEFAULT 0,
    type VARCHAR(32) NOT NULL,
    habit_count INT NOT NULL DEFAULT 0,
    habit_total_count INT NOT NULL DEFAULT 0,
    calendar_event BIGINT NOT NULL DEFAULT -1,
    business_timestamp VARCHAR(64) NOT NULL,
    sub_tasks JSON NOT NULL,
    deleted TINYINT(1) NOT NULL DEFAULT 0,
    client_ver_ts VARCHAR(64) NULL,
    client_ver_ctr INT NULL,
    client_ver_dev VARCHAR(64) NULL,
    server_change_id BIGINT NOT NULL,
    committed_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_by_device_uid VARCHAR(64) NULL,
    UNIQUE KEY uk_deadline_items_account_uid (account_id, uid),
    KEY idx_deadline_items_account_change (account_id, server_change_id),
    CONSTRAINT fk_deadline_items_account_id FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE habit_docs (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    account_id BIGINT NOT NULL,
    ddl_uid VARCHAR(128) NOT NULL,
    payload JSON NULL,
    deleted TINYINT(1) NOT NULL DEFAULT 0,
    client_ver_ts VARCHAR(64) NULL,
    client_ver_ctr INT NULL,
    client_ver_dev VARCHAR(64) NULL,
    server_change_id BIGINT NOT NULL,
    committed_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_by_device_uid VARCHAR(64) NULL,
    UNIQUE KEY uk_habit_docs_account_ddl_uid (account_id, ddl_uid),
    KEY idx_habit_docs_account_change (account_id, server_change_id),
    CONSTRAINT fk_habit_docs_account_id FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE sync_changes (
    change_id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    account_id BIGINT NOT NULL,
    device_uid VARCHAR(64) NOT NULL,
    mutation_id VARCHAR(128) NOT NULL,
    entity_kind VARCHAR(32) NOT NULL,
    entity_uid VARCHAR(128) NOT NULL,
    action VARCHAR(32) NOT NULL,
    payload JSON NULL,
    committed_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    UNIQUE KEY uk_sync_changes_account_device_mutation (account_id, device_uid, mutation_id),
    KEY idx_sync_changes_account_change (account_id, change_id),
    CONSTRAINT fk_sync_changes_account_id FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE mutation_receipts (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    account_id BIGINT NOT NULL,
    device_uid VARCHAR(64) NOT NULL,
    mutation_id VARCHAR(128) NOT NULL,
    entity_kind VARCHAR(32) NOT NULL,
    entity_uid VARCHAR(128) NOT NULL,
    status VARCHAR(32) NOT NULL,
    replayed TINYINT(1) NOT NULL DEFAULT 0,
    result_change_id BIGINT NULL,
    result_payload JSON NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    UNIQUE KEY uk_mutation_receipts_account_device_mutation (account_id, device_uid, mutation_id),
    KEY idx_mutation_receipts_account_device (account_id, device_uid, mutation_id),
    CONSTRAINT fk_mutation_receipts_account_id FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

