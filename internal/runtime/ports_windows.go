//go:build windows

package runtime

import "fmt"

// KillPort is not yet implemented on Windows.
// A future version will use `taskkill /PID <pid>` via `netstat -ano`.
func KillPort(port int, protocol string, force bool) (int, error) {
	return 0, fmt.Errorf("port kill is not supported on Windows yet")
}
