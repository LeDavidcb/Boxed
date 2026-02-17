package controllers

import (
	"errors"
	"net/http"

	boxed "github.com/David/Boxed"
	"github.com/David/Boxed/internal/common/types"
	"github.com/David/Boxed/repositories"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
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
		if errors.As(err, &echo.ErrNonExistentKey) {
			e := &types.ErrorResponse{
				Code:    types.InternalServerError, // Couldn't get jwt, so it's a middleware error.
				Message: "Error while getting user from jwt, please try again.",
			}
			return c.JSON(http.StatusInternalServerError, &e)
		}
	}
	uid, err := uuid.Parse(user.Subject)
	if err != nil {
		e := &types.ErrorResponse{
			Code:    types.InvalidFields,
			Message: "`uuid` provided is not valid.",
		}
		return c.JSON(http.StatusBadRequest, &e)
	}
	frepo := repositories.NewFilesRepo(boxed.GetInstance().DbConn)
	files, err := frepo.GetByOwnerID(uid)
	if err != nil {
		var pge *pgconn.PgError
		if errors.As(err, &pge) {
			e := &types.ErrorResponse{
				Code:    types.ResourceNotFound,
				Message: "Couldn't get any user with this id : " + uid.String(),
			}
			return c.JSON(http.StatusNotFound, &e)
		} else {

			e := &types.ErrorResponse{
				Code:    types.InternalServerError,
				Message: "Internal server error while getting files, Please try later.",
			}
			return c.JSON(http.StatusInternalServerError, &e)
		}
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
