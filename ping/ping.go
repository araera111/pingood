package ping

import (
"net"
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
// URLからホスト名を抽出
target = ExtractHostFromURL(target)

// DNSルックアップ
if _, err := net.LookupHost(target); err != nil {
return nil, err
}

// ping commandの生成
cmd, err := CreatePingCommand(target)
if err != nil {
return nil, err
}

start := time.Now()
if err := cmd.Run(); err != nil {
return nil, err
}
rtt := time.Since(start)

return &PingResult{
Target:    target,
RTT:       rtt,
Timestamp: time.Now(),
}, nil
}
