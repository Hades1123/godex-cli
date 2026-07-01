# godex

Go developer toolbox — manage Claude Code presets and network ports from the command line.

[![Go 1.26+](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://go.dev)
[![Cross-platform](https://img.shields.io/badge/platform-Linux%20%7C%20Windows-blue)]()

## Install

### Fedora / Linux (from source)

```bash
# Prerequisites
sudo dnf install golang         # Fedora — or install from go.dev/dl/

# Build & install
git clone https://github.com/Hades1123/godex-cli.git
cd godex-cli

# Option A: one-command build + install
./build.sh

# Option B: Makefile (more control)
make build          # just build the binary
make install        # build + copy to ~/.local/bin/
make install PREFIX=/usr  # system-wide install
```

Verify it works:

```bash
godex --help
godex ports list
```

### Windows (cross-compile)

```bash
./build.sh windows   # produces godex.exe
```

> Windows setup details → [docs/windows-setup.md](docs/windows-setup.md)

### Advanced: custom install script

The `install.sh` script auto-detects your platform and installs the binary plus shell completions:

```bash
./install.sh
```

### Shell completions (fish + bash)

`make install` installs completions automatically. To do it manually:

```bash
make completions               # generate completion scripts
make install-completions       # install to fish + bash completion dirs
```

- **fish** → `~/.config/fish/completions/godex.fish` (auto-loaded)
- **bash** → `~/.local/share/bash-completion/completions/godex` (needs `bash-completion` package)

Now `godex <Tab>` lists commands instead of files, and `godex cl<Tab>` completes `claude`.

## Commands

```
godex claude          Manage Claude Code settings presets
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

### `godex ports`

```bash
godex ports list              # List all listening ports
godex ports search nginx      # Search ports by process name
godex ports kill 3000         # Kill process on port 3000 (SIGTERM)
godex ports kill 3000 -f      # Kill with SIGKILL
godex ports kill 3000 -p udp  # Kill UDP port
```

## Project structure

```
├── cmd/                      # Cobra commands
│   ├── claude.go             #   godex claude [list|use|current|api|template]
│   ├── claude_helpers.go     #   settings I/O + preset management
│   ├── ports.go              #   godex ports [list|search|kill]
│   └── root.go               #   Root command + registration
├── internal/runtime/
│   ├── ports.go              #   Port listing via ss(8)
│   ├── ports_unix.go         #   Kill port (Linux/macOS, fuser + syscall)
│   └── ports_windows.go      #   Kill port stub (Windows)
├── completions/              # Generated shell completions
│   ├── godex.bash            #   bash completions (needs bash-completion)
│   └── godex.fish            #   fish completions (auto-loaded, press Tab!)
├── templates/                # Preset templates (downloaded via GitHub)
├── docs/                     # Documentation
├── build.sh                  # Build script (Linux + Windows cross-compile)
├── install.sh                # Auto-detecting installer (platform, completions)
├── Makefile                  # Build automation (build, install, test, clean)
└── main.go                   # Entry point
```

## Dependencies

| Package | Purpose |
|---------|---------|
| [cobra](https://github.com/spf13/cobra) v1.10 | CLI framework |

**Build-time only.** The binary has zero runtime dependencies — ports use built-in `ss`/`fuser`.

## License

MIT
