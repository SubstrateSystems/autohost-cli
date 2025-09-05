-- Seeder inicial para catalog_apps
-- Inserta las aplicaciones disponibles en el catálogo

INSERT OR IGNORE INTO catalog_apps (name, description) VALUES
  ('mysql', 'MySQL is an open-source relational database management system based on SQL'),
  ('postgres', 'PostgreSQL is a powerful, open source object-relational database system');