package v1

import (
	"fmt"
	"strconv"

	"github.com/viswas163/MarvelousShipt/db"
	"github.com/viswas163/MarvelousShipt/models"
)

var (
	allComicsAPI       = "comics"
	allComicsJSONDBKey = "comicsJSON"
)

func GetAllComicsOfChar(char models.Character) (models.ComicsDataWrapper, error) {

	comics, err := getComicsByCharIDOffset(char.ID, 0)
	if err != nil {
		return models.ComicsDataWrapper{}, err
	}

	for i := 100; i < comics.Data.Total; i += 100 {
		cs, err := getComicsByCharIDOffset(char.ID, i)
		if err != nil {
			return models.ComicsDataWrapper{}, err
		}
		comics.Data.Results = append(comics.Data.Results, cs.Data.Results...)
	}

	return comics, nil
}

func getComicsByCharIDOffset(id int, offset int) (models.ComicsDataWrapper, error) {
	coms := []byte{}
	db.GetInstance().Get(fmt.Sprint(allComicsJSONDBKey, offset), &coms)
	var err error
	if len(coms) <= 0 {
		fmt.Println("Running all comics API")
		coms, err = RunAPIWithParam(allCharactersAPI, offset, []string{strconv.Itoa(id), allComicsAPI})
		if err != nil {
			return models.ComicsDataWrapper{}, err
		}
		db.GetInstance().Put(fmt.Sprint(allComicsJSONDBKey, offset), coms)
	} else {
		fmt.Println("Getting comics response from cache!")
	}
	comics, err := models.GetAllComics(coms)
	if err != nil {
		return models.ComicsDataWrapper{}, err
	}
	return comics, nil
}
