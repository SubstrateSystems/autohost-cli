// internal/di/deps.go
package di

import (
	"autohost-cli/internal/adapters/storage/sqlite"
	"autohost-cli/internal/app"
	"database/sql"
)

type Deps struct {
	DB       *sql.DB
	Repos    Repos
	Services Services
}

type Repos struct {
	Installed *sqlite.InstalledRepo
}

type Services struct {
	App app.AppService
}
