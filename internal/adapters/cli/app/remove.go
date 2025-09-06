package app

import (
	appKit "autohost-cli/internal/adapters/cli/app/appkit"
	"autohost-cli/internal/platform/di"
	"autohost-cli/utils"
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func appRemoveCmd(deps di.Deps) *cobra.Command {
	return &cobra.Command{
		Use:   "remove [nombre]",
		Short: "Elimina una aplicación",
		Args:  cobra.ExactArgs(1),
		Run: utils.WithAppName(func(ctx context.Context, appName string) {
			if utils.Confirm(fmt.Sprintf("¿Estás seguro que quieres eliminar %s? [y/N]: ", appName)) {
				err := appKit.RemoveApp(appName)
				deps.Repos.Installed.Remove(ctx, appName)

				if err != nil {
					fmt.Printf("❌ No se pudo eliminar %s: %v\n", appName, err)
				} else {
					fmt.Printf("🧹 %s eliminada correctamente.\n", appName)
				}
			}
		}),
	}

}
