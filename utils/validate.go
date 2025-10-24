package util

import (
	"fmt"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func PasswordValidator(fl validator.FieldLevel) bool {
	return fl.Field().String() == "custom"
}

func RegisterValidation() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("password", PasswordValidator)
	}
}

func GetValidateErrors(err error) []string {
	errs := err.(validator.ValidationErrors)
	out := make([]string, len(errs))
	for i, fe := range errs {
		out[i] = customErrorMessage(fe)
	}
	return out
}

func customErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required.", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long.", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s must be no more than %s characters long.", fe.Field(), fe.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email address.", fe.Field())
	case "password":
		return fmt.Sprintf("%s must include at least one uppercase letter, one lowercase letter, and one number.", fe.Field())
	default:
		return fmt.Sprintf("%s is not valid.", fe.Field())
	}

}
