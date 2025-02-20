# pingood

pingoodは、複数のホストに対して同時にpingを実行し、結果を個別のログファイルに記録するGoで実装されたモニタリングツールです。
つくった理由は、客先でyahooと自社のサービスにpingを行い、うちのサービスが落ちているのかどうかを測定し、証拠として出すため。

## 特徴

- 複数ホストの同時監視
- URLおよびホスト名のサポート
- クロスプラットフォーム対応（Windows、Linux、macOS）
- ホストごとの個別ログファイル
- ログファイルの自動再作成機能
- 完全対話型モード対応
- ログファイル名の自動生成機能
- S3/MinIOへの自動アップロード機能
  - 定期的なログのアップロード
  - 既存ログファイルの選択的アップロード
  - AWS S3およびMinIO互換ストレージ対応

## インストール

```bash
go install github.com/araera111/pingood@latest
```

## 使用方法

### 基本的な使用方法

```bash
# コマンドライン引数を使用する場合
pingood -target "example.com,google.com" -interval 5 -log "example.log,google.log"

# 対話型モードを使用する場合
pingood
対象のURLまたはIPアドレスを入力してください（カンマ区切りで複数指定可能）: example.com,google.com
Ping実行間隔を秒単位で入力してください（デフォルト: 5）: 10
ログファイル名を自動生成しますか？（y/n、デフォルト: y）: y
```

### S3/MinIOアップロード機能の使用

```bash
# S3アップロード機能を有効にして実行
pingood -target example.com -upload -config config.toml

# 設定ファイルを指定して実行
pingood -target example.com -upload -config custom_config.toml
```

アップロード機能を有効にすると、以下の機能が利用可能になります：
- 起動時に既存ログファイルのアップロード（確認プロンプトあり）
- 定期的な自動アップロード（cron式またはHH:MM形式で指定）
- AWS S3またはMinIO互換ストレージへのアップロード

### 設定ファイル (config.toml)

```toml
# ログファイルの設定
log_files = [
    "ping.log",
    "error.log"
]

# S3アップロードの設定
[s3]
# 認証情報
access_key = "YOUR_ACCESS_KEY"
secret_key = "YOUR_SECRET_KEY"

# MinIO設定例
endpoint = "ext-dev-minio.example.com"  # MinIOサーバーのエンドポイント
force_path_style = true                # MinIOの場合は必須
tls = false                           # HTTPSを使用しない場合はfalse

# S3の基本設定
region = "ap-northeast-1"
bucket = "your-bucket-name"
key_prefix = "logs/ping"              # アップロード先のプレフィックス

# スケジュール設定（いずれか一方）
schedule = "*/10 * * * *"             # cron式（10分おき）
# または
upload_time = "23:00"                 # 時刻指定（HH:MM形式）

# アップロード後の処理
delete_after = false                  # アップロード後にログファイル削除
```

### オプション

- `-target`: ping対象のURLまたはIPアドレス（カンマ区切りで複数指定可能）
- `-interval`: ping実行間隔（秒単位、デフォルト: 5秒）
- `-log`: ログファイルのパス（カンマ区切りで複数指定可能、targetと同じ数が必要）
- `-upload`: S3/MinIOアップロード機能を有効化
- `-config`: アップロード設定ファイルのパス（デフォルト: config.toml）

### ログ形式

成功時のログ形式：
```
[2025-02-19 18:14:27] SUCCESS - Target: example.com, RTT: 123.456ms
```

エラー時のログ形式：
```
[2025-02-19 18:14:27] ERROR - Target: example.com, Error: failed to resolve host
```

## S3アップロードパス形式

アップロードされるログファイルは以下の形式で保存されます：
```
{key_prefix}/{YYYY}/{MM}/{DD}/{filename}
```
例：`logs/ping/2025/02/20/example.log`

## 要件

- Go 1.16以上
- AWS S3またはMinIO互換ストレージ（アップロード機能を使用する場合）

## ライセンス

MITライセンス