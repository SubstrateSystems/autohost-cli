package cc

import (
	"autohost-cli/internal/app"

	"github.com/spf13/cobra"
)

// CCCmd returns the root command for custom commands.
func CCCmd(svc *app.CCService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cc",
		Short: "Custom commands management",
		Long:  "Create and register custom bash scripts that can be executed remotely on this node.",
	}
	cmd.AddCommand(newCreateCmd(svc))
	cmd.AddCommand(newListCmd(svc))
	return cmd
}
