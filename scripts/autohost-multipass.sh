#!/usr/bin/env bash
set -euo pipefail

VM_NAME="autohost-test"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
BIN_DIR="${REPO_ROOT}/dist"
LOCAL_BIN="${BIN_DIR}/autohost"
MAIN_PKG="${REPO_ROOT}/cmd/autohost-cli"

log()  { echo -e "\033[1;34m[autohost]\033[0m $*"; }
warn() { echo -e "\033[1;33m[warn]\033[0m $*"; }
err()  { echo -e "\033[1;31m[err]\033[0m  $*" >&2; }

ensure_tools() {
  command -v multipass >/dev/null || { err "Multipass no está instalado."; exit 1; }
  command -v go >/dev/null || { err "Go no está instalado."; exit 1; }
}

build_binary() {
  log "Compilando autohost-cli..."
  mkdir -p "${BIN_DIR}"
  # Si quieres binario más chico: añade -ldflags="-s -w"
  GOOS=linux GOARCH="$(go env GOARCH)" go build -o "${LOCAL_BIN}" "${MAIN_PKG}"
  chmod +x "${LOCAL_BIN}"
  log "Binario generado en ${LOCAL_BIN}"
}

launch_vm() {
  if multipass info "${VM_NAME}" >/dev/null 2>&1; then
    log "La VM ${VM_NAME} ya existe. Omitiendo creación."
  else
    log "Creando VM ${VM_NAME}..."
    multipass launch --name "${VM_NAME}" --cpus 2 --memory 2G --disk 10G
  fi
}

ensure_vm_running() {
  if multipass info "${VM_NAME}" >/dev/null 2>&1; then
    multipass start "${VM_NAME}" >/dev/null 2>&1 || true
  else
    err "La VM ${VM_NAME} no existe. Ejecuta: $(basename "$0") run"
    exit 1
  fi
}

push_and_install() {
  log "Transfiriendo binario a la VM..."
  multipass transfer "${LOCAL_BIN}" "${VM_NAME}:/home/ubuntu/autohost"
  log "Moviendo a /usr/local/bin..."
  multipass exec "${VM_NAME}" -- bash -lc 'sudo mv /home/ubuntu/autohost /usr/local/bin/autohost && sudo chmod +x /usr/local/bin/autohost'
  log "Probando ejecución dentro de la VM..."
  multipass exec "${VM_NAME}" -- autohost --help || true
  multipass shell "${VM_NAME}"
}

update_binary() {
  ensure_vm_running
  build_binary

  # Hash previo (opcional, para ver cambios)
  OLD_HASH="$(multipass exec "${VM_NAME}" -- bash -lc 'sha256sum /usr/local/bin/autohost 2>/dev/null | awk "{print \$1}" || true')"
  NEW_HASH="$(sha256sum "${LOCAL_BIN}" | awk '{print $1}')"

  log "Transfiriendo nuevo binario..."
  multipass transfer "${LOCAL_BIN}" "${VM_NAME}:/home/ubuntu/autohost.new"

  log "Reemplazando /usr/local/bin/autohost (instalación atómica)..."
  multipass exec "${VM_NAME}" -- bash -lc '
    set -e
    sudo install -m 0755 /home/ubuntu/autohost.new /usr/local/bin/autohost
    rm -f /home/ubuntu/autohost.new
  '

  # Hash posterior (opcional)
  POST_HASH="$(multipass exec "${VM_NAME}" -- bash -lc 'sha256sum /usr/local/bin/autohost | awk "{print \$1}" || true')"

  log "Verificación:"
  echo "  Antes: ${OLD_HASH:-<sin binario previo>}"
  echo "  Ahora: ${POST_HASH}"
  if [[ -n "${OLD_HASH}" && "${OLD_HASH}" == "${POST_HASH}" ]]; then
    warn "El hash no cambió (¿compilación igual?)."
  else
    log "Actualización aplicada."
  fi

  # Prueba rápida
  multipass exec "${VM_NAME}" -- bash -lc 'command -v autohost && autohost --help >/dev/null 2>&1 || true'
}

cmd_run() {
  ensure_tools
  build_binary
  launch_vm
  push_and_install
}

cmd_update() {
  ensure_tools
  update_binary
}

cmd_delete() {
  ensure_tools
  log "Eliminando VM ${VM_NAME}..."
  multipass stop "${VM_NAME}" 2>/dev/null || true
  multipass delete "${VM_NAME}" 2>/dev/null || true
  multipass purge || true
  log "Listo."
}

usage() {
  cat <<EOF
Uso: $(basename "$0") <run|update|delete>

Comandos:
  run     Compila el binario, crea/arranca la VM '${VM_NAME}', instala el binario y abre una shell.
  update  Compila y reemplaza el binario dentro de la VM (requiere que la VM exista).
  delete  Detiene, elimina y purga la VM '${VM_NAME}'.

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
