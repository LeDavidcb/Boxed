package refresh

import (
	"context"
	"log"
	"net/http"

	boxed "github.com/David/Boxed"
	"github.com/David/Boxed/repositories"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v5"
)

// RefreshToken handles the process of refreshing a user's JWT by validating the provided refresh token.
// It generates a new access token (JWT) and a new refresh token, while ensuring transactional atomicity for database operations.
//
// Returns:
//   - Responds with HTTP 200 (OK) along with the new JWT and refresh token on success.
//   - Responds with appropriate HTTP status errors if token invalid, expired, or database operations fail.
//
// Errors:
//   - 400 Bad Request if refresh token is invalid or expired.
//   - 500 Internal Server Error for database or JWT generation failures.
func RefreshToken(c *echo.Context) error {
	// Get the refreshToken
	rt := c.Request().Header.Get("refresh-token")
	if rt == "" {
		return c.String(http.StatusBadRequest, "refresh-token header must be provided.")
	}
	// Set up transaction
	conn := boxed.GetInstance().DbConn
	t, err := conn.Begin(context.Background())
	if err != nil {
		log.Println("Error starting transaction at RefreshToken fn:", err)
		return c.String(http.StatusInternalServerError, "Internal server error, please try later.")
	}
	defer func() {
		if rollbackErr := t.Rollback(context.Background()); rollbackErr != nil && rollbackErr != pgx.ErrTxClosed {
			log.Printf("Rollback failed: %v", rollbackErr)
		}
	}()
	rtr := repositories.NewRefreshTokensRepo(conn)

	val, err := rtr.RegenerateToken(rt, &t)
	if err != nil {
		log.Println("Error while RegenerateToken:", err)
		return c.String(http.StatusBadRequest, "refresh-token is not valid or expired.")
	}
	// Get new JWT
	sig, err := ReSignJwt(val.Useruuid)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error generating, try later")
	}
	if commitErr := t.Commit(context.Background()); commitErr != nil {
		log.Println("Transaction commit error:", commitErr)
		return c.String(http.StatusInternalServerError, "Transaction failed")
	}
	return c.JSON(http.StatusOK, &struct {
		Jwt          string `json:"jwt"`
		RefreshToken string `json:"refresh-token"`
	}{sig, val.NewHash})
}
