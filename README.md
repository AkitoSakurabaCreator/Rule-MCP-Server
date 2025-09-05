# Rule MCP Server

[![npm version](https://badge.fury.io/js/rule-mcp-server.svg)](https://badge.fury.io/js/rule-mcp-server)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

AIエージェント（Cursor、Claude Desktop、Cline）が共通のルールを取得・適用できるMCP（Model Context Protocol）サーバーです。

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

## 🚀 クイックスタート

### 1. MCPサーバーのインストール

```bash
# pnpm dlx経由（推奨・インストール不要）
pnpm dlx rule-mcp-server

# またはグローバルインストール
pnpm add -g rule-mcp-server
```

### 2. AIエージェント設定

#### Cursor
```bash
# 設定テンプレートをコピー
cp config/pnpm-mcp-config.template.json ~/.cursor/mcp.json
```

#### Claude Desktop
```json
// ~/Library/Application Support/Claude/claude_desktop_config.json を作成
{
  "mcpServers": {
    "rule-mcp-server": {
      "command": "pnpm",
      "args": ["dlx", "rule-mcp-server"],
      "env": {
        "RULE_SERVER_URL": "http://localhost:18080",
        "MCP_API_KEY": "${MCP_API_KEY:-}"
      },
      "description": "Standard MCP Server for Rule Management - provides coding rules and validation tools for AI agents",
      "disabled": false,
      "autoApprove": []
    }
  }
}
```

### 3. 利用開始！

AIエージェント（Cursor/Claude Desktop）を再起動して、コーディングルールを自動取得・適用できるようになります。

**📦 npmパッケージ**: [rule-mcp-server](https://www.npmjs.com/package/rule-mcp-server)

### サーバー稼働の前提と起動手順（重要）

このMCPクライアント設定は、バックエンドのRule MCP Serverが稼働していることを前提としています。

#### 稼働確認
```bash
curl http://localhost:18080/api/v1/health
# -> {"status":"ok"} が返れば稼働中
```

#### サーバー未稼働の場合（ローカル起動: Docker）
```bash
# リポジトリを取得
git clone https://github.com/AkitoSakurabaCreator/Rule-MCP-Server.git
cd Rule-MCP-Server

# Dockerで起動（推奨）
docker compose up -d

# 停止
docker compose down
```

#### LAN 内公開の例（チーム運用）
- サーバーをLAN上のホストで起動し、クライアント側の環境変数をLAN IPに設定:
```json
{
  "mcpServers": {
    "rule-mcp-server": {
      "env": {
        "RULE_SERVER_URL": "http://192.168.1.20:18080",
        "MCP_API_KEY": "${MCP_API_KEY:-}"
      }
    }
  }
}
```

参考: Makefile を使う場合は `make docker-up` / `make docker-down`

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

### 標準的なMCPサーバー設定（推奨）

このプロジェクトは**標準的なMCP（Model Context Protocol）サーバー**を提供します。

#### **特徴**
- ✅ **標準MCP SDK使用**: `@modelcontextprotocol/sdk`による完全準拠
- ✅ **StdioServerTransport**: 標準的なMCPクライアントと互換
- ✅ **pnpmパッケージ対応**: `pnpm dlx`で簡単インストール
- ✅ **Docker対応**: 本番環境での安定動作
- ✅ **テンプレート設定**: 簡単セットアップ

#### **1. MCPサーバーのインストール**

##### **pnpm経由（推奨）**
```bash
# グローバルインストール
pnpm add -g rule-mcp-server

# またはpnpm dlx経由（インストール不要）
pnpm dlx rule-mcp-server
```

**📦 npmパッケージ**: [rule-mcp-server](https://www.npmjs.com/package/rule-mcp-server) として公開済み

##### **開発版ビルド**
```bash
# 依存関係のインストール
make install-mcp

# MCPサーバーのビルド
make build-mcp
```

#### **2. 環境別設定**

##### **pnpmパッケージ使用（推奨）**
```bash
# pnpmパッケージ用設定テンプレートを使用
cp config/pnpm-mcp-config.template.json ~/.cursor/mcp.json
```

設定例（pnpmパッケージ）:
```json
{
  "mcpServers": {
    "rule-mcp-server": {
      "command": "pnpm",
      "args": ["dlx", "rule-mcp-server"],
      "env": {
        "RULE_SERVER_URL": "http://localhost:18080",
        "MCP_API_KEY": ""
      },
      "description": "Standard MCP Server for Rule Management",
      "disabled": false,
      "autoApprove": []
    }
  }
}
```

##### **Docker環境**
```bash
# Docker環境を起動
make docker-up

# Docker用設定テンプレートを使用
cp config/docker-mcp-config.template.json ~/.cursor/mcp_settings.json
```

設定例（Docker環境）:
```json
{
  "mcpServers": {
    "rule-mcp-server": {
      "command": "node",
      "args": ["/path/to/your/RuleMCPServer/cmd/mcp-server/build/index.js"],
      "env": {
        "RULE_SERVER_URL": "http://localhost:18080",
        "MCP_API_KEY": ""
      },
      "description": "Standard MCP Server for Rule Management (Docker)",
      "disabled": false,
      "autoApprove": []
    }
  }
}
```

##### **開発環境**
```bash
# 開発サーバーを起動
make run

# 標準設定テンプレートを使用
cp config/standard-mcp-config.template.json ~/.cursor/mcp_settings.json
```

設定例（開発環境）:
```json
{
  "mcpServers": {
    "rule-mcp-server": {
      "command": "node",
      "args": ["/path/to/your/RuleMCPServer/cmd/mcp-server/build/index.js"],
      "env": {
        "RULE_SERVER_URL": "http://localhost:18081",
        "MCP_API_KEY": ""
      },
      "description": "Standard MCP Server for Rule Management (Development)",
      "disabled": false,
      "autoApprove": []
    }
  }
}
```

#### **3. 設定ファイルの配置**
- **Cursor**: `~/.cursor/mcp_settings.json`
- **Cline**: `~/.cline/mcp_settings.json`
- **Claude Desktop**: `~/Library/Application Support/Claude/claude_desktop_config.json`

#### **4. 利用可能なツール**
標準MCPサーバーは以下のツールを提供：

| ツール名            | 説明                         | 必須パラメータ         |
|---------------------|------------------------------|------------------------|
| `getRules`          | プロジェクトルール取得       | `project_id`           |
| `validateCode`      | コード検証                   | `project_id`, `code`   |
| `getProjectInfo`    | プロジェクト情報取得         | `project_id`           |
| `autoDetectProject` | プロジェクト自動検出         | `path`                 |
| `scanLocalProjects` | ローカルプロジェクトスキャン | `base_path` (optional) |
| `getGlobalRules`    | グローバルルール取得         | `language`             |

#### **5. 利用可能なリソース**
標準MCPサーバーは以下のリソースを提供：

| リソースURI                      | 説明                   |
|----------------------------------|------------------------|
| `rule://projects/list`           | プロジェクト一覧       |
| `rule://{project_id}/rules`      | プロジェクト固有ルール |
| `rule://{project_id}/info`       | プロジェクト情報       |
| `rule://global-rules/{language}` | 言語別グローバルルール |

### 従来のHTTP設定（互換性）

従来のHTTP APIを使用する場合：

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

#### Cursor
```json
// ~/.cursor/mcp.json
{
  "mcpServers": {
    "rule-mcp-server": {
      "command": "pnpm",
      "args": [
        "dlx",
        "rule-mcp-server"
      ],
      "env": {
        "RULE_SERVER_URL": "http://localhost:18080",
        "MCP_API_KEY": "${MCP_API_KEY:-}"
      },
      "description": "Standard MCP Server for Rule Management - provides coding rules and validation tools for AI agents",
      "disabled": false,
      "autoApprove": []
    }
  }
}
```

#### Claude Desktop
```json
// ~/Library/Application Support/Claude/claude_desktop_config.json
{
  "mcpServers": {
    "rule-mcp-server": {
      "command": "pnpm",
      "args": [
        "dlx",
        "rule-mcp-server"
      ],
      "env": {
        "RULE_SERVER_URL": "http://localhost:18080",
        "MCP_API_KEY": "${MCP_API_KEY:-}"
      },
      "description": "Standard MCP Server for Rule Management - provides coding rules and validation tools for AI agents",
      "disabled": false,
      "autoApprove": []
    }
  }
}
```

補足: `MCP_API_KEY` は未設定でも動作します（Publicアクセス）。チーム運用や管理APIを使う場合にのみ設定してください。

## 🚀 クイックスタート

### pnpmパッケージ使用（推奨）

```bash
# 1. MCPサーバーをインストール（またはpnpm dlxで自動インストール）
pnpm add -g rule-mcp-server

# 2. 設定ファイルを作成
cp config/pnpm-mcp-config.template.json ~/.cursor/mcp.json

# 3. AIエージェント（Cursor/Claude Desktop/Cline）で利用開始！
```

**📦 npmパッケージ**: [rule-mcp-server](https://www.npmjs.com/package/rule-mcp-server) として公開済み

### Docker環境

```bash
# 1. リポジトリをクローン
git clone https://github.com/AkitoSakurabaCreator/Rule-MCP-Server.git
cd Rule-MCP-Server

# 2. Docker環境を起動
make docker-up

# 3. MCPサーバーをビルド
make install-mcp && make build-mcp

# 4. 設定ファイルを作成
cp config/docker-mcp-config.template.json ~/.cursor/mcp_settings.json
# パスを編集: ${PROJECT_PATH} → /path/to/your/Rule-MCP-Server

# 5. AIエージェント（Cursor/Cline）で利用開始！
```

### 開発環境

```bash
# 1. 依存関係のインストール
go mod tidy
cd frontend && npm install && cd ..

# 2. サーバー起動
make run        # バックエンド（ポート18081）
make run-frontend  # フロントエンド（ポート3000）

# 3. MCPサーバーをビルド
make install-mcp && make build-mcp

# 4. 設定ファイルを作成
cp config/standard-mcp-config.template.json ~/.cursor/mcp_settings.json
```

## 🌟 主な特徴

### ✅ 標準MCP準拠
- **完全なMCP互換性**: `@modelcontextprotocol/sdk`使用
- **標準的なツール・リソース**: AIエージェントとの完璧な統合
- **StdioServerTransport**: 標準的な通信プロトコル

### ✅ 豊富な機能
- **6つのMCPツール**: ルール取得、コード検証、プロジェクト自動検出など
- **5つのMCPリソース**: プロジェクト一覧、ルール情報、グローバルルール
- **プロジェクト自動検出**: AIが自動的に適切なルールを適用

### ✅ 本番対応
- **Docker環境**: 安定した本番運用
- **PostgreSQL**: 高性能データベース
- **認証・認可**: APIキー、セッション認証
- **多言語対応**: 6言語サポート（日本語、英語、中国語など）

### ✅ 開発者フレンドリー
- **クリーンアーキテクチャ**: 保守性の高い設計
- **Web UI**: React製の管理画面
- **テンプレート設定**: 簡単セットアップ
- **豊富なドキュメント**: 詳細な使用方法

## 🎯 使用ケース

### 個人開発者
```bash
# 認証なしで簡単利用
curl http://localhost:18080/api/v1/rules?project_id=my-project
```

### チーム開発
```bash
# APIキーでチーム管理
curl -H "X-API-Key: team_key" http://localhost:18080/api/v1/projects
```

### AIエージェント統合
```json
// Cursor/Clineで自動ルール適用
{
  "mcpServers": {
    "rule-mcp-server": {
      "command": "node",
      "args": ["/path/to/Rule-MCP-Server/cmd/mcp-server/build/index.js"],
      "env": {
        "RULE_SERVER_URL": "http://localhost:18080"
      }
    }
  }
}
```

## 📊 システム要件

### 最小要件
- **OS**: Linux, macOS, Windows
- **Go**: 1.21+
- **Node.js**: 18+
- **Docker**: 20.10+ (推奨)

### 推奨要件
- **メモリ**: 2GB以上
- **ストレージ**: 1GB以上
- **CPU**: 2コア以上

## 🔧 開発・貢献

### 開発環境のセットアップ

```bash
# 1. フォーク & クローン
git clone https://github.com/your-username/Rule-MCP-Server.git
cd Rule-MCP-Server

# 2. 開発用ブランチ作成
git checkout -b feature/your-feature

# 3. 依存関係インストール
make deps
make install-mcp

# 4. 開発サーバー起動
make run
make run-frontend
```

### テスト実行

```bash
# バックエンドテスト
make test

# フロントエンドテスト
cd frontend && npm test

# MCPサーバーテスト
cd cmd/mcp-server && npm test
```

### コード品質チェック

```bash
# Go言語
make fmt
make lint

# TypeScript
cd cmd/mcp-server && npm run lint
```

## 🤝 貢献ガイドライン

### 貢献の流れ

1. **Issue作成**: バグ報告や機能提案
2. **フォーク**: リポジトリをフォーク
3. **ブランチ作成**: `feature/your-feature` または `fix/your-fix`
4. **開発**: コード変更とテスト追加
5. **プルリクエスト**: 詳細な説明と共に提出

### コミットメッセージ規約

```bash
# 機能追加
feat: add project auto-detection feature

# バグ修正
fix: resolve MCP server connection issue

# ドキュメント更新
docs: update README with Docker setup

# リファクタリング
refactor: improve error handling in MCP handlers
```

### 開発のベストプラクティス

- **テスト**: 新機能には必ずテストを追加
- **ドキュメント**: APIの変更は必ずドキュメント更新
- **型安全性**: TypeScriptの型定義を適切に使用
- **エラーハンドリング**: 適切なエラーメッセージとログ出力

## 🏆 コミュニティ

### 貢献者

このプロジェクトは以下の素晴らしい貢献者によって支えられています：

- [@AkitoSakurabaCreator](https://github.com/AkitoSakurabaCreator) - プロジェクト創設者・メンテナー

### 謝辞

- **Model Context Protocol**: 標準的なAIエージェント統合を可能にする素晴らしいプロトコル
- **Go Community**: 高性能なバックエンド開発環境
- **React Community**: 優れたフロントエンド開発体験
- **Docker**: 一貫した開発・本番環境

## 📄 ライセンス

MIT License - 詳細は [LICENSE](LICENSE) ファイルを参照してください。

## 🆘 サポート・コミュニティ

### 問題報告・質問

- **🐛 バグ報告**: [GitHub Issues](https://github.com/AkitoSakurabaCreator/Rule-MCP-Server/issues)
- **💡 機能提案**: [GitHub Issues](https://github.com/AkitoSakurabaCreator/Rule-MCP-Server/issues)
- **❓ 質問・議論**: [GitHub Discussions](https://github.com/AkitoSakurabaCreator/Rule-MCP-Server/discussions)

### ドキュメント

- **📚 詳細ドキュメント**: [GitHub Wiki](https://github.com/AkitoSakurabaCreator/Rule-MCP-Server/wiki)

### コミュニティ

- **💬 Discord**: [Rule MCP Server Community](https://discord.gg/dCAUC8m6dw)
- **🐦 X (旧Twitter)**: [@_sakuraba_akito](https://x.com/_sakuraba_akito)

### プロジェクト支援

- **💖 スポンサー**: [GitHub Sponsors](https://github.com/sponsors/AkitoSakurabaCreator)

---

**⭐ このプロジェクトが役に立ったら、GitHubでスターをお願いします！**

**🚀 AIエージェントとの統合で、より良いコード品質を実現しましょう！**
