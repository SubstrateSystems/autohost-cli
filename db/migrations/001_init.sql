PRAGMA foreign_keys=ON;

CREATE TABLE IF NOT EXISTS catalog_apps (
  name TEXT PRIMARY KEY,
  description TEXT NOT NULL,
  created_at TEXT DEFAULT (datetime('now')),
  updated_at TEXT DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS installed_apps (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL UNIQUE,
  port TEXT,
  port_db TEXT,
  http_url TEXT,
  catalog_app_id TEXT NOT NULL REFERENCES catalog_apps(name) ON DELETE RESTRICT,
  created_at TEXT DEFAULT (datetime('now')),
  updated_at TEXT DEFAULT (datetime('now'))
);

