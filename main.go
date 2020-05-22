package main

import (
	"fmt"
	"net/http"

	"github.com/viswas163/MarvelousShipt/routes"
)

func main() {

	fmt.Println("Starting Server...")

	routes.HandleRoutes()

	fmt.Println("Started Server... OK!")

	http.ListenAndServe(":3001", nil)
}
