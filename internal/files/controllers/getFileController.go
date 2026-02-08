package controllers

import (
	"fmt"
	"net/http"

	boxed "github.com/David/Boxed"
	"github.com/David/Boxed/repositories"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

// GetFile retrieves metadata for a specific file identified by the UUID provided in the request header.
//
// Returns:
//   - Responds with HTTP 200 (OK) and the file metadata as JSON.
//   - Responds with HTTP 400 (Bad Request) if the UUID is invalid or the file could not be found.
//   - Returns an error if there are issues querying the database.
func GetFileController(c *echo.Context) error {
	id := c.Request().Header.Get("uuid")
	if id == "" {
		c.String(http.StatusBadRequest, "No uuid was provided.")
		return fmt.Errorf("No uuid file was provided.")
	}
	conn := boxed.GetInstance().DbConn
	fileRepo := repositories.NewFilesRepo(conn)
	ui, e := uuid.Parse(id)
	if e != nil {
		c.String(http.StatusBadRequest, "The uuid provided was not valid.")
		return fmt.Errorf("Bad uuid. %v", e)
	}
	f, e := fileRepo.GetByID(ui)
	if e != nil {
		c.String(http.StatusBadRequest, "Not Found")
		return e
	}
	return c.JSON(http.StatusOK, f)
}
