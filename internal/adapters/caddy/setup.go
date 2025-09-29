package caddy

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

func EnsureCaddySnippetsSetup(ctx context.Context) error {
	const (
		groupName      = "autohost"
		snippetsDir    = "/etc/caddy/autohost"
		caddyfilePath  = "/etc/caddy/Caddyfile"
		sudoersPath    = "/etc/sudoers.d/autohost-caddy"
		importLine     = "import /etc/caddy/autohost/*.caddy"
		sudoersContent = "%autohost ALL=NOPASSWD:/usr/bin/systemctl reload caddy,/usr/bin/systemctl restart caddy,/usr/bin/caddy reload\n"
	)

	u, err := user.Current()
	if err != nil {
		return fmt.Errorf("user.Current: %w", err)
	}
	username := u.Username

	runSudo := func(args ...string) error {
		c := exec.CommandContext(ctx, "sudo", args...)
		c.Stdout, c.Stderr = os.Stdout, os.Stderr
		return c.Run()
	}

	outSudo := func(args ...string) (string, error) {
		c := exec.CommandContext(ctx, "sudo", args...)
		var buf bytes.Buffer
		c.Stdout, c.Stderr = &buf, &buf
		err := c.Run()
		return buf.String(), err
	}

	// 1) Grupo autohost
	if _, err := outSudo("getent", "group", groupName); err != nil {
		if err := runSudo("groupadd", groupName); err != nil {
			return fmt.Errorf("groupadd %s: %w", groupName, err)
		}
	}

	// 2) Usuario en el grupo
	groupsOut, _ := exec.Command("id", "-nG", username).Output()
	if !strings.Contains(" "+string(groupsOut)+" ", " "+groupName+" ") {
		if err := runSudo("usermod", "-aG", groupName, username); err != nil {
			return fmt.Errorf("usermod -aG %s %s: %w", groupName, username, err)
		}
	}

	// 3) Directorio snippets con root:autohost y 2775
	if err := runSudo("mkdir", "-p", snippetsDir); err != nil {
		return fmt.Errorf("mkdir %s: %w", snippetsDir, err)
	}
	if err := runSudo("chown", "-R", "root:"+groupName, snippetsDir); err != nil {
		return fmt.Errorf("chown %s: %w", snippetsDir, err)
	}
	// 2775 = setgid + rwxrwxr-x
	if err := runSudo("chmod", "2775", snippetsDir); err != nil {
		return fmt.Errorf("chmod 2775 %s: %w", snippetsDir, err)
	}

	// 4) Asegurar import en /etc/caddy/Caddyfile (si falta, lo agregamos al final)
	//    Estrategia: crear archivo temp y usar `sudo install` para copiar con modo 0644.
	tmp := filepath.Join(os.TempDir(), "Caddyfile.autohost.tmp")
	content, readErr := os.ReadFile(caddyfilePath)
	if readErr != nil {
		return fmt.Errorf("leer %s: %w", caddyfilePath, readErr)
	}
	if !strings.Contains(string(content), importLine) {
		newContent := strings.TrimRight(string(content), "\n") + "\n\n" + importLine + "\n"
		if err := os.WriteFile(tmp, []byte(newContent), 0o644); err != nil {
			return fmt.Errorf("escribir tmp Caddyfile: %w", err)
		}
		if err := runSudo("install", "-m", "0644", tmp, caddyfilePath); err != nil {
			return fmt.Errorf("actualizar %s: %w", caddyfilePath, err)
		}
	}

	// 5) Sudoers para recargar Caddy sin contraseña
	//    Creamos/actualizamos de forma segura con permisos 0440.
	tmpS := filepath.Join(os.TempDir(), "autohost-caddy.sudoers.tmp")
	if err := os.WriteFile(tmpS, []byte(sudoersContent), 0o440); err != nil {
		return fmt.Errorf("escribir tmp sudoers: %w", err)
	}
	if err := runSudo("install", "-m", "0440", tmpS, sudoersPath); err != nil {
		return fmt.Errorf("instalar sudoers: %w", err)
	}

	// 6) Recargar Caddy (para tomar el import si lo agregamos recién)
	_ = runSudo("systemctl", "reload", "caddy")

	// 7) Advertir sobre re-login si el grupo no está activo en esta sesión
	groupsOut2, _ := exec.Command("id", "-nG", username).Output()
	if !strings.Contains(" "+string(groupsOut2)+" ", " "+groupName+" ") {
		fmt.Println("ℹ️  Se agregó tu usuario al grupo 'autohost'. Para que surta efecto, cierra sesión y entra de nuevo,")
		fmt.Println("    o ejecuta:   newgrp autohost")
	}

	fmt.Printf("✅ Snippets de Caddy listos en %s (root:%s 2775) y sudoers configurado.\n", snippetsDir, groupName)
	return nil
}
