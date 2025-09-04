/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"autohost-cli/cmd/app"
	"autohost-cli/cmd/docker"
	"autohost-cli/cmd/initializer"
	"autohost-cli/cmd/setup"
	"autohost-cli/db"
	"autohost-cli/internal/di"
	"autohost-cli/internal/repo"
	"autohost-cli/internal/services"
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

	deps = buildDeps(dbc.DB)
	rootCmd.AddCommand(app.AppCmd(deps))
	rootCmd.AddCommand(initializer.InitCommand())
	rootCmd.AddCommand(setup.SetupCmd())
	rootCmd.AddCommand(docker.DockerCmd())
	// rootCmd.AddCommand(caddy.CaddyCmd())

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func buildDeps(sqlDB *sql.DB) di.Deps {
	installedRepo := repo.NewInstalledRepo(sqlDB)

	return di.Deps{
		DB: sqlDB,
		Repos: di.Repos{
			Installed: installedRepo,
		},
		Services: di.Services{
			App: services.AppService{Installed: installedRepo},
		},
	}
}
