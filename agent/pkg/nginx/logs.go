package nginx

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type LogManager struct {
	logDir string
	isDev  bool
}

func NewLogManager(isDev bool) *LogManager {
	if isDev || runtime.GOOS != "linux" {
		mockLogDir := filepath.Join(os.TempDir(), "kodepreneur", "logs")
		_ = os.MkdirAll(mockLogDir, 0755)
		return &LogManager{
			logDir: mockLogDir,
			isDev:  true,
		}
	}

	return &LogManager{
		logDir: "/var/log/nginx",
		isDev:  false,
	}
}

// ReadLogLines returns the last N lines from the given domain's access or error log.
func (lm *LogManager) ReadLogLines(domain, logType string, lineCount int) ([]string, error) {
	if lineCount <= 0 {
		lineCount = 100
	}
	if lineCount > 1000 {
		lineCount = 1000
	}

	logFile := filepath.Join(lm.logDir, fmt.Sprintf("%s.%s.log", domain, logType))

	if lm.isDev || runtime.GOOS != "linux" {
		// If file doesn't exist in dev, generate sample log entries
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			return lm.generateSampleLogs(domain, logType, lineCount), nil
		}
	}

	file, err := os.Open(logFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{fmt.Sprintf("[%s] No %s log entries recorded yet for %s", time.Now().Format(time.RFC3339), logType, domain)}, nil
		}
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading log file: %w", err)
	}

	// Keep only the last `lineCount` entries
	if len(lines) > lineCount {
		lines = lines[len(lines)-lineCount:]
	}

	return lines, nil
}

func (lm *LogManager) generateSampleLogs(domain, logType string, count int) []string {
	now := time.Now()
	var logs []string
	if logType == "access" {
		paths := []string{"/", "/api/v1/status", "/assets/app.js", "/assets/app.css", "/login"}
		for i := 0; i < 15; i++ {
			t := now.Add(time.Duration(-i*30) * time.Second).Format("02/Jan/2006:15:04:05 -0000")
			p := paths[i%len(paths)]
			logs = append(logs, fmt.Sprintf(`127.0.0.1 - - [%s] "GET %s HTTP/1.1" 200 4520 "https://%s" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)"`, t, p, domain))
		}
	} else {
		logs = append(logs,
			fmt.Sprintf(`%s [notice] 1420#1420: using the "epoll" event method for %s`, now.Format("2006/01/02 15:04:05"), domain),
			fmt.Sprintf(`%s [notice] 1420#1420: start worker processes`, now.Format("2006/01/02 15:04:05")),
			fmt.Sprintf(`%s [info] 1421#1421: client 127.0.0.1 connected to %s`, now.Format("2006/01/02 15:04:05"), domain),
		)
	}
	return logs
}
