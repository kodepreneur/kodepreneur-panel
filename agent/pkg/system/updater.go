package system

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type UpdateRequest struct {
	Repository string `json:"repository"`
	Branch     string `json:"branch"`
	Force      bool   `json:"force"`
}

type UpdateResult struct {
	Success         bool   `json:"success"`
	ExitCode        int    `json:"exit_code"`
	LogOutput       string `json:"log_output"`
	DurationSeconds int    `json:"duration_seconds"`
	CommitHash      string `json:"commit_hash"`
	CommitMessage   string `json:"commit_message"`
}

type UpdateRunner struct {
	isDev bool
}

func NewUpdateRunner(isDev bool) *UpdateRunner {
	return &UpdateRunner{
		isDev: isDev || runtime.GOOS != "linux",
	}
}

func (u *UpdateRunner) Execute(req UpdateRequest) (*UpdateResult, error) {
	if req.Branch == "" {
		req.Branch = "main"
	}
	if req.Repository == "" {
		req.Repository = "https://github.com/kodepreneur/kodepreneur-panel.git"
	}

	start := time.Now()

	if u.isDev {
		// Mock execution in development/mac environment
		time.Sleep(150 * time.Millisecond)
		duration := int(time.Since(start).Seconds())
		if duration < 1 {
			duration = 1
		}

		var logs bytes.Buffer
		logs.WriteString(fmt.Sprintf("[%s] 🚀 Starting Kodepreneur Panel Automated Update Pipeline...\n", time.Now().Format("15:04:05")))
		logs.WriteString(fmt.Sprintf("[%s] [1/6] Fetching latest release from %s (branch: %s)...\n", time.Now().Format("15:04:05"), req.Repository, req.Branch))
		logs.WriteString("  - Cloned latest commits successfully (HEAD -> 8706c6c)\n")
		logs.WriteString(fmt.Sprintf("[%s] [2/6] Synchronizing Control Plane files and updating Composer packages...\n", time.Now().Format("15:04:05")))
		logs.WriteString("  - Installing dependencies (no-dev, --optimize-autoloader)\n  - Generated optimized autoload files containing 3,420 classes\n")
		logs.WriteString(fmt.Sprintf("[%s] [3/6] Running database migrations...\n", time.Now().Format("15:04:05")))
		logs.WriteString("  - Nothing to migrate / Migrations applied up to date\n")
		logs.WriteString(fmt.Sprintf("[%s] [4/6] Compiling frontend assets via Vite...\n", time.Now().Format("15:04:05")))
		logs.WriteString("  - vite v6.4.3 building for production...\n  - ✓ built in 1.45s\n")
		logs.WriteString(fmt.Sprintf("[%s] [5/6] Optimizing Laravel application caches and permissions...\n", time.Now().Format("15:04:05")))
		logs.WriteString("  - Configuration cache cleared and rebuilt.\n  - Route cache optimized.\n  - Set ownership to www-data:www-data\n")
		logs.WriteString(fmt.Sprintf("[%s] [6/6] Recompiling Go Agent daemon and reloading Nginx...\n", time.Now().Format("15:04:05")))
		logs.WriteString("  - kodepreneur-agent binary updated: /usr/local/bin/kodepreneur-agent\n  - Restarted kodepreneur-agent.service\n  - Nginx configuration valid. Reloaded nginx.service\n")
		logs.WriteString(fmt.Sprintf("[%s] 🎉 Kodepreneur Panel successfully updated in %ds!\n", time.Now().Format("15:04:05"), duration))

		return &UpdateResult{
			Success:         true,
			ExitCode:        0,
			LogOutput:       logs.String(),
			DurationSeconds: duration,
			CommitHash:      "8706c6c49832",
			CommitMessage:   "feat(updates): add automated in-panel GitHub upgrade manager",
		}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	var logs bytes.Buffer
	logs.WriteString(fmt.Sprintf("[%s] [kodepreneur-update] Initializing root update runner...\n", time.Now().Format("15:04:05")))

	// Try running installer/update.sh if it exists, or run inline update script
	updateScriptPath := "/var/www/kodepreneur-panel/installer/update.sh"
	cmd := exec.CommandContext(ctx, "/bin/bash", updateScriptPath)
	out, err := cmd.CombinedOutput()
	logs.Write(out)

	duration := int(time.Since(start).Seconds())
	if duration < 1 {
		duration = 1
	}

	exitCode := 0
	success := true
	if err != nil {
		success = false
		exitCode = 1
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
		logs.WriteString(fmt.Sprintf("\n[error] Update script exited with code %d: %v\n", exitCode, err))
	} else {
		logs.WriteString(fmt.Sprintf("\n[success] Panel update successfully completed in %ds.\n", duration))
	}

	// Extract updated commit metadata
	commitHash := ""
	commitMsg := ""
	gitRevCmd := exec.CommandContext(ctx, "git", "-C", "/var/www/kodepreneur-panel", "rev-parse", "--short", "HEAD")
	if revOut, revErr := gitRevCmd.Output(); revErr == nil {
		commitHash = strings.TrimSpace(string(revOut))
	}
	gitMsgCmd := exec.CommandContext(ctx, "git", "-C", "/var/www/kodepreneur-panel", "log", "-1", "--pretty=%B")
	if msgOut, msgErr := gitMsgCmd.Output(); msgErr == nil {
		commitMsg = strings.TrimSpace(string(msgOut))
	}

	return &UpdateResult{
		Success:         success,
		ExitCode:        exitCode,
		LogOutput:       logs.String(),
		DurationSeconds: duration,
		CommitHash:      commitHash,
		CommitMessage:   commitMsg,
	}, nil
}
