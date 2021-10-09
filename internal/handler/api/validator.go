package api

import (
	"github.com/go-playground/validator/v10"
	"unicode"
)

func ValidateStruct(data interface{}) error {
	validate := validator.New()
	validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		return isValidPassword(password)
	})
	err := validate.Struct(data)
	return err
}

func isValidPassword(password string) bool {
	if len(password) < 6 {
		return false
	}
	var hasNumber, hasUpper, hasSpecial, hasLetter bool
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		case unicode.IsLetter(c) || c == ' ':
			hasLetter = true
		}
	}
	return hasNumber && hasSpecial && hasUpper && hasLetter
}
