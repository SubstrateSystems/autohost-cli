package app

import (
	"autohost-cli/internal/adapters/docker"
	"autohost-cli/internal/app"
	"autohost-cli/internal/platform/di"
	"autohost-cli/utils"
	"fmt"

	"github.com/spf13/cobra"
)

func appRemoveCmd(deps di.Deps) *cobra.Command {

	var svc = &app.AppService{
		Docker: docker.New(),
	}

	cmd := &cobra.Command{
		Use:   "remove [nombre]",
		Short: "Elimina una aplicación",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appName := args[0]
			if utils.Confirm(fmt.Sprintf("¿Estás seguro que quieres eliminar %s? [y/N]: ", appName)) {
				err := svc.RemoveApp(appName)
				// deps.Repos.Installed.Remove(ctx, appName)

				if err != nil {
					fmt.Printf("❌ No se pudo eliminar %s: %v\n", appName, err)
				} else {
					fmt.Printf("🧹 %s eliminada correctamente.\n", appName)
				}
			}
			return nil
		},
	}

	return cmd
}
