// tests/e2e/e2e_setup_test.go
//go:build e2e

package e2e

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func runInRepoRoot(t *testing.T, cmd string) string {
	t.Helper()
	// Localiza la raíz del módulo (directorio del go.mod)
	gomod := exec.Command("go", "env", "GOMOD")
	b, err := gomod.Output()
	if err != nil {
		t.Fatalf("no pude obtener go env GOMOD: %v", err)
	}
	repoRoot := strings.TrimSuffix(strings.TrimSpace(string(b)), "/go.mod")

	c := exec.Command("bash", "-lc", cmd)
	c.Dir = repoRoot
	var out bytes.Buffer
	var errb bytes.Buffer
	c.Stdout, c.Stderr = &out, &errb
	if err := c.Run(); err != nil {
		t.Fatalf("cmd falló (dir=%s): %s\nstderr:\n%s", repoRoot, cmd, errb.String())
	}
	return out.String()
}

func Test_SetupE2E(t *testing.T) {
	vm := getenv("VM_NAME", "autohost-test")
	ci := getenv("CI", "true")
	assume := getenv("ASSUME_YES", "-y")

	// 1) Levanta VM con tu script (sin shell interactiva)
	runInRepoRoot(t, "NO_SHELL=1 scripts/autohost-multipass.sh run")

	// 2) Verifica que la VM esté Running (si no, fallar)
	list := runInRepoRoot(t, "multipass list --format csv || true")
	if !strings.Contains(list, vm+",Running") {
		t.Fatalf("la VM %q no está en Running.\n'multipass list' devolvió:\n%s", vm, list)
	}

	// 3) Detecta el nombre del binario dentro de la VM (autohost )
	bin := "autohost"
	which := runInRepoRoot(t, "multipass exec "+vm+" -- bash -lc 'command -v autohost || true'")

	// 4) Ejecuta setup dentro de la VM (CI/ASSUME_YES para no interactivo)
	out := runInRepoRoot(t, "multipass exec "+vm+" -- bash -lc 'CI="+ci+" ASSUME_YES="+assume+" "+bin+" setup'")
	if !(strings.Contains(strings.ToLower(out), "docker") ||
		strings.Contains(strings.ToLower(out), "instalado") ||
		strings.Contains(strings.ToLower(out), "configured") ||
		strings.Contains(out, "✅")) {
		t.Fatalf("salida inesperada del setup:\n%s", out)
	}
	t.Logf("✅ setup completado:\n%s", out)

	// 5) Idempotencia: segunda ejecución
	out2 := runInRepoRoot(t, "multipass exec "+vm+" -- bash -lc 'CI="+ci+" ASSUME_YES="+assume+" "+bin+" setup'")
	if !(strings.Contains(strings.ToLower(out2), "ya está instalado") ||
		strings.Contains(strings.ToLower(out2), "already configured") ||
		strings.Contains(strings.ToLower(out2), "configured")) {
		t.Fatalf("segunda ejecución no mostró mensaje de idempotencia:\n%s", out2)
	}
	t.Log("✅ idempotencia confirmada")

	// 6) Limpieza (siempre)
	t.Cleanup(func() {
		_ = runInRepoRoot(t, "NO_SHELL=1 scripts/autohost-multipass.sh delete || true")
	})
}
