PRAGMA foreign_keys=ON;

CREATE TABLE IF NOT EXISTS catalog_apps (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL UNIQUE,
  description TEXT NOT NULL,
  default_port TEXT NOT NULL,
  default_port_db TEXT NOT NULL,
  client_db TEXT,
  created_at DEFAULT (datetime('now')),
  updated_at DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS installed_apps (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL UNIQUE,
  port TEXT NOT NULL UNIQUE,
  port_db TEXT,
  http_url TEXT,
  catalog_app_id INTEGER NOT NULL REFERENCES catalog_apps(id) ON DELETE RESTRICT,
  created_at DEFAULT (datetime('now')),
  updated_at DEFAULT (datetime('now'))
);

