package mappers

import (
	"autohost-cli/internal/adapters/storage/sqlite/models"
	"autohost-cli/internal/domain"
	"database/sql"
	"time"
)

func ToDomainCatalogApp(r models.CatalogAppRow) domain.CatalogApp {
	// Parse date strings into time.Time
	createdAt := time.Time{}
	updatedAt := time.Time{}

	if r.CreatedAt.Valid {
		t, err := time.Parse(time.RFC3339, r.CreatedAt.String)
		if err == nil {
			createdAt = t
		}
	}

	if r.UpdatedAt.Valid {
		t, err := time.Parse(time.RFC3339, r.UpdatedAt.String)
		if err == nil {
			updatedAt = t
		}
	}

	return domain.CatalogApp{
		Name:          r.Name,
		Description:   r.Description,
		DefaultPort:   r.DefaultPort,
		DefaultPortDB: r.DefaultPortDB,
		ClientDB:      nullString(r.ClientDB),
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func nullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
