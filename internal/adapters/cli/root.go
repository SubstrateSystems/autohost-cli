/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"autohost-cli/db"
	"autohost-cli/internal/adapters/cli/app"
	"autohost-cli/internal/adapters/cli/docker"
	"autohost-cli/internal/adapters/cli/setup"
	"autohost-cli/internal/adapters/storage/sqlite"
	appInternal "autohost-cli/internal/app"
	"autohost-cli/internal/platform/di"
	"autohost-cli/utils"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "autohost-cli",
		Short: "CLI para autohosting con Docker/Tailscale/Cloudflare/Caddy",
	}
	deps di.Deps
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	if !utils.IsInitialized() {
		err := ensureAutohostDirs()
		if err != nil {
			println("❌ Error al crear estructura de carpetas:", err.Error())
			os.Exit(1)
		}
		println("✅ Entorno de AutoHost creado")
	}

	sqlitePath := filepath.Join(utils.GetAutohostDir(), "autohost.db")

	dbc, err := db.Open(sqlitePath)
	if err != nil {
		fmt.Println("DB open error:", err)
		os.Exit(1)
	}
	if err := db.Migrate(dbc); err != nil {
		fmt.Println("DB migrate error:", err)
		os.Exit(1)
	}

	// Ejecutar seeding después de migraciones
	if err := db.Seed(dbc); err != nil {
		fmt.Println("DB seed error:", err)
		os.Exit(1)
	}

	// rootCmd.AddCommand(initializer.InitCommand())
	deps = buildDeps(dbc.DB)
	rootCmd.AddCommand(app.AppCmd(deps))
	rootCmd.AddCommand(setup.SetupCmd())
	rootCmd.AddCommand(docker.DockerCmd())

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func buildDeps(sqlDB *sql.DB) di.Deps {
	installedRepo := sqlite.NewInstalledRepo(sqlDB)
	catalogRepo := sqlite.NewCatalogRepo(sqlDB)

	return di.Deps{
		DB: sqlDB,
		Repos: di.Repos{
			Installed: installedRepo,
			Catalog:   catalogRepo,
		},
		Services: di.Services{
			App:     appInternal.AppService{Installed: installedRepo},
			Catalog: appInternal.CatalogService{Catalog: catalogRepo},
		},
	}
}

func ensureAutohostDirs() error {
	subdirs := []string{
		"config",
		"templates",
		"apps",
		"logs",
		"state",
		"backups",
	}

	for _, sub := range subdirs {
		if err := os.MkdirAll(utils.GetSubdir(sub), 0755); err != nil {
			return err
		}
	}
	return nil
}
