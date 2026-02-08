package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// File model represents the structure of the "files" table.
type File struct {
	ID           uuid.UUID `db:"id"`
	OwnerID      uuid.UUID `db:"owner_id"`
	OriginalName string    `db:"original_name"`
	StoragePath  string    `db:"storage_path"`
	Size         int64     `db:"size"`
	MimeType     string    `db:"mime_type"`
	ThumbnailId  uuid.UUID `db:"thumbnail_id"`
	CreatedAt    time.Time `db:"created_at"`
}

// FilesRepository interface exposes CRUD operations for files.
type FilesRepository interface {
	Create(file *File) error
	GetByID(id uuid.UUID) (*File, error)
	GetByOwnerID(ownerID uuid.UUID) ([]File, error)
	Delete(id uuid.UUID) error
}

// FilesRepo implements the FilesRepository interface using pgx for PostgreSQL interaction.
type FilesRepo struct {
	db *pgxpool.Pool
}

// NewFilesRepo initializes a new instance of FilesRepo.
func NewFilesRepo(db *pgxpool.Pool) *FilesRepo {
	return &FilesRepo{db: db}
}

// Create inserts a new file record in the "files" table.
// Create adds a new file entry into the `files` table.
//
// Parameters:
//   - file (*File): The file metadata to insert.
//
// Returns:
//   - error: An error if database insertion fails.
func (r *FilesRepo) Create(file *File) error {
	if file.ID == uuid.Nil {
		file.ID = uuid.New()
	}
	query := `
        INSERT INTO files (id, owner_id, original_name, storage_path, size, mime_type, thumbnail_id, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(context.Background(), query, file.ID, file.OwnerID, file.OriginalName, file.StoragePath,
		file.Size, file.MimeType, file.ThumbnailId, file.CreatedAt)
	return err
}

// GetByID retrieves a file record by its ID.
// GetByID retrieves a file entry by its unique ID.
//
// Parameters:
//   - id (uuid.UUID): The unique identifier for the file record.
//
// Returns:
//   - (*File, error): A pointer to the file's metadata if found; otherwise, an error.
func (r *FilesRepo) GetByID(id uuid.UUID) (*File, error) {
	file := &File{}
	query := `SELECT id, owner_id, original_name, storage_path, size, mime_type, thumbnail_id, created_at
              FROM files WHERE id = $1`
	err := r.db.QueryRow(context.Background(), query, id).
		Scan(&file.ID, &file.OwnerID, &file.OriginalName, &file.StoragePath, &file.Size,
			&file.MimeType, &file.ThumbnailId, &file.CreatedAt)
	return file, err
}

// GetByOwnerID retrieves all files owned by a specific user ID.
func (r *FilesRepo) GetByOwnerID(ownerID uuid.UUID) ([]File, error) {
	query := `SELECT id, owner_id, original_name, storage_path, size, mime_type, thumbnail_id, created_at
              FROM files WHERE owner_id = $1`
	rows, err := r.db.Query(context.Background(), query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := []File{}
	for rows.Next() {
		file := File{}
		err := rows.Scan(&file.ID, &file.OwnerID, &file.OriginalName, &file.StoragePath, &file.Size,
			&file.MimeType, &file.ThumbnailId, &file.CreatedAt)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

// Delete removes a file record by its ID.
func (r *FilesRepo) Delete(id uuid.UUID) error {
	query := "DELETE FROM files WHERE id = $1"
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}
