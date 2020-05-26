package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	v1 "github.com/viswas163/MarvelousShipt/api/v1"
	"github.com/viswas163/MarvelousShipt/db"
	"github.com/viswas163/MarvelousShipt/routes"
)

func main() {

	router := &routes.Router{R: mux.NewRouter()}
	router.HandleRoutes()

	server := &http.Server{
		Addr:    ":3001",
		Handler: router.R,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")

	// Initialize Auth Client
	v1.InitAuthClient()

	log.Print("Private Key Initialized!")
	log.Print("Initializing DB...")

	// Initialize DB
	_, err := db.Open("db/marvel.db")
	if err != nil {
		done <- os.Interrupt
	}

	v1.GetAllCharacters()

	log.Print("DB Initialized!")
	log.Print("Server Ready for incoming requests...")

	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cleanup()
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}

func cleanup() {
	fmt.Println("\r- Cleaning up...")
	if _, bool := os.LookupEnv(v1.MarvelPrivateAPIKeyEnvKey); bool {
		os.Unsetenv(v1.MarvelPrivateAPIKeyEnvKey)
	}
	if err := db.GetInstance().Close(); err != nil {
		fmt.Println("Error closing DB")
	}
	fmt.Println("\r- Good Bye!")
}
