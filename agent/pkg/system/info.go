package system

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var startTime = time.Now()

type SystemInfo struct {
	Hostname      string `json:"hostname"`
	OS            string `json:"os"`
	Kernel        string `json:"kernel"`
	Architecture  string `json:"architecture"`
	PublicIP      string `json:"public_ip"`
	UptimeSeconds int64  `json:"uptime_seconds"`
}

func GetSystemInfo() SystemInfo {
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "unknown-host"
	}

	osName := "Ubuntu 24.04 LTS"
	if runtime.GOOS == "darwin" {
		osName = "macOS Development Host"
	} else {
		// Read /etc/os-release on Linux
		if data, err := os.ReadFile("/etc/os-release"); err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "PRETTY_NAME=") {
					osName = strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), `"`)
					break
				}
			}
		}
	}

	kernel := runtime.GOOS + " " + runtime.GOARCH
	if out, err := exec.Command("uname", "-r").Output(); err == nil {
		kernel = strings.TrimSpace(string(out))
	}

	uptime := int64(time.Since(startTime).Seconds())
	if data, err := os.ReadFile("/proc/uptime"); err == nil {
		parts := strings.Fields(string(data))
		if len(parts) > 0 {
			var up float64
			if _, err := parseNumber(parts[0], &up); err == nil {
				uptime = int64(up)
			}
		}
	}

	return SystemInfo{
		Hostname:      hostname,
		OS:            osName,
		Kernel:        kernel,
		Architecture:  runtime.GOARCH,
		PublicIP:      "127.0.0.1",
		UptimeSeconds: uptime,
	}
}

func parseNumber(s string, out *float64) (bool, error) {
	var val float64
	_, err := strings.NewReader(s).Read(nil)
	if err != nil {
		return false, err
	}
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			val = val*10 + float64(ch-'0')
		} else if ch == '.' {
			break
		}
	}
	*out = val
	return true, nil
}
