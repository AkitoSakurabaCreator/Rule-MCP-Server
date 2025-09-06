# Contributing / コントリビュート

本プロジェクトへの貢献にあたっては、以下のガイドラインに従ってください。

## 開発環境 / Dev Setup
- Go 1.21+ / Node.js 18+
- フロントエンドは pnpm 推奨
- Docker / docker compose
- 起動: `docker compose up -d`

## Pull Request
- fork → feature ブランチ作成 → PR
- Lint・テストを通してからPR
- コミットは Conventional Commits を推奨

## コードスタイル / Code Style
- Go: gofmt, go vet（staticcheck 推奨）
- Frontend: ESLint + Prettier, TypeScript strict

## テスト / Tests
- ハンドラ（users/roles/languages/rule-options/settings）と権限の網羅
- 失敗時のエラーレスポンス（code/message/requestId/timestamp）も検証

## セキュリティ / Security
- 秘密情報はコミット禁止（.env と GitHub Secrets を使用）
- 脆弱性は SECURITY.md の手順で報告

## ライセンス / License
- 貢献内容は MIT License の下で提供されます
