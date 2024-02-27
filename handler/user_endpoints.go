package handler

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// Register a new user to the service, using the specified full name, phone number, and password
// (POST /users)
func (s *Server) RegisterUser(c echo.Context) error {
	var payload RegisterUserValidator
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	fieldErrors := payload.Validate()
	if len(fieldErrors) > 0 {
		return c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message:          "field validation failed",
			ValidationErrors: (*[]generated.FieldError)(&fieldErrors),
		})
	}

	ctx := c.Request().Context()

	// check if phone number already exists in DB
	existingUser, err := s.Repository.GetUserByPhoneNumber(ctx, payload.PhoneNumber)
	if err != nil && err != sql.ErrNoRows {
		return c.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}
	if existingUser.ID != "" {
		return c.JSON(http.StatusConflict, generated.ErrorResponse{
			Message: "phone number already registered",
		})
	}

	// generate password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	id := uuid.NewString()
	userInput := repository.CreateUserInput{
		ID:             id,
		FullName:       payload.FullName,
		PhoneNumber:    payload.PhoneNumber,
		HashedPassword: string(hashedPassword),
	}

	// save
	output, err := s.Repository.CreateUser(ctx, userInput)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, generated.UserResponse{
		Id:          output.ID,
		FullName:    payload.FullName,
		PhoneNumber: payload.PhoneNumber,
	})
}

// Logs a user in to the system & return the logged in user ID & generated jwt token
// (POST /auth)
func (s *Server) AuthenticateUser(c echo.Context) error {
	var payload AuthenticateUserValidator
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	ctx := c.Request().Context()

	// check for user
	existingUser, err := s.Repository.GetUserByPhoneNumber(ctx, payload.PhoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusBadRequest, generated.ErrorResponse{
				Message: "user not valid",
			})
		}

		return c.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	// check password correctness
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.HashedPassword), []byte(payload.Password))
	if err != nil {
		return c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "user not valid",
		})
	}

	// password correct: return
	token, err := s.GenerateJWT(existingUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	err = s.Repository.IncrementLoginCount(ctx, existingUser.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, generated.AuthenticateUserResponse{
		Id:    existingUser.ID,
		Token: token,
	})
}

func (s *Server) GenerateJWT(user repository.UserOutput) (string, error) {
	claims := JWTCustomClaims{
		ID:          user.ID,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}

	return s.JWT.CreateToken(claims)
}

// Get the logged in user info
// (GET /me)
func (s *Server) GetLoggedInUser(c echo.Context) error {
	id, err := s.ValidateLoggedInUser(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, generated.ErrorResponse{
			Message: "not logged in",
		})
	}

	user, err := s.Repository.GetUserByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, generated.UserResponse{
		Id:          user.ID,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		LoginCount:  int64(user.LoginCount),
	})
}

func (s *Server) ValidateLoggedInUser(c echo.Context) (string, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("unauthorized")
	}

	tokens := strings.Split(authHeader, " ")
	if len(tokens) != 2 {
		return "", errors.New("unauthorized")
	}

	return s.JWT.ValidateToken(tokens[1])
}

// Updates the logged in user info. Allows for phone number & full name update
// (PATCH /users)
func (s *Server) UpdateUser(c echo.Context) error {
	id, err := s.ValidateLoggedInUser(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, generated.ErrorResponse{
			Message: "not logged in",
		})
	}

	var payload UpdateUserValidator
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	fieldErrors := payload.Validate()
	if len(fieldErrors) > 0 {
		return c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: fieldErrors.Error(),
		})
	}

	ctx := c.Request().Context()

	// get existing user data to retrieve the initial full name & phone number
	existingUser, err := s.Repository.GetUserByID(ctx, id)
	if err != nil && err != sql.ErrNoRows {
		return c.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	updateInput := repository.UpdateUserInput{
		FullName:    existingUser.FullName,
		PhoneNumber: existingUser.PhoneNumber,
	}

	// first, populate phone number if it is filled & does not exist
	if payload.PhoneNumber != nil {
		userWithSamePhone, err := s.Repository.GetUserByPhoneNumber(ctx, *payload.PhoneNumber)
		if err != nil && err != sql.ErrNoRows {
			return c.JSON(http.StatusInternalServerError, generated.ErrorResponse{
				Message: err.Error(),
			})
		}

		if userWithSamePhone.ID != "" && userWithSamePhone.ID != existingUser.ID {
			return c.JSON(http.StatusConflict, generated.ErrorResponse{
				Message: "phone number already registered",
			})
		}

		updateInput.PhoneNumber = *payload.PhoneNumber
	}

	// then populate full name only if it exists
	if payload.FullName != nil {
		updateInput.FullName = *payload.FullName
	}

	err = s.Repository.UpdateUser(ctx, id, updateInput)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, generated.UserResponse{
		Id:          existingUser.ID,
		FullName:    updateInput.FullName,
		PhoneNumber: updateInput.PhoneNumber,
		LoginCount:  int64(existingUser.LoginCount),
	})
}
