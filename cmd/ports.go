package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hades/godex/internal/runtime"
	"github.com/spf13/cobra"
)

var (
	killForce bool
	killProto string
)

var portsCmd = &cobra.Command{
	Use:     "ports",
	Aliases: []string{"port"},
	Short:   "List and manage listening ports",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listPorts(cmd)
	},
}

var portsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all listening ports",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listPorts(cmd)
	},
}

var portsSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search ports by port number, process, or address",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ports, err := runtime.ListPorts()
		if err != nil {
			return err
		}
		q := strings.ToLower(args[0])
		var found bool
		for _, p := range ports {
			if strings.Contains(strconv.Itoa(p.Port), q) ||
				strings.Contains(strings.ToLower(p.Process), q) ||
				strings.Contains(strings.ToLower(p.Address), q) {
				found = true
				printPort(cmd, p)
			}
		}
		if !found {
			fmt.Fprintf(cmd.OutOrStdout(), "No ports match %q.\n", args[0])
		}
		return nil
	},
}

var portsKillCmd = &cobra.Command{
	Use:   "kill <port>",
	Short: "Kill the process listening on a port",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		port, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid port number: %s", args[0])
		}
		proto := strings.ToLower(killProto)
		if proto == "" {
			proto = "tcp"
		}
		pid, err := runtime.KillPort(port, proto, killForce)
		if err != nil {
			return err
		}
		sig := "SIGTERM"
		if killForce {
			sig = "SIGKILL"
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Killed PID %d on port %d/%s (%s).\n", pid, port, proto, sig)
		return nil
	},
}

func listPorts(cmd *cobra.Command) error {
	ports, err := runtime.ListPorts()
	if err != nil {
		return err
	}
	if len(ports) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No listening ports found.")
		return nil
	}
	fmt.Fprintf(cmd.OutOrStdout(), "%-6s %-7s %-22s %s\n", "PROTO", "PORT", "ADDRESS", "PROCESS")
	fmt.Fprintln(cmd.OutOrStdout(), strings.Repeat("-", 60))
	for _, p := range ports {
		printPort(cmd, p)
	}
	return nil
}

func printPort(cmd *cobra.Command, p runtime.PortInfo) {
	proc := p.Process
	if proc == "" {
		proc = "-"
	}
	if p.PID > 0 {
		proc = fmt.Sprintf("%s (%d)", proc, p.PID)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "%-6s %-7d %-22s %s\n",
		strings.ToUpper(p.Protocol), p.Port, p.Address, proc)
}

func init() {
	portsKillCmd.Flags().BoolVarP(&killForce, "force", "f", false, "Use SIGKILL instead of SIGTERM")
	portsKillCmd.Flags().StringVarP(&killProto, "proto", "p", "tcp", "Protocol (tcp or udp)")

	portsCmd.AddCommand(portsListCmd)
	portsCmd.AddCommand(portsSearchCmd)
	portsCmd.AddCommand(portsKillCmd)
	rootCmd.AddCommand(portsCmd)
}
