package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	v1 "github.com/viswas163/MarvelousShipt/api/v1"
	"github.com/viswas163/MarvelousShipt/routes"
)

func main() {

	fmt.Println("\nStarting Server...")

	setupCloseHandler()

	v1.InitAuthClient()

	resp, _ := v1.RunAuth("characters")
	fmt.Println(string(resp))

	routes.HandleRoutes()
	fmt.Println("Started Server... OK!")
	http.ListenAndServe(":3001", nil)
}

func setupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Server shutdown by user")
		cleanup()
		os.Exit(0)
	}()
}

func cleanup() {
	fmt.Println("\r- Cleaning up...")
	os.Unsetenv(v1.MarvelPrivateAPIKeyEnvKey)
	fmt.Println("\r- Good Bye!")
}
