package models

import "time"

type CatalogAppRow struct {
	Name          string    `db:"name"`
	Description   string    `db:"description"`
	DefaultPort   string    `db:"default_port"`
	DefaultPortDB string    `db:"default_port_db"`
	ClientDB      string    `db:"client_db"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
