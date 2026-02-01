package refresh

import (
	"time"

	boxed "github.com/David/Boxed"
	"github.com/David/Boxed/internal/common/types"
	"github.com/David/Boxed/repositories"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

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
