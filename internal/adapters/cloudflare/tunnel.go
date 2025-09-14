package cloudflare

import (
	"autohost-cli/utils"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func TunnelCloudflare(domain string) {
	fmt.Println("🌐 Iniciando túnel de Cloudflare...")

	if !utils.IsInitialized() {
		fmt.Println("⚠️ Ejecuta `autohost init` primero.")
		return
	}

	fmt.Printf("⚙️ Creando túnel para %s...\n", domain)

	// Crear el túnel
	createCmd := exec.Command("cloudflared", "tunnel", "create", "autohost-tunnel")
	createCmd.Stdout = os.Stdout
	createCmd.Stderr = os.Stderr
	err := createCmd.Run()
	if err != nil {
		fmt.Println("❌ Error al crear túnel:", err)
		return
	}

	// Mover archivo del túnel
	tunnelFile := filepath.Join(os.Getenv("HOME"), ".cloudflared", "autohost-tunnel.json")
	target := filepath.Join(utils.GetAutohostDir(), "cloudflare", "tunnel.json")

	if err := utils.CopyFile(tunnelFile, target); err != nil {
		fmt.Println("⚠️ No se pudo mover el archivo del túnel:", err)
	}

	// Enlazar túnel al dominio
	routeCmd := exec.Command("cloudflared", "tunnel", "route", "dns", "autohost-tunnel", domain)
	routeCmd.Stdout = os.Stdout
	routeCmd.Stderr = os.Stderr
	err = routeCmd.Run()
	if err != nil {
		fmt.Println("❌ Error al configurar ruta DNS:", err)
	} else {
		fmt.Println("✅ Túnel creado y vinculado al dominio:", domain)
	}

	// Guardar config
	cfg := utils.Config{
		Tunnel: "cloudflare",
		Domain: domain,
	}
	if err := utils.SaveConfig(cfg); err != nil {
		fmt.Println("⚠️ Error al guardar config:", err)
	}

	// utils.SaveStatus("cloudflare_tunnel", true)
	// utils.SaveStatus("cloudflare_domain", domain)
}
