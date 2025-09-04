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

### 初期管理アカウント

システム初回起動時に以下の初期管理アカウントが自動的に作成されます：

- **ユーザー名**: `admin`
- **パスワード**: `admin123`
- **メール**: `admin@rulemcp.com`
- **権限**: 管理者（admin）

**重要**: 初回ログイン後は必ずパスワードを変更してください。

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

#### 開発環境

```bash
# バックエンド（安全ポート18081）
PORT=18081 go run ./cmd/server

# フロントエンド
cd frontend
npm start

# Makefileを使用
make run        # 開発環境（ポート18081）
make run-frontend
```

#### 本番環境（Docker）

```bash
# 本番環境用の環境変数ファイルを作成
cp env.prod.example .env.prod
# .env.prodファイルの値を本番環境に合わせて編集

# 本番環境をデプロイ
make -f Makefile.prod deploy

# 本番環境のステータス確認
make -f Makefile.prod status

# 本番環境のログ確認
make -f Makefile.prod logs

# 本番環境を停止
make -f Makefile.prod down

# 本番環境をクリーンアップ
make -f Makefile.prod clean
```

#### 本番環境の特徴

- **セキュリティ強化**: 非rootユーザーでの実行、環境変数による設定
- **パフォーマンス最適化**: マルチステージビルド、軽量なAlpine Linux
- **ヘルスチェック**: 各サービスの健全性監視
- **ログ管理**: 構造化されたログ出力とローテーション
- **バックアップ**: データベースの自動バックアップ機能
- **スケーラビリティ**: Docker SwarmやKubernetes対応の準備

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

## MCP（Model Context Protocol）機能

### 基本機能

- **`getRules`**: プロジェクトIDとオプションの言語を指定してルールを取得
- **`validateCode`**: コードのルール違反を検証
- **`getProjectInfo`**: プロジェクト情報を取得

### プロジェクト自動検出機能 🆕

AIエージェントが**自動的にプロジェクトを認識**し、適切なルールを取得できる高度な機能です。

#### **自動検出の優先順位**

1. **ディレクトリ名ベース検出**（信頼度95%）
   - ディレクトリ名をプロジェクトIDとして検索
   - 除外ディレクトリ: `node_modules`, `vendor`, `dist`, `build`, `target`, `.git`, `.vscode`

2. **Gitリポジトリ名検出**（信頼度90%）
   - `.git/config`からorigin URLを解析
   - SSH形式: `git@github.com:username/repo-name.git`
   - HTTPS形式: `https://github.com/username/repo-name.git`

3. **言語固有ファイル検出**（信頼度85%）
   - `go.mod` → Goプロジェクト
   - `package.json` → Node.jsプロジェクト
   - `requirements.txt` → Pythonプロジェクト
   - `pom.xml` → Javaプロジェクト
   - `Cargo.toml` → Rustプロジェクト
   - `composer.json` → PHPプロジェクト
   - `Gemfile` → Rubyプロジェクト

4. **デフォルトプロジェクト**（信頼度70%）
   - 検出できない場合のフォールバック

#### **新しいMCPメソッド**

##### **`autoDetectProject`**
指定されたパスからプロジェクトを自動検出します。

```json
{
  "id": "auto-detect",
  "method": "autoDetectProject",
  "params": {
    "path": "/path/to/your/project"
  }
}
```

**レスポンス例:**
```json
{
  "id": "auto-detect",
  "result": {
    "project": {
      "project_id": "web-app",
      "name": "Web Application",
      "language": "javascript"
    },
    "rules": [...],
    "detection_method": "directory_name",
    "confidence": 0.95,
    "message": "ディレクトリ名 'web-app' からプロジェクトを検出しました"
  }
}
```

##### **`scanLocalProjects`**
ローカルディレクトリを再帰的にスキャンしてプロジェクトを検出します。

```json
{
  "id": "scan-local",
  "method": "scanLocalProjects",
  "params": {
    "base_path": "/home/user/projects"
  }
}
```

**レスポンス例:**
```json
{
  "id": "scan-local",
  "result": {
    "projects": [
      {
        "project": {...},
        "rules": [...],
        "detection_method": "git_repository",
        "confidence": 0.90,
        "message": "Gitリポジトリ名からプロジェクトを検出しました"
      }
    ],
    "count": 3
  }
}
```

#### **使用例**

##### **Cursor/Clineでの設定**

###### **1. cursor-mcp-config.json（推奨）**
```json
{
  "mcpServers": {
    "rule-mcp": {
      "command": "go",
      "args": ["run", "./cmd/server"],
      "env": {
        "PORT": "18081",
        "DB_HOST": "localhost",
        "DB_PORT": "15432",
        "DB_USER": "rule_mcp_user",
        "DB_PASSWORD": "rule_mcp_password",
        "DB_NAME": "rule_mcp_db"
      }
    }
  }
}
```

