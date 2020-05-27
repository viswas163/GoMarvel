package v1

import (
	"encoding/json"
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
)

// RetrieveAllCharacters : Gets all characters
func RetrieveAllCharacters(w http.ResponseWriter, r *http.Request) {
	characters, err := GetAllCharacters()
	if err != nil {
		http.Error(w, "Error getting characters", 500)
	}

	json, err := json.Marshal(characters.Data.Results)
	if err != nil {
		http.Error(w, "Error marshaling json", 500)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// GetAllCharacters : Gets all characters from API or response cache
func GetAllCharacters() (models.CharacterDataWrapper, error) {
	chars, err := GetCharactersByOffset(0)
	if err != nil {
		return models.CharacterDataWrapper{}, err
	}

	for i := 100; i < chars.Data.Total; i += 100 {
		ch, err := GetCharactersByOffset(i)
		if err != nil {
			return models.CharacterDataWrapper{}, err
		}
		chars.Data.Results = append(chars.Data.Results, ch.Data.Results...)
	}
	models.SetAllCharactersByName(chars.Data.Results)

	return chars, nil
}

// GetCharactersByOffset : Retrieves characters from API or response cache with offset
func GetCharactersByOffset(offset int) (models.CharacterDataWrapper, error) {

	characters := []byte{}
	db.GetInstance().Get(fmt.Sprint(allCharactersJSONDBKey, offset), &characters)
	var err error

	if len(characters) <= 0 {
		fmt.Println("Fetching API results...")
		characters, err = RunAPIWithoutParam(allCharactersAPI, offset)
		if err != nil {
			return models.CharacterDataWrapper{}, err
		}
		db.GetInstance().Put(fmt.Sprint(allCharactersJSONDBKey, offset), characters)
	} else {
		fmt.Println("Fetching result from Cache!")
	}

	chars, err := models.GetAllCharacters(characters)
	if err != nil {
		return models.CharacterDataWrapper{}, err
	}

	return chars, err
}
