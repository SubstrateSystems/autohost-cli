/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"autohost-cli/cmd/autohost-cli/agent"
	"autohost-cli/cmd/autohost-cli/cc"
	"autohost-cli/internal/plugins/enroll"
	"os"

	"github.com/spf13/cobra"
)

// Inyectado en build time por goreleaser
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

var (
	rootCmd = &cobra.Command{
		Use:     "autohost",
		Short:   "CLI para autohosting con Docker/Tailscale/Cloudflare/Caddy",
		Version: Version,
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// if !utils.IsInitialized() {
	// 	err := ensureAutohostDirs()
	// 	if err != nil {
	// 		println("❌ Error al crear estructura de carpetas:", err.Error())
	// 		os.Exit(1)
	// 	}
	// 	println("✅ Entorno de AutoHost creado")
	// }

	// rootCmd.AddCommand(initializer.InitCommand())
	// deps = di.Build(dbc.DB)
	// rootCmd.AddCommand(app.AppCmd())
	// rootCmd.AddCommand(install.InstallCmd())
	// rootCmd.AddCommand(setup.SetupCmd())
	// rootCmd.AddCommand(expose.ExposeCmd())
	rootCmd.AddCommand(agent.AgentCmd())
	rootCmd.AddCommand(enroll.EnrollCmd())
	rootCmd.AddCommand(cc.CCCmd())
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
