package docker

import (
	"autohost-cli/internal/adapters/docker"
	"autohost-cli/utils"
	"fmt"

	"github.com/spf13/cobra"
)

func dockerInstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Instala Docker y lo configura",
		Run: func(cmd *cobra.Command, args []string) {
			if !utils.IsInitialized() {
				fmt.Println("⚠️ Ejecuta `autohost init` primero.")
				return
			}

			if docker.DockerInstalled() {
				fmt.Println("✅ Docker ya está instalado.")
			} else {
				fmt.Println("🔧 Instalando Docker...")
				docker.InstallDocker()
			}

			if utils.Confirm("¿Agregar usuario al grupo docker? [y/N]: ") {
				docker.AddUserToDockerGroup()
			}
		},
	}
}
