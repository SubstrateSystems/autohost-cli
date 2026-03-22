package app

import (
	"autohost-cli/internal/app"
	"fmt"

	"github.com/spf13/cobra"
)

func appStartCmd(svc *app.AppService) *cobra.Command {
	return &cobra.Command{
		Use:   "start [nombre]",
		Short: "Inicia una aplicación",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := svc.StartApp(args[0]); err != nil {
				fmt.Printf("❌ No se pudo iniciar %s: %v\n", args[0], err)
			} else {
				fmt.Printf("🚀 %s iniciada correctamente.\n", args[0])
			}
			return nil
		},
	}
}
