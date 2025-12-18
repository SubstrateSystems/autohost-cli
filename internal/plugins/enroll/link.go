package enroll

import (
	"autohost-cli/internal/plugins/enroll/config"
	"autohost-cli/internal/plugins/enroll/http"
	"autohost-cli/internal/plugins/enroll/services"
	"autohost-cli/internal/plugins/enroll/types"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func NewLinkCmd() *cobra.Command {
	var api, token, name string

	cmd := &cobra.Command{
		Use:   "link",
		Short: "Enlaza este nodo con la API de AutoHost",
		RunE: func(cmd *cobra.Command, args []string) error {
			if api == "" || token == "" {
				return errors.New("--api y --token son obligatorios")
			}

			client := http.NewUserClient(api, token)
			nodeData := services.GetAgentData()

			if name != "" {
				nodeData.HostName = name
			}

			req := types.NodeRquest{
				ErollToken:   strings.TrimSpace(token),
				HostName:     strings.TrimSpace(nodeData.HostName),
				IPLocal:      strings.TrimSpace(nodeData.IPLocal),
				OS:           strings.TrimSpace(nodeData.OS),
				Arch:         strings.TrimSpace(nodeData.Arch),
				VersionAgent: strings.TrimSpace(nodeData.VersionAgent),
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			var resp types.NodeResponse

			status, err := client.PostJSON(ctx, "/v1/enrollments/enroll", req, &resp)
			if err != nil {
				log.Fatalf("Error en la petición: %v", err)
			}

			fmt.Println("Código HTTP:", status)
			fmt.Printf("Respuesta del servidor: {%s}\n", resp.ApiToken)

			fmt.Println()
			fmt.Println("💾 Guardando configuración...")
			cfg := config.AgentConfig{
				ApiToken: resp.ApiToken,
				ApiURL:   api,
			}
			err = config.Save(cfg)
			if err != nil {
				log.Fatalf("Error guardando configuración: %v", err)
			}

			fmt.Println("✅ Configuración actualizada en /etc/autohost/config.yaml")
			fmt.Println()
			fmt.Println("📝 Próximos pasos:")
			fmt.Println("  1. Habilitar servicio: sudo systemctl enable autohost-agent")
			fmt.Println("  2. Iniciar servicio:   sudo systemctl start autohost-agent")
			fmt.Println("  3. Verificar estado:   sudo systemctl status autohost-agent")

			return nil
		},
	}

	cmd.Flags().StringVar(&api, "api", "", "URL base del API (https://...)")
	cmd.Flags().StringVar(&token, "token", "", "Token de sesión de usuario")
	cmd.Flags().StringVar(&name, "name", "", "Nombre lógico del nodo")

	_ = cmd.MarkFlagRequired("api")
	_ = cmd.MarkFlagRequired("token")
	return cmd
}
