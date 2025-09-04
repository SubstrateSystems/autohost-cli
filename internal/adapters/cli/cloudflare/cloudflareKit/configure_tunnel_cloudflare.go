package cloudflarekit

import (
	"autohost-cli/utils"
	"fmt"
)

func ConfigureCloudflareTunnel(domain string) {
	fmt.Println("⚙️ Configurando Cloudflare Tunnel para:", domain)
	utils.ExecShell("cloudflared tunnel create autohost-tunnel")
	utils.ExecShell(fmt.Sprintf("cloudflared tunnel route dns autohost-tunnel %s", domain))
	fmt.Println("✅ Túnel configurado correctamente.")
}
