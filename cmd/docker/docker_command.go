package docker

import "github.com/spf13/cobra"

func DockerCmd() *cobra.Command {
	dockerCmd := &cobra.Command{
		Use:   "docker",
		Short: "Comandos para configurar Docker",
	}

	dockerCmd.AddCommand(dockerInstallCmd())

	return dockerCmd
}
