package infra

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Platform struct {
	GOOS     string // linux, darwin, windows
	GOARCH   string // amd64, arm64, arm
	DistroID string // ubuntu, debian, raspbian, etc. (best-effort)
}

func DetectPlatform() (*Platform, error) {
	p := &Platform{
		GOOS:   runtime.GOOS,
		GOARCH: runtime.GOARCH,
	}
	if p.GOOS == "linux" {
		// best-effort: detecta distro
		if data, err := os.ReadFile("/etc/os-release"); err == nil {
			for _, ln := range strings.Split(string(data), "\n") {
				if strings.HasPrefix(ln, "ID=") {
					p.DistroID = strings.Trim(strings.TrimPrefix(ln, "ID="), `"`)
					break
				}
			}
		}
	}
	return p, nil
}

// ArchKey para mapear a las claves del TOML (linux_amd64, linux_arm64, linux_arm)
func (p *Platform) ArchKey() (string, error) {
	if p.GOOS != "linux" {
		return "", fmt.Errorf("solo Linux soportado por ahora (Ubuntu y Raspberry Pi OS)")
	}
	switch p.GOARCH {
	case "amd64":
		return "linux_amd64", nil
	case "arm64":
		return "linux_arm64", nil
	case "arm":
		return "linux_arm", nil
	default:
		return "", fmt.Errorf("arquitectura no soportada: %s", p.GOARCH)
	}
}

func HasCmd(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
