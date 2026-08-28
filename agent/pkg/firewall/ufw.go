package firewall

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

type Rule struct {
	Port        int    `json:"port"`
	Protocol    string `json:"protocol"` // "tcp", "udp", "both"
	Action      string `json:"action"`   // "allow", "deny"
	SourceIP    string `json:"source_ip"`
	Description string `json:"description"`
}

type Manager struct {
	isDev bool
}

func NewManager(isDev bool) *Manager {
	return &Manager{
		isDev: isDev || runtime.GOOS != "linux",
	}
}

func (m *Manager) AddRule(r Rule) error {
	if r.Port <= 0 || r.Port > 65535 {
		return fmt.Errorf("invalid port number: %d", r.Port)
	}
	if r.Protocol == "" {
		r.Protocol = "tcp"
	}
	if r.Action == "" {
		r.Action = "allow"
	}

	if m.isDev {
		return nil
	}

	var args []string
	args = append(args, r.Action)

	if r.SourceIP != "" && r.SourceIP != "0.0.0.0/0" {
		args = append(args, "from", r.SourceIP, "to", "any", "port", fmt.Sprintf("%d", r.Port))
		if r.Protocol != "both" {
			args = append(args, "proto", r.Protocol)
		}
	} else {
		if r.Protocol == "both" {
			args = append(args, fmt.Sprintf("%d", r.Port))
		} else {
			args = append(args, fmt.Sprintf("%d/%s", r.Port, r.Protocol))
		}
	}

	cmd := exec.Command("ufw", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ufw rule addition failed: %s: %s", err.Error(), strings.TrimSpace(string(out)))
	}

	return nil
}

func (m *Manager) DeleteRule(port int, protocol, action string) error {
	if m.isDev {
		return nil
	}

	target := fmt.Sprintf("%d/%s", port, protocol)
	if protocol == "both" || protocol == "" {
		target = fmt.Sprintf("%d", port)
	}

	cmd := exec.Command("ufw", "delete", action, target)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ufw rule deletion failed: %s: %s", err.Error(), strings.TrimSpace(string(out)))
	}

	return nil
}
