package internal

import (
	"strings"

	boxed "github.com/David/Boxed"
	"github.com/David/Boxed/internal/auth/login"
	"github.com/David/Boxed/internal/auth/register"
	"github.com/David/Boxed/internal/common/types"
	"github.com/David/Boxed/internal/files"
	"github.com/David/Boxed/internal/health"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func SetupControllers() *echo.Echo {

	key := strings.Trim(boxed.GetInstance().JwtSecret, " ")
	router := echo.New()

	router.Use(middleware.RequestLogger())
	router.GET("/auth/login", login.LoginController)
	router.GET("/auth/register", register.RegisterController)

	jwtMiddleware := types.NewJwtMiddleware(key, jwt.SigningMethodHS256)
	validated := router.Group("/api") // Temporarily commented out
	validated.Use(jwtMiddleware.Middleware)
	validated.GET("/h", health.Hello)
	validated.POST("/upload-file", files.SendFile)
	validated.POST("/upload-files", files.SendFiles)
	return router

}
