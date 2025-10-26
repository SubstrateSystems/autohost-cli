#!/usr/bin/env bash
set -euo pipefail

REPO="mazapanuwu13/autohost-cli"
BIN_NAME="autohost"

# Detect PATH-friendly BIN_DIR (prioriza /usr/local/bin, si no, ~/.local/bin)
default_bin_dir() {
  if [ -w "/usr/local/bin" ]; then
    echo "/usr/local/bin"
  else
    # Mensaje SOLO a stderr, no a stdout
    echo "🔒 Instalación global requiere permisos de administrador." >&2
    # Devuelve una ruta de usuario (sin sudo ni mkdir aquí)
    echo "$HOME/.local/bin"
  fi
}



PREFIX="${PREFIX:-/usr/local}"
BIN_DIR="${BIN_DIR:-$(default_bin_dir)}"
VERSION="${VERSION:-}"   # opcional: export VERSION=v0.1.0 para fijar un tag

# Detectar OS/ARCH
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"   # linux, darwin
ARCH_RAW="$(uname -m)"                          # x86_64, arm64/aarch64, etc.
case "$ARCH_RAW" in
  x86_64|amd64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "❌ Arquitectura no soportada: $ARCH_RAW"; exit 1 ;;
esac

TMP_DIR="$(mktemp -d)"
cleanup() { rm -rf "$TMP_DIR"; }
trap cleanup EXIT

ua() {
  # User-Agent para evitar bloqueos de la API anónima
  echo "autohost-installer/1.0 (+https://github.com/${REPO})"
}

sha256_cmd() {
  if command -v sha256sum >/dev/null 2>&1; then
    echo "sha256sum"
  elif command -v shasum >/dev/null 2>&1; then
    echo "shasum -a 256"
  else
    echo ""
  fi
}

fetch_latest_tag() {
  curl -fsSL -H "User-Agent: $(ua)" \
    "https://api.github.com/repos/${REPO}/releases/latest" \
    | sed -n 's/.*"tag_name":[[:space:]]*"\([^"]*\)".*/\1/p' | head -n1
}

fetch_release_bin_and_checksums() {
  local tag="$1"
  local asset="${BIN_NAME}-${OS}-${ARCH}"
  local url_bin="https://github.com/${REPO}/releases/download/${tag}/${asset}"
  local url_sum="https://github.com/${REPO}/releases/download/${tag}/checksums_${tag}.txt"

  echo "⬇️  Descargando binario: $url_bin"
  curl -fLsS -H "User-Agent: $(ua)" -o "${TMP_DIR}/${BIN_NAME}" "$url_bin"

  echo "⬇️  Descargando checksums: $url_sum"
  curl -fLsS -H "User-Agent: $(ua)" -o "${TMP_DIR}/checksums.txt" "$url_sum" || {
    echo "⚠️  No se encontró archivo de checksums para ${tag} (continuando sin verificación)."
    return 0
  }

  # Verificar checksum si tenemos herramienta disponible
  local shacmd
  shacmd="$(sha256_cmd)"
  if [ -n "$shacmd" ]; then
    echo "🔐 Verificando checksum..."
    (
      cd "$TMP_DIR"
      # Extrae el checksum esperado para el asset (líneas tipo: <hash>  autohost-linux-amd64)
      expected="$(grep -E "[[:space:]]${asset}$" checksums.txt | awk '{print $1}' || true)"
      if [ -z "${expected}" ]; then
        echo "⚠️  No se encontró checksum para ${asset} en checksums_${tag}.txt (continuando sin verificación)."
      else
        actual="$($shacmd "${BIN_NAME}" | awk '{print $1}')"
        if [ "$expected" != "$actual" ]; then
          echo "❌ Checksum inválido. Esperado: ${expected}  Actual: ${actual}"
          exit 1
        fi
        echo "✅ Checksum verificado."
      fi
    )
  else
    echo "ℹ️  No se encontró 'sha256sum' ni 'shasum'; omitiendo verificación."
  fi
}

install_binary() {
  chmod +x "${TMP_DIR}/${BIN_NAME}"
  echo "🚚 Instalando en ${BIN_DIR}..."

  # Crear BIN_DIR si no existe (intenta sin sudo, luego con sudo)
  if ! mkdir -p "${BIN_DIR}" 2>/dev/null; then
    echo "🔒 Se requieren permisos elevados para instalar en ${BIN_DIR}."
    sudo mkdir -p "${BIN_DIR}"
  fi

  # Mover binario
  if ! mv "${TMP_DIR}/${BIN_NAME}" "${BIN_DIR}/${BIN_NAME}" 2>/dev/null; then
    echo "🔒 Se requieren permisos elevados para instalar en ${BIN_DIR}."
    sudo mv "${TMP_DIR}/${BIN_NAME}" "${BIN_DIR}/${BIN_NAME}"
  fi

  # Sugerir PATH si hace falta
  if ! command -v "${BIN_NAME}" >/dev/null 2>&1; then
    if ! echo "$PATH" | grep -qE "(^|:)${BIN_DIR}(:|$)"; then
      echo "ℹ️  Agrega '${BIN_DIR}' a tu PATH:"
      echo "   echo 'export PATH=\$PATH:${BIN_DIR}' >> ~/.bashrc && source ~/.bashrc"
    fi
  fi

  echo "✅ Instalación completa: $(command -v ${BIN_NAME} || echo "${BIN_DIR}/${BIN_NAME}")"
}


install_from_release() {
  local tag
  if [ -n "${VERSION}" ]; then
    tag="${VERSION}"
  else
    echo "🔎 Obteniendo la última versión (release estable)..."
    tag="$(fetch_latest_tag || true)"
    if [ -z "${tag:-}" ]; then
      echo "ℹ️  No hay releases publicados o la API no respondió."
      return 1
    fi
  fi

  fetch_release_bin_and_checksums "${tag}"
  install_binary
}

install_from_source() {
  echo "🛠  Compilando desde código (go install)..."
  if ! command -v go >/dev/null 2>&1; then
    echo "❌ Necesitas Go instalado para esta ruta (ej: sudo apt-get install -y golang)."
    exit 1
  fi

  local mod="github.com/${REPO}"
  local target="${mod}/cmd/${BIN_NAME}"

  if [ -n "${VERSION}" ]; then
    GO111MODULE=on go install "${target}@${VERSION}"
  else
    GO111MODULE=on go install "${target}@latest"
  fi

  # GOPATH/bin o GOBIN
  local bin_src
  bin_src="$(go env GOBIN || true)"
  if [ -z "$bin_src" ]; then
    bin_src="$(go env GOPATH)/bin"
  fi

  if [ ! -f "${bin_src}/${BIN_NAME}" ]; then
    echo "❌ No se encontró ${BIN_NAME} en ${bin_src}. ¿compiló bien?"
    exit 1
  fi

  echo "🚚 Moviendo a ${BIN_DIR}..."
  mkdir -p "${BIN_DIR}"
  if ! mv "${bin_src}/${BIN_NAME}" "${BIN_DIR}/${BIN_NAME}" 2>/dev/null; then
    echo "🔒 Se requieren permisos elevados para instalar en ${BIN_DIR}."
    sudo mv "${bin_src}/${BIN_NAME}" "${BIN_DIR}/${BIN_NAME}"
  fi

  echo "✅ Instalación completa: $(command -v ${BIN_NAME} || echo "${BIN_DIR}/${BIN_NAME}")"
}

# Intento con release; si falla, compilo desde fuente
if ! install_from_release; then
  install_from_source
fi

echo "👉 Ejecuta: ${BIN_NAME} --help"
