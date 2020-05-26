package models

// DataWrapper provides the common wrapper attributes to unmarshal the API's response.
// It is used to compose more specific wrappers, e.g., CharacterDataWrapper.
type DataWrapper struct {
	Code            int    `json:"code,omitempty"`
	Status          string `json:"status,omitempty"`
	Copyright       string `json:"copyright,omitempty"`
	AttributionText string `json:"attributionText,omitempty"`
	AttributionHTML string `json:"attributionHTML,omitempty"`
	ETag            string `json:"etag,omitempty"`
}

// DataContainer provides the common container attributes to unmarshal the API's response.
// It is used to compose more specific containers, e.g., CharacterDataContainer.
type DataContainer struct {
	Offset int `json:"offset,omitempty"`
	Limit  int `json:"limit,omitempty"`
	Total  int `json:"total,omitempty"`
	Count  int `json:"count,omitempty"`
}

// Summary provides the common summary attributes to unmarshal the API's response.
// It is used to compose more specific containers, e.g., CharacterSummary.
type Summary struct {
	ResourceURI string `json:"resourceURI,omitempty"`
	Name        string `json:"name,omitempty"`
}

// List provides the common list attributes to unmarshal the API's response.
// It is used to compose more specific containers, e.g., CharacterList.
type List struct {
	Available     int    `json:"available,omitempty"`
	Returned      int    `json:"returned,omitempty"`
	CollectionURI string `json:"collectionURI,omitempty"`
}
