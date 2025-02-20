package logger

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestLogError(t *testing.T) {
	// テスト用の設定
	config := &Config{
		ErrorLogMode: "both",
	}

	// テスト用のファイルを作成
	logFilePath := "test.log"
	errorLogFilePath := "test.error.log"

	// ファイルが存在する場合は削除
	os.Remove(logFilePath)
	os.Remove(errorLogFilePath)

	// Loggerを初期化
	logger := &Logger{
		files:      []*os.File{},
		errorFiles: map[string]*os.File{},
		paths:      []string{logFilePath},
		config:     config,
	}

	// ファイルを作成
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("ログファイルのオープンに失敗しました: %v", err)
	}
	logger.files = append(logger.files, file)

	errorFile, err := os.OpenFile(errorLogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("エラーログファイルのオープンに失敗しました: %v", err)
	}
	logger.errorFiles[logFilePath] = errorFile

	// エラーをログに記録
	target := "example.com"
	testErr := fmt.Errorf("test error")
	err = logger.LogError(0, target, testErr)
	if err != nil {
		t.Fatalf("LogErrorに失敗しました: %v", err)
	}

	// ログファイルの内容を読み込む
	logFileContent, err := os.ReadFile(logFilePath)
	if err != nil {
		t.Fatalf("ログファイルの読み込みに失敗しました: %v", err)
	}

	// エラーログファイルの内容を読み込む
	errorLogFileContent, err := os.ReadFile(errorLogFilePath)
	if err != nil {
		t.Fatalf("エラーログファイルの読み込みに失敗しました: %v", err)
	}

	// テスト結果を検証
	logLine := fmt.Sprintf("[%s] ERROR - Target: %s, Error: %v\n",
		time.Now().Format("2006-01-02 15:04:05"),
		target,
		testErr)

	if string(logFileContent) != logLine {
		t.Errorf("ログファイルの内容が期待値と異なります: expected %q, got %q", logLine, string(logFileContent))
	}

	if string(errorLogFileContent) != logLine {
		t.Errorf("エラーログファイルの内容が期待値と異なります: expected %q, got %q", logLine, string(errorLogFileContent))
	}

	// テスト用のファイルを削除
	os.Remove(logFilePath)
	os.Remove(errorLogFilePath)
}

func TestLogErrorModeSame(t *testing.T) {
	// テスト用の設定
	config := &Config{
		ErrorLogMode: "same",
	}

	// テスト用のファイルを作成
	logFilePath := "test.log"
	errorLogFilePath := "test.error.log"

	// ファイルが存在する場合は削除
	os.Remove(logFilePath)
	os.Remove(errorLogFilePath)

	// Loggerを初期化
	logger := &Logger{
		files:      []*os.File{},
		errorFiles: map[string]*os.File{},
		paths:      []string{logFilePath},
		config:     config,
	}

	// ファイルを作成
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("ログファイルのオープンに失敗しました: %v", err)
	}
	logger.files = append(logger.files, file)

	errorFile, err := os.OpenFile(errorLogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("エラーログファイルのオープンに失敗しました: %v", err)
	}
	logger.errorFiles[logFilePath] = errorFile

	// エラーをログに記録
	target := "example.com"
	testErr := fmt.Errorf("test error")
	err = logger.LogError(0, target, testErr)
	if err != nil {
		t.Fatalf("LogErrorに失敗しました: %v", err)
	}

	// ログファイルの内容を読み込む
	logFileContent, err := os.ReadFile(logFilePath)
	if err != nil {
		t.Fatalf("ログファイルの読み込みに失敗しました: %v", err)
	}

	// エラーログファイルの内容を読み込む
	errorLogFileContent, err := os.ReadFile(errorLogFilePath)
	if err != nil {
		t.Fatalf("エラーログファイルの読み込みに失敗しました: %v", err)
	}

	// テスト結果を検証
	logLine := fmt.Sprintf("[%s] ERROR - Target: %s, Error: %v\n",
		time.Now().Format("2006-01-02 15:04:05"),
		target,
		testErr)

	if string(logFileContent) != logLine {
		t.Errorf("ログファイルの内容が期待値と異なります: expected %q, got %q", logLine, string(logFileContent))
	}

	if string(errorLogFileContent) != "" {
		t.Errorf("エラーログファイルの内容が期待値と異なります: expected %q, got %q", "", string(errorLogFileContent))
	}

	// テスト用のファイルを削除
	os.Remove(logFilePath)
	os.Remove(errorLogFilePath)
}

