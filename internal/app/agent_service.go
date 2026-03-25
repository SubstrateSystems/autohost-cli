package app

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	agentBinaryName  = "autohost-agent"
	agentVersion     = "v0.1.0"
	agentInstallPath = "/usr/local/bin"
	agentConfigPath  = "/etc/autohost"
	agentServicePath = "/etc/systemd/system"
)

// AgentService manages installation of the AutoHost agent binary.
type AgentService struct{}

// Install downloads the agent binary, example config, and systemd service file,
// then installs them to the appropriate system paths.
func (s *AgentService) Install() error {
	needsSudo := os.Geteuid() != 0
	if needsSudo {
		if _, err := exec.LookPath("sudo"); err != nil {
			return fmt.Errorf("este comando requiere privilegios de root. Por favor, ejecuta como root o instala sudo")
		}
	}

	fmt.Println("📦 Instalando AutoHost Agent...")
	fmt.Println()

	tmpDir, err := os.MkdirTemp("", "autohost-agent-*")
	if err != nil {
		return fmt.Errorf("error creando directorio temporal: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	binaryPath := filepath.Join(tmpDir, agentBinaryName)
	configFile := filepath.Join(tmpDir, "agent.yaml")
	serviceFile := filepath.Join(tmpDir, "autohost-agent.service")

	fmt.Println("1️⃣  Descargando binario...")
	downloadURL := fmt.Sprintf("https://github.com/SubstrateSystems/autohost-agent/releases/download/%s/autohost-agent-linux-amd64", agentVersion)
	if err := agentDownloadFile(downloadURL, binaryPath); err != nil {
		return fmt.Errorf("error descargando binario: %w", err)
	}

	fmt.Println("2️⃣  Descargando configuración de ejemplo...")
	configURL := fmt.Sprintf("https://raw.githubusercontent.com/SubstrateSystems/autohost-agent/%s/configs/agent.yaml", agentVersion)
	if err := agentDownloadFile(configURL, configFile); err != nil {
		return fmt.Errorf("error descargando configuración: %w", err)
	}

	fmt.Println("3️⃣  Descargando archivo de servicio systemd...")
	serviceURL := fmt.Sprintf("https://raw.githubusercontent.com/SubstrateSystems/autohost-agent/%s/autohost-agent.service", agentVersion)
	if err := agentDownloadFile(serviceURL, serviceFile); err != nil {
		return fmt.Errorf("error descargando servicio systemd: %w", err)
	}

	fmt.Println("4️⃣  Creando directorios...")
	if err := agentRunCmd(needsSudo, "mkdir", "-p", agentConfigPath); err != nil {
		return fmt.Errorf("error creando directorio de configuración: %w", err)
	}

	fmt.Println("4️⃣.1️⃣  Verificando usuario del servicio...")
	if err := ensureAgentSystemUser(needsSudo); err != nil {
		return fmt.Errorf("error configurando usuario del servicio: %w", err)
	}

	fmt.Println("5️⃣  Instalando binario en", agentInstallPath, "...")
	destBinary := filepath.Join(agentInstallPath, agentBinaryName)
	if err := agentRunCmd(needsSudo, "cp", binaryPath, destBinary); err != nil {
		return fmt.Errorf("error copiando binario: %w", err)
	}
	if err := agentRunCmd(needsSudo, "chmod", "+x", destBinary); err != nil {
		return fmt.Errorf("error estableciendo permisos de ejecución: %w", err)
	}

	fmt.Println("6️⃣  Instalando configuración...")
	destConfig := filepath.Join(agentConfigPath, "config.yaml")
	if err := agentRunCmd(needsSudo, "cp", configFile, destConfig); err != nil {
		return fmt.Errorf("error copiando configuración: %w", err)
	}
	if err := agentRunCmd(needsSudo, "chown", "root:autohost", destConfig); err != nil {
		return fmt.Errorf("error estableciendo propietario del config: %w", err)
	}
	if err := agentRunCmd(needsSudo, "chmod", "640", destConfig); err != nil {
		return fmt.Errorf("error estableciendo permisos de configuración: %w", err)
	}

	fmt.Println("7️⃣  Instalando servicio systemd...")
	destService := filepath.Join(agentServicePath, "autohost-agent.service")
	if err := agentRunCmd(needsSudo, "cp", serviceFile, destService); err != nil {
		return fmt.Errorf("error copiando servicio: %w", err)
	}

	fmt.Println("8️⃣  Recargando systemd...")
	if err := agentRunCmd(needsSudo, "systemctl", "daemon-reload"); err != nil {
		return fmt.Errorf("error recargando systemd: %w", err)
	}

	sudoStr := ""
	if needsSudo {
		sudoStr = "sudo "
	}
	fmt.Println()
	fmt.Println("✅ Instalación completada exitosamente!")
	fmt.Println()
	fmt.Println("📍 Ubicaciones:")
	fmt.Printf("  • Binario:        %s/%s\n", agentInstallPath, agentBinaryName)
	fmt.Printf("  • Configuración:  %s/config.yaml\n", agentConfigPath)
	fmt.Printf("  • Servicio:       %s/autohost-agent.service\n", agentServicePath)
	fmt.Println()
	fmt.Println("📝 Próximos pasos:")
	fmt.Printf("  1. Editar configuración:  %snano %s/config.yaml\n", sudoStr, agentConfigPath)
	fmt.Printf("  2. Habilitar servicio:    %ssystemctl enable autohost-agent\n", sudoStr)
	fmt.Printf("  3. Iniciar servicio:      %ssystemctl start autohost-agent\n", sudoStr)
	fmt.Printf("  4. Verificar estado:      %ssystemctl status autohost-agent\n", sudoStr)
	fmt.Println()
	return nil
}

// agentDownloadFile downloads url to destination using curl or wget.
func agentDownloadFile(url, destination string) error {
	if _, err := exec.LookPath("curl"); err == nil {
		cmd := exec.Command("curl", "-L", "-o", destination, url)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	if _, err := exec.LookPath("wget"); err == nil {
		cmd := exec.Command("wget", "-O", destination, url)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	return fmt.Errorf("ni curl ni wget están disponibles. Por favor instala uno de ellos:\n  • apt install curl  (Debian/Ubuntu)\n  • yum install curl  (RHEL/CentOS)\n  • apk add curl      (Alpine)")
}

// agentRunCmd runs name with args, prepending sudo when sudo=true.
func agentRunCmd(sudo bool, name string, args ...string) error {
	cmdName, cmdArgs := name, args
	if sudo {
		cmdArgs = append([]string{name}, args...)
		cmdName = "sudo"
	}
	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func agentRunCmdQuiet(sudo bool, name string, args ...string) error {
	cmdName, cmdArgs := name, args
	if sudo {
		cmdArgs = append([]string{name}, args...)
		cmdName = "sudo"
	}
	cmd := exec.Command(cmdName, cmdArgs...)
	return cmd.Run()
}

func ensureAgentSystemUser(needsSudo bool) error {
	if err := agentRunCmdQuiet(needsSudo, "id", "-u", "autohost"); err == nil {
		return nil
	}

	nologin := "/usr/sbin/nologin"
	if _, err := os.Stat(nologin); err != nil {
		nologin = "/usr/bin/false"
	}

	if err := agentRunCmd(needsSudo,
		"useradd",
		"--system",
		"--no-create-home",
		"--shell", nologin,
		"--comment", "Autohost Agent",
		"autohost",
	); err != nil {
		return err
	}

	if err := agentRunCmdQuiet(needsSudo, "getent", "group", "docker"); err == nil {
		_ = agentRunCmd(needsSudo, "usermod", "-aG", "docker", "autohost")
	}

	return nil
}
