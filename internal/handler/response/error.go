package response

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type ErrorResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Errors  []ValidationError `json:"errors,omitempty"`
}

func FormatValidationErrors(err error) []ValidationError {
	var errors []ValidationError

	validationErrors, ok := err.(validator.ValidationErrors)

	if !ok {
		return nil
	}

	for _, fe := range validationErrors {
		var message string
		switch fe.Tag() {
		case "required":
			message = "必須項目です"
		case "email":
			message = "有効なメールアドレス形式で入力してください"
		case "phone":
			message = "電話番号が正しくありません"
		case "max":
			message = fmt.Sprintf("%s文字以内で入力してください", fe.Param())
		case "len":
			if fe.Param() == "36" {
				message = "不正なID形式です"
			} else {
				message = fmt.Sprintf("%s文字で入力してください", fe.Param())
			}
		case "oneof":
			message = "有効な値を指定してください"
		case "min":
			message = fmt.Sprintf("%s文字以上で入力してください", fe.Param())
		case "datetime":
			message = "正しい日付形式で入力してください"
		default:
			message = "入力内容に誤りがあります"
		}

		errors = append(errors, ValidationError{
			Field: fe.Field(),
			Error: message,
		})
	}

	return errors
}
