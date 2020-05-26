package routes

import (
	"github.com/gorilla/mux"
	v1 "github.com/viswas163/MarvelousShipt/api/v1"
)

var (
	// BaseURL : Marvel API base URL
	BaseURL = "https://gateway.marvel.com"
)

type Router struct {
	R *mux.Router
}

// HandleRoutes : Handles all api routes
func (router *Router) HandleRoutes() *Router {
	router.R.HandleFunc("/", v1.Hello)
	router.R.HandleFunc("/hello", v1.Hello)
	router.R.HandleFunc("/headers", v1.Headers)
	router.R.HandleFunc("/image", v1.Image)
	router.R.HandleFunc("/allcharacters", v1.GetAllCharacters)

	return router
}
