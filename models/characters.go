package models

import (
	"encoding/json"
	"fmt"
	"sync"
)

var (
	// CharactersByName : map[character_name][characterModel]  { Map[string][Character] }
	CharactersByName sync.Map
	characterCount   = 0
)

// AllCharacters : Model for all characters
type AllCharacters []Character

// Character : Model for character
type Character struct {
	ID          int       `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	ResourceURI string    `json:"resourceURI,omitempty"`
	Comics      ComicList `json:"comics,omitempty"`
}

// CharacterParams are optional parameters to narrow the character results returned
type CharacterParams struct {
	Name    string `url:"name,omitempty"`
	OrderBy string `url:"orderBy,omitempty"`
	Limit   int    `url:"limit,omitempty"`
	Offset  int    `url:"offset,omitempty"`
}

// CharacterDataWrapper provides character wrapper information returned by the API.
type CharacterDataWrapper struct {
	DataWrapper
	Data CharacterDataContainer `json:"data,omitempty"`
}

// CharacterDataContainer provides character container information returned by the API.
type CharacterDataContainer struct {
	DataContainer
	Results []Character `json:"results,omitempty"`
}

// CharacterList provides characters related to the parent entity.
type CharacterList struct {
	List
	Items []CharacterSummary `json:"items,omitempty"`
}

// CharacterSummary provides the summary for a character related to the parent entity.
type CharacterSummary struct {
	Summary
}

// GetAllCharacters : Gets all characters from response
func GetAllCharacters(charactersJSON []byte) (CharacterDataWrapper, error) {
	var allCharWrapper CharacterDataWrapper
	err := json.Unmarshal(charactersJSON, &allCharWrapper)
	if err != nil {
		fmt.Println("Error unmarshaling all chars")
		return CharacterDataWrapper{}, err
	}
	return allCharWrapper, nil
}

func CheckCharsExist(allChars []Character) bool {
	firstChName := allChars[0].Name
	lastChName := allChars[len(allChars)-1].Name

	_, firstOK := CharactersByName.Load(firstChName)
	_, lastOK := CharactersByName.Load(lastChName)

	return firstOK && lastOK
}

// SetAllCharactersByName : Sets all characters to map by name if not already existing
func SetAllCharactersByName(allChars []Character) error {
	if CheckCharsExist(allChars) {
		return nil
	}
	for _, character := range allChars {
		CharactersByName.LoadOrStore(character.Name, character)
		characterCount++
	}
	fmt.Printf("Added %d new characters!", characterCount)
	return nil
}

// GetCharacterCount : Returns the charactersbyname count
func GetCharacterCount() int {
	return characterCount
}
