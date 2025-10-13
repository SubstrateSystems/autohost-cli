package domain

// type AppService struct {
// 	Installed *sqlite.InstalledRepo
// }

// func (s AppService) ListInstalled(ctx context.Context) ([]domain.InstalledApp, error) {
// 	return s.Installed.List(ctx)
// }

// func (s AppService) RemoveApp(ctx context.Context, name string) error {
// 	return s.Installed.Remove(ctx, name)
// }

// func (s AppService) IsAppInstalled(ctx context.Context, name string) (bool, error) {
// 	return s.Installed.IsInstalledApp(ctx, name)
// }

// type CatalogService struct{ Catalog sqlite.CatalogRepo }

// func (s CatalogService) List(ctx context.Context) (any, error) {
// 	return s.Catalog.ListApps(ctx)
// }

// type CatalogApp struct {
// 	Name        string `db:"name"`
// 	Description string `db:"description"`
// 	CreatedAt   string `db:"created_at"`
// 	UpdatedAt   string `db:"updated_at"`
// }

// type InstalledApp struct {
// 	ID           int64  `db:"id" `
// 	Name         string `db:"name"`
// 	CatalogAppID string `db:"catalog_app_id"`
// 	CreatedAt    string `db:"created_at"`
// }

import "time"

type CatalogApp struct {
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type InstalledApp struct {
	ID           int64
	Name         string
	CatalogAppID string
	CreatedAt    time.Time
}

// Configuración: pertenece al dominio si expresa reglas del negocio
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