**ファイル配置場所:**
- **Cursor**: `~/.cursor/mcp-servers/rule-mcp.json`
- **Cline**: `~/.cline/mcp-servers/rule-mcp.json`

###### **2. mcp-client-config.json（完全版）**
```json
{
  "mcpServers": {
    "rule-mcp": {
      "command": "go",
      "args": ["run", "./cmd/server"],
      "env": {
        "PORT": "18081",
        "DB_HOST": "localhost",
        "DB_PORT": "15432",
        "DB_USER": "rule_mcp_user",
        "DB_PASSWORD": "rule_mcp_password",
        "DB_NAME": "rule_mcp_db"
      },
      "cwd": "/path/to/your/RuleMCPServer"
    }
  }
}
```

###### **3. 環境変数なし（JSONファイルモード）**
```json
{
  "mcpServers": {
    "rule-mcp": {
      "command": "go",
      "args": ["run", "./cmd/server"],
      "env": {
        "PORT": "18081"
      }
    }
  }
}
```

##### **設定の説明**

| パラメータ | 説明             | 必須 | デフォルト                |
|------------|------------------|------|---------------------------|
| `command`  | 実行するコマンド | ✅    | `go`                      |
| `args`     | コマンドの引数   | ✅    | `["run", "./cmd/server"]` |
| `env`      | 環境変数         | ❌    | なし                      |
| `cwd`      | 作業ディレクトリ | ❌    | カレントディレクトリ      |

##### **環境別の設定例**

###### **開発環境（JSONファイルモード）**
```json
{
  "mcpServers": {
    "rule-mcp": {
      "command": "go",
      "args": ["run", "./cmd/server"],
      "env": {
        "PORT": "18081"
      }
    }
  }
}
```

###### **本番環境（PostgreSQL接続）**
```json
{
  "mcpServers": {
    "rule-mcp": {
      "command": "go",
      "args": ["run", "./cmd/server"],
      "env": {
        "PORT": "18081",
        "DB_HOST": "your-db-host",
        "DB_PORT": "5432",
        "DB_USER": "your-db-user",
        "DB_PASSWORD": "your-db-password",
        "DB_NAME": "your-db-name"
      }
    }
  }
}
```

##### **設定ファイルの優先順位**

1. **プロジェクト固有**: `./cursor-mcp-config.json`
2. **ユーザー設定**: `~/.cursor/mcp-servers/rule-mcp.json`
3. **グローバル設定**: `~/.cursor/mcp-servers/rule-mcp.json`

##### **トラブルシューティング**

###### **よくある問題と解決方法**

| 問題                         | 原因              | 解決方法                 |
|------------------------------|-------------------|--------------------------|
| `Method not found`           | 古いMCPハンドラー | サーバーを再起動         |
| `Connection refused`         | ポートが使用中    | `lsof -i :18081`で確認   |
| `Database connection failed` | DB接続設定ミス    | 環境変数を確認           |
| `Permission denied`          | ファイル権限      | `chmod +x`で実行権限付与 |

##### **AIエージェントの自動検出**
AIエージェントは以下のように自動的にプロジェクトを認識できます：

1. **作業ディレクトリの自動検出**
2. **適切なルールセットの取得**
3. **言語固有のルール適用**
4. **プロジェクト固有のルール適用**

これにより、AIエージェントは**コンテキストを理解せずに**、常に適切なルールでコードレビューや提案を行うことができます。

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

#### 5. 認証エラー

```bash
# APIキーの確認
curl -H "X-API-Key: your_api_key" http://localhost:18081/api/v1/projects

# 権限レベルの確認
curl -H "X-API-Key: your_api_key" http://localhost:18081/api/v1/auth/me
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

## 権限管理システム

### アクセスレベル

Rule MCP Serverは3段階のアクセスレベルを提供します：

#### **Public（公開）**
- **認証**: 不要
- **権限**: 公開ルール・プロジェクトの閲覧、コード検証
- **用途**: 個人使用、オープンソースプロジェクト
- **MCPアクセス**: 制限なし（レート制限あり）

#### **User（ユーザー）**
- **認証**: APIキー必須
- **権限**: 個人ルール・プロジェクトの作成・編集・削除
- **用途**: 個人開発者、小規模チーム
- **MCPアクセス**: 制限なし

#### **Admin（管理者）**
- **認証**: 管理者APIキー必須
- **権限**: 全権限（ユーザー管理、グローバルルール管理）
- **用途**: チームリーダー、システム管理者
- **MCPアクセス**: 制限なし

### 認証方式

#### **APIキー認証**
```bash
# ヘッダーでの認証
curl -H "X-API-Key: your_api_key" http://localhost:18081/api/v1/projects

