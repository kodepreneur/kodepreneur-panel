package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	Server struct {
		ListenAddress string `json:"listen_address" yaml:"listen_address"`
		SocketPath    string `json:"socket_path" yaml:"socket_path"`
		UseSocket     bool   `json:"use_socket" yaml:"use_socket"`
		Port          int    `json:"port" yaml:"port"`
	} `json:"server" yaml:"server"`
	Security struct {
		SecretKey       string `json:"secret_key" yaml:"secret_key"`
		AllowedDriftSec int    `json:"allowed_drift_sec" yaml:"allowed_drift_sec"`
	} `json:"security" yaml:"security"`
	Environment struct {
		IsDev bool `json:"is_dev" yaml:"is_dev"`
	} `json:"environment" yaml:"environment"`
}

func DefaultConfig() *Config {
	cfg := &Config{}
	cfg.Server.ListenAddress = "127.0.0.1"
	cfg.Server.Port = 8765
	cfg.Server.SocketPath = "/run/kodepreneur/agent.sock"
	cfg.Server.UseSocket = false
	cfg.Security.SecretKey = "kodepreneur-default-insecure-dev-secret-change-me"
	cfg.Security.AllowedDriftSec = 30
	cfg.Environment.IsDev = false
	return cfg
}

func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()

	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			if !os.IsNotExist(err) {
				return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
			}
		} else {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					val := strings.TrimSpace(parts[1])
					val = strings.Trim(val, `"'`)

					switch key {
					case "listen_address", "host":
						cfg.Server.ListenAddress = val
					case "port":
						if p, err := strconv.Atoi(val); err == nil {
							cfg.Server.Port = p
						}
					case "socket_path":
						cfg.Server.SocketPath = val
					case "use_socket":
						cfg.Server.UseSocket = (val == "true" || val == "1" || val == "yes")
					case "secret_key":
						cfg.Security.SecretKey = val
					case "allowed_drift_sec":
						if d, err := strconv.Atoi(val); err == nil {
							cfg.Security.AllowedDriftSec = d
						}
					case "is_dev":
						cfg.Environment.IsDev = (val == "true" || val == "1" || val == "yes")
					}
				}
			}
		}
	}

	// Environment variable overrides
	if envSecret := os.Getenv("KODEPRENEUR_AGENT_SECRET"); envSecret != "" {
		cfg.Security.SecretKey = envSecret
	}
	if envPort := os.Getenv("KODEPRENEUR_AGENT_PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			cfg.Server.Port = p
		}
	}
	if envDev := os.Getenv("KODEPRENEUR_AGENT_DEV"); envDev != "" {
		cfg.Environment.IsDev = (envDev == "true" || envDev == "1")
	}

	return cfg, nil
}

func EnsureDirectory(path string) error {
	dir := filepath.Dir(path)
	return os.MkdirAll(dir, 0755)
}
