package apperr

// 認証関連のエラーメッセージ
const (
	MsgInvalidCredentials      = "認証情報が正しくありません"
	MsgUserNotAuthenticated    = "認証されていません"
	MsgInvalidToken            = "トークンが無効です"
	MsgTokenGenerationFailed   = "トークン生成に失敗しました"
	MsgAccountPendingApproval  = "アカウントは管理者の承認待ちです。承認後にログインできます。"
	MsgCurrentPasswordWrong    = "現在のパスワードが正しくありません"
	MsgPasswordChanged         = "パスワードを変更しました"
	MsgPasswordHashFailed      = "パスワードのハッシュ化に失敗しました"
	MsgPasswordUpdateFailed    = "パスワードの更新に失敗しました"
)

// ユーザー関連のエラーメッセージ
const (
	MsgUserNotFound          = "ユーザーが見つかりません"
	MsgUserCreationFailed    = "ユーザーの作成に失敗しました"
	MsgUsernameExists        = "このユーザー名は既に使用されています"
	MsgEmailExists           = "このメールアドレスは既に登録されています"
	MsgAccountCreated        = "アカウントが作成されました。管理者の承認をお待ちください。"
	MsgUserApproved          = "ユーザーを承認しました"
	MsgUserRejected          = "ユーザーを拒否しました"
	MsgUserApprovalFailed    = "ユーザー承認処理に失敗しました"
	MsgUserListFailed        = "ユーザー一覧の取得に失敗しました"
)

// 権限関連のエラーメッセージ
const (
	MsgAdminRequired = "管理者権限が必要です"
	MsgForbidden     = "権限がありません"
)

// データベース関連のエラーメッセージ
const (
	MsgDBConnectionNotEstablished = "データベース接続が確立されていません。環境変数を確認してください。"
	MsgDBError                    = "データベースエラーが発生しました。しばらく待ってから再度お試しください。"
	MsgNotFound                   = "対象が見つかりません"
)

// バリデーション関連のエラーメッセージ
const (
	MsgInvalidRequest   = "リクエストデータが不正です"
	MsgValidationFailed = "入力値が不正です"
)

// サーバー関連のエラーメッセージ
const (
	MsgInternalError = "サーバ内部でエラーが発生しました"
)
