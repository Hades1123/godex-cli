//go:build !windows

package runtime

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

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

	// fuser output: "8080/tcp:            48518\n"
	// Extract PID from after the colon.
	raw := strings.TrimSpace(string(out))
	parts := strings.SplitN(raw, ":", 2)
	pidStr := raw
	if len(parts) == 2 {
		pidStr = strings.TrimSpace(parts[1])
	}
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return 0, fmt.Errorf("could not parse PID from fuser output: %s", raw)
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
