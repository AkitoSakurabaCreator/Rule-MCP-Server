# 設定ファイル ガイド

Rule MCP Serverの設定ファイルとその使用方法について説明します。

## 設定ファイル一覧

### 1. `auth.yaml` - 認証・権限設定

権限レベル、レート制限、セキュリティ設定を定義するYAMLファイルです。

#### 主要設定項目
- **APIキー設定**: 有効期限、ハッシュアルゴリズム、レート制限
- **アクセスレベル**: public、user、adminの権限定義
- **MCPプロトコル**: 保護されたメソッドとパブリックメソッド
- **チーム協働**: メンバー数制限、デフォルト権限
- **セキュリティ**: パスワードポリシー、セッション管理
- **監査ログ**: ログ記録と保持期間

#### 使用例
```yaml
auth:
  access_levels:
    public:
      permissions:
        - "read:public_rules"
        - "validate:public_code"
      mcp_access: true
      rate_limit_multiplier: 0.5
```

### 2. `environment.md` - 環境変数設定

サーバー起動時に必要な環境変数の一覧と設定例を提供します。

#### 主要カテゴリ
- **サーバー設定**: ポート、ホスト、環境
- **データベース設定**: 接続情報、SSL設定
- **認証設定**: 有効化、デフォルトレベル、セッション
- **APIキー設定**: 有効期限、長さ、ハッシュアルゴリズム
- **レート制限**: リクエスト数制限、バースト制限
- **セキュリティ設定**: パスワード要件、HTTPS要件
- **チーム協働**: メンバー数制限、可視性設定
- **監査ログ**: ログ記録、保持期間
- **MCP設定**: 有効化、アクセスレベル、メソッド制限
- **CORS設定**: 許可されたオリジン、メソッド、ヘッダー

#### 使用例
```bash
# 開発環境
PORT=18081
ENVIRONMENT=development
AUTH_ENABLED=true
AUTH_DEFAULT_ACCESS_LEVEL=public

# 本番環境
PORT=18080
ENVIRONMENT=production
AUTH_REQUIRE_HTTPS=true
LOG_LEVEL=warn
```

### 3. `mcp-client-config.json` - MCPクライアント設定

Cursor/ClineなどのMCPクライアント用の設定ファイルです。

#### 主要機能
- **認証対応**: APIキーによる権限管理
- **拡張メソッド**: チーム協働機能を含む
- **権限別アクセス**: アクセスレベルに応じた機能制限
- **環境変数**: 動的な設定値の注入

#### 設定例
```json
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

### 4. `init.sql` - データベーススキーマ

権限管理システムを含む完全なデータベーススキーマです。

#### 主要テーブル
- **projects**: プロジェクト情報（アクセスレベル、作成者）
- **rules**: ルール定義（アクセスレベル、作成者）
- **global_rules**: グローバルルール（言語別、アクセスレベル）
- **users**: ユーザー管理（ロール、アクティブ状態）
- **api_keys**: APIキー管理（ハッシュ、有効期限、権限レベル）
- **project_members**: チームメンバー管理（権限、ロール）
- **rule_violations**: ルール違反記録

#### サンプルデータ
- デフォルトプロジェクト（public）
- チームプロジェクト（user）
- サンプルユーザー（admin、developer1、developer2）
- サンプルAPIキー（各権限レベル）
- サンプルルール（各言語、アクセスレベル）

## 設定の優先順位

1. **環境変数** (最高優先度)
2. **設定ファイル** (auth.yaml)
3. **デフォルト値** (最低優先度)

## 環境別設定

### 開発環境
```bash
ENVIRONMENT=development
DEV_MODE=true
DEV_SKIP_AUTH=false
AUTH_REQUIRE_HTTPS=false
LOG_LEVEL=debug
```

### 本番環境
```bash
ENVIRONMENT=production
DEV_MODE=false
DEV_SKIP_AUTH=false
AUTH_REQUIRE_HTTPS=true
LOG_LEVEL=warn
```

### テスト環境
```bash
ENVIRONMENT=testing
DEV_MODE=true
DEV_SKIP_AUTH=true
DEV_MOCK_DATA=true
LOG_LEVEL=info
```

## セキュリティ設定

### 本番環境での必須設定
```bash
AUTH_REQUIRE_HTTPS=true
AUTH_ENABLED=true
API_KEY_ENABLED=true
RATE_LIMIT_ENABLED=true
AUDIT_LOG_ENABLED=true
```

### 開発環境での推奨設定
```bash
AUTH_REQUIRE_HTTPS=false
DEV_MODE=true
LOG_LEVEL=debug
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:13000
```

## トラブルシューティング

### 設定ファイルの読み込みエラー
```bash
# 設定ファイルの存在確認
ls -la config/

# YAML構文チェック
yamllint config/auth.yaml

# 環境変数の確認
env | grep -E "(AUTH_|API_|MCP_)"
```

### 権限設定の確認
```bash
# 現在の権限レベル確認
curl -H "X-API-Key: your_key" http://localhost:18081/api/v1/auth/me

# 権限テスト
curl -H "X-API-Key: user_key" http://localhost:18081/api/v1/rules
```

### MCP接続の確認
```bash
# MCPエンドポイントのテスト
curl -X POST http://localhost:18081/mcp/request \
  -H "Content-Type: application/json" \
  -d '{"id":"test","method":"getRules","params":{"project_id":"web-app"}}'

# 認証付きMCPテスト
curl -X POST http://localhost:18081/mcp/request \
  -H "X-API-Key: your_key" \
  -H "Content-Type: application/json" \
  -d '{"id":"test","method":"createRule","params":{...}}'
```

## カスタマイズ

### 新しい権限レベルの追加
```yaml
# auth.yaml
access_levels:
  moderator:
    description: "Moderator access - limited admin permissions"
    permissions:
      - "read:*"
      - "create:rules"
      - "update:rules"
      - "delete:own_rules"
      - "moderate:content"
    mcp_access: true
    rate_limit_multiplier: 1.5
```

### カスタム権限の定義
```yaml
# auth.yaml
custom_permissions:
  - "manage:own_projects"
  - "invite:team_members"
  - "export:rules"
  - "import:rules"
```

### レート制限の調整
```yaml
# auth.yaml
rate_limit:
  custom_levels:
    premium_user:
      requests_per_minute: 500
      burst_limit: 100
    enterprise:
      requests_per_minute: 1000
      burst_limit: 200
```

## ベストプラクティス

1. **環境変数の使用**: 機密情報は環境変数で管理
2. **設定ファイルのバージョン管理**: 設定変更履歴を追跡
3. **段階的な権限付与**: 必要最小限の権限から開始
4. **定期的な監査**: 権限設定の見直しと更新
5. **セキュリティテスト**: 権限制限の動作確認
6. **ドキュメント化**: 設定変更の記録と説明

## サポート

設定に関する質問や問題がある場合は、以下の方法でサポートを受けることができます：

- **Issues**: GitHubのIssuesページで問題を報告
- **Discussions**: GitHubのDiscussionsページで質問・議論
- **Wiki**: 詳細な設定ガイドやチュートリアル
