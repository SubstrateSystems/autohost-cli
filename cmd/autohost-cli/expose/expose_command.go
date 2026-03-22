package expose

import (
	"autohost-cli/internal/app"

	"github.com/spf13/cobra"
)

// ExposeCmd returns the root expose command.
func ExposeCmd(svc *app.ExposeService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "expose",
		Short: "Comandos relacionados con la exposición de servicios",
	}
	cmd.AddCommand(exposeSetupCmd(svc))
	cmd.AddCommand(exposeAppCmd(svc))
	return cmd
}
