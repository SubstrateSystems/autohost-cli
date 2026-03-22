package agent

import (
	"autohost-cli/internal/app"

	"github.com/spf13/cobra"
)

// AgentCmd returns the root agent command.
func AgentCmd(svc *app.AgentService) *cobra.Command {
	agentCmd := &cobra.Command{
		Use:   "agent",
		Short: "Manage AutoHost agents",
	}
	agentCmd.AddCommand(agentInstallCmd(svc))
	return agentCmd
}
