package handlers

import "net/http"

type restrictedHandler struct {

}

func NewRestrictedHandler() *restrictedHandler {
	return &restrictedHandler{}
}

func (h *restrictedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("hello from the secret area!"))
	w.WriteHeader(http.StatusOK)
}