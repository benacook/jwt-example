package controller

import (
	"github.com/benacook/jwt-example/auth"
	"github.com/benacook/jwt-example/controller/handlers"
	"net/http"
)

func RegisterControllers(){
	rh := handlers.NewRestrictedHandler()
	ph := handlers.NewPublicHandler()
	http.Handle("/secret-area",auth.Middleware(rh))
	http.Handle("/new-token",ph)
}
