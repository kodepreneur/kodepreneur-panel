package process

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

type ProcessConfig struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Command     string `json:"command"`
	SystemUser  string `json:"system_user"`
	WorkingDir  string `json:"working_dir"`
	AutoRestart bool   `json:"auto_restart"`
}

const serviceUnitTemplate = `[Unit]
Description=Kodepreneur Managed Process - {{ .Name }} ({{ .ID }})
After=network.target

[Service]
Type=simple
User={{ .SystemUser }}
Group={{ .SystemUser }}
WorkingDirectory={{ .WorkingDir }}
ExecStart={{ .Command }}
Restart={{ if .AutoRestart }}always{{ else }}on-failure{{ end }}
RestartSec=5s
StandardOutput=journal
StandardError=journal
Environment=PATH=/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin

[Install]
WantedBy=multi-user.target
`

type Supervisor struct {
	systemdDir string
	isDev      bool
}

func NewSupervisor(isDev bool) *Supervisor {
	if isDev || runtime.GOOS != "linux" {
		mockDir := filepath.Join(os.TempDir(), "kodepreneur", "systemd")
		_ = os.MkdirAll(mockDir, 0755)
		return &Supervisor{
			systemdDir: mockDir,
			isDev:      true,
		}
	}

	return &Supervisor{
		systemdDir: "/etc/systemd/system",
		isDev:      false,
	}
}

func (s *Supervisor) unitName(id string) string {
	return fmt.Sprintf("kp-proc-%s.service", id)
}

func (s *Supervisor) CreateProcess(cfg ProcessConfig) (string, error) {
	if cfg.ID == "" || cfg.Command == "" {
		return "", fmt.Errorf("id and command are required")
	}
	if cfg.SystemUser == "" {
		cfg.SystemUser = "kodepreneur"
	}
	if cfg.WorkingDir == "" {
		cfg.WorkingDir = "/var/www"
	}

	tmpl, err := template.New("unit").Parse(serviceUnitTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse unit template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, cfg); err != nil {
		return "", fmt.Errorf("failed to render unit template: %w", err)
	}

	unitPath := filepath.Join(s.systemdDir, s.unitName(cfg.ID))
	if err := os.WriteFile(unitPath, buf.Bytes(), 0644); err != nil {
		return "", fmt.Errorf("failed to write unit file %s: %w", unitPath, err)
	}

	if !s.isDev {
		_ = exec.Command("systemctl", "daemon-reload").Run()
		_ = exec.Command("systemctl", "enable", s.unitName(cfg.ID)).Run()
		_ = exec.Command("systemctl", "start", s.unitName(cfg.ID)).Run()
	}

	return unitPath, nil
}

func (s *Supervisor) ManageProcess(id, action string) (string, error) {
	unit := s.unitName(id)

	if s.isDev {
		return fmt.Sprintf("Mock executed %s on %s", action, unit), nil
	}

	cmd := exec.Command("systemctl", action, unit)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("systemctl %s %s failed: %s: %s", action, unit, err.Error(), strings.TrimSpace(string(out)))
	}
	return strings.TrimSpace(string(out)), nil
}

func (s *Supervisor) DeleteProcess(id string) error {
	unit := s.unitName(id)
	unitPath := filepath.Join(s.systemdDir, unit)

	if !s.isDev {
		_ = exec.Command("systemctl", "stop", unit).Run()
		_ = exec.Command("systemctl", "disable", unit).Run()
		_ = os.Remove(unitPath)
		_ = exec.Command("systemctl", "daemon-reload").Run()
	} else {
		_ = os.Remove(unitPath)
	}

	return nil
}
