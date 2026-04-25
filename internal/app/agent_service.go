package app

import (
	"fmt"
	"os"
	"os/exec"
)

const agentInstallScriptURL = "https://raw.githubusercontent.com/mazapanuwu13/autohost-agent/main/scripts/install.sh"

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

// Update re-runs the install script to replace the binary with the latest
// (or a specific) version, then restarts the systemd service.
// Accepts the VERSION env var to pin a specific release.
func (s *AgentService) Update() error {
	if _, err := exec.LookPath("curl"); err != nil {
		return fmt.Errorf("curl no está disponible. Instálalo con: apt install curl")
	}
	if _, err := exec.LookPath("bash"); err != nil {
		return fmt.Errorf("bash no está disponible en el sistema")
	}

	fmt.Println("🔄 Actualizando AutoHost Agent...")
	fmt.Println()

	// Re-run the install script — it replaces the binary in place.
	cmd := exec.Command("bash", "-c", fmt.Sprintf("curl -fsSL %s | bash", agentInstallScriptURL))
	cmd.Stdin = nil
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("actualización del agente falló: %w", err)
	}

	// Restart the systemd service if available so the new binary is loaded
	// and the version gets reported to the API immediately.
	if _, err := exec.LookPath("systemctl"); err == nil {
		fmt.Println("🔁 Reiniciando servicio autohost-agent...")
		restart := exec.Command("systemctl", "restart", "autohost-agent")
		restart.Stdout = os.Stdout
		restart.Stderr = os.Stderr
		if err := restart.Run(); err != nil {
			fmt.Printf("⚠️  No se pudo reiniciar el servicio: %v\n", err)
			fmt.Println("   Reinicia manualmente con: systemctl restart autohost-agent")
		} else {
			fmt.Println("✅ Agente actualizado y reiniciado.")
		}
	} else {
		fmt.Println("✅ Agente actualizado. Reinicia el servicio manualmente para aplicar los cambios.")
	}

	return nil
}
