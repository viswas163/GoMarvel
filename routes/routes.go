package routes

import (
	"net/http"

	v1 "github.com/viswas163/MarvelousShipt/api/v1"
)

var (
	// BaseURL : Marvel API base URL
	BaseURL = "https://gateway.marvel.com"
)

// HandleRoutes : Handles all api routes
func HandleRoutes() {

	http.HandleFunc("/", v1.Hello)
	http.HandleFunc("/hello", v1.Hello)
	http.HandleFunc("/headers", v1.Headers)
	http.HandleFunc("/image", v1.Image)

	return
}
