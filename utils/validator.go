package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ExtractValidationError(req interface{}) error {
	var message error
	var v = validator.New()
	// get json tag
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}

		return name
	})

	err := v.Struct(req)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var e error
			switch err.Tag() {
			case "required":
				e = fmt.Errorf("'%s' harus diisi", err.Field())
			case "email":
				e = fmt.Errorf("'%s' harus mengikuti format email", err.Field())
			case "len":
				e = fmt.Errorf("panjang '%s'  %v karakter", err.Field(), err.Param())
			case "datetime":
				e = fmt.Errorf("'%s' harus mengikuti format %v ", err.Field(), err.Param())
			default:
				e = fmt.Errorf("'%s': '%v' harus berkriteria '%s' '%v' ", err.Field(), err.Value(), err.Tag(), err.Param())
			}

			message = e
		}

		return message
	}

	return nil
}
