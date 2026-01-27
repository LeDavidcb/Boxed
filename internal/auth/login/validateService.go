package login

import (
	"log"
	"strings"
	"time"

	boxed "github.com/David/Boxed"
	"github.com/David/Boxed/internal/common/types"
	"github.com/David/Boxed/repositories"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func validate(u *userLoginRequest, c *pgxpool.Pool) (*loginResponse, error) {
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
	return &loginResponse{SignedJwt: sig}, nil
}
