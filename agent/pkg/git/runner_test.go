package git

import (
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
