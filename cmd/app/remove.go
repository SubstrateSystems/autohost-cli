package app

import (
	"autohost-cli/internal/helpers/app"
	"autohost-cli/utils"
	"fmt"

	"github.com/spf13/cobra"
)

func appRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove [nombre]",
		Short: "Elimina una aplicación",
		Args:  cobra.ExactArgs(1),
		Run: utils.WithAppName(func(appName string) {
			if utils.Confirm(fmt.Sprintf("¿Estás seguro que quieres eliminar %s? [y/N]: ", appName)) {
				err := app.RemoveApp(appName)
				if err != nil {
					fmt.Printf("❌ No se pudo eliminar %s: %v\n", appName, err)
				} else {
					fmt.Printf("🧹 %s eliminada correctamente.\n", appName)
				}
			}
		}),
	}

}
