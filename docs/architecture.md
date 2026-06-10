# godex — Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                 godex CLI                                    │
│                          Java & Node.js version manager                      │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                    ┌─────────────────┼─────────────────┐
                    ▼                 ▼                 ▼
              ┌──────────┐    ┌────────────┐    ┌────────────┐
              │  godex   │    │ godex node │    │ godex tui  │
              │   java   │    │            │    │   (gui)    │
              └────┬─────┘    └─────┬──────┘    └─────┬──────┘
                   │                │                 │
      ┌────────────┼───────┐       │                 │
      ▼            ▼       ▼       ▼                 ▼
  ┌───────┐ ┌─────────┐ ┌────┐ ┌───────┐    ┌──────────────┐
  │ list  │ │ current │ │use │ │ list  │    │  Bubble Tea  │
  │       │ │         │ │    │ │current│    │    TUI       │
  └───┬───┘ └────┬────┘ └──┬─┘ │ use   │    │              │
      │          │         │   └───┬───┘    │ ┌──────────┐ │
      │          │         │       │        │ │ Java tab │ │
      │          │         │       │        │ │ Node tab │ │
      │          │         │       │        │ │ Ports tab│ │
      └──────────┼─────────┴───────┘        │ └────┬─────┘ │
                 │                           │      │       │
                 ▼                           │  ┌───┴────┐  │
    ┌────────────────────────┐              │  │Filter  │  │
    │   internal/runtime/    │◄─────────────┤  │Search  │  │
    │                        │              │  │Kill    │  │
    │  install.go            │              │  └────────┘  │
    │    Install struct      │              └──────┬───────┘
    │    listDirs()          │                     │
    │    findInstall()       │                     ▼
    │                        │         ┌──────────────────┐
    │  java.go               │         │ internal/ui/     │
    │    ListJava()           │         │   tui.go         │
    │    FindJava()           │         │                  │
    │    CurrentJava()        │         │ model struct     │
    │                        │         │ Update()  View() │
    │  node.go               │         │ tabs()  panels() │
    │    ListNode()           │         └──────────────────┘
    │    FindNode()           │
    │    CurrentNode()        │
    │                        │
    │  ports.go              │
    │    PortInfo struct      │
    │    ListPorts()          │
    │    KillPort()           │
    │      ├── ss -tulnp ────┤
    │      ├── fuser ────────┤
    │      └── syscall.Kill ─┤
    └────────────────────────┘
```

## Layer Pattern

```
 cmd/           ──  CLI layer (cobra commands, user-facing)
     ▲
     │  calls
     │
 internal/
     │
     ├── ui/      ──  TUI layer (Bubble Tea, lipgloss styles)
     │
     └── runtime/ ──  Business logic (no UI deps)
                      Pure Go, calls external tools:
                      java, node, ss, fuser, kill
```

## TUI State Flow

```
 ┌────────┐  tab   ┌────────┐  tab   ┌────────┐
 │  Java  │◄──────►│  Node  │◄──────►│ Ports  │
 │  tab   │        │  tab   │        │  tab   │
 └────────┘        └────────┘        └────────┘
                                            │
                                   ┌────────┴────────┐
                                   ▼                  ▼
                             / filter          d kill
                             │                  │
                             ▼                  ▼
                        ┌─────────┐    ┌──────────────┐
                        │ typing  │    │ pendingKill  │
                        │ j/k nav │    │ y=confirm    │
                        │ enter   │    │ any=cancel   │
                        │ esc     │    └──────────────┘
                        └─────────┘
```

## Keybindings

| Context     | Key          | Action                            |
|------------ |------------- |---------------------------------- |
| Global      | tab / ←→     | switch tab                        |
| Global      | j/k / ↑↓     | navigate list                     |
| Global      | q / esc      | quit                              |
| Java/Node   | enter        | generate activation command       |
| Java/Node   | c            | copy command to clipboard         |
| Ports       | /            | enter filter mode                 |
| Ports       | enter        | show port detail                  |
| Ports       | c            | copy port number                  |
| Ports       | d → y        | kill process on port (confirm)    |
| Ports       | r            | refresh port list                 |

## File Map

```
.
├── main.go                     # entry point
├── cmd/
│   ├── root.go                 # cobra root (godex)
│   ├── java.go                 # godex java {list,current,use}
│   ├── node.go                 # godex node {list,current,use}
│   └── gui.go                  # godex tui (alias: gui)
├── internal/
│   ├── runtime/
│   │   ├── install.go          # shared: Install struct, listDirs, findInstall
│   │   ├── java.go             # Java detection (sdkman, jenv, /usr/lib/jvm)
│   │   ├── node.go             # Node detection (nvm, fnm)
│   │   └── ports.go            # Port detection (ss) + kill (fuser, syscall)
│   └── ui/
│       └── tui.go              # Bubble Tea TUI (model, tabs, panels, filter)
└── docs/
    └── architecture.md         # this file
```
