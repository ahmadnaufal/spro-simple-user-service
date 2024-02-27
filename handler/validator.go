package handler

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

type FieldErrors []generated.FieldError

func (fe FieldErrors) Error() string {
	if len(fe) < 1 {
		return ""
	}

	errStrings := []string{}
	for _, e := range fe {
		errStrings = append(errStrings, fmt.Sprintf("%s: %s", e.Field, e.Validation))
	}

	return strings.Join(errStrings, "; ")
}

func validationMessages(tag, param string) string {
	switch tag {
	case "required":
		return "field is required"
	case "min":
		return fmt.Sprintf("length is less than minimum allowed length of %s", param)
	case "max":
		return fmt.Sprintf("length is more than maximum allowed length of %s", param)
	case "startswith":
		return fmt.Sprintf("value should begin with %s", param)
	default:
		return tag
	}
}

type RegisterUserValidator struct {
	FullName    string `json:"full_name" validate:"required,min=3,max=60"`
	PhoneNumber string `json:"phone_number" validate:"required,min=10,max=13,startswith=+62"`
	Password    string `json:"password" validate:"required,min=6,max=64"`
}

func (v RegisterUserValidator) Validate() FieldErrors {
	fieldErrors := FieldErrors{}

	err := validate.Struct(v)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)

		for _, validationErr := range validationErrors {
			fieldErrors = append(fieldErrors, generated.FieldError{
				Field:      validationErr.Field(),
				Validation: validationMessages(validationErr.Tag(), validationErr.Param()),
			})
		}
	}

	// passwords are validated using separate case from regex (not from validator, since it doesn't support regex validations)
	if v.Password != "" && !isPasswordValid(v.Password) {
		fieldErrors = append(fieldErrors, generated.FieldError{
			Field:      "Password",
			Validation: "must contain at least 1 capital characters, 1 number, and 1 special (nonalpha-numeric) characters",
		})
	}

	return fieldErrors
}

// helper function to check if password met the following criteria:
// containing at least 1 digit, 1 symbol, and 1 uppercase letter.
// this is done since regexp golang does not support lookaheads.
func isPasswordValid(password string) bool {
	var (
		containUppercase = false
		containSymbol    = false
		containDigit     = false
	)

	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			containDigit = true
		case unicode.IsUpper(c):
			containUppercase = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			containSymbol = true
		}
	}

	return containUppercase && containDigit && containSymbol
}

type AuthenticateUserValidator struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
	Password    string `json:"password" validate:"required"`
}

type UpdateUserValidator struct {
	FullName    *string `json:"full_name" validate:"omitempty,min=3,max=60"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=13,startswith=+62"`
}

func (v UpdateUserValidator) Validate() FieldErrors {
	fieldErrors := FieldErrors{}

	err := validate.Struct(v)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)

		for _, validationErr := range validationErrors {
			fieldErrors = append(fieldErrors, generated.FieldError{
				Field:      validationErr.Field(),
				Validation: validationErr.Error(),
			})
		}
	}

	return fieldErrors
}
