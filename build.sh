#!/usr/bin/env bash
set -e

cd "$(dirname "$0")"

echo "✨ Building godex…"
go build -o "$HOME/.local/bin/godex" .

echo "✨ Running vet…"
go vet ./...

hash -r

echo "✅ godex installed to ~/.local/bin/godex"
