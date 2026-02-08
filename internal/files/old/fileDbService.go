package files

import (
	"mime/multipart"
	"time"

	"github.com/David/Boxed/repositories"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// saveFileToDb will save the metadata into the database and return the corresponding id.
func saveFileToDatabase(c *pgxpool.Pool, file *multipart.FileHeader, fid, uid uuid.UUID, fpath, tpath string) error {
	mime := file.Header.Get("Content-Type")
	originalName := file.Filename
	fr := repositories.NewFilesRepo(c)
	return fr.Create(&repositories.File{
		ID:            fid,
		OwnerID:       uid,
		OriginalName:  originalName,
		StoragePath:   fpath,
		Size:          file.Size,
		MimeType:      mime,
		ThumbnailPath: &tpath,
		CreatedAt:     time.Now(),
	})
}
