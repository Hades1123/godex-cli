#!/usr/bin/env bash
# ─── godex build script ─────────────────────────────────────
# Delegates to Makefile for consistency.
# Usage: ./build.sh [linux|windows]
#   linux   — build + install to ~/.local/bin/ (default)
#   windows — cross-compile to godex.exe
set -e

cd "$(dirname "$0")"

case "${1:-linux}" in
    linux|lin)
        echo "✨ Building godex (Linux)…"
        make build
        make install
        make vet
        hash -r 2>/dev/null || true
        echo "✅ godex installed to ~/.local/bin/godex"
        echo ""
        echo "   Run: godex --help"
        ;;
    windows|win)
        echo "✨ Building godex.exe (Windows)…"
        make build/windows
        echo "✅ godex.exe created ($(du -h godex.exe | cut -f1))"
        echo ""
        echo "🪟  Copy godex.exe to your Windows machine:"
        echo "   C:\Users\<you>\bin\godex.exe --help"
        ;;
    *)
        echo "Usage: $0 [linux|windows]"
        echo "  linux   — build + install (default)"
        echo "  windows — cross-compile to godex.exe"
        exit 1
        ;;
esac
