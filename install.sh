#!/usr/bin/env bash
#
# Meo Reduce Token (MRT) installer.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/meocode-labs/mrt/main/install.sh | bash
#
# Options (via env vars):
#   MRT_INSTALL_DIR    Where to place the `meo` binary (default: /usr/local/bin)
#   MRT_VERSION        Specific version to install (default: latest release)
#   MRT_SKIP_POSTINSTALL  Set to 1 to skip any postinstall hooks
#
set -euo pipefail

REPO="meocode-labs/mrt"
BINARY="meo"

# -------- helpers -----------------------------------------------------------

say()  { printf '%b\n' "▸ $*"; }
warn() { printf '%b\n' "⚠ $*"; }
die()  { printf '%b\n' "✗ $*" >&2; exit 1; }

need() {
  command -v "$1" >/dev/null 2>&1 || die "required tool '$1' not found in PATH"
}

detect_os() {
  case "$(uname -s)" in
    Darwin) echo "darwin" ;;
    Linux)  echo "linux"  ;;
    MINGW*|MSYS*|CYGWIN*) echo "windows" ;;
    *) die "unsupported OS: $(uname -s)" ;;
  esac
}

detect_arch() {
  case "$(uname -m)" in
    x86_64|amd64)  echo "amd64" ;;
    arm64|aarch64) echo "arm64" ;;
    *) die "unsupported architecture: $(uname -m)" ;;
  esac
}

resolve_latest_version() {
  local url="https://api.github.com/repos/${REPO}/releases/latest"
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$url" \
      | grep '"tag_name"' \
      | head -1 \
      | sed -E 's/.*"v?([^"]+)".*/\1/'
  else
    die "curl is required to resolve the latest version"
  fi
}

# -------- preflight ---------------------------------------------------------

need curl
need uname
need mkdir
need mv
need chmod

OS="$(detect_os)"
ARCH="$(detect_arch)"
EXT=""
[[ "$OS" == "windows" ]] && EXT=".exe"

VERSION="${MRT_VERSION:-}"
if [[ -z "$VERSION" ]]; then
  say "resolving latest ${BINARY} version..."
  VERSION="$(resolve_latest_version)"
  [[ -n "$VERSION" ]] || die "could not determine latest version"
fi

INSTALL_DIR="${MRT_INSTALL_DIR:-/usr/local/bin}"
ASSET="${BINARY}_${OS}_${ARCH}${EXT}"
URL="https://github.com/${REPO}/releases/download/v${VERSION}/${ASSET}"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

say "installing ${BINARY} v${VERSION} (${OS}/${ARCH})"
say "from  ${URL}"
say "to    ${INSTALL_DIR}/${BINARY}"

# -------- install -----------------------------------------------------------

mkdir -p "$TMP"
curl -fsSL --retry 3 -o "${TMP}/${ASSET}" "$URL" \
  || die "download failed. Check that v${VERSION} exists at ${URL}"

chmod +x "${TMP}/${ASSET}"

# Move into place. If INSTALL_DIR is not writable, try with sudo.
TARGET="${INSTALL_DIR}/${BINARY}"
if [[ -w "$INSTALL_DIR" ]]; then
  mv "${TMP}/${ASSET}" "$TARGET"
else
  warn "${INSTALL_DIR} is not writable; attempting with sudo"
  need sudo
  sudo mv "${TMP}/${ASSET}" "$TARGET"
  sudo chmod +x "$TARGET"
fi

# -------- verify ------------------------------------------------------------

if ! command -v "$BINARY" >/dev/null 2>&1; then
  warn "${BINARY} installed to ${TARGET} but is not on PATH"
  warn "add ${INSTALL_DIR} to your PATH, or call ${TARGET} directly"
else
  say "✓ ${BINARY} $(command -v "$BINARY") is ready"
  "$BINARY" --version 2>/dev/null || true
fi
