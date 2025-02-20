package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
	"path/filepath"

	"github.com/robfig/cron/v3"
	"pingood/ping"
)

// Logger handles logging of ping results with optional S3 uploads
type Logger struct {
	files       []*os.File
	errorFiles  map[string]*os.File // エラーログファイル
	paths       []string
	uploader    *S3Uploader // オプショナル
	config      *Config     // オプショナル
	cron        *cron.Cron  // オプショナル
	mu       sync.Mutex
}

// LoggerOptions はロガーの設定オプションを定義します
type LoggerOptions struct {
	ConfigPath     string // S3アップロード用の設定ファイルパス
	UploadExisting bool   // 起動時に既存のログファイルをアップロードするか
}

// NewLogger creates a new Logger instance
func NewLogger(paths []string, opts *LoggerOptions) (*Logger, error) {
	var files []*os.File
	var errorFiles = make(map[string]*os.File)
	var uploader *S3Uploader
	var config *Config
	var cronJob *cron.Cron

	// 複数のログファイルを開く
	for _, path := range paths {
		file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			// エラーが発生した場合、既に開いたファイルを全て閉じる
			for _, f := range files {
				f.Close()
			}
			return nil, fmt.Errorf("ログファイルのオープンに失敗しました %s: %v", path, err)
		}
		files = append(files, file)

		// エラーログファイルを開く
		errorFilePath := getErrorLogFilePath(path)
		errorFile, err := os.OpenFile(errorFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			for _, f := range files {
				f.Close()
			}
			return nil, fmt.Errorf("エラーログファイルのオープンに失敗しました %s: %v", errorFilePath, err)
		}
		errorFiles[path] = errorFile
	}

	// S3アップロード機能の初期化（オプショナル）
	if opts != nil && opts.ConfigPath != "" {
		var err error
		// 設定ファイルの読み込み
		config, err = LoadConfig(opts.ConfigPath)
		if err != nil {
			for _, f := range files {
				f.Close()
			}
			return nil, fmt.Errorf("設定の読み込みに失敗しました: %v", err)
		}

		// S3アップローダーの初期化
		uploader, err = NewS3Uploader(config.S3)
		if err != nil {
			for _, f := range files {
				f.Close()
			}
			return nil, fmt.Errorf("S3アップローダーの初期化に失敗しました: %v", err)
		}

		// cronの初期化
		cronJob = cron.New()

		// 既存ファイルのアップロード確認
		if opts.UploadExisting {
			fmt.Println("既存のログファイルをアップロードしています...")
			for _, path := range paths {
				if err := uploader.UploadFile(path); err != nil {
					fmt.Fprintf(os.Stderr, "既存ファイルのアップロードに失敗しました %s: %v\n", path, err)
				} else {
					fmt.Printf("アップロード完了: %s\n", path)
				}
			}
		}
	}

	l := &Logger{
		files:       files,
		errorFiles:  errorFiles,
		paths:       paths,
		uploader:    uploader,
		config:      config,
		cron:        cronJob,
	}

	// S3アップロード機能が有効な場合のみスケジュール設定
	if cronJob != nil {
		if err := l.scheduleUpload(); err != nil {
			l.Close()
			return nil, fmt.Errorf("アップロードスケジュールの設定に失敗しました: %v", err)
		}
		cronJob.Start()
	}

	return l, nil
}

// scheduleUpload configures the upload schedule
func (l *Logger) scheduleUpload() error {
	var schedule string

	// scheduleが設定されている場合はそちらを優先
	if l.config.S3.Schedule != "" {
		schedule = l.config.S3.Schedule
	} else {
		// 後方互換性のためにupload_timeを使用
		uploadTime, err := l.uploader.ParseUploadTime()
		if err != nil {
			return fmt.Errorf("スケジュール時刻の解析に失敗しました: %v", err)
		}
		schedule = fmt.Sprintf("%d %d * * *", uploadTime.Minute(), uploadTime.Hour())
	}

	// cron式の検証
	if _, err := cron.ParseStandard(schedule); err != nil {
		return fmt.Errorf("不正なcron式です: %v", err)
	}

	_, err := l.cron.AddFunc(schedule, func() {
		l.uploadLogs()
	})

	return err
}

