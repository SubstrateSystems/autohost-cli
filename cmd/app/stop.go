package app

import (
	"autohost-cli/internal/helpers/app_helper"
	"autohost-cli/utils"
	"fmt"

	"github.com/spf13/cobra"
)

func appStopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop [nombre]",
		Short: "Detiene una aplicación",
		Args:  cobra.ExactArgs(1),
		Run: utils.WithAppName(func(appName string) {
			err := app_helper.StopApp(appName)

			if err != nil {
				fmt.Printf("❌ No se pudo detener %s: %v\n", appName, err)
			} else {
				fmt.Printf("🛑 %s detenida.\n", appName)
			}
		}),
	}

}
