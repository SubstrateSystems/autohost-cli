package app

import (
	"autohost-cli/internal/app"
	"fmt"

	"github.com/spf13/cobra"
)

func appStopCmd(svc *app.AppService) *cobra.Command {
	return &cobra.Command{
		Use:   "stop [nombre]",
		Short: "Detiene una aplicación",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := svc.StopApp(args[0]); err != nil {
				fmt.Printf("❌ No se pudo detener %s: %v\n", args[0], err)
			} else {
				fmt.Printf("🛑 %s detenida.\n", args[0])
			}
			return nil
		},
	}
}
