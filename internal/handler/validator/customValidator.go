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

		matched, _ := regexp.MatchString(`^[0-9\-]+$`, phoneStr)
		return matched
	})

	return v
}
