package path

import (
"strings"
)

// SanitizeTargetForFilename はターゲットURLやIPアドレスからログファイル名を生成します
func SanitizeTargetForFilename(target string) string {
return strings.ReplaceAll(target, "/", "-") + ".log"
}

// SanitizePaths はパスのスライスから空白を除去します
func SanitizePaths(paths []string) []string {
sanitized := make([]string, len(paths))
for i, path := range paths {
sanitized[i] = strings.TrimSpace(path)
}
return sanitized
}

// GenerateLogPaths は複数のターゲットに対するログファイルパスを生成します
func GenerateLogPaths(targets []string) []string {
paths := make([]string, len(targets))
for i, target := range targets {
paths[i] = SanitizeTargetForFilename(target)
}
return paths
}