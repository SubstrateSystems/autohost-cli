package agent

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	binaryName  = "autohost-agent"
	version     = "v0.1.0"
	installPath = "/usr/local/bin"
	configPath  = "/etc/autohost"
	servicePath = "/etc/systemd/system"
)

func agentInstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install the AutoHost agent",
		Long:  "Download and install the AutoHost agent binary, configuration, and systemd service",
		Example: `  autohost-cli agent install
  
  # After installation, configure and start:
  sudo nano /etc/autohost/config.yaml
  sudo systemctl enable autohost-agent
  sudo systemctl start autohost-agent`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return installAgent()
		},
	}

	return cmd
}

func installAgent() error {
	// Detectar si necesitamos sudo
	needsSudo := os.Geteuid() != 0
	sudoPrefix := ""
	if needsSudo {
		if _, err := exec.LookPath("sudo"); err != nil {
			return fmt.Errorf("este comando requiere privilegios de root. Por favor, ejecuta como root o instala sudo")
		}
		sudoPrefix = "sudo "
	}

	fmt.Println("📦 Instalando AutoHost Agent...")
	fmt.Println()

	// 1. Crear directorio temporal
	tmpDir, err := os.MkdirTemp("", "autohost-agent-*")
	if err != nil {
		return fmt.Errorf("error creando directorio temporal: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	binaryPath := filepath.Join(tmpDir, binaryName)
	configFile := filepath.Join(tmpDir, "agent.yaml")
	serviceFile := filepath.Join(tmpDir, "autohost-agent.service")

	// 2. Descargar binario
	fmt.Println("1️⃣  Descargando binario...")
	downloadURL := fmt.Sprintf("https://github.com/SubstrateSystems/autohost-agent/releases/download/%s/autohost-agent-linux-amd64", version)
	if err := downloadFile(downloadURL, binaryPath); err != nil {
		return fmt.Errorf("error descargando binario: %w", err)
	}

	// 3. Descargar configuración de ejemplo
	fmt.Println("2️⃣  Descargando configuración de ejemplo...")
	configURL := fmt.Sprintf("https://raw.githubusercontent.com/SubstrateSystems/autohost-agent/%s/configs/agent.yaml", version)
	if err := downloadFile(configURL, configFile); err != nil {
		return fmt.Errorf("error descargando configuración: %w", err)
	}

	// 4. Descargar archivo de servicio systemd
	fmt.Println("3️⃣  Descargando archivo de servicio systemd...")
	serviceURL := fmt.Sprintf("https://raw.githubusercontent.com/SubstrateSystems/autohost-agent/%s/autohost-agent.service", version)
	if err := downloadFile(serviceURL, serviceFile); err != nil {
		return fmt.Errorf("error descargando servicio systemd: %w", err)
	}

	// 5. Crear directorio de configuración
	fmt.Println("4️⃣  Creando directorios...")
	if err := runCommand(sudoPrefix+"mkdir", "-p", configPath); err != nil {
		return fmt.Errorf("error creando directorio de configuración: %w", err)
	}

	// 6. Instalar binario
	fmt.Println("5️⃣  Instalando binario en", installPath, "...")
	destBinary := filepath.Join(installPath, binaryName)
	if err := runCommand(sudoPrefix+"cp", binaryPath, destBinary); err != nil {
		return fmt.Errorf("error copiando binario: %w", err)
	}
	if err := runCommand(sudoPrefix+"chmod", "+x", destBinary); err != nil {
		return fmt.Errorf("error estableciendo permisos de ejecución: %w", err)
	}

	// 7. Instalar configuración
	fmt.Println("6️⃣  Instalando configuración...")
	destConfig := filepath.Join(configPath, "config.yaml")
	if err := runCommand(sudoPrefix+"cp", configFile, destConfig); err != nil {
		return fmt.Errorf("error copiando configuración: %w", err)
	}
	if err := runCommand(sudoPrefix+"chmod", "600", destConfig); err != nil {
		return fmt.Errorf("error estableciendo permisos de configuración: %w", err)
	}

	// 8. Instalar servicio systemd
	fmt.Println("7️⃣  Instalando servicio systemd...")
	destService := filepath.Join(servicePath, "autohost-agent.service")
	if err := runCommand(sudoPrefix+"cp", serviceFile, destService); err != nil {
		return fmt.Errorf("error copiando servicio: %w", err)
	}

	// 9. Recargar systemd
	fmt.Println("8️⃣  Recargando systemd...")
	if err := runCommand(sudoPrefix+"systemctl", "daemon-reload"); err != nil {
		return fmt.Errorf("error recargando systemd: %w", err)
	}

	// Mensajes finales
	fmt.Println()
	fmt.Println("✅ Instalación completada exitosamente!")
	fmt.Println()
	fmt.Println("📍 Ubicaciones:")
	fmt.Printf("  • Binario:        %s/%s\n", installPath, binaryName)
	fmt.Printf("  • Configuración:  %s/config.yaml\n", configPath)
	fmt.Printf("  • Servicio:       %s/autohost-agent.service\n", servicePath)
	fmt.Println()
	fmt.Println("📝 Próximos pasos:")
	fmt.Printf("  1. Editar configuración:  %snano %s/config.yaml\n", sudoPrefix, configPath)
	fmt.Printf("  2. Habilitar servicio:    %ssystemctl enable autohost-agent\n", sudoPrefix)
	fmt.Printf("  3. Iniciar servicio:      %ssystemctl start autohost-agent\n", sudoPrefix)
	fmt.Printf("  4. Verificar estado:      %ssystemctl status autohost-agent\n", sudoPrefix)
	fmt.Println()

	return nil
}

func downloadFile(url, destination string) error {
	// Intentar con curl primero
	if _, err := exec.LookPath("curl"); err == nil {
		cmd := exec.Command("curl", "-L", "-o", destination, url)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	// Si no hay curl, intentar con wget
	if _, err := exec.LookPath("wget"); err == nil {
		cmd := exec.Command("wget", "-O", destination, url)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	// Si no hay ninguno, dar error claro
	return fmt.Errorf("ni curl ni wget están disponibles. Por favor instala uno de ellos:\n  • apt install curl  (Debian/Ubuntu)\n  • yum install curl  (RHEL/CentOS)\n  • apk add curl      (Alpine)")
}

func runCommand(name string, args ...string) error {
	// Si el comando contiene "sudo ", separarlo
	if len(name) > 5 && name[:5] == "sudo " {
		actualCmd := name[5:]
		args = append([]string{actualCmd}, args...)
		name = "sudo"
	}

	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
