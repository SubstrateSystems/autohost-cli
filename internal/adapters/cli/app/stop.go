package app

import (
	appKit "autohost-cli/internal/adapters/cli/app/appkit"
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
			err := appKit.StopApp(appName)

			if err != nil {
				fmt.Printf("❌ No se pudo detener %s: %v\n", appName, err)
			} else {
				fmt.Printf("🛑 %s detenida.\n", appName)
			}
		}),
	}

}
