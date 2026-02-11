package controllers

import (
	"fmt"
	"log"
	"net/http"

	boxed "github.com/David/Boxed"
	"github.com/David/Boxed/internal/files/services"
	"github.com/David/Boxed/repositories"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

// DeleteFile removes a file identified by its UUID, along with its metadata, from both the database and storage.
//
// Returns:
//   - Responds with HTTP 200 (OK) for successful deletion.
//   - Responds with HTTP 400 (Bad Request) if the UUID is invalid or the file does not exist.
func DeleteFileController(c *echo.Context) error {
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
	// get by id
	f, e := fileRepo.GetByID(ui)
	if e != nil {
		c.String(http.StatusBadRequest, "Couldn't get any file with that id.")
		return fmt.Errorf("No file with the uuid: %v", e)
	}
	if e := fileRepo.Delete(ui); e != nil {
		c.String(http.StatusBadRequest, "Error while deleting the file, please try later.")
		return fmt.Errorf("Error while deleting file: %v", e)
	}
	// Get Thumbnail path
	tr := repositories.NewThumbnailRepository(boxed.GetInstance().DbConn)
	t, err := tr.GetByID(f.ThumbnailId)
	if err != nil {
		log.Println("No thumnbail by this id:", f.ThumbnailId)
	} else {
		tr.DeleteByID(t.ID)
		go services.DeleteFile(t.StoragePath)
	}
	go services.DeleteFile(f.StoragePath)
	return c.NoContent(http.StatusOK)
}
