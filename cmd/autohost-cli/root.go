/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"autohost-cli/cmd/autohost-cli/agent"
	"autohost-cli/internal/plugins/enroll"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "autohost-cli",
		Short: "CLI para autohosting con Docker/Tailscale/Cloudflare/Caddy",
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
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// func ensureAutohostDirs() error {
// 	subdirs := []string{
// 		"config",
// 		"templates",
// 		"apps",
// 		"logs",
// 		"state",
// 		"backups",
// 		"config",
// 	}

// 	for _, sub := range subdirs {
// 		if err := os.MkdirAll(utils.GetSubdir(sub), 0755); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
