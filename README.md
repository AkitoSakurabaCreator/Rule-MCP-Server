# Rule MCP Server

AIエージェント（Cursor、Cline）が共通のルールを取得・適用できるMCP（Model Context Protocol）サーバーです。

## 機能

- プロジェクトごとのルール管理
- 言語別のグローバルルール管理
- プロジェクトごとのグローバルルール適用設定
- コードのルール違反検証
- MCP経由でのルール配布
- RESTful API エンドポイント
- Web UI によるルール管理

## 技術スタック

- Go 1.21+
- Gin Web Framework
- JSON形式でのルール定義

## セットアップ

### 前提条件

- Go 1.21以上がインストールされていること
- Docker と Docker Compose がインストールされていること
- PostgreSQL 15以上（Docker使用時は自動でインストール）

### インストール

```bash
git clone <repository-url>
cd RuleMCPServer
go mod tidy
```

### 起動

```bash
go run .
```

サーバーはデフォルトでポート8080で起動します。環境変数`PORT`で変更可能です。

## API エンドポイント

### ヘルスチェック

```bash
GET /api/v1/health
```

### ルール取得

```bash
GET /api/v1/rules?project={project_id}
```

### コード検証

```bash
POST /api/v1/rules/validate
Content-Type: application/json

{
  "project_id": "web-app",
  "code": "console.log('test')"
}
```

### プロジェクト管理

```bash
# プロジェクト一覧取得
GET /api/v1/projects

# プロジェクト作成
POST /api/v1/projects
Content-Type: application/json

{
  "project_id": "new-project",
  "name": "New Project",
  "description": "A new project description",
  "language": "javascript",
  "apply_global_rules": true
}
```

# ルール作成
POST /api/v1/rules
Content-Type: application/json

{
  "project_id": "new-project",
  "rule_id": "no-debug-code",
  "name": "No Debug Code",
  "description": "Debug code should not be in production",
  "type": "style",
  "severity": "warning",
  "pattern": "debugger",
  "message": "Debug code detected. Remove before production."
}

# ルール削除
DELETE /api/v1/rules/{project_id}/{rule_id}

### グローバルルール管理

```bash
# 言語別グローバルルール取得
GET /api/v1/global-rules/{language}

# グローバルルール作成
POST /api/v1/global-rules
Content-Type: application/json

{
  "language": "javascript",
  "rule_id": "no-console-log",
  "name": "No Console Log",
  "description": "Console.log statements should not be in production code",
  "type": "style",
  "severity": "warning",
  "pattern": "console.log",
  "message": "Console.log detected. Use proper logging framework in production."
}

# グローバルルール削除
DELETE /api/v1/global-rules/{language}/{rule_id}

# 利用可能な言語一覧取得
GET /api/v1/languages
```

## ルール定義

`rules.json`ファイルでプロジェクトごとのルールを定義します。

```json
{
  "web-app": {
    "project_id": "web-app",
    "rules": [
      {
        "id": "no-console-log",
        "name": "No Console Log",
        "description": "Console.log statements should not be in production code",
        "type": "style",
        "severity": "warning",
        "pattern": "console.log",
        "message": "Console.log detected. Use proper logging framework in production."
      }
    ]
  }
}
```

## MCP Server

このサーバーはMCP（Model Context Protocol）サーバーとしても動作し、CursorやClineから直接ルールを取得できます。

### MCP メソッド

- `getRules`: プロジェクトのルールを取得
- `validateCode`: コードのルール違反を検証

## 設定

### 環境変数

- `PORT`: サーバーのポート番号（デフォルト: 8080）
- `HOST`: サーバーのホストアドレス（デフォルト: 0.0.0.0）
- `ENVIRONMENT`: 実行環境（development/production、デフォルト: development）
- `LOG_LEVEL`: ログレベル（デフォルト: info）

### ポート設定

開発者向けにポートの重複を避けるため、以下のポートを使用します：

- **Rule MCP Server**: 18080 (ホスト) → 8080 (コンテナ)
- **PostgreSQL**: 15432 (ホスト) → 5432 (コンテナ)
- **Web UI**: 13000 (ホスト) → 80 (コンテナ)

これにより、一般的なポート（8080, 5432, 3000）との競合を避けられます。

### 起動例

#### ローカル開発

```bash
# デフォルトポート（8080）で起動
go run .

# カスタムポートで起動
PORT=3000 go run .

# 本番環境で起動
ENVIRONMENT=production PORT=80 go run .

# Makefileを使用
make run-port  # ポート番号を入力して起動
make run-prod  # 本番環境で起動
```

#### Docker を使用

```bash
# Docker イメージをビルド
make docker-build

# サービスを起動
make docker-up

# サービスを停止
make docker-down

# ログを確認
make docker-logs

# リソースをクリーンアップ
make docker-clean
```

## 開発

### テスト

```bash
go test ./...
```

### ビルド

```bash
go build -o rule-mcp-server .
```

## ライセンス

MIT License
