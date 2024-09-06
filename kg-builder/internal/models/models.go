package models

type Concept struct {
	Name      string `json:"name"`
	Relation  string `json:"relation"`
	RelatedTo string `json:"relatedTo"`
}
