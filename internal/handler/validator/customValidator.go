package validator

import (
	"reflect"
	"regexp"

	"github.com/go-playground/validator/v10"
)

func NewValidator() *validator.Validate {
	v := validator.New()

	v.RegisterValidation("phone", func(fl validator.FieldLevel) bool {

		field := fl.Field()

		var phoneStr string

		// ポインタの場合
		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				return true
			}
			phoneStr = field.Elem().String()
		} else {
			// 値の場合
			phoneStr = field.String()
		}

		if phoneStr == "" {
			return true
		}
		// 固定電話と携帯電話
		matched, _ := regexp.MatchString(`^(\d{4}-\d{2}-\d{4}|\d{3}-\d{4}-\d{4})$`, phoneStr)
		return matched
	})

	return v
}
