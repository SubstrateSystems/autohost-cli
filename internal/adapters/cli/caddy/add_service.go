package caddy

import (
	"autohost-cli/internal/adapters/caddy"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	serviceName string
	servicePort int
	serviceHost string
)

func caddyAddServiceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add-service",
		Short: "Agrega un nuevo servicio al archivo Caddyfile",
		Run: func(cmd *cobra.Command, args []string) {
			caddy.AddServiceToCaddyfile(serviceName, serviceHost, servicePort)

			fmt.Printf("✅ Servicio '%s' agregado exitosamente a Caddyfile.\n", serviceName)
		},
	}

}
