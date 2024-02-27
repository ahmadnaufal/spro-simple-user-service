package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterUserValidator_Validate(t *testing.T) {
	type fields struct {
		FullName    string
		PhoneNumber string
		Password    string
	}
	tests := []struct {
		name   string
		fields fields
		want   FieldErrors
	}{
		{
			name: "successfully validates all field",
			fields: fields{
				FullName:    "john smith",
				PhoneNumber: "+62812345678",
				Password:    "TestingNewUser123!",
			},
			want: FieldErrors{},
		},
		{
			name: "full name is not present",
			fields: fields{
				PhoneNumber: "+62812345678",
				Password:    "TestingNewUser123!",
			},
			want: FieldErrors{
				{
					Field:      "FullName",
					Validation: "field is required",
				},
			},
		},
		{
			name: "full name is less than minimum",
			fields: fields{
				FullName:    "su",
				PhoneNumber: "+62812345678",
				Password:    "TestingNewUser123!",
			},
			want: FieldErrors{
				{
					Field:      "FullName",
					Validation: "length is less than minimum allowed length of 3",
				},
			},
		},
		{
			name: "full name is more than maximum",
			fields: fields{
				FullName:    "this is long string which results in length more than 60. what do you guess the result will be?",
				PhoneNumber: "+62812345678",
				Password:    "TestingNewUser123!",
			},
			want: FieldErrors{
				{
					Field:      "FullName",
					Validation: "length is more than maximum allowed length of 60",
				},
			},
		},
		{
			name: "phone number is not present",
			fields: fields{
				FullName: "john smith",
				Password: "TestingNewUser123!",
			},
			want: FieldErrors{
				{
					Field:      "PhoneNumber",
					Validation: "field is required",
				},
			},
		},
		{
			name: "phone number is less than 10",
			fields: fields{
				FullName:    "john smith",
				PhoneNumber: "+62812345",
				Password:    "TestingNewUser123!",
			},
			want: FieldErrors{
				{
					Field:      "PhoneNumber",
					Validation: "length is less than minimum allowed length of 10",
				},
			},
		},
		{
			name: "phone number is more than 13",
			fields: fields{
				FullName:    "john smith",
				PhoneNumber: "+6281234567899999",
				Password:    "TestingNewUser123!",
			},
			want: FieldErrors{
				{
					Field:      "PhoneNumber",
					Validation: "length is more than maximum allowed length of 13",
				},
			},
		},
		{
			name: "phone number does not start with +62",
			fields: fields{
				FullName:    "john smith",
				PhoneNumber: "8123456789",
				Password:    "TestingNewUser123!",
			},
			want: FieldErrors{
				{
					Field:      "PhoneNumber",
					Validation: "value should begin with +62",
				},
			},
		},
		{
			name: "password does not exist",
			fields: fields{
				FullName:    "john smith",
				PhoneNumber: "+628123456789",
			},
			want: FieldErrors{
				{
					Field:      "Password",
					Validation: "field is required",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := RegisterUserValidator{
				FullName:    tt.fields.FullName,
				PhoneNumber: tt.fields.PhoneNumber,
				Password:    tt.fields.Password,
			}

			assert.Equal(t, tt.want, v.Validate())
		})
	}
}
