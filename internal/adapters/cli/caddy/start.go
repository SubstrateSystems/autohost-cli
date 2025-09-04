package caddy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func caddyStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Inicia Caddy con el archivo de configuración de AutoHost",
		Run: func(cmd *cobra.Command, args []string) {
			homeDir, _ := os.UserHomeDir()
			caddyfilePath := filepath.Join(homeDir, ".autohost", "caddy", "Caddyfile")

			fmt.Println("🚀 Iniciando servidor Caddy...")
			startCmd := exec.Command("caddy", "run", "--config", caddyfilePath)
			startCmd.Stdout = os.Stdout
			startCmd.Stderr = os.Stderr
			err := startCmd.Run()
			if err != nil {
				fmt.Println("❌ Error al iniciar Caddy:", err)
			} else {
				fmt.Println("✅ Caddy iniciado correctamente.")
			}
		},
	}
}
