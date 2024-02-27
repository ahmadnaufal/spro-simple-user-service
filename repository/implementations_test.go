package repository

import (
	"context"
	"reflect"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestRepository_CreateUser(t *testing.T) {
	type fields struct {
		Db *sqlx.DB
	}
	type args struct {
		ctx   context.Context
		input CreateUserInput
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantOutput CreateUserOutput
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				Db: tt.fields.Db,
			}
			gotOutput, err := r.CreateUser(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("Repository.CreateUser() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}
