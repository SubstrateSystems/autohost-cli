package coredns

import (
	"autohost-cli/utils"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	coreDNSContainer = "coredns-autohost"
	coreDNSImage     = "coredns/coredns:latest"
)

func InstallAndRun(tailIP string) (string, error) {

	if _, err := exec.LookPath("docker"); err != nil {
		return "", fmt.Errorf("docker no está instalado en PATH: %w", err)
	}

	coreDir := utils.GetSubdir("coredns")

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
