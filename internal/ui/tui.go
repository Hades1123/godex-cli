package ui

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	devruntime "github.com/hades/cli/internal/runtime"
)

type runtimeKind int

const (
	javaRuntime runtimeKind = iota
	nodeRuntime
	portsRuntime

	javaTabLabel  = " Java"
	nodeTabLabel  = "󰎙 Node.js"
	portsTabLabel = "󰈁 Ports"
)

type model struct {
	active      runtimeKind
	cursor      int
	java        []devruntime.Install
	node        []devruntime.Install
	ports       []devruntime.PortInfo
	portFilter  textinput.Model
	filtering   bool
	currentJava string
	currentNode string
	command     string
	status      string
	err         error
	width       int
	height      int
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("62")).
			Padding(0, 1)

	tabStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(lipgloss.Color("244"))

	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 2).
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("31"))

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39"))

	mutedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("203"))

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("238")).
			Padding(1, 2)

	filterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39"))
)

func Run(_ io.Writer) error {
	m, err := newModel()
	if err != nil {
		return err
	}

	program := tea.NewProgram(m, tea.WithAltScreen())
	_, err = program.Run()
	return err
}

func newModel() (model, error) {
	javaInstalls, err := devruntime.ListJava()
	if err != nil {
		return model{}, err
	}
	nodeInstalls, err := devruntime.ListNode()
	if err != nil {
		return model{}, err
	}

	currentJava, _ := devruntime.CurrentJava()
	currentNode, _ := devruntime.CurrentNode()

	ports, _ := devruntime.ListPorts()

	filter := textinput.New()
	filter.Placeholder = "type port or process..."
	filter.Prompt = "🔍 "
	filter.CharLimit = 32
	filter.Width = 30

	return model{
		active:      javaRuntime,
		java:        javaInstalls,
		node:        nodeInstalls,
		ports:       ports,
		portFilter:  filter,
		currentJava: strings.TrimSpace(currentJava),
		currentNode: strings.TrimSpace(currentNode),
	}, nil
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		// When filtering, delegate most keys to the textinput.
		if m.filtering {
			switch msg.String() {
			case "esc":
				m.filtering = false
				m.portFilter.Reset()
				m.cursor = 0
				m.status = ""
				return m, nil
			case "enter":
				m.filtering = false
				m.cursor = 0
				m.status = ""
				return m, nil
			case "ctrl+c", "q":
				return m, tea.Quit
			default:
				var cmd tea.Cmd
				m.portFilter, cmd = m.portFilter.Update(msg)
				m.cursor = 0
				m.command = ""
				return m, cmd
			}
		}

		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "tab", "right", "l":
			m.switchRuntime()
		case "shift+tab", "left", "h":
			m.switchRuntime()
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < m.listLength()-1 {
				m.cursor++
			}
		case "enter":
			if m.active == portsRuntime {
				m.command = m.portDetail()
				m.status = ""
			} else {
				m.command = m.activationCommand()
				m.status = "Activation command generated."
			}
		case "c":
			if m.active == portsRuntime {
				port := m.selectedPort()
				if port != nil {
					portStr := strconv.Itoa(port.Port)
					if err := clipboard.WriteAll(portStr); err != nil {
						m.status = "Copy failed: " + err.Error()
					} else {
						m.status = fmt.Sprintf("Port %s copied.", portStr)
					}
				}
			} else if m.command != "" {
				if err := clipboard.WriteAll(m.command); err != nil {
					m.status = "Copy failed: " + err.Error()
				} else {
					m.status = "Activation command copied."
				}
			}
		case "/":
			if m.active == portsRuntime {
				m.filtering = true
				m.portFilter.Focus()
				m.status = ""
				return m, textinput.Blink
			}
		case "r":
			if m.active == portsRuntime && !m.filtering {
				ports, err := devruntime.ListPorts()
				if err != nil {
					m.status = "Refresh failed: " + err.Error()
				} else {
					m.ports = ports
					m.status = "Ports refreshed."
				}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return errorStyle.Render(m.err.Error()) + "\n"
	}

	width := m.width
	if width < 72 {
		width = 72
	}
	contentWidth := width - 6
	if contentWidth < 60 {
		contentWidth = 60
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render("devvm"))
	b.WriteString(" ")
	b.WriteString(mutedStyle.Render("Java and Node.js version manager"))
	b.WriteString("\n\n")
	b.WriteString(m.tabs())
	b.WriteString("\n\n")
	b.WriteString(panelStyle.Width(contentWidth).Render(m.runtimePanel(contentWidth - 4)))
	b.WriteString("\n\n")
	b.WriteString(panelStyle.Width(contentWidth).Render(m.commandPanel(contentWidth - 4)))
	b.WriteString("\n\n")
	if m.status != "" {
		b.WriteString(mutedStyle.Render(m.status))
		b.WriteString("\n")
	}
	if m.active == portsRuntime {
		b.WriteString(mutedStyle.Render("tab switch  j/k move  / filter  enter detail  c copy  r refresh  q quit"))
	} else {
		b.WriteString(mutedStyle.Render("tab switch runtime  j/k move  enter generate  c copy  q quit"))
	}
	b.WriteString("\n")
	return b.String()
}

