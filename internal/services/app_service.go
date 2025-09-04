// internal/services/app_service.go
package services

import (
	"context"

	"autohost-cli/internal/models"
	"autohost-cli/internal/repo"
)

type AppService struct {
	Installed *repo.InstalledRepo
}

func (s AppService) ListInstalled(ctx context.Context) ([]models.InstalledApp, error) {
	return s.Installed.List(ctx)
}
