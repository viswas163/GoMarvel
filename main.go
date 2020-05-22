package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/viswas163/MarvelousShipt/routes"
)

// MarvelAPIKeyEnvKey : The OS environment key to fetch the Marvel Developer API Key string
var MarvelAPIKeyEnvKey = "MARVEL_API_KEY"

func main() {
	hasKey := ""

	fmt.Println("Starting Server...")

	for !strings.EqualFold(hasKey, "y") && !strings.EqualFold(hasKey, "n") && !strings.EqualFold(hasKey, "yes") && !strings.EqualFold(hasKey, "no") {
		fmt.Println("Do you have a Marvel Developer API Key? (y/n)")
		fmt.Scanln(&hasKey)
	}

	if strings.EqualFold(hasKey, "y") || strings.EqualFold(hasKey, "yes") {
		userAPIKey := ""
		fmt.Println("Enter the API Key please...")
		fmt.Scanln(&userAPIKey)
		os.Setenv(MarvelAPIKeyEnvKey, userAPIKey)
	}

	routes.HandleRoutes()

	fmt.Println("Started Server... OK!")

	http.ListenAndServe(":3001", nil)
}
