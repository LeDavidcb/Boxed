package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/David/Boxed/internal/common/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

type JwtMiddleware struct {
	Key           string
	SigningMethod jwt.SigningMethod
}

// Return a JwtMiddleware struct based on the struct's parameters
func NewJwtMiddleware(k string, sm jwt.SigningMethod) *JwtMiddleware {
	r := new(JwtMiddleware)
	r.Key = k
	r.SigningMethod = sm
	return r
}
func (self *JwtMiddleware) Middleware(n echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		// Get the header
		rv := c.Request().Header.Get("Authorization")
		if rv == "" {
			e := &types.ErrorResponse{
				Code:    types.MissingFields,
				Message: "Authorization token must be provided.",
			}
			return c.JSON(http.StatusBadRequest, &e)
		}
		// Parse the unparsedToken
		var unparsedToken string

		if t := strings.Split(rv, " "); len(t) < 2 {
			e := &types.ErrorResponse{
				Code:    types.InvalidFormat,
				Message: "A non-valid form of authorization was provided. Provide a header entry like so: `Authorization: Bearer <TOKEN>`",
			}
			return c.JSON(http.StatusBadRequest, &e)
		} else {
			unparsedToken = t[1]
		}
		// Validate the token
		token, err := jwt.ParseWithClaims(unparsedToken, &types.ResponseClaims{}, func(t *jwt.Token) (any, error) {
			return []byte(self.Key), nil
		})
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				e := &types.ErrorResponse{
					Code:    types.AuthTokenExpired,
					Message: "Token provided is expired, please refresh it or log in.",
				}
				return c.JSON(http.StatusUnauthorized, &e)
			}
			if errors.Is(err, jwt.ErrSignatureInvalid) || errors.Is(err, jwt.ErrTokenMalformed) {
				e := &types.ErrorResponse{
					Code:    types.AuthTokenInvalid,
					Message: "The provided token is invalid. It may be malformed or unsigned by the server. Please log in again.",
				}
				return c.JSON(http.StatusUnauthorized, &e)
			}
			// Default
			e := &types.ErrorResponse{
				Code:    types.AuthTokenInvalid,
				Message: err.Error(),
			}
			return c.JSON(http.StatusUnauthorized, &e)
		}
		if !token.Valid {
			e := &types.ErrorResponse{
				Code:    types.AuthTokenInvalid,
				Message: "The token provided in the authorization header is not valid.",
			}
			return c.JSON(http.StatusUnauthorized, &e)
		}
		// add the token to the echo.context storage with a Key of: "user"
		c.Set("user", token.Claims)
		return n(c)
	}
}
