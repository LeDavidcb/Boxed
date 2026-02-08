package services

import (
	"log"
	"strings"
	"time"

	boxed "github.com/David/Boxed"
	"github.com/David/Boxed/internal/common/types"
	"github.com/David/Boxed/internal/common/utils"
	"github.com/David/Boxed/repositories"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// validate checks user credentials by comparing the provided password and stored hashed password.
// If credentials are valid, it generates a signed JWT and a refresh token.
//
// Parameters:
//   - u (*userLoginRequest): The user-provided login details, which include an email and password.
//   - c (*pgxpool.Pool): The database connection pool.
//
// Returns:
//   - (*loginResponse, error): A struct containing the signed JWT and refresh token on success; error otherwise.
//
// Errors:
//   - Returns an error if user credentials do not match or if database access fails.
func Validate(u *types.UserLoginRequest, c *pgxpool.Pool) (*types.LoginResponse, error) {
	repo := repositories.NewUserRepo(c)
	user, err := repo.GetByEmail(u.Email)
	if err != nil {
		return nil, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(u.Password)); err != nil {
		return nil, err
	}
	// Generate the jwt
	claims := &types.ResponseClaims{
		Name: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Subject:   user.ID.String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key := strings.Trim(boxed.GetInstance().JwtSecret, " ")
	sig, err := token.SignedString([]byte(key))
	if err != nil {
		log.Printf("[ERROR] Token signing failed: %v\n", err)
		return nil, err
	}
	// Generate the refreshToken and save it to the database
	hash, err := utils.GenerateRTHash(32)
	if err != nil {
		return nil, err
	}
	rtr := repositories.NewRefreshTokensRepo(boxed.GetInstance().DbConn)
	err = rtr.Create(&repositories.RefreshToken{
		ID:        uuid.New(),
		TokenHash: hash,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
		UserID:    user.ID,
		Revoked:   false,
	})
	if err != nil {
		return nil, err
	}

	return &types.LoginResponse{SignedJwt: sig, RefreshToken: hash}, nil
}
