package caddy

import (
	"autohost-cli/internal/adapters/caddy"
	"autohost-cli/utils"
	"fmt"

	"github.com/spf13/cobra"
)

func caddyInstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Instala el servidor Caddy y prepara su configuración",
		Run: func(cmd *cobra.Command, args []string) {
			if !utils.IsInitialized() {
				fmt.Println("⚠️ AutoHost no está inicializado. Ejecuta `autohost init` primero.")
				return
			}

			caddy.InstallCaddy()
			caddy.CreateCaddyfile()

		},
	}
}
