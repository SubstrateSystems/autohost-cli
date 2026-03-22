package app

import (
	"autohost-cli/internal/app"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func appLsCmd(svc *app.AppService) *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "List installed applications",
		RunE: func(cmd *cobra.Command, args []string) error {
			apps, err := svc.ListInstalled(cmd.Context())
			if err != nil {
				return fmt.Errorf("could not get list of apps: %w", err)
			}
			if len(apps) == 0 {
				fmt.Println("No applications installed yet.")
				return nil
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME\tINSTALLED AT")
			for _, a := range apps {
				fmt.Fprintf(w, "%d\t%s\t%s\n", a.ID, a.Name, a.CreatedAt.Format("2006-01-02 15:04:05"))
			}
			_ = w.Flush()
			return nil
		},
	}
}