# MCPリクエストでの認証
curl -X POST http://localhost:18081/mcp/request \
  -H "X-API-Key: your_api_key" \
  -H "Content-Type: application/json" \
  -d '{"id":"test","method":"createRule","params":{...}}'
```

#### **セッション認証**
```bash
# ログイン
curl -X POST http://localhost:18081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user","password":"password"}'

# セッションでの認証
curl -H "Cookie: session=session_token" http://localhost:18081/api/v1/projects
```

### チーム協働機能

#### **プロジェクトメンバー管理**
```bash
# チームメンバーの追加
curl -X POST http://localhost:18081/api/v1/projects/team-project/members \
  -H "X-API-Key: admin_key" \
  -H "Content-Type: application/json" \
  -d '{"username":"developer1","role":"member"}'

# チームメンバーの権限確認
curl http://localhost:18081/api/v1/projects/team-project/members \
  -H "X-API-Key: user_key"
```

#### **プロジェクト可視性**
- **Public**: 全ユーザーが閲覧可能
- **Team**: チームメンバーのみ閲覧可能
- **Private**: プロジェクト所有者のみ閲覧可能

#### **権限の細分化**
```json
{
  "permissions": {
    "read": true,    // ルール・プロジェクトの閲覧
    "write": true,   // ルール・プロジェクトの作成・編集
    "delete": false, // ルール・プロジェクトの削除
    "admin": false   // メンバー管理
  }
}
```

## 設定ファイル

### **認証設定** (`config/auth.yaml`)
権限レベル、レート制限、セキュリティ設定を定義

### **環境変数** (`config/environment.md`)
サーバー、データベース、認証、MCP設定の環境変数一覧

### **MCPクライアント設定**
- **`config/simple-mcp-config.json`**: シンプル版MCP設定（認証なし、初心者向け）
- **`config/mcp-client-config.json`**: 完全版MCP設定（認証・チーム機能対応）

### **データベーススキーマ** (`init.sql`)
権限管理テーブル、ユーザー管理、チーム協働機能を含む完全なスキーマ

## セキュリティ機能

### **レート制限**
- **Public**: 50 req/min（制限付き）
- **User**: 100 req/min
- **Admin**: 200 req/min

### **監査ログ**
- 認証試行の記録
- 権限変更の記録
- ルール変更の記録
- 保持期間: 365日

### **パスワードポリシー**
- 最小8文字
- 大文字・小文字・数字・特殊文字必須
- セッションタイムアウト: 8時間

### **APIキーセキュリティ**
- bcryptハッシュ化
- 有効期限設定
- HTTPS必須（本番環境）
- 使用ログ記録

## 使用例

### **個人使用（Public）**
```bash
# 認証なしでルール取得
curl http://localhost:18081/api/v1/rules?project_id=web-app

# MCP経由でコード検証
curl -X POST http://localhost:18081/mcp/request \
  -H "Content-Type: application/json" \
  -d '{"id":"test","method":"validateCode","params":{"project_id":"web-app","code":"console.log(\"test\")"}}'
```

### **チーム使用（User/Admin）**
```bash
# APIキーでルール作成
curl -X POST http://localhost:18081/api/v1/rules \
  -H "X-API-Key: user_key" \
  -H "Content-Type: application/json" \
  -d '{"project_id":"team-project","rule_id":"no-todo","name":"No TODO","type":"quality","severity":"warning","pattern":"TODO:","message":"TODO comment detected"}'

# チームメンバー管理
curl -X POST http://localhost:18081/api/v1/projects/team-project/members \
  -H "X-API-Key: admin_key" \
  -H "Content-Type: application/json" \
  -d '{"username":"new_dev","role":"member"}'
```

### **MCPクライアント設定**
```json
// ~/.cursor/mcp.json
{
  "mcpServers": {
    "rule-mcp-server": {
      "command": "curl",
      "args": [
        "-X", "POST",
        "-H", "Content-Type: application/json",
        "-H", "X-API-Key: ${MCP_API_KEY}",
        "-d", "{\"id\":\"${requestId}\",\"method\":\"${method}\",\"params\":${params}}",
        "${MCP_SERVER_URL}/mcp/request"
      ],
      "env": {
        "MCP_SERVER_URL": "http://localhost:18081",
        "MCP_API_KEY": "your_api_key_here",
        "AUTO_INJECT": "true"
      }
    }
  }
}
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
