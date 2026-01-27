package files

import (
	"fmt"
	"mime"
	"net/http"
	"path"

	boxed "github.com/David/Boxed"
	"github.com/David/Boxed/internal/common/types"
	"github.com/David/Boxed/repositories"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

func SendFile(c *echo.Context) error {
	db := boxed.GetInstance().DbConn
	ur := repositories.NewUserRepo(db)
	//fr := repositories.NewFilesRepo(db)
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "File must be provided.")
	}
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

	// metadata info
	m := file.Header.Get("Content-Type")
	typet, _ := mime.ExtensionsByType(m)
	fileId := uuid.New()
	filename := fmt.Sprintf("%v%v", fileId.String(), typet[0])
	filePath := path.Join(user.FolderPath, filename)
	thumbnailPath := path.Join(user.FolderPath, fmt.Sprintf("/thumbnail/%v", filename))
	// Create the file to the os
	err = saveFile(filePath, file)
	if err != nil {
		return err
	}
	// TODO BEHAVIOUR
	go createAndSaveThumbnail(thumbnailPath)
	err = saveFileToDb(db, file, fileId, id, filePath, thumbnailPath)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}
