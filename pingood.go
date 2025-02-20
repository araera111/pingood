package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"pingood/logger"
	"pingood/ping"
)

func readInput(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return scanner.Text()
	}
	return ""
}

func main() {
	// コマンドライン引数の定義
	target := flag.String("target", "", "Target URLs or IP addresses to ping (comma-separated)")
	interval := flag.Int("interval", 5, "Ping interval in seconds")
	logPath := flag.String("log", "", "Paths to log files (comma-separated)")
	upload := flag.Bool("upload", false, "Enable S3 upload with config.toml")
	configPath := flag.String("config", "config.toml", "Path to config.toml for S3 upload settings")
	flag.Parse()

	// 引数がない場合は対話的に入力を受け付ける
	if *target == "" {
		*target = readInput("対象のURLまたはIPアドレスを入力してください（カンマ区切りで複数指定可能）: ")
		if *target == "" {
			fmt.Println("Error: target is required")
			flag.Usage()
			os.Exit(1)
		}
	}

	// カンマ区切りの文字列をスライスに分割
	targets := strings.Split(*target, ",")
	for i := range targets {
		targets[i] = strings.TrimSpace(targets[i])
	}

	// intervalが未指定の場合
	if flag.Lookup("interval").Value.String() == "5" && len(os.Args) == 1 {
		intervalStr := readInput("Ping実行間隔を秒単位で入力してください（デフォルト: 5）: ")
		if intervalStr != "" {
			fmt.Sscanf(intervalStr, "%d", interval)
		}
	}

	// logPathが未指定の場合
	var logPaths []string
	if *logPath == "" {
		if flag.Lookup("log").Value.String() == "" {
			autoGenerate := readInput("ログファイル名を自動生成しますか？（y/n、デフォルト: y）: ")
			if autoGenerate == "" || strings.ToLower(autoGenerate) == "y" {
				// ターゲットごとにログファイル名を自動生成
				for _, t := range targets {
					logFileName := strings.ReplaceAll(t, "/", "-") + ".log"
					logPaths = append(logPaths, logFileName)
				}
			} else {
				*logPath = readInput("ログファイルのパスを入力してください（カンマ区切りで複数指定可能）: ")
				if *logPath == "" {
					fmt.Println("Error: log path is required")
					flag.Usage()
					os.Exit(1)
				}
				logPaths = strings.Split(*logPath, ",")
			}
		}
	} else {
		logPaths = strings.Split(*logPath, ",")
	}

	// 各要素からスペースを削除
	for i := range logPaths {
		logPaths[i] = strings.TrimSpace(logPaths[i])
	}

	// targetsとlogPathsの数が一致することを確認
	if len(targets) != len(logPaths) {
		log.Fatalf("Error: number of targets (%d) must match number of log files (%d)", len(targets), len(logPaths))
	}

	// ログファイルのディレクトリを作成
	for _, p := range logPaths {
		dir := filepath.Dir(p)
		if dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				log.Fatalf("Failed to create directory for log file: %v", err)
			}
		}
	}

	// ロガーの初期化
	var opts *logger.LoggerOptions
	if *upload {
		opts = &logger.LoggerOptions{
			ConfigPath: *configPath,
		}

		// 既存のログファイルをチェック
		for _, logPath := range logPaths {
			if info, err := os.Stat(logPath); err == nil && info.Size() > 0 {
				fmt.Printf("既存のログファイル '%s' が見つかりました（サイズ: %d bytes）\n", logPath, info.Size())
				fmt.Print("このログファイルをアップロードしますか？ [y/N]: ")
				scanner := bufio.NewScanner(os.Stdin)
				if scanner.Scan() && strings.ToLower(scanner.Text()) == "y" {
					opts.UploadExisting = true
					break
				}
			}
		}
	}

	l, err := logger.NewLogger(logPaths, opts)
	if err != nil {
		log.Fatalf("ロガーの初期化に失敗しました: %v", err)
	}
	defer l.Close()

	if *upload {
		log.Printf("S3アップロードが有効です（設定ファイル: %s）\n", *configPath)
	}

	// メインループ
	ticker := time.NewTicker(time.Duration(*interval) * time.Second)
	defer ticker.Stop()

	fmt.Printf("Starting ping to %s (interval: %d seconds)\n", strings.Join(targets, ", "), *interval)
	fmt.Printf("Logging to: %s\n", strings.Join(logPaths, ", "))

	for range ticker.C {
		for i, t := range targets {
			result, err := ping.Ping(t)
			if err != nil {
				l.LogError(i, t, err)
				continue
			}
			l.LogSuccess(i, t, result)
		}
	}
}
