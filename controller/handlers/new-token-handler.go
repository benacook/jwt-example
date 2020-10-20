package handlers

import (
	"github.com/benacook/jwt-example/auth"
	"log"
	"net/http"
)

type publicHandler struct {
	s int
}

func NewGenerateTokenHandler() *publicHandler {
	return &publicHandler{}
}

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