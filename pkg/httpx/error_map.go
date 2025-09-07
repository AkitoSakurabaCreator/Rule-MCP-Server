package httpx

import (
	"errors"
	"net/http"

	"github.com/AkitoSakurabaCreator/Rule-MCP-Server/pkg/apperr"
)

func MapAppError(err error) (int, string, string, interface{}) {
	if err == nil {
		return http.StatusOK, "", "", nil
	}

	var wd *apperr.WithDetails
	if errors.As(err, &wd) {
		switch {
		case errors.Is(err, apperr.ErrValidation):
			return http.StatusBadRequest, CodeValidation, wd.Message, wd.Details
		case errors.Is(err, apperr.ErrUnauthorized):
			return http.StatusUnauthorized, CodeUnauthorized, wd.Message, wd.Details
		case errors.Is(err, apperr.ErrForbidden):
			return http.StatusForbidden, CodeForbidden, wd.Message, wd.Details
		case errors.Is(err, apperr.ErrNotFound):
			return http.StatusNotFound, CodeNotFound, wd.Message, wd.Details
		case errors.Is(err, apperr.ErrConflict):
			return http.StatusConflict, CodeConflict, wd.Message, wd.Details
		case errors.Is(err, apperr.ErrUnprocessable):
			return http.StatusUnprocessableEntity, CodeUnprocessable, wd.Message, wd.Details
		}
	}

	switch {
	case errors.Is(err, apperr.ErrValidation):
		return http.StatusBadRequest, CodeValidation, "入力値が不正です", nil
	case errors.Is(err, apperr.ErrUnauthorized):
		return http.StatusUnauthorized, CodeUnauthorized, "認証が必要です", nil
	case errors.Is(err, apperr.ErrForbidden):
		return http.StatusForbidden, CodeForbidden, "権限がありません", nil
	case errors.Is(err, apperr.ErrNotFound):
		return http.StatusNotFound, CodeNotFound, "対象が見つかりません", nil
	case errors.Is(err, apperr.ErrConflict):
		return http.StatusConflict, CodeConflict, "リソースが競合しています", nil
	case errors.Is(err, apperr.ErrUnprocessable):
		return http.StatusUnprocessableEntity, CodeUnprocessable, "処理できない内容です", nil
	default:
		return http.StatusInternalServerError, CodeInternal, "サーバ内部でエラーが発生しました", nil
	}
}
