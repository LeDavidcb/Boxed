package controllers

import (
	"fmt"
	"net/http"

	boxed "github.com/David/Boxed"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
)

type userRegisterRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"Password"`
}

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
	var user userRegisterRequest
	if c.Request().Header.Get("Content-Type") != "application/json" {
		return c.NoContent(http.StatusUnsupportedMediaType)
	}
	err := echo.BindBody(c, &user)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Error at the provided body: %v", err))
		return echo.NewHTTPError(http.StatusBadRequest, "")
	}
	if user.Nickname == "" || user.Email == "" || user.Password == "" {
		c.String(http.StatusBadRequest, fmt.Sprintf("Error at the provided body: %v", err))
		return c.NoContent(http.StatusBadRequest)
	}
	if err := createUserDb(con, &user); err != nil {
		c.String(http.StatusInternalServerError, "Couldn't create the user.")
		return err
	}
	return c.NoContent(http.StatusCreated)
}
