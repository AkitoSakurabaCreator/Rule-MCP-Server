package mcpx

import (
	"errors"

	"github.com/AkitoSakurabaCreator/Rule-MCP-Server/pkg/apperr"
)

const (
	CodeValidation    = 4000
	CodeUnauthorized  = 4001
	CodeForbidden     = 4003
	CodeNotFound      = 4040
	CodeConflict      = 4090
	CodeUnprocessable = 4220
	CodeInternal      = 5000
)

func MapAppErrorToMCP(err error) (int, string) {
	if err == nil {
		return 0, ""
	}
	if errors.Is(err, apperr.ErrValidation) {
		return CodeValidation, "入力値が不正です"
	}
	if errors.Is(err, apperr.ErrUnauthorized) {
		return CodeUnauthorized, "認証が必要です"
	}
	if errors.Is(err, apperr.ErrForbidden) {
		return CodeForbidden, "権限がありません"
	}
	if errors.Is(err, apperr.ErrNotFound) {
		return CodeNotFound, "対象が見つかりません"
	}
	if errors.Is(err, apperr.ErrConflict) {
		return CodeConflict, "リソースが競合しています"
	}
	if errors.Is(err, apperr.ErrUnprocessable) {
		return CodeUnprocessable, "処理できない内容です"
	}
	return CodeInternal, "サーバ内部でエラーが発生しました"
}
