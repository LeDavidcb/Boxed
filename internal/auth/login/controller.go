package login

import (
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v5"
)

type userLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginController(c *echo.Context) error {
	defer c.Request().Body.Close()
	content, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	fmt.Println(string(content))
	return c.NoContent(http.StatusCreated)
}
