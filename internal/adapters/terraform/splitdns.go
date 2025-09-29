package terraform

import (
	"autohost-cli/internal/ports"
	"autohost-cli/utils"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Aplica Split DNS usando un único main.tf autogenerado
func ApplySplitDNS(ctx context.Context, workProfile string, cfg ports.SplitDNSConfig) error {
	// 1) Asegurar Terraform instalado (ya tienes Install)
	// if err := Install(ctx); err != nil {
	// 	return fmt.Errorf("install terraform: %w", err)
	// }

	// 2) Workdir estable (por tailnet/perfil)
	baseDir := filepath.Join(os.Getenv("HOME"), ".autohost", "terraform", "splitdns", workProfile)
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return err
	}

	// 3) Escribir main.tf (si no existe o forzar overwrite)
	mainTF := filepath.Join(baseDir, "main.tf")
	if err := os.WriteFile(mainTF, []byte(mainTFContent), 0o644); err != nil {
		return err
	}

	// 4) Escribir variables (terraform.tfvars.json)
	tfvars := map[string]any{
		"magic_dns":         cfg.MagicDNS,
		"search_paths":      cfg.SearchPaths,
		"split_nameservers": cfg.SplitNameservers,
	}

	tfvarsBytes, _ := json.MarshalIndent(tfvars, "", "  ")
	if err := os.WriteFile(filepath.Join(baseDir, "terraform.tfvars.json"), tfvarsBytes, 0o644); err != nil {
		return err
	}
	bin := filepath.Join(utils.GetAutohostDir(), "terraform", "terraform")
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}

	tf := func(args ...string) *exec.Cmd {
		all := append([]string{"-chdir=" + baseDir}, args...)
		cmd := exec.CommandContext(ctx, bin, all...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		// Si necesitas variables (TAILSCALE_*, etc.)
		cmd.Env = os.Environ()
		return cmd
	}

	if err := tf("init", "-input=false").Run(); err != nil {
		return fmt.Errorf("terraform init: %w", err)
	}

	// Opcional: mostrar plan
	_ = tf("plan").Run()

	if err := tf("apply", "-auto-approve").Run(); err != nil {
		return fmt.Errorf("terraform apply: %w", err)
	}

	fmt.Println("✔ Split DNS aplicado con Terraform.")
	return nil
}

// Plantilla del main.tf (corrigida)
const mainTFContent = `terraform {
  required_providers {
    tailscale = {
      source  = "tailscale/tailscale"
      version = "~> 0.22"
    }
  }
}

provider "tailscale" {}

variable "magic_dns" {
  type = bool
}

variable "search_paths" {
  type    = list(string)
  default = []
}

variable "split_nameservers" {
  description = "Mapa dominio => lista de nameservers (IPs)"
  type        = map(list(string))
  default     = {}
}

# Activa/Desactiva MagicDNS
resource "tailscale_dns_preferences" "prefs" {
  magic_dns = var.magic_dns
}

# Agrega search paths si se pasan
resource "tailscale_dns_search_paths" "search" {
  count        = length(var.search_paths) > 0 ? 1 : 0
  search_paths = var.search_paths
}

# Aplica Split DNS por dominio
resource "tailscale_dns_split_nameservers" "split" {
  for_each   = var.split_nameservers
  domain     = each.key
  nameservers = each.value
}
`
