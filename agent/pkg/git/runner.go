package git

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type DeploymentRequest struct {
	SystemUser    string   `json:"system_user"`
	WorkingDir    string   `json:"working_dir"`
	Repository    string   `json:"repository"`
	Branch        string   `json:"branch"`
	Commands      []string `json:"commands"`
	TimeoutSec    int      `json:"timeout_sec"`
	AuthType      string   `json:"auth_type"`
	SshPrivateKey string   `json:"ssh_private_key"`
	GitToken      string   `json:"git_token"`
	GitTokenUser  string   `json:"git_token_user"`
}

type DeploymentResult struct {
	Success         bool   `json:"success"`
	ExitCode        int    `json:"exit_code"`
	LogOutput       string `json:"log_output"`
	DurationSeconds int    `json:"duration_seconds"`
	CommitHash      string `json:"commit_hash"`
	CommitMessage   string `json:"commit_message"`
}

type CloneOptions struct {
	RepoUrl       string `json:"repo_url"`
	Branch        string `json:"branch"`
	TargetDir     string `json:"target_dir"`
	SystemUser    string `json:"system_user"`
	AuthType      string `json:"auth_type"` // "none", "ssh_key", "token"
	SshPrivateKey string `json:"ssh_private_key"`
	GitToken      string `json:"git_token"`
	GitTokenUser  string `json:"git_token_user"`
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
		if req.AuthType == "ssh_key" {
			logs.WriteString(fmt.Sprintf("[%s] [git] Using SSH Deploy Key authentication\n", time.Now().Format("15:04:05")))
		} else if req.AuthType == "token" {
			logs.WriteString(fmt.Sprintf("[%s] [git] Using Personal Access Token authentication\n", time.Now().Format("15:04:05")))
		}
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
	scriptBuilder.WriteString("export PATH=\"/usr/local/bin:/usr/bin:/bin:$PATH\"\n")
	scriptBuilder.WriteString("export GIT_TERMINAL_PROMPT=0\n")
	scriptBuilder.WriteString("export GIT_SSH_COMMAND=\"ssh -o StrictHostKeyChecking=accept-new -o BatchMode=yes\"\n")
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

// Clone clones a git repository according to CloneOptions.
func (r *Runner) Clone(opts CloneOptions) error {
	if opts.Branch == "" {
		opts.Branch = "main"
	}
	if opts.AuthType == "" {
		opts.AuthType = "none"
	}

	effectiveUrl := opts.RepoUrl

	// If token authentication, format the authenticated URL
	if opts.AuthType == "token" && opts.GitToken != "" {
		effectiveUrl = buildAuthenticatedUrl(opts.RepoUrl, opts.GitToken, opts.GitTokenUser)
	}

	if r.isDev {
		_ = os.MkdirAll(opts.TargetDir, 0755)
		sampleIndex := filepath.Join(opts.TargetDir, "public", "index.php")
		_ = os.MkdirAll(filepath.Dir(sampleIndex), 0755)
		displayUrl := maskSecretUrl(effectiveUrl)
		_ = os.WriteFile(sampleIndex, []byte("<?php echo 'Cloned from "+displayUrl+" ("+opts.Branch+")';"), 0644)
		_ = os.WriteFile(filepath.Join(opts.TargetDir, "artisan"), []byte("#!/usr/bin/env php\n<?php // artisan"), 0755)
		return nil
	}

	// 1. Ensure target directory exists and is owned by systemUser:www-data
	if err := os.MkdirAll(opts.TargetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory %s: %w", opts.TargetDir, err)
	}
	chownCmd := exec.Command("chown", "-R", fmt.Sprintf("%s:www-data", opts.SystemUser), opts.TargetDir)
	_ = chownCmd.Run()
	chmodCmd := exec.Command("chmod", "0755", opts.TargetDir)
	_ = chmodCmd.Run()

	// 2. Setup SSH Key if SSH Auth is specified
	if opts.AuthType == "ssh_key" && strings.TrimSpace(opts.SshPrivateKey) != "" {
		if err := setupUserSshKey(opts.SystemUser, opts.TargetDir, opts.SshPrivateKey); err != nil {
			return fmt.Errorf("failed to configure SSH deploy key: %w", err)
		}
	}

	// 3. Clone or initialize git repository in place as the system user
	cloneScript := fmt.Sprintf(`set -e
export GIT_TERMINAL_PROMPT=0
export GIT_SSH_COMMAND="ssh -o StrictHostKeyChecking=accept-new -o BatchMode=yes"
git config --global --add safe.directory "%s" 2>/dev/null || true
cd "%s"
if [ ! -d ".git" ]; then
    git init -b "%s" 2>/dev/null || (git init && git checkout -B "%s" 2>/dev/null || true)
    git remote add origin "%s" 2>/dev/null || git remote set-url origin "%s"
    git fetch --depth 1 origin "%s"
    git reset --hard "origin/%s" || git checkout -f -B "%s" "origin/%s"
else
    git remote set-url origin "%s"
    git fetch --depth 1 origin "%s"
    git reset --hard "origin/%s"
fi
`, opts.TargetDir, opts.TargetDir, opts.Branch, opts.Branch, effectiveUrl, effectiveUrl, opts.Branch, opts.Branch, opts.Branch, opts.Branch, effectiveUrl, opts.Branch, opts.Branch)

	cmd := exec.Command("su", "-", opts.SystemUser, "-s", "/bin/bash", "-c", cloneScript)
	if out, err := cmd.CombinedOutput(); err != nil {
		outStr := strings.TrimSpace(string(out))
		if opts.GitToken != "" {
			outStr = strings.ReplaceAll(outStr, opts.GitToken, "******")
		}
		return fmt.Errorf("git clone failed: %s: %s", err.Error(), outStr)
	}

	// 4. Ensure permissions across all cloned files
	postChownCmd := exec.Command("chown", "-R", fmt.Sprintf("%s:www-data", opts.SystemUser), opts.TargetDir)
	_ = postChownCmd.Run()

	return nil
}

// CloneRepo provides backward compatibility.
func (r *Runner) CloneRepo(repoUrl, branch, targetDir, systemUser string) error {
	return r.Clone(CloneOptions{
		RepoUrl:    repoUrl,
		Branch:     branch,
		TargetDir:  targetDir,
		SystemUser: systemUser,
		AuthType:   "none",
	})
}

// setupUserSshKey writes the SSH deploy key and configures known_hosts for the system user.
func setupUserSshKey(systemUser, userHome, privateKeyContent string) error {
	sshDir := filepath.Join(userHome, ".ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return err
	}

	keyFile := filepath.Join(sshDir, "id_rsa")
	cleanKey := strings.TrimSpace(privateKeyContent) + "\n"
	if err := os.WriteFile(keyFile, []byte(cleanKey), 0600); err != nil {
		return err
	}

	// Setup known hosts for common git hosts
	knownHostsFile := filepath.Join(sshDir, "known_hosts")
	if _, err := os.Stat(knownHostsFile); os.IsNotExist(err) {
		_ = os.WriteFile(knownHostsFile, []byte(""), 0644)
	}

	// Setup SSH config
	configFile := filepath.Join(sshDir, "config")
	configContent := "Host github.com gitlab.com bitbucket.org *\n    StrictHostKeyChecking accept-new\n    IdentityFile ~/.ssh/id_rsa\n    BatchMode yes\n"
	_ = os.WriteFile(configFile, []byte(configContent), 0644)

	// Ensure ownership is systemUser:www-data
	_ = exec.Command("chown", "-R", fmt.Sprintf("%s:www-data", systemUser), sshDir).Run()
	_ = exec.Command("chmod", "700", sshDir).Run()
	_ = exec.Command("chmod", "600", keyFile).Run()

	return nil
}

// buildAuthenticatedUrl embeds authentication tokens into HTTPS git URLs.
func buildAuthenticatedUrl(repoUrl, token, user string) string {
	repoUrl = strings.TrimSpace(repoUrl)
	if !strings.HasPrefix(repoUrl, "http://") && !strings.HasPrefix(repoUrl, "https://") {
		return repoUrl
	}

	parsed, err := url.Parse(repoUrl)
	if err != nil {
		return repoUrl
	}

	// Determine token user
	tokenUser := strings.TrimSpace(user)
	if tokenUser == "" {
		hostLower := strings.ToLower(parsed.Host)
		if strings.Contains(hostLower, "gitlab") {
			tokenUser = "oauth2"
		} else if strings.Contains(hostLower, "github") {
			tokenUser = "x-access-token"
		} else {
			tokenUser = ""
		}
	}

	if tokenUser != "" {
		parsed.User = url.UserPassword(tokenUser, token)
	} else {
		parsed.User = url.User(token)
	}

	return parsed.String()
}

// maskSecretUrl removes credentials from URLs for display/logging.
func maskSecretUrl(rawUrl string) string {
	parsed, err := url.Parse(rawUrl)
	if err != nil || parsed.User == nil {
		return rawUrl
	}
	parsed.User = url.User("******")
	return parsed.String()
}
