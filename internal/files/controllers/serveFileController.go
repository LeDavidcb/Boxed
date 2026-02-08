package controllers

import (
	"fmt"
	"net/http"

	boxed "github.com/David/Boxed"
	"github.com/David/Boxed/internal/common/types"
	"github.com/David/Boxed/repositories"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

// ServeFile streams a requested file to the user based on its UUID.
// Authorization is validated to ensure the requesting user owns the file.
//
// Returns:
//   - Responds with the file content if successful.
//   - Responds with HTTP 403 (Forbidden) if the user is not authorized.
//   - Responds with HTTP 400 (Bad Request) or HTTP 401 (Unauthorized) based on validation errors.
func ServeFileController(c *echo.Context) error {
	// Extract file UUID from path parameter
	id := c.Request().Header.Get("uuid")
	if id == "" {
		c.String(http.StatusBadRequest, "No uuid was provided.")
		return fmt.Errorf("No uuid file was provided.")
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		c.String(http.StatusBadRequest, "The uuid provided was not valid.")
		return fmt.Errorf("Bad uuid. %v", err)
	}
	// Verify user authorization
	userClaims, err := echo.ContextGet[*types.ResponseClaims](c, "user")
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	userID, _ := uuid.Parse(userClaims.Subject)
	// Query file details
	fileRepo := repositories.NewFilesRepo(boxed.GetInstance().DbConn)
	file, err := fileRepo.GetByID(uid)
	if file.OwnerID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
	}
	if err != nil {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "file not found"})
	}
	// Serve the file
	return c.File(file.StoragePath)
}
