package routes

import (
	"bytes"
	"encoding/json"
	"github.com/marmorag/supateam/internal"
	"github.com/marmorag/supateam/internal/models"
	"github.com/marmorag/supateam/internal/repository"
)

func MustAuthUser(identity string) string {
	userRepository := repository.NewUserRepository("")
	user, err := userRepository.FindOneByIdentity(identity)
	if err != nil {
		panic(err)
	}

	token := models.BuildToken(*user)

	t, err := token.SignedString([]byte(internal.GetConfig().ApplicationSecret))
	if err != nil {
		panic(err)
	}

	return t
}

func MustSerializeReader(b interface{}) *bytes.Reader {
	jsonBody, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(jsonBody)
}
