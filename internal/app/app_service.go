package app

import (
	"context"

	"autohost-cli/internal/adapters/storage/sqlite"
	"autohost-cli/internal/domain"
)

type AppService struct {
	Installed *sqlite.InstalledRepo
}

func (s AppService) ListInstalled(ctx context.Context) ([]domain.InstalledApp, error) {
	return s.Installed.List(ctx)
}

func (s AppService) IsAppInstalled(ctx context.Context, name string) (bool, error) {
	return s.Installed.IsInstalledApp(ctx, name)
}
