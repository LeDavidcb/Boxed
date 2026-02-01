package refresh

import (
	"log"

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
		log.Println("ERRORR WhiLE GETTING ID:", err)
		return "", err
	}
	claims := new(types.ResponseClaims)
	claims.Name = user.Username
	claims.Subject = user.ID.String()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	sig, err := token.SignedString([]byte(boxed.GetInstance().JwtSecret))
	log.Println("ERRORR WhiLE GETTING ID:", err)
	return sig, err

}
