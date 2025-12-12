package agent

import (
	"autohost-cli/internal/plugins/agent/config"
	"autohost-cli/internal/plugins/agent/http"
	"autohost-cli/internal/plugins/agent/services"
	"autohost-cli/internal/plugins/agent/types"
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
			fmt.Println("Respuesta del servidor:", resp)
			cfg := config.AgentConfig{
				ApiToken: resp.ApiToken,
			}
			err = config.Save(cfg)
			if err != nil {
				log.Fatalf("Error guardando configuración: %v", err)
			}

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
