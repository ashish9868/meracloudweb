package helpers

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func PrepareValidationError(e error) map[string]string {
	list := make(map[string]string)
	for _, err := range e.(validator.ValidationErrors) {
		field := strings.ToLower(err.Field())
		list[field] = field + " field is " + err.Tag() + "."
		fmt.Println(err.Namespace()) // can differ when a custom TagNameFunc is registered or
		fmt.Println(err.Field())     // by passing alt name to ReportError like below
		fmt.Println(err.StructNamespace())
		fmt.Println(err.StructField())
		fmt.Println(err.Tag())
		fmt.Println(err.ActualTag())
		fmt.Println(err.Kind())
		fmt.Println(err.Type())
		fmt.Println(err.Value())
		fmt.Println(err.Param())
		fmt.Println()
	}
	return list
}
