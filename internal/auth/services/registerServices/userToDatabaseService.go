package registerservices

import (
	"time"

	"github.com/David/Boxed/internal/auth/types"
	"github.com/David/Boxed/repositories"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// createUserDb creates a new user in the database with the given registration details.
// This includes hashing the password, assigning a unique ID, generating a folder path, and saving the user record.
//
// Parameters:
//   - c (*pgxpool.Pool): A database connection pool.
//   - u (*userRegisterRequest): Struct containing the nickname, email, and raw password of the new user.
//
// Returns:
//   - error: An error if the user creation process (e.g., password hashing, folder creation, or database insertion) fails.
func CreateUserInDatabase(c *pgxpool.Pool, u *types.UserRegisterRequest) error {
	user := new(repositories.User)
	user.ID = uuid.New()
	user.Username = u.Nickname
	user.Email = u.Email
	user.CreatedAt = time.Now()
	// Encript the password using bycript
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hash)
	// Create the folder
	path, err := CreateDirectory(user.ID)
	if err != nil {
		return err
	}
	user.FolderPath = path
	// Save user in the database
	if err = repositories.NewUserRepo(c).Create(user); err != nil {
		return err
	}
	return nil

}
