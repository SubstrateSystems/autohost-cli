package app

import (
	"autohost-cli/internal/app"
	"fmt"

	"github.com/spf13/cobra"
)

func appStatusCmd(svc *app.AppService) *cobra.Command {
	return &cobra.Command{
		Use:   "status [nombre]",
		Short: "Muestra el estado de una aplicación",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			status, err := svc.GetAppStatus(args[0])
			if err != nil {
				fmt.Printf("❌ Error al obtener el estado de %s: %v\n", args[0], err)
			} else {
				fmt.Printf("📊  Estado de %s: %s\n", args[0], status)
			}
			return nil
		},
	}
}
