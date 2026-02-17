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

// ServeThumbnailController streams a requested thumbnail bassed on its UUID.
func ServeThumbnailController(c *echo.Context) error {
	uid := c.Request().Header.Get("uuid")

	if uid == "" {
		e := &types.ErrorResponse{
			Code:    types.MissingFields,
			Message: "You must provide a `uuid` entry to serve a thumbnail",
		}
		return c.JSON(http.StatusBadRequest, &e)
	}
	thumbnailUUID, err := uuid.Parse(uid)
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

	repository := repositories.NewThumbnailRepository(boxed.GetInstance().DbConn)
	thumbnail, err := repository.GetByID(thumbnailUUID)
	if err != nil {

		var pge *pgconn.PgError
		if errors.As(err, &pge) || errors.As(err, &pgx.ErrNoRows) {
			e := &types.ErrorResponse{
				Code:    types.ResourceNotFound,
				Message: fmt.Sprintf("No thumbnail with uuid: %v", uid),
			}
			return c.JSON(http.StatusBadRequest, &e)
		}
	}
	if thumbnail.StoragePath == "" {
		e := &types.ErrorResponse{
			Code:    types.ResourceNotFound,
			Message: fmt.Sprintf("No thumbnail to serve with uuid: %v", uid),
		}
		return c.JSON(http.StatusBadRequest, &e)
	}

	if err != nil {
		e := &types.ErrorResponse{
			Code:    types.ResourceNotFound,
			Message: fmt.Sprintf("Couldn't get any thumbnail with uuid: %v", thumbnailUUID.String()),
		}
		return c.JSON(http.StatusBadRequest, &e)
	}
	if thumbnail.OwnerId != userID {
		e := &types.ErrorResponse{
			Code:    types.WrongOwner,
			Message: "This user don't own this resource.",
		}
		return c.JSON(http.StatusForbidden, &e)
	}
	return c.File(thumbnail.StoragePath)
}
