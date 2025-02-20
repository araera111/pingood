package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

// Config はアプリケーション全体の設定を保持します
type Config struct {
	LogFiles      []string `toml:"log_files"`
	S3            S3Config `toml:"s3"`
	ErrorLogMode string   `toml:"error_log_mode"`
}

// LoadConfig は指定されたパスから設定を読み込みます
func LoadConfig(path string) (*Config, error) {
	var config Config

	// 設定ファイルが存在しない場合はデフォルト設定を返す
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Config{
			LogFiles: []string{"ping.log"},
			S3: S3Config{
				Region:      "ap-northeast-1",
				Bucket:      "default-bucket",
				KeyPrefix:   "logs",
				UploadTime:  "00:00",
				DeleteAfter: false,
			},
		}, nil
	}

	// TOMLファイルを読み込む
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, fmt.Errorf("設定ファイルの読み込みに失敗しました: %v", err)
	}

	// 設定の検証
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// validateConfig は設定内容を検証します
func validateConfig(config *Config) error {
	if len(config.LogFiles) == 0 {
		return fmt.Errorf("少なくとも1つのログファイルパスを指定してください")
	}

	// scheduleまたはupload_timeのどちらかが設定されている場合のみS3の設定を検証
	if config.S3.Schedule != "" || config.S3.UploadTime != "" {
		if config.S3.Bucket == "" {
			return fmt.Errorf("S3バケットを指定してください")
		}

		if config.S3.Region == "" {
			return fmt.Errorf("AWSリージョンを指定してください")
		}

		if config.S3.AccessKey == "" {
			return fmt.Errorf("AWS Access Keyを指定してください")
		}

		if config.S3.SecretKey == "" {
			return fmt.Errorf("AWS Secret Keyを指定してください")
		}

		// scheduleが設定されていない場合のみupload_timeを検証
		if config.S3.Schedule == "" && config.S3.UploadTime != "" {
			if _, err := time.Parse("15:04", config.S3.UploadTime); err != nil {
				return fmt.Errorf("アップロード時刻のフォーマットが不正です（HH:MM形式で指定してください）: %v", err)
			}
		}
	}

	return nil
}

// WriteDefaultConfig はデフォルトの設定ファイルを生成します
func WriteDefaultConfig(path string) error {
	config := Config{
		LogFiles: []string{"ping.log"},
		S3: S3Config{
			Region:      "ap-northeast-1",
			Bucket:      "your-bucket-name",
			KeyPrefix:   "logs",
			UploadTime:  "00:00",
			DeleteAfter: false,
		},
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("設定ファイルの作成に失敗しました: %v", err)
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("設定ファイルの書き込みに失敗しました: %v", err)
	}

	return nil
}
