package middleware

import (
	"net/http"
	"strings"

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
			return c.String(http.StatusUnauthorized, "No authorization token was provided.")
		}
		// Parse the unparsedToken
		var unparsedToken string

		if t := strings.Split(rv, " "); len(t) < 2 {
			return c.String(http.StatusBadRequest, "A non valid token was provided.")
		} else {
			unparsedToken = t[1]
		}
		// Validate the token
		token, err := jwt.ParseWithClaims(unparsedToken, &ResponseClaims{}, func(t *jwt.Token) (any, error) {
			return []byte(self.Key), nil
		})
		if err != nil {
			return c.String(http.StatusUnauthorized, err.Error())
		}
		if !token.Valid {
			return c.String(http.StatusUnauthorized, "The token provided in the authorization header is not valid.")
		}
		// add the token to the echo.context storage with a Key of: "user"
		c.Set("user", token.Claims)
		return n(c)
	}
}
