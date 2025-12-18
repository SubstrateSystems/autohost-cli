package enroll

import "github.com/spf13/cobra"

func EnrollCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enroll",
		Short: "Enrollment commands for AutoHost agent",
	}
	cmd.AddCommand(NewLinkCmd())
	return cmd
}
