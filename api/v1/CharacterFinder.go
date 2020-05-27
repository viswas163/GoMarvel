package v1

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
	"github.com/juliangruber/go-intersect"
	"github.com/viswas163/MarvelousShipt/models"
)

// RetrieveCharacter : Retrieves the character with the name given in URL
func RetrieveCharacter(w http.ResponseWriter, r *http.Request) {
	key, ok := mux.Vars(r)["characterId"]
	if !ok {
		http.Error(w, "Invalid param", 400)
	}

	// fmt.Print(key)
	character, _ := models.GetCharacter(key)

	json, err := json.Marshal(character)
	if err != nil {
		http.Error(w, "Error marshaling json", 500)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// GetCommonComics : Retrieves the common comics between characters passed in URL
func GetCommonComics(w http.ResponseWriter, r *http.Request) {
	key1, ok := mux.Vars(r)["characterId"]
	if !ok {
		http.Error(w, "Invalid param", 400)
	}
	key2, ok := mux.Vars(r)["characterId2"]
	if !ok {
		http.Error(w, "Invalid param", 400)
	}

	char1, _ := models.GetCharacter(key1)
	char2, _ := models.GetCharacter(key2)

	PushChars(char1, char2)
	commonComics := GetCommon(char1, char2)

	json, err := json.Marshal(commonComics)
	if err != nil {
		http.Error(w, "Error marshaling json", 500)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// PushChars : Pushes the comics of given characters
func PushChars(char1 models.Character, char2 models.Character) {
	PushComicsOfChar(char1)
	PushComicsOfChar(char2)
}

// GetCommon : Gets the common comics between two characters
func GetCommon(char1 models.Character, char2 models.Character) []models.Comic {
	l1, ok1 := models.ComicsByCharacter.Load(char1.Name)
	l2, ok2 := models.ComicsByCharacter.Load(char2.Name)
	if !ok1 || !ok2 {
		return nil
	}
	ids := intersect.Hash(l1, l2)
	s := reflect.ValueOf(ids)

	var comics []models.Comic

	for i := 0; i < s.Len(); i++ {
		id := s.Index(i).Interface().(int)
		com, ok := models.ComicsByID.Load(id)
		if !ok {
			continue
		}

		comic, ok := com.(models.Comic)
		if !ok {
			continue
		}
		comics = append(comics, comic)
	}
	return comics
}

// PushComicsOfChar : Retrieves and pushes all comics of character to storage
func PushComicsOfChar(char models.Character) error {
	comics, err := GetAllComicsOfChar(char)
	if err != nil {
		return err
	}
	if err := models.SetAllComicsByCharName(char.Name, comics.Data.Results); err != nil {
		return err
	}
	return nil
}

// EmmaStorm : Retrieves the common comics between Storm and Emma Frost
func EmmaStorm(w http.ResponseWriter, r *http.Request) {
	var emma models.Character
	var storm models.Character

	t1, _ := models.CharactersByName.Load("Emma Frost")
	emma = t1.(models.Character)
	t2, _ := models.CharactersByName.Load("Storm")
	storm = t2.(models.Character)

	PushChars(emma, storm)
	common := GetCommon(emma, storm)

	json, err := json.Marshal(common)
	if err != nil {
		http.Error(w, "Error marshaling json", 500)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
