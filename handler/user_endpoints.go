package handler

import (
	"fmt"
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (s *Server) RegisterUser(ctx echo.Context) error {
	var payload generated.RegisterUserRequest
	if err := ctx.Bind(&payload); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	validateRegisterUserRequest(payload)

	return ctx.JSON(http.StatusCreated, generated.UserResponse{
		Id:          "123-456",
		FullName:    "jokowi",
		PhoneNumber: "+6281234567890",
	})
}

type RegisterUserValidationPayload struct {
	FullName    string `validate:"required,min=3,max=60"`
	PhoneNumber string `validate:"required,min=10,max=13,startswith=+62"`
	Password    string `validate:"required,min=6,max=64"`
}

func validateRegisterUserRequest(req generated.RegisterUserRequest) {
	var payload RegisterUserValidationPayload
	payload.FullName = req.FullName
	payload.PhoneNumber = req.PhoneNumber
	payload.Password = req.Password

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(payload)

	validationErrors := err.(validator.ValidationErrors)
	fmt.Println(validationErrors)
}

// Logs a user in to the system & return the logged in user ID & generated jwt token
// (POST /auth)
func (s *Server) AuthenticateUser(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, generated.AuthenticateUserResponse{
		Id:    uuid.NewString(),
		Token: "jwt-token",
	})
}

// Get the logged in user info
// (GET /me)
func (s *Server) GetLoggedInUser(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, generated.UserResponse{})
}

// Updates the logged in user info. Allows for phone number & full name update
// (PATCH /users)
func (s *Server) UpdateUser(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, generated.UserResponse{})
}