func (m model) tabs() string {
	java := tabStyle.Render(javaTabLabel)
	node := tabStyle.Render(nodeTabLabel)
	ports := tabStyle.Render(portsTabLabel)
	switch m.active {
	case javaRuntime:
		java = activeTabStyle.Render(javaTabLabel)
	case nodeRuntime:
		node = activeTabStyle.Render(nodeTabLabel)
	case portsRuntime:
		ports = activeTabStyle.Render(portsTabLabel)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, java, node, ports)
}

func (m model) runtimePanel(width int) string {
	var b strings.Builder

	switch m.active {
	case javaRuntime:
		b.WriteString(selectedStyle.Render(javaTabLabel))
		b.WriteString("\n")
		b.WriteString(mutedStyle.Render(firstLine(m.currentJava)))
	case nodeRuntime:
		b.WriteString(selectedStyle.Render(nodeTabLabel))
		b.WriteString("\n")
		b.WriteString(mutedStyle.Render(firstLine(m.currentNode)))
	case portsRuntime:
		b.WriteString(selectedStyle.Render(portsTabLabel))
		b.WriteString("\n")
		total := len(m.ports)
		filtered := len(m.filteredPorts())
		if m.filtering && m.portFilter.Value() != "" || total != filtered {
			b.WriteString(mutedStyle.Render(fmt.Sprintf("%d of %d listening ports", filtered, total)))
		} else {
			b.WriteString(mutedStyle.Render(fmt.Sprintf("%d listening ports", total)))
		}
	}
	b.WriteString("\n\n")

	if m.active == portsRuntime {
		return m.portsPanel(&b, width)
	}

	installs := m.installs()
	if len(installs) == 0 {
		b.WriteString(mutedStyle.Render("No installations found."))
		return b.String()
	}

	for i, install := range installs {
		prefix := "  "
		nameStyle := lipgloss.NewStyle()
		if i == m.cursor {
			prefix = "> "
			nameStyle = selectedStyle
		}
		b.WriteString(prefix)
		b.WriteString(nameStyle.Render(install.Name))
		b.WriteString("\n")
		b.WriteString("  ")
		b.WriteString(mutedStyle.Width(width-2).Render(install.Path))
		if i < len(installs)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

func (m model) portsPanel(b *strings.Builder, width int) string {
	// Filter input
	if m.filtering {
		b.WriteString(filterStyle.Render(m.portFilter.View()))
	} else {
		b.WriteString(mutedStyle.Render("Press / to filter by port or process"))
	}
	b.WriteString("\n\n")

	ports := m.filteredPorts()
	if len(ports) == 0 {
		if m.portFilter.Value() != "" {
			b.WriteString(mutedStyle.Render("No ports match the filter."))
		} else {
			b.WriteString(mutedStyle.Render("No listening ports found."))
		}
		return b.String()
	}

	// Header
	b.WriteString(fmt.Sprintf("  %-6s %-7s %-22s %s\n", "Proto", "Port", "Address", "Process"))
	b.WriteString(mutedStyle.Render(strings.Repeat("─", width)))
	b.WriteString("\n")

	for i, p := range ports {
		prefix := "  "
		lineStyle := lipgloss.NewStyle()
		if i == m.cursor {
			prefix = "> "
			lineStyle = selectedStyle
		}

		proto := strings.ToUpper(p.Protocol)
		portStr := strconv.Itoa(p.Port)

		procStr := p.Process
		if procStr == "" {
			procStr = "-"
		}
		if p.PID > 0 {
			procStr = fmt.Sprintf("%s (%d)", procStr, p.PID)
		}

		b.WriteString(prefix)
		b.WriteString(lineStyle.Render(fmt.Sprintf("%-6s %-7s %-22s %s",
			proto, portStr, p.Address, procStr)))
		b.WriteString("\n")
	}

	return b.String()
}

func (m model) commandPanel(width int) string {
	if m.active == portsRuntime {
		if m.filtering {
			return mutedStyle.Render("esc to clear filter  enter to apply")
		}
		if m.command == "" {
			return mutedStyle.Render("Select a port and press enter for details, c to copy the port number, r to refresh.")
		}
		return lipgloss.NewStyle().
			Width(width).
			Foreground(lipgloss.Color("120")).
			Render(m.command)
	}

	if m.command == "" {
		return mutedStyle.Render("Select a version and press enter to generate the shell activation command.")
	}

	return lipgloss.NewStyle().
		Width(width).
		Foreground(lipgloss.Color("120")).
		Render(m.command)
}

func (m *model) switchRuntime() {
	switch m.active {
	case javaRuntime:
		m.active = nodeRuntime
	case nodeRuntime:
		m.active = portsRuntime
	case portsRuntime:
		m.active = javaRuntime
	}
	m.cursor = 0
	m.command = ""
	m.filtering = false
	m.portFilter.Reset()
}

func (m model) installs() []devruntime.Install {
	if m.active == javaRuntime {
		return m.java
	}
	return m.node
}

func (m model) filteredPorts() []devruntime.PortInfo {
	if !m.filtering || m.portFilter.Value() == "" {
		return m.ports
	}
	q := strings.ToLower(m.portFilter.Value())
	var out []devruntime.PortInfo
	for _, p := range m.ports {
		if strings.Contains(strconv.Itoa(p.Port), q) ||
			strings.Contains(strings.ToLower(p.Process), q) ||
			strings.Contains(strings.ToLower(p.Address), q) {
			out = append(out, p)
		}
	}
	return out
}

func (m model) listLength() int {
	switch m.active {
	case javaRuntime:
		return len(m.java)
	case nodeRuntime:
		return len(m.node)
	case portsRuntime:
		return len(m.filteredPorts())
	default:
		return 0
	}
}

func (m model) activationCommand() string {
	installs := m.installs()
	if len(installs) == 0 || m.cursor < 0 || m.cursor >= len(installs) {
		return ""
	}

	path := installs[m.cursor].Path
	binPath := path + string(os.PathSeparator) + "bin"
	if m.active == javaRuntime {
		return fmt.Sprintf("export JAVA_HOME=%q\nexport PATH=%q:$PATH", path, binPath)
	}
	return fmt.Sprintf("export PATH=%q:$PATH", binPath)
}

func (m model) selectedPort() *devruntime.PortInfo {
	ports := m.filteredPorts()
	if m.cursor < 0 || m.cursor >= len(ports) {
		return nil
	}
	return &ports[m.cursor]
}

func (m model) portDetail() string {
	p := m.selectedPort()
	if p == nil {
		return ""
	}
	proc := p.Process
	if proc == "" {
		proc = "unknown"
	}
	if p.PID > 0 {
		proc = fmt.Sprintf("%s (pid %d)", proc, p.PID)
	}
	return fmt.Sprintf("%s :%d — %s (%s)", proc, p.Port, p.Address, strings.ToUpper(p.Protocol))
}

func firstLine(value string) string {
	line, _, _ := strings.Cut(value, "\n")
	if strings.TrimSpace(line) == "" {
		return "Current version is not available."
	}
	return line
}
