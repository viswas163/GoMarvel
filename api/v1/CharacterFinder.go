package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
	"github.com/juliangruber/go-intersect"
	"github.com/viswas163/MarvelousShipt/models"
)

func RetrieveCharacter(w http.ResponseWriter, r *http.Request) {
	key, ok := mux.Vars(r)["characterId"]
	if !ok {
		http.Error(w, "Invalid param", 400)
	}

	fmt.Print(key)
	character, _ := GetCharacter(key)

	json, err := json.Marshal(character)
	if err != nil {
		http.Error(w, "Error marshaling json", 500)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func GetCharacter(key string) (models.Character, error) {
	character := models.Character{}
	ch, ok := models.CharactersByName.Load(key)
	if ok {
		character, ok = ch.(models.Character)
		if !ok {
			return models.Character{}, errors.New("Error getting character ")
		}
	}
	return character, nil
}

func GetCommonComics(w http.ResponseWriter, r *http.Request) {
	key1, ok := mux.Vars(r)["characterId"]
	if !ok {
		http.Error(w, "Invalid param", 400)
	}
	key2, ok := mux.Vars(r)["characterId2"]
	if !ok {
		http.Error(w, "Invalid param", 400)
	}

	char1, _ := GetCharacter(key1)
	char2, _ := GetCharacter(key2)

	commonComics := GetCommon(char1, char2)

	json, err := json.Marshal(commonComics)
	if err != nil {
		http.Error(w, "Error marshaling json", 500)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func GetCommon(char1 models.Character, char2 models.Character) interface{} {
	PushComicsOfChar(char1)
	PushComicsOfChar(char2)
	l1, ok1 := models.ComicsByCharacter.Load(char1.Name)
	l2, ok2 := models.ComicsByCharacter.Load(char2.Name)
	if !ok1 || !ok2 {
		return nil
	}
	ids := intersect.Hash(l1, l2)
	s := reflect.ValueOf(ids)

	var comics []interface{}
	for i := 0; i < s.Len(); i++ {
		id := s.Index(0).Interface().(int)
		comic, ok := models.ComicsByID.Load(id)
		if !ok {
			continue
		}
		comics = append(comics, comic)
	}
	return comics
}

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

func EmmaStorm(w http.ResponseWriter, r *http.Request) {
	var emma models.Character
	var storm models.Character

	t1, _ := models.CharactersByName.Load("Emma Frost")
	emma = t1.(models.Character)
	t2, _ := models.CharactersByName.Load("Storm")
	storm = t2.(models.Character)

	common := GetCommon(emma, storm)

	json, err := json.Marshal(common)
	if err != nil {
		http.Error(w, "Error marshaling json", 500)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
