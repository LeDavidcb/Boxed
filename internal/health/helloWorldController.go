package health

import (
	"github.com/David/Boxed/internal/common/types"
	"github.com/labstack/echo/v5"
)

func Hello(c *echo.Context) error {
	user, err := echo.ContextGet[*types.ResponseClaims](c, "user")
	if err != nil {
		c.String(500, err.Error())
		return echo.ErrUnauthorized.Wrap(err)
	}
	return c.JSON(200, user)
}
