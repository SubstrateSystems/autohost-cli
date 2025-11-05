package models

import (
	"database/sql"
)

type CatalogAppRow struct {
	ID            int64          `db:"id"`
	Name          string         `db:"name"`
	Description   string         `db:"description"`
	DefaultPort   string         `db:"default_port"`
	DefaultPortDB string         `db:"default_port_db"`
	ClientDB      sql.NullString `db:"client_db"`
	CreatedAt     sql.NullString `db:"created_at"`
	UpdatedAt     sql.NullString `db:"updated_at"`
}
