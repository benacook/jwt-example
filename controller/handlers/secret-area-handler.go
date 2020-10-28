package handlers

import "net/http"

//======================================================================================

//blank type for class

type restrictedHandler struct {

}

//======================================================================================

//generates an instance of the private handler

func NewRestrictedHandler() *restrictedHandler {
	return &restrictedHandler{}
}

//======================================================================================

//the "secret" data

func (h *restrictedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("hello from the secret area!"))
	w.WriteHeader(http.StatusOK)
}