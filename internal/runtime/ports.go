package runtime

import (
	"bufio"
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

// PortInfo holds information about a listening port.
type PortInfo struct {
	Protocol string // "tcp", "tcp6", "udp", "udp6"
	Port     int
	Address  string
	Process  string // process name, empty if unavailable
	PID      int    // 0 if unavailable
}

// ListPorts returns all listening TCP and UDP ports.
func ListPorts() ([]PortInfo, error) {
	// Try with process info first, fall back without.
	out, err := exec.Command("ss", "-tulnp").CombinedOutput()
	if err != nil {
		// Fallback: try without process info.
		out, err = exec.Command("ss", "-tuln").CombinedOutput()
		if err != nil {
			return nil, fmt.Errorf("ss command failed: %w", err)
		}
	}
	return parseSS(string(out)), nil
}

func parseSS(output string) []PortInfo {
	var ports []PortInfo
	scanner := bufio.NewScanner(strings.NewReader(output))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "Netid") || strings.HasPrefix(line, "State") {
			continue
		}

		// Columns: Netid State Recv-Q Send-Q Local Address:Port Peer Address:Port [Process]
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		netid := strings.ToLower(fields[0])
		if netid != "tcp" && netid != "udp" {
			continue
		}

		// Only show listening / unconnected (for UDP) ports.
		state := fields[1]
		if netid == "tcp" && state != "LISTEN" {
			continue
		}

		local := fields[4]
		addr, portStr, ok := parseAddrPort(local)
		if !ok {
			continue
		}
		port, err := strconv.Atoi(portStr)
		if err != nil {
			continue
		}

		pi := PortInfo{
			Protocol: netid,
			Port:     port,
			Address:  addr,
		}

		// Parse process info if present (last field: users:(("name",pid=N,fd=N)))
		if len(fields) >= 7 && strings.HasPrefix(fields[6], "users:((") {
			procField := strings.Join(fields[6:], " ")
			pi.Process, pi.PID = parseProcess(procField)
		}

		ports = append(ports, pi)
	}

	sort.Slice(ports, func(i, j int) bool {
		if ports[i].Protocol != ports[j].Protocol {
			return ports[i].Protocol < ports[j].Protocol
		}
		return ports[i].Port < ports[j].Port
	})

	return ports
}

func parseAddrPort(local string) (addr, port string, ok bool) {
	// IPv6 addresses use [...] notation: [::1]:631
	if strings.HasPrefix(local, "[") {
		idx := strings.LastIndex(local, "]:")
		if idx < 0 {
			return "", "", false
		}
		return local[:idx+1], local[idx+2:], true
	}
	// IPv4: 127.0.0.1:631 or 0.0.0.0:5353
	idx := strings.LastIndex(local, ":")
	if idx < 0 {
		return "", "", false
	}
	return local[:idx], local[idx+1:], true
}

// KillPort kills the process listening on the given port and protocol.
// If force is true, sends SIGKILL instead of SIGTERM.
// Returns the PID that was killed, or an error.
func KillPort(port int, protocol string, force bool) (int, error) {
	proto := strings.ToLower(protocol)
	if proto == "tcp6" {
		proto = "tcp"
	}
	if proto == "udp6" {
		proto = "udp"
	}

	spec := fmt.Sprintf("%d/%s", port, proto)
	out, err := exec.Command("fuser", spec).CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("no process found on port %d/%s", port, proto)
	}

	pidStr := strings.TrimSpace(string(out))
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return 0, fmt.Errorf("could not parse PID from fuser output: %s", pidStr)
	}

	sig := syscall.SIGTERM
	if force {
		sig = syscall.SIGKILL
	}
	if err := syscall.Kill(pid, sig); err != nil {
		return 0, fmt.Errorf("failed to kill PID %d: %w", pid, err)
	}

	return pid, nil
}

func parseProcess(field string) (name string, pid int) {
	// Example: users:(("cupsd",pid=1234,fd=5))
	field = strings.TrimPrefix(field, "users:((")
	field = strings.TrimSuffix(field, "))")

	parts := strings.Split(field, ",")
	if len(parts) >= 2 {
		name = strings.Trim(parts[0], "\"")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if strings.HasPrefix(p, "pid=") {
				pidStr := strings.TrimPrefix(p, "pid=")
				if n, err := strconv.Atoi(pidStr); err == nil {
					pid = n
				}
			}
		}
	}
	return
}
