package helpers

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator"
)

// validatePassword uses regex to enforce length ≥8, symbols, digits, upper and lower case letters.
func validatePassword(fl validator.FieldLevel) bool {
	pwd := fl.Field().String()
	if len(pwd) < 8 {
		return false
	}
	// each of these must match at least once
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(pwd)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(pwd)
	hasDigit := regexp.MustCompile(`\d`).MatchString(pwd)
	// common symbols—adjust the character class to your needs
	hasSymbol := regexp.MustCompile(`[!@#~\$%\^&\*\(\)\-_\+=\[\]\{\}|\\;:'",<\.>\/\?]`).MatchString(pwd)

	return hasUpper && hasLower && hasDigit && hasSymbol
}

func Init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("pwd", validatePassword)
	}
}
