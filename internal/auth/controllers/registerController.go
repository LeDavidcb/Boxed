package controllers

import (
	"errors"
	"net/http"

	boxed "github.com/David/Boxed"
	registerservices "github.com/David/Boxed/internal/auth/services/registerServices"
	"github.com/David/Boxed/internal/auth/types"
	commonTypes "github.com/David/Boxed/internal/common/types"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
)

// RegisterController handles user registration by validating incoming data and creating a new user in the database.
//
// Returns:
//   - Responds with HTTP 201 (Created) upon successful user registration.
//   - Responds with appropriate HTTP error codes for validation or database failures.
//
// Errors:
//   - 400 Bad Request for missing or invalid fields.
//   - 415 Unsupported Media Type for incorrect Content-Type header.
//   - 500 Internal Server Error for database-related errors.
func RegisterController(c *echo.Context) error {
	defer c.Request().Body.Close()
	var con *pgxpool.Pool = boxed.GetInstance().DbConn
	var user types.UserRegisterRequest
	if c.Request().Header.Get("Content-Type") != "application/json" {
		e := &commonTypes.ErrorResponse{
			Code:    commonTypes.InvalidFormat,
			Message: "The Content-Type request must be `application/json`",
		}
		return c.JSON(http.StatusUnsupportedMediaType, &e)
	}
	err := echo.BindBody(c, &user)

	if err != nil {
		e := &commonTypes.ErrorResponse{
			Code:    commonTypes.InvalidFormat,
			Message: "No body provided. Please provide a valid body for register process",
		}
		return c.JSON(http.StatusBadRequest, &e)
	}
	if user.Nickname == "" || user.Email == "" || user.Password == "" {
		e := &commonTypes.ErrorResponse{
			Code:    commonTypes.MissingFields,
			Message: "All properties (Nickname, Email and Password) MUST be provided.",
		}
		return c.JSON(http.StatusBadRequest, &e)
	}
	if err := registerservices.CreateUserInDatabase(con, &user); err != nil {
		var pge *pgconn.PgError
		if errors.As(err, &pge) {
			e := &commonTypes.ErrorResponse{
				Code:    commonTypes.UserEmailAlreadyExists,
				Message: "This email already exists in the database.",
			}
			return c.JSON(http.StatusBadRequest, &e)
		}

		e := &commonTypes.ErrorResponse{
			Code:    commonTypes.InternalServerError,
			Message: "Error while processing register, Please try again later.",
		}
		return c.JSON(http.StatusInternalServerError, &e)

	}
	return c.NoContent(http.StatusCreated)
}
