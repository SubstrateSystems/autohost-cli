-- Seeder inicial para catalog_apps
-- Inserta las aplicaciones disponibles en el catálogo

INSERT OR IGNORE INTO catalog_apps (name, description) VALUES
  ('bookstack', 'BookStack is a self-hosted, easy-to-use platform for organising and storing information'),
  ('nextcloud', 'Nextcloud is a suite of client-server software for creating and using file hosting services'),
  ('redis', 'Redis is an open source, in-memory data structure store, used as a database, cache, and message broker');
