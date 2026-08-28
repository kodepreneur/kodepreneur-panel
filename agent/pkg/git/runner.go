package git

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type DeploymentRequest struct {
	SystemUser string   `json:"system_user"`
	WorkingDir string   `json:"working_dir"`
	Repository string   `json:"repository"`
	Branch     string   `json:"branch"`
	Commands   []string `json:"commands"`
	TimeoutSec int      `json:"timeout_sec"`
}

type DeploymentResult struct {
	Success         bool   `json:"success"`
	ExitCode        int    `json:"exit_code"`
	LogOutput       string `json:"log_output"`
	DurationSeconds int    `json:"duration_seconds"`
	CommitHash      string `json:"commit_hash"`
	CommitMessage   string `json:"commit_message"`
}

type Runner struct {
	isDev bool
}

func NewRunner(isDev bool) *Runner {
	return &Runner{
		isDev: isDev || runtime.GOOS != "linux",
	}
}

func (r *Runner) Execute(req DeploymentRequest) (*DeploymentResult, error) {
	if req.TimeoutSec <= 0 {
		req.TimeoutSec = 600 // 10 minutes default
	}
	if req.Branch == "" {
		req.Branch = "main"
	}

	start := time.Now()

	if r.isDev {
		// Mock execution
		time.Sleep(100 * time.Millisecond)
		duration := int(time.Since(start).Seconds())
		if duration < 1 {
			duration = 1
		}

		var logs bytes.Buffer
		logs.WriteString(fmt.Sprintf("[%s] Starting unprivileged deployment as %s in %s\n", time.Now().Format("15:04:05"), req.SystemUser, req.WorkingDir))
		logs.WriteString(fmt.Sprintf("[%s] [git] Checking out branch %s...\n", time.Now().Format("15:04:05"), req.Branch))
		logs.WriteString(fmt.Sprintf("[%s] [git] HEAD is now at 8706c6c (Release update: v1.0.0-prod)\n", time.Now().Format("15:04:05")))

		for _, cmdStr := range req.Commands {
			logs.WriteString(fmt.Sprintf("[%s] [exec] %s\n", time.Now().Format("15:04:05"), cmdStr))
			if strings.Contains(cmdStr, "composer") {
				logs.WriteString("  - Installing dependencies (no-dev, optimize-autoloader)\n  - Generated optimized autoload files containing 3420 classes\n")
			} else if strings.Contains(cmdStr, "artisan") {
				logs.WriteString("  - Configuration cache cleared!\n  - Configuration cached successfully.\n  - Route cache cleared!\n  - Routes cached successfully.\n")
			} else if strings.Contains(cmdStr, "npm") {
				logs.WriteString("  - vite v6.4.3 building for production...\n  - ✓ built in 1.45s\n")
			}
		}
		logs.WriteString(fmt.Sprintf("[%s] Deployment completed successfully in %ds.\n", time.Now().Format("15:04:05"), duration))

		return &DeploymentResult{
			Success:         true,
			ExitCode:        0,
			LogOutput:       logs.String(),
			DurationSeconds: duration,
			CommitHash:      "8706c6c49832",
			CommitMessage:   "Release update: v1.0.0-prod",
		}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(req.TimeoutSec)*time.Second)
	defer cancel()

	var logs bytes.Buffer
	logs.WriteString(fmt.Sprintf("[%s] [kodepreneur-deploy] Executing pipeline for %s (user: %s)\n", time.Now().Format("15:04:05"), req.WorkingDir, req.SystemUser))

	// Build shell script executed as system_user
	var scriptBuilder strings.Builder
	scriptBuilder.WriteString("set -e\n")
	scriptBuilder.WriteString(fmt.Sprintf("cd %s\n", req.WorkingDir))

	for _, cmdStr := range req.Commands {
		cmdStr = strings.TrimSpace(cmdStr)
		if cmdStr != "" {
			scriptBuilder.WriteString(fmt.Sprintf("echo '[exec] %s'\n", cmdStr))
			scriptBuilder.WriteString(fmt.Sprintf("%s\n", cmdStr))
		}
	}

	// Run via `su - <user> -s /bin/bash -c '<script>'`
	cmd := exec.CommandContext(ctx, "su", "-", req.SystemUser, "-s", "/bin/bash", "-c", scriptBuilder.String())
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
		logs.WriteString(fmt.Sprintf("\n[error] Deployment command failed with exit code %d: %v\n", exitCode, err))
	} else {
		logs.WriteString(fmt.Sprintf("\n[success] Deployment completed successfully in %ds.\n", duration))
	}

	// Try extracting git commit metadata
	commitHash := ""
	commitMsg := ""
	gitRevCmd := exec.CommandContext(ctx, "su", "-", req.SystemUser, "-s", "/bin/bash", "-c", fmt.Sprintf("cd %s && git rev-parse --short HEAD 2>/dev/null || true", req.WorkingDir))
	if revOut, revErr := gitRevCmd.Output(); revErr == nil {
		commitHash = strings.TrimSpace(string(revOut))
	}
	gitMsgCmd := exec.CommandContext(ctx, "su", "-", req.SystemUser, "-s", "/bin/bash", "-c", fmt.Sprintf("cd %s && git log -1 --pretty=%%B 2>/dev/null || true", req.WorkingDir))
	if msgOut, msgErr := gitMsgCmd.Output(); msgErr == nil {
		commitMsg = strings.TrimSpace(string(msgOut))
	}

	return &DeploymentResult{
		Success:         success,
		ExitCode:        exitCode,
		LogOutput:       logs.String(),
		DurationSeconds: duration,
		CommitHash:      commitHash,
		CommitMessage:   commitMsg,
	}, nil
}

// CloneRepo clones a git repository into targetDir as the systemUser.
func (r *Runner) CloneRepo(repoUrl, branch, targetDir, systemUser string) error {
	if branch == "" {
		branch = "main"
	}

	if r.isDev {
		_ = os.MkdirAll(targetDir, 0755)
		sampleIndex := filepath.Join(targetDir, "public", "index.php")
		_ = os.MkdirAll(filepath.Dir(sampleIndex), 0755)
		_ = os.WriteFile(sampleIndex, []byte("<?php echo 'Cloned from "+repoUrl+" ("+branch+")';"), 0644)
		_ = os.WriteFile(filepath.Join(targetDir, "artisan"), []byte("#!/usr/bin/env php\n<?php // artisan"), 0755)
		return nil
	}

	_ = os.MkdirAll(targetDir, 0755)
	cloneCmd := fmt.Sprintf("git clone --branch %s --depth 1 %s %s", branch, repoUrl, targetDir)
	cmd := exec.Command("su", "-", systemUser, "-s", "/bin/bash", "-c", cloneCmd)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git clone failed: %s: %s", err.Error(), strings.TrimSpace(string(out)))
	}

	return nil
}
