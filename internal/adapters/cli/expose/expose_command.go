package expose

import (
	"github.com/spf13/cobra"
)

func ExposeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "expose",
		Short: "Comandos relacionados con la exposición de servicios",
	}

	cmd.AddCommand(exposeSetupCmd())
	cmd.AddCommand(exposeAppCmd())

	return cmd
}
