package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersRepository interface {
	Create(user *User) error
	GetByID(id uuid.UUID) (*User, error)
	Update(user *User) error
	Delete(id uuid.UUID) error
}
type User struct {
	ID           uuid.UUID `db:"id"`
	Username     string    `db:"username"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	FolderPath   string    `db:"folder_path"`
	CreatedAt    time.Time `db:"created_at"`
}

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(d *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: d}
}

// Creates an user in the `db`.
// Create inserts a new user into the `users` table.
//
// Parameters:
//   - u (*User): The user to insert with fields for ID, username, email, hashed password, etc.
//
// Returns:
//   - error: An error if the operation does not succeed.
func (s *UserRepo) Create(u *User) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	query := `INSERT INTO users (id, username, email, password_hash, folder_path, created_at) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.Exec(context.Background(), query, u.ID, u.Username, u.Email, u.PasswordHash, u.FolderPath, u.CreatedAt)
	return err

}

// GetByID retrieves a user from the `users` table by their unique ID.
//
// Parameters:
//   - id (uuid.UUID): The unique identifier for the user.
//
// Returns:
//   - (*User, error): A reference to the user object if found; otherwise, an error.
func (s *UserRepo) GetByID(id uuid.UUID) (*User, error) {
	user := &User{}
	query := `
        SELECT id, username, email, password_hash, created_at, folder_path
        FROM users
        WHERE id = $1`
	err := s.db.QueryRow(context.Background(), query, id).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.FolderPath)
	return user, err
}

func (s *UserRepo) GetByEmail(email string) (*User, error) {
	user := &User{}
	query := `SELECT id, username, email, password_hash FROM users WHERE email = $1`
	err := s.db.QueryRow(context.Background(), query, email).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
	return user, err
}