func TestLogErrorModeError(t *testing.T) {
	// テスト用の設定
	config := &Config{
		ErrorLogMode: "error",
	}

	// テスト用のファイルを作成
	logFilePath := "test.log"
	errorLogFilePath := "test.error.log"

	// ファイルが存在する場合は削除
	os.Remove(logFilePath)
	os.Remove(errorLogFilePath)

	// Loggerを初期化
	logger := &Logger{
		files:      []*os.File{},
		errorFiles: map[string]*os.File{},
		paths:      []string{logFilePath},
		config:     config,
	}

	// ファイルを作成
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("ログファイルのオープンに失敗しました: %v", err)
	}
	logger.files = append(logger.files, file)

	errorFile, err := os.OpenFile(errorLogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("エラーログファイルのオープンに失敗しました: %v", err)
	}
	logger.errorFiles[logFilePath] = errorFile

	// エラーをログに記録
	target := "example.com"
	testErr := fmt.Errorf("test error")
	err = logger.LogError(0, target, testErr)
	if err != nil {
		t.Fatalf("LogErrorに失敗しました: %v", err)
	}

	// ログファイルの内容を読み込む
	logFileContent, err := os.ReadFile(logFilePath)
	if err != nil {
		t.Fatalf("ログファイルの読み込みに失敗しました: %v", err)
	}

	// エラーログファイルの内容を読み込む
	errorLogFileContent, err := os.ReadFile(errorLogFilePath)
	if err != nil {
		t.Fatalf("エラーログファイルの読み込みに失敗しました: %v", err)
	}

	// テスト結果を検証
	logLine := fmt.Sprintf("[%s] ERROR - Target: %s, Error: %v\n",
		time.Now().Format("2006-01-02 15:04:05"),
		target,
		testErr)

	if string(logFileContent) != "" {
		t.Errorf("ログファイルの内容が期待値と異なります: expected %q, got %q", "", string(logFileContent))
	}

	if string(errorLogFileContent) != logLine {
		t.Errorf("エラーログファイルの内容が期待値と異なります: expected %q, got %q", logLine, string(errorLogFileContent))
	}

	// テスト用のファイルを削除
	os.Remove(logFilePath)
	os.Remove(errorLogFilePath)
}

func TestLogErrorModeOther(t *testing.T) {
	// テスト用の設定
	config := &Config{
		ErrorLogMode: "other",
	}

	// テスト用のファイルを作成
	logFilePath := "test.log"
	errorLogFilePath := "test.error.log"

	// ファイルが存在する場合は削除
	os.Remove(logFilePath)
	os.Remove(errorLogFilePath)

	// Loggerを初期化
	logger := &Logger{
		files:      []*os.File{},
		errorFiles: map[string]*os.File{},
		paths:      []string{logFilePath},
		config:     config,
	}

	// ファイルを作成
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("ログファイルのオープンに失敗しました: %v", err)
	}
	logger.files = append(logger.files, file)

	errorFile, err := os.OpenFile(errorLogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("エラーログファイルのオープンに失敗しました: %v", err)
	}
	logger.errorFiles[logFilePath] = errorFile

	// エラーをログに記録
	target := "example.com"
	testErr := fmt.Errorf("test error")
	err = logger.LogError(0, target, testErr)
	if err != nil {
		t.Fatalf("LogErrorに失敗しました: %v", err)
	}

	// ログファイルの内容を読み込む
	logFileContent, err := os.ReadFile(logFilePath)
	if err != nil {
		t.Fatalf("ログファイルの読み込みに失敗しました: %v", err)
	}

	// エラーログファイルの内容を読み込む
	errorLogFileContent, err := os.ReadFile(errorLogFilePath)
	if err != nil {
		t.Fatalf("エラーログファイルの読み込みに失敗しました: %v", err)
	}

	// テスト結果を検証
	logLine := fmt.Sprintf("[%s] ERROR - Target: %s, Error: %v\n",
		time.Now().Format("2006-01-02 15:04:05"),
		target,
		testErr)

	if string(logFileContent) != logLine {
		t.Errorf("ログファイルの内容が期待値と異なります: expected %q, got %q", logLine, string(logFileContent))
	}

	if string(errorLogFileContent) != logLine {
		t.Errorf("エラーログファイルの内容が期待値と異なります: expected %q, got %q", logLine, string(errorLogFileContent))
	}

	// テスト用のファイルを削除
	os.Remove(logFilePath)
	os.Remove(errorLogFilePath)
}

func TestGetErrorLogFilePath(t *testing.T) {
	logFilePath := "test.log"
	expectedErrorLogFilePath := "test.error.log"
	errorLogFilePath := getErrorLogFilePath(logFilePath)
	if errorLogFilePath != expectedErrorLogFilePath {
		t.Errorf("エラーログファイルパスが期待値と異なります: expected %q, got %q", expectedErrorLogFilePath, errorLogFilePath)
	}

	logFilePath = "path/to/test.log"
	expectedErrorLogFilePath = "path/to/test.error.log"
	errorLogFilePath = getErrorLogFilePath(logFilePath)
	if errorLogFilePath != expectedErrorLogFilePath {
		t.Errorf("エラーログファイルパスが期待値と異なります: expected %q, got %q", expectedErrorLogFilePath, errorLogFilePath)
	}
}
