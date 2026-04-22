package app

import (
	"autohost-cli/internal/adapters/agentconfig"
	"autohost-cli/internal/adapters/enrollapi"
	"autohost-cli/internal/domain"
	"context"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"
)

// EnrollService handles node enrollment with the AutoHost cloud API.
type EnrollService struct{}

// Link gathers local node data, enrolls this node with the AutoHost API using
// the provided user token, and persists the resulting agent token to disk.
func (s *EnrollService) Link(api, token, name string) error {
	nd := s.gatherNodeData()
	if name != "" {
		nd.HostName = name
	}

	req := domain.NodeRequest{
		EnrollToken:  strings.TrimSpace(token),
		HostName:     strings.TrimSpace(nd.HostName),
		IPLocal:      strings.TrimSpace(nd.IPLocal),
		OS:           strings.TrimSpace(nd.OS),
		Arch:         strings.TrimSpace(nd.Arch),
		VersionAgent: strings.TrimSpace(nd.VersionAgent),
	}

	client := enrollapi.NewUserClient(api, token)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var resp domain.NodeResponse
	_, err := client.PostJSON(ctx, "/v1/enrollments/enroll", req, &resp)
	if err != nil {
		return fmt.Errorf("enrollment fallido: %w", err)
	}
	fmt.Println("💾 Guardando configuración...")

	if err := agentconfig.Save(agentconfig.AgentConfig{
		ApiToken:     resp.ApiToken,
		RefreshToken: resp.RefreshToken,
		ApiURL:       api,
		NodeID:       resp.NodeID,
	}); err != nil {
		return fmt.Errorf("error guardando configuración: %w", err)
	}

	fmt.Println("✅ Configuración actualizada en /etc/autohost/config.yaml")
	fmt.Println()
	fmt.Println("📝 Próximos pasos:")
	fmt.Println("  1. Habilitar servicio: sudo systemctl enable autohost-agent")
	fmt.Println("  2. Iniciar servicio:   sudo systemctl start autohost-agent")
	fmt.Println("  3. Verificar estado:   sudo systemctl status autohost-agent")
	return nil
}

func (s *EnrollService) gatherNodeData() *domain.NodeData {
	nd := &domain.NodeData{}
	if out, err := exec.Command("hostname").Output(); err == nil {
		nd.HostName = strings.TrimSpace(string(out))
	}
	if out, err := exec.Command("uname", "-o").Output(); err == nil {
		nd.OS = strings.TrimSpace(string(out))
	}
	if out, err := exec.Command("uname", "-m").Output(); err == nil {
		nd.Arch = strings.TrimSpace(string(out))
	}
	nd.IPLocal = detectLocalIP()
	nd.VersionAgent = "0.1.0"
	return nd
}

func detectLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "unknown"
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}
