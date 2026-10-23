package main

import (
	"log"
	"net/http"

	"github.com/brunoantunesc/sezzle-calculator/backend/api"
)

func main() {
	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", api.NewHandler()))
}
