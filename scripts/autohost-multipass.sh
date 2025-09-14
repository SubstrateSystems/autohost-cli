#!/usr/bin/env bash
set -euo pipefail

VM_NAME="autohost-test"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
BIN_DIR="${REPO_ROOT}/dist"
LOCAL_BIN="${BIN_DIR}/autohost"
MAIN_PKG="${REPO_ROOT}/cmd/autohost-cli"

log()  { echo -e "\033[1;34m[autohost]\033[0m $*"; }
err()  { echo -e "\033[1;31m[err]\033[0m  $*" >&2; }

ensure_tools() {
  command -v multipass >/dev/null || { err "Multipass no está instalado."; exit 1; }
  command -v go >/dev/null || { err "Go no está instalado."; exit 1; }
}

build_binary() {
  log "Compilando autohost-cli..."
  mkdir -p "${BIN_DIR}"
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

push_and_install() {
  log "Transfiriendo binario a la VM..."
  multipass transfer "${LOCAL_BIN}" "${VM_NAME}:/home/ubuntu/autohost"
  log "Moviendo a /usr/local/bin..."
  multipass exec "${VM_NAME}" -- bash -lc 'sudo mv /home/ubuntu/autohost /usr/local/bin/autohost && sudo chmod +x /usr/local/bin/autohost'
  log "Probando ejecución dentro de la VM..."
  multipass exec "${VM_NAME}" -- autohost --help || true
  multipass shell "${VM_NAME}"
}

cmd_run() {
  ensure_tools
  build_binary
  launch_vm
  push_and_install
}

cmd_delete() {
  ensure_tools
  log "Eliminando VM ${VM_NAME}..."
  multipass stop "${VM_NAME}" 2>/dev/null || true
  multipass delete "${VM_NAME}" 2>/dev/null || true
  multipass purge || true
}

usage() {
  echo "Uso: $0 <run|delete>"
}

main() {
  case "${1:-}" in
    run) cmd_run ;;
    delete) cmd_delete ;;
    *) usage; exit 1 ;;
  esac
}

main "$@"
