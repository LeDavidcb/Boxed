package controllers

import (
	"errors"
	"net/http"

	boxed "github.com/David/Boxed"
	"github.com/David/Boxed/internal/auth/services"
	"github.com/David/Boxed/internal/auth/types"
	commonTypes "github.com/David/Boxed/internal/common/types"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
)

// LoginController manages user login by validating credentials and returning a signed JWT along with a refresh token.
//
// Returns:
//   - Responds with HTTP 200 (OK) and `SignedJwt` and `RefreshToken` on successful authentication.
//   - Responds with appropriate HTTP error codes if credentials are invalid or required fields are missing.
//
// Errors:
//   - 400 Bad Request for invalid fields or missing data.
//   - 415 Unsupported Media Type for missing or incorrect Content-Type header.
func LoginController(c *echo.Context) error {
	defer c.Request().Body.Close()
	var con *pgxpool.Pool = boxed.GetInstance().DbConn
	var (
		user     *types.UserLoginRequest
		response *types.LoginResponse
	)

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
			Code:    commonTypes.MissingFields,
			Message: "No body provided. Please provide a valid body for login process",
		}
		return c.JSON(http.StatusBadRequest, &e)
	}

	if user.Email == "" || user.Password == "" {
		e := &commonTypes.ErrorResponse{
			Code:    commonTypes.InvalidFormat,
			Message: "All properties (Email and Password) MUST be provided.",
		}
		return c.JSON(http.StatusBadRequest, &e)
	}

	response, err = services.Validate(user, con)
	// TODO: check the type of error
	if err != nil {
		var pge *pgconn.PgError
		if errors.As(err, &pge) || errors.As(err, &pgx.ErrNoRows) {
			e := &commonTypes.ErrorResponse{
				Code:    commonTypes.AuthInvalidCredentials,
				Message: "Invalid credentials provided.",
			}
			return c.JSON(http.StatusBadRequest, &e)
		} else {
			e := &commonTypes.ErrorResponse{
				Code:    commonTypes.InternalServerError,
				Message: "Internal error when performing login, please try alter.",
			}
			return c.JSON(http.StatusInternalServerError, &e)
		}
		return c.NoContent(http.StatusBadRequest)
	}
	return c.JSON(http.StatusOK, response)
}
