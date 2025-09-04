package app

import (
	"autohost-cli/internal/helpers/app_helper"
	"autohost-cli/utils"
	"fmt"

	"github.com/spf13/cobra"
)

func appStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status [nombre]",
		Short: "Muestra el estado de una aplicación",
		Args:  cobra.ExactArgs(1),
		Run: utils.WithAppName(func(appName string) {
			status, err := app_helper.GetAppStatus(appName)
			if err != nil {
				fmt.Printf("❌ Error al obtener el estado de %s: %v\n", appName, err)
			} else {
				fmt.Printf("📊  Estado de %s: %s\n", appName, status)
			}
		}),
	}
}
