package agent

import "github.com/spf13/cobra"

func AgentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Agente de AutoHost",
	}
	cmd.AddCommand(NewLinkCmd())
	return cmd
}
