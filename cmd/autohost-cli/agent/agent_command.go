package agent

import "github.com/spf13/cobra"

func AgentCmd() *cobra.Command {
	agentCmd := &cobra.Command{
		Use:   "agent",
		Short: "Manage AutoHost agents",
	}

	agentCmd.AddCommand(agentInstallCmd())
	// agentCmd.AddCommand(agentStartCmd())
	// agentCmd.AddCommand(agentRestartCmd())
	// agentCmd.AddCommand(agentStopCmd())

	return agentCmd
}
