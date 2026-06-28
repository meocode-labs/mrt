#!/usr/bin/env bash
#
# Meo Reduce Token (MRT) installer.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/meocode-labs/mrt/main/install.sh | bash
#
# Options (via env vars):
#   MRT_INSTALL_DIR    Where to place the `mrt` binary (default: /usr/local/bin)
#   MRT_VERSION        Specific version to install (default: latest release)
#   MRT_SKIP_POSTINSTALL  Set to 1 to skip any postinstall hooks
#   MRT_AUTO_MIGRATE   Set to 1 to auto-remove old `meo` binary if found
#
set -euo pipefail

REPO="${MRT_REPO:-meocode-labs/mrt}"
# Source asset prefix on GitHub releases: the auto-release workflow
# produces binaries named mrt_<os>_<arch> (e.g. mrt_darwin_arm64).
# The installed command is also `mrt` (renamed from `meo` in v1.3.0).
ASSET_PREFIX="mrt"
BINARY="mrt"
BASE_URL="${MRT_BASE_URL:-https://github.com}"
API_URL="${MRT_API_URL:-https://api.github.com}"

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
  local url="${API_URL}/repos/${REPO}/releases/latest"
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
# Normalize: accept "1.2.0" or "v1.2.0" — strip any leading "v" so URL
# construction below always produces exactly one "v" prefix.
VERSION="${VERSION#v}"

INSTALL_DIR="${MRT_INSTALL_DIR:-/usr/local/bin}"
ASSET="${ASSET_PREFIX}_${OS}_${ARCH}${EXT}"
URL="${BASE_URL}/${REPO}/releases/download/v${VERSION}/${ASSET}"
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

# -------- migration notice (v1.3.0: `meo` renamed to `mrt`) -----------------

OLD_BINARY="${INSTALL_DIR}/meo"
if [[ -f "$OLD_BINARY" && "$OLD_BINARY" != "$TARGET" ]]; then
  warn "v1.3.0 renamed the command from \`meo\` to \`${BINARY}\`."
  warn "  old binary found: ${OLD_BINARY}"
  if [[ "${MRT_AUTO_MIGRATE:-0}" == "1" ]]; then
    if rm -f "$OLD_BINARY" 2>/dev/null || sudo rm -f "$OLD_BINARY" 2>/dev/null; then
      say "✓ removed old \`meo\` binary"
    else
      warn "  could not remove ${OLD_BINARY}; remove it manually"
    fi
  else
    warn "  remove it with: rm ${OLD_BINARY}"
    warn "  or re-run with MRT_AUTO_MIGRATE=1 to remove automatically"
  fi
fi

# -------- verify ------------------------------------------------------------

# Always report the path we just installed to. If a different mrt is
# already on PATH (e.g. from npm), the user can decide whether to
# adjust PATH or uninstall the old copy.
say "✓ ${BINARY} installed to ${TARGET}"
"${TARGET}" --version 2>/dev/null || true

if ! command -v "$BINARY" >/dev/null 2>&1 \
   || [[ "$(command -v "$BINARY")" != "$TARGET" ]]; then
  warn "${TARGET} is not the ${BINARY} on PATH"
  warn "  on PATH: $(command -v "$BINARY" 2>/dev/null || echo '(none)')"
  warn "  add ${INSTALL_DIR} to PATH, or call ${TARGET} directly"
fi
