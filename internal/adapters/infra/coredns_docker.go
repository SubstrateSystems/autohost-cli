// internal/infra/coredns_docker.go
package infra

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	coreDNSContainer = "coredns-autohost"
	coreDNSImage     = "coredns/coredns:latest"
)

func InstallAndRunCoreDNSWithDocker(tailIP string) (string, error) {

	if _, err := exec.LookPath("docker"); err != nil {
		return "", fmt.Errorf("docker no está instalado en PATH: %w", err)
	}

	// Preparar directorio y Corefile
	home, _ := os.UserHomeDir()
	coreDir := filepath.Join(home, ".autohost", "coredns")

	if err := os.MkdirAll(coreDir, 0o755); err != nil {
		return "", err
	}

	corefilePath := filepath.Join(coreDir, "Corefile")
	println("Directorio CoreDNS:", corefilePath)

	// Si no existe, crear uno base con la zona + bloque global "."
	if _, err := os.Stat(corefilePath); os.IsNotExist(err) {
		base := createTemplateCorefile(tailIP)

		if err := os.WriteFile(corefilePath, []byte(base), 0644); err != nil {
			return "", fmt.Errorf("no pude escribir Corefile inicial: %w", err)
		}
	} else if err != nil {
		return "", fmt.Errorf("error al verificar Corefile: %w", err)
	}

	if err := runCoreDNSContainer(corefilePath); err != nil {
		return "", err
	}

	return corefilePath, nil
}

func createTemplateCorefile(tailIP string) string {
	return fmt.Sprintf(`# CoreDNS (Docker) para AutoHost
.:53 {
    bind %s
    hosts {
        
        fallthrough
    }
    forward . 1.1.1.1
    log
    errors
}
`, tailIP)
}

func runCoreDNSContainer(corefilePath string) error {
	cmd := exec.Command(
		"docker", "run", "-d",
		"--name", coreDNSContainer,
		"--restart", "unless-stopped",
		"--network", "host",
		"-v", corefilePath+":/Corefile:ro",
		coreDNSImage, "-conf", "/Corefile",
	)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("no se pudo iniciar el contenedor CoreDNS: %w", err)
	}
	return nil
}
