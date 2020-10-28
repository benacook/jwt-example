package handlers

import (
	"github.com/benacook/jwt-example/auth"
	"log"
	"net/http"
)

//======================================================================================

//blank type for class

type publicHandler struct {

}

//======================================================================================

//generates an instance of the public handler

func NewGenerateTokenHandler() *publicHandler {
	return &publicHandler{}
}

//======================================================================================

//generates and returns a new token using the secret setup in the init of the auth
//module

func (h *publicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	newToken, err := auth.GenerateToken()
	tokenIdentifier := "token="
	responseBody := tokenIdentifier + newToken
	if err != nil {
		log.Println(err)
		_, _ = w.Write([]byte("error making token"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte(responseBody))
	w.WriteHeader(http.StatusOK)
}