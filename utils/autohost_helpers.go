package utils

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Exec(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Stdout, c.Stderr = os.Stdout, os.Stderr
	return c.Run()
}
func ExecShell(script string) error {
	// bash con -e (stop on error) y -o pipefail
	return Exec("bash", "-eo", "pipefail", "-c", script)
}

func ExecWithDir(dir string, cmdName string, args ...string) error {
	cmd := exec.Command(cmdName, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Confirm(prompt string) bool {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))
	return input == "y" || input == "yes"
}

func AskOption(prompt string, options []string) string {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println(prompt)
		for i, opt := range options {
			fmt.Printf("[%d] %s\n", i+1, opt)
		}
		fmt.Print("Elige una opción: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if i, err := strconv.Atoi(input); err == nil && i >= 1 && i <= len(options) {
			return options[i-1]
		}
		fmt.Println("❌ Opción inválida, intenta de nuevo.")
	}
}

func ValidPort(portStr string) (int, error) {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, fmt.Errorf("puerto inválido: %s", portStr)
	}

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return 0, fmt.Errorf("el puerto %d ya está en uso", port)
	}
	ln.Close()

	return port, nil
}

func AskAppPort(reader *bufio.Reader, message string, port string) string {
	for {
		port = AskInput(reader, message, port)
		// Validar que el puerto sea un número válido y esté libre
		if _, err := ValidPort(port); err == nil {
			break
		} else {
			fmt.Println("❌", err)
		}
	}
	return port
}
