package handler

import (
	"reflect"
	"testing"
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := RegisterUserValidator{
				FullName:    tt.fields.FullName,
				PhoneNumber: tt.fields.PhoneNumber,
				Password:    tt.fields.Password,
			}
			if got := v.Validate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegisterUserValidator.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
