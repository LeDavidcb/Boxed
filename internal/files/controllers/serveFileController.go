package controllers

import (
	"fmt"
	"net/http"

	boxed "github.com/David/Boxed"
	"github.com/David/Boxed/internal/common/types"
	"github.com/David/Boxed/repositories"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

// ServeFile streams a requested file to the user based on its UUID.
// Authorization is validated to ensure the requesting user owns the file.
//
// Returns:
//   - Responds with the file content if successful.
//   - Responds with HTTP 403 (Forbidden) if the user is not authorized.
//   - Responds with HTTP 400 (Bad Request) or HTTP 401 (Unauthorized) based on validation errors.
func ServeFileController(c *echo.Context) error {
	// Extract file UUID from path parameter
	id := c.Request().Header.Get("uuid")
	if id == "" {
		e := &types.ErrorResponse{
			Code:    types.MissingFields,
			Message: "You must provide a `uuid` entry to serve file.",
		}
		return c.JSON(http.StatusBadRequest, &e)
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		e := &types.ErrorResponse{
			Code:    types.InvalidFields,
			Message: "`uuid` provided is not valid.",
		}
		return c.JSON(http.StatusBadRequest, &e)
	}
	// Verify user authorization
	userClaims, err := echo.ContextGet[*types.ResponseClaims](c, "user")
	if err != nil {
		e := &types.ErrorResponse{
			Code:    types.InternalServerError, // Couldn't get jwt, so it's a middleware error.
			Message: "Error while getting user from jwt, please try again.",
		}
		return c.JSON(http.StatusInternalServerError, &e)
	}
	userID, err := uuid.Parse(userClaims.Subject)
	if err != nil {
		e := &types.ErrorResponse{
			Code:    types.InternalServerError,
			Message: "Internal error while parsing user uuid, please try again.",
		}
		return c.JSON(http.StatusInternalServerError, &e)
	}
	// Query file details
	fileRepo := repositories.NewFilesRepo(boxed.GetInstance().DbConn)
	file, err := fileRepo.GetByID(uid)
	if err != nil {
		e := &types.ErrorResponse{
			Code:    types.ResourceNotFound,
			Message: fmt.Sprintf("Couldn't get any file with uuid: %v", uid.String()),
		}
		return c.JSON(http.StatusBadRequest, &e)
	}
	if file.OwnerID != userID {
		e := &types.ErrorResponse{
			Code:    types.WrongOwner,
			Message: "This user don't own this resource.",
		}
		return c.JSON(http.StatusForbidden, &e)
	}
	// Serve the file
	return c.File(file.StoragePath)
}
