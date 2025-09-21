package terraform

import (
	"autohost-cli/internal/adapters/infra"
	"autohost-cli/internal/platform/config"
	"autohost-cli/utils"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func Install(ctx context.Context) error {
	p, err := infra.DetectPlatform()
	if err != nil {
		return err
	}

	if p.GOOS != "linux" {
		return fmt.Errorf("por ahora solo Linux (Ubuntu / Raspberry Pi OS Lite)")
	}

	if infra.HasCmd("terraform") {
		ver, _ := version(ctx)
		if ver != "" {
			fmt.Printf("✔ Terraform ya instalado (%s). Nada que hacer.\n", ver)
			return nil
		}
	}

	key, err := p.ArchKey()
	if err != nil {
		return err
	}

	url := ""
	switch key {
	case "linux_amd64":
		url = config.MustString("urls.toml", "terraform", "linux_amd64")
	case "linux_arm64":
		url = config.MustString("urls.toml", "terraform", "linux_arm64")
	case "linux_arm":
		url = config.MustString("urls.toml", "terraform", "linux_arm")
	}
	if url == "" {
		return fmt.Errorf("no hay URL configurada para %s", key)
	}

	versionTerraform := config.MustString("urls.toml", "terraform", "version")

	fmt.Printf("↓ Descargando Terraform (%s) para %s…\n", versionTerraform, key)
	zipPath, err := infra.DownloadToTemp(ctx, url)
	if err != nil {
		return err
	}

	tmpDir := os.TempDir()
	extracted, err := infra.UnzipSingleBinary(zipPath, "terraform", tmpDir)
	if err != nil {
		return err
	}

	finalPath := filepath.Join(utils.GetAutohostDir(), "terraform")
	if err := os.MkdirAll(finalPath, 0755); err != nil {
		return err
	}

	if err := moveFile(extracted, finalPath); err != nil {
		return err
	}
	if err := os.Chmod(finalPath, 0o755); err != nil {
		return err
	}

	fmt.Printf("✔ Terraform instalado en %s\n", finalPath)
	ver, _ := version(ctx)
	if ver != "" {
		fmt.Printf("→ terraform -version => %s\n", ver)
	}
	// ensureOnPathMessage(targetDir)
	return nil
}

func version(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "terraform", "-version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	// salida típica: "Terraform v1.9.5\n…"
	s := string(out)
	for _, ln := range []string{s} {
		if len(ln) > 0 {
			// devolvemos la primera línea entera (simple)
			return stringsTrimFirstLine(s), nil
		}
	}
	return "", nil
}

func stringsTrimFirstLine(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' || s[i] == '\r' {
			return s[:i]
		}
	}
	return s
}

func canWrite(dir string) bool {
	test := filepath.Join(dir, ".autohost_write_test")
	if err := os.WriteFile(test, []byte("ok"), 0o644); err != nil {
		return false
	}
	_ = os.Remove(test)
	return true
}

func moveFile(src, dst string) error {
	// intenta rename; si falla por cross-device, copia
	if err := os.Rename(src, dst); err == nil {
		return nil
	}
	// copia manual
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return os.Remove(src)
}
