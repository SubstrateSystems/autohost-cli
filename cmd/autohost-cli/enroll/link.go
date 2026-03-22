package enroll

import (
	"autohost-cli/internal/app"
	"errors"

	"github.com/spf13/cobra"
)

func newLinkCmd(svc *app.EnrollService) *cobra.Command {
	var api, token, name string

	cmd := &cobra.Command{
		Use:   "link",
		Short: "Enlaza este nodo con la API de AutoHost",
		Example: `  autohost enroll link --api https://api.autohost.dev --token <user-token>
  autohost enroll link --api https://api.autohost.dev --token <user-token> --name mi-servidor`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if api == "" || token == "" {
				return errors.New("--api y --token son obligatorios")
			}
			return svc.Link(api, token, name)
		},
	}

	cmd.Flags().StringVar(&api, "api", "", "URL base del API (https://...)")
	cmd.Flags().StringVar(&token, "token", "", "Token de sesión de usuario")
	cmd.Flags().StringVar(&name, "name", "", "Nombre lógico del nodo (por defecto: hostname)")

	_ = cmd.MarkFlagRequired("api")
	_ = cmd.MarkFlagRequired("token")
	return cmd
}
