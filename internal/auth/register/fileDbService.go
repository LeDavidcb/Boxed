package register

import (
	"time"

	"github.com/David/Boxed/repositories"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func createUserDb(c *pgxpool.Pool, u *user) error {
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
