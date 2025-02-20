package ping

import (
"runtime"
"testing"
)

func TestExtractHostFromURL(t *testing.T) {
tests := []struct {
name   string
target string
want   string
}{
{"Simple URL", "http://example.com", "example.com"},
{"URL with path", "https://example.com/path", "example.com"},
{"URL with port", "http://example.com:8080", "example.com:8080"},
{"Non-URL host", "example.com", "example.com"},
{"IP address", "192.168.1.1", "192.168.1.1"},
{"Complex URL", "https://api.example.com/v1/test", "api.example.com"},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
got := ExtractHostFromURL(tt.target)
if got != tt.want {
t.Errorf("ExtractHostFromURL(%q) = %v, want %v", tt.target, got, tt.want)
}
})
}
}

func TestCreatePingCommand(t *testing.T) {
tests := []struct {
name    string
target  string
wantErr bool
}{
{"Valid target", "example.com", false},
{"IP address", "192.168.1.1", false},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
cmd, err := CreatePingCommand(tt.target)
if (err != nil) != tt.wantErr {
t.Errorf("CreatePingCommand() error = %v, wantErr %v", err, tt.wantErr)
return
}
if !tt.wantErr {
if cmd == nil {
t.Error("CreatePingCommand() returned nil command")
return
}

var expectedArgs []string
switch runtime.GOOS {
case "windows":
expectedArgs = []string{"ping", "-n", "1", "-w", "1000", tt.target}
case "linux":
expectedArgs = []string{"ping", "-c", "1", "-W", "1", tt.target}
case "darwin":
expectedArgs = []string{"ping", "-c", "1", "-W", "1000", tt.target}
}

if len(cmd.Args) != len(expectedArgs) {
t.Errorf("CreatePingCommand() args count = %v, want %v", len(cmd.Args), len(expectedArgs))
return
}

for i, arg := range cmd.Args {
if arg != expectedArgs[i] {
t.Errorf("CreatePingCommand() arg[%d] = %v, want %v", i, arg, expectedArgs[i])
}
}
}
})
}
}