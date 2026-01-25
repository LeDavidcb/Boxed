package internal

import (
	"github.com/David/Boxed/internal/auth/login"
	"github.com/David/Boxed/internal/auth/register"
	"github.com/David/Boxed/internal/health"
	"github.com/labstack/echo/v5"
)

func SetupControllers() *echo.Echo {

	echo := echo.New()

	echo.GET("/health", health.Hello)
	echo.GET("/auth/login", login.LoginController)
	echo.GET("/auth/register", register.RegisterController)
	return echo

}
