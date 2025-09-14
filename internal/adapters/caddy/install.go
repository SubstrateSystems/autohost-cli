package caddy

import (
	"autohost-cli/utils"
	"fmt"
)

func InstallCaddy() {
	fmt.Println("🚀 Instalando Caddy...")
	utils.ExecShell(`
	sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https curl &&
		curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg &&
		curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list &&
		sudo apt update && sudo apt install caddy
	`)
	utils.ExecShell("sudo systemctl enable caddy")
	utils.ExecShell("sudo systemctl start caddy")
	fmt.Println("✅ Caddy instalado y activado correctamente.")
}
