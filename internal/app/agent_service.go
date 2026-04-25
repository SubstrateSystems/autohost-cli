package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

const (
	agentInstallScriptURL = "https://raw.githubusercontent.com/mazapanuwu13/autohost-agent/main/scripts/install.sh"
	agentRepo             = "mazapanuwu13/autohost-agent"
	agentBinName          = "autohost-agent"
	agentBinDir           = "/usr/local/bin"
)

// AgentService manages installation of the AutoHost agent binary.
type AgentService struct{}

// Install runs the upstream install.sh script via curl | bash.
// Accepts the same environment variables as the script:
//
//	VERSION          — specific release tag (e.g. v0.2.0)
//	AUTOHOST_API_URL — API URL
//	AUTOHOST_TOKEN   — enrollment token
//	AUTOHOST_NODE_ID — node identifier (default: hostname)
//	AUTOHOST_TAGS    — comma-separated tags
func (s *AgentService) Install() error {
	if _, err := exec.LookPath("curl"); err != nil {
		return fmt.Errorf("curl no está disponible. Instálalo con: apt install curl")
	}
	if _, err := exec.LookPath("bash"); err != nil {
		return fmt.Errorf("bash no está disponible en el sistema")
	}

	fmt.Println("📦 Instalando AutoHost Agent...")
	fmt.Println()

	cmd := exec.Command("bash", "-c", fmt.Sprintf("curl -fsSL %s | bash", agentInstallScriptURL))
	cmd.Stdin = nil
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("instalación del agente falló: %w", err)
	}
	return nil
}

// Update downloads only the latest agent binary and restarts the systemd
// service. Config files are left untouched.
// Set the VERSION env var to pin a specific release (e.g. VERSION=v0.3.0).
func (s *AgentService) Update() error {
	version := os.Getenv("VERSION")

	fmt.Println("🔄 Actualizando AutoHost Agent...")

	if version == "" {
		fmt.Print("   Obteniendo última versión... ")
		var err error
		version, err = fetchLatestAgentTag()
		if err != nil {
			return fmt.Errorf("no se pudo obtener la última versión: %w", err)
		}
		fmt.Println(version)
	} else {
		fmt.Printf("   Versión solicitada: %s\n", version)
	}

	arch := runtime.GOARCH // "amd64" or "arm64"
	assetName := fmt.Sprintf("%s-linux-%s", agentBinName, arch)
	downloadURL := fmt.Sprintf(
		"https://github.com/%s/releases/download/%s/%s",
		agentRepo, version, assetName,
	)

	fmt.Printf("   Descargando %s...\n", downloadURL)

	tmp, err := os.CreateTemp("", agentBinName+"-*")
	if err != nil {
		return fmt.Errorf("creando archivo temporal: %w", err)
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	resp, err := http.Get(downloadURL) //nolint:gosec // URL built from known constants + version tag
	if err != nil {
		return fmt.Errorf("descargando binario: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("descargando binario: HTTP %d", resp.StatusCode)
	}

	if _, err := io.Copy(tmp, resp.Body); err != nil {
		tmp.Close()
		return fmt.Errorf("guardando binario: %w", err)
	}
	tmp.Close()

	if err := os.Chmod(tmpPath, 0o755); err != nil {
		return fmt.Errorf("chmod binario: %w", err)
	}

	destPath := agentBinDir + "/" + agentBinName

	// Use sudo only if needed
	moveCmd := buildMoveCmd(tmpPath, destPath)
	moveCmd.Stdout = os.Stdout
	moveCmd.Stderr = os.Stderr
	if err := moveCmd.Run(); err != nil {
		return fmt.Errorf("instalando binario en %s: %w", destPath, err)
	}

	fmt.Printf("✅ Binario instalado en %s\n", destPath)

	// Restart the systemd service if available.
	if _, err := exec.LookPath("systemctl"); err == nil {
		fmt.Print("🔁 Reiniciando servicio autohost-agent... ")

		var restartCmd *exec.Cmd
		if os.Geteuid() == 0 {
			restartCmd = exec.Command("systemctl", "restart", agentBinName)
		} else {
			restartCmd = exec.Command("sudo", "systemctl", "restart", agentBinName)
		}
		restartCmd.Stdout = os.Stdout
		restartCmd.Stderr = os.Stderr

		if err := restartCmd.Run(); err != nil {
			fmt.Printf("\n⚠️  No se pudo reiniciar el servicio: %v\n", err)
			fmt.Println("   Reinicia manualmente: sudo systemctl restart autohost-agent")
			return nil
		}

		fmt.Println("✅")
		fmt.Println("✅ Agente actualizado y reiniciado correctamente.")
	} else {
		fmt.Println("✅ Agente actualizado. Reinicia el servicio para aplicar los cambios.")
	}

	return nil
}

// fetchLatestAgentTag queries the GitHub API for the latest release tag.
func fetchLatestAgentTag() (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", agentRepo)
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned HTTP %d", resp.StatusCode)
	}
	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}
	if release.TagName == "" {
		return "", fmt.Errorf("no se encontró ningún release publicado")
	}
	return release.TagName, nil
}

// buildMoveCmd returns an mv command, prefixed with sudo if the current user
// cannot write to the destination directory.
func buildMoveCmd(src, dst string) *exec.Cmd {
	if os.Geteuid() == 0 {
		return exec.Command("mv", src, dst)
	}
	return exec.Command("sudo", "mv", src, dst)
}
