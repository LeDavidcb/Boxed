package refresh

import (
	"log"
	"net/http"

	boxed "github.com/David/Boxed"
	"github.com/David/Boxed/repositories"
	"github.com/labstack/echo/v5"
)

func RefreshToken(c *echo.Context) error {
	// Get the refreshToken
	rt := c.Request().Header.Get("refresh-token")
	if rt == "" {
		return c.String(http.StatusBadRequest, "refresh-token header must be provided.")
	}
	// Create refresh token repository
	conn := boxed.GetInstance().DbConn
	rtr := repositories.NewRefreshTokensRepo(conn)
	val, err := rtr.RegenerateToken(rt)
	log.Println(val)
	if err != nil {
		return c.String(http.StatusBadRequest, "refresh-token is not valid or expired.")
	}
	// Get new JWT
	log.Println("INPUT VAL IN RT controller:", val)
	sig, err := ReSignJwt(val.Useruuid)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error generating, try later")
		// if this fails, please make that the old refreshToken gets in the database.
	}
	return c.JSON(http.StatusOK, &struct {
		Jwt          string `json:"jwt"`
		RefreshToken string `json:"refresh-token"`
	}{sig, val.NewHash})
}
