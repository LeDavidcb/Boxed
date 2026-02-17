package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"

	boxed "github.com/David/Boxed"
	"github.com/David/Boxed/internal/auth/services"
	"github.com/David/Boxed/internal/common/types"
	"github.com/David/Boxed/repositories"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
func RefreshTokenController(c *echo.Context) error {
	// Get the refreshToken
	rt := c.Request().Header.Get("refresh-token")
	if rt == "" {
		e := &types.ErrorResponse{
			Code:    types.RefreshTokenMissing,
			Message: "refresh-token must be priovided to perform this.",
		}
		return c.JSON(http.StatusBadRequest, &e)
	}
	// Set up transaction
	conn := boxed.GetInstance().DbConn
	t, err := conn.Begin(context.Background())
	if err != nil {
		e := &types.ErrorResponse{
			Code:    types.InternalServerError,
			Message: "Internal server error while performing refresh, Please try later.",
		}
		log.Printf("Error while starting transaction, More info: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, &e)
	}

	defer func() {
		if rollbackErr := t.Rollback(context.Background()); rollbackErr != nil && rollbackErr != pgx.ErrTxClosed {
			log.Printf("Rollback failed: %v", rollbackErr)
		}
	}()

	refrshTokenRepository := repositories.NewRefreshTokensRepo(conn)

	val, err := refrshTokenRepository.RegenerateToken(rt, &t)
	if err != nil {
		var pge *pgconn.PgError
		if errors.As(err, &pge) || errors.As(err, &pgx.ErrNoRows) {
			e := &types.ErrorResponse{
				Code:    types.RefreshTokenExpiredOrInvalid,
				Message: "The refresh token provided was either invalid or had expired.",
			}
			return c.JSON(http.StatusBadRequest, &e)
		} else {
			log.Println("Error while RegenerateToken:", err)
			e := &types.ErrorResponse{
				Code:    types.InternalServerError,
				Message: "Internal error while regenerating token, please try later.",
			}
			return c.JSON(http.StatusInternalServerError, &e)
		}
	}
	// Get new JWT
	sig, err := services.ReSignJwt(val.Useruuid)
	if err != nil {
		e := &types.ErrorResponse{
			Code:    types.InternalServerError,
			Message: "Internal error while signing jwt, please try later.",
		}
		return c.JSON(http.StatusInternalServerError, &e)
	}
	if commitErr := t.Commit(context.Background()); commitErr != nil {
		log.Println("Transaction commit error:", commitErr)
		e := &types.ErrorResponse{
			Code:    types.InternalServerError,
			Message: "Internal error while performing refresh operations, please try later.",
		}
		return c.JSON(http.StatusInternalServerError, &e)
	}
	return c.JSON(http.StatusOK, &struct {
		Jwt          string `json:"jwt"`
		RefreshToken string `json:"refresh-token"`
	}{sig, val.NewHash})
}
