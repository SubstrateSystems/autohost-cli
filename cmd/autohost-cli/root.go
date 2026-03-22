package cli

import (
	"autohost-cli/cmd/autohost-cli/agent"
	"autohost-cli/cmd/autohost-cli/app"
	"autohost-cli/cmd/autohost-cli/cc"
	"autohost-cli/cmd/autohost-cli/enroll"
	"autohost-cli/cmd/autohost-cli/expose"
	"autohost-cli/cmd/autohost-cli/install"
	"autohost-cli/cmd/autohost-cli/setup"
	"autohost-cli/internal/adapters/caddy"
	"autohost-cli/internal/adapters/catalog"
	"autohost-cli/internal/adapters/cloudflare"
	coredns "autohost-cli/internal/adapters/coreDNS"
	"autohost-cli/internal/adapters/docker"
	"autohost-cli/internal/adapters/installed"
	"autohost-cli/internal/adapters/tailscale"
	"autohost-cli/internal/adapters/terraform"
	appSvc "autohost-cli/internal/app"
	"os"

	"github.com/spf13/cobra"
)

// Inyectado en build time por goreleaser
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:     "autohost",
	Short:   "CLI para autohosting con Docker/Tailscale/Cloudflare/Caddy",
	Version: Version,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Composition root: all services are constructed and injected here.

	dockerAdapter := docker.New()

	appService := &appSvc.AppService{
		Docker:    dockerAdapter,
		Catalog:   catalog.New(),
		Installed: installed.New(),
	}

	exposeService := &appSvc.ExposeService{
		Caddy:      caddy.New(),
		Tailscale:  tailscale.New(),
		CoreDNS:    coredns.New(),
		Cloudflare: cloudflare.New(),
		Terraform:  terraform.New(),
	}

	rootCmd.AddCommand(agent.AgentCmd(&appSvc.AgentService{}))
	rootCmd.AddCommand(enroll.EnrollCmd(&appSvc.EnrollService{}))
	rootCmd.AddCommand(cc.CCCmd(&appSvc.CCService{}))
	rootCmd.AddCommand(app.AppCmd(appService))
	rootCmd.AddCommand(install.InstallCmd(appService))
	rootCmd.AddCommand(expose.ExposeCmd(exposeService))
	rootCmd.AddCommand(setup.SetupCmd(&appSvc.SetupService{
		Docker:    dockerAdapter,
		Tailscale: tailscale.New(),
	}))
}
