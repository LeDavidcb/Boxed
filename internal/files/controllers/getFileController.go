package controllers

import (
	"errors"
	"fmt"
	"net/http"

	boxed "github.com/David/Boxed"
	"github.com/David/Boxed/internal/common/types"
	"github.com/David/Boxed/repositories"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
		e := &types.ErrorResponse{
			Code:    types.MissingFields,
			Message: "`uuid` must be provided.",
		}
		return c.JSON(http.StatusBadRequest, &e)
	}
	conn := boxed.GetInstance().DbConn
	fileRepo := repositories.NewFilesRepo(conn)
	ui, e := uuid.Parse(id)
	if e != nil {
		e := &types.ErrorResponse{
			Code:    types.InvalidFields,
			Message: "`uuid` provided is not valid.",
		}
		return c.JSON(http.StatusBadRequest, &e)
	}
	f, e := fileRepo.GetByID(ui)
	if e != nil {
		var pge *pgconn.PgError
		if errors.As(e, &pge) || errors.As(e, &pgx.ErrNoRows) {
			em := &types.ErrorResponse{
				Code:    types.ResourceNotFound,
				Message: fmt.Sprintf("Couldn't get any file with %v to serve.", id),
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
	return c.JSON(http.StatusOK, f)
}
