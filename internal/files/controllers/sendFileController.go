package controllers

import (
	"errors"
	"fmt"
	"mime"
	"net/http"
	"path"
	"strings"

	boxed "github.com/David/Boxed"
	"github.com/David/Boxed/internal/common/types"
	"github.com/David/Boxed/internal/files/services"
	"github.com/David/Boxed/repositories"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v5"
)

// SendFile uploads a single file for the authenticated user and saves both the file and its metadata to disk and the database.
//
// Returns:
//   - Responds with HTTP 201 (Created) on success.
//   - Responds with HTTP 400 (Bad Request) if the file or user info is invalid.
//   - Returns an error if saving the file or metadata fails.
func SendFileController(c *echo.Context) error {
	db := boxed.GetInstance().DbConn
	ur := repositories.NewUserRepo(db)

	file, err := c.FormFile("file")
	if err != nil {
		e := &types.ErrorResponse{
			Code:    types.MissingFields,
			Message: "Multipart with an entry `file` must be provided.",
		}
		return c.JSON(http.StatusBadRequest, &e)
	}
	claims, err := echo.ContextGet[*types.ResponseClaims](c, "user")
	if err != nil {
		c.String(400, err.Error())
		return echo.ErrBadRequest.Wrap(err)
	}
	// Get the user
	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		e := &types.ErrorResponse{
			Code:    types.InvalidFields,
			Message: "`uuid` provided is not valid.",
		}
		return c.JSON(http.StatusBadRequest, &e)
	}
	user, err := ur.GetByID(id)
	// ??????????????????????????????

	if err != nil {
		var pge *pgconn.PgError
		if errors.As(err, &pge) || errors.As(err, &pgx.ErrNoRows) {
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
	// metadata info
	m := file.Header.Get("Content-Type")
	typet, _ := mime.ExtensionsByType(m)
	fileId := uuid.New()
	filename := fmt.Sprintf("%v%v", fileId.String(), typet[0])
	originalName := strings.Split(file.Filename, ".")[0]
	filePath := path.Join(user.FolderPath, filename)
	// Create the file to the os
	err = services.SaveFile(filePath, file)
	if err != nil {
		e := &types.ErrorResponse{
			Code:    types.InternalServerError,
			Message: "Error while trying to save a file to the server. Please try later.",
		}
		return c.JSON(http.StatusInternalServerError, &e)
	}
	// Setup thumbnail
	thumbnailRepository := repositories.NewThumbnailRepository(db)
	thumbnailUUID := uuid.New()
	thumbnailRepository.Create(&repositories.Thumbnail{
		ID:      thumbnailUUID,
		OwnerId: user.ID,
	})
	thumbnailPath := path.Join(user.FolderPath, fmt.Sprintf("/thumbnail/%v.jpg", thumbnailUUID))
	// Generate Thumbnail
	go services.CreateAndSaveThumbnail(filePath, thumbnailPath, m, originalName, thumbnailUUID, thumbnailRepository)
	err = services.SaveFileToDatabase(db, file, fileId, id, filePath, thumbnailUUID)
	if err != nil {
		e := &types.ErrorResponse{
			Code:    types.InternalServerError,
			Message: "Error while trying to save a file to the database. Please try later.",
		}
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.NoContent(http.StatusCreated)
}
