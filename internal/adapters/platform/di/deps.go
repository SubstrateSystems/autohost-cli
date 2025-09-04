// internal/di/deps.go
package di

import (
	"database/sql"

	"autohost-cli/internal/repo"
	"autohost-cli/internal/services"
)

type Deps struct {
	DB       *sql.DB
	Repos    Repos
	Services Services
}

type Repos struct {
	Installed *repo.InstalledRepo
}

type Services struct {
	App services.AppService
}
