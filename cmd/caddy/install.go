package caddy

import (
	"autohost-cli/utils"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

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

			homeDir, _ := os.UserHomeDir()
			caddyDir := filepath.Join(homeDir, ".autohost", "caddy")
			caddyfilePath := filepath.Join(caddyDir, "Caddyfile")

			err := os.MkdirAll(caddyDir, 0755)
			if err != nil {
				fmt.Println("❌ No se pudo crear el directorio de configuración de Caddy:", err)
				return
			}

			fmt.Println("📦 Instalando Caddy...")

			installScript := `
		sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https curl &&
		curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg &&
		curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list &&
		sudo apt update &&
		sudo apt install caddy
	`

			installCmd := exec.Command("bash", "-c", installScript)
			installCmd.Stdout = os.Stdout
			installCmd.Stderr = os.Stderr

			if err := installCmd.Run(); err != nil {
				fmt.Println("❌ Error al instalar Caddy:", err)
				return
			}

			// Crear Caddyfile si no existe
			if _, err := os.Stat(caddyfilePath); os.IsNotExist(err) {
				base := `# Archivo de configuración de Caddy para AutoHost
# Ejemplo:
# plex.localhost {
#     reverse_proxy 127.0.0.1:32400
# }
`
				os.WriteFile(caddyfilePath, []byte(base), 0644)
			}

			fmt.Println("✅ Caddy instalado y configurado. Puedes editar tu archivo en:")
			fmt.Println("   ", caddyfilePath)
		},
	}

}
