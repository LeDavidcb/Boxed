package controllers

import (
	"fmt"
	"mime"
	"net/http"
	"path"

	boxed "github.com/David/Boxed"
	"github.com/David/Boxed/internal/common/types"
	"github.com/David/Boxed/internal/files/services"
	"github.com/David/Boxed/repositories"
	"github.com/google/uuid"
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
		c.String(400, err.Error())
		return echo.ErrBadRequest.Wrap(err)
	}

	// Get the user
	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		c.String(400, err.Error())
		return echo.ErrBadRequest.Wrap(err)
	}
	user, err := ur.GetByID(id)
	if err != nil {
		return echo.ErrInternalServerError.Wrap(err)
	}

	// Process multiple files
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid form data.")
		return echo.ErrBadRequest.Wrap(err)
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.String(http.StatusBadRequest, "No files provided.")
		return echo.ErrBadRequest
	}
	// Iterate over files
	for _, file := range files {
		m := file.Header.Get("Content-Type")
		typet, _ := mime.ExtensionsByType(m)
		fileId := uuid.New()
		filename := fmt.Sprintf("%v%v", fileId.String(), typet[0])
		filePath := path.Join(user.FolderPath, filename)
		thumbnailPath := path.Join(user.FolderPath, fmt.Sprintf("/thumbnail/%v", filename))

		saveErr := services.SaveFile(filePath, file)
		if saveErr != nil {
			c.Logger().Error(saveErr.Error())
			continue
		}

		go services.CreateAndSaveThumbnail(thumbnailPath)

		dbErr := services.SaveFileToDatabase(db, file, fileId, id, filePath, thumbnailPath)
		if dbErr != nil {
			c.Logger().Error(dbErr.Error())
			continue
		}
	}

	return c.NoContent(http.StatusCreated)
}
