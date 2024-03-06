package handler_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"testing"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/testutil"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// matchers
type createUserInputMatcher struct {
	expected repository.CreateUserInput
}

func (m createUserInputMatcher) Matches(x interface{}) bool {
	actualInput, ok := x.(repository.CreateUserInput)
	if !ok {
		return false
	}

	return actualInput.ID != "" &&
		m.expected.FullName == actualInput.FullName &&
		m.expected.PhoneNumber == actualInput.PhoneNumber &&
		actualInput.HashedPassword != ""
}

func (m createUserInputMatcher) String() string {
	return fmt.Sprintf("{CreateUserInput - ID:%s FullName:%s PhoneNumber:%s}", m.expected.ID, m.expected.FullName, m.expected.PhoneNumber)
}

func TestServer_RegisterUser(t *testing.T) {
	type fields struct {
		Repository repository.RepositoryInterface
		JWT        handler.JWT
	}
	type args struct {
		payload map[string]interface{}
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
	}{
		{
			name: "successfully register new user",
			fields: fields{
				Repository: func() *repository.MockRepositoryInterface {
					ctrl := gomock.NewController(t)
					mockRepo := repository.NewMockRepositoryInterface(ctrl)

					expected := repository.CreateUserInput{
						FullName:    "test",
						PhoneNumber: "+62812345678",
					}

					mockRepo.EXPECT().
						GetUserByPhoneNumber(gomock.Any(), expected.PhoneNumber).
						Return(repository.UserOutput{}, sql.ErrNoRows)

					mockRepo.EXPECT().
						CreateUser(gomock.Any(), createUserInputMatcher{expected}).
						Return(repository.CreateUserOutput{
							ID: "test-user-id",
						}, nil)

					return mockRepo
				}(),
				JWT: handler.JWT{},
			},
			args: args{
				payload: map[string]interface{}{
					"full_name":    "test",
					"phone_number": "+62812345678",
					"password":     "Enter123!",
				},
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "failed to parse payload",
			fields: fields{
				Repository: func() *repository.MockRepositoryInterface {
					ctrl := gomock.NewController(t)
					mockRepo := repository.NewMockRepositoryInterface(ctrl)
					return mockRepo
				}(),
				JWT: handler.JWT{},
			},
			args: args{
				payload: map[string]interface{}{
					"full_name":    123,
					"phone_number": "+62812345678",
					"password":     "Enter123!",
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "some fields failed validation",
			fields: fields{
				Repository: func() *repository.MockRepositoryInterface {
					ctrl := gomock.NewController(t)
					mockRepo := repository.NewMockRepositoryInterface(ctrl)
					return mockRepo
				}(),
				JWT: handler.JWT{},
			},
			args: args{
				payload: map[string]interface{}{
					"full_name":    "test",
					"phone_number": "12345678",
					"password":     "Enter123!",
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "error phone number already registered",
			fields: fields{
				Repository: func() *repository.MockRepositoryInterface {
					ctrl := gomock.NewController(t)
					mockRepo := repository.NewMockRepositoryInterface(ctrl)

					expected := repository.CreateUserInput{
						FullName:    "test",
						PhoneNumber: "+62812345678",
					}

					mockRepo.EXPECT().
						GetUserByPhoneNumber(gomock.Any(), expected.PhoneNumber).
						Return(repository.UserOutput{ID: "existing-user", PhoneNumber: "+62812345678"}, nil)

					return mockRepo
				}(),
				JWT: handler.JWT{},
			},
			args: args{
				payload: map[string]interface{}{
					"full_name":    "test",
					"phone_number": "+62812345678",
					"password":     "Enter123!",
				},
			},
			wantStatus: http.StatusConflict,
		},
		{
			name: "error when creating user",
			fields: fields{
				Repository: func() *repository.MockRepositoryInterface {
					ctrl := gomock.NewController(t)
					mockRepo := repository.NewMockRepositoryInterface(ctrl)

					expected := repository.CreateUserInput{
						FullName:    "test",
						PhoneNumber: "+62812345678",
					}

					mockRepo.EXPECT().
						GetUserByPhoneNumber(gomock.Any(), expected.PhoneNumber).
						Return(repository.UserOutput{}, sql.ErrNoRows)

					mockRepo.EXPECT().
						CreateUser(gomock.Any(), createUserInputMatcher{expected}).
						Return(repository.CreateUserOutput{}, assert.AnError)

					return mockRepo
				}(),
				JWT: handler.JWT{},
			},
			args: args{
				payload: map[string]interface{}{
					"full_name":    "test",
					"phone_number": "+62812345678",
					"password":     "Enter123!",
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			s := handler.NewServer(handler.NewServerOptions{
				Repository: tt.fields.Repository,
				JWT:        tt.fields.JWT,
			})

			generated.RegisterHandlers(e, s)

			response := testutil.NewRequest().Post("/users").
				WithAcceptJson().
				WithJsonBody(tt.args.payload).
				GoWithHTTPHandler(t, e)

			assert.Equal(t, tt.wantStatus, response.Code())
		})
	}
}

func TestServer_AuthenticateUser(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpass!"), bcrypt.DefaultCost)
	user := repository.UserOutput{
		ID:             "abc123-def456",
		FullName:       "test",
		PhoneNumber:    "+62812345678",
		HashedPassword: string(hashedPassword),
	}

	type fields struct {
		Repository repository.RepositoryInterface
		JWT        handler.JWT
	}
	type args struct {
		payload map[string]interface{}
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
	}{
		{
			name: "successfully logs in",
			fields: fields{
				Repository: func() *repository.MockRepositoryInterface {
					ctrl := gomock.NewController(t)
					mockRepo := repository.NewMockRepositoryInterface(ctrl)

					mockRepo.EXPECT().
						GetUserByPhoneNumber(gomock.Any(), user.PhoneNumber).
						Return(user, nil)

					mockRepo.EXPECT().
						IncrementLoginCount(gomock.Any(), user.ID).
						Return(nil)

					return mockRepo
				}(),
				JWT: handler.NewJWT([]byte(RsaPublicKey), []byte(RsaPrivateKey)),
			},
			args: args{
				payload: map[string]interface{}{
					"phone_number": "+62812345678",
					"password":     "testpass!",
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "user with the phone number not found",
			fields: fields{
				Repository: func() *repository.MockRepositoryInterface {
					ctrl := gomock.NewController(t)
					mockRepo := repository.NewMockRepositoryInterface(ctrl)

					mockRepo.EXPECT().
						GetUserByPhoneNumber(gomock.Any(), user.PhoneNumber).
						Return(repository.UserOutput{}, sql.ErrNoRows)

					return mockRepo
				}(),
				JWT: handler.NewJWT([]byte(RsaPublicKey), []byte(RsaPrivateKey)),
			},
			args: args{
				payload: map[string]interface{}{
					"phone_number": "+62812345678",
					"password":     "testpass!",
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid password",
			fields: fields{
				Repository: func() *repository.MockRepositoryInterface {
					ctrl := gomock.NewController(t)
					mockRepo := repository.NewMockRepositoryInterface(ctrl)

					mockRepo.EXPECT().
						GetUserByPhoneNumber(gomock.Any(), user.PhoneNumber).
						Return(user, nil)

					return mockRepo
				}(),
				JWT: handler.NewJWT([]byte(RsaPublicKey), []byte(RsaPrivateKey)),
			},
			args: args{
				payload: map[string]interface{}{
					"phone_number": "+62812345678",
					"password":     "testpass!invalid",
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "failed when generating jwt",
			fields: fields{
				Repository: func() *repository.MockRepositoryInterface {
					ctrl := gomock.NewController(t)
					mockRepo := repository.NewMockRepositoryInterface(ctrl)

					mockRepo.EXPECT().
						GetUserByPhoneNumber(gomock.Any(), user.PhoneNumber).
						Return(user, nil)

					return mockRepo
				}(),
				JWT: handler.NewJWT([]byte("random_public_key"), []byte("random_private_key")),
			},
			args: args{
				payload: map[string]interface{}{
					"phone_number": "+62812345678",
					"password":     "testpass!",
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "failed when incrementing login counter",
			fields: fields{
				Repository: func() *repository.MockRepositoryInterface {
					ctrl := gomock.NewController(t)
					mockRepo := repository.NewMockRepositoryInterface(ctrl)

					mockRepo.EXPECT().
						GetUserByPhoneNumber(gomock.Any(), user.PhoneNumber).
						Return(user, nil)

					mockRepo.EXPECT().
						IncrementLoginCount(gomock.Any(), user.ID).
						Return(assert.AnError)

					return mockRepo
				}(),
				JWT: handler.NewJWT([]byte(RsaPublicKey), []byte(RsaPrivateKey)),
			},
			args: args{
				payload: map[string]interface{}{
					"phone_number": "+62812345678",
					"password":     "testpass!",
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			s := handler.NewServer(handler.NewServerOptions{
				Repository: tt.fields.Repository,
				JWT:        tt.fields.JWT,
			})

			generated.RegisterHandlers(e, s)

			response := testutil.NewRequest().Post("/auth").
				WithAcceptJson().
				WithJsonBody(tt.args.payload).
				GoWithHTTPHandler(t, e)

			assert.Equal(t, tt.wantStatus, response.Code())
		})
	}
}

func TestServer_GetLoggedInUser(t *testing.T) {
	user := repository.UserOutput{
		ID:          "c118a1a9-28f1-4137-9093-87487d24e5d9",
		FullName:    "test",
		PhoneNumber: "+62812345678",
		LoginCount:  1,
	}

	type fields struct {
		Repository repository.RepositoryInterface
		JWT        handler.JWT
	}
	type args struct {
		header map[string]string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
	}{
		{
			name: "successfully get the logged in user",
			fields: fields{
				Repository: func() *repository.MockRepositoryInterface {
					ctrl := gomock.NewController(t)
					mockRepo := repository.NewMockRepositoryInterface(ctrl)

					mockRepo.EXPECT().
						GetUserByID(gomock.Any(), user.ID).
						Return(user, nil)

					return mockRepo
				}(),
				JWT: handler.NewJWT([]byte(RsaPublicKey), []byte(RsaPrivateKey)),
			},
			args: args{
				header: map[string]string{
					"Authorization": fmt.Sprintf("Bearer %s", dummyJWT),
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "user is not logged in",
			fields: fields{
				Repository: func() *repository.MockRepositoryInterface {
					ctrl := gomock.NewController(t)
					mockRepo := repository.NewMockRepositoryInterface(ctrl)
					return mockRepo
				}(),
				JWT: handler.NewJWT([]byte(RsaPublicKey), []byte(RsaPrivateKey)),
			},
			args:       args{},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "error fetching user data",
			fields: fields{
				Repository: func() *repository.MockRepositoryInterface {
					ctrl := gomock.NewController(t)
					mockRepo := repository.NewMockRepositoryInterface(ctrl)

					mockRepo.EXPECT().
						GetUserByID(gomock.Any(), user.ID).
						Return(repository.UserOutput{}, assert.AnError)

					return mockRepo
				}(),
				JWT: handler.NewJWT([]byte(RsaPublicKey), []byte(RsaPrivateKey)),
			},
			args: args{
				header: map[string]string{
					"Authorization": fmt.Sprintf("Bearer %s", dummyJWT),
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			s := handler.NewServer(handler.NewServerOptions{
				Repository: tt.fields.Repository,
				JWT:        tt.fields.JWT,
			})

			generated.RegisterHandlers(e, s)

			req := testutil.NewRequest().Get("/me").
				WithAcceptJson()

			for key, value := range tt.args.header {
				req = req.WithHeader(key, value)
			}

			response := req.GoWithHTTPHandler(t, e)

			assert.Equal(t, tt.wantStatus, response.Code())
		})
	}
}

func TestServer_UpdateUser(t *testing.T) {
	user := repository.UserOutput{
		ID:          "c118a1a9-28f1-4137-9093-87487d24e5d9",
		FullName:    "test",
		PhoneNumber: "+62812345678",
		LoginCount:  1,
	}

	type fields struct {
		Repository repository.RepositoryInterface
		JWT        handler.JWT
	}
	type args struct {
		header  map[string]string
		payload map[string]interface{}
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
	}{
		{
			name: "successfully updates both full name & phone number",
			fields: fields{
				Repository: func() *repository.MockRepositoryInterface {
					ctrl := gomock.NewController(t)
					mockRepo := repository.NewMockRepositoryInterface(ctrl)

					mockRepo.EXPECT().
						GetUserByID(gomock.Any(), user.ID).
						Return(user, nil)

					mockRepo.EXPECT().
						GetUserByPhoneNumber(gomock.Any(), "+62812345677").
						Return(repository.UserOutput{}, sql.ErrNoRows)

					mockRepo.EXPECT().
						UpdateUser(gomock.Any(), user.ID, repository.UpdateUserInput{
							FullName:    "test new",
							PhoneNumber: "+62812345677",
						}).
						Return(nil)

					return mockRepo
				}(),
				JWT: handler.NewJWT([]byte(RsaPublicKey), []byte(RsaPrivateKey)),
			},
			args: args{
				header: map[string]string{
					"Authorization": fmt.Sprintf("Bearer %s", dummyJWT),
				},
				payload: map[string]interface{}{
					"full_name":    "test new",
					"phone_number": "+62812345677",
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "user is not logged in",
			fields: fields{
				Repository: func() *repository.MockRepositoryInterface {
					ctrl := gomock.NewController(t)
					mockRepo := repository.NewMockRepositoryInterface(ctrl)
					return mockRepo
				}(),
				JWT: handler.NewJWT([]byte(RsaPublicKey), []byte(RsaPrivateKey)),
			},
			args: args{
				payload: map[string]interface{}{
					"full_name":    "test new",
					"phone_number": "+62812345677",
				},
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "successfully updates only phone",
			fields: fields{
				Repository: func() *repository.MockRepositoryInterface {
					ctrl := gomock.NewController(t)
					mockRepo := repository.NewMockRepositoryInterface(ctrl)

					mockRepo.EXPECT().
						GetUserByID(gomock.Any(), user.ID).
						Return(user, nil)

					mockRepo.EXPECT().
						GetUserByPhoneNumber(gomock.Any(), "+62812345677").
						Return(repository.UserOutput{}, sql.ErrNoRows)

					mockRepo.EXPECT().
						UpdateUser(gomock.Any(), user.ID, repository.UpdateUserInput{
							FullName:    user.FullName,
							PhoneNumber: "+62812345677",
						}).
						Return(nil)

					return mockRepo
				}(),
				JWT: handler.NewJWT([]byte(RsaPublicKey), []byte(RsaPrivateKey)),
			},
			args: args{
				header: map[string]string{
					"Authorization": fmt.Sprintf("Bearer %s", dummyJWT),
				},
				payload: map[string]interface{}{
					"phone_number": "+62812345677",
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "successfully updates only full name",
			fields: fields{
				Repository: func() *repository.MockRepositoryInterface {
					ctrl := gomock.NewController(t)
					mockRepo := repository.NewMockRepositoryInterface(ctrl)

					mockRepo.EXPECT().
						GetUserByID(gomock.Any(), user.ID).
						Return(user, nil)

					mockRepo.EXPECT().
						UpdateUser(gomock.Any(), user.ID, repository.UpdateUserInput{
							FullName:    "test new",
							PhoneNumber: user.PhoneNumber,
						}).
						Return(nil)

					return mockRepo
				}(),
				JWT: handler.NewJWT([]byte(RsaPublicKey), []byte(RsaPrivateKey)),
			},
			args: args{
				header: map[string]string{
					"Authorization": fmt.Sprintf("Bearer %s", dummyJWT),
				},
				payload: map[string]interface{}{
					"full_name": "test new",
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "some field failed validation",
			fields: fields{
				Repository: func() *repository.MockRepositoryInterface {
					ctrl := gomock.NewController(t)
					mockRepo := repository.NewMockRepositoryInterface(ctrl)
					return mockRepo
				}(),
				JWT: handler.NewJWT([]byte(RsaPublicKey), []byte(RsaPrivateKey)),
			},
			args: args{
				header: map[string]string{
					"Authorization": fmt.Sprintf("Bearer %s", dummyJWT),
				},
				payload: map[string]interface{}{
					"phone_number": "12345677",
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "error updated phone number already exists",
			fields: fields{
				Repository: func() *repository.MockRepositoryInterface {
					ctrl := gomock.NewController(t)
					mockRepo := repository.NewMockRepositoryInterface(ctrl)

					mockRepo.EXPECT().
						GetUserByID(gomock.Any(), user.ID).
						Return(user, nil)

					mockRepo.EXPECT().
						GetUserByPhoneNumber(gomock.Any(), "+62812345677").
						Return(repository.UserOutput{ID: "exists"}, nil)

					return mockRepo
				}(),
				JWT: handler.NewJWT([]byte(RsaPublicKey), []byte(RsaPrivateKey)),
			},
			args: args{
				header: map[string]string{
					"Authorization": fmt.Sprintf("Bearer %s", dummyJWT),
				},
				payload: map[string]interface{}{
					"phone_number": "+62812345677",
				},
			},
			wantStatus: http.StatusConflict,
		},
		{
			name: "failed update user",
			fields: fields{
				Repository: func() *repository.MockRepositoryInterface {
					ctrl := gomock.NewController(t)
					mockRepo := repository.NewMockRepositoryInterface(ctrl)

					mockRepo.EXPECT().
						GetUserByID(gomock.Any(), user.ID).
						Return(user, nil)

					mockRepo.EXPECT().
						UpdateUser(gomock.Any(), user.ID, repository.UpdateUserInput{
							FullName:    "test new",
							PhoneNumber: user.PhoneNumber,
						}).
						Return(assert.AnError)

					return mockRepo
				}(),
				JWT: handler.NewJWT([]byte(RsaPublicKey), []byte(RsaPrivateKey)),
			},
			args: args{
				header: map[string]string{
					"Authorization": fmt.Sprintf("Bearer %s", dummyJWT),
				},
				payload: map[string]interface{}{
					"full_name": "test new",
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			s := handler.NewServer(handler.NewServerOptions{
				Repository: tt.fields.Repository,
				JWT:        tt.fields.JWT,
			})

			generated.RegisterHandlers(e, s)

			req := testutil.NewRequest().Patch("/users").
				WithAcceptJson().
				WithJsonBody(tt.args.payload)

			for key, value := range tt.args.header {
				req = req.WithHeader(key, value)
			}

			response := req.GoWithHTTPHandler(t, e)

			assert.Equal(t, tt.wantStatus, response.Code())
		})
	}
}
