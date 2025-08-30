package utils

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

const (
	yellow = "\033[33m"
	reset  = "\033[0m"
)

// AskInput muestra un prompt y lee la entrada del usuario, usando un valor por defecto si está vacío.
func AskInput(reader *bufio.Reader, prompt, def string) string {
	fmt.Printf("%s [%s%s%s]: ", prompt, yellow, def, reset)
	val, _ := reader.ReadString('\n')
	val = strings.TrimSpace(val)
	if val == "" {
		return def
	}
	return val
}

func WithAppName(fn func(appName string)) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		appName := args[0]
		fn(appName)
	}
}
