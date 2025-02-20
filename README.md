# pingood

pingoodは、複数のホストに対して同時にpingを実行し、結果を個別のログファイルに記録するGoで実装されたモニタリングツールです。

## 特徴

- 複数ホストの同時監視
- URLおよびホスト名のサポート
- クロスプラットフォーム対応（Windows、Linux、macOS）
- ホストごとの個別ログファイル
- ログファイルの自動再作成機能
- 完全対話型モード対応
- ログファイル名の自動生成機能

## インストール

```bash
go install github.com/araera111/pingood@latest
```

## 使用方法

### コマンドライン引数を使用する場合

```bash
pingood -target "example.com,google.com" -interval 5 -log "example.log,google.log"
```

### 対話型モードを使用する場合

引数なしで実行すると、対話型モードで必要な情報を入力できます：

```bash
pingood
対象のURLまたはIPアドレスを入力してください（カンマ区切りで複数指定可能）: example.com,google.com
Ping実行間隔を秒単位で入力してください（デフォルト: 5）: 10
ログファイル名を自動生成しますか？（y/n、デフォルト: y）: y
```

または一部の引数のみを指定することもできます：

```bash
pingood -interval 5
対象のURLまたはIPアドレスを入力してください（カンマ区切りで複数指定可能）: example.com,google.com
ログファイル名を自動生成しますか？（y/n、デフォルト: y）: y
```

### オプション

- `-target`: ping対象のURLまたはIPアドレス（カンマ区切りで複数指定可能）。省略時は対話的に入力を促します
- `-interval`: ping実行間隔（秒単位、デフォルト: 5秒）。省略時は対話的に入力を促します
- `-log`: ログファイルのパス（カンマ区切りで複数指定可能、targetと同じ数が必要）。省略時は自動生成するかどうかを確認します

### ログ形式

成功時のログ形式：
```
[2025-02-19 18:14:27] SUCCESS - Target: example.com, RTT: 123.456ms
```

エラー時のログ形式：
```
[2025-02-19 18:14:27] ERROR - Target: example.com, Error: failed to resolve host
```

## 要件

- Go 1.16以上

## ライセンス

MITライセンス