#!/usr/bin/env bash
set -euo pipefail

NO_SHELL="${NO_SHELL:-0}"

VM_NAME="autohost-test"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
BIN_DIR="${REPO_ROOT}/dist"
LOCAL_BIN="${BIN_DIR}/autohost"
MAIN_PKG="${REPO_ROOT}/"

log()  { echo -e "\033[1;34m[autohost]\033[0m $*"; }
warn() { echo -e "\033[1;33m[warn]\033[0m $*"; }
err()  { echo -e "\033[1;31m[err]\033[0m  $*" >&2; }

ensure_tools() {
  command -v incus >/dev/null || { err "Incus is not installed."; exit 1; }
  command -v go >/dev/null || { err "Go is not installed."; exit 1; }
  
  # Verificar permisos con incus
  if ! incus list >/dev/null 2>&1; then
    err "No tienes permisos para usar Incus."
    echo ""
    echo "Solución:"
    echo "  1. Añádete al grupo incus-admin:"
    echo "     sudo usermod -aG incus-admin \$USER"
    echo ""
    echo "  2. Cierra sesión y vuelve a entrar (o reinicia)"
    echo ""
    echo "  3. Verifica con: groups | grep incus-admin"
    echo ""
    exit 1
  fi
}

build_binary() {
  log "Building autohost-cli..."
  mkdir -p "${BIN_DIR}"
  # Si quieres binario más chico: añade -ldflags="-s -w"
  GOOS=linux GOARCH="$(go env GOARCH)" go build -o "${LOCAL_BIN}" "${MAIN_PKG}"
  chmod +x "${LOCAL_BIN}"
  log "Binary generated at ${LOCAL_BIN}"
}

launch_vm() {
  if incus info "${VM_NAME}" >/dev/null 2>&1; then
    log "The instance ${VM_NAME} already exists. Skipping creation."
  else
    log "Creating instance ${VM_NAME}..."
    # Usa images:ubuntu/24.04 o images:ubuntu/22.04 según prefieras
    incus launch images:ubuntu/24.04 "${VM_NAME}" \
      --config limits.cpu=2 \
      --config limits.memory=2GiB
    
    # Esperar a que la instancia esté lista
    log "Waiting for instance to be ready..."
    for i in {1..30}; do
      if incus exec "${VM_NAME}" -- systemctl is-system-running --wait 2>/dev/null | grep -qE "running|degraded"; then
        break
      fi
      sleep 2
    done
  fi
}

ensure_vm_running() {
  if incus info "${VM_NAME}" >/dev/null 2>&1; then
    incus start "${VM_NAME}" >/dev/null 2>&1 || true
  else
    err "The instance ${VM_NAME} does not exist. Run: $(basename "$0") run"
    exit 1
  fi
}

push_and_install() {
  log "Transfer binary to the instance..."
  incus file push "${LOCAL_BIN}" "${VM_NAME}/home/ubuntu/autohost"

  log "Moving to /usr/local/bin..."
  incus exec "${VM_NAME}" -- bash -lc 'sudo mv /home/ubuntu/autohost /usr/local/bin/autohost && sudo chmod +x /usr/local/bin/autohost'

  log "Testing execution inside the instance..."
  # usa bash -lc para asegurar PATH, y tolera exit codes
  incus exec "${VM_NAME}" -- bash -lc 'autohost --help || true'

  if [[ "${NO_SHELL}" = "1" ]]; then
    log "NO_SHELL=1 -> skipping interactive shell."
    return
  fi

  incus exec "${VM_NAME}" -- bash -l
}

update_binary() {
  ensure_vm_running
  build_binary

  # Hash previo (opcional, para ver cambios)
  OLD_HASH="$(incus exec "${VM_NAME}" -- bash -lc 'sha256sum /usr/local/bin/autohost 2>/dev/null | awk "{print \$1}" || true')"
  NEW_HASH="$(sha256sum "${LOCAL_BIN}" | awk '{print $1}')"

  log "transfer new binary..."
  incus file push "${LOCAL_BIN}" "${VM_NAME}/home/ubuntu/autohost.new"

  log "Replacing /usr/local/bin/autohost (atomic installation)..."
  incus exec "${VM_NAME}" -- bash -lc '
    set -e
    sudo install -m 0755 /home/ubuntu/autohost.new /usr/local/bin/autohost
    rm -f /home/ubuntu/autohost.new
  '

  # Hash posterior (opcional)
  POST_HASH="$(incus exec "${VM_NAME}" -- bash -lc 'sha256sum /usr/local/bin/autohost | awk "{print \$1}" || true')"

  log " Verification:"
  echo "  Before: ${OLD_HASH:-<no previous binary>}"
  echo "  Now: ${POST_HASH}"
  if [[ -n "${OLD_HASH}" && "${OLD_HASH}" == "${POST_HASH}" ]]; then
    warn "The hash did not change (was the build identical?)."
  else
    log "Update applied."
  fi

  # Quick test
  incus exec "${VM_NAME}" -- bash -lc 'command -v autohost && autohost --help >/dev/null 2>&1 || true'
}

cmd_run() {
  set -euo pipefail

  ensure_tools
  build_binary
  launch_vm

  # (opcional) verifica que la instancia esté corriendo antes de seguir
  if ! incus list --format csv | grep -q "^${VM_NAME},"; then
    log "ERROR: Instance ${VM_NAME} not found"; exit 1
  fi

  push_and_install
}

cmd_update() {
  ensure_tools
  update_binary
}

cmd_delete() {
  ensure_tools
  log "Deleting instance ${VM_NAME}..."
  incus stop "${VM_NAME}" 2>/dev/null || true
  incus delete "${VM_NAME}" 2>/dev/null || true
  log "Done."
}

usage() {
  cat <<EOF
Usage: $(basename "$0") <run|update|delete>

Commands:
  run     Compiles the binary, creates/starts the instance '${VM_NAME}', installs the binary, and opens a shell.
  update  Compiles and replaces the binary inside the instance (requires the instance to exist).
  delete  Stops and deletes the instance '${VM_NAME}'.

EOF
}

main() {
  case "${1:-}" in
    run)    cmd_run ;;
    update) cmd_update ;;
    delete) cmd_delete ;;
    *)      usage; exit 1 ;;
  esac
}

main "$@"
