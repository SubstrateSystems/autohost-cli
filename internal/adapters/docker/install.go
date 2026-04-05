package docker

import (
	"autohost-cli/utils"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type osRelease struct {
	ID     string
	IDLike string
}

func Install() error {
	if runningInContainer() {
		fmt.Println("⚠️  Detecté contenedor. No instalo Docker aquí. Usa el socket del host o dind para pruebas.")
		return nil
	}
	if dockerAvailable() {
		fmt.Println("✅ Docker ya está instalado.")
		return nil
	}
	fmt.Println("🔄 Instalando Docker...")

	// Asegura curl
	if err := ensureCurl(); err != nil {
		panic("❌ No pude instalar/ubicar curl: " + err.Error())
	}

	// Script oficial SIN pipe ciego
	if err := utils.ExecShell(`
set -e
tmp="$(mktemp)"
curl -fsSL https://get.docker.com -o "$tmp"
sh "$tmp"
rm -f "$tmp"
`); err != nil {
		panic("❌ Error ejecutando el instalador de Docker: " + err.Error())
	}

	// Arrancar/enable del daemon (si hay systemd)
	if systemctlAvailable() {
		_ = utils.Exec("sudo", "systemctl", "enable", "--now", "docker")
	} else {
		// fallback best-effort
		_ = utils.Exec("sudo", "service", "docker", "start")
	}

	// Verificar CLI + daemon
	if err := exec.Command("sudo", "docker", "--version").Run(); err != nil {
		panic("❌ Docker CLI no quedó instalado correctamente.")
	}
	if err := exec.Command("sudo", "docker", "info").Run(); err != nil {
		fmt.Println("⚠️  Docker instalado, pero el daemon no responde aún. Revisa el servicio o reinicia el host.")
	} else {
		fmt.Println("✅ Docker instalado y en ejecución.")
	}
	return nil
}

func ensureCurl() error {
	osr := readOSRelease()
	id := osr.ID + " " + osr.IDLike

	switch {
	case strings.Contains(id, "debian") || strings.Contains(id, "ubuntu"):
		return utils.ExecShell(`sudo apt-get update -y && sudo apt-get install -y curl ca-certificates && sudo update-ca-certificates`)
	case strings.Contains(id, "rhel") || strings.Contains(id, "centos") || strings.Contains(id, "rocky") || strings.Contains(id, "almalinux"):
		return utils.ExecShell(`sudo yum install -y curl ca-certificates || sudo dnf install -y curl ca-certificates`)
	case strings.Contains(id, "fedora"):
		return utils.ExecShell(`sudo dnf install -y curl ca-certificates`)
	case strings.Contains(id, "amzn"): // Amazon Linux
		return utils.ExecShell(`sudo yum install -y curl ca-certificates || sudo dnf install -y curl ca-certificates`)
	case strings.Contains(id, "alpine"):
		return utils.ExecShell(`sudo apk add --no-cache curl ca-certificates && sudo update-ca-certificates`)
	case strings.Contains(id, "suse") || strings.Contains(id, "sles") || strings.Contains(id, "opensuse"):
		return utils.ExecShell(`sudo zypper --non-interactive install -y curl ca-certificates`)
	default:
		return utils.Exec("which", "curl")
	}
}

func systemctlAvailable() bool { return exec.Command("which", "systemctl").Run() == nil }

func readOSRelease() osRelease {
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return osRelease{}
	}
	defer f.Close()
	kv := map[string]string{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		k := parts[0]
		v := strings.Trim(parts[1], `"'`)
		kv[k] = v
	}
	return osRelease{ID: kv["ID"], IDLike: kv["ID_LIKE"]}
}

func runningInContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	// opcional: variable para forzar
	return os.Getenv("AUTOHOST_IN_CONTAINER") == "true"
}

func dockerAvailable() bool { return exec.Command("docker", "version").Run() == nil }
