package domain

import (
	"fmt"
	"regexp"
	"time"
)

type AppName string

func (n AppName) Validate() error {
	if n == "" {
		return fmt.Errorf("app name cannot be empty")
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(string(n)) {
		return fmt.Errorf("invalid app name format")
	}
	return nil
}

type CatalogApp struct {
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type InstalledApp struct {
	ID           int64
	Name         string
	CatalogAppID string
	CreatedAt    time.Time
}

// Configuración: pertenece al dominio si expresa reglas del negocio
type AppConfig struct {
	Name     string
	Template string
	Port     string
	MySQL    *MySQLConfig
	Postgres *PostgresConfig
}

type MySQLConfig struct {
	RootPassword string
	User         string
	Password     string
	Database     string
	Port         string
}

type PostgresConfig struct {
	User     string
	Password string
	Database string
	Port     string
}
