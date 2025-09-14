package setup

import (
	"autohost-cli/internal/adapters/caddy"
	cloudflarekit "autohost-cli/internal/adapters/cli/cloudflare/cloudflareKit"
	initializerkit "autohost-cli/internal/adapters/cli/initializer/initializerKit"
	"autohost-cli/internal/adapters/docker"
	tailscale "autohost-cli/internal/adapters/tilscale"
	"autohost-cli/utils"
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func SetupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "Configura tu servidor para autohospedar servicios",
		Long: `Este comando instala Docker, Caddy, configura dominios,
		y prepara túneles seguros para desplegar tus apps autohospedadas.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("\n🔧 Iniciando configuración del servidor...")

			initializerkit.EnsureAutohostDirs()

			if !docker.DockerInstalled() {
				if utils.Confirm("⚠️ Docker no está instalado. ¿Deseas instalarlo automáticamente? [y/N]: ") {
					docker.InstallDocker()
					docker.CreateDockerNetwork()
					fmt.Println("✅ Docker instalado correctamente.")
					fmt.Println("✅ Red Docker 'autohost_net' creada.")
					fmt.Println("🔄 Reiniciando sesión para aplicar cambios de grupo...")
				} else {
					fmt.Println("🚫 Instalación cancelada. Instala Docker manualmente y vuelve a ejecutar el setup.")
					return
				}
			} else {
				fmt.Println("✅ Docker ya está instalado.")
			}

			if utils.Confirm("¿Deseas agregar tu usuario al grupo 'docker' para usar Docker sin sudo? [y/N]: ") {
				docker.AddUserToDockerGroup()
			}

			if utils.Confirm("¿Deseas instalar y configurar Caddy como reverse proxy? [y/N]: ") {
				caddy.InstallCaddy()
				caddy.CreateCaddyfile()
			}

			option := utils.AskOption("🔒 ¿Qué tipo de acceso quieres configurar?", []string{"Tailscale (privado)", "Cloudflare Tunnel (público con dominio)"})
			switch option {
			case "Tailscale (privado)":
				tailscale.InstallTailscale()
			case "Cloudflare Tunnel (público con dominio)":
				cloudflarekit.InstallCloudflared()
				fmt.Print("Introduce el subdominio para el túnel (ej: blog.misitio.com): ")
				reader := bufio.NewReader(os.Stdin)
				domain, _ := reader.ReadString('\n')
				domain = strings.TrimSpace(domain)
				cloudflarekit.ConfigureCloudflareTunnel(domain)
			}

			fmt.Println("\n✅ Configuración inicial completa.")
		},
	}
}
