package handler

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

type FieldErrors []FieldError

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

type FieldError struct {
	Field      string
	Validation string
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
			fieldErrors = append(fieldErrors, FieldError{
				Field:      validationErr.Field(),
				Validation: validationErr.Error(),
			})
		}
	}

	// passwords are validated using separate case from regex (not from validator, since it doesn't support regex validations)
	// pattern := regexp.MustCompile(`^(?=.*[A-Z])(?=.*\d)(?=.*\W).+$`)
	// if v.Password != "" && !pattern.MatchString(v.Password) {
	// 	fieldErrors = append(fieldErrors, FieldError{
	// 		Field:      "password",
	// 		Validation: "password must contain at least 1 capital characters, 1 number, and 1 special (nonalpha-numeric) characters",
	// 	})
	// }

	return fieldErrors
}

type AuthenticateUserValidator struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
	Password    string `json:"password" validate:"required"`
}

func (v AuthenticateUserValidator) Validate() FieldErrors {
	fieldErrors := FieldErrors{}

	err := validate.Struct(v)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)

		for _, validationErr := range validationErrors {
			fieldErrors = append(fieldErrors, FieldError{
				Field:      validationErr.Field(),
				Validation: validationErr.Error(),
			})
		}
	}

	return fieldErrors
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
			fieldErrors = append(fieldErrors, FieldError{
				Field:      validationErr.Field(),
				Validation: validationErr.Error(),
			})
		}
	}

	return fieldErrors
}
