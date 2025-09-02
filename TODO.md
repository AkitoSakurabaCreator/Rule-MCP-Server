# Rule MCP Server 開発 TODO

## 1. 基本設計
- [ ] サーバの役割を定義する
  - 各プロジェクトやユーザーに対してルールを配布
  - MCP 経由で Cursor / Cline に統制ルールを注入
- [ ] 管理するルールの種類を決める
  - コーディング規約（lint / prettier / naming）
  - 使用禁止ライブラリ一覧
  - セキュリティ制限（env参照禁止、直SQL禁止など）
  - チーム固有ルール（例: AWSリソースはCDK経由）
- [ ] API 設計
  - `GET /rules?project=xxx` : プロジェクトに適用するルールを返す
  - `POST /rules/validate` : ルールに従ってコードや提案を検証
  - `GET /health` : サーバ状態チェック

## 2. 技術スタック
- [ ] Node.js + Express または Go + Gin で API 実装
- [ ] データストアを決定（JSON/YAML, PostgreSQL, etc.）
- [ ] MCP SDK を利用してサーバ化
- [ ] 認証認可
  - APIキー or GitHub OAuth or 社内SSO

## 3. 機能実装
- [ ] ルール管理
  - プロジェクトごと／ユーザーごとのルールを保存
  - 継承ルール（共通ルール＋個別ルール）
- [ ] MCP ハンドラ
  - Cursor/Cline からのリクエストを処理
  - ルールをプロンプトに注入できる形式で返却
- [ ] ログと監査
  - 誰がいつどのルールを取得したか記録

## 4. クライアント側対応
- [ ] Cursor / Cline の設定に `ruleServer` を追加
- [ ] ルール取得リクエストを行い、システムプロンプトに反映
- [ ] ルール検証 API を呼び出し、AI 提案のバリデーションを行う

## 5. 運用
- [ ] CI/CD に組み込み
- [ ] ルール更新の通知（Webhook or Slack連携）
- [ ] バージョン管理（ルールの履歴を追えるように）

---
