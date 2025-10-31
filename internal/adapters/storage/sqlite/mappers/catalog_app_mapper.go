package mappers

import (
	"autohost-cli/internal/adapters/storage/sqlite/models"
	"autohost-cli/internal/domain"
	"database/sql"
)

func ToDomainCatalogApp(r models.CatalogAppRow) domain.CatalogApp {
	return domain.CatalogApp{
		Name:          r.Name,
		Description:   r.Description,
		DefaultPort:   r.DefaultPort,
		DefaultPortDB: r.DefaultPortDB,
		ClientDB:      nullString(r.ClientDB),
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
	}
}

func nullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
