package system

import (
	"bufio"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type CPUStats struct {
	UsagePercent float64 `json:"usage_percent"`
	Cores        int     `json:"cores"`
}

type MemoryStats struct {
	TotalBytes   uint64  `json:"total_bytes"`
	UsedBytes    uint64  `json:"used_bytes"`
	FreeBytes    uint64  `json:"free_bytes"`
	UsagePercent float64 `json:"usage_percent"`
}

type DiskStats struct {
	Mount        string  `json:"mount"`
	TotalBytes   uint64  `json:"total_bytes"`
	UsedBytes    uint64  `json:"used_bytes"`
	FreeBytes    uint64  `json:"free_bytes"`
	UsagePercent float64 `json:"usage_percent"`
}

type LoadAvgStats struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

type SystemMetrics struct {
	CPU         CPUStats     `json:"cpu"`
	Memory      MemoryStats  `json:"memory"`
	Disk        DiskStats    `json:"disk"`
	LoadAverage LoadAvgStats `json:"load_average"`
	Timestamp   int64        `json:"timestamp"`
}

var (
	prevCPUTotal uint64
	prevCPUIdle  uint64
	cpuMu        sync.Mutex
)

func GetSystemMetrics() SystemMetrics {
	return SystemMetrics{
		CPU:         getCPUStats(),
		Memory:      getMemoryStats(),
		Disk:        getDiskStats("/"),
		LoadAverage: getLoadAverage(),
		Timestamp:   time.Now().Unix(),
	}
}

func getCPUStats() CPUStats {
	cores := runtime.NumCPU()
	usage := 0.0

	if runtime.GOOS == "linux" {
		usage = readLinuxCPUUsage()
	} else {
		// Mock CPU calculation for macOS/Dev
		usage = 12.5 + math.Sin(float64(time.Now().Second()))*5.0
	}

	return CPUStats{
		UsagePercent: math.Round(usage*100) / 100,
		Cores:        cores,
	}
}

func readLinuxCPUUsage() float64 {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return 0.0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu ") {
			fields := strings.Fields(line)
			if len(fields) < 5 {
				return 0.0
			}

			var total, idle uint64
			for i := 1; i < len(fields); i++ {
				val, _ := strconv.ParseUint(fields[i], 10, 64)
				total += val
				if i == 4 { // idle is 4th column in /proc/stat
					idle = val
				}
			}

			cpuMu.Lock()
			defer cpuMu.Unlock()

			if prevCPUTotal == 0 {
				prevCPUTotal = total
				prevCPUIdle = idle
				return 0.0
			}

			totalDiff := total - prevCPUTotal
			idleDiff := idle - prevCPUIdle

			prevCPUTotal = total
			prevCPUIdle = idle

			if totalDiff == 0 {
				return 0.0
			}

			usage := float64(totalDiff-idleDiff) / float64(totalDiff) * 100.0
			if usage < 0 {
				usage = 0
			}
			if usage > 100 {
				usage = 100
			}
			return usage
		}
	}
	return 0.0
}

func getMemoryStats() MemoryStats {
	if runtime.GOOS == "linux" {
		file, err := os.Open("/proc/meminfo")
		if err == nil {
			defer file.Close()
			var memTotal, memAvailable uint64
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "MemTotal:") {
					memTotal = parseMeminfoLine(line) * 1024
				} else if strings.HasPrefix(line, "MemAvailable:") {
					memAvailable = parseMeminfoLine(line) * 1024
				}
			}

			if memTotal > 0 {
				used := memTotal - memAvailable
				usagePct := float64(used) / float64(memTotal) * 100.0
				return MemoryStats{
					TotalBytes:   memTotal,
					UsedBytes:    used,
					FreeBytes:    memAvailable,
					UsagePercent: math.Round(usagePct*100) / 100,
				}
			}
		}
	}

	// Dev / Fallback calculation
	total := uint64(2 * 1024 * 1024 * 1024) // 2GB
	used := uint64(800 * 1024 * 1024)       // 800MB
	free := total - used
	return MemoryStats{
		TotalBytes:   total,
		UsedBytes:    used,
		FreeBytes:    free,
		UsagePercent: math.Round(float64(used)/float64(total)*10000) / 100,
	}
}

func parseMeminfoLine(line string) uint64 {
	fields := strings.Fields(line)
	if len(fields) >= 2 {
		val, _ := strconv.ParseUint(fields[1], 10, 64)
		return val
	}
	return 0
}

func getDiskStats(mount string) DiskStats {
	out, err := exec.Command("df", "-k", mount).Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(lines) >= 2 {
			fields := strings.Fields(lines[1])
			if len(fields) >= 4 {
				totalKB, _ := strconv.ParseUint(fields[1], 10, 64)
				usedKB, _ := strconv.ParseUint(fields[2], 10, 64)
				freeKB, _ := strconv.ParseUint(fields[3], 10, 64)

				total := totalKB * 1024
				used := usedKB * 1024
				free := freeKB * 1024
				usagePct := 0.0
				if total > 0 {
					usagePct = float64(used) / float64(total) * 100.0
				}

				return DiskStats{
					Mount:        mount,
					TotalBytes:   total,
					UsedBytes:    used,
					FreeBytes:    free,
					UsagePercent: math.Round(usagePct*100) / 100,
				}
			}
		}
	}

	// Fallback
	total := uint64(50 * 1024 * 1024 * 1024)
	used := uint64(12 * 1024 * 1024 * 1024)
	return DiskStats{
		Mount:        "/",
		TotalBytes:   total,
		UsedBytes:    used,
		FreeBytes:    total - used,
		UsagePercent: 24.0,
	}
}

func getLoadAverage() LoadAvgStats {
	if data, err := os.ReadFile("/proc/loadavg"); err == nil {
		fields := strings.Fields(string(data))
		if len(fields) >= 3 {
			l1, _ := strconv.ParseFloat(fields[0], 64)
			l5, _ := strconv.ParseFloat(fields[1], 64)
			l15, _ := strconv.ParseFloat(fields[2], 64)
			return LoadAvgStats{Load1: l1, Load5: l5, Load15: l15}
		}
	}

	return LoadAvgStats{
		Load1:  0.22,
		Load5:  0.18,
		Load15: 0.14,
	}
}
