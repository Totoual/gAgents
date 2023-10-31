package chatgpt

type SearchResult struct {
	Object              string   `json:"object"`
	Characteristics     []string `json:"characteristics"`
	Category            string   `json:"category"`
	PriceRange          float32  `json:"price_range"`
	IntendedUse         string   `json:"intended_use"`
	MaterialPreferences []string `json:"material_preferences"`
	RelevantTopics      []string `json:"relevant_topics"`
}
