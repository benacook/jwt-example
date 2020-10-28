package controller

import (
	"github.com/benacook/jwt-example/auth"
	"github.com/benacook/jwt-example/controller/handlers"
	"net/http"
)

//======================================================================================

//exaple of using the middleware to resrtict access to a handler function

func RegisterControllers(){
	rh := handlers.NewRestrictedHandler()
	ph := handlers.NewGenerateTokenHandler()
	http.Handle("/secret-area",auth.Middleware(rh))
	http.Handle("/new-token",ph)
}
