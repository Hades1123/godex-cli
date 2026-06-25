# godex

Go developer toolbox — manage Java/Node.js versions, Claude Code presets, and network ports from the command line.

[![Go 1.26+](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://go.dev)
[![Cross-platform](https://img.shields.io/badge/platform-Linux%20%7C%20Windows-blue)]()

## Install

```bash
git clone https://github.com/Hades1123/godex-cli.git
cd godex-cli
./build.sh          # Linux: build + install to ~/.local/bin/
./build.sh windows  # Cross-compile to godex.exe for Windows
```

> Windows setup details → [docs/windows-setup.md](docs/windows-setup.md)

## Commands

```
godex claude          Manage Claude Code settings presets
godex java            Manage Java versions
godex node            Manage Node.js versions
godex ports           List and kill processes by port
```

### `godex claude`

Quickly switch between Claude Code model/API presets (DeepSeek ↔ GLM ↔ …).

```bash
godex claude list                       # List available presets
godex claude current                    # Show active model + API URL
godex claude use glm                    # Switch to the "glm" preset
godex claude api <your-key>             # Change active API key
godex claude template list              # List downloadable templates
godex claude template install deepseek  # Download DeepSeek template from GitHub
```

> Details → [docs/config.md](docs/config.md)

### `godex java`

```bash
godex java list       # List installed Java runtimes
godex java current    # Show the active Java version
godex java use 21     # Print exports for switching (eval with shell wrapper)
```

### `godex node`

```bash
godex node list       # List installed Node.js runtimes
godex node current    # Show the active Node.js version
godex node use 20     # Print exports for switching
```

### `godex ports`

```bash
godex ports list              # List all listening ports
godex ports search nginx      # Search ports by process name
godex ports kill 3000         # Kill process on port 3000 (SIGTERM)
godex ports kill 3000 -f      # Kill with SIGKILL
godex ports kill 3000 -p udp  # Kill UDP port
```

## Shell integration (zsh)

Source `shell/godex.zsh` in your `.zshrc` so `godex java use` / `godex node use` auto-eval:

```bash
echo 'source ~/path/to/godex/shell/godex.zsh' >> ~/.zshrc
source ~/.zshrc
```

## Project structure

```
├── cmd/                      # Cobra commands
│   ├── claude.go             #   godex claude [list|use|current|api|template]
│   ├── claude_helpers.go     #   settings I/O + preset management
│   ├── java.go               #   godex java [list|current|use]
│   ├── node.go               #   godex node [list|current|use]
│   ├── ports.go              #   godex ports [list|search|kill]
│   └── root.go               #   Root command + registration
├── internal/runtime/
│   ├── java.go               #   Java discovery + version switching
│   ├── node.go               #   Node.js discovery + version switching
│   ├── ports.go              #   Port listing via ss(8)
│   ├── ports_unix.go         #   Kill port (Linux/macOS, fuser + syscall)
│   └── ports_windows.go      #   Kill port stub (Windows)
├── templates/                # Preset templates (downloaded via GitHub)
├── shell/godex.zsh           # Zsh wrapper to auto-eval godex * use
├── docs/                     # Documentation
├── build.sh                  # Build script (Linux + Windows cross-compile)
└── main.go                   # Entry point
```

## Dependencies

| Package | Purpose |
|---------|---------|
| [cobra](https://github.com/spf13/cobra) v1.10 | CLI framework |

No external runtime dependencies — Java/Node detection uses syscall + `PATH` scanning, ports use built-in `ss`/`fuser`.

## License

MIT
