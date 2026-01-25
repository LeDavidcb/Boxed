package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RefreshToken represents the structure of the "refresh_tokens" table.
type RefreshToken struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	TokenHash string    `db:"token_hash"`
	ExpiresAt time.Time `db:"expires_at"`
	Revoked   bool      `db:"revoked"`
	CreatedAt time.Time `db:"created_at"`
}

// RefreshTokensRepository defines CRUD operations for the "refresh_tokens" table.
type RefreshTokensRepository interface {
	Create(token *RefreshToken) error
	GetByUserID(userID uuid.UUID) ([]RefreshToken, error)
	DeleteByID(id uuid.UUID) error
	RevokeByID(id uuid.UUID) error
}

// RefreshTokensRepo implements the RefreshTokensRepository interface.
type RefreshTokensRepo struct {
	db *pgxpool.Pool
}

// NewRefreshTokensRepo initializes a new instance of RefreshTokensRepo.
func NewRefreshTokensRepo(db *pgxpool.Pool) *RefreshTokensRepo {
	return &RefreshTokensRepo{db: db}
}

// Create inserts a new refresh token record in the "refresh_tokens" table.
func (r *RefreshTokensRepo) Create(token *RefreshToken) error {
	if token.ID == uuid.Nil {
		token.ID = uuid.New()
	}
	query := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, revoked, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(context.Background(), query, token.ID, token.UserID, token.TokenHash, token.ExpiresAt, token.Revoked, token.CreatedAt)
	return err
}

// GetByUserID retrieves all refresh tokens for a specific user.
func (r *RefreshTokensRepo) GetByUserID(userID uuid.UUID) ([]RefreshToken, error) {
	query := `SELECT id, user_id, token_hash, expires_at, revoked, created_at
			  FROM refresh_tokens WHERE user_id = $1`
	rows, err := r.db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tokens := []RefreshToken{}
	for rows.Next() {
		token := RefreshToken{}
		err := rows.Scan(&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt, &token.Revoked, &token.CreatedAt)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

// DeleteByID deletes a refresh token by its ID.
func (r *RefreshTokensRepo) DeleteByID(id uuid.UUID) error {
	query := "DELETE FROM refresh_tokens WHERE id = $1"
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}

// RevokeByID marks a refresh token as revoked by its ID.
func (r *RefreshTokensRepo) RevokeByID(id uuid.UUID) error {
	query := "UPDATE refresh_tokens SET revoked = true WHERE id = $1"
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}
