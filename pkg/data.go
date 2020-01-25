package exago

type Data struct {
	Results  Results           `json:"results"`
	Metadata Metadata          `json:"metadata"`
	Score    Score             `json:"score"`
	Errors   map[string]string `json:"errors,omitempty"`
}
