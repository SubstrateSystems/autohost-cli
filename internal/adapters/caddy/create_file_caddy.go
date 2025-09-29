package caddy

import (
	"fmt"
	"os"
	"os/exec"
)

func CreateCaddyfile() error {
	caddyfilePath := "/etc/caddy/Caddyfile"

	if _, err := os.Stat(caddyfilePath); err == nil {
		fmt.Println("📄 Ya existe un Caddyfile, no se modificará.")
		return nil
	}

	content := `
http://localhost {
	respond \"🚀 AutoHost CLI: Caddy instalado y funcionando\"
}
`
	err := os.WriteFile(caddyfilePath, []byte(content), 0644)
	if err != nil {
		fmt.Println("❌ Error creando Caddyfile:", err)
		os.Exit(1)
	}

	fmt.Println("✅ Caddyfile creado en /etc/caddy/Caddyfile")

	reloadCmd := exec.Command("sudo", "systemctl", "reload", "caddy")
	reloadCmd.Stdout = os.Stdout
	reloadCmd.Stderr = os.Stderr
	if err := reloadCmd.Run(); err != nil {
		fmt.Println("⚠️ No se pudo recargar Caddy automáticamente. Hazlo manualmente con: sudo systemctl reload caddy")
	} else {
		fmt.Println("🔁 Caddy recargado con éxito.")
	}
	return nil
}
