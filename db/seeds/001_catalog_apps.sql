-- Seeder initial to catalog_apps
INSERT OR IGNORE INTO catalog_apps 
(name, description, default_port, default_port_db, client_db)
VALUES
  ('bookstack', 'BookStack is a self-hosted, easy-to-use platform for organising and storing information', '8080', '3306', 'mysql'),
  ('nextcloud', 'Nextcloud is a suite of client-server software for creating and using file hosting services', '8081', '3306', 'mysql'),
  ('redis', 'Redis is an open source, in-memory data structure store, used as a database, cache, and message broker', '6379', '0', NULL),
  ('mysql', 'MySQL is an open-source relational database management system based on SQL', '3306', '3306', NULL),
  ('postgres', 'PostgreSQL is a powerful, open source object-relational database system', '5432', '5432', NULL),
  ('joplin', 'Joplin Server is an open-source note-taking and synchronization service that allows users to connect their Joplin desktop and mobile applications.', '22300', '3306', 'postgres');
