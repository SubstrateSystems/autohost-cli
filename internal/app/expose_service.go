package app

import (
	"autohost-cli/internal/adapters/terraform"
	"autohost-cli/internal/ports"
	"context"
	"fmt"
)

type ExposeService struct {
	Caddy      ports.Caddy
	Tailscale  ports.Tailscale
	CoreDNS    ports.CoreDNS
	SplitDNS   ports.SplitDNS
	Cloudflare ports.Cloudflare
	Terraform  ports.Terraform
}

func (s *ExposeService) SetupPrivate(ctx context.Context, domain string) error {
	if err := s.Caddy.Install(); err != nil {
		return fmt.Errorf("caddy install: %w", err)
	}

	if err := s.Terraform.Install(ctx); err != nil {
		return fmt.Errorf("install terraform: %w", err)
	}

	if err := s.Caddy.CreateCaddyfile(); err != nil {
		return fmt.Errorf("caddyfile: %w", err)
	}
	// setup caddy snippets dir and sudoers
	if err := s.Caddy.EnsureCaddySnippetsSetup(ctx); err != nil {
		return fmt.Errorf("caddy setup snippets: %w", err)
	}

	if err := s.Tailscale.Install(); err != nil {
		return fmt.Errorf("tailscale install: %w", err)
	}

	if err := s.Tailscale.Login(); err != nil {
		return fmt.Errorf("tailscale login: %w", err)
	}

	ip, err := s.Tailscale.IP()
	if err != nil || ip == "" {
		return fmt.Errorf("tailscale ip: %w", err)
	}

	corefilePath, err := s.CoreDNS.InstallAndRun(ip)
	if err != nil {
		return fmt.Errorf("coredns: %w", err)
	}
	fmt.Println("🧩 CoreDNS listo. Corefile:", corefilePath)

	if s.SplitDNS != nil && domain != "" {
		if err := s.SplitDNS.Ensure(domain, []string{ip}); err != nil {
			return fmt.Errorf("split-dns: %w", err)
		}
	}
	return nil
}

func (s *ExposeService) SetupPublic(domain string) error {
	if domain == "" {
		return fmt.Errorf("domain requerido en modo public")
	}
	if err := s.Cloudflare.Install(); err != nil {
		return fmt.Errorf("cloudflare install: %w", err)
	}
	if err := s.Cloudflare.Login(); err != nil {
		return fmt.Errorf("cloudflare login: %w", err)
	}
	return nil
}

func (s *ExposeService) ExposeApp(ctx context.Context, subdomain string, nameApp string, port int) error {
	tailscaleIP, err := s.Tailscale.IP()
	if err != nil {
		return fmt.Errorf("no se pudo obtener la IP de Tailscale: %w", err)
	}
	// create splitDns in Tailscale
	cfg := ports.SplitDNSConfig{
		MagicDNS:    true,                // opcional pero útil
		SearchPaths: []string{subdomain}, // opcional; permite resolver "maza-server" como "maza-server.test" o "maza-server.test2"
		SplitNameservers: map[string][]string{
			subdomain: {tailscaleIP},
		},
	}
	if err := terraform.ApplySplitDNS(ctx, nameApp, cfg); err != nil {
		fmt.Printf("⚠️  No se pudo aplicar Split DNS en Tailscale: %v\n", err)
	}

	// update CoreFile and restart
	name, err := s.Tailscale.GetMachineName()
	if err != nil {
		return fmt.Errorf("no se pudo obtener el nombre de la máquina en Tailscale: %w", err)
	}
	nameWithSubdomain := fmt.Sprintf("%s.%s", subdomain, name)
	fmt.Printf("🔍 La IP de Tailscale es %q y el nombre de la máquina es %q (usando %q)\n", tailscaleIP, name, nameWithSubdomain)
	s.CoreDNS.UpdateCorefile(nameWithSubdomain, tailscaleIP)

	// update Caddyfile and restart
	machineName, _ := s.Tailscale.GetMachineName()

	if err := s.Caddy.AddService(machineName, 8080); err != nil {
		return fmt.Errorf("no se pudo actualizar Caddyfile: %w", err)
	}
	return nil
}
