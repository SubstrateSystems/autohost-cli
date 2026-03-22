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
	ID            int64
	Name          string
	Description   string
	DefaultPort   string
	DefaultPortDB string
	ClientDB      string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type InstalledApp struct {
	ID           int64
	Name         string
	Port         string
	PortDB       string
	HttpURL      string
	Template     string
	CatalogAppID int64
	CreatedAt    time.Time
}

// Configuración: pertenece al dominio si expresa reglas del negocio
type AppConfig struct {
	AppSettings InstalledApp
	MySQL       *MySQLConfig
	Postgres    *PostgresConfig
	Minio       *MinioConfig
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

// MinioConfig holds credentials and storage configuration for a MinIO instance.
type MinioConfig struct {
	User        string
	Password    string
	ConsolePort string
	// DataPath is the host path where MinIO will persist data.
	// When backed by an external disk it is the disk's mount point.
	DataPath string
}
