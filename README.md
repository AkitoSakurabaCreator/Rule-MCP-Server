# Rule MCP Server

AIエージェント（Cursor、Cline）が共通のルールを取得・適用できるMCP（Model Context Protocol）サーバーです。

## 機能

- プロジェクトごとのルール管理
- 言語別のグローバルルール管理
- プロジェクトごとのグローバルルール適用設定
- コードのルール違反検証
- MCP経由でのルール配布
- RESTful API エンドポイント
- **ReactベースのWeb UI**によるルール管理
- **多言語対応（i18n）**: 英語、日本語、中国語、ヒンディー語、スペイン語、アラビア語
- **ダークテーマ対応**: ライト/ダークモード切り替え
- **クリーンアーキテクチャ**による保守性の高い設計

## 技術スタック

### バックエンド

- **Go 1.21+** + **Gin Web Framework**
- **クリーンアーキテクチャ**（Domain, Usecase, Interface, Infrastructure）
- **PostgreSQL** データベース
- **MCP（Model Context Protocol）** サポート

### フロントエンド

- **React 18** + **TypeScript**
- **Material-UI (MUI)** コンポーネントライブラリ
- **React Router** によるルーティング
- **i18next** による多言語対応
- **Axios** によるAPI通信

## セットアップ

### 前提条件

- Go 1.21以上がインストールされていること
- Node.js 18以上がインストールされていること
- Docker と Docker Compose がインストールされていること
- PostgreSQL 15以上（Docker使用時は自動でインストール）

### インストール

```bash
git clone <repository-url>
cd RuleMCPServer

# バックエンド依存関係
go mod tidy

# フロントエンド依存関係
cd frontend
npm install
cd ..
```

### 起動

#### バックエンド

```bash
# 開発環境（安全ポート18081）
PORT=18081 go run ./cmd/server

# 本番環境（安全ポート18080）
ENVIRONMENT=production PORT=18080 go run ./cmd/server

# Makefileを使用
make run        # 開発環境（ポート18081）
make run-prod   # 本番環境（ポート18080）
```

#### フロントエンド

```bash
cd frontend
npm start

# または Makefileを使用
make run-frontend
```

## アーキテクチャ

### クリーンアーキテクチャ

```
cmd/server/          # エントリーポイント
├── internal/
│   ├── domain/      # エンティティ、リポジトリインターフェース
│   ├── usecase/     # ビジネスロジック
│   ├── interface/   # HTTPハンドラー、MCPハンドラー
│   └── infrastructure/ # データベース実装
└── frontend/        # Reactフロントエンド
```

### レイヤー構成

- **Domain Layer**: ビジネスエンティティとルール
- **Usecase Layer**: アプリケーションのビジネスロジック
- **Interface Layer**: HTTP API、MCPプロトコル
- **Infrastructure Layer**: PostgreSQL、ファイルシステム

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

### ルール管理

```bash
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
```

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

## MCP（Model Context Protocol）Server

このサーバーはMCPサーバーとして動作し、CursorやClineから直接ルールを取得できます。

### MCP エンドポイント

```
POST /mcp/request    # HTTP MCPリクエスト
GET  /mcp/ws         # WebSocket MCP接続
```

### MCP メソッド

- **`getRules`**: プロジェクトのルールを取得
- **`validateCode`**: コードのルール違反を検証
- **`getProjectInfo`**: プロジェクト情報を取得

### Cursor設定

`~/.cursor/mcp.json`に以下を設定：

```json
{
  "mcpServers": {
    "rule-mcp-server": {
      "command": "curl",
      "args": [
        "-X", "POST",
        "-H", "Content-Type: application/json",
        "-d", "{\"id\":\"${requestId}\",\"method\":\"${method}\",\"params\":${params}}",
        "http://localhost:18081/mcp/request"
      ],
      "env": {
        "MCP_SERVER_URL": "http://localhost:18081"
      },
      "description": "Rule MCP Server for AI agents to retrieve and apply common rules"
    }
  }
}
```

## フロントエンド機能

### 多言語対応（i18n）

以下の言語をサポート：

- **英語 (en)**: デフォルト言語
- **日本語 (ja)**: 完全対応
- **中国語 (zh-CN)**: 完全対応
- **ヒンディー語 (hi)**: 完全対応
- **スペイン語 (es)**: 完全対応
- **アラビア語 (ar)**: 完全対応（RTL対応）

### テーマ切り替え

- **ライトテーマ**: 明るく読みやすい
- **ダークテーマ**: 目に優しい夜間モード
- **設定の自動保存**: ブラウザのlocalStorageに保存

