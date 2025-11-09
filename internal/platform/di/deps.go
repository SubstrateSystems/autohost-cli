package di

import (
	"autohost-cli/internal/adapters/storage/sqlite/repository"
	"autohost-cli/internal/app"
	"autohost-cli/internal/ports"

	"github.com/jmoiron/sqlx"
)

type Deps struct {
	DB       *sqlx.DB
	Repos    Repos
	Services Services
}

type Repos struct {
	Installed ports.InstalledRepository
	Catalog   ports.CatalogRepository
}

type Services struct {
	App     app.AppService
	Catalog app.CatalogService
}

func Build(db *sqlx.DB) Deps {

	installedRepo := repository.NewInstalledRepo(db)
	catalogRepo := repository.NewCatalogRepo(db)

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
