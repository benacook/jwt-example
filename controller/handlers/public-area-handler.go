package handlers

import (
	"github.com/benacook/jwt-example/auth"
	"log"
	"net/http"
)

type publicHandler struct {

}

func NewPublicHandler() *publicHandler {
	return &publicHandler{}
}

func (h *publicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	newToken, err := auth.GenerateToken()
	if err != nil {
		log.Println(err)
		_, _ = w.Write([]byte("error making token"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte(newToken))
	w.WriteHeader(http.StatusOK)
}