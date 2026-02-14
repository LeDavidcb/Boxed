package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	boxed "github.com/David/Boxed"
	"github.com/David/Boxed/internal/common/types"
	"github.com/David/Boxed/internal/files/services"
	"github.com/David/Boxed/repositories"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
		e := &types.ErrorResponse{
			Code:    types.MissingFields,
			Message: "`uuid` must be provided.",
		}
		c.JSON(http.StatusBadRequest, &e)
	}

	conn := boxed.GetInstance().DbConn
	fileRepo := repositories.NewFilesRepo(conn)

	ui, e := uuid.Parse(id)
	if e != nil {
		e := &types.ErrorResponse{
			Code:    types.InvalidFields,
			Message: "`uuid` provided is not valid.",
		}
		c.JSON(http.StatusBadRequest, &e)
	}
	// get by id
	f, e := fileRepo.GetByID(ui)
	if e != nil {
		var pge *pgconn.PgError
		if errors.As(e, &pge) || errors.As(e, &pgx.ErrNoRows) {
			em := &types.ErrorResponse{
				Code:    types.ResourceNotFound,
				Message: fmt.Sprintf("Couldn't get any file with %v to delete.", id),
			}
			return c.JSON(http.StatusBadRequest, &em)
		} else {
			em := &types.ErrorResponse{
				Code:    types.ResourceNotFound,
				Message: fmt.Sprintf("Internal error while getting file with id: %v", id),
			}
			return c.JSON(http.StatusInternalServerError, &em)
		}
	}
	if e := fileRepo.Delete(ui); e != nil {

		em := &types.ErrorResponse{
			Code:    types.ResourceDeleteFailed,
			Message: fmt.Sprintf("Internal error while deleting file `%v`. Please try later.", id),
		}
		return c.JSON(http.StatusInternalServerError, &em)
	}
	// Get Thumbnail path
	tr := repositories.NewThumbnailRepository(boxed.GetInstance().DbConn)
	t, err := tr.GetByID(f.ThumbnailId)
	if err != nil {
		log.Println("No thumnbail by this id:", f.ThumbnailId)
	} else {
		err := tr.DeleteByID(t.ID)
		if err != nil {
			var pge *pgconn.PgError
			if errors.As(e, &pge) || errors.As(e, &pgx.ErrNoRows) {
				em := &types.ErrorResponse{
					Code:    types.ResourceNotFound,
					Message: fmt.Sprintf("Couldn't get any thumnbail with %v to delete.", id),
				}
				return c.JSON(http.StatusBadRequest, &em)
			} else {
				em := &types.ErrorResponse{
					Code:    types.ResourceDeleteFailed,
					Message: fmt.Sprintf("Internal error while deleting thumnbail with id: %v", id),
				}
				return c.JSON(http.StatusInternalServerError, &em)
			}
		}
		go services.DeleteFile(t.StoragePath)
	}
	go services.DeleteFile(f.StoragePath)
	return c.NoContent(http.StatusOK)
}
