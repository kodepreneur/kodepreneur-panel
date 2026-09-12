package git

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGitRunnerDevExecution(t *testing.T) {
	runner := NewRunner(true)

	req := DeploymentRequest{
		SystemUser: "kp_test",
		WorkingDir: "/var/www/test.com",
		Branch:     "main",
		Commands: []string{
			"git pull origin main",
			"composer install",
			"php artisan config:cache",
		},
	}

	res, err := runner.Execute(req)
	if err != nil {
		t.Fatalf("Runner execute failed: %v", err)
	}

	if !res.Success {
		t.Errorf("Expected success, got false")
	}

	if !strings.Contains(res.LogOutput, "Deployment completed successfully") {
		t.Errorf("Expected completion message in log output, got:\n%s", res.LogOutput)
	}
}

func TestGitRunnerCloneSshAndToken(t *testing.T) {
	runner := NewRunner(true)

	tmpDir := filepath.Join(os.TempDir(), "kp_test_git_clone")
	defer os.RemoveAll(tmpDir)

	// 1. Test SSH clone
	sshOpts := CloneOptions{
		RepoUrl:       "git@github.com:example/private-app.git",
		Branch:        "main",
		TargetDir:     filepath.Join(tmpDir, "ssh-site"),
		SystemUser:    "kp_ssh",
		AuthType:      "ssh_key",
		SshPrivateKey: "-----BEGIN OPENSSH PRIVATE KEY-----\ntest\n-----END OPENSSH PRIVATE KEY-----",
	}

	if err := runner.Clone(sshOpts); err != nil {
		t.Fatalf("SSH Clone failed: %v", err)
	}

	// 2. Test Token clone
	tokenOpts := CloneOptions{
		RepoUrl:      "https://github.com/example/token-app.git",
		Branch:       "staging",
		TargetDir:    filepath.Join(tmpDir, "token-site"),
		SystemUser:   "kp_token",
		AuthType:     "token",
		GitToken:     "ghp_secret12345",
		GitTokenUser: "x-access-token",
	}

	if err := runner.Clone(tokenOpts); err != nil {
		t.Fatalf("Token Clone failed: %v", err)
	}
}

func TestBuildAuthenticatedUrlAndMasking(t *testing.T) {
	// GitHub token URL
	ghUrl := buildAuthenticatedUrl("https://github.com/org/repo.git", "ghp_12345", "")
	if !strings.Contains(ghUrl, "x-access-token:ghp_12345@github.com") {
		t.Errorf("Expected x-access-token in github url, got %s", ghUrl)
	}

	// GitLab token URL
	glUrl := buildAuthenticatedUrl("https://gitlab.com/org/repo.git", "glpat-12345", "")
	if !strings.Contains(glUrl, "oauth2:glpat-12345@gitlab.com") {
		t.Errorf("Expected oauth2 in gitlab url, got %s", glUrl)
	}

	// Custom user token URL
	bbUrl := buildAuthenticatedUrl("https://bitbucket.org/org/repo.git", "pass123", "bbuser")
	if !strings.Contains(bbUrl, "bbuser:pass123@bitbucket.org") {
		t.Errorf("Expected bbuser in bitbucket url, got %s", bbUrl)
	}

	// Masking
	masked := maskSecretUrl(ghUrl)
	if strings.Contains(masked, "ghp_12345") {
		t.Errorf("Expected token to be masked, got %s", masked)
	}
}

func TestSetupUserSshKey_CrlfNormalization(t *testing.T) {
	tmpHome := filepath.Join(os.TempDir(), "kp_test_ssh_user_home")
	defer os.RemoveAll(tmpHome)

	crlfKey := "-----BEGIN OPENSSH PRIVATE KEY-----\r\nb3BlbnNzaC1rZXktdjEAAAAA\r\n-----END OPENSSH PRIVATE KEY-----\r\n"
	err := setupUserSshKey("kp_user", tmpHome, crlfKey)
	if err != nil {
		t.Fatalf("setupUserSshKey failed: %v", err)
	}

	keyPath := filepath.Join(tmpHome, ".ssh", "id_rsa")
	content, err := os.ReadFile(keyPath)
	if err != nil {
		t.Fatalf("Failed to read id_rsa: %v", err)
	}

	if strings.Contains(string(content), "\r") {
		t.Errorf("Found carriage return in written key file: %q", string(content))
	}

	if !strings.HasSuffix(string(content), "\n") {
		t.Errorf("Key file should end with a newline")
	}

	// Verify id_ed25519 was also created
	edPath := filepath.Join(tmpHome, ".ssh", "id_ed25519")
	if _, err := os.Stat(edPath); os.IsNotExist(err) {
		t.Errorf("id_ed25519 was not created")
	}
}

func TestGenerateDeployKey(t *testing.T) {
	runner := NewRunner(true)
	pub, priv, err := runner.GenerateDeployKey("ed25519")
	if err != nil {
		t.Fatalf("GenerateDeployKey failed: %v", err)
	}

	if !strings.HasPrefix(pub, "ssh-ed25519 ") {
		t.Errorf("Expected ssh-ed25519 prefix in public key, got %s", pub)
	}

	if !strings.Contains(priv, "PRIVATE KEY") {
		t.Errorf("Expected PRIVATE KEY in private key, got %s", priv)
	}
}

