package models

import (
	"github.com/dgrijalva/jwt-go"
	gojwt "github.com/golang-jwt/jwt/v4"
	"github.com/marmorag/supateam/internal/middleware/auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	Id            primitive.ObjectID  `bson:"_id" json:"id,omitempty"`
	Identity      string              `json:"-"`
	Name          string              `json:"name,omitempty"`
	Authorization map[string][]string `json:"authorization"`
}

type CreateUserRequest struct {
	Name          string `validate:"required"`
	Identity      string `validate:"required"`
	Authorization map[string][]string
}

type UpdateUserRequest struct {
	Name          string              `validate:"required"`
	Identity      string              `validate:"alphanum"`
	Authorization map[string][]string `validate:"required"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func BuildToken(user User) *jwt.Token {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &auth.ApplicationClaim{})
	claims := token.Claims.(*auth.ApplicationClaim)
	claims.UserName = user.Name
	claims.UserId = user.Id.Hex()
	claims.UserAuthorization = user.Authorization
	claims.ExpiresAt = time.Now().Add(time.Hour * 96).Unix()

	token.Claims = claims
	return token
}

func GetUserIdFromToken(token *gojwt.Token) string {
	claims := token.Claims.(*auth.ApplicationClaim)
	return claims.UserId
}
