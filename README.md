![autohost-cli](./autohost-cli.png)

# AutoHost CLI

CLI para automatizar el self-hosting en Linux: instala aplicaciones, gestiona el agente, y conecta nodos con la [autohost-cloud-api](../autohost-cloud-api).

## Stack

- **Go** 1.23.0
- **Cobra** (comandos CLI)
- **SQLite** (estado local via modernc/sqlite)
- **Docker** (backend de apps)
- **Caddy / CoreDNS / Tailscale / Cloudflare Tunnel** (exposición de servicios)

## Requisitos previos

- Linux (amd64 o arm64)
- Go 1.23+ (solo para compilar desde fuente)
- Docker instalado y en ejecución
- `sudo` o permisos de root para operaciones de sistema

---

## Instalación

### Desde GitHub Releases (recomendado)

```bash
curl -fsSL https://raw.githubusercontent.com/mazapanuwu13/autohost-cli/main/scripts/install.sh | bash
```

The script auto-detects your OS and architecture, downloads the latest binary from GitHub Releases, verifies the checksum, and installs it.

If the binary lands in `~/.local/bin`, add it to your PATH:

```bash
export PATH="$PATH:$HOME/.local/bin"
source ~/.bashrc
```

### Specific version

```bash
VERSION=v0.4.0 curl -fsSL https://raw.githubusercontent.com/mazapanuwu13/autohost-cli/main/scripts/install.sh | bash
```

### Verify installation

```bash
autohost --version
```

---

## 🛠 Basic Usage

### Example workflow

```bash
export PATH="$PATH:$HOME/.local/bin"   # añadir también a ~/.bashrc
```

### Compilar desde fuente

```bash
git clone https://github.com/mazapanuwu13/autohost-cli.git
cd autohost-cli
go build -o autohost main.go
```

---

## Comandos disponibles

### `autohost agent` — Gestión del agente

```bash
autohost agent install   # Descarga e instala el autohost-agent como servicio systemd
```

Tras la instalación, edita `/etc/autohost/config.yaml` y arranca el servicio:

```bash
sudo nano /etc/autohost/config.yaml
sudo systemctl enable autohost-agent
sudo systemctl start autohost-agent
```

### `autohost enroll` — Enrollment de nodos

```bash
autohost enroll link --api <URL_API> --token <ENROLLMENT_TOKEN> [--name <nombre>]
```

Enlaza el nodo actual con la cloud API. El comando:
1. Auto-detecta hostname, IP local, OS y arquitectura
2. Llama a `POST /v1/enrollments/enroll` con el token de enrollment
3. Guarda el `agent_token` resultante en `/etc/autohost/config.yaml`

Ejemplo:

```bash
autohost enroll link \
  --api http://192.168.1.10:8080 \
  --token autohost-enroll_xxxx \
  --name mi-servidor
```

---

## Comandos en desarrollo

Los siguientes comandos están implementados pero no están activos en el binario actual (comentados en `root.go`). Se activarán progresivamente:

| Comando | Descripción |
|---------|-------------|
| `autohost setup` | Instala Docker, Caddy y configura el servidor |
| `autohost app install [nombre]` | Instala una app del catálogo |
| `autohost app ls` | Lista apps instaladas |
| `autohost app start/stop/status/remove [nombre]` | Gestiona el ciclo de vida de una app |
| `autohost expose app` | Expone una app via Tailscale (privado) o Cloudflare (público) |

### Apps en el catálogo

| App | Puerto por defecto |
|-----|-------------------|
| Nextcloud | 8081 |
| BookStack | 6875 |
| Joplin | — |
| Excalidraw | — |
| MySQL | 3306 |
| PostgreSQL | 5432 |
| Redis | 6379 |

---

## Estructura del proyecto

```
.
├── cmd/autohost-cli/
│   ├── root.go           # Punto de entrada, registro de comandos
│   ├── agent/            # autohost agent install
│   ├── app/              # autohost app (install, ls, start, stop, status, remove)
│   ├── expose/           # autohost expose app
│   ├── install/          # autohost install
│   └── setup/            # autohost setup
├── internal/
│   ├── adapters/         # Integraciones externas
│   │   ├── caddy/
│   │   ├── cloudflare/
│   │   ├── coreDNS/
│   │   ├── docker/
│   │   ├── tailscale/
│   │   ├── terraform/
│   │   └── storage/      # Repositorios SQLite
│   ├── app/              # Servicios de aplicación (AppService, ExposeService, SetupService)
│   ├── domain/           # Modelos e interfaces de dominio
│   ├── platform/         # DI, configuración, filesystem
│   └── plugins/enroll/   # Plugin de enrollment de nodos
├── assets/docker/        # Plantillas docker-compose embebidas por app
├── utils/                # Helpers (paths, env, secrets, URLs…)
├── scripts/
│   ├── install.sh        # Instalador del CLI
│   ├── autohost-multipass.sh
│   └── autohost-incus.sh
└── Makefile
```

---

## Comandos Make (entornos de prueba)

```bash
# Multipass
make vm-run        # Crea VM de prueba (autohost-test) con Multipass
make vm-update     # Actualiza el binario en la VM
make vm-delete     # Elimina la VM

# Incus
make incus-run     # Crea instancia de prueba con Incus
make incus-update  # Actualiza el binario en la instancia
make incus-delete  # Elimina la instancia
```

---

## Convenciones de código

- Go 1.23.0, `gofmt` / `goimports` obligatorio
- Errores envueltos con contexto: `fmt.Errorf("contexto: %w", err)`
- Comandos idempotentes: crear si no existe, actualizar si cambió, no duplicar
- Sin secretos en logs ni mensajes de error
- Detectar root con `os.Geteuid() == 0`; usar `sudo` como fallback si está disponible

## Licencia

MIT
