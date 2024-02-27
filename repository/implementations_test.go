package repository_test

import (
	"context"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestRepository_CreateUser(t *testing.T) {
	type mockExec struct {
		data repository.CreateUserOutput
		err  error
	}
	type args struct {
		ctx   context.Context
		input repository.CreateUserInput
	}
	tests := []struct {
		name     string
		mockExec mockExec
		args     args
		want     repository.CreateUserOutput
		wantErr  bool
	}{
		{
			name: "successfully inserts data",
			mockExec: mockExec{
				data: repository.CreateUserOutput{
					ID: "abc123-def456",
				},
				err: nil,
			},
			args: args{
				ctx: context.Background(),
				input: repository.CreateUserInput{
					ID:             "abc123-def456",
					FullName:       "Test",
					PhoneNumber:    "+62812345567",
					HashedPassword: "test-hashed-password",
				},
			},
			want: repository.CreateUserOutput{
				ID: "abc123-def456",
			},
			wantErr: false,
		},
		{
			name: "error when inserting data",
			mockExec: mockExec{
				err: assert.AnError,
			},
			args: args{
				ctx: context.Background(),
				input: repository.CreateUserInput{
					ID:             "abc123-def456",
					FullName:       "Test",
					PhoneNumber:    "+62812345567",
					HashedPassword: "test-hashed-password",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, m, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatal("unexpected error")
			}
			defer db.Close()

			query := `
				INSERT INTO
					users
					(id, full_name, phone_number, hashed_password)
				VALUES
					($1, $2, $3, $4)
				RETURNING id
			`

			expectExec := m.ExpectQuery(query).WithArgs(tt.args.input.ID, tt.args.input.FullName, tt.args.input.PhoneNumber, tt.args.input.HashedPassword)
			if tt.mockExec.err != nil {
				expectExec.WillReturnError(tt.mockExec.err)
			} else {
				rows := sqlmock.NewRows([]string{"id"})
				mockRow := tt.mockExec.data
				rows.AddRow(mockRow.ID)

				expectExec.WillReturnRows(rows)
			}

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			r := repository.Repository{Db: sqlxDB}

			got, err := r.CreateUser(tt.args.ctx, tt.args.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestRepository_GetUserByPhoneNumber(t *testing.T) {
	type mockExec struct {
		data repository.UserOutput
		err  error
	}
	type args struct {
		ctx         context.Context
		phoneNumber string
	}
	tests := []struct {
		name     string
		mockExec mockExec
		args     args
		want     repository.UserOutput
		wantErr  bool
	}{
		{
			name: "successfully fetches a single user by phone number",
			mockExec: mockExec{
				data: repository.UserOutput{
					ID:          "abc123-def456",
					PhoneNumber: "+62812345567",
				},
				err: nil,
			},
			args: args{
				ctx:         context.Background(),
				phoneNumber: "+62812345567",
			},
			want: repository.UserOutput{
				ID:          "abc123-def456",
				PhoneNumber: "+62812345567",
			},
			wantErr: false,
		},
		{
			name: "error when fetching user by phone number",
			mockExec: mockExec{
				err: assert.AnError,
			},
			args: args{
				ctx:         context.Background(),
				phoneNumber: "+62812345567",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, m, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatal("unexpected error")
			}
			defer db.Close()

			query := `
				SELECT
					id,
					full_name,
					phone_number,
					hashed_password,
					login_count
				FROM
					users
				WHERE
					phone_number = $1
				LIMIT 1
			`

			expectExec := m.ExpectQuery(query).WithArgs(tt.args.phoneNumber)
			if tt.mockExec.err != nil {
				expectExec.WillReturnError(tt.mockExec.err)
			} else {
				rows := sqlmock.NewRows([]string{"id", "full_name", "phone_number", "hashed_password", "login_count"})
				mockRow := tt.mockExec.data
				rows.AddRow(mockRow.ID, mockRow.FullName, mockRow.PhoneNumber, mockRow.HashedPassword, mockRow.LoginCount)

				expectExec.WillReturnRows(rows)
			}

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			r := repository.Repository{Db: sqlxDB}

			got, err := r.GetUserByPhoneNumber(tt.args.ctx, tt.args.phoneNumber)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestRepository_GetUserByID(t *testing.T) {
	type mockExec struct {
		data repository.UserOutput
		err  error
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name     string
		mockExec mockExec
		args     args
		want     repository.UserOutput
		wantErr  bool
	}{
		{
			name: "successfully fetches a single user by id",
			mockExec: mockExec{
				data: repository.UserOutput{
					ID:          "abc123-def456",
					PhoneNumber: "+62812345567",
				},
				err: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  "abc123-def456",
			},
			want: repository.UserOutput{
				ID:          "abc123-def456",
				PhoneNumber: "+62812345567",
			},
			wantErr: false,
		},
		{
			name: "error when fetching user by id",
			mockExec: mockExec{
				err: assert.AnError,
			},
			args: args{
				ctx: context.Background(),
				id:  "abc123-def456",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, m, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatal("unexpected error")
			}
			defer db.Close()

			query := `
				SELECT
					id,
					full_name,
					phone_number,
					hashed_password,
					login_count
				FROM
					users
				WHERE
					id = $1
				LIMIT 1
			`

			expectExec := m.ExpectQuery(query).WithArgs(tt.args.id)
			if tt.mockExec.err != nil {
				expectExec.WillReturnError(tt.mockExec.err)
			} else {
				rows := sqlmock.NewRows([]string{"id", "full_name", "phone_number", "hashed_password", "login_count"})
				mockRow := tt.mockExec.data
				rows.AddRow(mockRow.ID, mockRow.FullName, mockRow.PhoneNumber, mockRow.HashedPassword, mockRow.LoginCount)

				expectExec.WillReturnRows(rows)
			}

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			r := repository.Repository{Db: sqlxDB}

			got, err := r.GetUserByID(tt.args.ctx, tt.args.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestRepository_UpdateUser(t *testing.T) {
	type mockExec struct {
		err          error
		affectedRows int
	}
	type args struct {
		ctx   context.Context
		id    string
		input repository.UpdateUserInput
	}
	tests := []struct {
		name     string
		mockExec mockExec
		args     args
		wantErr  bool
	}{
		{
			name: "successfully updates a single user by id",
			mockExec: mockExec{
				err:          nil,
				affectedRows: 1,
			},
			args: args{
				ctx: context.Background(),
				id:  "abc123-def456",
				input: repository.UpdateUserInput{
					FullName:    "full_name",
					PhoneNumber: "+62812345567",
				},
			},
			wantErr: false,
		},
		{
			name: "error when update user by id",
			mockExec: mockExec{
				err: assert.AnError,
			},
			args: args{
				ctx: context.Background(),
				id:  "abc123-def456",
				input: repository.UpdateUserInput{
					FullName:    "full_name",
					PhoneNumber: "+62812345567",
				},
			},
			wantErr: true,
		},
		{
			name: "error when update user due to affected rows is 0",
			mockExec: mockExec{
				affectedRows: 0,
			},
			args: args{
				ctx: context.Background(),
				id:  "abc123-def456",
				input: repository.UpdateUserInput{
					FullName:    "full_name",
					PhoneNumber: "+62812345567",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, m, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatal("unexpected error")
			}
			defer db.Close()

			query := `
				UPDATE
					users
				SET
					full_name = $1,
					phone_number = $2
				WHERE
					id = $3
			`

			expectExec := m.ExpectExec(query).WithArgs(tt.args.input.FullName, tt.args.input.PhoneNumber, tt.args.id)
			if tt.mockExec.err != nil {
				expectExec.WillReturnError(tt.mockExec.err)
			} else {
				expectExec.WillReturnResult(sqlmock.NewResult(0, int64(tt.mockExec.affectedRows)))
			}

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			r := repository.Repository{Db: sqlxDB}

			err = r.UpdateUser(tt.args.ctx, tt.args.id, tt.args.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRepository_IncrementLoginCount(t *testing.T) {
	type mockExec struct {
		err          error
		affectedRows int
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name     string
		mockExec mockExec
		args     args
		wantErr  bool
	}{
		{
			name: "successfully increment single user login count by id",
			mockExec: mockExec{
				err:          nil,
				affectedRows: 1,
			},
			args: args{
				ctx: context.Background(),
				id:  "abc123-def456",
			},
			wantErr: false,
		},
		{
			name: "error when increment user login count by id",
			mockExec: mockExec{
				err: assert.AnError,
			},
			args: args{
				ctx: context.Background(),
				id:  "abc123-def456",
			},
			wantErr: true,
		},
		{
			name: "error when increment user login count due to affected rows is 0",
			mockExec: mockExec{
				affectedRows: 0,
			},
			args: args{
				ctx: context.Background(),
				id:  "abc123-def456",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, m, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatal("unexpected error")
			}
			defer db.Close()

			query := `
				UPDATE
					users
				SET
					login_count = login_count + 1
				WHERE
					id = $1
			`

			expectExec := m.ExpectExec(query).WithArgs(tt.args.id)
			if tt.mockExec.err != nil {
				expectExec.WillReturnError(tt.mockExec.err)
			} else {
				expectExec.WillReturnResult(sqlmock.NewResult(0, int64(tt.mockExec.affectedRows)))
			}

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			r := repository.Repository{Db: sqlxDB}

			err = r.IncrementLoginCount(tt.args.ctx, tt.args.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
