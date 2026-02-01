package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/David/Boxed/internal/common/fn"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
func (r *RefreshTokensRepo) GetByHashToken(h string) (*RefreshToken, error) {
	if h == "" {
		return nil, errors.New("token hash must not be empty")
	}
	log.Println("PROVIDED H", h)
	query := "SELECT id, user_id, token_hash, expires_at, created_at FROM refresh_tokens WHERE token_hash = $1"
	row := r.db.QueryRow(context.Background(), query, h)
	response := &RefreshToken{}
	err := row.Scan(&response.ID, &response.UserID, &response.TokenHash, &response.ExpiresAt, &response.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			// No rows found; handle differently
			return nil, fmt.Errorf("refresh token not found for hash: %s", h)
		}
		// Other database errors
		log.Printf("Database error fetching token by hash: %v", err)
		return nil, err
	}
	return response, nil
}
func (r *RefreshTokensRepo) RegenerateToken(h string) (*struct {
	Useruuid uuid.UUID
	NewHash  string
}, error) {
	token, err := r.GetByHashToken(h)
	if err != nil {
		return nil, err
	}
	hash, err := fn.GenerateRTHash(32)
	if err != nil {
		return nil, err
	}
	// Set
	token.TokenHash = hash
	token.CreatedAt = time.Now()
	token.ExpiresAt = time.Now().Add(time.Hour * 24 * 7)
	// Update it to the database.
	_, err = r.db.Exec(context.Background(), "UPDATE refresh_tokens SET token_hash = $1, created_at = $2, expires_at = $3 WHERE id = $4", token.TokenHash, token.CreatedAt, token.ExpiresAt, token.ID)
	return &struct {
		Useruuid uuid.UUID
		NewHash  string
	}{Useruuid: token.UserID, NewHash: hash}, err
}
