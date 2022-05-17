package main

import (
	"log"
	"net/http"

	"github.com/mysticis/go-dcktst-demo/middleware"
)

func main() {

	srv := middleware.NewServer()

	log.Println("listening on port :8080")
	log.Fatal(http.ListenAndServe(":8080", srv))
}
