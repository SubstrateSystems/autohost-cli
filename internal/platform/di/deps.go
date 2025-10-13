package di

import (
	"autohost-cli/internal/adapters/storage/sqlite"
	"autohost-cli/internal/domain"
	"database/sql"
)

type Deps struct {
	DB       *sql.DB
	Repos    Repos
	Services Services
}

type Repos struct {
	Installed *sqlite.InstalledRepo
	Catalog   sqlite.CatalogRepo
}

type Services struct {
	App     domain.AppService
	Catalog domain.CatalogService
}
