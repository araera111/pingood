package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Config はS3アップロードの設定を保持します
type S3Config struct {
	Region         string `toml:"region"`
	Bucket         string `toml:"bucket"`
	KeyPrefix      string `toml:"key_prefix"`
	AccessKey      string `toml:"access_key"`       // AWS Access Key
	SecretKey      string `toml:"secret_key"`       // AWS Secret Key
	Endpoint       string `toml:"endpoint"`         // カスタムエンドポイント（オプション）
	ForcePathStyle bool   `toml:"force_path_style"` // パススタイルアクセスを強制
	TLS            *bool  `toml:"tls"`              // TLS使用の有無（nilの場合はデフォルト）
	Schedule       string `toml:"schedule"`         // cron式でのスケジュール
	UploadTime     string `toml:"upload_time"`      // HH:MM形式（後方互換性用）
	DeleteAfter    bool   `toml:"delete_after"`
}

// S3Uploader はS3へのアップロード機能を提供します
type S3Uploader struct {
	client *s3.Client
	config S3Config
}

// NewS3Uploader は新しいS3Uploaderインスタンスを作成します
func NewS3Uploader(cfg S3Config) (*S3Uploader, error) {
	// 認証情報を設定
	creds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		cfg.AccessKey,
		cfg.SecretKey,
		"",
	))

	// AWS設定のオプションを準備
	opts := []func(*config.LoadOptions) error{
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(creds),
	}

	// カスタムエンドポイントが指定されている場合は追加
	if cfg.Endpoint != "" {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			// TLSの設定（デフォルトはtrue）
			useTLS := true
			if cfg.TLS != nil {
				useTLS = *cfg.TLS
			}

			// エンドポイントのURLを解析
			endpoint := cfg.Endpoint
			if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
				// プロトコルが指定されていない場合、TLS設定に基づいて追加
				if useTLS {
					endpoint = "https://" + endpoint
				} else {
					endpoint = "http://" + endpoint
				}
			}

			return aws.Endpoint{
				URL:               endpoint,
				SigningRegion:     cfg.Region,
				Source:            aws.EndpointSourceCustom,
				HostnameImmutable: true,
			}, nil
		})
		opts = append(opts, config.WithEndpointResolverWithOptions(customResolver))
	}

	// AWS設定を作成
	awsCfg, err := config.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("AWS設定の読み込みに失敗しました: %v", err)
	}

	// S3クライアントを初期化（パススタイル設定を含む）
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = cfg.ForcePathStyle
	})
	return &S3Uploader{
		client: client,
		config: cfg,
	}, nil
}

// UploadFile は指定されたファイルをS3にアップロードします
func (u *S3Uploader) UploadFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("ファイルのオープンに失敗しました: %v", err)
	}
	defer file.Close()

	// S3のキーを生成（プレフィックス + 日付 + ファイル名_タイムスタンプ）
	now := time.Now()
	baseFileName := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	ext := filepath.Ext(filePath)
	key := fmt.Sprintf("%s/%s/%s_%s%s",
		u.config.KeyPrefix,
		now.Format("2006/01/02"),
		baseFileName,
		now.Format("2006_01_02_15_04_05"),
		ext,
	)

	// S3にアップロード
	_, err = u.client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: &u.config.Bucket,
		Key:    &key,
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("S3へのアップロードに失敗しました: %v", err)
	}

	// アップロード後に削除が指定されている場合
	if u.config.DeleteAfter {
		if err := os.Remove(filePath); err != nil {
			return fmt.Errorf("ファイルの削除に失敗しました: %v", err)
		}
	}
	return nil
}

// ParseUploadTime は設定された時刻をパースします
func (u *S3Uploader) ParseUploadTime() (time.Time, error) {
	now := time.Now()
	timeStr := u.config.UploadTime

	// HH:MM形式をパース
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("時刻のパースに失敗しました: %v", err)
	}

	// 現在の日付と組み合わせる
	uploadTime := time.Date(
		now.Year(), now.Month(), now.Day(),
		t.Hour(), t.Minute(), 0, 0,
		now.Location(),
	)

	// 指定時刻が現在時刻より前の場合、翌日の同時刻とする
	if uploadTime.Before(now) {
		uploadTime = uploadTime.Add(24 * time.Hour)
	}

	return uploadTime, nil
}
