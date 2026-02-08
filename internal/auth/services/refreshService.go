package services

import (
	"time"

	boxed "github.com/David/Boxed"
	"github.com/David/Boxed/internal/common/types"
	"github.com/David/Boxed/repositories"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// ReSignJwt creates a new JWT for a user identified by their UUID.
// This function retrieves the user's information from the database, constructs the JWT claims,
// and signs the token with a secret key.
//
// Parameters:
//   - id (uuid.UUID): The UUID of the user for whom the JWT is created.
//
// Returns:
//   - (string, error): The newly signed JWT as a string; an error if token signing or database access fails.
//
// Errors:
//   - Returns an error if the user does not exist in the database.
//   - Returns an error if signing the token fails.
func ReSignJwt(id uuid.UUID) (string, error) {
	con := boxed.GetInstance().DbConn
	ur := repositories.NewUserRepo(con)

	user, err := ur.GetByID(id)
	if err != nil {
		return "", err
	}
	claims := &types.ResponseClaims{
		Name: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Subject:   user.ID.String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	sig, err := token.SignedString([]byte(boxed.GetInstance().JwtSecret))
	return sig, err

}
