package util

import (
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func NoSpaceValidator(fl validator.FieldLevel) bool {
	match, _ := regexp.MatchString(`\s`, fl.Field().String())
	return !match
}

func NoSpecialValidator(fl validator.FieldLevel) bool {
	match, _ := regexp.MatchString(`[^a-zA-Z0-9]`, fl.Field().String())
	return !match
}

func HaveLowerValidator(fl validator.FieldLevel) bool {
	match, _ := regexp.MatchString(`[a-z]+`, fl.Field().String())
	return match
}

func HaveUpperValidator(fl validator.FieldLevel) bool {
	match, _ := regexp.MatchString(`[A-Z]+`, fl.Field().String())
	return match
}

func HaveNumberValidator(fl validator.FieldLevel) bool {
	match, _ := regexp.MatchString(`[0-9]+`, fl.Field().String())
	return match
}

func RegisterValidation() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("nospace", NoSpaceValidator)
		v.RegisterValidation("nospecial", NoSpecialValidator)
		v.RegisterValidation("havelower", HaveLowerValidator)
		v.RegisterValidation("haveupper", HaveUpperValidator)
		v.RegisterValidation("havenumber", HaveNumberValidator)
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
	case "nospace":
		return fmt.Sprintf("%s must not contain any spaces.", fe.Field())
	case "nospecial":
		return fmt.Sprintf("%s must not contain special characters.", fe.Field())
	case "havelower":
		return fmt.Sprintf("%s must include at least one lowercase letter.", fe.Field())
	case "haveupper":
		return fmt.Sprintf("%s must include at least one uppercase letter.", fe.Field())
	case "havenumber":
		return fmt.Sprintf("%s must include at least one number.", fe.Field())
	default:
		return fmt.Sprintf("%s is not valid.", fe.Field())
	}

}
