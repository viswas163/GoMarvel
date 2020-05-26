package v1

import (
	"fmt"
	"net/http"

	"github.com/viswas163/MarvelousShipt/db"
	"github.com/viswas163/MarvelousShipt/models"
)

var (
	// AllCharacters : Contains all characters
	AllCharacters models.AllCharacters

	allCharactersAPI       = "characters"
	allCharactersJSONDBKey = "charactersJSON"
	allCharactersDBKey     = "characters"
)

// GetAllCharacters : Gets all characters
func GetAllCharacters(w http.ResponseWriter, r *http.Request) {
	characters := []byte{}
	var err error
	db.GetInstance().Get(allCharactersJSONDBKey, &characters)

	if len(characters) <= 0 {
		fmt.Println("Fetching Characters from API...")

		characters, err = RunAPI(allCharactersAPI)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		db.GetInstance().Put(allCharactersJSONDBKey, characters)
	} else {
		fmt.Println(allCharactersAPI, "Fetched from Cache!")
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(characters)
}
