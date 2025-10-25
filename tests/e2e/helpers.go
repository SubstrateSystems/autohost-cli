package e2e

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

// Ejecuta comandos y devuelve salida o falla el test.
func run(t *testing.T, name string, args ...string) string {
	t.Helper()
	cmd := exec.Command(name, args...)
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	if err := cmd.Run(); err != nil {
		t.Fatalf("❌ comando falló: %s %v\nstderr:\n%s", name, args, errb.String())
	}
	return strings.TrimSpace(out.String())
}

// Ejecuta comandos dentro de la VM.
func execVM(t *testing.T, vm string, command string) string {
	t.Helper()
	return run(t, "multipass", "exec", vm, "--", "bash", "-lc", command)
}
