package tailscale

// import (
// 	"autohost-cli/internal/adapters/infra"
// 	"fmt"
// 	"strings"
// )

// func TailscaleSplitDNS(subdomain, appIP string) error {
// 	domain, _ := cmd.Flags().GetString("domain")
// 	nsStr, _ := cmd.Flags().GetString("nameservers")
// 	searchStr, _ := cmd.Flags().GetString("search-paths")
// 	tailnet, _ := cmd.Flags().GetString("tailnet")

// 	if domain == "" || nsStr == "" {
// 		return fmt.Errorf("flags requeridas: --domain y --nameservers (separados por coma si son varios)")
// 	}

// 	nameservers := splitAndTrim(nsStr)
// 	searchPaths := splitAndTrim(searchStr)

// 	fmt.Println("⚙️  Configurando Split DNS con Terraform...")
// 	err := infra.ConfigureSplitDNSWithTerraform(infra.SplitDNSOpts{
// 		// Tailnet:      tailnet,
// 		Domain: subdomain,
// 		// Nameservers:  nameservers,
// 		// SearchPaths:  searchPaths,
// 		APIKeyEnvVar: "TAILSCALE_API_KEY",
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println("✅ Split DNS aplicado.")
// 	return nil
// }

// func splitAndTrim(s string) []string {
// 	if strings.TrimSpace(s) == "" {
// 		return nil
// 	}
// 	parts := strings.Split(s, ",")
// 	out := make([]string, 0, len(parts))
// 	for _, p := range parts {
// 		t := strings.TrimSpace(p)
// 		if t != "" {
// 			out = append(out, t)
// 		}
// 	}
// 	return out
// }
