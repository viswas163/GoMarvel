package models

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
