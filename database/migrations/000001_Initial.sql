-- Represents a user and employee of company
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS link_types (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS magic_links (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    link_type TEXT NOT NULL,
    expired_at TEXT NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(link_type) REFERENCES link_types(id)
);

CREATE TABLE IF NOT EXISTS tasks (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    due_date TEXT NOT NULL,
    assigned_to TEXT,
    status TEXT NOT NULL,
    FOREIGN KEY(assigned_to) REFERENCES users(id)
);

-- Rerpesents a entity(client) in the system
CREATE TABLE IF NOT EXISTS entities (
    id TEXT PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL,
    phone TEXT NOT NULL,
    status TEXT NOT NULL,
    assigned_to TEXT,
    created_at TEXT NOT NULL,
    converted_at TEXT NOT NULL,
    FOREIGN KEY(assigned_to) REFERENCES users(id)
);
