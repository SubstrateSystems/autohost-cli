package app

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const DefaultCloudURL = "https://cloud.autohst.dev"
const DefaultAPIURL = "https://api.autohst.dev"

// const DefaultCloudURL = "http://192.168.101.2:3000"
// const DefaultAPIURL = "http://192.168.101.2:8080"

// UpService handles the interactive node enrollment flow (Tailscale-style).
type UpService struct{}

type upCallbackResult struct {
	token  string
	apiURL string
	err    error
}

// Up opens a browser for authentication and registers this node with the
// AutoHost cloud. If cloudURL is empty, DefaultCloudURL is used.
func (s *UpService) Up(cloudURL, name string) error {
	if cloudURL == "" {
		cloudURL = DefaultCloudURL
	}
	cloudURL = strings.TrimRight(cloudURL, "/")

	// Gather node hostname up-front so we can include it in the auth URL.
	nd := (&EnrollService{}).gatherNodeData()
	if name != "" {
		nd.HostName = name
	}

	// Generate a random state nonce to prevent CSRF on the local callback.
	stateBytes := make([]byte, 16)
	if _, err := rand.Read(stateBytes); err != nil {
		return fmt.Errorf("generando state nonce: %w", err)
	}
	state := hex.EncodeToString(stateBytes)

	// Bind to all interfaces so the callback is reachable from the host
	// when the CLI runs inside a container or remote machine.
	listener, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		return fmt.Errorf("iniciando servidor local de callback: %w", err)
	}
	localPort := listener.Addr().(*net.TCPAddr).Port

	// Determine the IP address the browser should reach for the callback.
	// detectLocalIP uses a UDP dial trick to find the outbound interface IP.
	callbackIP := detectLocalIP()
	if callbackIP == "unknown" {
		callbackIP = "127.0.0.1"
	}
	callbackBase := fmt.Sprintf("http://%s:%d", callbackIP, localPort)

	resultCh := make(chan upCallbackResult, 1)
	var once sync.Once

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		receivedState := q.Get("state")
		token := q.Get("token")
		apiURL := q.Get("api_url")

		if receivedState != state {
			http.Error(w, "state mismatch", http.StatusBadRequest)
			once.Do(func() {
				resultCh <- upCallbackResult{err: fmt.Errorf("state mismatch: posible ataque CSRF")}
			})
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, successHTML)
		once.Do(func() {
			resultCh <- upCallbackResult{token: token, apiURL: apiURL}
		})
	})

	srv := &http.Server{
		Handler:     mux,
		ReadTimeout: 10 * time.Second,
	}
	go func() { _ = srv.Serve(listener) }()
	defer srv.Close()

	// Build the cloud CLI-auth URL.
	authURL := fmt.Sprintf("%s/cli-auth?state=%s&callback=%s&node=%s&api_url=%s",
		cloudURL,
		url.QueryEscape(state),
		url.QueryEscape(callbackBase),
		url.QueryEscape(nd.HostName),
		url.QueryEscape(DefaultAPIURL),
	)

	fmt.Println()
	fmt.Println("🌐 Abriendo el navegador para autorizar este nodo...")
	fmt.Printf("   Nodo: %s\n", nd.HostName)
	fmt.Println()
	fmt.Println("   URL:", authURL)
	fmt.Println()
	fmt.Println("💡 Si el navegador no se abre automáticamente, copia y pega la URL.")
	fmt.Println()

	fmt.Println("⏳ Esperando autorización en el navegador (5 min)...")

	select {
	case res := <-resultCh:
		if res.err != nil {
			return fmt.Errorf("autorización fallida: %w", res.err)
		}
		fmt.Println()
		fmt.Println("✅ Autorización recibida. Registrando nodo...")
		return (&EnrollService{}).Link(res.apiURL, res.token, nd.HostName)
	case <-time.After(5 * time.Minute):
		return fmt.Errorf("tiempo de espera agotado (5 min). Vuelve a ejecutar `autohost up`")
	}
}

const successHTML = `<!DOCTYPE html>
<html lang="es">
<head><meta charset="UTF-8"><title>AutoHost CLI — Autorizado</title>
<style>
  body{font-family:system-ui,sans-serif;display:flex;align-items:center;justify-content:center;min-height:100vh;margin:0;background:#0a0a0a;color:#e5e5e5}
  .card{text-align:center;padding:2.5rem 3rem;border:1px solid #27272a;border-radius:1rem;background:#18181b;max-width:440px}
  h2{font-size:1.5rem;margin:0 0 .75rem}
  p{color:#a1a1aa;margin:0}
</style>
</head>
<body>
  <div class="card">
    <h2>✅ Nodo autorizado</h2>
    <p>Puedes cerrar esta ventana y volver a la terminal.</p>
  </div>
</body>
</html>`
