#!/usr/bin/env bash
set -e

# ╔══════════════════════════════════════════════════════════════╗
# ║  godex — Go developer toolbox installer                      ║
# ║  Usage: curl -sfL https://... | bash                         ║
# ║         ./install.sh                                         ║
# ╚══════════════════════════════════════════════════════════════╝

CD="$(dirname "$0")"
BINARY="godex"
OPT_INSTALL_COMPLETIONS=true

# ─── Color helpers ──────────────────────────────────────────────
if [ -t 1 ] && command -v tput >/dev/null 2>&1; then
    GREEN="$(tput setaf 2)"
    BOLD="$(tput bold)"
    RESET="$(tput sgr0)"
else
    GREEN=""; BOLD=""; RESET=""
fi
info()  { printf "  ${GREEN}%s${RESET}\n" "$*"; }
step()  { printf "  • %s ... " "$*"; }
ok()    { printf "${GREEN}ok${RESET}\n"; }

# ─── Detect OS & arch ───────────────────────────────────────────
detect_platform() {
    local os arch

    case "$(uname -s)" in
        Linux)  os="linux" ;;
        Darwin) os="darwin" ;;
        MINGW*|MSYS*) os="windows" ;;
        *)      echo "Unsupported OS: $(uname -s)"; exit 1 ;;
    esac

    case "$(uname -m)" in
        x86_64|amd64) arch="amd64" ;;
        aarch64|arm64) arch="arm64" ;;
        *) echo "Unsupported arch: $(uname -m)"; exit 1 ;;
    esac

    echo "${os}/${arch}"
}

# ─── Detect best install prefix ─────────────────────────────────
detect_prefix() {
    # $PREFIX env var wins
    if [ -n "${PREFIX:-}" ]; then
        echo "$PREFIX"
        return
    fi

    # ~/.local/bin is the XDG convention — check if it's on PATH
    if echo "$PATH" | tr ':' '\n' | grep -q "${HOME}/.local/bin"; then
        echo "${HOME}/.local"
        return
    fi

    # Fallback
    echo "${HOME}/.local"
}

# ─── Build binary ───────────────────────────────────────────────
build() {
    local os="${1%/*}"
    local arch="${1#*/}"

    step "Building ${BINARY} for ${os}/${arch}"
    if [ "$os" = "windows" ]; then
        GOOS=windows GOARCH="$arch" go build -o "${BINARY}.exe" .
    else
        go build -o "${BINARY}" .
    fi
    ok
}

# ─── Install binary ─────────────────────────────────────────────
install_binary() {
    local bindir="$1/bin"
    mkdir -p "$bindir"

    step "Installing to ${bindir}/${BINARY}"
    install -m 0755 "${BINARY}" "${bindir}/${BINARY}"
    ok

    # Add to PATH hint
    if ! echo "$PATH" | tr ':' '\n' | grep -qx "$bindir"; then
        local shell_name
        shell_name="$(basename "${SHELL:-bash}")"
        echo ""
        echo "  ⚠️  ${bindir} is not in your PATH."
        echo "     Add this line to your ~/.${shell_name}rc / config.fish:"
        echo ""
        echo "        export PATH=\"\$PATH:${bindir}\""
        echo ""
    fi
}

# ─── Install shell completions (bash, fish) ────────────────────
install_completions() {
    local comp_dir

    # Fish — auto-loaded from ~/.config/fish/completions/<name>.fish
    comp_dir="${HOME}/.config/fish/completions"
    if [ -d "/usr/share/fish/completions" ] && [ -w "/usr/share/fish/completions" ]; then
        comp_dir="/usr/share/fish/completions"
    fi
    mkdir -p "$comp_dir"
    step "Installing fish completions → ${comp_dir}"
    ./"${BINARY}" completion fish > "${comp_dir}/${BINARY}.fish" 2>/dev/null && ok || echo ""

    # Bash — bash-completion loads from BASH_COMPLETION_USER_DIR (default ~/.local/share/bash-completion)
    comp_dir="${HOME}/.local/share/bash-completion/completions"
    if [ -d "/usr/share/bash-completion/completions" ] && [ -w "/usr/share/bash-completion/completions" ]; then
        comp_dir="/usr/share/bash-completion/completions"
    fi
    mkdir -p "$comp_dir"
    step "Installing bash completions → ${comp_dir}"
    ./"${BINARY}" completion bash > "${comp_dir}/${BINARY}" 2>/dev/null && ok || echo ""

    info ""
    info " ✅ Completions installed for fish and bash"
    info "    (bash needs the 'bash-completion' package)"
}

# ─── Cleanup ────────────────────────────────────────────────────
cleanup() {
    rm -f "${BINARY}" "${BINARY}.exe"
}

# ─── Main ───────────────────────────────────────────────────────
main() {
    echo ""
    info "${BOLD}┌──────────────────────────────────────────┐${RESET}"
    info "${BOLD}│  godex — Go developer toolbox installer  │${RESET}"
    info "${BOLD}└──────────────────────────────────────────┘${RESET}"
    echo ""

    # Check prerequisites
    if ! command -v go &>/dev/null; then
        echo "❌ Go is not installed."
        echo "   Install it first:"
        echo "     Fedora: sudo dnf install golang"
        echo "     Ubuntu: sudo apt install golang-go"
        echo "     macOS:  brew install go"
        echo "     Or download from https://go.dev/dl/"
        exit 1
    fi

    local platform
    platform="$(detect_platform)"
    local prefix
    prefix="$(detect_prefix)"

    build "$platform"
    install_binary "$prefix"

    if [ "$OPT_INSTALL_COMPLETIONS" = true ]; then
        echo ""
        install_completions "$prefix"
    fi

    cleanup

    echo ""
    info " 🎉 ${BOLD}godex installed!${RESET}"
    info ""
    info "     godex --help"
    info "     godex ports list"
    info "     godex claude list"
    echo ""
}

main "$@"
