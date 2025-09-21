package app

import (
	"autohost-cli/internal/ports"
	"fmt"
)

type ExposeService struct {
	Caddy      ports.Caddy
	Tailscale  ports.Tailscale
	CoreDNS    ports.CoreDNS
	SplitDNS   ports.SplitDNS
	Cloudflare ports.Cloudflare
}

func (s *ExposeService) SetupPrivate(domain string) error {
	if err := s.Caddy.Install(); err != nil {
		return fmt.Errorf("caddy install: %w", err)
	}
	if err := s.Caddy.CreateCaddyfile(); err != nil {
		return fmt.Errorf("caddyfile: %w", err)
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
