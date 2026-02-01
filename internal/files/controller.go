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
		return echo.ErrBadRequest
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

func SendFiles(c *echo.Context) error {
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

		saveErr := saveFile(filePath, file)
		if saveErr != nil {
			c.Logger().Error(saveErr.Error())
			continue
		}

		go createAndSaveThumbnail(thumbnailPath)

		dbErr := saveFileToDb(db, file, fileId, id, filePath, thumbnailPath)
		if dbErr != nil {
			c.Logger().Error(dbErr.Error())
			continue
		}
	}

	return c.NoContent(http.StatusCreated)
}
func GetFile(c *echo.Context) error {
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
func DeleteFile(c *echo.Context) error {
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
	go deleteFile(f.StoragePath)
	return c.NoContent(http.StatusOK)
}
func GetFiles(c *echo.Context) error {
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
func ServeFile(c *echo.Context) error {
	// Extract file UUID from path parameter
	id := c.Request().Header.Get("uuid")
	if id == "" {
		c.String(http.StatusBadRequest, "No uuid was provided.")
		return fmt.Errorf("No uuid file was provided.")
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		c.String(http.StatusBadRequest, "The uuid provided was not valid.")
		return fmt.Errorf("Bad uuid. %v", err)
	}
	// Verify user authorization
	userClaims, err := echo.ContextGet[*types.ResponseClaims](c, "user")
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	userID, _ := uuid.Parse(userClaims.Subject)
	// Query file details
	fileRepo := repositories.NewFilesRepo(boxed.GetInstance().DbConn)
	file, err := fileRepo.GetByID(uid)
	if err != nil || file.OwnerID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied or file not found"})
	}
	// Serve the file
	return c.File(file.StoragePath)
}
