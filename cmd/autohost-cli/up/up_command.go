package up

import (
	appSvc "autohost-cli/internal/app"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// UpCmd returns the `autohost up` command.
func UpCmd(svc *appSvc.UpService) *cobra.Command {
	var cloudURL string
	var name string

	cmd := &cobra.Command{
		Use:   "up",
		Short: "Conectar este nodo al cloud de AutoHost",
		Long: `Abre el navegador para autenticarse en el cloud de AutoHost,
genera un token de enrollment y registra este nodo automáticamente.

Similar a 'tailscale up': un solo comando para conectar el nodo.`,
		Example: `  # Usar la nube por defecto (https://cloud.autohost.dev)
  autohost up

  # Usar una instancia cloud personalizada
  autohost up --cloud https://mycloud.example.com

  # Personalizar el nombre del nodo
  autohost up --name mi-servidor`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := svc.Up(cloudURL, name); err != nil {
				return fmt.Errorf("autohost up: %w", err)
			}
			return nil
		},
	}

	// --cloud flag default: env var AUTOHOST_CLOUD_URL > compiled-in DefaultCloudURL
	cloudDefault := os.Getenv("AUTOHOST_CLOUD_URL")
	if cloudDefault == "" {
		cloudDefault = appSvc.DefaultCloudURL
	}
	cmd.Flags().StringVar(&cloudURL, "cloud", cloudDefault, "URL del cloud de AutoHost (env: AUTOHOST_CLOUD_URL)")
	cmd.Flags().StringVar(&name, "name", "", "Nombre del nodo (default: hostname del sistema)")

	return cmd
}
