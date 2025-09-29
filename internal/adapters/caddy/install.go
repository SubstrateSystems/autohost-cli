package caddy

import (
	"autohost-cli/internal/platform/config"
	"autohost-cli/utils"
	"fmt"
)

func InstallCaddy() error {

	fmt.Println("🚀 Instalando Caddy...")

	// Lee configuración desde TOML embebido
	gpgKeyURL := config.MustString("urls.toml", "caddy", "gpg_key_url")
	repoListURL := config.MustString("urls.toml", "caddy", "repo_list_url")
	packageName := config.MustString("urls.toml", "caddy", "package_name")
	serviceName := config.MustString("urls.toml", "caddy", "service_name")

	utils.ExecShell(fmt.Sprintf(`
		set -e
		sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https curl &&
		curl -1sLf '%s' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg &&
		curl -1sLf '%s' | sudo tee /etc/apt/sources.list.d/caddy-stable.list >/dev/null &&
		sudo apt update && sudo apt install -y %s
	`, gpgKeyURL, repoListURL, packageName))

	utils.ExecShell(fmt.Sprintf("sudo systemctl enable %s", serviceName))
	utils.ExecShell(fmt.Sprintf("sudo systemctl start %s", serviceName))

	fmt.Println("✅ Caddy instalado y activado correctamente.")
	return nil
}
