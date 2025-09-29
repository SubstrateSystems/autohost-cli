package domain

type CatalogApp struct {
	Name        string `db:"name"`
	Description string `db:"description"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
}

type InstalledApp struct {
	ID           int64  `db:"id" `
	Name         string `db:"name"`
	CatalogAppID string `db:"catalog_app_id"`
	CreatedAt    string `db:"created_at"`
}

type AppConfig struct {
	Name     string
	Template string
	Port     string
	MySQL    *MySQLConfig
	Postgres *PostgresConfig
}

type MySQLConfig struct {
	RootPassword string
	User         string
	Password     string
	Database     string
	Port         string
}

type PostgresConfig struct {
	User     string
	Password string
	Database string
	Port     string
}
