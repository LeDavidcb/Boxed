package health

import "github.com/labstack/echo/v5"

func Hello(c *echo.Context) error {
	return c.String(200, "Hello world")
}
