#!/usr/bin/env bash
set -e

cd "$(dirname "$0")"

case "${1:-linux}" in
    linux|lin)
        echo "✨ Building godex (Linux)…"
        go build -o "$HOME/.local/bin/godex" .

        echo "✨ Running vet…"
        go vet ./...

        hash -r
        echo "✅ godex installed to ~/.local/bin/godex"
        ;;
    windows|win)
        echo "✨ Building godex.exe (Windows)…"
        GOOS=windows GOARCH=amd64 go build -o godex.exe .

        echo "✅ godex.exe created ($(du -h godex.exe | cut -f1))"
        echo ""
        echo "🪟  Copy godex.exe sang Windows và chạy:"
        echo "   C:\Users\&lt;bạn&gt;\bin\godex.exe config list"
        ;;
    *)
        echo "Usage: $0 [linux|windows]"
        echo "  linux   — build + cài vào PATH (mặc định)"
        echo "  windows — cross-compile ra godex.exe"
        exit 1
        ;;
esac
