package enroll

import (
	"autohost-cli/internal/app"

	"github.com/spf13/cobra"
)

// EnrollCmd returns the root enroll command.
func EnrollCmd(svc *app.EnrollService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enroll",
		Short: "Enrollment commands for AutoHost agent",
	}
	cmd.AddCommand(newLinkCmd(svc))
	return cmd
}