// uploadLogs uploads all log files to S3
func (l *Logger) uploadLogs() {
 l.mu.Lock()
 defer l.mu.Unlock()

 for _, path := range l.paths {
  if err := l.uploader.UploadFile(path); err != nil {
   fmt.Fprintf(os.Stderr, "ログファイルのアップロードに失敗しました %s: %v\n", path, err)
  }

  // エラーログファイルのアップロード
  errorFilePath := getErrorLogFilePath(path)
  if _, err := os.Stat(errorFilePath); err == nil { // ファイルが存在する場合のみアップロード
   if err := l.uploader.UploadFile(errorFilePath); err != nil {
    fmt.Fprintf(os.Stderr, "エラーログファイルのアップロードに失敗しました %s: %v\n", errorFilePath, err)
   }
  }
 }
}

// UploadNow triggers an immediate upload of all log files to S3
func (l *Logger) UploadNow() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	var lastErr error
	for _, path := range l.paths {
		if err := l.uploader.UploadFile(path); err != nil {
			lastErr = err
			fmt.Fprintf(os.Stderr, "ログファイルのアップロードに失敗しました %s: %v\n", path, err)
		}
	}
	return lastErr
}

// Close stops the cron job and closes all log files
func (l *Logger) Close() error {
	// cronジョブを停止
	if l.cron != nil {
		l.cron.Stop()
	}

	var lastErr error
	for _, file := range l.files {
		if err := file.Close(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// ensureFileExists checks if the file exists and recreates it if necessary
func (l *Logger) ensureFileExists(index int) error {
	if _, err := os.Stat(l.paths[index]); os.IsNotExist(err) {
		// 既存のファイルハンドルを閉じる
		if l.files[index] != nil {
			l.files[index].Close()
		}

		// 新しいファイルを作成
		file, err := os.OpenFile(l.paths[index], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to recreate log file %s: %v", l.paths[index], err)
		}
		l.files[index] = file
	}
	return nil
}

// LogSuccess logs a successful ping result to the specified file index
func (l *Logger) LogSuccess(index int, target string, result *ping.PingResult) error {
	if index < 0 || index >= len(l.files) {
		return fmt.Errorf("invalid file index: %d", index)
	}

	// ファイルの存在を確認し、必要に応じて再作成
	if err := l.ensureFileExists(index); err != nil {
		return err
	}

	logLine := fmt.Sprintf("[%s] SUCCESS - Target: %s, RTT: %v\n",
		result.Timestamp.Format("2006-01-02 15:04:05"),
		target,
		result.RTT)

	_, err := l.files[index].WriteString(logLine)
	return err
}

func getErrorLogFilePath(logFilePath string) string {
	dir, file := filepath.Split(logFilePath)
	ext := filepath.Ext(file)
	base := file[:len(file)-len(ext)]
	return filepath.Join(dir, base+".error.log")
}

// LogError logs a failed ping attempt to the specified file index
func (l *Logger) LogError(index int, target string, err error) error {
	if index < 0 || index >= len(l.files) {
		return fmt.Errorf("invalid file index: %d", index)
	}

	// ファイルの存在を確認し、必要に応じて再作成
	if err := l.ensureFileExists(index); err != nil {
		return err
	}

	logLine := fmt.Sprintf("[%s] ERROR - Target: %s, Error: %v\n",
		time.Now().Format("2006-01-02 15:04:05"),
		target,
		err)

	errorLogMode := l.config.ErrorLogMode
	if errorLogMode == "" {
		errorLogMode = "both" // デフォルトはboth
	}

	switch errorLogMode {
	case "same":
	  _, err = l.files[index].WriteString(logLine)
	  if err != nil {
	   return err
	  }
	case "both":
	  _, err = l.files[index].WriteString(logLine)
	  if err != nil {
	   return err
	  }
	  _, err = l.errorFiles[l.paths[index]].WriteString(logLine)
		if err != nil {
			return err
		}
	case "error":
		_, err = l.errorFiles[l.paths[index]].WriteString(logLine)
		if err != nil {
			return err
		}
	default:
		// 不正な設定の場合は、両方のログに書き出す
		_, err = l.files[index].WriteString(logLine)
		if err != nil {
			return err
		}
		_, err = l.errorFiles[l.paths[index]].WriteString(logLine)
		if err != nil {
			return err
		}
	}

	return err
}
