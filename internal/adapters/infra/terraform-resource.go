// // // file: internal/infra/terraform_splitdns.go
package infra

// import (
// 	"fmt"
// 	"os"
// )

// // import (
// // 	"archive/zip"
// // 	"bytes"
// // 	"compress/gzip"
// // 	"errors"
// // 	"fmt"
// // 	"io"
// // 	"net/http"
// // 	"os"
// // 	"os/exec"
// // 	"path/filepath"
// // 	"runtime"
// // 	"strings"
// // 	"time"
// // )

// // const (
// // 	tfVersion      = "1.9.8" // fija una versión estable
// // 	tfBinDirRel    = ".autohost/bin"
// // 	tfStateDirRel  = ".autohost/state/tailscale"
// // 	tfProviderVer  = "~> 0.21"
// // 	tfDownloadBase = "https://releases.hashicorp.com/terraform"
// // )

// type SplitDNSOpts struct {
// 	Tailnet      string   // ej: "tu-org.ts.net"  (si vacío, usa env TAILSCALE_TAILNET o "-")
// 	Domain       string   // ej: "maza-server"
// 	Nameservers  []string // ej: ["100.112.92.90"]
// 	SearchPaths  []string // opcional: ej ["maza-server"]
// 	APIKeyEnvVar string   // por defecto "TAILSCALE_API_KEY"
// }

// // // ConfigureSplitDNSWithTerraform genera el .tf y aplica con terraform.
// func ConfigureSplitDNSWithTerraform(opts SplitDNSOpts) error {
// 	if opts.Domain == "" || len(opts.Nameservers) == 0 {
// 		return fmt.Errorf("domain y al menos un nameserver son obligatorios")
// 	}

// 	// 1) Validar API key (Terraform provider la usa)
// 	apiEnv := opts.APIKeyEnvVar
// 	if apiEnv == "" {
// 		apiEnv = "TAILSCALE_API_KEY"
// 	}
// 	if os.Getenv(apiEnv) == "" {
// 		return fmt.Errorf("%s no está definido en el entorno", apiEnv)
// 	}

// 	// 2) Resolver tailnet
// 	tailnet := opts.Tailnet
// 	if tailnet == "" {
// 		tailnet = os.Getenv("TAILSCALE_TAILNET")
// 		if tailnet == "" {
// 			tailnet = "-" // tailnet por defecto del token
// 		}
// 	}

// 	// 3) Asegurar terraform binario
// 	tfPath, err := ensureTerraform()
// 	if err != nil {
// 		return fmt.Errorf("no se pudo asegurar terraform: %w", err)
// 	}

// 	// 4) Preparar workspace: ~/.autohost/state/tailscale/<tailnet>/split-dns-<domain>/
// 	ws, err := prepareWorkspace(tailnet, opts.Domain)
// 	if err != nil {
// 		return err
// 	}

// 	// 5) Escribir main.tf
// 	if err := writeMainTF(ws, opts.Domain, opts.Nameservers, opts.SearchPaths); err != nil {
// 		return err
// 	}

// 	// 6) terraform init
// 	if err := runCmd(ws, tfPath, "init", "-upgrade"); err != nil {
// 		return fmt.Errorf("terraform init falló: %w", err)
// 	}

// 	// 7) terraform apply
// 	if err := runCmd(ws, tfPath, "apply", "-auto-approve"); err != nil {
// 		return fmt.Errorf("terraform apply falló: %w", err)
// 	}

// 	return nil
// }

// // 		_ = os.Chmod(filepath.Join(binDir, "terraform"), 0o755)
// // 	}
// // 	return filepath.Join(binDir, "terraform"), nil
// // }


// // func writeMainTF(dir, domain string, nameservers, searchPaths []string) error {
// // 	if len(nameservers) == 0 {
// // 		return errors.New("nameservers requerido")
// // 	}
// // 	tf := fmt.Sprintf(`
// // terraform {
// //   required_providers {
// //     tailscale = {
// //       source  = "tailscale/tailscale"
// //       version = "%s"
// //     }
// //   }
// // }

// // # Usa TAILSCALE_API_KEY y (opcional) TAILSCALE_TAILNET desde el entorno
// // provider "tailscale" {}

// // resource "tailscale_dns_split_nameservers" "split" {
// //   domain      = %q
// //   nameservers = [%s]
// // }
// // `, tfProviderVer, domain, quoteJoin(nameservers))

// // 	if len(searchPaths) > 0 {
// // 		tf += fmt.Sprintf(`
// // resource "tailscale_dns_search_paths" "paths" {
// //   search_paths = [%s]
// // }
// // `, quoteJoin(searchPaths))
// // 	}

// // 	return os.WriteFile(filepath.Join(dir, "main.tf"), []byte(strings.TrimSpace(tf)+"\n"), 0o644)
// // }

// // func runCmd(workdir, bin string, args ...string) error {
// // 	cmd := exec.Command(bin, args...)
// // 	cmd.Dir = workdir
// // 	cmd.Stdout = os.Stdout
// // 	cmd.Stderr = os.Stderr
// // 	return cmd.Run()
// // }

// // func quoteJoin(items []string) string {
// // 	qs := make([]string, 0, len(items))
// // 	for _, s := range items {
// // 		qs = append(qs, fmt.Sprintf("%q", s))
// // 	}
// // 	return strings.Join(qs, ", ")
// // }

// // func httpGet(url string, timeout time.Duration) ([]byte, error) {
// // 	client := &http.Client{Timeout: timeout}
// // 	resp, err := client.Get(url)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	defer resp.Body.Close()
// // 	if resp.StatusCode >= 300 {
// // 		return nil, fmt.Errorf("descarga falló: %s", resp.Status)
// // 	}
// // 	return io.ReadAll(resp.Body)
// // }

// // func unzipTerraform(zipBytes []byte, dest string) error {
// // 	r, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
// // 	if err != nil {
// // 		return err
// // 	}
// // 	for _, f := range r.File {
// // 		if f.FileInfo().IsDir() {
// // 			continue
// // 		}
// // 		rc, err := f.Open()
// // 		if err != nil {
// // 			return err
// // 		}
// // 		defer rc.Close()
// // 		out := filepath.Join(dest, f.Name)
// // 		w, err := os.Create(out)
// // 		if err != nil {
// // 			return err
// // 		}
// // 		if _, err := io.Copy(w, rc); err != nil {
// // 			w.Close()
// // 			return err
// // 		}
// // 		w.Close()
// // 	}
// // 	return nil
// // }
