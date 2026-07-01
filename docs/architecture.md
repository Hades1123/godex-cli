# godex — Architecture

```
┌──────────────────────────────────────────────────────────┐
│                        godex CLI                         │
│        Claude Code presets & network port manager        │
└──────────────────────────────────────────────────────────┘
                              │
                ┌─────────────┴─────────────┐
                ▼                           ▼
          ┌──────────┐                ┌──────────┐
          │  godex   │                │  godex   │
          │  claude  │                │  ports   │
          └────┬─────┘                └────┬─────┘
               │                           │
      ┌────────┼────────┐         ┌────────┼─────────┐
      ▼        ▼        ▼         ▼        ▼          ▼
   ┌─────┐ ┌───────┐ ┌──────┐ ┌──────┐ ┌───────┐ ┌──────┐
   │list │ │current│ │ use  │ │ list │ │search │ │ kill │
   │current│ │ api  │ │ api  │ │      │ │       │ │      │
   │ api │ │template│ │      │ │      │ │       │ │      │
   └──┬──┘ └───┬───┘ └──┬───┘ └───┬──┘ └───┬───┘ └──┬───┘
      │        │        │         │        │        │
      ▼        ▼        ▼         ▼        ▼        ▼
   ┌──────────────────────┐   ┌─────────────────────────────┐
   │ internal/runtime/    │   │ internal/runtime/           │
   │  (claude preset I/O) │   │  ports.go  — ss -tulnp      │
   │                      │   │  ports_unix.go — fuser +    │
   │  settings.json read/ │   │     syscall.Kill            │
   │  write + presets     │   │  ports_windows.go — stub    │
   └──────────────────────┘   └─────────────────────────────┘
```

## Layer Pattern

```
 cmd/           ──  CLI layer (cobra commands, user-facing)
     ▲
     │  calls
     │
 internal/runtime/  ──  Business logic (no CLI deps)
                      Pure Go; calls external tools:
                      ss, fuser, kill
```

The `claude` command manages the Claude Code `settings.json` preset directly (read/write
`~/.claude/settings.json` and the preset store), plus downloads preset templates from GitHub.

The `ports` command shells out to `ss(8)` to list listeners, and to `fuser` + `syscall.Kill`
(on Unix) to terminate a process bound to a port.

## File Map

```
.
├── main.go                     # entry point
├── cmd/
│   ├── root.go                 # cobra root (godex)
│   ├── claude.go               # godex claude {list,current,use,api,template}
│   ├── claude_helpers.go       # settings.json I/O + preset management
│   └── ports.go                # godex ports {list,search,kill}
├── internal/runtime/
│   ├── ports.go                # Port listing via ss -tulnp
│   ├── ports_unix.go           # Kill port (Linux/macOS): fuser + syscall.Kill
│   └── ports_windows.go        # Kill port stub (Windows)
├── templates/                  # Preset templates (downloaded via GitHub)
├── completions/                # Generated shell completions (bash, fish)
└── docs/
    └── architecture.md         # this file
```
