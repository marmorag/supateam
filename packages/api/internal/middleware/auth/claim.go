package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type ApplicationClaim struct {
	ExpiresAt int64 `json:"exp"`
	UserName string `json:"userName"`
	UserId string `json:"userId"`
	UserAuthorization map[string][]string `json:"userAuthorization"`
}

func (a ApplicationClaim) Valid() error {
	validationErr := new(jwt.ValidationError)
	now := time.Now()

	if a.ExpiresAt == 0 || now.After(time.Unix(a.ExpiresAt, 0)) {
		delta := time.Unix(now.Unix(), 0).Sub(time.Unix(a.ExpiresAt, 0))
		validationErr.Inner = fmt.Errorf("token is expired by %v", delta)
		validationErr.Errors |= jwt.ValidationErrorExpired
	}

	return validationErr
}

func ParseToken()  {

}