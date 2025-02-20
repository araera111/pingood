package path

import (
"reflect"
"testing"
)

func TestSanitizeTargetForFilename(t *testing.T) {
tests := []struct {
name   string
target string
want   string
}{
{"Simple hostname", "example.com", "example.com.log"},
{"URL with path", "http://example.com/path", "http:--example.com-path.log"},
{"IP address", "192.168.1.1", "192.168.1.1.log"},
{"Complex URL", "https://api.example.com/v1/endpoint", "https:--api.example.com-v1-endpoint.log"},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
got := SanitizeTargetForFilename(tt.target)
if got != tt.want {
t.Errorf("SanitizeTargetForFilename(%q) = %v, want %v", tt.target, got, tt.want)
}
})
}
}

func TestSanitizePaths(t *testing.T) {
tests := []struct {
name  string
paths []string
want  []string
}{
{
name:  "Paths with spaces",
paths: []string{" path1 ", "path2  ", "  path3"},
want:  []string{"path1", "path2", "path3"},
},
{
name:  "Empty paths",
paths: []string{"", " ", "  "},
want:  []string{"", "", ""},
},
{
name:  "Mixed paths",
paths: []string{"path1", " path2 ", "path3"},
want:  []string{"path1", "path2", "path3"},
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
got := SanitizePaths(tt.paths)
if !reflect.DeepEqual(got, tt.want) {
t.Errorf("SanitizePaths() = %v, want %v", got, tt.want)
}
})
}
}

func TestGenerateLogPaths(t *testing.T) {
tests := []struct {
name    string
targets []string
want    []string
}{
{
name:    "Simple targets",
targets: []string{"example.com", "test.com"},
want:    []string{"example.com.log", "test.com.log"},
},
{
name:    "URLs with paths",
targets: []string{"http://example.com/path", "https://test.com/api"},
want:    []string{"http:--example.com-path.log", "https:--test.com-api.log"},
},
{
name:    "Mixed targets",
targets: []string{"192.168.1.1", "example.com", "https://test.com/api"},
want:    []string{"192.168.1.1.log", "example.com.log", "https:--test.com-api.log"},
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
got := GenerateLogPaths(tt.targets)
if !reflect.DeepEqual(got, tt.want) {
t.Errorf("GenerateLogPaths() = %v, want %v", got, tt.want)
}
})
}
}