package agent

import (
	"autohost-cli/internal/app"

	"github.com/spf13/cobra"
)

func agentInstallCmd(svc *app.AgentService) *cobra.Command {
	return &cobra.Command{
		Use:     "install",
		Short:   "Install the AutoHost agent",
		Long:    "Download and install the AutoHost agent binary, configuration, and systemd service",
		Example: "  autohost agent install",
		RunE: func(cmd *cobra.Command, args []string) error {
			return svc.Install()
		},
	}
}
