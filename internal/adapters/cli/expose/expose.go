package expose

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func exposeCmd() *cobra.Command {
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
			}
			fmt.Printf("Exponiendo %q en %q a través de %q\n", subdomain, nameApp, exposeType)
			return nil
		},
	}

	return cmd
}
