package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
)

type Thumbnail struct {
	ID           uuid.UUID `db:"id"`
	OwnerId      uuid.UUID `db:"owner_id"`
	OriginalName string    `db:"original_name"`
	StoragePath  string    `db:"storage_path"`
}

type ThumbnailRepositoryInterface interface {
	Create(t *Thumbnail) error
	GetByID(id uuid.UUID) (Thumbnail, error)
	DeleteByID(id uuid.UUID) error
	UpdateByID(t *Thumbnail) error
}
type ThumbnailRepository struct {
	db *pgxpool.Pool
}

func NewThumbnailRepository(db *pgxpool.Pool) *ThumbnailRepository {
	return &ThumbnailRepository{db: db}
}

func (r *ThumbnailRepository) Create(t *Thumbnail) error {
	query := "INSERT INTO thumbnails (id, owner_id, original_name, storage_path) VALUES ($1, $2, $3, $4)"
	_, err := r.db.Exec(context.Background(), query, t.ID, t.OwnerId, t.OriginalName, t.StoragePath)
	return err
}

func (r *ThumbnailRepository) GetByID(id uuid.UUID) (*Thumbnail, error) {
	query := "SELECT id, owner_id, original_name, storage_path FROM thumbnails WHERE id = $1"
	row := r.db.QueryRow(context.Background(), query, id)

	t := &Thumbnail{}
	err := row.Scan(&t.ID, &t.OwnerId, &t.OriginalName, &t.StoragePath)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *ThumbnailRepository) DeleteByID(id uuid.UUID) error {
	query := "DELETE FROM thumbnails WHERE id = $1"
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}

// UpdateByID updates the `thumbnails` table with the values provided in the `Thumbnail` struct.
//
// Parameters:
// - t (*Thumbnail): A pointer to a `Thumbnail` struct.
//   - `t.ID` (UUID): The unique identifier for the thumbnail. A valid non-nil UUID is required for this update to succeed.
//   - `t.original_name` and `t.storage_path` (strings): Fields to be updated. If either field is an empty string or uninitialized, it will be ignored.
//
// Returns:
// - error: Returns an error if:
//   - `t.ID` is not a valid UUID.
//   - The query execution fails (e.g., database connectivity issues).
//
// Notes:
// - All other fields in the `Thumbnail` struct are ignored by this operation.
// - Ensure that database connectivity is properly configured before calling this function.
func (r *ThumbnailRepository) UpdateByID(t *Thumbnail) error {
	if t.ID == uuid.Nil {
		return fmt.Errorf("UUID not provided.")
	}

	updates := []string{}
	args := []any{}

	if t.OriginalName != "" {
		updates = append(updates, "original_name = $1")
		args = append(args, t.OriginalName)
	}

	if t.StoragePath != "" {
		updates = append(updates, "storage_path = $2")
		args = append(args, t.StoragePath)
	}

	if len(updates) == 0 {
		return nil // No fields to update
	}

	query := fmt.Sprintf("UPDATE thumbnails SET %s WHERE id = $3", strings.Join(updates, ", "))
	args = append(args, t.ID)

	_, err := r.db.Exec(context.Background(), query, args...)
	return err
}
