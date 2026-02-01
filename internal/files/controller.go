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

// SendFile uploads a single file for the authenticated user and saves both the file and its metadata to disk and the database.
//
// Returns:
//   - Responds with HTTP 201 (Created) on success.
//   - Responds with HTTP 400 (Bad Request) if the file or user info is invalid.
//   - Returns an error if saving the file or metadata fails.
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

// SendFiles allows authenticated users to upload multiple files at once, saving file data and metadata to database and disk.
//
// Returns:
//   - Responds with HTTP 201 (Created) after successfully processing all files.
//   - Responds with HTTP 400 (Bad Request) if form data or file inputs are invalid.
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

// GetFile retrieves metadata for a specific file identified by the UUID provided in the request header.
//
// Returns:
//   - Responds with HTTP 200 (OK) and the file metadata as JSON.
//   - Responds with HTTP 400 (Bad Request) if the UUID is invalid or the file could not be found.
//   - Returns an error if there are issues querying the database.
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

// DeleteFile removes a file identified by its UUID, along with its metadata, from both the database and storage.
//
// Returns:
//   - Responds with HTTP 200 (OK) for successful deletion.
//   - Responds with HTTP 400 (Bad Request) if the UUID is invalid or the file does not exist.
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

// GetFiles retrieves all files (Metadata) owned by the authenticated user and returns their metadata.
//
// Returns:
//   - Responds with HTTP 200 (OK) along with a JSON payload containing file metadata.
//   - Responds with HTTP 404 (Not Found) if no files exist for the user.
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

// ServeFile streams a requested file to the user based on its UUID.
// Authorization is validated to ensure the requesting user owns the file.
//
// Returns:
//   - Responds with the file content if successful.
//   - Responds with HTTP 403 (Forbidden) if the user is not authorized.
//   - Responds with HTTP 400 (Bad Request) or HTTP 401 (Unauthorized) based on validation errors.
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
