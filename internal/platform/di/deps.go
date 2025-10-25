package di

import (
	"autohost-cli/internal/adapters/storage/sqlite/repository"
	"autohost-cli/internal/app"
	"autohost-cli/internal/domain"

	// "autohost-cli/internal/app"
	"database/sql"
)

type Deps struct {
	DB       *sql.DB
	Repos    Repos
	Services Services
}

type Repos struct {
	Installed domain.InstalledRepo
	Catalog   domain.CatalogRepo
}

type Services struct {
	App     app.AppService
	Catalog app.CatalogService
}

func Build(db *sql.DB) Deps {
	// Adapters
	installedRepo := repository.NewInstalledRepo(db)
	catalogRepo := repository.NewCatalogRepo(db)

	// Services (inyectando interfaces, aunque uses concretos aquí)
	appSvc := app.AppService{Installed: installedRepo}
	catSvc := app.CatalogService{Catalog: catalogRepo}

	return Deps{
		DB: db,
		Repos: Repos{
			Installed: installedRepo,
			Catalog:   catalogRepo,
		},
		Services: Services{
			App:     appSvc,
			Catalog: catSvc,
		},
	}
}
