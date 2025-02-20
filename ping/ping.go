package ping

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// PingResult represents the result of a ping operation
type PingResult struct {
	Target    string
	RTT       time.Duration
	Timestamp time.Time
}

// Ping executes a ping command to the specified target and returns the result
func Ping(target string) (*PingResult, error) {
	// ターゲットの検証
	if strings.Contains(target, "://") {
		host := strings.Split(strings.Split(target, "://")[1], "/")[0]
		target = host
	}

	// DNSルックアップ
	if _, err := net.LookupHost(target); err != nil {
		return nil, fmt.Errorf("failed to resolve host: %v", err)
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("ping", "-n", "1", "-w", "1000", target)
	case "linux":
		cmd = exec.Command("ping", "-c", "1", "-W", "1", target)
	case "darwin":
		cmd = exec.Command("ping", "-c", "1", "-W", "1000", target)
	default:
		return nil, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	start := time.Now()
	err := cmd.Run()
	rtt := time.Since(start)

	if err != nil {
		return nil, fmt.Errorf("ping failed: %v", err)
	}

	return &PingResult{
		Target:    target,
		RTT:       rtt,
		Timestamp: time.Now(),
	}, nil
}
