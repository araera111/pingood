package ping

import (
	"os/exec"
	"runtime"
	"strings"
)

// ExtractHostFromURL はURLからホスト部分を抽出します
func ExtractHostFromURL(target string) string {
	if !strings.Contains(target, "://") {
		return target
	}
	host := strings.Split(strings.Split(target, "://")[1], "/")[0]
	return host
}

// CreatePingCommand はOSに応じたping commandを生成します
func CreatePingCommand(target string) (*exec.Cmd, error) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("ping", "-n", "1", "-w", "1000", target)
	case "linux":
		cmd = exec.Command("ping", "-c", "1", "-W", "1", target)
	case "darwin":
		cmd = exec.Command("ping", "-c", "1", "-W", "1000", target)
	default:
		return nil, ErrUnsupportedOS
	}
	return cmd, nil
}
