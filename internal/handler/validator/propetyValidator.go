package validator

import (
	"strconv"
	"strings"

	"github.com/yamu-studio/profact-simulated-practical-go/internal/handler/request"
)

func ParseInt(s string) int {
	val, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return 0
	}
	return val
}
func ValidateCSRow(row []string, lineNumber int) []request.ValidationError {
	var errors []request.ValidationError
	if len(row) < 5 {
		errors = append(errors, request.ValidationError{
			Row:     lineNumber,
			Field:   "row",
			Message: "カラム数が不足しています",
		})
	}

	name := strings.TrimSpace(row[0])
	if name == "" {
		errors = append(errors, request.ValidationError{
			Row:     lineNumber,
			Field:   "name",
			Message: "物件名は必須です",
		})
	} else if len(name) > 255 {
		errors = append(errors, request.ValidationError{
			Row:     lineNumber,
			Field:   "name",
			Message: "物件名は255文字以内で入力してください",
		})
	}

	rentStr := strings.TrimSpace(row[1])
	if rentStr == "" {
		errors = append(errors, request.ValidationError{
			Row:     lineNumber,
			Field:   "rent",
			Message: "賃料は必須です",
		})
	} else {
		rent, err := strconv.Atoi(rentStr)
		if err != nil {
			errors = append(errors, request.ValidationError{
				Row:     lineNumber,
				Field:   "rent",
				Message: "賃料は数値で入力してください",
			})
		} else if rent < 0 {
			errors = append(errors, request.ValidationError{
				Row:     lineNumber,
				Field:   "rent",
				Message: "賃料は0以上の数値を入力してください",
			})
		} else if rent > 1000000000 {
			errors = append(errors, request.ValidationError{
				Row:     lineNumber,
				Field:   "rent",
				Message: "賃料が大きすぎます",
			})
		}
	}

	address := strings.TrimSpace(row[2])
	if address == "" {
		errors = append(errors, request.ValidationError{
			Row:     lineNumber,
			Field:   "address",
			Message: "住所は必須です",
		})
	} else if len(address) > 255 {
		errors = append(errors, request.ValidationError{
			Row:     lineNumber,
			Field:   "address",
			Message: "住所は255文字以内で入力してください",
		})
	}

	layout := strings.TrimSpace(row[3])
	if len(layout) > 50 {
		errors = append(errors, request.ValidationError{
			Row:     lineNumber,
			Field:   "layout",
			Message: "間取りは50文字以内で入力してください",
		})
	}
	status := strings.TrimSpace(row[4])
	validStatuses := map[string]bool{
		"available": true, "contracted": true, "hidden": true,
	}
	if status != "" && !validStatuses[status] {
		errors = append(errors, request.ValidationError{
			Row:     lineNumber,
			Field:   "status",
			Message: "ステータスが不正です（有効: available/contracted/hidden）",
		})
	}

	return errors

}
