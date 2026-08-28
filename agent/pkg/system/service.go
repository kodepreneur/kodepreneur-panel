package system

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

var AllowedServices = map[string]bool{
	"nginx":               true,
	"php8.3-fpm":          true,
	"php8.4-fpm":          true,
	"mysql":               true,
	"mariadb":             true,
	"postgresql":          true,
	"kodepreneur-agent":   true,
	"kodepreneur-worker":  true,
	"ufw":                 true,
	"fail2ban":            true,
}

var AllowedActions = map[string]bool{
	"status":  true,
	"start":   true,
	"stop":    true,
	"restart": true,
	"reload":  true,
}

type ServiceStatus struct {
	Name      string `json:"name"`
	IsActive  bool   `json:"is_active"`
	IsEnabled bool   `json:"is_enabled"`
	Status    string `json:"status"`
}

func ManageService(serviceName, action string) (string, error) {
	if !AllowedServices[serviceName] {
		return "", fmt.Errorf("service '%s' is not in the managed whitelist", serviceName)
	}

	if !AllowedActions[action] {
		return "", fmt.Errorf("action '%s' is not allowed", action)
	}

	if runtime.GOOS != "linux" {
		return fmt.Sprintf("[DEV-MOCK] Executed systemctl %s %s successfully", action, serviceName), nil
	}

	cmd := exec.Command("systemctl", action, serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("systemctl %s %s failed: %w (output: %s)", action, serviceName, err, strings.TrimSpace(string(output)))
	}

	return strings.TrimSpace(string(output)), nil
}

func GetServiceStatus(serviceName string) (ServiceStatus, error) {
	if !AllowedServices[serviceName] {
		return ServiceStatus{}, errors.New("service not in whitelist")
	}

	if runtime.GOOS != "linux" {
		return ServiceStatus{
			Name:      serviceName,
			IsActive:  true,
			IsEnabled: true,
			Status:    "active (running) [mock]",
		}, nil
	}

	isActive := false
	if err := exec.Command("systemctl", "is-active", "--quiet", serviceName).Run(); err == nil {
		isActive = true
	}

	isEnabled := false
	if err := exec.Command("systemctl", "is-enabled", "--quiet", serviceName).Run(); err == nil {
		isEnabled = true
	}

	out, _ := exec.Command("systemctl", "status", serviceName, "--no-pager").CombinedOutput()

	return ServiceStatus{
		Name:      serviceName,
		IsActive:  isActive,
		IsEnabled: isEnabled,
		Status:    strings.TrimSpace(string(out)),
	}, nil
}
