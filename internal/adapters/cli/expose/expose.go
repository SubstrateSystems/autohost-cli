package expose

import (
	"autohost-cli/internal/adapters/infra"
	"autohost-cli/internal/adapters/tailscale"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func exposeAppCmd() *cobra.Command {
	var (
		exposeType string
		subdomain  string
		nameApp    string
	)

	cmd := &cobra.Command{
		Use:   "expose",
		Short: "Configura la exposición de servicios",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			exposeType = strings.ToLower(strings.TrimSpace(exposeType))
			switch exposeType {
			case "private", "public":
				// ok
			default:
				return fmt.Errorf("tipo inválido: %q (usa: private|public)", exposeType)
			}

			if subdomain == "" {
				return fmt.Errorf("subdominio no puede estar vacío")
			}
			if nameApp == "" {
				return fmt.Errorf("nombre de la app no puede estar vacío")
			}
			return nil

		},
		RunE: func(cmd *cobra.Command, args []string) error {
			switch exposeType {
			case "public":

				fmt.Println("🌐 Exposición vía Cloudflare seleccionada (no implementado aún).")
			case "private":
				fmt.Println("🛰️  Exposición vía Tailscale seleccionada (no implementado aún).")
				// create splitDns in Tailscale

				// update CoreFile and restart
				tailscaleIP, err := tailscale.TailscaleIP()
				if err != nil {
					return fmt.Errorf("no se pudo obtener la IP de Tailscale: %w", err)
				}
				name, err := tailscale.GetMachineName()
				if err != nil {
					return fmt.Errorf("no se pudo obtener el nombre de la máquina en Tailscale: %w", err)
				}
				nameWithSubdomain := fmt.Sprintf("%s.%s", subdomain, name)
				fmt.Printf("🔍 La IP de Tailscale es %q y el nombre de la máquina es %q (usando %q)\n", tailscaleIP, name, nameWithSubdomain)
				infra.UpdateCorefile(nameWithSubdomain, tailscaleIP)

				// update Caddyfile and restart

			}
			fmt.Printf("Exponiendo %q en %q a través de %q\n", subdomain, nameApp, exposeType)
			return nil
		},
	}

	return cmd
}
