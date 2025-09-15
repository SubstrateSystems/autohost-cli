package expose

import (
	"autohost-cli/internal/adapters/caddy"
	"autohost-cli/internal/adapters/cloudflare"
	"autohost-cli/internal/adapters/infra"
	tailscale "autohost-cli/internal/adapters/tilscale"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func exposeSetupCmd() *cobra.Command {
	var (
		mode   string
		domain string
		yes    bool
	)

	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Configura la exposición de servicios",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Normaliza y valida mode
			mode = strings.ToLower(strings.TrimSpace(mode))
			switch mode {
			case "private", "public":
				// ok
			default:
				return fmt.Errorf("modo inválido: %q (usa: private|public)", mode)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			switch mode {
			case "private":
				fmt.Println("🔒 Modo PRIVATE: sólo tailnet (DNS interno, sin exposición pública).")
				caddy.InstallCaddy()
				caddy.CreateCaddyfile()

				tailIP, err := tailscale.TailscaleIP()
				if err != nil || tailIP == "" {
					return fmt.Errorf("no se pudo obtener IP de tailscale (¿logueado?): %w", err)
				}
				fmt.Printf("🛰️  IP tailnet local: %s\n", tailIP)

				corefilePath, err := infra.InstallAndRunCoreDNSWithDocker(tailIP)
				if err != nil {
					return fmt.Errorf("no se pudo instalar CoreDNS con Docker: %w", err)
				}
				fmt.Println("🧩 CoreDNS (Docker) listo. Corefile:", corefilePath)

			case "public":
				fmt.Printf("🌍 Modo PUBLIC: exponiendo públicamente con dominio %s\n", domain)
				cloudflare.InstallCloudflare()
				fmt.Println("🌐 Asegúrate de tener un dominio registrado en Cloudflare.")
				cloudflare.LoginCloudflare()

			}

			return nil
		},
	}

	// Flags
	cmd.Flags().StringVarP(&mode, "mode", "m", "private", "Modo de exposición: private | public")
	cmd.Flags().StringVar(&domain, "domain", "", "Dominio base para exposición pública (requerido en --mode=public)")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "No preguntar (auto-confirmar)")

	// Autocompletado de valores para --mode
	_ = cmd.RegisterFlagCompletionFunc("mode", func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		return []string{"private", "public"}, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
