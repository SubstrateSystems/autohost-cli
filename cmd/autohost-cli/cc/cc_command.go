package cc

import "github.com/spf13/cobra"

// CCCmd returns the root command for custom commands.
func CCCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cc",
		Short: "Custom commands management",
		Long:  "Create and register custom bash scripts that can be executed remotely on this node.",
	}

	cmd.AddCommand(newCreateCmd())
	cmd.AddCommand(newListCmd())

	return cmd
}
