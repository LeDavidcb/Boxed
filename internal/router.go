package internal

import (
	"strings"

	boxed "github.com/David/Boxed"
	auth "github.com/David/Boxed/internal/auth/controllers"
	jwtMiddleware "github.com/David/Boxed/internal/auth/middleware"
	files "github.com/David/Boxed/internal/files/controllers"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func SetupControllers() *echo.Echo {

	key := strings.Trim(boxed.GetInstance().JwtSecret, " ")
	router := echo.New()

	router.Use(middleware.RequestLogger())
	router.GET("/auth/login", auth.LoginController)
	router.GET("/auth/register", auth.RegisterController)
	router.GET("/auth/refresh", auth.RefreshTokenController)

	jwtMiddleware := jwtMiddleware.NewJwtMiddleware(key, jwt.SigningMethodHS256)
	validated := router.Group("/api") // Temporarily commented out
	validated.Use(jwtMiddleware.Middleware)
	validated.POST("/upload-file", files.SendFileController)
	validated.POST("/upload-files", files.SendFilesController)
	validated.GET("/get-file", files.GetFileController)
	validated.GET("/get-files", files.GetFilesController)
	validated.GET("/serve-file", files.ServeFileController)
	validated.GET("/serve-thumbnail", files.ServeThumbnailController)
	validated.DELETE("/delete-file", files.DeleteFileController)
	return router

}
