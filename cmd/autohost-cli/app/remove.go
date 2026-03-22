package app

import (
	"autohost-cli/internal/app"
	"autohost-cli/internal/domain"
	"autohost-cli/utils"
	"fmt"

	"github.com/spf13/cobra"
)

func appRemoveCmd(svc *app.AppService) *cobra.Command {
	return &cobra.Command{
		Use:   "remove [nombre]",
		Short: "Elimina una aplicación",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appName := args[0]
			if utils.Confirm(fmt.Sprintf("¿Estás seguro que quieres eliminar %s? [y/N]: ", appName)) {
				if err := svc.RemoveApp(cmd.Context(), domain.AppName(appName)); err != nil {
					fmt.Printf("❌ No se pudo eliminar %s: %v\n", appName, err)
				} else {
					fmt.Printf("🧹 %s eliminada correctamente.\n", appName)
				}
			}
			return nil
		},
	}
}
