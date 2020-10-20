package main

import (
	"github.com/benacook/jwt-example/controller"
	"log"
	"net/http"
)

func main() {
	controller.RegisterControllers()
	println("serving on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}

