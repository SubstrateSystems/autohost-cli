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
  command -v multipass >/dev/null || { err "Multipass is not installed."; exit 1; }
  command -v go >/dev/null || { err "Go is not installed."; exit 1; }
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

  if multipass info "${VM_NAME}" >/dev/null 2>&1; then
    log "The VM ${VM_NAME} already exists. Skipping creation."
  else
    log "Creating VM ${VM_NAME}..."
    multipass launch --name "${VM_NAME}" --cpus 2 --memory 2G --disk 10G
  fi
}

ensure_vm_running() {
  if multipass info "${VM_NAME}" >/dev/null 2>&1; then
    multipass start "${VM_NAME}" >/dev/null 2>&1 || true
  else
    err "The VM ${VM_NAME} does not exist. Run: $(basename "$0") run"
    exit 1
  fi
}

push_and_install() {
  log "Transfer binary to the VM..."
  multipass transfer "${LOCAL_BIN}" "${VM_NAME}:/home/ubuntu/autohost"

  log "Moving to /usr/local/bin..."
  multipass exec "${VM_NAME}" -- bash -lc 'sudo mv /home/ubuntu/autohost /usr/local/bin/autohost && sudo chmod +x /usr/local/bin/autohost'

  log "Testing execution inside the VM..."
  # usa bash -lc para asegurar PATH, y tolera exit codes
  multipass exec "${VM_NAME}" -- bash -lc 'autohost --help || true'

  if [[ "${NO_SHELL}" = "1" ]]; then
    log "NO_SHELL=1 -> skipping interactive shell."
    return
  fi

  multipass shell "${VM_NAME}"
}

update_binary() {
  ensure_vm_running
  build_binary

  # Hash previo (opcional, para ver cambios)
  OLD_HASH="$(multipass exec "${VM_NAME}" -- bash -lc 'sha256sum /usr/local/bin/autohost 2>/dev/null | awk "{print \$1}" || true')"
  NEW_HASH="$(sha256sum "${LOCAL_BIN}" | awk '{print $1}')"

  log "transfer new binary..."
  multipass transfer "${LOCAL_BIN}" "${VM_NAME}:/home/ubuntu/autohost.new"

  log "Replacing /usr/local/bin/autohost (atomic installation)..."
  multipass exec "${VM_NAME}" -- bash -lc '
    set -e
    sudo install -m 0755 /home/ubuntu/autohost.new /usr/local/bin/autohost
    rm -f /home/ubuntu/autohost.new
  '

  # Hash posterior (opcional)
  POST_HASH="$(multipass exec "${VM_NAME}" -- bash -lc 'sha256sum /usr/local/bin/autohost | awk "{print \$1}" || true')"

  log " Verification:"
  echo "  Before: ${OLD_HASH:-<no previous binary>}"
  echo "  Now: ${POST_HASH}"
  if [[ -n "${OLD_HASH}" && "${OLD_HASH}" == "${POST_HASH}" ]]; then
    warn "The hash did not change (was the build identical?)."
  else
    log "Update applied."
  fi

  # Quick test
  multipass exec "${VM_NAME}" -- bash -lc 'command -v autohost && autohost --help >/dev/null 2>&1 || true'
}

cmd_run() {
  set -euo pipefail

  ensure_tools
  build_binary
  launch_vm

  # (opcional) verifica que la VM esté corriendo antes de seguir
  if ! multipass list --format csv | grep -q "^${VM_NAME},Running"; then
    log "ERROR: VM ${VM_NAME} no está en Running"; exit 1
  fi

  push_and_install
}

cmd_update() {
  ensure_tools
  update_binary
}

cmd_delete() {
  ensure_tools
  log "Deleting VM ${VM_NAME}..."
  multipass stop "${VM_NAME}" 2>/dev/null || true
  multipass delete "${VM_NAME}" 2>/dev/null || true
  multipass purge || true
  log "Done."
}

usage() {
  cat <<EOF
Usage: $(basename "$0") <run|update|delete>

Commands:
  run     Compiles the binary, creates/starts the VM '${VM_NAME}', installs the binary, and opens a shell.
  update  Compiles and replaces the binary inside the VM (requires the VM to exist).
  delete  Stops, deletes, and purges the VM '${VM_NAME}'.

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