### Web UI機能

- **プロジェクト管理**: 作成、編集、削除
- **ルール管理**: プロジェクト固有ルール、グローバルルール
- **コード検証**: リアルタイムでのルール違反チェック
- **レスポンシブデザイン**: モバイル・タブレット対応

## ルール定義

### ルール作成の項目

1. **Rule ID**: ルールの一意識別子（例：`no-console-log`）
2. **Name**: ルールの表示名（例：`No Console Log`）
3. **Description**: ルールの内容・目的・理由（例：`Console.log statements should not be in production code`）
4. **Type**: ルールのカテゴリ（naming, formatting, security, performance等）
5. **Severity**: 重要度（error, warning, info）
6. **Pattern**: 検出する正規表現パターン（例：`console\.log`）
7. **Message**: 違反時の修正指示（例：`Console.log detected. Use proper logging framework in production.`）
8. **Active**: ルールの有効/無効

### サンプルルール

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

## 設定

### 環境変数

- `PORT`: サーバーのポート番号（デフォルト: 8080）
- `HOST`: サーバーのホストアドレス（デフォルト: 0.0.0.0）
- `ENVIRONMENT`: 実行環境（development/production、デフォルト: development）
- `LOG_LEVEL`: ログレベル（デフォルト: info）

### データベース設定

- `DB_HOST`: データベースホスト（デフォルト: localhost）
- `DB_PORT`: データベースポート（デフォルト: 5432）
- `DB_NAME`: データベース名（デフォルト: rule_mcp_db）
- `DB_USER`: データベースユーザー（デフォルト: rule_mcp_user）
- `DB_PASSWORD`: データベースパスワード

### ポート設定

開発者向けにポートの重複を避けるため、以下のポートを使用します：

- **Rule MCP Server**: 18081 (開発) / 18080 (本番)
- **PostgreSQL**: 15432 (ホスト) → 5432 (コンテナ)
- **Web UI**: 13000 (ホスト) → 80 (コンテナ)

これにより、一般的なポート（8080, 5432, 3000）との競合を避けられます。

### 起動例

#### ローカル開発

```bash
# バックエンド（開発環境）
PORT=18081 go run ./cmd/server

# フロントエンド
cd frontend && npm start

# Makefileを使用
make run           # バックエンド（ポート18081）
make run-frontend  # フロントエンド
make run-prod      # 本番環境（ポート18080）
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
# バックエンドテスト
go test ./...

# フロントエンドテスト
cd frontend && npm test
```

### ビルド

```bash
# バックエンド
go build -o rule-mcp-server ./cmd/server

# フロントエンド
cd frontend && npm run build
```

### コード品質

```bash
# Go言語の品質チェック
make glb

# 全リポジトリの品質チェック
make glb_repo
```

## トラブルシューティング

### よくある問題

#### 1. データベース接続エラー

```bash
# データベースコンテナの再起動
docker-compose restart postgres

# データベースの完全再作成
docker-compose down -v
docker-compose up -d postgres
```

#### 2. ポート競合

```bash
# 使用中のポートを確認
lsof -i :18081
lsof -i :13000

# 別のポートで起動
PORT=18082 go run ./cmd/server
```

#### 3. MCPサーバーが応答しない

```bash
# サーバーの起動確認
curl http://localhost:18081/api/v1/health

# MCPエンドポイントのテスト
curl -X POST http://localhost:18081/mcp/request \
  -H "Content-Type: application/json" \
  -d '{"id":"test","method":"getRules","params":{"project_id":"web-app"}}'
```

#### 4. フロントエンドのビルドエラー

```bash
# 依存関係の再インストール
cd frontend
rm -rf node_modules package-lock.json
npm install
```

### ログ確認

```bash
# バックエンドログ
docker logs rule-mcp-server

# データベースログ
docker logs rule-mcp-postgres

# フロントエンドログ
cd frontend && npm start
```

## 貢献

1. このリポジトリをフォーク
2. フィーチャーブランチを作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'Add some amazing feature'`)
4. ブランチにプッシュ (`git push origin feature/amazing-feature`)
5. プルリクエストを作成

## ライセンス

MIT License

## サポート

問題が発生した場合や質問がある場合は、以下の方法でサポートを受けることができます：

- **Issues**: GitHubのIssuesページで問題を報告
- **Discussions**: GitHubのDiscussionsページで質問・議論
- **Wiki**: 詳細なドキュメントやチュートリアル
