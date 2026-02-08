package services

import (
	"mime/multipart"
	"strings"
	"time"

	"github.com/David/Boxed/repositories"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// saveFileToDb will save the metadata into the database and return the corresponding id.
func SaveFileToDatabase(c *pgxpool.Pool, file *multipart.FileHeader, fid, uid uuid.UUID, fpath string, thumbnail uuid.UUID) error {
	mime := file.Header.Get("Content-Type")
	originalName := strings.Split(file.Filename, ".")[0]
	fr := repositories.NewFilesRepo(c)
	return fr.Create(&repositories.File{
		ID:           fid,
		OwnerID:      uid,
		OriginalName: originalName,
		StoragePath:  fpath,
		Size:         file.Size,
		MimeType:     mime,
		ThumbnailId:  thumbnail,
		CreatedAt:    time.Now(),
	})
}
