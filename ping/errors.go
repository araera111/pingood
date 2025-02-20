package ping

import (
	"fmt"
	"runtime"
)

var (
	// ErrUnsupportedOS は未サポートのOSでpingを実行しようとした場合のエラーです
	ErrUnsupportedOS = fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
)
