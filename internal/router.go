package internal

import (
	"github.com/David/Boxed/internal/health"
	"github.com/labstack/echo/v5"
)

func SetupControllers() *echo.Echo {

	echo := echo.New()

	echo.GET("/health", health.Hello)
	return echo

}
