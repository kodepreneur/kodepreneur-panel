package cron

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

type JobItem struct {
	Schedule string `json:"schedule"` // e.g. "* * * * *"
	Command  string `json:"command"`  // e.g. "php /var/www/site/artisan schedule:run"
	IsActive bool   `json:"is_active"`
}

type Manager struct {
	isDev bool
}

func NewManager(isDev bool) *Manager {
	return &Manager{
		isDev: isDev || runtime.GOOS != "linux",
	}
}

// ValidateSchedule checks if standard 5-field cron format is given.
func ValidateSchedule(schedule string) error {
	fields := strings.Fields(strings.TrimSpace(schedule))
	if len(fields) != 5 {
		return fmt.Errorf("invalid cron schedule: must contain 5 space-separated fields, got %d", len(fields))
	}
	return nil
}

// SyncUserCrontab rewrites crontab for system user.
func (m *Manager) SyncUserCrontab(systemUser string, jobs []JobItem) error {
	if systemUser == "" {
		return fmt.Errorf("system user cannot be empty")
	}

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("# Managed by Kodepreneur Control Plane for user %s\n", systemUser))
	buf.WriteString("SHELL=/bin/bash\nPATH=/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin\n\n")

	for _, job := range jobs {
		if !job.IsActive {
			buf.WriteString("# (Disabled) ")
		}
		buf.WriteString(fmt.Sprintf("%s %s\n", job.Schedule, job.Command))
	}

	if m.isDev {
		return nil
	}

	cmd := exec.Command("crontab", "-u", systemUser, "-")
	cmd.Stdin = &buf
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to write crontab for %s: %s: %s", systemUser, err.Error(), strings.TrimSpace(string(out)))
	}

	return nil
}
