package controllers

import (
	"net/http"

	boxed "github.com/David/Boxed"
	"github.com/David/Boxed/internal/common/types"
	"github.com/David/Boxed/repositories"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

// GetFiles retrieves all files (Metadata) owned by the authenticated user and returns their metadata.
//
// Returns:
//   - Responds with HTTP 200 (OK) along with a JSON payload containing file metadata.
//   - Responds with HTTP 404 (Not Found) if no files exist for the user.
func GetFilesController(c *echo.Context) error {
	user, err := echo.ContextGet[*types.ResponseClaims](c, "user")
	if err != nil {
		c.NoContent(http.StatusInternalServerError)
		return echo.ErrUnauthorized.Wrap(err)
	}
	uid, err := uuid.Parse(user.Subject)
	if err != nil {
		c.String(http.StatusNotFound, "No User with that uuid.")
		return echo.ErrUnauthorized.Wrap(err)
	}
	frepo := repositories.NewFilesRepo(boxed.GetInstance().DbConn)
	files, err := frepo.GetByOwnerID(uid)
	if err != nil {
		c.NoContent(http.StatusNotFound)
		return echo.ErrUnauthorized.Wrap(err)
	}
	content := struct {
		Length int `json:"length"`
		Files  any `json:"files"`
	}{
		Length: len(files),
		Files:  files,
	}
	return c.JSON(200, content)
}
