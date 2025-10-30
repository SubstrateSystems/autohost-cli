PRAGMA foreign_keys=ON;

CREATE TABLE IF NOT EXISTS catalog_apps (
  name TEXT PRIMARY KEY,
  description TEXT NOT NULL,
  default_port TEXT NOT NULL UNIQUE,
  default_port_db TEXT NOT NULL UNIQUE,
  client_db TEXT,
  created_at INTEGER DEFAULT (unixepoch()),
  updated_at INTEGER DEFAULT (unixepoch())
);

CREATE TABLE IF NOT EXISTS installed_apps (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL UNIQUE,
  port TEXT NOT NULL UNIQUE,
  port_db TEXT,
  http_url TEXT,
  catalog_app_id TEXT NOT NULL REFERENCES catalog_apps(name) ON DELETE RESTRICT,
  created_at INTEGER DEFAULT (unixepoch()),
  updated_at INTEGER DEFAULT (unixepoch())
);

