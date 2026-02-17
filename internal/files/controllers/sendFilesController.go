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

// SendFiles allows authenticated users to upload multiple files at once, saving file data and metadata to database and disk.
// Returns:
//   - Responds with HTTP 201 (Created) after successfully processing all files.
//   - Responds with HTTP 400 (Bad Request) if form data or file inputs are invalid.
func SendFilesController(c *echo.Context) error {
	db := boxed.GetInstance().DbConn
	ur := repositories.NewUserRepo(db)

	claims, err := echo.ContextGet[*types.ResponseClaims](c, "user")
	if err != nil {
		e := &types.ErrorResponse{
			Code:    types.InternalServerError, // Couldn't get jwt, so it's a middleware error.
			Message: "Error while getting user from jwt, please try again.",
		}
		return c.JSON(http.StatusInternalServerError, &e)
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

	// Process multiple files
	form, err := c.MultipartForm()
	if err != nil {
		e := &types.ErrorResponse{
			Code:    types.InvalidFormat,
			Message: "body must be a MultipartForm",
		}
		return c.JSON(http.StatusBadRequest, &e)
	}

	files := form.File["files"]
	if len(files) == 0 {
		e := &types.ErrorResponse{
			Code:    types.MissingFields,
			Message: "Multipart with one or many `files` entries must be provided.",
		}
		return c.JSON(http.StatusBadRequest, &e)
	}
	// Track of files that failed to update.
	var failed []string
	// Iterate over files
	for _, file := range files {
		m := file.Header.Get("Content-Type")
		typet, _ := mime.ExtensionsByType(m)
		fileId := uuid.New()
		filename := fmt.Sprintf("%v%v", fileId.String(), typet[0])
		filePath := path.Join(user.FolderPath, filename)
		originalName := strings.Split(file.Filename, ".")[0]
		thumbnailPath := path.Join(user.FolderPath, fmt.Sprintf("/thumbnail/%v", filename))

		saveErr := services.SaveFile(filePath, file)
		if saveErr != nil {
			c.Logger().Error(saveErr.Error())
			continue
		}

		// Setup thumbnail
		thumbnailRepository := repositories.NewThumbnailRepository(db)
		thumbnailUUID := uuid.New()
		thumbnailRepository.Create(&repositories.Thumbnail{
			ID:      thumbnailUUID,
			OwnerId: user.ID,
		})
		// Generate Thumbnail
		go services.CreateAndSaveThumbnail(filePath, thumbnailPath, m, originalName, thumbnailUUID, thumbnailRepository)
		err = services.SaveFileToDatabase(db, file, fileId, id, filePath, thumbnailUUID)
		if err != nil {
			failed = append(failed, file.Filename)
			continue
		}
	}
	if len(failed) > 0 {
		e := &types.ErrorResponse{
			Code:    types.FileUploadFailed,
			Message: fmt.Sprintf("Some files failed to upload: %v \n", failed),
		}
		return c.JSON(http.StatusInternalServerError, &e)
	}
	return c.NoContent(http.StatusCreated)
}
