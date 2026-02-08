package controllers

import (
	"fmt"
	"net/http"

	boxed "github.com/David/Boxed"
	"github.com/David/Boxed/repositories"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

// ServeThumbnailController streams a requested thumbnail bassed on its UUID.
type thumbnailUUID struct {
	ID string `query:"id"`
}

func ServeThumbnailController(c *echo.Context) error {
	var rawThumbnailUUID thumbnailUUID
	err := c.Bind(&rawThumbnailUUID)
	if err != nil {
		return c.String(http.StatusBadRequest, "UUID for thumbnail must be provided.")
	}
	thumbnailUUID, err := uuid.Parse(rawThumbnailUUID.ID)
	if err != nil {

		return c.String(http.StatusBadRequest, "An invalid UUID was provided. ")
	}

	repository := repositories.NewThumbnailRepository(boxed.GetInstance().DbConn)
	thumbnail, err := repository.GetByID(thumbnailUUID)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("No thumbnail with %v UUID.", thumbnailUUID))
	}
	if thumbnail.StoragePath == "" {

		return c.String(http.StatusBadRequest, fmt.Sprintf("File { name: %v} Does not have a thumbnail.", thumbnail.OriginalName))
	}
	return c.File(thumbnail.StoragePath)
}
