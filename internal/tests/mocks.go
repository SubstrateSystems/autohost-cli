package tests

import (
	"context"

	"autohost-cli/internal/domain"
)

// MockInstalledRepo provides a mock implementation for testing
type MockInstalledRepo struct {
	Apps      []domain.InstalledApp
	AddError  error
	ListError error
	RemoveError error
	IsInstalledError error
}

func NewMockInstalledRepo() *MockInstalledRepo {
	return &MockInstalledRepo{
		Apps: []domain.InstalledApp{},
	}
}

func (m *MockInstalledRepo) Add(ctx context.Context, app domain.InstalledApp) error {
	if m.AddError != nil {
		return m.AddError
	}
	
	// Simulate adding to the list
	app.ID = int64(len(m.Apps) + 1)
	if app.CreatedAt == "" {
		app.CreatedAt = "2023-01-01 10:00:00"
	}
	m.Apps = append(m.Apps, app)
	return nil
}

func (m *MockInstalledRepo) List(ctx context.Context) ([]domain.InstalledApp, error) {
	if m.ListError != nil {
		return nil, m.ListError
	}
	return m.Apps, nil
}

func (m *MockInstalledRepo) Remove(ctx context.Context, name string) error {
	if m.RemoveError != nil {
		return m.RemoveError
	}
	
	// Simulate removing from the list
	for i, app := range m.Apps {
		if app.Name == name {
			m.Apps = append(m.Apps[:i], m.Apps[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *MockInstalledRepo) IsInstalledApp(ctx context.Context, name string) (bool, error) {
	if m.IsInstalledError != nil {
		return false, m.IsInstalledError
	}
	
	for _, app := range m.Apps {
		if app.Name == name {
			return true, nil
		}
	}
	return false, nil
}

// MockCatalogRepo provides a mock implementation for testing
type MockCatalogRepo struct {
	Apps      []domain.CatalogApp
	ListError error
}

func NewMockCatalogRepo() MockCatalogRepo {
	return MockCatalogRepo{
		Apps: []domain.CatalogApp{
			{
				Name:        "nextcloud",
				Description: "Suite de software para servicios de hosting de archivos",
				CreatedAt:   "2023-01-01 10:00:00",
				UpdatedAt:   "2023-01-01 10:00:00",
			},
			{
				Name:        "bookstack",
				Description: "Plataforma para organizar y almacenar información",
				CreatedAt:   "2023-01-01 10:00:00", 
				UpdatedAt:   "2023-01-01 10:00:00",
			},
		},
	}
}

func (m MockCatalogRepo) ListApps(ctx context.Context) ([]domain.CatalogApp, error) {
	if m.ListError != nil {
		return nil, m.ListError
	}
	return m.Apps, nil
}