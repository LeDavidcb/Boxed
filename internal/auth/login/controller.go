package login

import (
	"net/http"

	boxed "github.com/David/Boxed"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
)

type userLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type loginResponse struct {
	SignedJwt    string `json:"signed-jwt"`
	RefreshToken string `json:"refresh-token"`
}

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
		user     *userLoginRequest
		response *loginResponse
	)

	if c.Request().Header.Get("Content-Type") != "application/json" {
		return c.NoContent(http.StatusUnsupportedMediaType)
	}

	err := echo.BindBody(c, &user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "")
	}
	if user.Email == "" || user.Password == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	response, err = validate(user, con)
	// TODO: check the type of error
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	return c.JSON(http.StatusOK, response)
}
