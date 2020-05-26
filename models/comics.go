package models

// Comic represents a Marvel comic.
type Comic struct {
	ID                 int            `json:"id,omitempty"`
	DigitalID          int            `json:"digitalId,omitempty"`
	Title              string         `json:"title,omitempty"`
	IssueNumber        int            `json:"issueNumber,omitempty"`
	VariantDescription string         `json:"variantDescription,omitempty"`
	Description        string         `json:"description,omitempty"`
	Format             string         `json:"format,omitempty"`
	PageCount          int            `json:"pageCount,omitempty"`
	ResourceURI        string         `json:"resourceURI,omitempty"`
	Variants           []ComicSummary `json:"variants,omitempty"`
	Collections        []ComicSummary `json:"collections,omitempty"`
	CollectedIssues    []ComicSummary `json:"collectedIssues,omitempty"`
	Characters         CharacterList  `json:"characters,omitempty"`
}

// ComicList provides comics related to the parent entity.
type ComicList struct {
	List
	Items []ComicSummary `json:"items,omitempty"`
}

// ComicSummary provides the summary for a comic related to the parent entity.
type ComicSummary struct {
	Summary
}
